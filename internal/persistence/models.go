package persistence

import "time"

type ContentModel struct {
	ID                   uint   `gorm:"primaryKey"`
	Slug                 string `gorm:"uniqueIndex;not null"`
	Title                string `gorm:"not null"`
	Summary              string
	Body                 string `gorm:"not null"`
	Status               string `gorm:"not null;default:'draft'"`
	Published            bool   `gorm:"not null;default:false"`
	CategoryID           uint
	Category             CategoryModel `gorm:"foreignKey:CategoryID"`
	Badges               []BadgeModel  `gorm:"many2many:content_badges;joinForeignKey:ContentID;joinReferences:BadgeID"`
	CreatedAt, UpdatedAt time.Time
}

func (ContentModel) TableName() string { return "content" }

type CategoryModel struct {
	ID   uint   `gorm:"primaryKey"`
	Slug string `gorm:"uniqueIndex;not null"`
	Name string `gorm:"not null"`
}

func (CategoryModel) TableName() string { return "categories" }

type BadgeModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

func (BadgeModel) TableName() string { return "badges" }

type ContentHistoryModel struct {
	ID        uint   `gorm:"primaryKey"`
	ContentID uint   `gorm:"index;not null"`
	Operation string `gorm:"not null"`
	CreatedAt time.Time
	Snapshot  string `gorm:"type:jsonb;not null"`
}

func (ContentHistoryModel) TableName() string { return "content_history" }
