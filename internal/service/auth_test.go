package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/identity"
)

type fakeUsers struct{}

func (fakeUsers) FindOrCreateByIdentity(_ context.Context, u domain.User) (domain.User, error) {
	u.ID = 7
	return u, nil
}

type fakeAdmins struct{}

func (fakeAdmins) FindByEmail(context.Context, string) (domain.Admin, error) {
	return domain.Admin{ID: 2, Email: "admin@example.com"}, nil
}

type fakeSessions struct{ created domain.Session }

func (f *fakeSessions) Create(_ context.Context, s domain.Session) error { f.created = s; return nil }
func (f *fakeSessions) FindValid(context.Context, string) (domain.Session, error) {
	return f.created, nil
}

func TestAuthServiceCreatesUserSession(t *testing.T) {
	sessions := &fakeSessions{}
	s := NewAuthService(fakeUsers{}, fakeAdmins{}, sessions, time.Hour)
	u, token, err := s.Login(context.Background(), identity.Identity{Subject: "sub", Email: "user@example.com"})
	if err != nil || token == "" || u.ID != 7 || sessions.created.UserID == nil {
		t.Fatalf("login failed: user=%+v token=%q session=%+v err=%v", u, token, sessions.created, err)
	}
}
