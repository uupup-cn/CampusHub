package model

type MerchantApplication struct {
    BaseModel
    DiscourseUserID int64   `gorm:"not null" json:"discourse_user_id"`
    ShopName        string  `gorm:"type:varchar(128);not null" json:"shop_name"`
    Description     *string `gorm:"type:text" json:"description"`
    Status          string  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
    ReviewedBy      *int64  `json:"reviewed_by"`
    ReviewComment   *string `gorm:"type:text" json:"review_comment"`
}
