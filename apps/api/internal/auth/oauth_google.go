package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

// GoogleProfile is the subset of a Google account used to provision or link a user.
type GoogleProfile struct {
	Sub           string
	Email         string
	Name          string
	EmailVerified bool
}

// GoogleAuthenticator abstracts the Google OAuth flow so the HTTP layer and tests do
// not depend on a live Google call. The real implementation is GoogleOAuth.
type GoogleAuthenticator interface {
	// AuthCodeURL returns Google's consent URL for the Authorization Code flow, binding
	// the given anti-CSRF state.
	AuthCodeURL(state string) string
	// ExchangeAndFetchProfile exchanges an authorization code for tokens and returns the
	// user's Google profile.
	ExchangeAndFetchProfile(ctx context.Context, code string) (GoogleProfile, error)
}

// GoogleOAuth is the production GoogleAuthenticator backed by golang.org/x/oauth2.
type GoogleOAuth struct {
	config *oauth2.Config
}

// NewGoogleOAuth builds an authenticator. It returns nil when credentials are not
// configured, so callers can treat Google login as disabled.
func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil
	}
	return &GoogleOAuth{config: &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email", "profile"},
	}}
}

func (g *GoogleOAuth) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (g *GoogleOAuth) ExchangeAndFetchProfile(ctx context.Context, code string) (GoogleProfile, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("exchange google code: %w", err)
	}

	resp, err := g.config.Client(ctx, token).Get(googleUserInfoURL)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("fetch google profile: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return GoogleProfile{}, fmt.Errorf("google userinfo status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("read google profile: %w", err)
	}

	var payload struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return GoogleProfile{}, fmt.Errorf("decode google profile: %w", err)
	}
	if payload.Sub == "" {
		return GoogleProfile{}, fmt.Errorf("google profile missing subject id")
	}

	return GoogleProfile{
		Sub:           payload.Sub,
		Email:         payload.Email,
		Name:          payload.Name,
		EmailVerified: payload.EmailVerified,
	}, nil
}

var _ GoogleAuthenticator = (*GoogleOAuth)(nil)
