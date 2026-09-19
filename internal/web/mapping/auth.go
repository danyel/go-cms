package mapping

import "github.com/example/cms/internal/identity"

type LoginRequest struct{ Subject, Email, Name string }

func (r LoginRequest) Identity() identity.Identity {
	return identity.Identity{Subject: r.Subject, Email: r.Email, Name: r.Name}
}

type LoginResponse struct {
	Token string `json:"token"`
	User  any    `json:"user,omitempty"`
	Admin any    `json:"admin,omitempty"`
}
