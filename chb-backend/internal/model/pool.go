package model

type Pool struct {
    BaseModel
    PoolType    string `gorm:"type:varchar(32);uniqueIndex;not null" json:"pool_type"`
    TotalSupply int64  `gorm:"not null" json:"total_supply"`
    Balance     int64  `gorm:"not null" json:"balance"`
}
