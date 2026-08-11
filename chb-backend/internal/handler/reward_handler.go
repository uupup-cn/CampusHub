package handler

import (
	"github.com/campushub/chb-backend/internal/service"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	svc *service.RewardService
}

func NewRewardHandler(svc *service.RewardService) *RewardHandler {
	return &RewardHandler{svc: svc}
}

// GrantReward - internal API for reward plugin
func (h *RewardHandler) GrantReward(c *gin.Context) {
	var req struct {
		Action          string `json:"action"`
		DiscourseUserID int64  `json:"discourse_user_id"`
		RefType         string `json:"ref_type"`
		RefID           int64  `json:"ref_id"`
		IdempotencyKey  string `json:"idempotency_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Action == "" || req.DiscourseUserID == 0 || req.IdempotencyKey == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	result, err := h.svc.Grant(&service.RewardRequest{
		Action:          req.Action,
		DiscourseUserID: req.DiscourseUserID,
		RefType:         req.RefType,
		RefID:           req.RefID,
		IdempotencyKey:  req.IdempotencyKey,
		IPAddress:       c.ClientIP(),
	})
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

// Checkin - user daily checkin
func (h *RewardHandler) Checkin(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}

	result, err := h.svc.Checkin(userID, c.ClientIP())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

// CheckinStatus - query today's checkin status
func (h *RewardHandler) CheckinStatus(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}

	status, err := h.svc.CheckinStatus(userID)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, status)
}

// ListRewardRules - admin API
func (h *RewardHandler) ListRewardRules(c *gin.Context) {
	// TODO: admin auth
	rules, err := h.svc.ListRewardRules()
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, rules)
}

// SyncTrustLevel - internal API for trust level sync
func (h *RewardHandler) SyncTrustLevel(c *gin.Context) {
	var req struct {
		DiscourseUserID int64 `json:"discourse_user_id"`
		TrustLevel      int16 `json:"trust_level"`
		IdempotencyKey  string `json:"idempotency_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	// TODO: implement trust level sync logic
	response.Success(c, gin.H{
		"old_trust_level": 0,
		"new_trust_level": req.TrustLevel,
	})
}
