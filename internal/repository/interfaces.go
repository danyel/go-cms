package repository

import (
	"context"
	"github.com/example/cms/internal/domain"
)

type IUserRepository interface {
	FindOrCreateByIdentity(context.Context, domain.User) (domain.User, error)
}
type IUserLookup interface {
	FindByID(context.Context, uint) (domain.User, error)
}
type IAdminRepository interface {
	FindByEmail(context.Context, string) (domain.Admin, error)
}
type ISessionRepository interface {
	Create(context.Context, domain.Session) error
	FindValid(context.Context, string) (domain.Session, error)
}
type IContentRepository interface {
	List(context.Context, int, int, string, []string) ([]domain.Content, error)
	ListCategories(context.Context) ([]domain.Category, error)
	ListBadges(context.Context) ([]domain.Badge, error)
	FindBySlug(context.Context, string) (domain.Content, error)
	Update(context.Context, domain.Content, domain.ContentHistory) error
}
