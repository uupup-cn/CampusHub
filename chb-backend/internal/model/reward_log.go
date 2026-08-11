package model

import "time"

type RewardLog struct {
    ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    DiscourseUserID int64     `gorm:"not null" json:"discourse_user_id"`
    Action          string    `gorm:"type:varchar(32);not null" json:"action"`
    Amount          int64     `gorm:"not null" json:"amount"`
    RefID           int64     `gorm:"not null" json:"ref_id"`
    RefType         string    `gorm:"type:varchar(32);not null" json:"ref_type"`
    TrustLevel      int16     `gorm:"not null" json:"trust_level"`
    Multiplier      float64   `gorm:"type:decimal(3,2);not null" json:"multiplier"`
    IPAddress      *string   `gorm:"type:varchar(45)" json:"ip_address"`
    Status          string    `gorm:"type:varchar(20);not null;default:completed" json:"status"`
    RejectReason   *string   `gorm:"type:varchar(255)" json:"reject_reason"`
    CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}
