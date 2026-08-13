package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type DisputeService struct {
	db          *gorm.DB
	disputeRepo *repository.DisputeRepo
	orderRepo   *repository.MarketplaceRepo
	balanceRepo *repository.UserBalanceRepo
	txRepo      *repository.TransactionRepo
	poolRepo    *repository.PoolRepo
}

func NewDisputeService(db *gorm.DB, dr *repository.DisputeRepo, or *repository.MarketplaceRepo, br *repository.UserBalanceRepo, tr *repository.TransactionRepo, pr *repository.PoolRepo) *DisputeService {
	return &DisputeService{db: db, disputeRepo: dr, orderRepo: or, balanceRepo: br, txRepo: tr, poolRepo: pr}
}

type CreateDisputeReq struct {
	Reason string   `json:"reason"`
	Images []string `json:"images"`
}

type RejectDisputeReq struct {
	Reason string   `json:"reason"`
	Images []string `json:"images"`
}

func (s *DisputeService) CreateDispute(buyerID int64, orderID int64, req CreateDisputeReq) (*repository.DisputeModel, error) {
	// 获取订单
	order, err := s.orderRepo.GetOrderByID(orderID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if order.BuyerID != buyerID {
		return nil, errcode.ErrPermissionDenied
	}
	// 检查3天内可发起
	if time.Since(order.CreatedAt) > 72*time.Hour {
		return nil, errcode.ErrParamInvalid
	}
	// 检查是否已有争议
	existing, _ := s.disputeRepo.GetByOrderID(orderID)
	if existing != nil && existing.ID > 0 {
		return nil, errcode.ErrDuplicateRequest
	}

	var imagesJSON *string
	if len(req.Images) > 0 {
		b, _ := json.Marshal(req.Images)
		s := string(b)
		imagesJSON = &s
	}

	dispute := &repository.DisputeModel{
		OrderID:      orderID,
		BuyerID:      buyerID,
		SellerID:     order.SellerID,
		Status:       "pending",
		BuyerReason:  req.Reason,
		BuyerImages:  imagesJSON,
		RefundAmount: order.TotalAmount,
	}

	err = s.disputeRepo.Create(dispute)
	if err != nil {
		return nil, errcode.ErrInternal
	}

	// 更新订单争议状态
	s.db.Model(&repository.MarketplaceOrderModel{}).Where("id = ?", orderID).
		Update("dispute_status", "pending")

	return dispute, nil
}

func (s *DisputeService) AcceptDispute(sellerID int64, disputeID int64) error {
	d, err := s.disputeRepo.GetByID(disputeID)
	if err != nil || d.SellerID != sellerID {
		return errcode.ErrPermissionDenied
	}
	if d.Status != "pending" {
		return errcode.ErrParamInvalid
	}

	return s.doRefund(d, "accepted")
}

func (s *DisputeService) RejectDispute(sellerID int64, disputeID int64, req RejectDisputeReq) error {
	d, err := s.disputeRepo.GetByID(disputeID)
	if err != nil || d.SellerID != sellerID {
		return errcode.ErrPermissionDenied
	}
	if d.Status != "pending" {
		return errcode.ErrParamInvalid
	}

	var imagesJSON *string
	if len(req.Images) > 0 {
		b, _ := json.Marshal(req.Images)
		s := string(b)
		imagesJSON = &s
	}

	return s.disputeRepo.UpdateSellerResponse(disputeID, "reject", req.Reason, imagesJSON)
}

func (s *DisputeService) AutoRefundExpired() error {
	list, err := s.disputeRepo.ListExpiredPending()
	if err != nil {
		return err
	}
	for _, d := range list {
		s.doRefund(&d, "auto_refunded")
		s.disputeRepo.UpdateStatus(d.ID, "auto_refunded")
	}
	return nil
}

func (s *DisputeService) ReleaseSellerWin() error {
	list, err := s.disputeRepo.ListSellerWinExpired()
	if err != nil {
		return err
	}
	for _, d := range list {
		order, _ := s.orderRepo.GetOrderByID(d.OrderID)
		if order == nil {
			continue
		}
		tx := s.db.Begin()
		s.balanceRepo.TransferPendingToAvailable(tx, d.SellerID, order.NetAmount)
		tx.Commit()
		s.disputeRepo.UpdateStatus(d.ID, "closed")
	}
	return nil
}

func (s *DisputeService) AdminDecide(adminID int64, disputeID int64, decision, note string) error {
	d, err := s.disputeRepo.GetByID(disputeID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if d.Status != "rejected" {
		return errcode.ErrParamInvalid
	}

	err = s.disputeRepo.UpdateAdminDecision(disputeID, adminID, decision, note)
	if err != nil {
		return errcode.ErrInternal
	}

	if decision == "buyer_win" {
		return s.doRefund(d, "admin_buyer_win")
	}
	// seller_win: 48h后自动转可用，scheduler处理
	return nil
}

func (s *DisputeService) doRefund(d *repository.DisputeModel, newStatus string) error {
	order, err := s.orderRepo.GetOrderByID(d.OrderID)
	if err != nil || order == nil {
		return errcode.ErrNotFound
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 买家可用积分恢复全额(含手续费)
	buyer, err := s.balanceRepo.GetByUserIDWithLock(tx, d.BuyerID)
	if err != nil {
		tx.Rollback()
		return errcode.ErrInternal
	}
	s.balanceRepo.UpdateBalance(tx, d.BuyerID, buyer.Balance+d.RefundAmount, buyer.Version+1)

	// 卖家未来积分扣减净额
	s.balanceRepo.DeductPendingBalance(tx, d.SellerID, order.NetAmount)

	// 写入交易记录 - 买家退款
	buyerDesc := fmt.Sprintf("争议退款-订单#%d", d.OrderID)
	tx2 := &repository.Transaction{
		TxType:         "refund",
		DiscourseUserID: d.BuyerID,
		Amount:         d.RefundAmount,
		NetAmount:      d.RefundAmount,
		FromType:       "seller_pending",
		ToType:         "user",
		IdempotencyKey: fmt.Sprintf("refund_buyer_%d", d.ID),
		Description:    &buyerDesc,
		Status:         "completed",
	}
	s.txRepo.Create(tx, tx2)

	// 写入交易记录 - 卖家扣减
	sellerDesc := fmt.Sprintf("争议扣减-订单#%d", d.OrderID)
	tx3 := &repository.Transaction{
		TxType:         "refund_deduction",
		DiscourseUserID: d.SellerID,
		Amount:         -order.NetAmount,
		NetAmount:      -order.NetAmount,
		FromType:       "order",
		ToType:         "user",
		IdempotencyKey: fmt.Sprintf("refund_seller_%d", d.ID),
		Description:    &sellerDesc,
		Status:         "completed",
	}
	s.txRepo.Create(tx, tx3)

	// 更新争议状态
	now := time.Now()
	tx.Model(&repository.DisputeModel{}).Where("id = ?", d.ID).
		Updates(map[string]interface{}{"status": newStatus, "resolved_at": now})

	// 更新订单争议状态
	tx.Model(&repository.MarketplaceOrderModel{}).Where("id = ?", d.OrderID).
		Update("dispute_status", "closed")

	if err := tx.Commit().Error; err != nil {
		return errcode.ErrInternal
	}
	return nil
}

func (s *DisputeService) GetDispute(disputeID int64) (*repository.DisputeModel, error) {
	return s.disputeRepo.GetByID(disputeID)
}

func (s *DisputeService) ListMyDisputes(userID int64, role string, page, pageSize int) ([]repository.DisputeModel, int64, error) {
	return s.disputeRepo.ListByUser(userID, role, page, pageSize)
}

func (s *DisputeService) ListAllDisputes(page, pageSize int, status string) ([]DisputeDetail, int64, error) {
	list, total, err := s.disputeRepo.ListAll(page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}
	var result []DisputeDetail
	for _, d := range list {
		order, _ := s.orderRepo.GetOrderByID(d.OrderID)
		result = append(result, DisputeDetail{Dispute: d, ItemTitle: item_title})
	}
	return result, total, nil
}

type DisputeDetail struct {
	Dispute   repository.DisputeModel `json:"dispute"`
	ItemTitle string                  `json:"item_title"`
}
