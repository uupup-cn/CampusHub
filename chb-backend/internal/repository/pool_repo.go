package repository

import (
	"gorm.io/gorm"
)

type PoolRepo struct {
	db *gorm.DB
}

func NewPoolRepo(db *gorm.DB) *PoolRepo {
	return &PoolRepo{db: db}
}

func (r *PoolRepo) GetByType(poolType string) (*Pool, error) {
	var p Pool
	err := r.db.Where("pool_type = ?", poolType).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateBalance updates pool balance using the repo default connection (outside transactions).
func (r *PoolRepo) UpdateBalance(poolID, newBalance int64) error {
	return r.db.Model(&Pool{}).Where("id = ?", poolID).Update("balance", newBalance).Error
}

// UpdateBalanceTx updates pool balance within a transaction, avoiding lock conflicts
// when the row was previously locked with FOR UPDATE in the same transaction.
func (r *PoolRepo) UpdateBalanceTx(tx *gorm.DB, poolID, newBalance int64) error {
	return tx.Model(&Pool{}).Where("id = ?", poolID).Update("balance", newBalance).Error
}

func (r *PoolRepo) GetPublicPool() (*Pool, error) {
	return r.GetByType("public")
}

func (r *PoolRepo) GetOfficialPool() (*Pool, error) {
	return r.GetByType("official")
}

func (r *PoolRepo) GetPublicPoolWithLock(tx *gorm.DB) (*Pool, error) {
	var p Pool
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("pool_type = ?", "public").First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PoolRepo) GetOfficialPoolWithLock(tx *gorm.DB) (*Pool, error) {
	var p Pool
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("pool_type = ?", "official").First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Pool is a shadow model for repository queries
type Pool struct {
	ID          int64
	PoolType    string
	TotalSupply int64
	Balance     int64
}
