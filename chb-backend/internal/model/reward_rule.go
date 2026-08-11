package model

type RewardRule struct {
    BaseModel
    Action          string `gorm:"type:varchar(32);uniqueIndex;not null" json:"action"`
    DisplayName     string `gorm:"type:varchar(64);not null" json:"display_name"`
    Amount          int64  `gorm:"not null" json:"amount"`
    CooldownSeconds int    `gorm:"not null;default:0" json:"cooldown_seconds"`
    DailyCapPerUser int64  `gorm:"not null" json:"daily_cap_per_user"`
    IsActive        bool   `gorm:"not null;default:true" json:"is_active"`
}
