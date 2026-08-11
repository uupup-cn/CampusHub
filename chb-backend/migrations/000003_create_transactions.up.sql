CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    tx_type VARCHAR(32) NOT NULL,
    discourse_user_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    fee BIGINT NOT NULL DEFAULT 0,
    net_amount BIGINT NOT NULL,
    from_type VARCHAR(20) NOT NULL,
    to_type VARCHAR(20) NOT NULL,
    from_id BIGINT,
    to_id BIGINT,
    idempotency_key VARCHAR(128) NOT NULL,
    ref_type VARCHAR(32),
    ref_id BIGINT,
    description VARCHAR(512),
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_tx_idempotency_key ON transactions(idempotency_key);
CREATE INDEX idx_tx_user_id ON transactions(discourse_user_id);
CREATE INDEX idx_tx_type ON transactions(tx_type);
CREATE INDEX idx_tx_ref ON transactions(ref_type, ref_id);
CREATE INDEX idx_tx_created_at ON transactions(created_at);
