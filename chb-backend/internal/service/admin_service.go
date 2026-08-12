package service

import (
	"crypto/rand"
	"encoding/hex"
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

// ===== System Config =====

type SystemSetting struct {
	MarketplaceFeeRate float64 `json:"marketplace_fee_rate"`
	AutoReleaseEnabled bool    `json:"auto_release_enabled"`
	AutoReleaseThreshold int   `json:"auto_release_threshold"`
	AutoReleaseRatio   int     `json:"auto_release_ratio"`
	AutoReleaseMonthlyCap int64 `json:"auto_release_monthly_cap"`
}

func (s *AdminService) GetSettings() *SystemSetting {
	return &SystemSetting{
		MarketplaceFeeRate:   10.0,
		AutoReleaseEnabled:   false,
		AutoReleaseThreshold: 80,
		AutoReleaseRatio:     50,
		AutoReleaseMonthlyCap: 10000000,
	}
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
	return s.appRepo.Create(app)
}

func (s *AdminService) UpdateApp(app *repository.AppModel) error {
	return s.appRepo.Update(app)
}

func (s *AdminService) DeleteApp(id int64) error {
	return s.appRepo.Delete(id)
}

// ===== Marketplace Review =====

func (s *AdminService) ListPendingApplications(page, pageSize int) ([]repository.MerchantApplicationModel, int64, error) {
	return s.marketRepo.ListPendingApplications(page, pageSize)
}

func (s *AdminService) ReviewApplication(id int64, status, comment string, reviewerID int64) error {
	app, err := s.marketRepo.GetItem(id)
	if err != nil {
		return errcode.ErrNotFound
	}
	_ = app
	update := &repository.MerchantApplicationModel{
		ID:            id,
		Status:        status,
		ReviewedBy:    &reviewerID,
		ReviewComment: &comment,
	}
	return s.marketRepo.UpdateApplication(update)
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
	return s.marketRepo.UpdateItem(item)
}

// ===== User Management =====

func (s *AdminService) ListUsers(page, pageSize int, keyword string) ([]repository.UserBalance, int64, error) {
	// For now, just use the balance repo list
	return s.balanceRepo.List(page, pageSize)
}

func (s *AdminService) FreezeUser(userID int64) error {
	return s.balanceRepo.SetStatus(s.db, userID, "frozen")
}

func (s *AdminService) UnfreezeUser(userID int64) error {
	return s.balanceRepo.SetStatus(s.db, userID, "active")
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
	// TODO: implement audit log query
	return []AuditLogEntry{}, 0, nil
}

// ===== User Recover =====

func (s *AdminService) RecoverPoints(userID int64, amount int64, reason string, operatorID int64) error {
	// TODO: implement point recovery
	return nil
}
