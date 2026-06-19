// Package alipan implements the Aliyun Drive (阿里云盘) cloud drive adapter.
package alipan

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"warmisle/internal/plugin"
)

const alipanAPIBase = "https://openapi.alipan.com"

// Ensure AlipanOAuth implements plugin.TokenProvider at compile time.
var _ plugin.TokenProvider = (*AlipanOAuth)(nil)

// AlipanToken represents the OAuth token response from AliDrive.
type AlipanToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AlipanOAuth manages OAuth2 authorization with Aliyun Drive.
// Supports both confidential client (client_secret + refresh_token) and
// public client PKCE mode (no client_secret, 30-day access token).
type AlipanOAuth struct {
	AppID        string
	AppSecret    string // Decrypted plain-text secret; empty for PKCE mode.
	RedirectURI  string
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
	codeVerifier string // PKCE code_verifier, generated during GetAuthURL
	httpClient   *http.Client
}

// NewAlipanOAuth creates a new AlipanOAuth instance with a 10-second HTTP client timeout.
func NewAlipanOAuth(appID, appSecret, redirectURI string) *AlipanOAuth {
	return &AlipanOAuth{
		AppID:       appID,
		AppSecret:   appSecret,
		RedirectURI: redirectURI,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetAuthURL constructs the authorization URL that the user visits to grant access.
// state is an anti-CSRF token returned in the OAuth callback.
// For PKCE mode (no client_secret), generates code_verifier and code_challenge.
func (a *AlipanOAuth) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", a.AppID)
	params.Set("redirect_uri", a.RedirectURI)
	params.Set("scope", "file:all:read,file:all:write")
	params.Set("response_type", "code")
	params.Set("state", state)

	// PKCE: generate code_verifier and code_challenge when no client_secret
	if a.AppSecret == "" {
		a.codeVerifier = generateCodeVerifier()
		params.Set("code_challenge", computeCodeChallenge(a.codeVerifier))
		params.Set("code_challenge_method", "S256")
	}

	return fmt.Sprintf("%s/oauth/authorize?%s", alipanAPIBase, params.Encode())
}

// GetCodeVerifier returns the PKCE code_verifier for use in the token exchange.
// Only valid when in PKCE mode (AppSecret is empty).
func (a *AlipanOAuth) GetCodeVerifier() string {
	return a.codeVerifier
}

// SetCodeVerifier sets the PKCE code_verifier (for restoring state across callbacks).
func (a *AlipanOAuth) SetCodeVerifier(v string) {
	a.codeVerifier = v
}

// ExchangeCode exchanges an authorization code for an access token.
// In PKCE mode (no client_secret), uses code_verifier instead.
func (a *AlipanOAuth) ExchangeCode(code string) error {
	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("client_id", a.AppID)

	if a.AppSecret != "" {
		// Confidential client mode: use client_secret
		body.Set("client_secret", a.AppSecret)
	} else {
		// PKCE public client mode: use code_verifier
		body.Set("code_verifier", a.codeVerifier)
	}
	body.Set("redirect_uri", a.RedirectURI)

	return a.doTokenRequest(body.Encode())
}

// refreshTokenInternal refreshes the access token using the refresh token.
// Not applicable for PKCE mode (30-day access token, no refresh).
func (a *AlipanOAuth) refreshTokenInternal() error {
	if a.RefreshToken == "" {
		return fmt.Errorf("PKCE模式不支持刷新令牌，请重新授权")
	}
	log.Printf("[alipan] refreshing access token...")

	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", a.RefreshToken)
	body.Set("client_id", a.AppID)
	if a.AppSecret != "" {
		body.Set("client_secret", a.AppSecret)
	}

	return a.doTokenRequest(body.Encode())
}

// generateCodeVerifier creates a cryptographically random PKCE code_verifier (43-128 chars).
func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	// base64url without padding = 43 chars
	return base64.RawURLEncoding.EncodeToString(b)
}

// computeCodeChallenge computes the S256 PKCE code_challenge from a code_verifier.
func computeCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// doTokenRequest sends a token request and parses the common OAuth response.
func (a *AlipanOAuth) doTokenRequest(formBody string) error {
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/oauth/access_token", alipanAPIBase),
		strings.NewReader(formBody))
	if err != nil {
		return fmt.Errorf("创建令牌请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("[alipan] POST %s/oauth/access_token", alipanAPIBase)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("令牌请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取令牌响应失败: %w", err)
	}

	log.Printf("[alipan] token response status=%d body=%s", resp.StatusCode, string(respBytes))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("令牌请求返回错误状态 %d: %s", resp.StatusCode, string(respBytes))
	}

	var token AlipanToken
	if err := json.Unmarshal(respBytes, &token); err != nil {
		return fmt.Errorf("解析令牌响应失败: %w", err)
	}

	if token.AccessToken == "" {
		return fmt.Errorf("令牌响应中缺少 access_token")
	}

	a.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		a.RefreshToken = token.RefreshToken
	}
	a.TokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	log.Printf("[alipan] token obtained, expires in %ds", token.ExpiresIn)
	return nil
}

// GetAccessToken returns a valid access token, refreshing if it will expire
// within the next 5 minutes.
func (a *AlipanOAuth) GetAccessToken() (string, error) {
	if a.AccessToken == "" {
		return "", fmt.Errorf("尚未授权，请先完成 OAuth 授权流程")
	}

	if time.Now().Add(5 * time.Minute).After(a.TokenExpiry) {
		if err := a.refreshTokenInternal(); err != nil {
			return "", fmt.Errorf("刷新令牌失败: %w", err)
		}
	}

	return a.AccessToken, nil
}

// LoadTokens loads existing token information (e.g., after decrypting from DB).
func (a *AlipanOAuth) LoadTokens(accessToken, refreshToken string, tokenExpiry time.Time) {
	a.AccessToken = accessToken
	a.RefreshToken = refreshToken
	a.TokenExpiry = tokenExpiry
}

// GetTokenInfo returns the current token state for serialization and storage.
func (a *AlipanOAuth) GetTokenInfo() *AlipanToken {
	expiresIn := int(time.Until(a.TokenExpiry).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}
	return &AlipanToken{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}
}
