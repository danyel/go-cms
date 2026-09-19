package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/security"
)

type fakeContent struct{ updated bool }

func (f *fakeContent) List(context.Context, int, int, string, string, []string) ([]domain.Content, error) {
	return nil, nil
}
func (f *fakeContent) ListCategories(context.Context) ([]domain.Category, error) { return nil, nil }
func (f *fakeContent) ListBadges(context.Context) ([]domain.Badge, error)        { return nil, nil }
func (f *fakeContent) Get(context.Context, string) (domain.Content, error) {
	return domain.Content{Slug: "post", Title: "Post", Body: "body"}, nil
}
func (f *fakeContent) Update(_ context.Context, c domain.Content) (domain.Content, error) {
	f.updated = true
	return c, nil
}

func newTestRoutes(content *fakeContent) http.Handler {
	return NewHandler(content, security.NewSSOTokenVerifier("X-SSO-Token", "", "")).Routes("http://localhost:5173")
}

func TestAnonymousCanReadButNotEdit(t *testing.T) {
	content := &fakeContent{}
	routes := newTestRoutes(content)

	get := httptest.NewRecorder()
	routes.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/content/post", nil))
	if get.Code != http.StatusOK {
		t.Fatalf("anonymous read status = %d, want %d", get.Code, http.StatusOK)
	}

	put := httptest.NewRecorder()
	routes.ServeHTTP(put, httptest.NewRequest(http.MethodPut, "/api/content/post", strings.NewReader(`{"title":"T","body":"B"}`)))
	if put.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous update status = %d, want %d", put.Code, http.StatusUnauthorized)
	}
	if content.updated {
		t.Fatal("anonymous request must not reach the content service")
	}

	session := httptest.NewRecorder()
	routes.ServeHTTP(session, httptest.NewRequest(http.MethodGet, "/api/auth/session", nil))
	if !strings.Contains(session.Body.String(), `"authenticated":false`) {
		t.Fatalf("anonymous session = %s", session.Body.String())
	}
}

func TestProxyIdentityCanEdit(t *testing.T) {
	content := &fakeContent{}
	routes := newTestRoutes(content)

	put := httptest.NewRequest(http.MethodPut, "/api/content/post", strings.NewReader(`{"title":"T","body":"B"}`))
	put.Header.Set("X-SSO-Token", "owner-token")
	putRec := httptest.NewRecorder()
	routes.ServeHTTP(putRec, put)
	if putRec.Code != http.StatusOK {
		t.Fatalf("owner update status = %d, body = %s", putRec.Code, putRec.Body.String())
	}
	if !content.updated {
		t.Fatal("owner request must reach the content service")
	}

	session := httptest.NewRequest(http.MethodGet, "/api/auth/session", nil)
	session.Header.Set("X-SSO-Token", "owner-token")
	sessionRec := httptest.NewRecorder()
	routes.ServeHTTP(sessionRec, session)
	if !strings.Contains(sessionRec.Body.String(), `"canEdit":true`) {
		t.Fatalf("owner session = %s", sessionRec.Body.String())
	}
}
