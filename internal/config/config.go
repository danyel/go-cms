package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, AllowedOrigins, GoogleClientID, GoogleClientSecret, GoogleRedirectURL string
	SessionTTL                                                                                   time.Duration
}

func Load() Config {
	return Config{HTTPAddr: get("CMS_HTTP_ADDR", ":8080"), DatabaseURL: get("CMS_DATABASE_URL", "./cms.db"), AllowedOrigins: get("CMS_ALLOWED_ORIGINS", "http://localhost:5173"), SessionTTL: duration("CMS_SESSION_TTL", 24*time.Hour), GoogleClientID: os.Getenv("CMS_GOOGLE_CLIENT_ID"), GoogleClientSecret: os.Getenv("CMS_GOOGLE_CLIENT_SECRET"), GoogleRedirectURL: os.Getenv("CMS_GOOGLE_REDIRECT_URL")}
}
func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if n, e := time.ParseDuration(v); e == nil {
			return n
		}
	}
	return d
}
