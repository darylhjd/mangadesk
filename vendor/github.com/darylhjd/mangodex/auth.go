package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	LoginPath        = "auth/login"
	CheckTokenPath   = "auth/check"
	LogoutPath       = "auth/logout"
	RefreshTokenPath = "auth/refresh"
)

type AuthResponse struct {
	Result  string  `json:"result"`
	Token   Token   `json:"token"`
	Message *string `json:"message"`
}

func (ar AuthResponse) GetResult() string {
	return ar.Result
}

type Token struct {
	Session string `json:"session"`
	Refresh string `json:"refresh"`
}

type TokenCheckResponse struct {
	OK              string   `json:"ok"`
	IsAuthenticated bool     `json:"isAuthenticated"`
	Roles           []string `json:"roles"`
	Permissions     []string `json:"permissions"`
}

// Login : Login to MangaDex.
// https://api.mangadex.org/docs.html#operation/post-auth-login
func (dc *DexClient) Login(user, pwd string) error {
	return dc.LoginContext(context.Background(), user, pwd)
}

// LoginContext : Login with custom context.
func (dc *DexClient) LoginContext(ctx context.Context, user, pwd string) error {
	// Create required request body.
	req := map[string]string{
		"username": user,
		"password": pwd,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	var ar AuthResponse
	if err := dc.responseOp(ctx, http.MethodPost, LoginPath, bytes.NewBuffer(rBytes), &ar); err != nil {
		return err
	}

	// Set client Token and header for authorization.
	dc.isLoggedIn = true
	dc.RefreshToken = ar.Token.Refresh
	dc.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	return nil
}

// CheckToken : Check session token validity.
// https://api.mangadex.org/docs.html#operation/get-auth-check
func (dc *DexClient) CheckToken() (bool, error) {
	return dc.CheckTokenContext(context.Background())
}

// CheckTokenContext : CheckToken with custom context.
func (dc *DexClient) CheckTokenContext(ctx context.Context) (bool, error) {
	u, _ := url.Parse(BaseAPI)
	u.Path = CheckTokenPath

	var c TokenCheckResponse
	_, err := dc.RequestAndDecode(ctx, http.MethodGet, u.String(), nil, &c)
	return c.IsAuthenticated, err
}

// Logout : Logout of MangaDex and invalidates all tokens.
// https://api.mangadex.org/docs.html#operation/post-auth-logout
func (dc *DexClient) Logout() error {
	return dc.LogoutContext(context.Background())
}

// LogoutContext : Logout with custom context.
func (dc *DexClient) LogoutContext(ctx context.Context) error {
	if err := dc.responseOp(ctx, http.MethodPost, LogoutPath, nil, nil); err != nil {
		return nil
	}

	// Remove the stored client token and also authorization header if ok.
	dc.isLoggedIn = false
	dc.RefreshToken = ""
	dc.header.Del("Authorization")
	return nil
}

// RefreshSessionToken : Refresh session token using refresh token.
// https://api.mangadex.org/docs.html#operation/post-auth-refresh
func (dc *DexClient) RefreshSessionToken() error {
	return dc.RefreshSessionTokenContext(context.Background())
}

// RefreshSessionTokenContext : RefreshToken with custom context.
func (dc *DexClient) RefreshSessionTokenContext(ctx context.Context) error {
	// Create required request body.
	req := map[string]string{
		"token": dc.RefreshToken,
	}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		dc.isLoggedIn = false
		return err
	}

	var ar AuthResponse
	if err := dc.responseOp(ctx, http.MethodPost, RefreshTokenPath, bytes.NewBuffer(rBytes), &ar); err != nil {
		dc.isLoggedIn = false
		return err
	}

	// Update tokens
	dc.isLoggedIn = true
	dc.RefreshToken = ar.Token.Refresh
	dc.header.Set("Authorization", fmt.Sprintf("Bearer %s", ar.Token.Session))
	return nil
}

// IsLoggedIn : Return true when client logged in and false otherwise.
func (dc *DexClient) IsLoggedIn() bool {
	return dc.isLoggedIn
}
