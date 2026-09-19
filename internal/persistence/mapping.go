package persistence

import "github.com/example/cms/internal/domain"

func UserToDomain(m UserModel) domain.User {
	return domain.User{ID: m.ID, Email: m.Email, Name: m.Name, Provider: m.Provider, ProviderSubject: m.ProviderSubject, Role: m.Role, Editor: m.Editor, CreatedAt: m.CreatedAt}
}
func UserFromDomain(d domain.User) UserModel {
	return UserModel{ID: d.ID, Email: d.Email, Name: d.Name, Provider: d.Provider, ProviderSubject: d.ProviderSubject, Role: d.Role, Editor: d.Editor, CreatedAt: d.CreatedAt}
}
func ContentToDomain(m ContentModel) domain.Content {
	return domain.Content{ID: m.ID, Slug: m.Slug, Title: m.Title, Summary: m.Summary, Body: m.Body, Status: m.Status, Published: m.Published, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, CreatedBy: m.CreatedBy, UpdatedBy: m.UpdatedBy}
}
func ContentFromDomain(d domain.Content) ContentModel {
	return ContentModel{ID: d.ID, Slug: d.Slug, Title: d.Title, Summary: d.Summary, Body: d.Body, Status: d.Status, Published: d.Published, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt, CreatedBy: d.CreatedBy, UpdatedBy: d.UpdatedBy}
}
func AdminToDomain(m AdminModel) domain.Admin {
	return domain.Admin{ID: m.ID, Email: m.Email, Name: m.Name, CreatedAt: m.CreatedAt}
}
func SessionToDomain(m SessionModel) domain.Session {
	return domain.Session{ID: m.ID, Token: m.Token, UserID: m.UserID, AdminID: m.AdminID, ExpiresAt: m.ExpiresAt, CreatedAt: m.CreatedAt}
}
func SessionFromDomain(d domain.Session) SessionModel {
	return SessionModel{ID: d.ID, Token: d.Token, UserID: d.UserID, AdminID: d.AdminID, ExpiresAt: d.ExpiresAt, CreatedAt: d.CreatedAt}
}
