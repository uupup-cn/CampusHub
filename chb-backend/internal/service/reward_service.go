package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"gorm.io/gorm"
)

type RewardService struct {
	db          *gorm.DB
	rewardRepo  *repository.RewardRepo
	poolRepo    *repository.PoolRepo
	balanceRepo *repository.UserBalanceRepo
	txRepo      *repository.TransactionRepo
	ledgerSvc   *LedgerService
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
	TrustLevel      int16
	IPAddress       string
}

type RewardResponse struct {
	Amount       int64   `json:"amount"`
	TrustLevel   int16   `json:"trust_level"`
	Multiplier   float64 `json:"multiplier"`
	FinalAmount  int64   `json:"final_amount"`
	EarnedToday  int64   `json:"earned_today"`
	ActionEarned int64   `json:"action_earned"`
	DailyCap     int64   `json:"daily_cap"`
	Status       string  `json:"status"`
	RejectReason string  `json:"reject_reason,omitempty"`
}

func (s *RewardService) Grant(req *RewardRequest) (*RewardResponse, error) {
	// 1. Idempotency check by ref
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
				ub = &repository.UserBalance{
					DiscourseUserID: req.DiscourseUserID,
					Username:        fmt.Sprintf("user_%d", req.DiscourseUserID),
					Balance:         0,
					Version:         1,
					TrustLevel:      req.TrustLevel,
					Status:          "active",
				}
				if err := s.balanceRepo.CreateWithTx(tx, ub); err != nil {
					return errcode.ErrDatabase
				}
			} else {
				return errcode.ErrDatabase
			}
		}

		// 4b. Update trust level if provided and different
		if req.TrustLevel > 0 && ub.TrustLevel != req.TrustLevel {
			if err := s.balanceRepo.UpdateTrustLevel(tx, req.DiscourseUserID, req.TrustLevel); err != nil {
				return errcode.ErrDatabase
			}
			ub.TrustLevel = req.TrustLevel
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

		// 7. Get today quota with row lock
		today := time.Now().Format("2006-01-02")
		quota, err := s.rewardRepo.GetQuotaForUpdate(tx, req.DiscourseUserID, today)
		totalEarnedToday := int64(0)
		actionEarnedMap := map[string]int64{}
		if err == nil && quota != nil {
			totalEarnedToday = quota.EarnedToday
			if quota.ActionCounts != "" {
				_ = json.Unmarshal([]byte(quota.ActionCounts), &actionEarnedMap)
			}
		}

		// 8. Check per-action daily cap
		actionEarned := actionEarnedMap[req.Action]
		actionCap := rule.DailyCapPerUser
		if actionEarned >= actionCap {
			return &RejectError{Reason: "action_daily_cap_reached"}
		}

		// 9. Check total daily cap (trust level cap)
		totalCap := cap.DailyCap
		if totalEarnedToday >= totalCap {
			return &RejectError{Reason: "total_daily_cap_reached"}
		}

		// 10. Calculate reward
		baseAmount := rule.Amount
		multiplier := cap.RewardMultiplier
		finalAmount := int64(float64(baseAmount) * multiplier)

		// 11. Ensure not exceeding per-action remaining
		actionRemaining := actionCap - actionEarned
		if finalAmount > actionRemaining {
			finalAmount = actionRemaining
		}

		// 12. Ensure not exceeding total remaining
		totalRemaining := totalCap - totalEarnedToday
		if finalAmount > totalRemaining {
			finalAmount = totalRemaining
		}

		// 13. Check public pool balance
		publicPool, err := s.poolRepo.GetPublicPoolWithLock(tx)
		if err != nil {
			return errcode.ErrDatabase
		}
		if publicPool.Balance < finalAmount {
			return errcode.ErrPoolInsufficient
		}

		// 14. Deduct from public pool
		if err := s.poolRepo.UpdateBalance(publicPool.ID, publicPool.Balance-finalAmount); err != nil {
			return errcode.ErrDatabase
		}

		// 15. Add to user balance
		newBalance := ub.Balance + finalAmount
		if err := s.balanceRepo.UpdateBalanceAndEarned(tx, req.DiscourseUserID, newBalance, ub.Version+1, ub.TotalEarned+finalAmount); err != nil {
			return err
		}

		// 16. Update daily quota
		newTotalEarned := totalEarnedToday + finalAmount
		actionEarnedMap[req.Action] = actionEarned + finalAmount
		actionCountsJSON, _ := json.Marshal(actionEarnedMap)
		if err := s.rewardRepo.UpsertQuota(tx, req.DiscourseUserID, today, newTotalEarned, string(actionCountsJSON)); err != nil {
			return errcode.ErrDatabase
		}

		// 17. Write transaction
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
		if err := s.txRepo.Create(tx, txModel); err != nil {
			return errcode.ErrDatabase
		}

		// 18. Write reward log
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
		if req.IPAddress != " + dq + dq + " {
			rewardLog.IPAddress = &req.IPAddress
		}
		if err := s.rewardRepo.CreateLog(tx, rewardLog); err != nil {
			return errcode.ErrDatabase
		}

		resp = &RewardResponse{
			Amount:       baseAmount,
			TrustLevel:   ub.TrustLevel,
			Multiplier:   multiplier,
			FinalAmount:  finalAmount,
			EarnedToday:  newTotalEarned,
			ActionEarned: actionEarned + finalAmount,
			DailyCap:     actionCap,
			Status:       "completed",
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
	return s.CheckinWithTrustLevel(userID, 0, ipAddress)
}

// CheckinWithTrustLevel records a daily checkin with trust level
func (s *RewardService) CheckinWithTrustLevel(userID int64, trustLevel int16, ipAddress string) (*RewardResponse, error) {
	today := time.Now().Format("2006-01-02")
	idempotencyKey := fmt.Sprintf("checkin_%d_%s", userID, today)

	// Check if already checked in today
	existingLog, _ := s.rewardRepo.GetLogByRef("checkin", userID, userID)
	if existingLog != nil {
		if existingLog.CreatedAt.Format("2006-01-02") == today {
			return &RewardResponse{Status: "rejected", RejectReason: "already_checked_in"}, nil
		}
	}

	return s.Grant(&RewardRequest{
		Action:          "checkin",
		DiscourseUserID: userID,
		RefType:         "checkin",
		RefID:           userID,
		IdempotencyKey:  idempotencyKey,
		TrustLevel:      trustLevel,
		IPAddress:       ipAddress,
	})
}

// CheckinStatus returns today checkin status
func (s *RewardService) CheckinStatus(userID int64) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")
	checkedIn := false
	streak := int64(0)

	quota, err := s.rewardRepo.GetQuota(userID, today)
	if err == nil && quota != nil {
		actionMap := map[string]int64{}
		if quota.ActionCounts != " + dq + dq + " {
			_ = json.Unmarshal([]byte(quota.ActionCounts), &actionMap)
		}
		if actionMap["checkin"] > 0 {
			checkedIn = true
		}
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

// SyncTrustLevel updates user trust level in the balance table
func (s *RewardService) SyncTrustLevel(discourseUserID int64, trustLevel int16) (map[string]interface{}, error) {
	ub, err := s.balanceRepo.GetByUserID(discourseUserID)
	oldLevel := int16(0)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ub = &repository.UserBalance{
				DiscourseUserID: discourseUserID,
				Username:        fmt.Sprintf("user_%d", discourseUserID),
				Balance:         0,
				Version:         1,
				TrustLevel:      trustLevel,
				Status:          "active",
			}
			if err := s.balanceRepo.Create(ub); err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"old_trust_level": 0,
				"new_trust_level": trustLevel,
				"created":         true,
			}, nil
		}
		return nil, err
	}
	oldLevel = ub.TrustLevel

	if ub.TrustLevel != trustLevel {
		if err := s.balanceRepo.UpdateTrustLevel(s.db, discourseUserID, trustLevel); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"old_trust_level": oldLevel,
		"new_trust_level": trustLevel,
		"updated":         oldLevel != trustLevel,
	}, nil
}
