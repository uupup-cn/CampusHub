CREATE TABLE system_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(128) NOT NULL UNIQUE,
    config_value JSONB NOT NULL,
    description VARCHAR(512),
    updated_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE release_logs (
    id BIGSERIAL PRIMARY KEY,
    amount BIGINT NOT NULL,
    from_pool VARCHAR(20) NOT NULL,
    to_pool VARCHAR(20) NOT NULL,
    operator_id BIGINT NOT NULL,
    reason VARCHAR(512),
    is_auto BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_release_logs_created ON release_logs(created_at);
