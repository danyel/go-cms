package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/identity"
	"github.com/example/cms/internal/repository"
	"strings"
	"time"
)

type IUserService interface {
	Login(context.Context, identity.Identity) (domain.User, string, error)
	Authenticate(context.Context, string) (domain.Session, error)
}
type IAdminAuthService interface {
	Login(context.Context, string) (domain.Admin, string, error)
	Authenticate(context.Context, string) (domain.Session, error)
}
type AuthService struct {
	users    repository.IUserRepository
	admins   repository.IAdminRepository
	sessions repository.ISessionRepository
	ttl      time.Duration
}

func NewAuthService(u repository.IUserRepository, a repository.IAdminRepository, s repository.ISessionRepository, ttl time.Duration) *AuthService {
	return &AuthService{u, a, s, ttl}
}
func (s *AuthService) Login(ctx context.Context, id identity.Identity) (domain.User, string, error) {
	if strings.TrimSpace(id.Email) == "" || strings.TrimSpace(id.Subject) == "" {
		return domain.User{}, "", errors.New("identity email and subject are required")
	}
	u, e := s.users.FindOrCreateByIdentity(ctx, domain.User{Email: id.Email, Name: id.Name, Provider: "google", ProviderSubject: id.Subject, CreatedAt: time.Now()})
	if e != nil {
		return domain.User{}, "", e
	}
	t, e := newToken()
	if e != nil {
		return domain.User{}, "", e
	}
	if e = s.sessions.Create(ctx, domain.Session{Token: t, UserID: &u.ID, ExpiresAt: time.Now().Add(s.ttl), CreatedAt: time.Now()}); e != nil {
		return domain.User{}, "", e
	}
	return u, t, nil
}
func (s *AuthService) Authenticate(ctx context.Context, t string) (domain.Session, error) {
	return s.sessions.FindValid(ctx, t)
}
func (s *AuthService) AdminLogin(ctx context.Context, email string) (domain.Admin, string, error) {
	a, e := s.admins.FindByEmail(ctx, email)
	if e != nil {
		return domain.Admin{}, "", errors.New("invalid admin credentials")
	}
	t, e := newToken()
	if e != nil {
		return domain.Admin{}, "", e
	}
	if e = s.sessions.Create(ctx, domain.Session{Token: t, AdminID: &a.ID, ExpiresAt: time.Now().Add(s.ttl), CreatedAt: time.Now()}); e != nil {
		return domain.Admin{}, "", e
	}
	return a, t, nil
}
func (s *AuthService) AdminAuthenticate(ctx context.Context, t string) (domain.Session, error) {
	return s.Authenticate(ctx, t)
}
func newToken() (string, error) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
