CREATE TABLE apps (
    id BIGSERIAL PRIMARY KEY,
    app_name VARCHAR(128) NOT NULL,
    app_description TEXT,
    client_id VARCHAR(64) NOT NULL UNIQUE,
    client_secret VARCHAR(128) NOT NULL,
    redirect_uris JSONB NOT NULL,
    scopes JSONB NOT NULL,
    min_trust_level SMALLINT NOT NULL DEFAULT 0,
    fee_rate DECIMAL(4,2) NOT NULL DEFAULT 10.00,
    bound_user_id BIGINT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_codes (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(128) NOT NULL UNIQUE,
    client_id VARCHAR(64) NOT NULL,
    discourse_user_id BIGINT NOT NULL,
    scopes JSONB NOT NULL,
    redirect_uri VARCHAR(512) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_codes_expires ON auth_codes(expires_at);

CREATE TABLE access_tokens (
    id BIGSERIAL PRIMARY KEY,
    access_token VARCHAR(256) NOT NULL UNIQUE,
    refresh_token VARCHAR(128) UNIQUE,
    client_id VARCHAR(64) NOT NULL,
    discourse_user_id BIGINT NOT NULL,
    scopes JSONB NOT NULL,
    access_expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_access_tokens_expires ON access_tokens(access_expires_at);
CREATE INDEX idx_access_tokens_user ON access_tokens(discourse_user_id);
