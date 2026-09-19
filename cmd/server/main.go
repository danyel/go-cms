package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/example/cms/internal/config"
	"github.com/example/cms/internal/security"
	"github.com/example/cms/internal/service"
	"github.com/example/cms/internal/web"
)

// web/dist holds the compiled React app. A placeholder file is committed so the
// embed directive always resolves, even when the frontend has not been built.
//
//go:embed all:web/dist
var ui embed.FS

func main() {
	cfg := config.Load()

	content, err := service.NewContentService(service.SeedContent(), cfg.ContentFile)
	if err != nil {
		log.Fatal(err)
	}

	assets, err := fs.Sub(ui, "web/dist")
	if err != nil {
		log.Fatal(err)
	}

	sso := security.NewSSOTokenVerifier(cfg.SSOHeader, cfg.SSOToken, cfg.SSOCookie)
	var demo *service.DemoAuth
	if cfg.Profile == "demo" {
		demo = service.NewDemoAuth()
		log.Printf("CMS profile: demo (in-memory, %d seeded articles)", len(service.SeedContent()))
	} else {
		log.Printf("CMS profile: %s (SSO identity required; PostgreSQL URL configured: %t)", cfg.Profile, cfg.PostgresURL != "")
	}
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           web.NewHandler(content, assets, sso, demo).Routes(cfg.AllowedOrigins),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("CMS listening on %s", cfg.HTTPAddr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
