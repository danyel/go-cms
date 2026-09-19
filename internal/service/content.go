package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/repository"
	"regexp"
	"strings"
	"time"
)

type IContentService interface {
	List(context.Context, int, int, string, []string) ([]domain.Content, error)
	ListCategories(context.Context) ([]domain.Category, error)
	ListBadges(context.Context) ([]domain.Badge, error)
	Get(context.Context, string) (domain.Content, error)
	Update(context.Context, domain.Content, domain.Session) (domain.Content, error)
}
type ContentService struct {
	repo repository.IContentRepository
	auth *AuthService
}

func NewContentService(r repository.IContentRepository, a *AuthService) *ContentService {
	return &ContentService{repo: r, auth: a}
}
func (s *ContentService) List(ctx context.Context, limit, offset int, category string, badges []string) ([]domain.Content, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset, category, badges)
}
func (s *ContentService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx)
}
func (s *ContentService) ListBadges(ctx context.Context) ([]domain.Badge, error) {
	return s.repo.ListBadges(ctx)
}
func (s *ContentService) Get(ctx context.Context, slug string) (domain.Content, error) {
	if !validSlug(slug) {
		return domain.Content{}, errors.New("invalid slug")
	}
	return s.repo.FindBySlug(ctx, slug)
}
func (s *ContentService) Update(ctx context.Context, c domain.Content, session domain.Session) (domain.Content, error) {
	if s.auth == nil || !s.auth.CanEdit(ctx, session) {
		return domain.Content{}, errors.New("forbidden")
	}
	if !validSlug(c.Slug) || strings.TrimSpace(c.Title) == "" || strings.TrimSpace(c.Body) == "" {
		return domain.Content{}, errors.New("invalid content")
	}
	old, err := s.repo.FindBySlug(ctx, c.Slug)
	if err != nil {
		return domain.Content{}, err
	}
	c.ID = old.ID
	c.CreatedAt = old.CreatedAt
	c.CreatedBy = old.CreatedBy
	c.UpdatedAt = time.Now()
	if session.UserID != nil {
		c.UpdatedBy = *session.UserID
	}
	snapshot, _ := json.Marshal(c)
	h := domain.ContentHistory{ContentID: c.ID, Operation: "updated", CreatedAt: c.UpdatedAt, Snapshot: string(snapshot), ActorID: session.UserID, ActorAdminID: session.AdminID}
	if err = s.repo.Update(ctx, c, h); err != nil {
		return domain.Content{}, err
	}
	return c, nil
}

var slugRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validSlug(s string) bool { return len(s) <= 200 && slugRE.MatchString(s) }
