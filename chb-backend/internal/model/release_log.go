package model

type ReleaseLog struct {
    BaseModel
    Amount     int64   `gorm:"not null" json:"amount"`
    FromPool   string  `gorm:"type:varchar(20);not null" json:"from_pool"`
    ToPool     string  `gorm:"type:varchar(20);not null" json:"to_pool"`
    OperatorID int64   `gorm:"not null" json:"operator_id"`
    Reason     *string `gorm:"type:varchar(512)" json:"reason"`
    IsAuto     bool    `gorm:"not null;default:false" json:"is_auto"`
}
