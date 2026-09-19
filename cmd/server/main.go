package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/example/cms/internal/config"
	"github.com/example/cms/internal/repository"
	"github.com/example/cms/internal/security"
	"github.com/example/cms/internal/service"
	"github.com/example/cms/internal/web"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/postgres/*.sql
var migrations embed.FS

func main() {
	cfg := config.Load()
	applicationDB, err := gorm.Open(postgres.Open(cfg.ApplicationDatabaseURL), &gorm.Config{})
	if err != nil {
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
	content := service.NewContentService(repository.NewContentRepository(applicationDB))
	sso := security.NewSSOTokenVerifier(cfg.SSOHeader, cfg.SSOToken, cfg.SSOCookie)
	log.Printf("CMS listening on %s", cfg.HTTPAddr)
	if err = http.ListenAndServe(cfg.HTTPAddr, web.NewHandler(content, sso).Routes(cfg.AllowedOrigins)); err != nil {
		log.Fatal(err)
	}
}
