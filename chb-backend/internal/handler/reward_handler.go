package handler

import (
	"fmt"

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

func (h *RewardHandler) GrantReward(c *gin.Context) {
	var req struct {
		Action          string `json:"action"`
		DiscourseUserID int64  `json:"discourse_user_id"`
		RefType         string `json:"ref_type"`
		RefID           int64  `json:"ref_id"`
		IdempotencyKey  string `json:"idempotency_key"`
		TrustLevel      int16  `json:"trust_level"`
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
		TrustLevel:      req.TrustLevel,
		IPAddress:       c.ClientIP(),
	})
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *RewardHandler) Checkin(c *gin.Context) {
	var req struct {
		DiscourseUserID int64 `json:"discourse_user_id"`
		TrustLevel      int16 `json:"trust_level"`
	}
	_ = c.ShouldBindJSON(&req)

	userID := req.DiscourseUserID
	trustLevel := req.TrustLevel
	if userID == 0 {
		userID = getUserID(c)
		trustLevel = getTrustLevel(c)
	}
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}

	result, err := h.svc.CheckinWithTrustLevel(userID, trustLevel, c.ClientIP())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *RewardHandler) CheckinStatus(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID int64
	if userIDStr != "" {
		fmt.Sscanf(userIDStr, "%d", &userID)
	}
	if userID == 0 {
		userID = getUserID(c)
	}
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

func (h *RewardHandler) ListRewardRules(c *gin.Context) {
	rules, err := h.svc.ListRewardRules()
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, rules)
}

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
	if req.DiscourseUserID == 0 {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	result, err := h.svc.SyncTrustLevel(req.DiscourseUserID, req.TrustLevel)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}
