package domain

import "time"

type User struct {
	ID                        uint
	Email, Name               string
	Provider, ProviderSubject string
	Role                      string
	Editor                    bool
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

type Content struct {
	ID                         uint
	Slug, Title, Summary, Body string
	Status                     string
	Published                  bool
	CategoryID                 uint
	Category                   string
	Badges                     []string
	CreatedAt, UpdatedAt       time.Time
	CreatedBy, UpdatedBy       uint
}

type Category struct {
	ID   uint
	Slug string
	Name string
}

type Badge struct {
	ID   uint
	Name string
}

type ContentHistory struct {
	ID           uint
	ContentID    uint
	Operation    string
	ActorID      *uint
	ActorAdminID *uint
	CreatedAt    time.Time
	Snapshot     string
}
