package model

type MarketplaceOrder struct {
    BaseModel
    OrderNo     string `gorm:"type:varchar(64);uniqueIndex;not null" json:"order_no"`
    ItemID      int64  `gorm:"not null" json:"item_id"`
    BuyerID     int64  `gorm:"not null" json:"buyer_id"`
    SellerID    int64  `gorm:"not null" json:"seller_id"`
    Quantity    int    `gorm:"not null;default:1" json:"quantity"`
    UnitPrice   int64  `gorm:"not null" json:"unit_price"`
    TotalAmount int64  `gorm:"not null" json:"total_amount"`
    Fee         int64  `gorm:"not null" json:"fee"`
    NetAmount   int64  `gorm:"not null" json:"net_amount"`
    Status      string `gorm:"type:varchar(20);not null;default:pending" json:"status"`
}
