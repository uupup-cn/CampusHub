package handler

import (
	"strconv"

	"github.com/campushub/chb-backend/internal/service"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type DisputeHandler struct {
	svc *service.DisputeService
}

func NewDisputeHandler(svc *service.DisputeService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

// 买家发起争议
func (h *DisputeHandler) CreateDispute(c *gin.Context) {
	buyerID := getUserID(c)
	if buyerID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req service.CreateDisputeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Reason == "" {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	dispute, err := h.svc.CreateDispute(buyerID, orderID, req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, dispute)
}

// 查询争议列表
func (h *DisputeHandler) ListDisputes(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	role := c.DefaultQuery("role", "buyer")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	list, total, err := h.svc.ListMyDisputes(userID, role, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, list, total, page, pageSize)
}

// 争议详情
func (h *DisputeHandler) GetDispute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	dispute, err := h.svc.GetDispute(id)
	if err != nil {
		response.Error(c, errcode.ErrNotFound)
		return
	}
	response.Success(c, dispute)
}

// 卖家同意退款
func (h *DisputeHandler) AcceptDispute(c *gin.Context) {
	sellerID := getUserID(c)
	if sellerID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.AcceptDispute(sellerID, id); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"status": "accepted"})
}

// 卖家拒绝退款
func (h *DisputeHandler) RejectDispute(c *gin.Context) {
	sellerID := getUserID(c)
	if sellerID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req service.RejectDisputeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Reason == "" {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.RejectDispute(sellerID, id, req); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"status": "rejected"})
}

// 管理员争议列表
func (h *DisputeHandler) AdminListDisputes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	list, total, err := h.svc.ListAllDisputes(page, pageSize, status)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, list, total, page, pageSize)
}

// 管理员争议详情
func (h *DisputeHandler) AdminGetDispute(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	dispute, err := h.svc.GetDispute(id)
	if err != nil {
		response.Error(c, errcode.ErrNotFound)
		return
	}
	response.Success(c, dispute)
}

// 管理员判定
type AdminDecideReq struct {
	Decision string `json:"decision"`
	Note     string `json:"note"`
}

func (h *DisputeHandler) AdminDecide(c *gin.Context) {
	adminID := int64(1) // 管理员ID固定为1
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	var req AdminDecideReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Decision != "buyer_win" && req.Decision != "seller_win" {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if err := h.svc.AdminDecide(adminID, id, req.Decision, req.Note); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"status": "decided"})
}
