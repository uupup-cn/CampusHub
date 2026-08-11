CREATE TABLE user_balances (
    id BIGSERIAL PRIMARY KEY,
    discourse_user_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(60) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    version BIGINT NOT NULL DEFAULT 1,
    trust_level SMALLINT NOT NULL DEFAULT 0,
    total_earned BIGINT NOT NULL DEFAULT 0,
    total_spent BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_balances_user_id ON user_balances(discourse_user_id);
CREATE INDEX idx_user_balances_status ON user_balances(status);
