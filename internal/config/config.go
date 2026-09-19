package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr, AdminDatabaseURL, ApplicationDatabaseURL, AllowedOrigins, GoogleClientID, GoogleClientSecret, GoogleRedirectURL string
	SessionTTL                                                                                                                time.Duration
}

func Load() (Config, error) {
	clientID, clientSecret, err := googleCredentials()
	if err != nil {
		return Config{}, err
	}
	return Config{
		HTTPAddr:               get("CMS_HTTP_ADDR", ":8080"),
		AdminDatabaseURL:       get("CMS_ADMIN_DATABASE_URL", "./cms-admin.db"),
		ApplicationDatabaseURL: get("CMS_APPLICATION_DATABASE_URL", "host=localhost user=cms password=cms dbname=cms port=5432 sslmode=disable"),
		AllowedOrigins:         get("CMS_ALLOWED_ORIGINS", "http://localhost:5173"),
		SessionTTL:             duration("CMS_SESSION_TTL", 24*time.Hour),
		GoogleClientID:         get("CMS_GOOGLE_CLIENT_ID", clientID),
		GoogleClientSecret:     get("CMS_GOOGLE_CLIENT_SECRET", clientSecret),
		GoogleRedirectURL:      get("CMS_GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
	}, nil
}

func googleCredentials() (string, string, error) {
	path := os.Getenv("CMS_GOOGLE_CLIENT_FILE")
	if path == "" {
		return "", "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("read Google OAuth client file: %w", err)
	}
	var credentials struct {
		Web struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"web"`
		Installed struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"installed"`
	}
	if err := json.Unmarshal(data, &credentials); err != nil {
		return "", "", fmt.Errorf("parse Google OAuth client file: %w", err)
	}
	if credentials.Web.ClientID != "" || credentials.Web.ClientSecret != "" {
		return credentials.Web.ClientID, credentials.Web.ClientSecret, nil
	}
	return credentials.Installed.ClientID, credentials.Installed.ClientSecret, nil
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
