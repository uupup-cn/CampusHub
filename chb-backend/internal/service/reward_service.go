package service

import (
	"fmt"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type RewardService struct {
	db         *gorm.DB
	rewardRepo *repository.RewardRepo
	poolRepo   *repository.PoolRepo
	balanceRepo *repository.UserBalanceRepo
	txRepo     *repository.TransactionRepo
	ledgerSvc  *LedgerService
}

func NewRewardService(
	db *gorm.DB,
	rewardRepo *repository.RewardRepo,
	poolRepo *repository.PoolRepo,
	balanceRepo *repository.UserBalanceRepo,
	txRepo *repository.TransactionRepo,
	ledgerSvc *LedgerService,
) *RewardService {
	return &RewardService{
		db:          db,
		rewardRepo:  rewardRepo,
		poolRepo:    poolRepo,
		balanceRepo: balanceRepo,
		txRepo:      txRepo,
		ledgerSvc:   ledgerSvc,
	}
}

type RewardRequest struct {
	Action          string
	DiscourseUserID int64
	RefType         string
	RefID           int64
	IdempotencyKey  string
	IPAddress       string
}

type RewardResponse struct {
	Amount      int64   `json:"amount"`
	TrustLevel  int16   `json:"trust_level"`
	Multiplier  float64 `json:"multiplier"`
	FinalAmount int64   `json:"final_amount"`
	EarnedToday int64   `json:"earned_today"`
	DailyCap    int64   `json:"daily_cap"`
	Status      string  `json:"status"`
	RejectReason string `json:"reject_reason,omitempty"`
}

func (s *RewardService) Grant(req *RewardRequest) (*RewardResponse, error) {
	// 1. Idempotency check
	existingLog, err := s.rewardRepo.GetLogByRef(req.RefType, req.RefID, req.DiscourseUserID)
	if err == nil && existingLog != nil {
		return &RewardResponse{Status: "duplicate", RejectReason: "duplicate"}, nil
	}

	// 2. Check idempotency key in transactions
	existingTx, err := s.txRepo.GetByIDempotencyKey(req.IdempotencyKey)
	if err == nil && existingTx != nil {
		return &RewardResponse{Status: "duplicate", RejectReason: "duplicate"}, nil
	}

	var resp *RewardResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 3. Get reward rule
		rule, err := s.rewardRepo.GetRuleByAction(req.Action)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errcode.ErrNotFound
			}
			return errcode.ErrDatabase
		}

		// 4. Get or create user balance
		ub, err := s.balanceRepo.GetByUserIDWithLock(tx, req.DiscourseUserID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 首次活跃的论坛用户自动建档，余额从 0 开始
				ub = &repository.UserBalance{
					DiscourseUserID: req.DiscourseUserID,
					Username:        fmt.Sprintf("user_%d", req.DiscourseUserID),
					Balance:         0,
					Version:         1,
					TrustLevel:      0,
					Status:          "active",
				}
				if err := s.balanceRepo.CreateWithTx(tx, ub); err != nil {
					return errcode.ErrDatabase
				}
			} else {
				return errcode.ErrDatabase
			}
		}

		// 5. Check account status
		if ub.Status != "active" {
			return errcode.ErrAccountFrozen
		}

		// 6. Get trust level cap
		cap, err := s.rewardRepo.GetCapByLevel(ub.TrustLevel)
		if err != nil {
			return errcode.ErrDatabase
		}

		// 7. Get today's quota
		today := time.Now().Format("2006-01-02")
		quota, err := s.rewardRepo.GetQuota(req.DiscourseUserID, today)
		earnedToday := int64(0)
		if err == nil {
			earnedToday = quota.EarnedToday
		}

		// 8. Check daily caps (user-level + level-level)
		userDailyCap := rule.DailyCapPerUser
		if userDailyCap > cap.DailyCap {
			userDailyCap = cap.DailyCap
		}
		if earnedToday >= userDailyCap {
			return &RejectError{Reason: "user_daily_cap_reached"}
		}

		// 9. Calculate reward
		baseAmount := rule.Amount
		multiplier := cap.RewardMultiplier
		finalAmount := int64(float64(baseAmount) * multiplier)

		// 10. Ensure not exceeding daily cap
		remaining := userDailyCap - earnedToday
		if finalAmount > remaining {
			finalAmount = remaining
		}

		// 11. Check public pool balance
		publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}
		if publicPool.Balance < finalAmount {
			return errcode.ErrPoolInsufficient
		}

		// 12. Deduct from public pool
		if err := s.poolRepo.UpdateBalance(publicPool.ID, publicPool.Balance-finalAmount); err != nil {
			return errcode.ErrDatabase
		}

		// 13. Add to user balance
		newBalance := ub.Balance + finalAmount
		if err := s.balanceRepo.UpdateBalanceAndEarned(tx, req.DiscourseUserID, newBalance, ub.Version+1, ub.TotalEarned+finalAmount); err != nil {
			return err
		}

		// 14. Update daily quota
		newEarned := earnedToday + finalAmount
		if err := s.rewardRepo.UpsertQuota(tx, req.DiscourseUserID, today, newEarned, fmt.Sprintf("{\"%s\":1}", req.Action)); err != nil {
			return errcode.ErrDatabase
		}

		// 15. Write transaction
		txModel := &repository.Transaction{
			TxType:          "reward",
			DiscourseUserID: req.DiscourseUserID,
			Amount:          finalAmount,
			Fee:             0,
			NetAmount:       finalAmount,
			FromType:        "pool",
			ToType:          "user",
			IdempotencyKey:  req.IdempotencyKey,
			RefType:         &req.RefType,
			RefID:           &req.RefID,
			Status:          "completed",
		}
		ip := req.IPAddress
		_ = ip

		if err := s.txRepo.Create(tx, txModel); err != nil {
			return errcode.ErrDatabase
		}

		// 16. Write reward log
		rewardLog := &repository.RewardLogRecord{
			DiscourseUserID: req.DiscourseUserID,
			Action:          req.Action,
			Amount:          finalAmount,
			RefID:           req.RefID,
			RefType:         req.RefType,
			TrustLevel:      ub.TrustLevel,
			Multiplier:      multiplier,
			Status:          "completed",
		}
		if req.IPAddress != "" {
			rewardLog.IPAddress = &req.IPAddress
		}
		if err := s.rewardRepo.CreateLog(tx, rewardLog); err != nil {
			return errcode.ErrDatabase
		}

		resp = &RewardResponse{
			Amount:      baseAmount,
			TrustLevel:  ub.TrustLevel,
			Multiplier:  multiplier,
			FinalAmount: finalAmount,
			EarnedToday: newEarned,
			DailyCap:    userDailyCap,
			Status:      "completed",
		}
		return nil
	})

	if err != nil {
		if reject, ok := err.(*RejectError); ok {
			return &RewardResponse{Status: "rejected", RejectReason: reject.Reason}, nil
		}
		if ec, ok := err.(errcode.ErrorCode); ok {
			return &RewardResponse{Status: "rejected", RejectReason: ec.Message}, nil
		}
		return nil, err
	}

	return resp, nil
}

type RejectError struct {
	Reason string
}

func (e *RejectError) Error() string {
	return e.Reason
}

// Checkin records a daily checkin
func (s *RewardService) Checkin(userID int64, ipAddress string) (*RewardResponse, error) {
	today := time.Now().Format("2006-01-02")
	idempotencyKey := fmt.Sprintf("checkin_%d_%s", userID, today)

	// Check if already checked in today
	quota, err := s.rewardRepo.GetQuota(userID, today)
	if err == nil {
		// Check if checkin action was already counted
		existingLog, _ := s.rewardRepo.GetLogByRef("checkin", userID, userID)
		if existingLog != nil {
			// Check if it was today
			if existingLog.CreatedAt.Format("2006-01-02") == today {
				return &RewardResponse{Status: "rejected", RejectReason: "already_checked_in"}, nil
			}
		}
		_ = quota
	}

	return s.Grant(&RewardRequest{
		Action:          "checkin",
		DiscourseUserID: userID,
		RefType:         "checkin",
		RefID:           userID,
		IdempotencyKey:  idempotencyKey,
		IPAddress:       ipAddress,
	})
}

// CheckinStatus returns today's checkin status
func (s *RewardService) CheckinStatus(userID int64) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")
	checkedIn := false
	streak := int64(0)

	quota, err := s.rewardRepo.GetQuota(userID, today)
	if err == nil && quota.EarnedToday > 0 {
		checkedIn = true
	}

	return map[string]interface{}{
		"checked_in_today": checkedIn,
		"streak":           streak,
		"last_checkin_date": today,
	}, nil
}
func (s *RewardService) ListRewardRules() ([]repository.RewardRule, error) {
	return s.rewardRepo.ListRules()
}
