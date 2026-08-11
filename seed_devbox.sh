#!/bin/bash
set -e
cd /mnt/c/Users/PC/Desktop/claudecode/xiaoyuan
export PATH="/home/pc/go-install/go/bin:$PATH"
export GOPROXY=https://goproxy.cn,direct
cd chb-backend
cat <<'SQL' | PGPASSWORD=chb_test_2024 psql -h localhost -p 5433 -U chb -d chb_test
INSERT INTO pools (pool_type, total_supply, balance) VALUES ('public', 50000000000, 50000000000) ON CONFLICT (pool_type) DO NOTHING;
INSERT INTO pools (pool_type, total_supply, balance) VALUES ('official', 0, 0) ON CONFLICT (pool_type) DO NOTHING;
INSERT INTO reward_rules (action, display_name, amount, cooldown_seconds, daily_cap_per_user, is_active) VALUES
  ('post', '发布主题', 100, 0, 500, true),
  ('reply', '回复帖子', 50, 0, 500, true),
  ('checkin', '每日签到', 30, 86400, 30, true),
  ('like_received', '被点赞', 20, 0, 200, true)
ON CONFLICT (action) DO NOTHING;
INSERT INTO trust_level_caps (trust_level, daily_cap, reward_multiplier) VALUES
  (0, 200, 0.50),
  (1, 500, 0.80),
  (2, 1000, 1.00),
  (3, 2000, 1.20),
  (4, 5000, 1.50)
ON CONFLICT (trust_level) DO NOTHING;
INSERT INTO system_configs (config_key, config_value, description) VALUES
  ('marketplace_fee_rate', '0.10', '集市交易手续费率'),
  ('auto_release_enabled', 'true', '自动释放阀开关'),
  ('auto_release_threshold', '80', '公共池水位阈值(%)'),
  ('auto_release_ratio', '10', '自动释放比例(%)'),
  ('auto_release_monthly_cap', '500000000', '月度释放上限')
ON CONFLICT (config_key) DO NOTHING;
SQL
echo "seed done"
