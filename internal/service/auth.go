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
	users         repository.IUserRepository
	admins        repository.IAdminRepository
	userSessions  repository.ISessionRepository
	adminSessions repository.ISessionRepository
	ttl           time.Duration
}

func NewAuthService(u repository.IUserRepository, a repository.IAdminRepository, sessions repository.ISessionRepository, ttl time.Duration) *AuthService {
	return NewAuthServiceWithSeparateSessions(u, a, sessions, sessions, ttl)
}

func NewAuthServiceWithSeparateSessions(u repository.IUserRepository, a repository.IAdminRepository, userSessions, adminSessions repository.ISessionRepository, ttl time.Duration) *AuthService {
	return &AuthService{users: u, admins: a, userSessions: userSessions, adminSessions: adminSessions, ttl: ttl}
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
	if e = s.userSessions.Create(ctx, domain.Session{Token: t, UserID: &u.ID, ExpiresAt: time.Now().Add(s.ttl), CreatedAt: time.Now()}); e != nil {
		return domain.User{}, "", e
	}
	return u, t, nil
}
func (s *AuthService) Authenticate(ctx context.Context, t string) (domain.Session, error) {
	return s.userSessions.FindValid(ctx, t)
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
	if e = s.adminSessions.Create(ctx, domain.Session{Token: t, AdminID: &a.ID, ExpiresAt: time.Now().Add(s.ttl), CreatedAt: time.Now()}); e != nil {
		return domain.Admin{}, "", e
	}
	return a, t, nil
}
func (s *AuthService) AdminAuthenticate(ctx context.Context, t string) (domain.Session, error) {
	return s.adminSessions.FindValid(ctx, t)
}
func (s *AuthService) UserByID(ctx context.Context, id uint) (domain.User, error) {
	if lookup, ok := s.users.(repository.IUserLookup); ok {
		return lookup.FindByID(ctx, id)
	}
	return domain.User{}, errors.New("user lookup unavailable")
}
func (s *AuthService) CanEdit(ctx context.Context, session domain.Session) bool {
	if session.AdminID != nil {
		return true
	}
	if session.UserID == nil {
		return false
	}
	u, err := s.UserByID(ctx, *session.UserID)
	return err == nil && (u.Editor || u.Role == "editor" || u.Role == "admin")
}
func newToken() (string, error) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
