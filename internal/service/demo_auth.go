package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
)

type DemoUser struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type DemoAuth struct {
	mu       sync.RWMutex
	users    []DemoUser
	sessions map[string]string
}

func NewDemoAuth() *DemoAuth {
	return &DemoAuth{
		users: []DemoUser{
			{Username: "urpi", Name: "Urpi"},
			{Username: "tux", Name: "Tux"},
			{Username: "gopher", Name: "Gopher"},
		},
		sessions: map[string]string{},
	}
}

func (a *DemoAuth) Users() []DemoUser {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]DemoUser, len(a.users))
	copy(out, a.users)
	for i := range out {
		out[i].Password = ""
	}
	return out
}

func (a *DemoAuth) Login(username, password string) (string, DemoUser, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	a.mu.RLock()
	var user DemoUser
	for _, candidate := range a.users {
		if candidate.Username == username {
			user = candidate
			break
		}
	}
	a.mu.RUnlock()
	if user.Username == "" || password != username {
		return "", DemoUser{}, errors.New("invalid demo credentials")
	}
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", DemoUser{}, err
	}
	token := hex.EncodeToString(raw)
	a.mu.Lock()
	a.sessions[token] = user.Username
	a.mu.Unlock()
	user.Password = ""
	return token, user, nil
}

func (a *DemoAuth) Authenticated(token string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	username, ok := a.sessions[strings.TrimSpace(token)]
	return username, ok
}
