package repository

import (
	"context"
	"github.com/example/cms/internal/domain"
)

type IContentRepository interface {
	List(context.Context, int, int, string, string, []string) ([]domain.Content, error)
	ListCategories(context.Context) ([]domain.Category, error)
	ListBadges(context.Context) ([]domain.Badge, error)
	FindBySlug(context.Context, string) (domain.Content, error)
	Update(context.Context, domain.Content, domain.ContentHistory) error
}
