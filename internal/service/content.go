// Package service holds the application use cases. The demo CMS keeps its
// content in memory: it is seeded at startup and optionally mirrored to a single
// JSON file, so the container needs no database server.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Content is a single article in the backlog.
type Content struct {
	ID          uint       `json:"id"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Body        string     `json:"body"`
	Status      string     `json:"status"`
	Published   bool       `json:"published"`
	Category    string     `json:"category,omitempty"`
	Badges      []string   `json:"badges"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Category groups content for the filter dropdown.
type Category struct {
	ID   uint   `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Badge is a reusable tag shared by any number of articles.
type Badge struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// contentRevision records one edit so the app keeps an audit trail.
type contentRevision struct {
	ContentID uint      `json:"contentId"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"createdAt"`
	Snapshot  string    `json:"snapshot"`
}

// IContentService exposes the content use cases. Reads are anonymous; Update is
// only reachable once the web layer verified the SSO identity.
type IContentService interface {
	List(context.Context, int, int, string, string, []string) ([]Content, error)
	ListCategories(context.Context) ([]Category, error)
	ListBadges(context.Context) ([]Badge, error)
	Get(context.Context, string) (Content, error)
	Update(context.Context, Content) (Content, error)
}

// ContentService keeps the demo content in memory behind a mutex. When file is
// set, the store is loaded from that JSON file on startup and written back after
// every edit, so a mounted volume keeps changes across restarts.
type ContentService struct {
	mu            sync.RWMutex
	items         []Content
	history       []contentRevision
	categories    map[string]string
	categoryOrder []string
	allBadges     []string
	nextID        uint
	file          string
}

// NewContentService seeds the store from the given articles. When file is not
// empty and already exists, the store is loaded from it instead, which lets a
// mounted volume keep edits across restarts.
func NewContentService(seed []Content, file string) (*ContentService, error) {
	s := &ContentService{file: file}
	if file != "" {
		if _, err := os.Stat(file); err == nil {
			if err := s.load(file); err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	if len(seed) == 0 {
		return nil, errors.New("no seed content available")
	}
	s.applySeed(seed)
	if err := s.persistLocked(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *ContentService) applySeed(seed []Content) {
	s.categories = map[string]string{}
	s.categoryOrder = nil
	s.items = make([]Content, 0, len(seed))
	s.history = nil
	known := map[string]struct{}{}
	for i, item := range seed {
		item.ID = uint(i + 1)
		item.Category = normalizeCategory(item.Category)
		if _, seen := s.categories[item.Category]; !seen && item.Category != "" {
			s.categories[item.Category] = categoryDisplay(item.Category)
			s.categoryOrder = append(s.categoryOrder, item.Category)
		}
		item.Badges = normalizeBadges(item.Badges)
		for _, badge := range item.Badges {
			known[badge] = struct{}{}
		}
		item.Status = normalizeStatus(item.Status)
		if item.CreatedAt.IsZero() {
			item.CreatedAt = time.Now()
		}
		if item.UpdatedAt.IsZero() {
			item.UpdatedAt = item.CreatedAt
		}
		s.items = append(s.items, item)
	}
	s.nextID = uint(len(seed) + 1)
	s.allBadges = sortedKeys(known)
	sortContent(s.items)
}

// List returns a page of articles filtered by category, status, and badges. The
// selected badges are combined with AND semantics.
func (s *ContentService) List(_ context.Context, limit, offset int, category, status string, badges []string) ([]Content, error) {
	limit, offset = clampPage(limit, offset)
	category = normalizeCategory(category)
	if strings.TrimSpace(status) != "" {
		status = normalizeStatus(status)
	}
	wanted := normalizeBadges(badges)

	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Content, 0, limit)
	for _, item := range s.items {
		if category != "" && item.Category != category {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		if !hasBadges(item, wanted) {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		if len(out) == limit {
			break
		}
		out = append(out, cloneContent(item))
	}
	return out, nil
}

// ListCategories returns the fixed category list used by the filter dropdown.
func (s *ContentService) ListCategories(context.Context) ([]Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Category, 0, len(s.categoryOrder))
	for i, slug := range s.categoryOrder {
		out = append(out, Category{ID: uint(i + 1), Slug: slug, Name: s.categories[slug]})
	}
	return out, nil
}

// ListBadges returns every badge carried by an article, so the filter and the
// editor never suggest a badge that would return nothing.
func (s *ContentService) ListBadges(context.Context) ([]Badge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	known := map[string]struct{}{}
	for _, badge := range s.allBadges {
		known[badge] = struct{}{}
	}
	for _, item := range s.items {
		for _, badge := range item.Badges {
			known[badge] = struct{}{}
		}
	}
	names := sortedKeys(known)
	out := make([]Badge, 0, len(names))
	for i, name := range names {
		out = append(out, Badge{ID: uint(i + 1), Name: name})
	}
	return out, nil
}

// Get returns one article by slug.
func (s *ContentService) Get(_ context.Context, slug string) (Content, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if !validSlug(slug) {
		return Content{}, errors.New("invalid slug")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.Slug == slug {
			return cloneContent(item), nil
		}
	}
	return Content{}, errors.New("content not found")
}

// Update saves an edited article, records a revision, and mirrors the store to
// disk when a content file is configured.
func (s *ContentService) Update(_ context.Context, c Content) (Content, error) {
	return s.update(c, "")
}

func (s *ContentService) UpdateAs(_ context.Context, c Content, actor string) (Content, error) {
	return s.update(c, actor)
}

func (s *ContentService) update(c Content, actor string) (Content, error) {
	c.Slug = strings.ToLower(strings.TrimSpace(c.Slug))
	if !validSlug(c.Slug) || strings.TrimSpace(c.Title) == "" || strings.TrimSpace(c.Body) == "" {
		return Content{}, errors.New("invalid content")
	}
	c.Badges = normalizeBadges(c.Badges)
	c.Status = normalizeStatus(c.Status)
	c.Published = c.Status == "published"

	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1
	for i := range s.items {
		if s.items[i].Slug == c.Slug {
			index = i
			break
		}
	}
	if index < 0 {
		return Content{}, errors.New("content not found")
	}
	current := s.items[index]
	if actor != "" && current.Author != actor {
		return Content{}, errors.New("you can only edit your own content")
	}
	c.ID = current.ID
	c.CreatedAt = current.CreatedAt
	c.UpdatedAt = time.Now().Truncate(time.Second)
	c.Category = normalizeCategory(current.Category)
	c.Author = current.Author
	if c.Status == "published" && c.PublishedAt == nil {
		now := time.Now().Truncate(time.Second)
		c.PublishedAt = &now
	}
	if c.Status != "published" {
		c.PublishedAt = nil
	}

	s.items[index] = c
	sortContent(s.items)
	snapshot, _ := json.Marshal(c)
	s.history = append(s.history, contentRevision{ContentID: c.ID, Operation: "updated", CreatedAt: c.UpdatedAt, Snapshot: string(snapshot)})
	if err := s.persistLocked(); err != nil {
		return Content{}, err
	}
	return cloneContent(c), nil
}

// storedContent is the on-disk shape of an article. Identifiers are positional
// so the JSON file stays readable and diff-friendly.
type storedContent struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	Published bool      `json:"published"`
	Category  string    `json:"category,omitempty"`
	Badges    []string  `json:"badges"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type storedFile struct {
	Categories map[string]string `json:"categories"`
	Content    []storedContent   `json:"content"`
	History    []contentRevision `json:"history,omitempty"`
}

func (s *ContentService) load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var payload storedFile
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	s.categories = payload.Categories
	if s.categories == nil {
		s.categories = map[string]string{}
	}
	s.categoryOrder = make([]string, 0, len(s.categories))
	for slug := range s.categories {
		s.categoryOrder = append(s.categoryOrder, slug)
	}
	sort.Strings(s.categoryOrder)
	s.history = payload.History

	known := map[string]struct{}{}
	s.items = make([]Content, 0, len(payload.Content))
	for i, item := range payload.Content {
		c := Content{
			ID: uint(i + 1), Slug: item.Slug, Title: item.Title, Summary: item.Summary, Body: item.Body,
			Status: normalizeStatus(item.Status), Category: normalizeCategory(item.Category),
			Badges: normalizeBadges(item.Badges), CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		}
		c.Published = c.Status == "published"
		for _, badge := range c.Badges {
			known[badge] = struct{}{}
		}
		s.items = append(s.items, c)
	}
	s.nextID = uint(len(s.items) + 1)
	s.allBadges = sortedKeys(known)
	sortContent(s.items)
	return nil
}

func (s *ContentService) persistLocked() error {
	if s.file == "" {
		return nil
	}
	payload := storedFile{Categories: s.categories, History: s.history}
	payload.Content = make([]storedContent, 0, len(s.items))
	for _, item := range s.items {
		payload.Content = append(payload.Content, storedContent{
			Slug: item.Slug, Title: item.Title, Summary: item.Summary, Body: item.Body,
			Status: item.Status, Published: item.Published, Category: item.Category,
			Badges: item.Badges, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	if dir := filepath.Dir(s.file); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := s.file + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.file)
}

func clampPage(limit, offset int) (int, int) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func hasBadges(item Content, wanted []string) bool {
	for _, badge := range wanted {
		found := false
		for _, have := range item.Badges {
			if have == badge {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// sortContent orders articles newest first, matching the list page ordering.
func sortContent(items []Content) {
	sort.SliceStable(items, func(i, j int) bool {
		if !items[i].UpdatedAt.Equal(items[j].UpdatedAt) {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		}
		if !items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].ID > items[j].ID
	})
}

func cloneContent(c Content) Content {
	c.Badges = append([]string(nil), c.Badges...)
	return c
}

func normalizeCategory(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if !validSlug(value) {
		return ""
	}
	return value
}

func normalizeStatus(value string) string {
	switch value = strings.ToLower(strings.TrimSpace(value)); value {
	case "published", "draft", "review", "archived":
		return value
	default:
		return "draft"
	}
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
	sort.Strings(result)
	return result
}

func sortedKeys(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

// categoryDisplay returns the readable name for a category slug, falling back
// to a title-cased slug for anything not in the seeded list.
func categoryDisplay(slug string) string {
	if name, ok := categoryNames[slug]; ok {
		return name
	}
	return strings.ToUpper(slug[:1]) + slug[1:]
}

var slugRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validSlug(s string) bool { return len(s) > 0 && len(s) <= 200 && slugRE.MatchString(s) }
