package handler

import (
	"strconv"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
	txRepo   *repository.TransactionRepo
	appRepo  *repository.AppRepo
}

func NewUserHandler(db *gorm.DB, txRepo *repository.TransactionRepo, appRepo *repository.AppRepo) *UserHandler {
	return &UserHandler{db: db, txRepo: txRepo, appRepo: appRepo}
}

// ListMyTransactions returns paginated transaction history for the current user
// GET /api/chb/me/transactions
func (h *UserHandler) ListMyTransactions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	txType := c.Query("type")

	items, total, err := h.txRepo.ListByUser(userID, txType, "", "", page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.SuccessPaginated(c, items, total, page, pageSize)
}

// ListMyApps returns OAuth2 apps the current user has authorized
// GET /api/oauth/my-apps
func (h *UserHandler) ListMyApps(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}

	var tokens []repository.AccessTokenModel
	err := h.db.Where("discourse_user_id = ? AND revoked = ?", userID, false).
		Order("created_at DESC").
		Find(&tokens).Error
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}

	type AppInfo struct {
		ID        int64  `json:"id"`
		AppName   string `json:"app_name"`
		ClientID  string `json:"client_id"`
		Scopes    string `json:"scopes"`
		CreatedAt string `json:"created_at"`
	}

	var result []AppInfo
	for _, t := range tokens {
		app, err := h.appRepo.GetByClientID(t.ClientID)
		if err != nil || app == nil {
			continue
		}
		result = append(result, AppInfo{
			ID:        app.ID,
			AppName:   app.AppName,
			ClientID:  t.ClientID,
			Scopes:    t.Scopes,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		})
	}
	if result == nil {
		result = []AppInfo{}
	}
	response.Success(c, gin.H{"items": result})
}
