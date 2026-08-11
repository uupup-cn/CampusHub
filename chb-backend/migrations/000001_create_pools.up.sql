CREATE TABLE pools (
    id BIGSERIAL PRIMARY KEY,
    pool_type VARCHAR(32) NOT NULL UNIQUE,
    total_supply BIGINT NOT NULL,
    balance BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO pools (pool_type, total_supply, balance) VALUES
('public', 50000000000, 50000000000),
('official', 0, 0);
