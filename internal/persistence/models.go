package persistence

import "time"

type UserModel struct {
	ID              uint   `gorm:"primaryKey"`
	Email           string `gorm:"uniqueIndex;not null"`
	Name            string
	Provider        string `gorm:"not null"`
	ProviderSubject string `gorm:"uniqueIndex;not null"`
	CreatedAt       time.Time
}

func (UserModel) TableName() string { return "users" }

type AdminModel struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Name      string
	CreatedAt time.Time
}

func (AdminModel) TableName() string { return "admins" }

type SessionModel struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    *uint
	AdminID   *uint
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (SessionModel) TableName() string { return "sessions" }
