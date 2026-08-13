package repository

import (
	"time"

	"gorm.io/gorm"
)

type DisputeModel struct {
	ID              int64      `json:"id"`
	OrderID         int64      `json:"order_id"`
	BuyerID         int64      `json:"buyer_id"`
	SellerID        int64      `json:"seller_id"`
	Status          string     `json:"status"`
	BuyerReason     string     `json:"buyer_reason"`
	BuyerImages     *string    `json:"buyer_images"`
	CreatedAt       time.Time  `json:"created_at"`
	SellerAction    *string    `json:"seller_action"`
	SellerReason    *string    `json:"seller_reason"`
	SellerImages    *string    `json:"seller_images"`
	SellerRespondedAt *time.Time `json:"seller_responded_at"`
	AutoRefunded    bool       `json:"auto_refunded"`
	AdminID         *int64     `json:"admin_id"`
	AdminDecision   *string    `json:"admin_decision"`
	AdminNote       *string    `json:"admin_note"`
	ResolvedAt      *time.Time `json:"resolved_at"`
	RefundAmount    int64      `json:"refund_amount"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (DisputeModel) TableName() string {
	return "disputes"
}

type DisputeRepo struct {
	db *gorm.DB
}

func NewDisputeRepo(db *gorm.DB) *DisputeRepo {
	return &DisputeRepo{db: db}
}

func (r *DisputeRepo) Create(d *DisputeModel) error {
	return r.db.Create(d).Error
}

func (r *DisputeRepo) GetByID(id int64) (*DisputeModel, error) {
	var d DisputeModel
	err := r.db.Where("id = ?", id).First(&d).Error
	return &d, err
}

func (r *DisputeRepo) GetByOrderID(orderID int64) (*DisputeModel, error) {
	var d DisputeModel
	err := r.db.Where("order_id = ?", orderID).Order("created_at DESC").First(&d).Error
	return &d, err
}

func (r *DisputeRepo) ListByBuyer(buyerID int64, page, pageSize int) ([]DisputeModel, int64, error) {
	var list []DisputeModel
	var total int64
	q := r.db.Model(&DisputeModel{}).Where("buyer_id = ?", buyerID)
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *DisputeRepo) ListBySeller(sellerID int64, page, pageSize int) ([]DisputeModel, int64, error) {
	var list []DisputeModel
	var total int64
	q := r.db.Model(&DisputeModel{}).Where("seller_id = ?", sellerID)
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *DisputeRepo) ListByUser(userID int64, role string, page, pageSize int) ([]DisputeModel, int64, error) {
	if role == "buyer" {
		return r.ListByBuyer(userID, page, pageSize)
	}
	return r.ListBySeller(userID, page, pageSize)
}

func (r *DisputeRepo) ListAll(page, pageSize int, status string) ([]DisputeModel, int64, error) {
	var list []DisputeModel
	var total int64
	q := r.db.Model(&DisputeModel{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *DisputeRepo) ListPendingBySeller(sellerID int64) ([]DisputeModel, error) {
	var list []DisputeModel
	err := r.db.Where("seller_id = ? AND status = ?", sellerID, "pending").Find(&list).Error
	return list, err
}

func (r *DisputeRepo) CountPendingBySeller(sellerID int64) (int64, error) {
	var count int64
	err := r.db.Model(&DisputeModel{}).Where("seller_id = ? AND status = ?", sellerID, "pending").Count(&count).Error
	return count, err
}

func (r *DisputeRepo) CountByBuyer(buyerID int64) (int64, error) {
	var count int64
	err := r.db.Model(&DisputeModel{}).Where("buyer_id = ?", buyerID).Count(&count).Error
	return count, err
}

func (r *DisputeRepo) UpdateStatus(id int64, status string) error {
	return r.db.Model(&DisputeModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *DisputeRepo) UpdateStatusWithResolve(id int64, status string, resolvedAt time.Time) error {
	return r.db.Model(&DisputeModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "resolved_at": resolvedAt, "updated_at": time.Now()}).Error
}

func (r *DisputeRepo) ListExpiredPending() ([]DisputeModel, error) {
	var list []DisputeModel
	err := r.db.Where("status = ? AND created_at <= ?", "pending", time.Now().Add(-72*time.Hour)).Find(&list).Error
	return list, err
}

func (r *DisputeRepo) ListSellerWinExpired() ([]DisputeModel, error) {
	var list []DisputeModel
	err := r.db.Where("status = ? AND resolved_at <= ?", "admin_seller_win", time.Now().Add(-48*time.Hour)).Find(&list).Error
	return list, err
}

func (r *DisputeRepo) UpdateSellerResponse(id int64, action, reason string, images *string) error {
	now := time.Now()
	statusValue := "rejected"
	if action == "accept" {
		statusValue = "accepted"
	}
	return r.db.Model(&DisputeModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"seller_action":      action,
			"seller_reason":      reason,
			"seller_images":      images,
			"seller_responded_at": now,
			"status":             statusValue,
			"updated_at":         now,
		}).Error
}

func (r *DisputeRepo) UpdateAdminDecision(id int64, adminID int64, decision, note string) error {
	now := time.Now()
	status := "admin_seller_win"
	if decision == "buyer_win" {
		status = "admin_buyer_win"
	}
	return r.db.Model(&DisputeModel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"admin_id":       adminID,
			"admin_decision": decision,
			"admin_note":     note,
			"status":         status,
			"resolved_at":    now,
			"updated_at":     now,
		}).Error
}
