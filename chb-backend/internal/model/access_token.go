package model

import "time"

type AccessToken struct {
    ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    AccessToken      string    `gorm:"type:varchar(256);uniqueIndex;not null" json:"access_token"`
    RefreshToken    *string    `gorm:"type:varchar(128);uniqueIndex" json:"refresh_token"`
    ClientID         string    `gorm:"type:varchar(64);not null" json:"client_id"`
    DiscourseUserID  int64     `gorm:"not null" json:"discourse_user_id"`
    Scopes           string    `gorm:"type:jsonb;not null" json:"scopes"`
    AccessExpiresAt  time.Time `gorm:"not null" json:"access_expires_at"`
    RefreshExpiresAt time.Time `gorm:"not null" json:"refresh_expires_at"`
    Revoked          bool      `gorm:"not null;default:false" json:"revoked"`
    CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}
