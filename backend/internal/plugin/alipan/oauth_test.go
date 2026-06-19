package alipan

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// TestGetAuthURL verifies that GetAuthURL returns a well-formed authorization
// URL containing all required OAuth parameters.
// ---------------------------------------------------------------------------
func TestGetAuthURL(t *testing.T) {
	oauth := NewAlipanOAuth("test-app-id", "test-secret", "https://example.com/callback")
	state := "test-state-token"

	authURL := oauth.GetAuthURL(state)

	// Parse the URL.
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("GetAuthURL() returned invalid URL %q: %v", authURL, err)
	}

	// Verify the base path.
	if !strings.HasPrefix(authURL, alipanAPIBase+"/oauth/authorize?") {
		t.Fatalf("URL does not start with expected prefix: %s", authURL)
	}

	// Verify required query parameters.
	params := parsed.Query()

	if params.Get("client_id") != "test-app-id" {
		t.Fatalf("client_id = %q, want %q", params.Get("client_id"), "test-app-id")
	}
	if params.Get("redirect_uri") != "https://example.com/callback" {
		t.Fatalf("redirect_uri = %q, want %q", params.Get("redirect_uri"), "https://example.com/callback")
	}
	if params.Get("response_type") != "code" {
		t.Fatalf("response_type = %q, want %q", params.Get("response_type"), "code")
	}
	if params.Get("state") != state {
		t.Fatalf("state = %q, want %q", params.Get("state"), state)
	}
	if params.Get("scope") == "" {
		t.Fatal("scope parameter is missing")
	}

	// Verify scope contains expected permissions.
	scope := params.Get("scope")
	if !strings.Contains(scope, "user:base") {
		t.Errorf("scope missing 'user:base': %s", scope)
	}
	if !strings.Contains(scope, "file:all:read") {
		t.Errorf("scope missing 'file:all:read': %s", scope)
	}
	if !strings.Contains(scope, "file:all:write") {
		t.Errorf("scope missing 'file:all:write': %s", scope)
	}
}

// ---------------------------------------------------------------------------
// TestLoadTokens verifies that LoadTokens correctly sets the access token,
// refresh token, and expiry on the AlipanOAuth instance.
// ---------------------------------------------------------------------------
func TestLoadTokens(t *testing.T) {
	oauth := NewAlipanOAuth("app-id", "secret", "https://cb.example.com")

	expiry := time.Now().Add(2 * time.Hour)
	oauth.LoadTokens("access-token-abc", "refresh-token-xyz", expiry)

	if oauth.AccessToken != "access-token-abc" {
		t.Fatalf("AccessToken = %q, want %q", oauth.AccessToken, "access-token-abc")
	}
	if oauth.RefreshToken != "refresh-token-xyz" {
		t.Fatalf("RefreshToken = %q, want %q", oauth.RefreshToken, "refresh-token-xyz")
	}
	if !oauth.TokenExpiry.Equal(expiry) {
		t.Fatalf("TokenExpiry = %v, want %v", oauth.TokenExpiry, expiry)
	}
}

// ---------------------------------------------------------------------------
// TestGetTokenInfo verifies that GetTokenInfo returns a correctly populated
// AlipanToken reflecting the current token state.
// ---------------------------------------------------------------------------
func TestGetTokenInfo(t *testing.T) {
	oauth := NewAlipanOAuth("app-id", "secret", "https://cb.example.com")

	expiry := time.Now().Add(1 * time.Hour)
	oauth.LoadTokens("my-access-token", "my-refresh-token", expiry)

	info := oauth.GetTokenInfo()

	if info.AccessToken != "my-access-token" {
		t.Fatalf("AccessToken = %q, want %q", info.AccessToken, "my-access-token")
	}
	if info.RefreshToken != "my-refresh-token" {
		t.Fatalf("RefreshToken = %q, want %q", info.RefreshToken, "my-refresh-token")
	}
	if info.TokenType != "Bearer" {
		t.Fatalf("TokenType = %q, want %q", info.TokenType, "Bearer")
	}
	if info.ExpiresIn <= 0 {
		t.Fatalf("ExpiresIn = %d, want positive", info.ExpiresIn)
	}
	// ExpiresIn should be within ~1 hour (give/take a few seconds for test execution).
	if info.ExpiresIn < 3590 || info.ExpiresIn > 3610 {
		t.Errorf("ExpiresIn = %d, expected ~3600", info.ExpiresIn)
	}
}

// ---------------------------------------------------------------------------
// TestGetTokenInfo_Expired verifies that GetTokenInfo returns ExpiresIn=0
// when the token has already expired.
// ---------------------------------------------------------------------------
func TestGetTokenInfo_Expired(t *testing.T) {
	oauth := NewAlipanOAuth("app-id", "secret", "https://cb.example.com")

	expiry := time.Now().Add(-1 * time.Hour) // already expired
	oauth.LoadTokens("expired-token", "expired-refresh", expiry)

	info := oauth.GetTokenInfo()

	if info.ExpiresIn != 0 {
		t.Fatalf("ExpiresIn = %d, want 0 for expired token", info.ExpiresIn)
	}
	if info.AccessToken != "expired-token" {
		t.Fatalf("AccessToken = %q, want %q", info.AccessToken, "expired-token")
	}
}
