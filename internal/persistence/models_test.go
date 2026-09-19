package persistence

import (
	"testing"
	"time"

	"github.com/example/cms/internal/domain"
)

func TestModelsUseMigrationTableNames(t *testing.T) {
	tests := map[string]string{
		"content":         (ContentModel{}).TableName(),
		"content_history": (ContentHistoryModel{}).TableName(),
		"categories":      (CategoryModel{}).TableName(),
		"badges":          (BadgeModel{}).TableName(),
	}
	for expected, actual := range tests {
		if actual != expected {
			t.Errorf("table name = %q, want %q", actual, expected)
		}
	}
}

func TestContentMappingRoundTrip(t *testing.T) {
	d := domain.Content{ID: 2, Slug: "hello-world", Title: "Hello", Body: "body", Status: "draft", Badges: []string{"go"}, CreatedAt: time.Now()}
	got := ContentToDomain(ContentFromDomain(d))
	if got.Slug != d.Slug || got.Title != d.Title || got.Status != d.Status {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}
