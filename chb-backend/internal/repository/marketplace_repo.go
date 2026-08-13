package repository

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceItemModel struct {
	ID          int64     `json:"id"`
	SellerID    int64     `json:"seller_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Price       int64     `json:"price"`
	Stock       int       `json:"stock"`
	Status      string    `json:"status"`
	Category    *string   `json:"category"`
	ImageURL    *string   `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MarketplaceItemModel) TableName() string {
	return "marketplace_items"
}

type MarketplaceOrderModel struct {
	ID          int64     `json:"id"`
	OrderNo     string    `json:"order_no"`
	ItemID      int64     `json:"item_id"`
	BuyerID     int64     `json:"buyer_id"`
	SellerID    int64     `json:"seller_id"`
	Quantity    int       `json:"quantity"`
	UnitPrice   int64     `json:"unit_price"`
	TotalAmount int64     `json:"total_amount"`
	Fee         int64     `json:"fee"`
	NetAmount   int64     `json:"net_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MarketplaceOrderModel) TableName() string {
	return "marketplace_orders"
}

type MerchantApplicationModel struct {
	ID              int64     `json:"id"`
	DiscourseUserID int64     `json:"discourse_user_id"`
	ShopName        string    `json:"shop_name"`
	Description     *string   `json:"description"`
	Status          string    `json:"status"`
	ReviewedBy      *int64    `json:"reviewed_by"`
	ReviewComment   *string   `json:"review_comment"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (MerchantApplicationModel) TableName() string {
	return "merchant_applications"
}

type MarketplaceRepo struct {
	db *gorm.DB
}

func NewMarketplaceRepo(db *gorm.DB) *MarketplaceRepo {
	return &MarketplaceRepo{db: db}
}

// ===== Items =====

func (r *MarketplaceRepo) ListItems(category, keyword, sort string, page, pageSize int) ([]MarketplaceItemModel, int64, error) {
	var list []MarketplaceItemModel
	var total int64
	query := r.db.Model(&MarketplaceItemModel{}).Where("status = ?", "approved")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	order := "created_at DESC"
	switch sort {
	case "price_asc":
		order = "price ASC"
	case "price_desc":
		order = "price DESC"
	case "newest":
		order = "created_at DESC"
	}
	err := query.Order(order).Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MarketplaceRepo) GetItem(id int64) (*MarketplaceItemModel, error) {
	var item MarketplaceItemModel
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketplaceRepo) GetItemWithLock(tx *gorm.DB, id int64) (*MarketplaceItemModel, error) {
	var item MarketplaceItemModel
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketplaceRepo) CreateItem(item *MarketplaceItemModel) error {
	return r.db.Create(item).Error
}

func (r *MarketplaceRepo) UpdateItem(item *MarketplaceItemModel) error {
	return r.db.Model(&MarketplaceItemModel{}).Where("id = ?", item.ID).Updates(item).Error
}

func (r *MarketplaceRepo) UpdateItemStock(tx *gorm.DB, itemID int64, newStock int) error {
	return tx.Model(&MarketplaceItemModel{}).Where("id = ?", itemID).Update("stock", newStock).Error
}

func (r *MarketplaceRepo) ListItemsBySeller(sellerID int64, page, pageSize int) ([]MarketplaceItemModel, int64, error) {
	var list []MarketplaceItemModel
	var total int64
	r.db.Model(&MarketplaceItemModel{}).Where("seller_id = ?", sellerID).Count(&total)
	err := r.db.Where("seller_id = ?", sellerID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MarketplaceRepo) ListPendingItems(page, pageSize int) ([]MarketplaceItemModel, int64, error) {
	var list []MarketplaceItemModel
	var total int64
	r.db.Model(&MarketplaceItemModel{}).Where("status = ?", "pending").Count(&total)
	err := r.db.Where("status = ?", "pending").Order("created_at ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// ===== Orders =====

func (r *MarketplaceRepo) CreateOrder(tx *gorm.DB, order *MarketplaceOrderModel) error {
	return tx.Create(order).Error
}

func (r *MarketplaceRepo) GetOrder(orderNo string) (*MarketplaceOrderModel, error) {
	var order MarketplaceOrderModel
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *MarketplaceRepo) GetOrderByID(id int64) (*MarketplaceOrderModel, error) {
	var order MarketplaceOrderModel
	err := r.db.Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *MarketplaceRepo) ListOrdersByUser(userID int64, role, status string, page, pageSize int) ([]MarketplaceOrderModel, int64, error) {
	var list []MarketplaceOrderModel
	var total int64
	query := r.db.Model(&MarketplaceOrderModel{})
	if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	} else {
		query = query.Where("buyer_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// ===== Merchant Applications =====

func (r *MarketplaceRepo) CreateApplication(app *MerchantApplicationModel) error {
	return r.db.Create(app).Error
}

func (r *MarketplaceRepo) GetApplicationByUser(userID int64) (*MerchantApplicationModel, error) {
	var app MerchantApplicationModel
	err := r.db.Where("discourse_user_id = ?", userID).Last(&app).Error
	return &app, err
}

func (r *MarketplaceRepo) ListApplications(page, pageSize int) ([]MerchantApplicationModel, int64, error) {
	var list []MerchantApplicationModel
	var total int64
	r.db.Model(&MerchantApplicationModel{}).Count(&total)
	err := r.db.Order("created_at ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MarketplaceRepo) ListPendingApplications(page, pageSize int) ([]MerchantApplicationModel, int64, error) {
	var list []MerchantApplicationModel
	var total int64
	r.db.Model(&MerchantApplicationModel{}).Where("status = ?", "pending").Count(&total)
	err := r.db.Where("status = ?", "pending").Order("created_at ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *MarketplaceRepo) UpdateApplication(app *MerchantApplicationModel) error {
	return r.db.Model(&MerchantApplicationModel{}).Where("id = ?", app.ID).Updates(app).Error
}
