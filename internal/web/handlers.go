package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/example/cms/internal/security"
	"github.com/example/cms/internal/service"
)

// Handler serves the CMS API and the embedded React app from one binary.
// Reading content is anonymous; every mutation requires the identity asserted by
// the upstream SSO proxy.
type Handler struct {
	content service.IContentService
	assets  fs.FS
	sso     security.ITokenVerifier
	demo    *service.DemoAuth
}

func NewHandler(content service.IContentService, assets fs.FS, sso security.ITokenVerifier, demo ...*service.DemoAuth) *Handler {
	h := &Handler{content: content, assets: assets, sso: sso}
	if len(demo) > 0 {
		h.demo = demo[0]
	}
	return h
}

func (h *Handler) Routes(origins string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		jsonWrite(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/auth/session", h.session)
	mux.HandleFunc("/api/auth/sudo", h.sudo)
	mux.HandleFunc("/api/auth/users", h.users)
	mux.HandleFunc("/api/protected", h.protected)
	mux.HandleFunc("/api/test/protected", h.protected)
	mux.HandleFunc("/api/content", h.contentList)
	mux.HandleFunc("/api/content/categories", h.contentCategories)
	mux.HandleFunc("/api/content/badges", h.contentBadges)
	mux.HandleFunc("/api/content/", h.contentDetail)
	mux.HandleFunc("/", h.singlePageApp)
	return cors(mux, origins, h.sso.HeaderName())
}

func (h *Handler) session(w http.ResponseWriter, r *http.Request) {
	username, authenticated := h.identity(r)
	jsonWrite(w, http.StatusOK, map[string]any{"authenticated": authenticated, "canEdit": authenticated, "username": username})
}

func (h *Handler) protected(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.identity(r); !ok {
		jsonWrite(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]string{"status": "authenticated"})
}

func (h *Handler) users(w http.ResponseWriter, r *http.Request) {
	if h.demo == nil || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"items": h.demo.Users()})
}

func (h *Handler) sudo(w http.ResponseWriter, r *http.Request) {
	if h.demo == nil || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !decode(r, &input) {
		jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	token, user, err := h.demo.Login(input.Username, input.Password)
	if err != nil {
		jsonWrite(w, http.StatusUnauthorized, map[string]string{"error": "invalid demo credentials"})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *Handler) contentList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	limit, offset := 20, 0
	_, _ = fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &limit)
	_, _ = fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)
	badges := r.URL.Query()["badge"]
	if len(badges) == 0 {
		badges = strings.Split(r.URL.Query().Get("badges"), ",")
	}
	items, err := h.content.List(r.Context(), limit, offset,
		r.URL.Query().Get("category"),
		r.URL.Query().Get("status"),
		badges,
	)
	if err != nil {
		jsonWrite(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"items": items, "limit": limit, "offset": offset})
}

func (h *Handler) contentCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	categories, err := h.content.ListCategories(r.Context())
	if err != nil {
		jsonWrite(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"items": categories})
}

func (h *Handler) contentBadges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
		return
	}
	badges, err := h.content.ListBadges(r.Context())
	if err != nil {
		jsonWrite(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	jsonWrite(w, http.StatusOK, map[string]any{"items": badges})
}

func (h *Handler) contentDetail(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/content/")
	switch r.Method {
	case http.MethodGet:
		item, err := h.content.Get(r.Context(), slug)
		if err != nil {
			jsonWrite(w, http.StatusNotFound, map[string]string{"error": "content not found"})
			return
		}
		jsonWrite(w, http.StatusOK, item)
	case http.MethodPut:
		username, authenticated := h.identity(r)
		if !authenticated {
			jsonWrite(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
			return
		}
		var item service.Content
		if !decode(r, &item) {
			jsonWrite(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		item.Slug = slug
		var saved service.Content
		var err error
		if owned, ok := h.content.(interface {
			UpdateAs(context.Context, service.Content, string) (service.Content, error)
		}); ok {
			saved, err = owned.UpdateAs(r.Context(), item, username)
		} else {
			saved, err = h.content.Update(r.Context(), item)
		}
		if err != nil {
			jsonWrite(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		jsonWrite(w, http.StatusOK, saved)
	default:
		jsonWrite(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) identity(r *http.Request) (string, bool) {
	if h.demo != nil {
		token := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if username, ok := h.demo.Authenticated(token); ok {
			return username, true
		}
	}
	if h.sso.Authenticated(r) {
		return "sso-owner", true
	}
	return "", false
}

// singlePageApp serves the embedded build and falls back to index.html so the
// client-side router owns every non-API path.
func (h *Handler) singlePageApp(w http.ResponseWriter, r *http.Request) {
	if h.assets == nil {
		http.NotFound(w, r)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name != "" {
		if file, err := h.assets.Open(name); err == nil {
			_ = file.Close()
			http.FileServer(http.FS(h.assets)).ServeHTTP(w, r)
			return
		}
	}
	index, err := fs.ReadFile(h.assets, "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(index)
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
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
