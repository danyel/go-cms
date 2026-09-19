// Package config carries the runtime settings for the demo CMS. There is no
// database to configure: content is seeded in memory and optionally mirrored to
// a single JSON file.
package config

import "os"

type Config struct {
	HTTPAddr, Profile, ContentFile, PostgresURL    string
	AllowedOrigins, SSOHeader, SSOToken, SSOCookie string
}

func Load() Config {
	return Config{
		HTTPAddr:       get("CMS_HTTP_ADDR", ":8080"),
		Profile:        get("CMS_PROFILE", "demo"),
		ContentFile:    get("CMS_CONTENT_FILE", ""),
		PostgresURL:    get("CMS_POSTGRES_URL", ""),
		AllowedOrigins: get("CMS_ALLOWED_ORIGINS", "http://localhost:5173"),
		SSOHeader:      get("CMS_SSO_HEADER", "X-SSO-Token"),
		SSOToken:       get("CMS_SSO_TOKEN", ""),
		SSOCookie:      get("CMS_SSO_COOKIE", ""),
	}
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
