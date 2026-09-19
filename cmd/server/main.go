package main

import (
	"context"
	"embed"
	"log"
	"net/http"

	"github.com/example/cms/internal/config"
	"github.com/example/cms/internal/identity"
	"github.com/example/cms/internal/repository"
	"github.com/example/cms/internal/service"
	"github.com/example/cms/internal/web"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	cfg := config.Load()
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, sqlDB, migrations)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = provider.Up(context.Background()); err != nil {
		log.Fatal(err)
	}
	auth := service.NewAuthService(repository.NewUserRepository(db), repository.NewAdminRepository(db), repository.NewSessionRepository(db), cfg.SessionTTL)
	google := identity.GoogleProvider{ClientID: cfg.GoogleClientID, ClientSecret: cfg.GoogleClientSecret, RedirectURL: cfg.GoogleRedirectURL}
	log.Printf("CMS listening on %s", cfg.HTTPAddr)
	if err = http.ListenAndServe(cfg.HTTPAddr, web.NewHandler(auth, google).Routes(cfg.AllowedOrigins)); err != nil {
		log.Fatal(err)
	}
}
