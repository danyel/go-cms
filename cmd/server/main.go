package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/example/cms/internal/config"
	"github.com/example/cms/internal/identity"
	"github.com/example/cms/internal/repository"
	"github.com/example/cms/internal/service"
	"github.com/example/cms/internal/web"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed migrations/admin/*.sql migrations/postgres/*.sql
var migrations embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	adminDB, err := gorm.Open(sqlite.Open(cfg.AdminDatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	applicationDB, err := gorm.Open(postgres.Open(cfg.ApplicationDatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	adminSQLDB, err := adminDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	adminMigrations, err := fs.Sub(migrations, "migrations/admin")
	if err != nil {
		log.Fatal(err)
	}
	adminStore, err := database.NewStore(goose.DialectSQLite3, "goose_admin_db_version")
	if err != nil {
		log.Fatal(err)
	}
	adminProvider, err := goose.NewProvider("", adminSQLDB, adminMigrations, goose.WithStore(adminStore))
	if err != nil {
		log.Fatal(err)
	}
	if _, err = adminProvider.Up(context.Background()); err != nil {
		log.Fatal(err)
	}
	applicationSQLDB, err := applicationDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	applicationMigrations, err := fs.Sub(migrations, "migrations/postgres")
	if err != nil {
		log.Fatal(err)
	}
	applicationStore, err := database.NewStore(goose.DialectPostgres, "goose_application_db_version")
	if err != nil {
		log.Fatal(err)
	}
	applicationProvider, err := goose.NewProvider("", applicationSQLDB, applicationMigrations, goose.WithStore(applicationStore))
	if err != nil {
		log.Fatal(err)
	}
	if _, err = applicationProvider.Up(context.Background()); err != nil {
		log.Fatal(err)
	}
	auth := service.NewAuthServiceWithSeparateSessions(
		repository.NewUserRepository(applicationDB),
		repository.NewAdminRepository(adminDB),
		repository.NewSessionRepository(applicationDB),
		repository.NewSessionRepository(adminDB),
		cfg.SessionTTL,
	)
	content := service.NewContentService(repository.NewContentRepository(applicationDB), auth)
	google := identity.GoogleProvider{ClientID: cfg.GoogleClientID, ClientSecret: cfg.GoogleClientSecret, RedirectURL: cfg.GoogleRedirectURL}
	log.Printf("CMS listening on %s", cfg.HTTPAddr)
	if err = http.ListenAndServe(cfg.HTTPAddr, web.NewHandlerWithContent(auth, google, content).Routes(cfg.AllowedOrigins)); err != nil {
		log.Fatal(err)
	}
}
