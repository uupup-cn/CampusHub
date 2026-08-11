#!/bin/bash
set -e
export PATH="/home/pc/go-install/go/bin:$PATH"
PGPASSWORD=chb_test_2024 psql -h localhost -p 5433 -U chb -d chb_test <<'SQL'
INSERT INTO apps (app_name, client_id, client_secret, redirect_uris, scopes, min_trust_level, fee_rate, status)
VALUES ('TestApp', 'test_client_id', 'test_client_secret',
        '["http://localhost:3000/callback"]',
        '["profile:read","chb:read","chb:spend"]',
        0, 10.00, 'active')
ON CONFLICT (client_id) DO NOTHING;
SQL
echo "app seeded"
