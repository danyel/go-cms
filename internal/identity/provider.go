package identity

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type Identity struct{ Subject, Email, Name string }
type IProvider interface {
	Name() string
	Authenticate(context.Context, string) (Identity, error)
}
type GoogleProvider struct{ ClientID, ClientSecret, RedirectURL string }

func (GoogleProvider) Name() string { return "google" }
func (p GoogleProvider) Authenticate(context.Context, string) (Identity, error) {
	if p.ClientID == "" {
		return Identity{}, errors.New("google identity provider is not configured")
	}
	return Identity{}, errors.New("google authentication requires the OAuth callback")
}

func (p GoogleProvider) OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: p.ClientID, ClientSecret: p.ClientSecret, RedirectURL: p.RedirectURL,
		Endpoint: google.Endpoint, Scopes: []string{"openid", "email", "profile"},
	}
}

func (p GoogleProvider) FetchIdentity(ctx context.Context, code string) (Identity, error) {
	if p.ClientID == "" || p.ClientSecret == "" || p.RedirectURL == "" {
		return Identity{}, errors.New("google identity provider is not configured")
	}
	token, err := p.OAuthConfig().Exchange(ctx, code)
	if err != nil {
		return Identity{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return Identity{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Identity{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Identity{}, errors.New("google userinfo request failed")
	}
	var profile struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return Identity{}, err
	}
	return Identity{Subject: profile.Subject, Email: profile.Email, Name: profile.Name}, nil
}
