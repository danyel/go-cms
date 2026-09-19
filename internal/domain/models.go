package domain

import "time"

type User struct {
	ID                        uint
	Email, Name               string
	Provider, ProviderSubject string
	CreatedAt                 time.Time
}
type Admin struct {
	ID          uint
	Email, Name string
	CreatedAt   time.Time
}
type Session struct {
	ID                   uint
	Token                string
	UserID               *uint
	AdminID              *uint
	ExpiresAt, CreatedAt time.Time
}
