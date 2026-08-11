package service

import (
	"fmt"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type MarketplaceService struct {
	db           *gorm.DB
	marketRepo   *repository.MarketplaceRepo
	balanceRepo  *repository.UserBalanceRepo
	poolRepo     *repository.PoolRepo
	txRepo       *repository.TransactionRepo
}

func NewMarketplaceService(
	db *gorm.DB,
	marketRepo *repository.MarketplaceRepo,
	balanceRepo *repository.UserBalanceRepo,
	poolRepo *repository.PoolRepo,
	txRepo *repository.TransactionRepo,
) *MarketplaceService {
	return &MarketplaceService{
		db:          db,
		marketRepo:  marketRepo,
		balanceRepo: balanceRepo,
		poolRepo:    poolRepo,
		txRepo:      txRepo,
	}
}

// ===== Items =====

type ItemQuery struct {
	Category string
	Keyword  string
	Sort     string
	Page     int
	PageSize int
}

func (s *MarketplaceService) ListItems(q ItemQuery) ([]repository.MarketplaceItemModel, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return s.marketRepo.ListItems(q.Category, q.Keyword, q.Sort, q.Page, q.PageSize)
}

func (s *MarketplaceService) GetItem(id int64) (*repository.MarketplaceItemModel, error) {
	return s.marketRepo.GetItem(id)
}

func (s *MarketplaceService) CreateItem(sellerID int64, title, description string, price int64, stock int, category, imageURL string) (*repository.MarketplaceItemModel, error) {
	// 确保卖家余额账户存在（首次入驻自动建档）
	if err := s.ensureUser(sellerID); err != nil {
		return nil, err
	}
	item := &repository.MarketplaceItemModel{
		SellerID: sellerID,
		Title:    title,
		Price:    price,
		Stock:    stock,
		Status:   "pending",
		Category: &category,
	}
	if description != "" {
		item.Description = &description
	}
	if imageURL != "" {
		item.ImageURL = &imageURL
	}
	if category == "" {
		item.Category = nil
	}
	if err := s.marketRepo.CreateItem(item); err != nil {
		return nil, errcode.ErrDatabase
	}
	return item, nil
}

func (s *MarketplaceService) ListMyItems(sellerID int64, page, pageSize int) ([]repository.MarketplaceItemModel, int64, error) {
	return s.marketRepo.ListItemsBySeller(sellerID, page, pageSize)
}

// ===== Orders =====

type OrderResult struct {
	OrderNo     string `json:"order_no"`
	TotalAmount int64  `json:"total_amount"`
	Fee         int64  `json:"fee"`
	NetAmount   int64  `json:"net_amount"`
	Status      string `json:"status"`
}

func (s *MarketplaceService) CreateOrder(buyerID, itemID int64, quantity int, idempotencyKey string) (*OrderResult, error) {
	if quantity <= 0 {
		quantity = 1
	}
	// 确保买卖双方账户存在
	if err := s.ensureUser(buyerID); err != nil {
		return nil, err
	}

	existing, err := s.txRepo.GetByIDempotencyKey(idempotencyKey)
	if err == nil && existing != nil {
		return &OrderResult{Status: "duplicate"}, nil
	}

	var result *OrderResult
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Lock item
		item, err := s.marketRepo.GetItemWithLock(tx, itemID)
		if err != nil {
			return errcode.ErrNotFound
		}
		if item.Status != "approved" {
			return errcode.ErrItemPending
		}
		if item.Stock < quantity {
			return errcode.ErrItemSoldOut
		}

		// Lock buyer
		buyer, err := s.balanceRepo.GetByUserIDWithLock(tx, buyerID)
		if err != nil {
			return errcode.ErrBalanceInsufficient
		}
		if buyer.Status != "active" {
			return errcode.ErrAccountFrozen
		}

		totalAmount := item.Price * int64(quantity)
		fee := totalAmount * 10 / 100 // 10% fee
		netAmount := totalAmount - fee

		if buyer.Balance < totalAmount {
			return errcode.ErrBalanceInsufficient
		}

		// Lock seller
		seller, err := s.balanceRepo.GetByUserIDWithLock(tx, item.SellerID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 卖家账户首次出现时自动建档
				seller = &repository.UserBalance{
					DiscourseUserID: item.SellerID,
					Username:        fmt.Sprintf("user_%d", item.SellerID),
					Balance:         0,
					Version:         1,
					TrustLevel:      0,
					Status:          "active",
				}
				if err := s.balanceRepo.CreateWithTx(tx, seller); err != nil {
					return errcode.ErrDatabase
				}
			} else {
				return errcode.ErrDatabase
			}
		}

		// Lock official pool for fee
		officialPool, err := s.poolRepo.GetOfficialPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}

		// Deduct from buyer
		if err := s.balanceRepo.UpdateBalanceAndSpent(tx, buyerID, buyer.Balance-totalAmount, buyer.Version+1, buyer.TotalSpent+totalAmount); err != nil {
			return err
		}

		// Add net to seller
		if err := s.balanceRepo.UpdateBalance(tx, item.SellerID, seller.Balance+netAmount, seller.Version+1); err != nil {
			return err
		}

		// Add fee to official pool
		if err := s.poolRepo.UpdateBalance(officialPool.ID, officialPool.Balance+fee); err != nil {
			return errcode.ErrDatabase
		}

		// Update stock
		if err := s.marketRepo.UpdateItemStock(tx, itemID, item.Stock-quantity); err != nil {
			return errcode.ErrDatabase
		}

		// Create order
		orderNo := fmt.Sprintf("ORD%d%s", time.Now().UnixMilli(), idempotencyKey[len(idempotencyKey)-6:])
		order := &repository.MarketplaceOrderModel{
			OrderNo:     orderNo,
			ItemID:      itemID,
			BuyerID:     buyerID,
			SellerID:    item.SellerID,
			Quantity:    quantity,
			UnitPrice:   item.Price,
			TotalAmount: totalAmount,
			Fee:         fee,
			NetAmount:   netAmount,
			Status:      "completed",
		}
		if err := s.marketRepo.CreateOrder(tx, order); err != nil {
			return errcode.ErrDatabase
		}

		// Write transaction
		desc := fmt.Sprintf("购买商品: %s x%d", item.Title, quantity)
		txModel := &repository.Transaction{
			TxType:          "transfer",
			DiscourseUserID: buyerID,
			Amount:          totalAmount,
			Fee:             fee,
			NetAmount:       netAmount,
			FromType:        "user",
			ToType:          "user",
			FromID:          &buyerID,
			ToID:            &item.SellerID,
			IdempotencyKey:  idempotencyKey,
			Description:     &desc,
			Status:          "completed",
		}
		if err := s.txRepo.Create(tx, txModel); err != nil {
			return errcode.ErrDatabase
		}

		result = &OrderResult{
			OrderNo:     orderNo,
			TotalAmount: totalAmount,
			Fee:         fee,
			NetAmount:   netAmount,
			Status:      "completed",
		}
		return nil
	})

	return result, err
}

// ensureUser 确保用户余额账户存在（首次活跃自动建档），不存在则创建 0 余额账户。
func (s *MarketplaceService) ensureUser(userID int64) error {
	if userID <= 0 {
		return errcode.ErrParamInvalid
	}
	_, err := s.balanceRepo.GetByUserID(userID)
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return errcode.ErrDatabase
	}
	return s.balanceRepo.Create(&repository.UserBalance{
		DiscourseUserID: userID,
		Username:        fmt.Sprintf("user_%d", userID),
		Balance:         0,
		Version:         1,
		TrustLevel:      0,
		Status:          "active",
	})
}

func (s *MarketplaceService) ListOrders(userID int64, role, status string, page, pageSize int) ([]repository.MarketplaceOrderModel, int64, error) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 || pageSize > 100 { pageSize = 20 }
	return s.marketRepo.ListOrdersByUser(userID, role, status, page, pageSize)
}

// ===== Merchant Applications =====

func (s *MarketplaceService) ApplyMerchant(userID int64, shopName, description string) error {
	existing, err := s.marketRepo.GetApplicationByUser(userID)
	if err == nil && existing != nil {
		if existing.Status == "pending" {
			return errcode.ErrMerchantAppExists
		}
	}
	app := &repository.MerchantApplicationModel{
		DiscourseUserID: userID,
		ShopName:        shopName,
		Status:          "pending",
	}
	if description != "" {
		app.Description = &description
	}
	return s.marketRepo.CreateApplication(app)
}
