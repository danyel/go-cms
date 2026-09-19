package repository

import (
	"context"
	"github.com/example/cms/internal/domain"
	"github.com/example/cms/internal/persistence"
	"gorm.io/gorm"
)

type GORMContentRepository struct{ db *gorm.DB }

func NewContentRepository(db *gorm.DB) IContentRepository { return &GORMContentRepository{db} }
func (r *GORMContentRepository) List(ctx context.Context, limit, offset int, category, status string, badges []string) ([]domain.Content, error) {
	var rows []persistence.ContentModel
	query := r.db.WithContext(ctx).Model(&persistence.ContentModel{}).Preload("Category").Preload("Badges")
	if category != "" {
		query = query.Where("category_id IN (SELECT id FROM categories WHERE slug = ?)", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	for _, badge := range badges {
		query = query.Where("EXISTS (SELECT 1 FROM content_badges cb JOIN badges b ON b.id = cb.badge_id WHERE cb.content_id = content.id AND b.name = ?)", badge)
	}
	err := query.Order("updated_at DESC, created_at DESC, id DESC").Limit(limit).Offset(offset).Find(&rows).Error
	out := make([]domain.Content, len(rows))
	for i := range rows {
		out[i] = persistence.ContentToDomain(rows[i])
	}
	return out, err
}
func (r *GORMContentRepository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	var rows []persistence.CategoryModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Category, len(rows))
	for i := range rows {
		out[i] = domain.Category{ID: rows[i].ID, Slug: rows[i].Slug, Name: rows[i].Name}
	}
	return out, nil
}
func (r *GORMContentRepository) ListBadges(ctx context.Context) ([]domain.Badge, error) {
	var rows []persistence.BadgeModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Badge, len(rows))
	for i := range rows {
		out[i] = domain.Badge{ID: rows[i].ID, Name: rows[i].Name}
	}
	return out, nil
}
func (r *GORMContentRepository) FindBySlug(ctx context.Context, slug string) (domain.Content, error) {
	var row persistence.ContentModel
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Badges").Where("slug = ?", slug).First(&row).Error; err != nil {
		return domain.Content{}, err
	}
	return persistence.ContentToDomain(row), nil
}
func (r *GORMContentRepository) Update(ctx context.Context, d domain.Content, h domain.ContentHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&persistence.ContentModel{}).Where("id = ?", d.ID).Updates(map[string]any{"slug": d.Slug, "title": d.Title, "summary": d.Summary, "body": d.Body, "status": d.Status, "published": d.Published, "updated_at": d.UpdatedAt}).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM content_badges WHERE content_id = ?", d.ID).Error; err != nil {
			return err
		}
		for _, badge := range d.Badges {
			if err := tx.Exec("INSERT INTO badges (name) VALUES (?) ON CONFLICT (name) DO NOTHING", badge).Error; err != nil {
				return err
			}
			if err := tx.Exec("INSERT INTO content_badges (content_id, badge_id) SELECT ?, id FROM badges WHERE name = ? ON CONFLICT DO NOTHING", d.ID, badge).Error; err != nil {
				return err
			}
		}
		return tx.Create(&persistence.ContentHistoryModel{ContentID: h.ContentID, Operation: h.Operation, CreatedAt: h.CreatedAt, Snapshot: h.Snapshot}).Error
	})
}
