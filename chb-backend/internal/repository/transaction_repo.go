package repository

import (
	"time"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(tx *gorm.DB, t *Transaction) error {
	return tx.Create(t).Error
}

func (r *TransactionRepo) GetByIDempotencyKey(key string) (*Transaction, error) {
	var t Transaction
	err := r.db.Where("idempotency_key = ?", key).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepo) ListByUser(discourseUserID int64, txType string, startDate, endDate string, page, pageSize int) ([]Transaction, int64, error) {
	var list []Transaction
	var total int64
	query := r.db.Model(&Transaction{}).Where("discourse_user_id = ?", discourseUserID)
	if txType != "" {
		query = query.Where("tx_type = ?", txType)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate+"T00:00:00+08:00")
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+"T23:59:59+08:00")
	}
	query.Count(&total)
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *TransactionRepo) GetTotalCount() (int64, error) {
	var total int64
	err := r.db.Model(&Transaction{}).Count(&total).Error
	return total, err
}

func (r *TransactionRepo) GetTodayCount() (int64, error) {
	var total int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&Transaction{}).
		Where("created_at >= ?", today+"T00:00:00+08:00").
		Count(&total).Error
	return total, err
}

// Transaction shadow model
type Transaction struct {
	ID              int64     `json:"id"`
	TxType          string    `json:"tx_type"`
	DiscourseUserID int64     `json:"discourse_user_id"`
	Amount          int64     `json:"amount"`
	Fee             int64     `json:"fee"`
	NetAmount       int64     `json:"net_amount"`
	FromType        string    `json:"from_type"`
	ToType          string    `json:"to_type"`
	FromID          *int64    `json:"from_id"`
	ToID            *int64    `json:"to_id"`
	IdempotencyKey  string    `json:"idempotency_key"`
	RefType         *string   `json:"ref_type"`
	RefID           *int64    `json:"ref_id"`
	Description     *string   `json:"description"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
