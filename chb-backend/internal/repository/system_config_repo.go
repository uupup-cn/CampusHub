package repository

import "time"

type SystemConfigModel struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigKey   string    `json:"config_key" gorm:"column:config_key"`
	ConfigValue string    `json:"config_value" gorm:"column:config_value"`
	Description string    `json:"description"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SystemConfigModel) TableName() string {
	return "system_configs"
}
