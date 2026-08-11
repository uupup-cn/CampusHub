package handler

import (
	"strconv"

	"github.com/campushub/chb-backend/internal/service"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	svc *service.LedgerService
}

func NewLedgerHandler(svc *service.LedgerService) *LedgerHandler {
	return &LedgerHandler{svc: svc}
}

func (h *LedgerHandler) GetBalance(c *gin.Context) {
	userID := getUserID(c)
	info, err := h.svc.GetBalance(userID)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, info)
}

func (h *LedgerHandler) Spend(c *gin.Context) {
	var req struct {
		Amount         int64  `json:"amount"`
		IdempotencyKey string `json:"idempotency_key"`
		Description    string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Amount <= 0 || req.IdempotencyKey == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	userID := getUserID(c)
	result, err := h.svc.Spend(userID, req.Amount, req.IdempotencyKey, req.Description, 10.0, nil)
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *LedgerHandler) ListTransactions(c *gin.Context) {
	userID := getUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	q := service.TransactionQuery{
		TxType:    c.Query("tx_type"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Page:      page,
		PageSize:  pageSize,
	}

	items, total, err := h.svc.ListTransactions(userID, q)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

func (h *LedgerHandler) GetPools(c *gin.Context) {
	info, err := h.svc.GetPools()
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, info)
}

func (h *LedgerHandler) Release(c *gin.Context) {
	var req struct {
		Amount int64  `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}

	result, err := h.svc.Release(req.Amount, req.Reason, 1)
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *LedgerHandler) Audit(c *gin.Context) {
	info, err := h.svc.Audit()
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, info)
}

func getUserID(c *gin.Context) int64 {
	if v, exists := c.Get("discourse_user_id"); exists {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	idStr := c.GetHeader("X-User-ID")
	if idStr != "" {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		return id
	}
	return 0
}

// getTrustLevel 读取当前用户信任等级：优先取中间件注入值，其次取 X-Trust-Level 请求头（开发/测试环境）。
func getTrustLevel(c *gin.Context) int16 {
	if v, exists := c.Get("trust_level"); exists {
		if tl, ok := v.(int16); ok {
			return tl
		}
	}
	tlStr := c.GetHeader("X-Trust-Level")
	if tlStr != "" {
		tl, _ := strconv.ParseInt(tlStr, 10, 16)
		return int16(tl)
	}
	return 0
}
