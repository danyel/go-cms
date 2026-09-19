package persistence

import "time"

type UserModel struct {
	ID              uint   `gorm:"primaryKey"`
	Email           string `gorm:"uniqueIndex;not null"`
	Name            string
	Provider        string `gorm:"not null"`
	ProviderSubject string `gorm:"uniqueIndex;not null"`
	Role            string `gorm:"not null;default:''"`
	Editor          bool   `gorm:"not null;default:false"`
	CreatedAt       time.Time
}

type ContentModel struct {
	ID                   uint   `gorm:"primaryKey"`
	Slug                 string `gorm:"uniqueIndex;not null"`
	Title                string `gorm:"not null"`
	Summary              string
	Body                 string `gorm:"not null"`
	Status               string `gorm:"not null;default:'draft'"`
	Published            bool   `gorm:"not null;default:false"`
	CreatedAt, UpdatedAt time.Time
	CreatedBy, UpdatedBy uint
}

func (ContentModel) TableName() string { return "content" }

type ContentHistoryModel struct {
	ID           uint   `gorm:"primaryKey"`
	ContentID    uint   `gorm:"index;not null"`
	Operation    string `gorm:"not null"`
	ActorID      *uint
	ActorAdminID *uint
	CreatedAt    time.Time
	Snapshot     string `gorm:"type:jsonb;not null"`
}

func (ContentHistoryModel) TableName() string { return "content_history" }

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
