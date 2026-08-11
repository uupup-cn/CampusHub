package repository

import (
	"time"

	"gorm.io/gorm"
)

// RewardRule - 奖励规则配置
type RewardRule struct {
	ID              int64     `json:"id"`
	Action          string    `json:"action"`
	DisplayName     string    `json:"display_name"`
	Amount          int64     `json:"amount"`
	CooldownSeconds int       `json:"cooldown_seconds"`
	DailyCapPerUser int64     `json:"daily_cap_per_user"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TrustLevelCap - 等级日上限
type TrustLevelCapModel struct {
	ID               int64     `json:"id"`
	TrustLevel       int16     `json:"trust_level"`
	DailyCap         int64     `json:"daily_cap"`
	RewardMultiplier float64   `json:"reward_multiplier"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (TrustLevelCapModel) TableName() string {
	return "trust_level_caps"
}

// DailyRewardQuota - 用户每日配额
type DailyRewardQuotaModel struct {
	ID              int64     `json:"id"`
	DiscourseUserID int64     `json:"discourse_user_id"`
	RewardDate      string    `json:"reward_date"`
	EarnedToday     int64     `json:"earned_today"`
	ActionCounts    string    `json:"action_counts"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DailyRewardQuotaModel) TableName() string {
	return "daily_reward_quotas"
}

// RewardLog - 奖励发放日志
type RewardLogRecord struct {
	ID              int64     `json:"id"`
	DiscourseUserID int64     `json:"discourse_user_id"`
	Action          string    `json:"action"`
	Amount          int64     `json:"amount"`
	RefID           int64     `json:"ref_id"`
	RefType         string    `json:"ref_type"`
	TrustLevel      int16     `json:"trust_level"`
	Multiplier      float64   `json:"multiplier"`
	IPAddress       *string   `json:"ip_address"`
	Status          string    `json:"status"`
	RejectReason    *string   `json:"reject_reason"`
	CreatedAt       time.Time `json:"created_at"`
}

func (RewardLogRecord) TableName() string {
	return "reward_logs"
}

type RewardRepo struct {
	db *gorm.DB
}

func NewRewardRepo(db *gorm.DB) *RewardRepo {
	return &RewardRepo{db: db}
}

// ===== Reward Rules =====

func (r *RewardRepo) GetRuleByAction(action string) (*RewardRule, error) {
	var rule RewardRule
	err := r.db.Where("action = ? AND is_active = ?", action, true).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RewardRepo) ListRules() ([]RewardRule, error) {
	var rules []RewardRule
	err := r.db.Order("action ASC").Find(&rules).Error
	return rules, err
}

func (r *RewardRepo) UpdateRule(rule *RewardRule) error {
	return r.db.Model(&RewardRule{}).Where("action = ?", rule.Action).Updates(map[string]interface{}{
		"amount":           rule.Amount,
		"cooldown_seconds": rule.CooldownSeconds,
		"daily_cap_per_user": rule.DailyCapPerUser,
		"is_active":        rule.IsActive,
	}).Error
}

// ===== Trust Level Caps =====

func (r *RewardRepo) GetCapByLevel(trustLevel int16) (*TrustLevelCapModel, error) {
	var cap TrustLevelCapModel
	err := r.db.Where("trust_level = ?", trustLevel).First(&cap).Error
	if err != nil {
		return nil, err
	}
	return &cap, nil
}

func (r *RewardRepo) ListCaps() ([]TrustLevelCapModel, error) {
	var caps []TrustLevelCapModel
	err := r.db.Order("trust_level ASC").Find(&caps).Error
	return caps, err
}

func (r *RewardRepo) UpdateCap(cap *TrustLevelCapModel) error {
	return r.db.Model(&TrustLevelCapModel{}).Where("trust_level = ?", cap.TrustLevel).Updates(map[string]interface{}{
		"daily_cap":          cap.DailyCap,
		"reward_multiplier":  cap.RewardMultiplier,
	}).Error
}

// ===== Daily Quotas =====

func (r *RewardRepo) GetQuota(userID int64, date string) (*DailyRewardQuotaModel, error) {
	var q DailyRewardQuotaModel
	err := r.db.Where("discourse_user_id = ? AND reward_date = ?", userID, date).First(&q).Error
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *RewardRepo) UpsertQuota(tx *gorm.DB, userID int64, date string, earnedToday int64, actionCounts string) error {
	result := tx.Exec(`
		INSERT INTO daily_reward_quotas (discourse_user_id, reward_date, earned_today, action_counts, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
		ON CONFLICT (discourse_user_id, reward_date) DO UPDATE SET
			earned_today = ?, action_counts = ?, updated_at = NOW()
	`, userID, date, earnedToday, actionCounts, earnedToday, actionCounts)
	return result.Error
}

// ===== Reward Logs =====

func (r *RewardRepo) CreateLog(tx *gorm.DB, log *RewardLogRecord) error {
	return tx.Create(log).Error
}

func (r *RewardRepo) GetLogByRef(refType string, refID int64, userID int64) (*RewardLogRecord, error) {
	var log RewardLogRecord
	err := r.db.Where("ref_type = ? AND ref_id = ? AND discourse_user_id = ?", refType, refID, userID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *RewardRepo) InsertDefaultRules(tx *gorm.DB) error {
	rules := []RewardRule{
		{Action: "post", DisplayName: "发帖", Amount: 100, CooldownSeconds: 300, DailyCapPerUser: 500, IsActive: true},
		{Action: "reply", DisplayName: "回复", Amount: 50, CooldownSeconds: 60, DailyCapPerUser: 300, IsActive: true},
		{Action: "checkin", DisplayName: "签到", Amount: 30, CooldownSeconds: 86400, DailyCapPerUser: 30, IsActive: true},
		{Action: "liked", DisplayName: "被点赞", Amount: 20, CooldownSeconds: 0, DailyCapPerUser: 200, IsActive: true},
	}
	for _, rule := range rules {
		if err := tx.Create(&rule).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RewardRepo) InsertDefaultCaps(tx *gorm.DB) error {
	caps := []TrustLevelCapModel{
		{TrustLevel: 0, DailyCap: 100, RewardMultiplier: 0.50},
		{TrustLevel: 1, DailyCap: 500, RewardMultiplier: 1.00},
		{TrustLevel: 2, DailyCap: 1000, RewardMultiplier: 1.00},
		{TrustLevel: 3, DailyCap: 2000, RewardMultiplier: 1.25},
		{TrustLevel: 4, DailyCap: 5000, RewardMultiplier: 1.50},
	}
	for _, cap := range caps {
		if err := tx.Create(&cap).Error; err != nil {
			return err
		}
	}
	return nil
}
