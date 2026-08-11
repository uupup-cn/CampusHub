package model

import "time"

type DailyRewardQuota struct {
    BaseModel
    DiscourseUserID int64     `gorm:"uniqueIndex:idx_daily_quota_user_date;not null" json:"discourse_user_id"`
    RewardDate      time.Time `gorm:"type:date;uniqueIndex:idx_daily_quota_user_date;not null" json:"reward_date"`
    EarnedToday     int64     `gorm:"not null;default:0" json:"earned_today"`
    ActionCounts    string    `gorm:"type:jsonb;not null;default:'{}'" json:"action_counts"`
}
