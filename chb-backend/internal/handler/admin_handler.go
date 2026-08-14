package handler

import (
	"strconv"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/internal/service"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// ===== Settings =====

func (h *AdminHandler) GetSettings(c *gin.Context) {
	response.Success(c, h.svc.GetSettings())
}

func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req struct {
		MarketplaceFeeRate    float64 `json:"marketplace_fee_rate"`
		AutoReleaseEnabled    bool    `json:"auto_release_enabled"`
		AutoReleaseThreshold  int     `json:"auto_release_threshold"`
		AutoReleaseRatio      int     `json:"auto_release_ratio"`
		AutoReleaseMonthlyCap int64   `json:"auto_release_monthly_cap"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	settings := &service.SystemSetting{
		MarketplaceFeeRate:    req.MarketplaceFeeRate,
		AutoReleaseEnabled:    req.AutoReleaseEnabled,
		AutoReleaseThreshold:  req.AutoReleaseThreshold,
		AutoReleaseRatio:      req.AutoReleaseRatio,
		AutoReleaseMonthlyCap: req.AutoReleaseMonthlyCap,
	}
	if err := h.svc.UpdateSettings(settings); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, req)
}

// ===== Trust Levels =====

func (h *AdminHandler) ListTrustLevels(c *gin.Context) {
	caps, err := h.svc.ListTrustLevelCaps()
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, caps)
}

func (h *AdminHandler) UpdateTrustLevel(c *gin.Context) {
	var req repository.TrustLevelCapModel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.UpdateTrustLevelCap(&req); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

// ===== Reward Rules =====

func (h *AdminHandler) UpdateRewardRule(c *gin.Context) {
	var req repository.RewardRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.UpdateRewardRule(&req); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

// ===== Apps =====

func (h *AdminHandler) ListApps(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	apps, total, err := h.svc.ListApps(page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.SuccessPaginated(c, apps, total, page, pageSize)
}

func (h *AdminHandler) CreateApp(c *gin.Context) {
	var req repository.AppModel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.CreateApp(&req); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	// 客户端凭据只在创建成功时返回一次，secret 不落入列表/详情接口
	response.Success(c, gin.H{
		"id":            req.ID,
		"app_name":      req.AppName,
		"client_id":     req.ClientID,
		"client_secret": req.ClientSecret,
		"min_trust_level": req.MinTrustLevel,
		"fee_rate":      req.FeeRate,
	})
}

func (h *AdminHandler) UpdateApp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req repository.AppModel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	req.ID = id
	if err := h.svc.UpdateApp(&req); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) DeleteApp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.DeleteApp(id); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

// ===== Marketplace Review =====

func (h *AdminHandler) ListApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	apps, total, err := h.svc.ListPendingApplications(page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.SuccessPaginated(c, apps, total, page, pageSize)
}

func (h *AdminHandler) ReviewApplication(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req struct {
		Status  string `json:"status"`
		Comment string `json:"review_comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.ReviewApplication(id, req.Status, req.Comment, 1); err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok { response.Error(c, ec); return }; response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) ListPendingItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.ListPendingItems(page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

func (h *AdminHandler) ReviewItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.ReviewItem(id, req.Status); err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok { response.Error(c, ec); return }; response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// ===== Users =====

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	users, total, err := h.svc.ListUsers(page, pageSize, keyword)
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.SuccessPaginated(c, users, total, page, pageSize)
}

func (h *AdminHandler) FreezeUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.FreezeUser(id); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) UnfreezeUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.UnfreezeUser(id); err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.Success(c, nil)
}

// ===== Stats =====

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok { response.Error(c, ec); return }; response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, stats)
}

// ===== Audit Logs =====

func (h *AdminHandler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	logs, total, err := h.svc.ListAuditLogs(
		c.Query("action"),
		0,
		c.Query("start_date"),
		c.Query("end_date"),
		page, pageSize,
	)
	if err != nil {
		response.Error(c, errcode.ErrDatabase)
		return
	}
	response.SuccessPaginated(c, logs, total, page, pageSize)
}

func (h *AdminHandler) RecoverPoints(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req struct {
		Amount int64  `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.RecoverPoints(id, req.Amount, req.Reason, 1); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}
func (h *AdminHandler) GetApp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	app, err := h.svc.GetAppByID(id)
	if err != nil {
		response.Error(c, errcode.ErrNotFound)
		return
	}
	response.Success(c, app)
}
func (h *AdminHandler) AdjustPoints(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req struct {
		Amount    int64  `json:"amount"`
		Direction string `json:"direction"`
		Reason    string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Amount <= 0 || req.Reason == "" || (req.Direction != "add" && req.Direction != "deduct") {
		response.Error(c, errcode.ErrParamMissing)
		return
	}
	if err := h.svc.AdjustPoints(id, req.Amount, req.Direction, req.Reason, 1); err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}
