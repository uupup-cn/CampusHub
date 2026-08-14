-- CampusHub 测试数据
-- 创建 5 个非管理员用户，覆盖不同业务状态
-- 注意：使用 SELECT 子查询引用 order_id，兼容已有数据场景

-- ============================================================
-- 1. 用户基础数据
-- ============================================================

-- 用户 1001: 张三 - 入驻申请（待审核）
INSERT INTO user_balances (discourse_user_id, username, balance, version, trust_level, total_earned, total_spent, status)
VALUES (1001, 'zhangsan', 10000, 1, 1, 5000, 0, 'active')
ON CONFLICT (discourse_user_id) DO NOTHING;

INSERT INTO merchant_applications (discourse_user_id, shop_name, description, status)
VALUES (1001, '张三的二手书店', '主要出售二手教材和参考书', 'pending')
ON CONFLICT DO NOTHING;

-- 用户 1002: 李四 - 入驻已通过，商品待审核
INSERT INTO user_balances (discourse_user_id, username, balance, version, trust_level, total_earned, total_spent, status)
VALUES (1002, 'lisi', 50000, 1, 2, 20000, 500, 'active')
ON CONFLICT (discourse_user_id) DO NOTHING;

INSERT INTO merchant_applications (discourse_user_id, shop_name, description, status, reviewed_by, review_comment)
VALUES (1002, '李四的数码铺', '出售各类数码配件和电子产品', 'approved', 1, '审核通过，欢迎入驻')
ON CONFLICT DO NOTHING;

INSERT INTO marketplace_items (seller_id, title, description, price, stock, status, category)
VALUES
  (1002, '二手 iPad 保护壳', '使用一个月，成色很新', 500, 3, 'pending', '数码配件'),
  (1002, 'USB-C 扩展坞', '支持 HDMI/PD/USB3.0', 1500, 5, 'pending', '数码配件')
ON CONFLICT DO NOTHING;

-- 用户 1003: 王五 - 入驻已通过，商品已上架
INSERT INTO user_balances (discourse_user_id, username, balance, version, trust_level, total_earned, total_spent, status)
VALUES (1003, 'wangwu', 80000, 2, 2, 30000, 10000, 'active')
ON CONFLICT (discourse_user_id) DO NOTHING;

INSERT INTO merchant_applications (discourse_user_id, shop_name, description, status, reviewed_by, review_comment)
VALUES (1003, '王五的学习资料', '各类学习资料和笔记分享', 'approved', 1, '审核通过')
ON CONFLICT DO NOTHING;

INSERT INTO marketplace_items (seller_id, title, description, price, stock, status, category)
VALUES
  (1003, '高等数学(上) 第八版', '几乎全新，带笔记', 2000, 1, 'active', '教材'),
  (1003, '考研英语词汇书', '红宝书，附赠App激活码', 1500, 2, 'active', '教材'),
  (1003, '计算机组成原理笔记', '手写笔记，重点清晰', 800, 10, 'active', '资料')
ON CONFLICT DO NOTHING;

-- 用户 1004: 赵六 - 买家，有购买订单
INSERT INTO user_balances (discourse_user_id, username, balance, version, trust_level, total_earned, total_spent, status)
VALUES (1004, 'zhaoliu', 20000, 1, 1, 15000, 3000, 'active')
ON CONFLICT (discourse_user_id) DO NOTHING;

-- 订单: 赵六购买王五的商品
INSERT INTO marketplace_orders (order_no, item_id, buyer_id, seller_id, quantity, unit_price, total_amount, fee, net_amount, status)
SELECT 'TEST20260814001', id, 1004, seller_id, 1, price, price, price/10, price - price/10, 'completed'
FROM marketplace_items WHERE title = '高等数学(上) 第八版' AND seller_id = 1003
ON CONFLICT (order_no) DO NOTHING;

INSERT INTO marketplace_orders (order_no, item_id, buyer_id, seller_id, quantity, unit_price, total_amount, fee, net_amount, status)
SELECT 'TEST20260814002', id, 1004, seller_id, 2, price, price*2, price*2/10, price*2 - price*2/10, 'pending'
FROM marketplace_items WHERE title = '考研英语词汇书' AND seller_id = 1003
ON CONFLICT (order_no) DO NOTHING;

-- 用户 1005: 孙七 - 买家，有争议订单
INSERT INTO user_balances (discourse_user_id, username, balance, version, trust_level, total_earned, total_spent, status)
VALUES (1005, 'sunqi', 5000, 1, 0, 2000, 1500, 'active')
ON CONFLICT (discourse_user_id) DO NOTHING;

-- 订单: 孙七购买，有争议
INSERT INTO marketplace_orders (order_no, item_id, buyer_id, seller_id, quantity, unit_price, total_amount, fee, net_amount, status, dispute_status, pending_release_at)
SELECT 'TEST20260814003', id, 1005, seller_id, 1, price, price, price/10, price - price/10, 'disputed', 'open', NOW() + INTERVAL '7 days'
FROM marketplace_items WHERE title = '计算机组成原理笔记' AND seller_id = 1003
ON CONFLICT (order_no) DO NOTHING;

-- 争议记录（使用子查询引用 order_id，避免硬编码）
INSERT INTO disputes (order_id, buyer_id, seller_id, status, buyer_reason, created_at)
SELECT mo.id, 1005, 1003, 'pending', '资料内容与描述不符，缺少关键章节', NOW()
FROM marketplace_orders mo WHERE mo.order_no = 'TEST20260814003'
ON CONFLICT DO NOTHING;

-- 应用数据
INSERT INTO apps (app_name, app_description, client_id, client_secret, redirect_uris, scopes, min_trust_level, fee_rate, status)
VALUES
  ('校园二手市场', '校内二手商品交易平台', 'campus_market_client', 'secret_market_key', '["http://localhost:3000/callback","https://market.example.com/callback"]', '["profile:read","chb:read","chb:spend"]', 0, 10.00, 'active'),
  ('学习助手', '学习资料分享与积分激励', 'study_helper_client', 'secret_study_key', '["http://localhost:3000/study/callback"]', '["profile:read","chb:read"]', 1, 5.00, 'active')
ON CONFLICT (client_id) DO NOTHING;

SELECT 'Test data seeded successfully' AS result;
