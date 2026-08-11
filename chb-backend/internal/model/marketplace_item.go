package model

type MarketplaceItem struct {
    BaseModel
    SellerID    int64   `gorm:"not null" json:"seller_id"`
    Title       string  `gorm:"type:varchar(256);not null" json:"title"`
    Description *string `gorm:"type:text" json:"description"`
    Price       int64   `gorm:"not null" json:"price"`
    Stock       int     `gorm:"not null;default:0" json:"stock"`
    Status      string  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
    Category    *string `gorm:"type:varchar(64)" json:"category"`
    ImageURL    *string `gorm:"type:varchar(512)" json:"image_url"`
}
