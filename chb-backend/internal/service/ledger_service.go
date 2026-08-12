package service

import (
"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type LedgerService struct {
	db          *gorm.DB
	poolRepo    *repository.PoolRepo
	balanceRepo *repository.UserBalanceRepo
	txRepo      *repository.TransactionRepo
}

func NewLedgerService(db *gorm.DB, poolRepo *repository.PoolRepo, balanceRepo *repository.UserBalanceRepo, txRepo *repository.TransactionRepo) *LedgerService {
	return &LedgerService{
		db:          db,
		poolRepo:    poolRepo,
		balanceRepo: balanceRepo,
		txRepo:      txRepo,
	}
}

type BalanceInfo struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Balance     int64  `json:"balance"`
	TrustLevel  int16  `json:"trust_level"`
	Status      string `json:"status"`
	TotalEarned int64  `json:"total_earned"`
	TotalSpent  int64  `json:"total_spent"`
}

func (s *LedgerService) GetBalance(discourseUserID int64) (*BalanceInfo, error) {
	ub, err := s.balanceRepo.GetByUserID(discourseUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &BalanceInfo{
				UserID:     discourseUserID,
				Balance:    0,
				TrustLevel: 0,
				Status:     "active",
			}, nil
		}
		return nil, errcode.ErrDatabase
	}

	return &BalanceInfo{
		UserID:      ub.DiscourseUserID,
		Username:    ub.Username,
		Balance:     ub.Balance,
		TrustLevel:  ub.TrustLevel,
		Status:      ub.Status,
		TotalEarned: ub.TotalEarned,
		TotalSpent:  ub.TotalSpent,
	}, nil
}

type SpendResult struct {
	TransactionID int64  `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee"`
	NetAmount     int64  `json:"net_amount"`
	NewBalance    int64  `json:"new_balance"`
	Status        string `json:"status"`
}

func (s *LedgerService) Spend(discourseUserID, amount int64, idempotencyKey, description string, appFeeRate float64, boundUserID *int64) (*SpendResult, error) {
	if amount <= 0 {
		return nil, errcode.ErrParamInvalid
	}

	// Check idempotency
	existing, err := s.txRepo.GetByIDempotencyKey(idempotencyKey)
	if err == nil && existing != nil {
		return &SpendResult{
			TransactionID: existing.ID,
			Amount:        existing.Amount,
			Fee:           existing.Fee,
			NetAmount:     existing.NetAmount,
			NewBalance:    0,
			Status:        "duplicate",
		}, nil
	}

	var result *SpendResult
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Lock user balance
		ub, err := s.balanceRepo.GetByUserIDWithLock(tx, discourseUserID)
		if err != nil {
			return errcode.ErrBalanceInsufficient
		}

		if ub.Status != "active" {
			return errcode.ErrAccountFrozen
		}

		if ub.Balance < amount {
			return errcode.ErrBalanceInsufficient
		}

		// Calculate fee
		fee := int64(float64(amount) * appFeeRate / 100.0)
		netAmount := amount - fee
		newBalance := ub.Balance - amount

		// Deduct from user
		if err := s.balanceRepo.UpdateBalanceAndSpent(tx, discourseUserID, newBalance, ub.Version+1, ub.TotalSpent+amount); err != nil {
			return err
		}

		// Add fee to official pool
		officialPool, err := s.poolRepo.GetOfficialPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}
		if err := s.poolRepo.UpdateBalanceTx(tx, officialPool.ID, officialPool.Balance+fee); err != nil {
			return errcode.ErrDatabase
		}

		// If bound user, transfer net amount
		var toID *int64
		if boundUserID != nil {
			toID = boundUserID
			bu, err := s.balanceRepo.GetByUserIDWithLock(tx, *boundUserID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errcode.ErrNotFound
				}
				return errcode.ErrDatabase
			}
			if err := s.balanceRepo.UpdateBalance(tx, *boundUserID, bu.Balance+netAmount, bu.Version+1); err != nil {
				return err
			}
		}

		// Write transaction
		txModel := &repository.Transaction{
			TxType:          "spend",
			DiscourseUserID: discourseUserID,
			Amount:          amount,
			Fee:             fee,
			NetAmount:       netAmount,
			FromType:        "user",
			ToType:          "app",
			ToID:            toID,
			IdempotencyKey:  idempotencyKey,
			Description:     &description,
			Status:          "completed",
		}
		if err := s.txRepo.Create(tx, txModel); err != nil {
			return errcode.ErrDatabase
		}

		result = &SpendResult{
			TransactionID: txModel.ID,
			Amount:        amount,
			Fee:           fee,
			NetAmount:     netAmount,
			NewBalance:    newBalance,
			Status:        "completed",
		}
		return nil
	})

	return result, err
}

type RewardResult struct {
	Amount      int64  `json:"amount"`
	TrustLevel  int16  `json:"trust_level"`
	Multiplier  float64 `json:"multiplier"`
	FinalAmount int64  `json:"final_amount"`
	EarnedToday int64  `json:"earned_today"`
	DailyCap    int64  `json:"daily_cap"`
	Status      string `json:"status"`
}

func (s *LedgerService) Reward(discourseUserID int64, action string, refType string, refID int64, idempotencyKey string) (*RewardResult, error) {
	// Check idempotency
	existing, err := s.txRepo.GetByIDempotencyKey(idempotencyKey)
	if err == nil && existing != nil {
		return &RewardResult{Status: "duplicate"}, nil
	}

	var result *RewardResult
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Get reward rule
		// TODO: load from reward service
		_ = tx
		result = &RewardResult{
			Amount:      100,
			TrustLevel:  0,
			Multiplier:  1.0,
			FinalAmount: 100,
			EarnedToday: 100,
			DailyCap:    1000,
			Status:      "completed",
		}
		return nil
	})

	return result, err
}

type PoolInfo struct {
	PublicPool  *PoolDetail `json:"public_pool"`
	OfficialPool *PoolDetail `json:"official_pool"`
}

type PoolDetail struct {
	TotalSupply int64   `json:"total_supply"`
	Balance     int64   `json:"balance"`
	WaterLevel  float64 `json:"water_level,omitempty"`
}

func (s *LedgerService) GetPools() (*PoolInfo, error) {
	publicPool, err := s.poolRepo.GetPublicPool()
	if err != nil {
		return nil, errcode.ErrDatabase
	}
	officialPool, err := s.poolRepo.GetOfficialPool()
	if err != nil {
		return nil, errcode.ErrDatabase
	}

	waterLevel := float64(publicPool.Balance) / float64(publicPool.TotalSupply)

	return &PoolInfo{
		PublicPool: &PoolDetail{
			TotalSupply: publicPool.TotalSupply,
			Balance:     publicPool.Balance,
			WaterLevel:  waterLevel,
		},
		OfficialPool: &PoolDetail{
			Balance: officialPool.Balance,
		},
	}, nil
}

type ReleaseResult struct {
	ReleaseID          int64  `json:"release_id"`
	Amount             int64  `json:"amount"`
	NewPublicBalance   int64  `json:"new_public_balance"`
	NewOfficialBalance int64  `json:"new_official_balance"`
}

func (s *LedgerService) Release(amount int64, reason string, operatorID int64) (*ReleaseResult, error) {
	if amount <= 0 {
		return nil, errcode.ErrParamInvalid
	}

	var result *ReleaseResult
	err := s.db.Transaction(func(tx *gorm.DB) error {
		officialPool, err := s.poolRepo.GetOfficialPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}
		if officialPool.Balance < amount {
			return errcode.ErrBalanceInsufficient
		}

		publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}

		if err := s.poolRepo.UpdateBalanceTx(tx, officialPool.ID, officialPool.Balance-amount); err != nil {
			return errcode.ErrDatabase
		}
		if err := s.poolRepo.UpdateBalanceTx(tx, publicPool.ID, publicPool.Balance+amount); err != nil {
			return errcode.ErrDatabase
		}

		now := time.Now()
		releaseLog := &repository.ReleaseLog{
			Amount:     amount,
			FromPool:   "official",
			ToPool:     "public",
			OperatorID: operatorID,
			Reason:     &reason,
			IsAuto:     false,
			CreatedAt:  now,
		}
		if err := releaseLog.Create(tx); err != nil {
			return errcode.ErrDatabase
		}

		result = &ReleaseResult{
			ReleaseID:          releaseLog.ID,
			Amount:             amount,
			NewPublicBalance:   publicPool.Balance + amount,
			NewOfficialBalance: officialPool.Balance - amount,
		}
		return nil
	})

	return result, err
}

type AuditInfo struct {
	TotalSupply      int64  `json:"total_supply"`
	UserBalancesSum  int64  `json:"user_balances_sum"`
	PublicPoolBalance int64 `json:"public_pool_balance"`
	OfficialPoolBalance int64 `json:"official_pool_balance"`
	IsBalanced       bool   `json:"is_balanced"`
}

func (s *LedgerService) Audit() (*AuditInfo, error) {
	publicPool, err := s.poolRepo.GetPublicPool()
	if err != nil {
		return nil, errcode.ErrDatabase
	}
	officialPool, err := s.poolRepo.GetOfficialPool()
	if err != nil {
		return nil, errcode.ErrDatabase
	}
	userSum, err := s.balanceRepo.GetTotalBalance()
	if err != nil {
		return nil, errcode.ErrDatabase
	}

	total := userSum + publicPool.Balance + officialPool.Balance
	isBalanced := total == publicPool.TotalSupply

	return &AuditInfo{
		TotalSupply:        publicPool.TotalSupply,
		UserBalancesSum:    userSum,
		PublicPoolBalance:  publicPool.Balance,
		OfficialPoolBalance: officialPool.Balance,
		IsBalanced:         isBalanced,
	}, nil
}

// TransactionQuery for listing transactions
type TransactionQuery struct {
	TxType    string
	StartDate string
	EndDate   string
	Page      int
	PageSize  int
}

func (s *LedgerService) ListTransactions(discourseUserID int64, q TransactionQuery) ([]repository.Transaction, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return s.txRepo.ListByUser(discourseUserID, q.TxType, q.StartDate, q.EndDate, q.Page, q.PageSize)
}
