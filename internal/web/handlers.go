package web

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/identity"
	"github.com/example/cms/internal/service"
	"github.com/example/cms/internal/web/mapping"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	auth    *service.AuthService
	google  identity.GoogleProvider
	content service.IContentService
}

func NewHandler(a *service.AuthService, google identity.GoogleProvider) *Handler {
	return &Handler{auth: a, google: google}
}
func NewHandlerWithContent(a *service.AuthService, google identity.GoogleProvider, content service.IContentService) *Handler {
	return &Handler{auth: a, google: google, content: content}
}
func (h *Handler) Routes(origins string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { jsonWrite(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/api/login", h.login)
	mux.HandleFunc("/api/auth/google", h.googleLogin)
	mux.HandleFunc("/api/auth/google/callback", h.googleCallback)
	mux.HandleFunc("/api/auth/session", h.session)
	mux.HandleFunc("/admin/login", h.adminLogin)
	mux.HandleFunc("/api/test/protected", h.protected)
	mux.HandleFunc("/api/protected", h.protected)
	mux.HandleFunc("/api/content", h.contentList)
	mux.HandleFunc("/api/content/categories", h.contentCategories)
	mux.HandleFunc("/api/content/", h.contentDetail)
	return cors(mux, origins)
}
func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	if h.google.ClientID == "" {
		jsonWrite(w, http.StatusServiceUnavailable, map[string]string{"error": "google identity provider is not configured"})
		return
	}
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		jsonWrite(w, http.StatusInternalServerError, map[string]string{"error": "unable to create oauth state"})
		return
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes)
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, HttpOnly: true, Secure: r.TLS != nil, SameSite: http.SameSiteLaxMode, Path: "/", MaxAge: 600})
	http.Redirect(w, r, h.google.OAuthConfig().AuthCodeURL(state), http.StatusFound)
}
func (h *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	id, err := h.google.FetchIdentity(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("Google OAuth callback failed: %v", err)
		http.Error(w, "google authentication failed; start a new sign-in attempt", http.StatusUnauthorized)
		return
	}
	user, token, err := h.auth.Login(r.Context(), id)
	if err != nil {
		http.Error(w, "unable to create session", http.StatusInternalServerError)
		return
	}
	setSessionCookie(w, token, r.TLS != nil)
	_ = user
	http.Redirect(w, r, "/test", http.StatusFound)
}
func (h *Handler) session(w http.ResponseWriter, r *http.Request) {
	token := bearerOrCookie(r)
	if token == "" {
		jsonWrite(w, http.StatusUnauthorized, map[string]string{"authenticated": "false"})
		return
	}
	s, err := h.auth.Authenticate(r.Context(), token)
	if err != nil || s.UserID == nil {
		jsonWrite(w, http.StatusUnauthorized, map[string]string{"authenticated": "false"})
		return
	}
	u, err := h.auth.UserByID(r.Context(), *s.UserID)
	if err != nil {
		jsonWrite(w, http.StatusUnauthorized, map[string]string{"authenticated": "false"})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"authenticated": true, "user": u, "canEdit": h.auth.CanEdit(r.Context(), s)})
}
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonWrite(w, 405, nil)
		return
	}
	var req mapping.LoginRequest
	if !decode(r, &req) {
		jsonWrite(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	u, t, e := h.auth.Login(r.Context(), req.Identity())
	if e != nil {
		jsonWrite(w, 400, map[string]string{"error": e.Error()})
		return
	}
	jsonWrite(w, 200, mapping.LoginResponse{Token: t, User: u})
}
func (h *Handler) adminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		jsonWrite(w, 405, nil)
		return
	}
	var req struct {
		Email string `json:"email"`
	}
	if !decode(r, &req) {
		jsonWrite(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	a, t, e := h.auth.AdminLogin(r.Context(), req.Email)
	if e != nil {
		jsonWrite(w, 401, map[string]string{"error": e.Error()})
		return
	}
	jsonWrite(w, 200, mapping.LoginResponse{Token: t, Admin: a})
}
func (h *Handler) protected(w http.ResponseWriter, r *http.Request) {
	t := bearerOrCookie(r)
	if t == "" {
		jsonWrite(w, 401, map[string]string{"error": "missing bearer token"})
		return
	}
	s, e := h.auth.Authenticate(r.Context(), t)
	if e != nil || s.UserID == nil {
		jsonWrite(w, 401, map[string]string{"error": "invalid session"})
		return
	}
	jsonWrite(w, 200, map[string]string{"status": "authenticated"})
}
func (h *Handler) sessionForRequest(r *http.Request) (domain.Session, bool) {
	t := bearerOrCookie(r)
	if t == "" {
		return domain.Session{}, false
	}
	if s, e := h.auth.Authenticate(r.Context(), t); e == nil {
		return s, true
	}
	if s, e := h.auth.AdminAuthenticate(r.Context(), t); e == nil {
		return s, true
	}
	return domain.Session{}, false
}
func (h *Handler) contentList(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	if _, ok := h.sessionForRequest(r); !ok {
		jsonWrite(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	if r.Method != "GET" {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	limit, offset := 20, 0
	_, _ = fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
	_, _ = fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)
	badges := r.URL.Query()["badge"]
	if len(badges) == 0 {
		badges = strings.Split(strings.TrimSpace(r.URL.Query().Get("badges")), ",")
	}
	filtered := badges[:0]
	for _, badge := range badges {
		if badge = strings.TrimSpace(strings.ToLower(badge)); badge != "" {
			filtered = append(filtered, badge)
		}
	}
	items, e := h.content.List(r.Context(), limit, offset, strings.TrimSpace(strings.ToLower(r.URL.Query().Get("category"))), filtered)
	if e != nil {
		jsonWrite(w, 500, map[string]string{"error": e.Error()})
		return
	}
	jsonWrite(w, 200, map[string]any{"items": items, "limit": limit, "offset": offset})
}
func (h *Handler) contentCategories(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	if _, ok := h.sessionForRequest(r); !ok {
		jsonWrite(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	categories, err := h.content.ListCategories(r.Context())
	if err != nil {
		jsonWrite(w, 500, map[string]string{"error": err.Error()})
		return
	}
	jsonWrite(w, 200, map[string]any{"items": categories})
}
func (h *Handler) contentDetail(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/api/content/")
	if r.Method == "GET" {
		if _, ok := h.sessionForRequest(r); !ok {
			jsonWrite(w, 401, map[string]string{"error": "authentication required"})
			return
		}
		c, e := h.content.Get(r.Context(), slug)
		if e != nil {
			jsonWrite(w, 404, map[string]string{"error": "content not found"})
			return
		}
		jsonWrite(w, 200, c)
		return
	}
	if r.Method != "PUT" {
		jsonWrite(w, 405, nil)
		return
	}
	s, ok := h.sessionForRequest(r)
	if !ok {
		jsonWrite(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	if !h.auth.CanEdit(r.Context(), s) {
		jsonWrite(w, 403, map[string]string{"error": "forbidden"})
		return
	}
	var c domain.Content
	if !decode(r, &c) {
		jsonWrite(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	c.Slug = slug
	out, e := h.content.Update(r.Context(), c, s)
	if e != nil {
		jsonWrite(w, 400, map[string]string{"error": e.Error()})
		return
	}
	jsonWrite(w, 200, out)
}
func bearerOrCookie(r *http.Request) string {
	header := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	if header != "" {
		return header
	}
	cookie, err := r.Cookie("cms_session")
	if err != nil {
		return ""
	}
	return cookie.Value
}
func setSessionCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{Name: "cms_session", Value: token, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, Path: "/", MaxAge: 86400})
}
func decode(r *http.Request, v any) bool { return json.NewDecoder(r.Body).Decode(v) == nil }
func jsonWrite(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
func cors(next http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
