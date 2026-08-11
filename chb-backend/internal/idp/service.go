package idp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type IdpService struct {
	db      *gorm.DB
	appRepo *repository.AppRepo
}

func NewIdpService(db *gorm.DB, appRepo *repository.AppRepo) *IdpService {
	return &IdpService{db: db, appRepo: appRepo}
}

type AuthRequest struct {
	ClientID    string
	RedirectURI string
	Scope       string
	State       string
	UserID      int64
	TrustLevel  int16
}

type AuthResponse struct {
	RedirectURI string
	Code        string
	State       string
	Error       string
}

type AppInfo struct {
	AppName        string   `json:"app_name"`
	AppDescription string   `json:"app_description"`
	Scopes         []string `json:"scopes"`
	MinTrustLevel  int16    `json:"min_trust_level"`
	RedirectURIs   []string `json:"redirect_uris"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

type IntrospectResponse struct {
	Active          bool   `json:"active"`
	ClientID        string `json:"client_id,omitempty"`
	DiscourseUserID int64  `json:"discourse_user_id,omitempty"`
	Username        string `json:"username,omitempty"`
	TrustLevel      int16  `json:"trust_level,omitempty"`
	Scope           string `json:"scope,omitempty"`
	Exp             int64  `json:"exp,omitempty"`
}

func (s *IdpService) Authorize(req *AuthRequest) *AuthResponse {
	// 1. Validate client
	app, err := s.appRepo.GetByClientID(req.ClientID)
	if err != nil {
		return &AuthResponse{
			RedirectURI: req.RedirectURI,
			Error:       "invalid_client",
			State:       req.State,
		}
	}

	// 2. Validate redirect URI against registered URIs
	if !appHasRedirectURI(app, req.RedirectURI) {
		return &AuthResponse{
			RedirectURI: req.RedirectURI,
			Error:       "invalid_redirect_uri",
			State:       req.State,
		}
	}

	// 3. Check trust level
	if req.TrustLevel < app.MinTrustLevel {
		return &AuthResponse{
			RedirectURI: req.RedirectURI,
			Error:       "insufficient_trust_level",
			State:       req.State,
		}
	}

	// 4. Generate authorization code
	codeBytes := make([]byte, 32)
	rand.Read(codeBytes)
	code := hex.EncodeToString(codeBytes)

	authCode := &repository.AuthCodeModel{
		Code:            code,
		ClientID:        req.ClientID,
		DiscourseUserID: req.UserID,
		Scopes:          scopesToJSON(req.Scope),
		RedirectURI:     req.RedirectURI,
		ExpiresAt:       time.Now().Add(10 * time.Minute),
		Used:            false,
	}

	if err := s.appRepo.CreateAuthCode(s.db, authCode); err != nil {
		return &AuthResponse{
			RedirectURI: req.RedirectURI,
			Error:       "server_error",
			State:       req.State,
		}
	}

	redirectURI := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, code, req.State)
	return &AuthResponse{
		RedirectURI: redirectURI,
		Code:        code,
		State:       req.State,
	}
}

// GetAppInfo 返回应用公开信息，供授权页展示。
func (s *IdpService) GetAppInfo(clientID string) (*AppInfo, error) {
	app, err := s.appRepo.GetByClientID(clientID)
	if err != nil {
		return nil, errcode.ErrParamInvalid
	}
	return appInfoFromModel(app)
}

func appHasRedirectURI(app *repository.AppModel, uri string) bool {
	info, err := appInfoFromModel(app)
	if err != nil {
		return false
	}
	for _, u := range info.RedirectURIs {
		if u == uri {
			return true
		}
	}
	return false
}

func appInfoFromModel(app *repository.AppModel) (*AppInfo, error) {
	var scopes []string
	if app.Scopes != "" {
		if err := json.Unmarshal([]byte(app.Scopes), &scopes); err != nil {
			scopes = []string{}
		}
	}
	var redirectURIs []string
	if app.RedirectURIs != "" {
		if err := json.Unmarshal([]byte(app.RedirectURIs), &redirectURIs); err != nil {
			redirectURIs = []string{}
		}
	}
	desc := ""
	if app.AppDescription != nil {
		desc = *app.AppDescription
	}
	return &AppInfo{
		AppName:        app.AppName,
		AppDescription: desc,
		Scopes:         scopes,
		MinTrustLevel:  app.MinTrustLevel,
		RedirectURIs:   redirectURIs,
	}, nil
}

func (s *IdpService) Token(req *TokenRequest) (*TokenResponse, error) {
	// Validate client
	app, err := s.appRepo.GetByClientID(req.ClientID)
	if err != nil {
		return nil, errcode.ErrParamInvalid
	}

	// Validate client secret
	if app.ClientSecret != req.ClientSecret {
		return nil, errcode.ErrPermissionDenied
	}

	switch req.GrantType {
	case "authorization_code":
		return s.handleAuthCodeGrant(app, req)
	case "refresh_token":
		return s.handleRefreshTokenGrant(app, req)
	default:
		return nil, errcode.ErrParamInvalid
	}
}

func (s *IdpService) handleAuthCodeGrant(app *repository.AppModel, req *TokenRequest) (*TokenResponse, error) {
	// Validate authorization code
	authCode, err := s.appRepo.GetAuthCode(req.Code)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}

	if authCode.ClientID != req.ClientID {
		return nil, errcode.ErrTokenInvalid
	}

	if authCode.RedirectURI != req.RedirectURI {
		return nil, errcode.ErrParamInvalid
	}

	var resp *TokenResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Mark code as used
		if err := s.appRepo.MarkAuthCodeUsed(tx, authCode.ID); err != nil {
			return errcode.ErrDatabase
		}

		// Generate tokens
		accessToken := generateToken(32)
		refreshToken := generateToken(24)
		now := time.Now()

		tokenModel := &repository.AccessTokenModel{
			AccessToken:      accessToken,
			RefreshToken:     &refreshToken,
			ClientID:         req.ClientID,
			DiscourseUserID:  authCode.DiscourseUserID,
			Scopes:           authCode.Scopes,
			AccessExpiresAt:  now.Add(2 * time.Hour),
			RefreshExpiresAt: now.Add(30 * 24 * time.Hour),
			Revoked:          false,
		}

		if err := s.appRepo.CreateAccessToken(tx, tokenModel); err != nil {
			return errcode.ErrDatabase
		}

		resp = &TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    7200,
			RefreshToken: refreshToken,
			Scope:        scopesFromJSON(authCode.Scopes),
		}
		return nil
	})

	return resp, err
}

func (s *IdpService) handleRefreshTokenGrant(app *repository.AppModel, req *TokenRequest) (*TokenResponse, error) {
	existing, err := s.appRepo.GetByRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}

	if existing.ClientID != req.ClientID {
		return nil, errcode.ErrPermissionDenied
	}

	var resp *TokenResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Revoke old token
		if err := s.appRepo.RevokeToken(tx, existing.ID); err != nil {
			return errcode.ErrDatabase
		}

		// Issue new tokens
		accessToken := generateToken(32)
		refreshToken := generateToken(24)
		now := time.Now()

		tokenModel := &repository.AccessTokenModel{
			AccessToken:      accessToken,
			RefreshToken:     &refreshToken,
			ClientID:         req.ClientID,
			DiscourseUserID:  existing.DiscourseUserID,
			Scopes:           existing.Scopes,
			AccessExpiresAt:  now.Add(2 * time.Hour),
			RefreshExpiresAt: now.Add(30 * 24 * time.Hour),
			Revoked:          false,
		}

		if err := s.appRepo.CreateAccessToken(tx, tokenModel); err != nil {
			return errcode.ErrDatabase
		}

		resp = &TokenResponse{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    7200,
			RefreshToken: refreshToken,
			Scope:        scopesFromJSON(existing.Scopes),
		}
		return nil
	})

	return resp, err
}

func (s *IdpService) Introspect(token string) (*IntrospectResponse, error) {
	at, err := s.appRepo.GetAccessToken(token)
	if err != nil {
		return &IntrospectResponse{Active: false}, nil
	}

	return &IntrospectResponse{
		Active:          true,
		ClientID:        at.ClientID,
		DiscourseUserID: at.DiscourseUserID,
		Scope:           scopesFromJSON(at.Scopes),
		Exp:             at.AccessExpiresAt.Unix(),
	}, nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

// scopesToJSON 将空格分隔的 scope 字符串转为 JSON 数组文本（存入 JSONB 列）。
func scopesToJSON(scope string) string {
	if scope == "" {
		return "[]"
	}
	var parts []string
	for _, s := range strings.Fields(scope) {
		parts = append(parts, s)
	}
	b, err := json.Marshal(parts)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// scopesFromJSON 将 JSONB 列中的 scope 数组文本还原为空格分隔字符串。
func scopesFromJSON(scopeJSON string) string {
	if scopeJSON == "" {
		return ""
	}
	var parts []string
	if err := json.Unmarshal([]byte(scopeJSON), &parts); err != nil {
		return scopeJSON
	}
	return strings.Join(parts, " ")
}
