package domain

import "time"

type Content struct {
	ID                         uint
	Slug, Title, Summary, Body string
	Status                     string
	Published                  bool
	CategoryID                 uint
	Category                   string
	Badges                     []string
	CreatedAt, UpdatedAt       time.Time
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
	ID        uint
	ContentID uint
	Operation string
	CreatedAt time.Time
	Snapshot  string
}
