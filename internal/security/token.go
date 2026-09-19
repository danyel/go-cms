// Package security decides whether a request carries the identity asserted by
// the upstream single-sign-on proxy. The CMS never signs users in itself: the
// proxy authenticates the owner and forwards a header on every accepted
// request. Requests without that header are anonymous.
package security

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// ITokenVerifier reports whether a request was authenticated by the SSO proxy.
type ITokenVerifier interface {
	Authenticated(*http.Request) bool
	HeaderName() string
}

// SSOTokenVerifier trusts the header injected by the identity-aware proxy. An
// empty Header disables header lookups; Cookie optionally names a cookie that
// carries the same value. When Token is non-empty the forwarded value must
// match it exactly, otherwise the presence of any non-empty value is accepted.
type SSOTokenVerifier struct {
	Header string
	Token  string
	Cookie string
}

// NewSSOTokenVerifier builds a verifier for the given header, optional expected
// token, and optional cookie fallback.
func NewSSOTokenVerifier(header, token, cookie string) SSOTokenVerifier {
	return SSOTokenVerifier{Header: header, Token: token, Cookie: cookie}
}

// HeaderName returns the proxy header carrying the SSO token.
func (v SSOTokenVerifier) HeaderName() string { return v.Header }

// Authenticated reports whether the request carries an accepted SSO token.
func (v SSOTokenVerifier) Authenticated(r *http.Request) bool {
	return v.accepted(v.presented(r))
}

func (v SSOTokenVerifier) presented(r *http.Request) string {
	if v.Header != "" {
		if value := trimBearer(r.Header.Get(v.Header)); value != "" {
			return value
		}
	}
	if v.Cookie != "" {
		if c, err := r.Cookie(v.Cookie); err == nil {
			return trimBearer(c.Value)
		}
	}
	return ""
}

func (v SSOTokenVerifier) accepted(presented string) bool {
	if presented == "" {
		return false
	}
	if v.Token == "" {
		return true
	}
	return subtle.ConstantTimeCompare([]byte(presented), []byte(v.Token)) == 1
}

func trimBearer(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= len("bearer ") && strings.EqualFold(value[:len("bearer ")], "bearer ") {
		return strings.TrimSpace(value[len("bearer "):])
	}
	return value
}
