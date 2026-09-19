package config

import (
	"os"
)

// Config carries the runtime settings for the CMS. The application database is
// the only datastore: identities and sessions belong to the upstream SSO proxy.
type Config struct {
	HTTPAddr, ApplicationDatabaseURL               string
	AllowedOrigins, SSOHeader, SSOToken, SSOCookie string
}

func Load() Config {
	return Config{
		HTTPAddr:               get("CMS_HTTP_ADDR", ":8080"),
		ApplicationDatabaseURL: get("CMS_APPLICATION_DATABASE_URL", "host=localhost user=cms password=cms dbname=cms port=5432 sslmode=disable"),
		AllowedOrigins:         get("CMS_ALLOWED_ORIGINS", "http://localhost:5173"),
		SSOHeader:              get("CMS_SSO_HEADER", "X-SSO-Token"),
		SSOToken:               get("CMS_SSO_TOKEN", ""),
		SSOCookie:              get("CMS_SSO_COOKIE", ""),
	}
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
