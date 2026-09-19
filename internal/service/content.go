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

// IContentService exposes the content use cases. Read operations are anonymous;
// Update is only reachable once the web layer verified the SSO proxy identity.
type IContentService interface {
	List(context.Context, int, int, string, string, []string) ([]domain.Content, error)
	ListCategories(context.Context) ([]domain.Category, error)
	ListBadges(context.Context) ([]domain.Badge, error)
	Get(context.Context, string) (domain.Content, error)
	Update(context.Context, domain.Content) (domain.Content, error)
}
type ContentService struct {
	repo repository.IContentRepository
}

func NewContentService(r repository.IContentRepository) *ContentService {
	return &ContentService{repo: r}
}
func (s *ContentService) List(ctx context.Context, limit, offset int, category, status string, badges []string) ([]domain.Content, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset, category, status, badges)
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
func (s *ContentService) Update(ctx context.Context, c domain.Content) (domain.Content, error) {
	if !validSlug(c.Slug) || strings.TrimSpace(c.Title) == "" || strings.TrimSpace(c.Body) == "" {
		return domain.Content{}, errors.New("invalid content")
	}
	old, err := s.repo.FindBySlug(ctx, c.Slug)
	if err != nil {
		return domain.Content{}, err
	}
	c.ID = old.ID
	c.CreatedAt = old.CreatedAt
	c.UpdatedAt = time.Now()
	c.Badges = normalizeBadges(c.Badges)

	snapshot, _ := json.Marshal(c)
	h := domain.ContentHistory{ContentID: c.ID, Operation: "updated", CreatedAt: c.UpdatedAt, Snapshot: string(snapshot)}
	if err = s.repo.Update(ctx, c, h); err != nil {
		return domain.Content{}, err
	}
	return c, nil
}

func normalizeBadges(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" || len(value) > 100 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

var slugRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validSlug(s string) bool { return len(s) <= 200 && slugRE.MatchString(s) }
