package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	dbOK := true
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbOK = false
	}

	status := "ok"
	if !dbOK {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"database": dbOK,
		"version":  "1.0.0",
		"service":  "chb-backend",
	})
}
