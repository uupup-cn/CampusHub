package handler

import (
	"github.com/campushub/chb-backend/internal/config"
	"github.com/campushub/chb-backend/internal/idp"
	"github.com/campushub/chb-backend/internal/middleware"
	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(&cfg.CORS))
	r.Use(middleware.RateLimit(cfg.RateLimit.Enabled, cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst))

	// Repositories
	poolRepo := repository.NewPoolRepo(db)
	balanceRepo := repository.NewUserBalanceRepo(db)
	txRepo := repository.NewTransactionRepo(db)
	rewardRepo := repository.NewRewardRepo(db)
	appRepo := repository.NewAppRepo(db)
	marketRepo := repository.NewMarketplaceRepo(db)

	// Services
	ledgerSvc := service.NewLedgerService(db, poolRepo, balanceRepo, txRepo)
	rewardSvc := service.NewRewardService(db, rewardRepo, poolRepo, balanceRepo, txRepo, ledgerSvc)
	idpSvc := idp.NewIdpService(db, appRepo)
	marketSvc := service.NewMarketplaceService(db, marketRepo, balanceRepo, poolRepo, txRepo)
	adminSvc := service.NewAdminService(db, balanceRepo, txRepo, poolRepo, rewardRepo, marketRepo, appRepo)

	// Handlers
	health := NewHealthHandler(db)
	ledger := NewLedgerHandler(ledgerSvc, appRepo)
	reward := NewRewardHandler(rewardSvc)
	oauth := NewOAuthHandler(idpSvc)
	market := NewMarketplaceHandler(marketSvc)
	admin := NewAdminHandler(adminSvc)
	auth := NewAuthHandler()
	userH := NewUserHandler(db, txRepo, appRepo)

	// Routes
	r.GET("/health", health.Health)
	r.GET("/api/health", health.Health)

	api := r.Group("/api")
	{
		// CHB API - dual auth: Bearer token (OAuth) or X-User-ID (internal)
		chb := api.Group("/chb")
		chb.Use(middleware.OptionalAuth(appRepo))
		{
			// Read operations - need "read" scope for OAuth
			chb.GET("/balance", ledger.GetBalance)
			chb.GET("/transactions", ledger.ListTransactions)
			chb.GET("/checkin/status", reward.CheckinStatus)
			chb.GET("/pools", ledger.GetPools)
			chb.GET("/audit", ledger.Audit)
			chb.GET("/me/transactions", userH.ListMyTransactions)

			// Write operations - need "spend" scope for OAuth
			chb.POST("/spend", ledger.Spend)
			chb.POST("/checkin", reward.Checkin)

			// Internal-only routes (plugin calls with X-API-Key)
			chb.POST("/reward", reward.GrantReward)
			chb.POST("/sync/trust-level", reward.SyncTrustLevel)
			chb.GET("/reward/rules", reward.ListRewardRules)
			chb.POST("/release", ledger.Release)
		}

		// Marketplace - dual auth
		marketplace := api.Group("/marketplace")
		marketplace.Use(middleware.OptionalAuth(appRepo))
		{
			marketplace.GET("/items", market.ListItems)
			marketplace.GET("/items/mine", market.ListMyItems)
			marketplace.GET("/items/:id", market.GetItem)
			marketplace.POST("/items", market.CreateItem)
			marketplace.POST("/orders", market.CreateOrder)
			marketplace.GET("/orders", market.ListOrders)
			marketplace.POST("/apply", market.ApplyMerchant)
			marketplace.GET("/my-status", market.MyStatus)
		}

		// Admin API - X-Admin-Key
		adminAPI := api.Group("/admin")
		adminAPI.Use(middleware.AdminAuth(cfg.AdminKey))
		{
			adminAPI.GET("/settings", admin.GetSettings)
			adminAPI.PUT("/settings", admin.UpdateSettings)

			adminAPI.GET("/trust-levels", admin.ListTrustLevels)
			adminAPI.PUT("/trust-levels", admin.UpdateTrustLevel)

			adminAPI.PUT("/reward/rules", admin.UpdateRewardRule)

			adminAPI.GET("/apps", admin.ListApps)
			adminAPI.POST("/apps", admin.CreateApp)
			adminAPI.PUT("/apps/:id", admin.UpdateApp)
			adminAPI.DELETE("/apps/:id", admin.DeleteApp)

			adminAPI.GET("/marketplace/applications", admin.ListApplications)
			adminAPI.PUT("/marketplace/applications/:id", admin.ReviewApplication)
			adminAPI.GET("/marketplace/items", admin.ListPendingItems)
			adminAPI.PUT("/marketplace/items/:id", admin.ReviewItem)

			adminAPI.GET("/users", admin.ListUsers)
			adminAPI.PUT("/users/:id/freeze", admin.FreezeUser)
			adminAPI.PUT("/users/:id/unfreeze", admin.UnfreezeUser)
			adminAPI.POST("/users/:id/recover", admin.RecoverPoints)

			adminAPI.GET("/audit-logs", admin.ListAuditLogs)
			adminAPI.GET("/stats", admin.GetStats)
		}

		// Auth API - Discourse session check
		authAPI := api.Group("/auth")
		{
			authAPI.GET("/me", auth.Me)
		}

		// OAuth API
		oauthAPI := api.Group("/oauth")
		{
			oauthAPI.GET("/app-info", oauth.AppInfo)
			oauthAPI.GET("/my-apps", userH.ListMyApps)
			oauthAPI.POST("/authorize/confirm", oauth.Confirm)
		}
	}

	r.GET("/oauth/authorize", oauth.Authorize)
	r.POST("/oauth/token", oauth.Token)
	r.POST("/oauth/introspect", oauth.Introspect)

	return r
}
