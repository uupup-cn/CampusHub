package model

type UserBalance struct {
    BaseModel
    DiscourseUserID int64  `gorm:"uniqueIndex;not null" json:"discourse_user_id"`
    Username        string `gorm:"type:varchar(60);not null" json:"username"`
    Balance         int64  `gorm:"not null;default:0" json:"balance"`
    Version         int64  `gorm:"not null;default:1" json:"version"`
    TrustLevel      int16  `gorm:"not null;default:0" json:"trust_level"`
    TotalEarned     int64  `gorm:"not null;default:0" json:"total_earned"`
    TotalSpent      int64  `gorm:"not null;default:0" json:"total_spent"`
    Status          string `gorm:"type:varchar(20);not null;default:active" json:"status"`
}
