package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type AdminService struct {
	db          *gorm.DB
	balanceRepo *repository.UserBalanceRepo
	txRepo      *repository.TransactionRepo
	poolRepo    *repository.PoolRepo
	rewardRepo  *repository.RewardRepo
	marketRepo  *repository.MarketplaceRepo
	appRepo     *repository.AppRepo
}

func NewAdminService(
	db *gorm.DB,
	balanceRepo *repository.UserBalanceRepo,
	txRepo *repository.TransactionRepo,
	poolRepo *repository.PoolRepo,
	rewardRepo *repository.RewardRepo,
	marketRepo *repository.MarketplaceRepo,
	appRepo *repository.AppRepo,
) *AdminService {
	return &AdminService{
		db:          db,
		balanceRepo: balanceRepo,
		txRepo:      txRepo,
		poolRepo:    poolRepo,
		rewardRepo:  rewardRepo,
		marketRepo:  marketRepo,
		appRepo:     appRepo,
	}
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("150405.000000000")
	}
	return hex.EncodeToString(buf)
}

func (s *AdminService) writeAuditLog(operatorID int64, action string, targetType string, targetID int64, detail string) {
	detailJSON := fmt.Sprintf(`{"msg": "%s"}`, detail)
	s.db.Exec(
		"INSERT INTO audit_logs (operator_id, action, target_type, target_id, detail, ip_address, created_at) VALUES (?, ?, ?, ?, ?::jsonb, ?, NOW())",
		operatorID, action, targetType, targetID, detailJSON, "",
	)
}

// ===== System Config =====

type SystemSetting struct {
	MarketplaceFeeRate float64 `json:"marketplace_fee_rate"`
	AutoReleaseEnabled bool    `json:"auto_release_enabled"`
	AutoReleaseThreshold int   `json:"auto_release_threshold"`
	AutoReleaseRatio   int     `json:"auto_release_ratio"`
	AutoReleaseMonthlyCap int64 `json:"auto_release_monthly_cap"`
}

func (s *AdminService) GetSettings() *SystemSetting {
	settings := &SystemSetting{
		MarketplaceFeeRate:   10.0,
		AutoReleaseEnabled:   false,
		AutoReleaseThreshold: 80,
		AutoReleaseRatio:     50,
		AutoReleaseMonthlyCap: 10000000,
	}

	var configs []repository.SystemConfigModel
	s.db.Find(&configs)
	for _, cfg := range configs {
		switch cfg.ConfigKey {
		case "marketplace_fee_rate":
			settings.MarketplaceFeeRate = parseFloat(cfg.ConfigValue)
		case "auto_release_enabled":
			settings.AutoReleaseEnabled = cfg.ConfigValue == "true"
		case "auto_release_threshold":
			settings.AutoReleaseThreshold = parseInt(cfg.ConfigValue)
		case "auto_release_ratio":
			settings.AutoReleaseRatio = parseInt(cfg.ConfigValue)
		case "auto_release_monthly_cap":
			settings.AutoReleaseMonthlyCap = parseInt64(cfg.ConfigValue)
		}
	}
	return settings
}

func (s *AdminService) UpdateSettings(settings *SystemSetting) error {
	updates := map[string]string{
		"marketplace_fee_rate":     fmt.Sprintf("%v", settings.MarketplaceFeeRate),
		"auto_release_enabled":     fmt.Sprintf("%v", settings.AutoReleaseEnabled),
		"auto_release_threshold":   fmt.Sprintf("%d", settings.AutoReleaseThreshold),
		"auto_release_ratio":       fmt.Sprintf("%d", settings.AutoReleaseRatio),
		"auto_release_monthly_cap": fmt.Sprintf("%d", settings.AutoReleaseMonthlyCap),
	}
	for key, val := range updates {
		s.db.Model(&repository.SystemConfigModel{}).
			Where("config_key = ?", key).
			Update("config_value", val)
	}
	return nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseInt64(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}

// ===== Trust Level Caps =====

func (s *AdminService) ListTrustLevelCaps() ([]repository.TrustLevelCapModel, error) {
	return s.rewardRepo.ListCaps()
}

func (s *AdminService) UpdateTrustLevelCap(cap *repository.TrustLevelCapModel) error {
	return s.rewardRepo.UpdateCap(cap)
}

// ===== Reward Rules =====

func (s *AdminService) UpdateRewardRule(rule *repository.RewardRule) error {
	return s.rewardRepo.UpdateRule(rule)
}

// ===== Apps =====

func (s *AdminService) ListApps(page, pageSize int) ([]repository.AppModel, int64, error) {
	return s.appRepo.List(page, pageSize)
}

func (s *AdminService) CreateApp(app *repository.AppModel) error {
	app.ClientID = "chb_" + randomHex(8)
	app.ClientSecret = randomHex(24)
	if app.Status == "" {
		app.Status = "active"
	}
	err := s.appRepo.Create(app)
	if err == nil {
		s.writeAuditLog(0, "app_create", "app", app.ID, "create app: "+app.AppName)
	}
	return err
}

func (s *AdminService) UpdateApp(app *repository.AppModel) error {
	return s.appRepo.Update(app)
}

func (s *AdminService) DeleteApp(id int64) error {
	err := s.appRepo.Delete(id)
	if err == nil {
		s.writeAuditLog(0, "app_delete", "app", id, "delete app")
	}
	return err
}

// ===== Marketplace Review =====

func (s *AdminService) ListPendingApplications(page, pageSize int) ([]repository.MerchantApplicationModel, int64, error) {
	return s.marketRepo.ListPendingApplications(page, pageSize)
}

func (s *AdminService) ReviewApplication(id int64, status, comment string, reviewerID int64) error {
	update := &repository.MerchantApplicationModel{
		ID:            id,
		Status:        status,
		ReviewedBy:    &reviewerID,
		ReviewComment: &comment,
	}
	err := s.marketRepo.UpdateApplication(update)
	if err == nil {
		s.writeAuditLog(0, "merchant_approve", "application", id, "review: "+status)
	}
	return err
}

func (s *AdminService) ListPendingItems(page, pageSize int) ([]repository.MarketplaceItemModel, int64, error) {
	return s.marketRepo.ListPendingItems(page, pageSize)
}

func (s *AdminService) ReviewItem(id int64, status string) error {
	item, err := s.marketRepo.GetItem(id)
	if err != nil {
		return errcode.ErrNotFound
	}
	item.Status = status
	err = s.marketRepo.UpdateItem(item)
	if err == nil {
		s.writeAuditLog(0, "item_review", "item", id, "review: "+status)
	}
	return err
}

// ===== User Management =====

func (s *AdminService) ListUsers(page, pageSize int, keyword string) ([]repository.UserBalance, int64, error) {
	// For now, just use the balance repo list
	return s.balanceRepo.List(page, pageSize)
}

func (s *AdminService) FreezeUser(userID int64) error {
	err := s.balanceRepo.SetStatus(s.db, userID, "frozen")
	if err == nil {
		s.writeAuditLog(0, "user_freeze", "user", userID, "freeze user")
	}
	return err
}

func (s *AdminService) UnfreezeUser(userID int64) error {
	err := s.balanceRepo.SetStatus(s.db, userID, "active")
	if err == nil {
		s.writeAuditLog(0, "user_unfreeze", "user", userID, "unfreeze user")
	}
	return err
}

// ===== Stats =====

type AdminStats struct {
	TotalUsers           int64  `json:"total_users"`
	ActiveUsersToday     int64  `json:"active_users_today"`
	TotalTransactions    int64  `json:"total_transactions"`
	TodayTransactions    int64  `json:"today_transactions"`
	TotalMarketplaceOrders int64 `json:"total_marketplace_orders"`
	PublicPoolWaterLevel float64 `json:"public_pool_water_level"`
	PendingApplications  int64  `json:"pending_applications"`
	PendingItems         int64  `json:"pending_items"`
}

func (s *AdminService) GetStats() (*AdminStats, error) {
	publicPool, err := s.poolRepo.GetPublicPool()
	if err != nil {
		return nil, errcode.ErrDatabase
	}
	txCount, _ := s.txRepo.GetTotalCount()
	todayTx, _ := s.txRepo.GetTodayCount()
	waterLevel := float64(publicPool.Balance) / float64(publicPool.TotalSupply)

	_, totalUsers, _ := s.balanceRepo.List(1, 1)

	var activeUsersToday int64
	s.db.Table("transactions").Select("COUNT(DISTINCT discourse_user_id)").
		Where("DATE(created_at) = CURRENT_DATE").Count(&activeUsersToday)

	var totalOrders int64
	s.db.Table("marketplace_orders").Count(&totalOrders)

	_, pendingApps, _ := s.marketRepo.ListPendingApplications(1, 1)
	_, pendingItems, _ := s.marketRepo.ListPendingItems(1, 1)

	stats := &AdminStats{
		TotalUsers:            totalUsers,
		ActiveUsersToday:      activeUsersToday,
		TotalTransactions:     txCount,
		TodayTransactions:     todayTx,
		TotalMarketplaceOrders: totalOrders,
		PublicPoolWaterLevel:  waterLevel,
		PendingApplications:   pendingApps,
		PendingItems:          pendingItems,
	}
	return stats, nil
}

// ===== Audit Log =====

type AuditLogEntry struct {
	ID         int64     `json:"id"`
	OperatorID *int64    `json:"operator_id"`
	Action     string    `json:"action"`
	TargetType *string   `json:"target_type"`
	TargetID   *int64    `json:"target_id"`
	Detail     *string   `json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *AdminService) ListAuditLogs(action string, operatorID int64, startDate, endDate string, page, pageSize int) ([]AuditLogEntry, int64, error) {
	query := s.db.Table("audit_logs")
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if operatorID > 0 {
		query = query.Where("operator_id = ?", operatorID)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	var total int64
	query.Count(&total)

	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 100 { pageSize = 20 }
	offset := (page - 1) * pageSize

	var logs []AuditLogEntry
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)
	return logs, total, nil
}

// ===== User Recover =====

func (s *AdminService) RecoverPoints(userID int64, amount int64, reason string, operatorID int64) error {
	if amount <= 0 {
		return errcode.ErrParamInvalid
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		ub, err := s.balanceRepo.GetByUserIDWithLock(tx, userID)
		if err != nil {
			return errcode.ErrNotFound
		}
		if ub.Balance < amount {
			return errcode.ErrBalanceInsufficient
		}
		newBalance := ub.Balance - amount
		if err := s.balanceRepo.UpdateBalance(tx, userID, newBalance, ub.Version+1); err != nil {
			return err
		}
		// 追回积分回注公共池（事务内操作，确保原子性）
		publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}
		if err := s.poolRepo.UpdateBalanceTx(tx, publicPool.ID, publicPool.Balance+amount); err != nil {
			return errcode.ErrDatabase
		}
		// 写交易记录（唯一幂等键防止重复）
		idemKey := "recover_" + randomHex(16)
		ttx := &repository.Transaction{
			TxType:          "recover",
			DiscourseUserID: userID,
			Amount:          amount,
			Fee:             0,
			NetAmount:       amount,
			FromType:        "user",
			ToType:          "pool",
			IdempotencyKey:  idemKey,
			Description:     &reason,
			Status:          "completed",
		}
		if err := s.txRepo.Create(tx, ttx); err != nil {
			return errcode.ErrDatabase
		}
		// 审计日志（事务内写入，保证一致性）
		detailMsg := "recover " + fmt.Sprintf("%d", amount) + " CHB: " + reason
		detailJSON := fmt.Sprintf(`{"msg": "%s"}`, detailMsg)
		tx.Exec(
			"INSERT INTO audit_logs (operator_id, action, target_type, target_id, detail, ip_address, created_at) VALUES (?, ?, ?, ?, ?::jsonb, ?, NOW())",
			operatorID, "points_recover", "user", userID, detailJSON, "",
		)
		return nil
	})
}
func (s *AdminService) GetAppByID(id int64) (*repository.AppModel, error) {
	return s.appRepo.GetByID(id)
}
func (s *AdminService) AdjustPoints(userID int64, amount int64, direction string, reason string, operatorID int64) error {
	if amount <= 0 || reason == "" {
		return errcode.ErrParamInvalid
	}
	if direction != "add" && direction != "deduct" {
		return errcode.ErrParamInvalid
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		ub, err := s.balanceRepo.GetByUserIDWithLock(tx, userID)
		if err != nil {
			return errcode.ErrNotFound
		}
		if direction == "deduct" {
			if ub.Balance < amount {
				return errcode.ErrBalanceInsufficient
			}
			newBalance := ub.Balance - amount
			if err := s.balanceRepo.UpdateBalance(tx, userID, newBalance, ub.Version+1); err != nil {
				return err
			}
			publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
			if err != nil {
				return errcode.ErrDatabase
			}
			if err := s.poolRepo.UpdateBalanceTx(tx, publicPool.ID, publicPool.Balance+amount); err != nil {
				return errcode.ErrDatabase
			}
		} else {
			publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
			if err != nil {
				return errcode.ErrDatabase
			}
			if publicPool.Balance < amount {
				return errcode.ErrBalanceInsufficient
			}
			if err := s.poolRepo.UpdateBalanceTx(tx, publicPool.ID, publicPool.Balance-amount); err != nil {
				return errcode.ErrDatabase
			}
			newBalance := ub.Balance + amount
			newEarned := ub.TotalEarned + amount
			if err := s.balanceRepo.UpdateBalanceAndEarned(tx, userID, newBalance, ub.Version+1, newEarned); err != nil {
				return err
			}
		}
		idemKey := "adjust_" + randomHex(16)
		txType := "admin_adjust"
		fromType := "pool"
		toType := "user"
		if direction == "deduct" {
			fromType = "user"
			toType = "pool"
		}
		adjAmount := amount
		if direction == "deduct" {
			adjAmount = -amount
		}
		ttx := &repository.Transaction{
			TxType:          txType,
			DiscourseUserID: userID,
			Amount:          adjAmount,
			NetAmount:       amount,
			FromType:        fromType,
			ToType:          toType,
			IdempotencyKey:  idemKey,
			Description:     &reason,
			Status:          "completed",
		}
		if err := s.txRepo.Create(tx, ttx); err != nil {
			return errcode.ErrDatabase
		}
		action := "points_add"
		if direction == "deduct" {
			action = "points_deduct"
		}
		detailJSON := fmt.Sprintf("{\"msg\": \"%s %d CHB: %s\"}", direction, amount, reason)
		tx.Exec(
			"INSERT INTO audit_logs (operator_id, action, target_type, target_id, detail, ip_address, created_at) VALUES (?, ?, ?, ?, ?::jsonb, ?, NOW())",
			operatorID, action, "user", userID, detailJSON, "",
		)
		return nil
	})
}
