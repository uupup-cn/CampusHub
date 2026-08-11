CREATE TABLE reward_rules (
    id BIGSERIAL PRIMARY KEY,
    action VARCHAR(32) NOT NULL UNIQUE,
    display_name VARCHAR(64) NOT NULL,
    amount BIGINT NOT NULL,
    cooldown_seconds INT NOT NULL DEFAULT 0,
    daily_cap_per_user BIGINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trust_level_caps (
    id BIGSERIAL PRIMARY KEY,
    trust_level SMALLINT NOT NULL UNIQUE,
    daily_cap BIGINT NOT NULL,
    reward_multiplier DECIMAL(3,2) NOT NULL DEFAULT 1.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE daily_reward_quotas (
    id BIGSERIAL PRIMARY KEY,
    discourse_user_id BIGINT NOT NULL,
    reward_date DATE NOT NULL,
    earned_today BIGINT NOT NULL DEFAULT 0,
    action_counts JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_daily_quota_user_date ON daily_reward_quotas(discourse_user_id, reward_date);
CREATE INDEX idx_daily_quota_date ON daily_reward_quotas(reward_date);

CREATE TABLE reward_logs (
    id BIGSERIAL PRIMARY KEY,
    discourse_user_id BIGINT NOT NULL,
    action VARCHAR(32) NOT NULL,
    amount BIGINT NOT NULL,
    ref_id BIGINT NOT NULL,
    ref_type VARCHAR(32) NOT NULL,
    trust_level SMALLINT NOT NULL,
    multiplier DECIMAL(3,2) NOT NULL,
    ip_address VARCHAR(45),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    reject_reason VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reward_logs_ref ON reward_logs(ref_type, ref_id, discourse_user_id);
CREATE INDEX idx_reward_logs_user ON reward_logs(discourse_user_id);
CREATE INDEX idx_reward_logs_action ON reward_logs(action, created_at);
