package persistence

import "github.com/example/cms/internal/domain"

func ContentToDomain(m ContentModel) domain.Content {
	badges := make([]string, len(m.Badges))
	for i := range m.Badges {
		badges[i] = m.Badges[i].Name
	}
	return domain.Content{ID: m.ID, Slug: m.Slug, Title: m.Title, Summary: m.Summary, Body: m.Body, Status: m.Status, Published: m.Published, CategoryID: m.CategoryID, Category: m.Category.Name, Badges: badges, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func ContentFromDomain(d domain.Content) ContentModel {
	return ContentModel{ID: d.ID, Slug: d.Slug, Title: d.Title, Summary: d.Summary, Body: d.Body, Status: d.Status, Published: d.Published, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}
