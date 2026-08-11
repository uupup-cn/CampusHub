package handler

import (
	"strconv"

	"github.com/campushub/chb-backend/internal/service"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type MarketplaceHandler struct {
	svc *service.MarketplaceService
}

func NewMarketplaceHandler(svc *service.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc}
}

func (h *MarketplaceHandler) ListItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	q := service.ItemQuery{
		Category: c.Query("category"),
		Keyword:  c.Query("keyword"),
		Sort:     c.Query("sort"),
		Page:     page,
		PageSize: pageSize,
	}
	items, total, err := h.svc.ListItems(q)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

func (h *MarketplaceHandler) GetItem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	item, err := h.svc.GetItem(id)
	if err != nil {
		response.Error(c, errcode.ErrNotFound)
		return
	}
	response.Success(c, item)
}

// ListMyItems 返回当前卖家发布的全部商品（含待审核/已拒绝），供“我的商品”页使用。
func (h *MarketplaceHandler) ListMyItems(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	items, total, err := h.svc.ListMyItems(userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

func (h *MarketplaceHandler) CreateItem(c *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Price       int64  `json:"price"`
		Stock       int    `json:"stock"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.Title == "" || req.Price <= 0 {
		response.Error(c, errcode.ErrParamMissing)
		return
	}
	sellerID := getUserID(c)
	if sellerID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	item, err := h.svc.CreateItem(sellerID, req.Title, req.Description, req.Price, req.Stock, req.Category, req.ImageURL)
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok { response.Error(c, ec); return }; response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, item)
}

func (h *MarketplaceHandler) CreateOrder(c *gin.Context) {
	var req struct {
		ItemID         int64  `json:"item_id"`
		Quantity       int    `json:"quantity"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.ItemID <= 0 || req.IdempotencyKey == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}
	buyerID := getUserID(c)
	if buyerID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	result, err := h.svc.CreateOrder(buyerID, req.ItemID, req.Quantity, req.IdempotencyKey)
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

func (h *MarketplaceHandler) ListOrders(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	role := c.DefaultQuery("role", "buyer")
	status := c.Query("status")
	items, total, err := h.svc.ListOrders(userID, role, status, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

func (h *MarketplaceHandler) ApplyMerchant(c *gin.Context) {
	var req struct {
		ShopName    string `json:"shop_name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.ShopName == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	err := h.svc.ApplyMerchant(userID, req.ShopName, req.Description)
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"status": "pending"})
}
