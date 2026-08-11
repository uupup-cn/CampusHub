package model

import "time"

type AuthCode struct {
    ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    Code            string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"code"`
    ClientID        string    `gorm:"type:varchar(64);not null" json:"client_id"`
    DiscourseUserID int64     `gorm:"not null" json:"discourse_user_id"`
    Scopes          string    `gorm:"type:jsonb;not null" json:"scopes"`
    RedirectURI     string    `gorm:"type:varchar(512);not null" json:"redirect_uri"`
    ExpiresAt       time.Time `gorm:"not null" json:"expires_at"`
    Used            bool      `gorm:"not null;default:false" json:"used"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}
