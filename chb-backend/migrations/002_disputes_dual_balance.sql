-- CampusHub 争议退款 + 双积分体系迁移
-- 执行前请备份数据库

-- 1. user_balances 增加未来积分字段
ALTER TABLE user_balances ADD COLUMN IF NOT EXISTS pending_balance BIGINT NOT NULL DEFAULT 0;

-- 2. marketplace_orders 增加争议和未来积分相关字段
ALTER TABLE marketplace_orders ADD COLUMN IF NOT EXISTS dispute_status VARCHAR(20) DEFAULT NULL;
ALTER TABLE marketplace_orders ADD COLUMN IF NOT EXISTS pending_release_at TIMESTAMPTZ;
ALTER TABLE marketplace_orders ADD COLUMN IF NOT EXISTS seller_pending_credited BOOLEAN NOT NULL DEFAULT FALSE;

-- 3. 新建 disputes 表
CREATE TABLE IF NOT EXISTS disputes (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES marketplace_orders(id),
    buyer_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    
    buyer_reason TEXT NOT NULL,
    buyer_images TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    seller_action VARCHAR(20),
    seller_reason TEXT,
    seller_images TEXT,
    seller_responded_at TIMESTAMPTZ,
    
    auto_refunded BOOLEAN NOT NULL DEFAULT FALSE,
    
    admin_id BIGINT,
    admin_decision VARCHAR(20),
    admin_note TEXT,
    resolved_at TIMESTAMPTZ,
    
    refund_amount BIGINT NOT NULL DEFAULT 0,
    
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_disputes_buyer ON disputes(buyer_id);
CREATE INDEX IF NOT EXISTS idx_disputes_seller ON disputes(seller_id);
CREATE INDEX IF NOT EXISTS idx_disputes_status ON disputes(status);
CREATE INDEX IF NOT EXISTS idx_disputes_order ON disputes(order_id);

-- 4. 为已有完成的订单设置 pending_release_at（7天后）
UPDATE marketplace_orders 
SET pending_release_at = NOW() + INTERVAL '7 days'
WHERE status = 'completed' AND pending_release_at IS NULL;
