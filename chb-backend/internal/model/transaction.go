package model

import "database/sql"

type Transaction struct {
    BaseModel
    TxType          string       `gorm:"type:varchar(32);not null" json:"tx_type"`
    DiscourseUserID int64        `gorm:"not null" json:"discourse_user_id"`
    Amount          int64        `gorm:"not null" json:"amount"`
    Fee             int64        `gorm:"not null;default:0" json:"fee"`
    NetAmount       int64        `gorm:"not null" json:"net_amount"`
    FromType        string       `gorm:"type:varchar(20);not null" json:"from_type"`
    ToType          string       `gorm:"type:varchar(20);not null" json:"to_type"`
    FromID          sql.NullInt64 `json:"from_id"`
    ToID            sql.NullInt64 `json:"to_id"`
    IdempotencyKey  string       `gorm:"type:varchar(128);uniqueIndex;not null" json:"idempotency_key"`
    RefType         sql.NullString `gorm:"type:varchar(32)" json:"ref_type"`
    RefID           sql.NullInt64 `json:"ref_id"`
    Description     sql.NullString `gorm:"type:varchar(512)" json:"description"`
    Status          string       `gorm:"type:varchar(20);not null;default:completed" json:"status"`
}
