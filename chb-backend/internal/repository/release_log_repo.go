package repository

import (
	"time"
	"gorm.io/gorm"
)

type ReleaseLog struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Amount     int64     `json:"amount"`
	FromPool   string    `json:"from_pool"`
	ToPool     string    `json:"to_pool"`
	OperatorID int64     `json:"operator_id"`
	Reason     *string   `json:"reason"`
	IsAuto     bool      `json:"is_auto"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ReleaseLog) TableName() string {
	return "release_logs"
}

func (r *ReleaseLog) Create(db *gorm.DB) error {
	return db.Create(r).Error
}
