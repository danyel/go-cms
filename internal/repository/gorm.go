package repository

import (
	"context"
	"errors"
	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/persistence"
	"gorm.io/gorm"
	"time"
)

type GORMUserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) IUserRepository { return &GORMUserRepository{db} }
func (r *GORMUserRepository) FindOrCreateByIdentity(ctx context.Context, d domain.User) (domain.User, error) {
	var m persistence.UserModel
	err := r.db.WithContext(ctx).Where("provider = ? AND provider_subject = ?", d.Provider, d.ProviderSubject).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		m = persistence.UserFromDomain(d)
		if err = r.db.WithContext(ctx).Create(&m).Error; err != nil {
			return domain.User{}, err
		}
		return persistence.UserToDomain(m), nil
	}
	if err != nil {
		return domain.User{}, err
	}
	return persistence.UserToDomain(m), nil
}
func (r *GORMUserRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var m persistence.UserModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return domain.User{}, err
	}
	return persistence.UserToDomain(m), nil
}

type GORMAdminRepository struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) IAdminRepository { return &GORMAdminRepository{db} }
func (r *GORMAdminRepository) FindByEmail(ctx context.Context, email string) (domain.Admin, error) {
	var m persistence.AdminModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return domain.Admin{}, err
	}
	return persistence.AdminToDomain(m), nil
}

type GORMSessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) ISessionRepository { return &GORMSessionRepository{db} }
func (r *GORMSessionRepository) Create(ctx context.Context, d domain.Session) error {
	return r.db.WithContext(ctx).Create(&persistence.SessionModel{ID: d.ID, Token: d.Token, UserID: d.UserID, AdminID: d.AdminID, ExpiresAt: d.ExpiresAt, CreatedAt: d.CreatedAt}).Error
}
func (r *GORMSessionRepository) FindValid(ctx context.Context, t string) (domain.Session, error) {
	var m persistence.SessionModel
	if err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", t, time.Now()).First(&m).Error; err != nil {
		return domain.Session{}, err
	}
	return persistence.SessionToDomain(m), nil
}

type GORMContentRepository struct{ db *gorm.DB }

func NewContentRepository(db *gorm.DB) IContentRepository { return &GORMContentRepository{db} }
func (r *GORMContentRepository) List(ctx context.Context, limit, offset int) ([]domain.Content, error) {
	var rows []persistence.ContentModel
	err := r.db.WithContext(ctx).Order("updated_at DESC, created_at DESC, id DESC").Limit(limit).Offset(offset).Find(&rows).Error
	out := make([]domain.Content, len(rows))
	for i := range rows {
		out[i] = persistence.ContentToDomain(rows[i])
	}
	return out, err
}
func (r *GORMContentRepository) FindBySlug(ctx context.Context, slug string) (domain.Content, error) {
	var row persistence.ContentModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&row).Error; err != nil {
		return domain.Content{}, err
	}
	return persistence.ContentToDomain(row), nil
}
func (r *GORMContentRepository) Update(ctx context.Context, d domain.Content, h domain.ContentHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&persistence.ContentModel{}).Where("id = ?", d.ID).Updates(map[string]any{"slug": d.Slug, "title": d.Title, "summary": d.Summary, "body": d.Body, "status": d.Status, "published": d.Published, "updated_at": d.UpdatedAt, "updated_by": d.UpdatedBy}).Error; err != nil {
			return err
		}
		return tx.Create(&persistence.ContentHistoryModel{ContentID: h.ContentID, Operation: h.Operation, ActorID: h.ActorID, ActorAdminID: h.ActorAdminID, CreatedAt: h.CreatedAt, Snapshot: h.Snapshot}).Error
	})
}
