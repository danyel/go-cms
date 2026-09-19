package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/security"
	"github.com/example/cms/internal/service"
)

// Handler serves the CMS API. Reading content is anonymous; every mutation
// requires the identity asserted by the upstream SSO proxy.
type Handler struct {
	content service.IContentService
	sso     security.ITokenVerifier
}

func NewHandler(content service.IContentService, sso security.ITokenVerifier) *Handler {
	return &Handler{content: content, sso: sso}
}
func (h *Handler) Routes(origins string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { jsonWrite(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/api/auth/session", h.session)
	mux.HandleFunc("/api/protected", h.protected)
	mux.HandleFunc("/api/test/protected", h.protected)
	mux.HandleFunc("/api/content", h.contentList)
	mux.HandleFunc("/api/content/categories", h.contentCategories)
	mux.HandleFunc("/api/content/badges", h.contentBadges)
	mux.HandleFunc("/api/content/", h.contentDetail)
	return cors(mux, origins, h.sso.HeaderName())
}
func (h *Handler) session(w http.ResponseWriter, r *http.Request) {
	authenticated := h.sso.Authenticated(r)
	jsonWrite(w, http.StatusOK, map[string]any{"authenticated": authenticated, "canEdit": authenticated})
}
func (h *Handler) protected(w http.ResponseWriter, r *http.Request) {
	if !h.sso.Authenticated(r) {
		jsonWrite(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	jsonWrite(w, 200, map[string]string{"status": "authenticated"})
}
func (h *Handler) contentList(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	if r.Method != http.MethodGet {
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
	items, e := h.content.List(r.Context(), limit, offset,
		strings.TrimSpace(strings.ToLower(r.URL.Query().Get("category"))),
		strings.TrimSpace(strings.ToLower(r.URL.Query().Get("status"))),
		filtered,
	)
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
func (h *Handler) contentBadges(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	badges, err := h.content.ListBadges(r.Context())
	if err != nil {
		jsonWrite(w, 500, map[string]string{"error": err.Error()})
		return
	}
	jsonWrite(w, 200, map[string]any{"items": badges})
}
func (h *Handler) contentDetail(w http.ResponseWriter, r *http.Request) {
	if h.content == nil {
		jsonWrite(w, 503, map[string]string{"error": "content unavailable"})
		return
	}
	slug := strings.TrimPrefix(r.URL.Path, "/api/content/")
	if r.Method == http.MethodGet {
		c, e := h.content.Get(r.Context(), slug)
		if e != nil {
			jsonWrite(w, 404, map[string]string{"error": "content not found"})
			return
		}
		jsonWrite(w, 200, c)
		return
	}
	if r.Method != http.MethodPut {
		jsonWrite(w, 405, nil)
		return
	}
	if !h.sso.Authenticated(r) {
		jsonWrite(w, 401, map[string]string{"error": "authentication required"})
		return
	}
	var c domain.Content
	if !decode(r, &c) {
		jsonWrite(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}
	c.Slug = slug
	out, e := h.content.Update(r.Context(), c)
	if e != nil {
		jsonWrite(w, 400, map[string]string{"error": e.Error()})
		return
	}
	jsonWrite(w, 200, out)
}
func decode(r *http.Request, v any) bool { return json.NewDecoder(r.Body).Decode(v) == nil }
func jsonWrite(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
func cors(next http.Handler, origin, ssoHeader string) http.Handler {
	allowed := "Content-Type, Authorization"
	if ssoHeader != "" && !strings.EqualFold(ssoHeader, "Authorization") {
		allowed += ", " + ssoHeader
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Headers", allowed)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
