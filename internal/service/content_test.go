package service

import (
	"context"
	"path/filepath"
	"testing"
)

func TestSeedStoreListsNewestFirst(t *testing.T) {
	s, err := NewContentService(SeedContent(), "")
	if err != nil {
		t.Fatal(err)
	}
	items, err := s.List(context.Background(), 5, 0, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 5 {
		t.Fatalf("limit ignored: got %d items", len(items))
	}
	if items[0].Slug != "built-with-ai" {
		t.Fatalf("newest item = %q, want built-with-ai", items[0].Slug)
	}
	if !items[0].UpdatedAt.After(items[1].UpdatedAt) && !items[0].UpdatedAt.Equal(items[1].UpdatedAt) {
		t.Fatalf("items are not ordered newest first: %v then %v", items[0].UpdatedAt, items[1].UpdatedAt)
	}
}

func TestListFiltersByCategoryAndBadge(t *testing.T) {
	s, err := NewContentService(SeedContent(), "")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	categories, err := s.ListCategories(ctx)
	if err != nil || len(categories) == 0 {
		t.Fatalf("categories = %v err = %v", categories, err)
	}

	byCategory, err := s.List(ctx, 100, 0, "linux", "", nil)
	if err != nil || len(byCategory) == 0 {
		t.Fatalf("linux category returned %d items (err %v)", len(byCategory), err)
	}
	for _, item := range byCategory {
		if item.Category != "linux" {
			t.Fatalf("category filter leaked %q", item.Category)
		}
	}

	// Badges combine with AND semantics, so this pair must be narrower than one.
	tagged, err := s.List(ctx, 100, 0, "", "", []string{"tooling"})
	if err != nil || len(tagged) < 2 {
		t.Fatalf("tooling badge returned %d items (err %v)", len(tagged), err)
	}
	both, err := s.List(ctx, 100, 0, "", "", []string{"tooling", "docker"})
	if err != nil {
		t.Fatal(err)
	}
	if len(both) == 0 || len(both) >= len(tagged) {
		t.Fatalf("badge AND filtering wrong: tooling=%d tooling+docker=%d", len(tagged), len(both))
	}
	for _, item := range both {
		if !contains(item.Badges, "tooling") || !contains(item.Badges, "docker") {
			t.Fatalf("badge filter leaked %v", item.Badges)
		}
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestUpdatePersistsAcrossRestart(t *testing.T) {
	file := filepath.Join(t.TempDir(), "content.json")
	s, err := NewContentService(SeedContent(), file)
	if err != nil {
		t.Fatal(err)
	}
	saved, err := s.Update(context.Background(), Content{
		Slug:   "built-with-ai",
		Title:  "Edited title",
		Body:   "Edited body",
		Status: "draft",
		Badges: []string{"Go", "go", "writing"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Title != "Edited title" || saved.Published {
		t.Fatalf("unexpected save result: %+v", saved)
	}
	if len(saved.Badges) != 2 {
		t.Fatalf("badges were not normalised and deduplicated: %v", saved.Badges)
	}

	// A new service pointed at the same file must load the edit, not the seed.
	restarted, err := NewContentService(SeedContent(), file)
	if err != nil {
		t.Fatal(err)
	}
	reloaded, err := restarted.Get(context.Background(), "built-with-ai")
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Title != "Edited title" {
		t.Fatalf("edit did not survive restart: %q", reloaded.Title)
	}
}

func TestUpdateRejectsInvalidContent(t *testing.T) {
	s, err := NewContentService(SeedContent(), "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Update(context.Background(), Content{Slug: "missing", Title: "T", Body: "B"}); err == nil {
		t.Fatal("expected unknown slug to be rejected")
	}
	if _, err = s.Update(context.Background(), Content{Slug: "built-with-ai", Title: "", Body: "B"}); err == nil {
		t.Fatal("expected empty title to be rejected")
	}
}
