package scheduler

import (
	"log"
	"time"

	"github.com/campushub/chb-backend/internal/repository"
	"github.com/campushub/chb-backend/internal/service"
	"gorm.io/gorm"
)

type Scheduler struct {
	db            *gorm.DB
	disputeSvc    *service.DisputeService
	orderRepo     *repository.MarketplaceRepo
	balanceRepo   *repository.UserBalanceRepo
	txRepo        *repository.TransactionRepo
}

func NewScheduler(db *gorm.DB, disputeSvc *service.DisputeService, orderRepo *repository.MarketplaceRepo, balanceRepo *repository.UserBalanceRepo, txRepo *repository.TransactionRepo) *Scheduler {
	return &Scheduler{db: db, disputeSvc: disputeSvc, orderRepo: orderRepo, balanceRepo: balanceRepo, txRepo: txRepo}
}

// Start 启动定时任务，每天下午6点(CST)执行
func (s *Scheduler) Start() {
	go func() {
		for {
			now := time.Now()
			// 计算下一个18:00 CST
			loc, _ := time.LoadLocation("Asia/Shanghai")
			nowCST := now.In(loc)
			next := time.Date(nowCST.Year(), nowCST.Month(), nowCST.Day(), 18, 0, 0, 0, loc)
			if nowCST.After(next) {
				next = next.Add(24 * time.Hour)
			}
			waitDuration := next.Sub(nowCST)
			log.Printf("[Scheduler] 下次执行时间: %s (等待 %s)", next.Format("2006-01-02 15:04:05"), waitDuration)
			time.Sleep(waitDuration)
			s.run()
		}
	}()
}

func (s *Scheduler) run() {
	log.Println("[Scheduler] 开始执行定时任务")

	// 任务1: 未来积分自动转入
	s.releasePendingBalances()

	// 任务2: 争议超时自动退款
	if err := s.disputeSvc.AutoRefundExpired(); err != nil {
		log.Printf("[Scheduler] 争议超时退款失败: %v", err)
	}

	// 任务3: 卖家胜诉积分转入
	if err := s.disputeSvc.ReleaseSellerWin(); err != nil {
		log.Printf("[Scheduler] 卖家胜诉积分转入失败: %v", err)
	}

	log.Println("[Scheduler] 定时任务执行完成")
}

// releasePendingBalances 将到期未来积分转为可用积分
func (s *Scheduler) releasePendingBalances() {
	// 查询到期且无争议的订单
	var orders []repository.MarketplaceOrderModel
	s.db.Where("status = ? AND pending_release_at <= ? AND seller_pending_credited = ? AND dispute_status IS NULL",
		"completed", time.Now(), false).Find(&orders)

	for _, order := range orders {
		tx := s.db.Begin()
		err := s.balanceRepo.TransferPendingToAvailable(tx, order.SellerID, order.NetAmount)
		if err != nil {
			tx.Rollback()
			log.Printf("[Scheduler] 订单#%d 积分转入失败: %v", order.ID, err)
			continue
		}
		tx.Model(&repository.MarketplaceOrderModel{}).Where("id = ?", order.ID).
			Update("seller_pending_credited", true)
		tx.Commit()
		log.Printf("[Scheduler] 订单#%d 卖家未来积分 %d 已转为可用积分", order.ID, order.NetAmount)
	}
}


// StartScheduler 创建并启动调度器
func StartScheduler(db *gorm.DB) {
	disputeRepo := repository.NewDisputeRepo(db)
	orderRepo := repository.NewMarketplaceRepo(db)
	balanceRepo := repository.NewUserBalanceRepo(db)
	txRepo := repository.NewTransactionRepo(db)
	poolRepo := repository.NewPoolRepo(db)
	disputeSvc := service.NewDisputeService(db, disputeRepo, orderRepo, balanceRepo, txRepo, poolRepo)

	s := NewScheduler(db, disputeSvc, orderRepo, balanceRepo, txRepo)
	s.Start()
	log.Println("[Scheduler] 定时任务已启动")
}
