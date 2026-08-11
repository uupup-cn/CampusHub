package model

type TrustLevelCap struct {
    BaseModel
    TrustLevel       int16   `gorm:"uniqueIndex;not null" json:"trust_level"`
    DailyCap         int64   `gorm:"not null" json:"daily_cap"`
    RewardMultiplier float64 `gorm:"type:decimal(3,2);not null;default:1.00" json:"reward_multiplier"`
}
