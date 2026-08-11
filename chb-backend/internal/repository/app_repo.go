package repository

import (
	"time"

	"gorm.io/gorm"
)

type AppModel struct {
	ID              int64     `json:"id"`
	AppName         string    `json:"app_name"`
	AppDescription  *string   `json:"app_description"`
	ClientID        string    `json:"client_id"`
	ClientSecret    string    `json:"-"`
	RedirectURIs    string    `json:"redirect_uris"`
	Scopes          string    `json:"scopes"`
	MinTrustLevel   int16     `json:"min_trust_level"`
	FeeRate         float64   `json:"fee_rate"`
	BoundUserID     *int64    `json:"bound_user_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (AppModel) TableName() string {
	return "apps"
}

type AuthCodeModel struct {
	ID              int64
	Code            string
	ClientID        string
	DiscourseUserID int64
	Scopes          string
	RedirectURI     string
	ExpiresAt       time.Time
	Used            bool
	CreatedAt       time.Time
}

func (AuthCodeModel) TableName() string {
	return "auth_codes"
}

type AccessTokenModel struct {
	ID               int64
	AccessToken      string
	RefreshToken     *string
	ClientID         string
	DiscourseUserID  int64
	Scopes           string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	Revoked          bool
	CreatedAt        time.Time
}

func (AccessTokenModel) TableName() string {
	return "access_tokens"
}

type AppRepo struct {
	db *gorm.DB
}

func NewAppRepo(db *gorm.DB) *AppRepo {
	return &AppRepo{db: db}
}

func (r *AppRepo) GetByClientID(clientID string) (*AppModel, error) {
	var app AppModel
	err := r.db.Where("client_id = ? AND status = ?", clientID, "active").First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *AppRepo) List(page, pageSize int) ([]AppModel, int64, error) {
	var list []AppModel
	var total int64
	r.db.Model(&AppModel{}).Count(&total)
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&list).Error
	return list, total, err
}

func (r *AppRepo) Create(app *AppModel) error {
	return r.db.Create(app).Error
}

func (r *AppRepo) Update(app *AppModel) error {
	return r.db.Model(&AppModel{}).Where("id = ?", app.ID).Updates(app).Error
}

func (r *AppRepo) Delete(id int64) error {
	return r.db.Model(&AppModel{}).Where("id = ?", id).Update("status", "disabled").Error
}

func (r *AppRepo) CreateAuthCode(tx *gorm.DB, code *AuthCodeModel) error {
	return tx.Create(code).Error
}

func (r *AppRepo) GetAuthCode(code string) (*AuthCodeModel, error) {
	var ac AuthCodeModel
	err := r.db.Where("code = ? AND used = ? AND expires_at > NOW()", code, false).First(&ac).Error
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *AppRepo) MarkAuthCodeUsed(tx *gorm.DB, codeID int64) error {
	return tx.Model(&AuthCodeModel{}).Where("id = ?", codeID).Update("used", true).Error
}

func (r *AppRepo) CreateAccessToken(tx *gorm.DB, token *AccessTokenModel) error {
	return tx.Create(token).Error
}

func (r *AppRepo) GetAccessToken(token string) (*AccessTokenModel, error) {
	var at AccessTokenModel
	err := r.db.Where("access_token = ? AND revoked = ? AND access_expires_at > NOW()", token, false).First(&at).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *AppRepo) GetByRefreshToken(refreshToken string) (*AccessTokenModel, error) {
	var at AccessTokenModel
	err := r.db.Where("refresh_token = ? AND revoked = ? AND refresh_expires_at > NOW()", refreshToken, false).First(&at).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *AppRepo) RevokeToken(tx *gorm.DB, tokenID int64) error {
	return tx.Model(&AccessTokenModel{}).Where("id = ?", tokenID).Update("revoked", true).Error
}
