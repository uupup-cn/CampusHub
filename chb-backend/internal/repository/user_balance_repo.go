package repository

import (
	"gorm.io/gorm"
)

type UserBalanceRepo struct {
	db *gorm.DB
}

func NewUserBalanceRepo(db *gorm.DB) *UserBalanceRepo {
	return &UserBalanceRepo{db: db}
}

func (r *UserBalanceRepo) GetByUserID(discourseUserID int64) (*UserBalance, error) {
	var ub UserBalance
	err := r.db.Where("discourse_user_id = ?", discourseUserID).First(&ub).Error
	if err != nil {
		return nil, err
	}
	return &ub, nil
}

func (r *UserBalanceRepo) GetByUserIDWithLock(tx *gorm.DB, discourseUserID int64) (*UserBalance, error) {
	var ub UserBalance
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("discourse_user_id = ?", discourseUserID).First(&ub).Error
	if err != nil {
		return nil, err
	}
	return &ub, nil
}

func (r *UserBalanceRepo) Create(user *UserBalance) error {
	return r.db.Create(user).Error
}

func (r *UserBalanceRepo) CreateWithTx(tx *gorm.DB, user *UserBalance) error {
	return tx.Create(user).Error
}

func (r *UserBalanceRepo) UpdateBalance(tx *gorm.DB, discourseUserID, newBalance, newVersion int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ? AND version = ?", discourseUserID, newVersion-1).
		Updates(map[string]interface{}{
			"balance": newBalance,
			"version": newVersion,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

func (r *UserBalanceRepo) UpdateBalanceAndEarned(tx *gorm.DB, discourseUserID, newBalance, newVersion, totalEarned int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ? AND version = ?", discourseUserID, newVersion-1).
		Updates(map[string]interface{}{
			"balance":      newBalance,
			"version":      newVersion,
			"total_earned": totalEarned,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

func (r *UserBalanceRepo) UpdateBalanceAndSpent(tx *gorm.DB, discourseUserID, newBalance, newVersion, totalSpent int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ? AND version = ?", discourseUserID, newVersion-1).
		Updates(map[string]interface{}{
			"balance":     newBalance,
			"version":     newVersion,
			"total_spent": totalSpent,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

func (r *UserBalanceRepo) UpdateTrustLevel(tx *gorm.DB, discourseUserID int64, trustLevel int16) error {
	return tx.Model(&UserBalance{}).Where("discourse_user_id = ?", discourseUserID).
		Update("trust_level", trustLevel).Error
}

func (r *UserBalanceRepo) GetTotalBalance() (int64, error) {
	var sum struct{ Total int64 }
	err := r.db.Model(&UserBalance{}).Select("COALESCE(SUM(balance), 0) as total").Scan(&sum).Error
	return sum.Total, err
}

func (r *UserBalanceRepo) List(page, pageSize int) ([]UserBalance, int64, error) {
	var list []UserBalance
	var total int64
	r.db.Model(&UserBalance{}).Count(&total)
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *UserBalanceRepo) SetStatus(tx *gorm.DB, discourseUserID int64, status string) error {
	return tx.Model(&UserBalance{}).Where("discourse_user_id = ?", discourseUserID).
		Update("status", status).Error
}

// UserBalance shadow model
type UserBalance struct {
	ID              int64  `json:"id"`
	DiscourseUserID int64  `json:"discourse_user_id"`
	Username        string `json:"username"`
	Balance         int64  `json:"balance"`
	PendingBalance  int64  `json:"pending_balance"`
	Version         int64  `json:"version"`
	TrustLevel      int16  `json:"trust_level"`
	TotalEarned     int64  `json:"total_earned"`
	TotalSpent      int64  `json:"total_spent"`
	Status          string `json:"status"`
}


func (r *UserBalanceRepo) AddPendingBalance(tx *gorm.DB, discourseUserID int64, amount int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ?", discourseUserID).
		Update("pending_balance", gorm.Expr("pending_balance + ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

func (r *UserBalanceRepo) DeductPendingBalance(tx *gorm.DB, discourseUserID int64, amount int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ? AND pending_balance >= ?", discourseUserID, amount).
		Update("pending_balance", gorm.Expr("pending_balance - ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

func (r *UserBalanceRepo) TransferPendingToAvailable(tx *gorm.DB, discourseUserID int64, amount int64) error {
	result := tx.Model(&UserBalance{}).
		Where("discourse_user_id = ? AND pending_balance >= ?", discourseUserID, amount).
		Updates(map[string]interface{}{
			"balance":         gorm.Expr("balance + ?", amount),
			"pending_balance": gorm.Expr("pending_balance - ?", amount),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOptimisticLock
	}
	return nil
}

var ErrOptimisticLock = &OptimisticLockError{}

type OptimisticLockError struct{}

func (e *OptimisticLockError) Error() string {
	return "optimistic lock conflict"
}
