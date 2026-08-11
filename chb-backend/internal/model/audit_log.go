package model

import "time"

type AuditLog struct {
    ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    OperatorID  *int64    `json:"operator_id"`
    Action      string    `gorm:"type:varchar(64);not null" json:"action"`
    TargetType  *string   `gorm:"type:varchar(32)" json:"target_type"`
    TargetID    *int64    `json:"target_id"`
    Detail      *string   `gorm:"type:jsonb" json:"detail"`
    IPAddress   *string   `gorm:"type:varchar(45)" json:"ip_address"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
