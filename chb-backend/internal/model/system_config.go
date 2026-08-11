package model

type SystemConfig struct {
    BaseModel
    ConfigKey   string  `gorm:"type:varchar(128);uniqueIndex;not null" json:"config_key"`
    ConfigValue string  `gorm:"type:jsonb;not null" json:"config_value"`
    Description *string `gorm:"type:varchar(512)" json:"description"`
    UpdatedBy   *int64  `json:"updated_by"`
}
