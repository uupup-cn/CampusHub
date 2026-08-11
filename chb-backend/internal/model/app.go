package model

type App struct {
    BaseModel
    AppName        string  `gorm:"type:varchar(128);not null" json:"app_name"`
    AppDescription *string `gorm:"type:text" json:"app_description"`
    ClientID       string  `gorm:"type:varchar(64);uniqueIndex;not null" json:"client_id"`
    ClientSecret   string  `gorm:"type:varchar(128);not null" json:"-"`
    RedirectURIs   string  `gorm:"type:jsonb;not null" json:"redirect_uris"`
    Scopes         string  `gorm:"type:jsonb;not null" json:"scopes"`
    MinTrustLevel  int16   `gorm:"not null;default:0" json:"min_trust_level"`
    FeeRate        float64 `gorm:"type:decimal(4,2);not null;default:10.00" json:"fee_rate"`
    BoundUserID    *int64  `json:"bound_user_id"`
    Status         string  `gorm:"type:varchar(20);not null;default:active" json:"status"`
}
