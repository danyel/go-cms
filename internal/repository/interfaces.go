package repository

import (
	"context"
	"github.com/example/cms/internal/domain"
)

type IUserRepository interface {
	FindOrCreateByIdentity(context.Context, domain.User) (domain.User, error)
}
type IAdminRepository interface {
	FindByEmail(context.Context, string) (domain.Admin, error)
}
type ISessionRepository interface {
	Create(context.Context, domain.Session) error
	FindValid(context.Context, string) (domain.Session, error)
}
