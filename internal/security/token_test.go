package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func requestWithHeader(name, value string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(name, value)
	return r
}

func TestSSOTokenVerifierRequiresProxyHeader(t *testing.T) {
	v := NewSSOTokenVerifier("X-SSO-Token", "", "")
	if v.Authenticated(httptest.NewRequest(http.MethodGet, "/", nil)) {
		t.Fatal("request without proxy header must stay anonymous")
	}
	if !v.Authenticated(requestWithHeader("X-SSO-Token", "oauth2-access-token")) {
		t.Fatal("request with proxy header must be authenticated")
	}
}

func TestSSOTokenVerifierChecksExpectedToken(t *testing.T) {
	v := NewSSOTokenVerifier("X-SSO-Token", "secret", "")
	if v.Authenticated(requestWithHeader("X-SSO-Token", "other")) {
		t.Fatal("mismatching token must stay anonymous")
	}
	if !v.Authenticated(requestWithHeader("X-SSO-Token", "Bearer secret")) {
		t.Fatal("bearer-prefixed token must be authenticated")
	}
}

func TestSSOTokenVerifierCookieFallback(t *testing.T) {
	v := NewSSOTokenVerifier("X-SSO-Token", "", "cms_session")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "cms_session", Value: "cookie-token"})
	if !v.Authenticated(r) {
		t.Fatal("cookie token must be authenticated")
	}
}
