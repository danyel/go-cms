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

type GORMAdminRepository struct{ db *gorm.DB }

func NewAdminRepository(db *gorm.DB) IAdminRepository { return &GORMAdminRepository{db} }
func (r *GORMAdminRepository) FindByEmail(ctx context.Context, email string) (domain.Admin, error) {
	var m persistence.AdminModel
	e := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if e != nil {
		return domain.Admin{}, e
	}
	return persistence.AdminToDomain(m), nil
}

type GORMSessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) ISessionRepository { return &GORMSessionRepository{db} }
func (r *GORMSessionRepository) Create(ctx context.Context, d domain.Session) error {
	return r.db.WithContext(ctx).Create(persistence.SessionFromDomain(d)).Error
}
func (r *GORMSessionRepository) FindValid(ctx context.Context, t string) (domain.Session, error) {
	var m persistence.SessionModel
	e := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", t, time.Now()).First(&m).Error
	if e != nil {
		return domain.Session{}, e
	}
	return persistence.SessionToDomain(m), nil
}
