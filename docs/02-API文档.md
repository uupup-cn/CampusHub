# CampusHub API 文档

> 本文档定义 CampusHub 所有 API 路由、参数、响应格式，包括 Discourse 自身 API 和自研服务 API。

## 文档作用

定义完整的 API 路由表和参数规范，供前端开发、第三方应用集成、联调测试时使用。

> 引用场景：前端开发时调用 API；第三方应用集成时使用 OAuth2 和积分 API；联调测试时核对请求和响应格式。安全规范（→ [AGENT.md §7](../AGENT.md#7-安全规范)）和数据库设计（→ [AGENT.md §6](../AGENT.md#6-数据库设计)）是本文档的输入。

## 1. 通用规范

### 1.1 基础 URL

| 服务 | 基础 URL | 说明 |
|------|---------|------|
| Discourse API | https://forum.campushub.com | 论坛自身 API |
| 核心后端 API | https://api.campushub.com | 自研服务 API |
| 前端页面 | https://campushub.com | 集市和管理后台 |

### 1.2 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 1.3 分页格式

请求参数：?page=1&page_size=20
响应 data：

```json
{
  "items": [],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### 1.4 鉴权方式

| 鉴权方式 | 请求头 | 说明 |
|---------|--------|------|
| Session Cookie | Cookie: _t=xxx | 论坛登录态 |
| Bearer Token | Authorization: Bearer xxx | OAuth2 access_token |
| API Key | X-API-Key: xxx | 内部服务调用（奖励引擎） |
| 管理后台 Key | X-Admin-Key: xxx | 管理后台 API |

### 1.5 错误码

| 错误码 | 说明 | HTTP 状态码 |
|--------|------|------------|
| 0 | 成功 | 200 |
| 1001 | 参数缺失 | 400 |
| 1002 | 参数格式错误 | 400 |
| 1003 | 资源不存在 | 404 |
| 2001 | 未登录 | 401 |
| 2002 | 权限不足 | 403 |
| 2003 | Token 无效 | 401 |
| 2004 | Token 过期 | 401 |
| 2005 | 等级不足 | 403 |
| 2006 | Scope 不足 | 403 |
| 3001 | 余额不足 | 400 |
| 3002 | 积分池余额不足 | 500 |
| 3003 | 超出日上限 | 400 |
| 3004 | 冷却中 | 429 |
| 3005 | 重复请求（幂等） | 200 |
| 3006 | 账户已冻结 | 403 |
| 4001 | 商品已售罄 | 400 |
| 4002 | 商品待审核 | 400 |
| 4003 | 入驻申请已存在 | 400 |
| 4004 | 非商家 | 403 |
| 5001 | 系统内部错误 | 500 |
| 5002 | 数据库错误 | 500 |
| 5003 | 服务不可用 | 503 |

---

## 2. Discourse 自身 API（仅列出相关部分）

### 2.1 用户信息

| 接口 | 方法 | 路由 | 说明 |
|------|------|------|------|
| 获取当前用户 | GET | /session/current.json | 获取当前登录用户信息 |
| 获取用户详情 | GET | /u/{username}.json | 获取指定用户信息 |
| 获取用户信任等级 | GET | /u/{username}.json | 返回 trust_level 字段 |

### 2.2 帖子相关

| 接口 | 方法 | 路由 | 说明 |
|------|------|------|------|
| 获取帖子详情 | GET | /t/{topic_id}.json | 获取主题帖详情 |
| 获取回复详情 | GET | /posts/{post_id}.json | 获取回复详情 |
| 获取点赞列表 | GET | /posts/{post_id}/likes.json | 获取点赞列表 |

### 2.3 OAuth2 端点

| 接口 | 方法 | 路由 | 说明 |
|------|------|------|------|
| 授权端点 | GET | /oauth/authorize | OAuth2 授权（自研 IdP） |
| Token 端点 | POST | /oauth/token | OAuth2 Token 交换（自研 IdP） |
| 用户信息 | GET | /oauth/userinfo | OAuth2 用户信息（自研 IdP） |

> 注意：CampusHub 的 OAuth2 应用接入使用自研 IdP 服务（Go 后端实现），不依赖 Discourse 的 OAuth2。以下路由属于自研 IdP。

### 2.4 Webhook 相关

| 接口 | 方法 | 路由 | 说明 |
|------|------|------|------|
| 创建 Webhook | POST | /admin/api/webhooks | 管理后台创建 Webhook |
| Webhook 事件列表 | GET | /admin/api/webhooks/events.json | 查看 Webhook 事件 |

---

## 3. 自研 API - 积分账本

### 3.1 余额查询

```
GET /api/chb/balance
```

**鉴权**：Bearer Token / Session Cookie

**请求参数**：无

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 123,
    "username": "testuser",
    "balance": 10000,
    "trust_level": 2,
    "status": "active",
    "total_earned": 50000,
    "total_spent": 40000
  }
}
```

### 3.2 积分扣减（第三方应用消费）

```
POST /api/chb/spend
```

**鉴权**：Bearer Token（需 chb:spend scope）

**请求参数**：
```json
{
  "amount": 1000,
  "idempotency_key": "app_order_20240101001",
  "description": "购买高级会员"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "transaction_id": 50001,
    "amount": 1000,
    "fee": 100,
    "net_amount": 900,
    "new_balance": 9000,
    "status": "completed"
  }
}
```

**错误码**：3001（余额不足）、3005（重复请求，幂等命中）、2006（Scope 不足）、2005（等级不足）

### 3.3 交易流水查询

```
GET /api/chb/transactions
```

**鉴权**：Bearer Token / Session Cookie

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20，最大 100 |
| tx_type | string | 否 | 筛选类型：reward/spend/transfer/fee |
| start_date | string | 否 | 开始日期 YYYY-MM-DD |
| end_date | string | 否 | 结束日期 YYYY-MM-DD |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 50001,
        "tx_type": "spend",
        "amount": 1000,
        "fee": 100,
        "net_amount": 900,
        "description": "购买高级会员",
        "created_at": "2024-01-01T12:00:00+08:00"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

### 3.4 积分池信息

```
GET /api/chb/pools
```

**鉴权**：管理后台 API Key

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "public_pool": {
      "total_supply": 50000000000,
      "balance": 48000000000,
      "water_level": 0.96
    },
    "official_pool": {
      "balance": 500000000
    }
  }
}
```

### 3.5 释放阀操作

```
POST /api/chb/release
```

**鉴权**：管理后台 API Key

**请求参数**：
```json
{
  "amount": 1000000,
  "reason": "公共池水位低于 80%，手动释放"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "release_id": 1,
    "amount": 1000000,
    "new_public_balance": 49000000000,
    "new_official_balance": 400000000
  }
}
```

### 3.6 对账查询

```
GET /api/chb/audit
```

**鉴权**：管理后台 API Key

**请求参数**：无

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_balance": 50000000000,
    "user_balances_sum": 49500000000,
    "public_pool_balance": 48000000000,
    "official_pool_balance": 500000000,
    "is_balanced": true
  }
}
```

---

## 4. 自研 API - 奖励引擎

### 4.1 奖励发放

```
POST /api/chb/reward
```

**鉴权**：API Key（内部服务调用）

**请求参数**：
```json
{
  "action": "post",
  "discourse_user_id": 123,
  "ref_type": "topic",
  "ref_id": 45678,
  "idempotency_key": "discourse_event_abc123"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "amount": 100,
    "trust_level": 2,
    "multiplier": 1.0,
    "final_amount": 100,
    "earned_today": 500,
    "daily_cap": 1000,
    "status": "completed"
  }
}
```

**拒绝原因**（status=rejected 时）：

| 原因 | 说明 |
|------|------|
| cooldown | 冷却时间内 |
| user_daily_cap_reached | 用户日上限已满 |
| level_daily_cap_reached | 等级日上限已满 |
| account_frozen | 账户已冻结 |
| duplicate | 重复请求（幂等命中） |

### 4.2 信任等级同步

```
POST /api/chb/sync/trust-level
```

**鉴权**：API Key（内部服务调用）

**请求参数**：
```json
{
  "discourse_user_id": 123,
  "trust_level": 3,
  "idempotency_key": "discourse_event_def456"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "old_trust_level": 2,
    "new_trust_level": 3
  }
}
```

### 4.3 查询奖励规则

```
GET /api/chb/reward/rules
```

**鉴权**：管理后台 API Key

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "action": "post",
      "display_name": "发帖",
      "amount": 100,
      "cooldown_seconds": 300,
      "daily_cap_per_user": 500,
      "is_active": true
    }
  ]
}
```

### 4.4 更新奖励规则

```
PUT /api/admin/reward/rules
```

**鉴权**：管理后台 API Key

**请求参数**：
```json
{
  "action": "post",
  "amount": 150,
  "cooldown_seconds": 600,
  "daily_cap_per_user": 1000,
  "is_active": true
}
```

---

## 5. 自研 API - OAuth2 IdP

### 5.1 授权端点

```
GET /oauth/authorize
```

**鉴权**：Session Cookie（论坛登录态）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| client_id | string | 是 | 应用 client_id |
| redirect_uri | string | 是 | 回调 URI |
| response_type | string | 是 | 固定 code |
| scope | string | 是 | 请求的 scope，空格分隔 |
| state | string | 是 | CSRF 保护，回调时原样返回 |

**流程**：
1. 检查用户是否登录（Session Cookie）
2. 校验 client_id 和 redirect_uri
3. 检查用户 trust_level >= 应用的 min_trust_level
4. 等级不足 → 显示"等级不足"页面，跳转 redirect_uri 带 error=insufficient_trust_level
5. 等级满足 → 显示授权同意页，用户确认后发放授权码

### 5.2 授权确认

```
POST /api/oauth/authorize/confirm
```

**鉴权**：Session Cookie（开发/测试环境通过 `X-User-ID` + `X-Trust-Level` 请求头模拟）

**说明**：授权页用户点击"同意授权"后由前端调用，服务端完成 client / redirect_uri / 等级校验并发放授权码，返回带 code 的回调地址。

**请求参数**：
```json
{
  "client_id": "chb_abc123",
  "redirect_uri": "https://app.example.com/callback",
  "response_type": "code",
  "scope": "profile:read chb:read chb:spend",
  "state": "csrf_state_123"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "redirect_uri": "https://app.example.com/callback?code=xxxx&state=csrf_state_123",
    "code": "xxxx",
    "state": "csrf_state_123"
  }
}
```

**错误场景**：

| 场景 | HTTP | code | message |
|------|------|------|---------|
| client_id 无效 | 403 | 2002 | invalid_client |
| redirect_uri 未注册 | 403 | 2002 | invalid_redirect_uri |
| 等级不足 | 403 | 2002 | insufficient_trust_level |
| 参数缺失 | 400 | 1001 | 参数缺失 |
| 参数格式错误 | 400 | 1002 | 参数格式错误 |
| 未登录 | 401 | 2001 | 未授权 |

### 5.3 Token 端点

```
POST /oauth/token
```

**鉴权**：无（client_id + client_secret 在请求体中验证）

**请求参数**（authorization_code）：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grant_type | string | 是 | 固定 authorization_code |
| code | string | 是 | 授权码 |
| redirect_uri | string | 是 | 必须与授权请求一致 |
| client_id | string | 是 | 应用 client_id |
| client_secret | string | 是 | 应用 client_secret |

**请求参数**（refresh_token）：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grant_type | string | 是 | 固定 refresh_token |
| refresh_token | string | 是 | refresh_token |
| client_id | string | 是 | 应用 client_id |
| client_secret | string | 是 | 应用 client_secret |

**响应示例**：
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 7200,
  "refresh_token": "def456...",
  "scope": "profile:read chb:read chb:spend"
}
```

### 5.4 Token 校验

```
POST /oauth/introspect
```

**鉴权**：API Key（内部服务调用）

**请求参数**：
```json
{
  "token": "eyJhbGciOiJSUzI1NiIs..."
}
```

**响应示例**：
```json
{
  "active": true,
  "client_id": "app_abc123",
  "discourse_user_id": 123,
  "username": "testuser",
  "trust_level": 2,
  "scope": "profile:read chb:read chb:spend",
  "exp": 1704100000
}
```

### 5.5 应用信息

```
GET /api/oauth/app-info?client_id=xxx
```

**鉴权**：无（公开接口）

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "app_name": "我的应用",
    "app_description": "应用描述",
    "scopes": ["profile:read", "chb:read", "chb:spend"],
    "min_trust_level": 1,
    "redirect_uris": ["https://app.example.com/callback"]
  }
}
```

---

## 6. 自研 API - 集市

### 6.1 商品列表

```
GET /api/marketplace/items
```

**鉴权**：无（公开接口）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| category | string | 否 | 分类筛选 |
| keyword | string | 否 | 搜索关键词 |
| sort | string | 否 | 排序：price_asc / price_desc / newest |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "seller_id": 42,
        "title": "CampusHub 定制笔记本",
        "description": "精美的 CampusHub 定制笔记本",
        "price": 500,
        "stock": 10,
        "category": "周边",
        "status": "approved",
        "image_url": "https://...",
        "created_at": "2024-01-01T12:00:00+08:00"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 20
  }
}
```

### 6.2 商品详情

```
GET /api/marketplace/items/:id
```

**鉴权**：无（公开接口）

### 6.3 我的商品

```
GET /api/marketplace/items/mine
```

**鉴权**：Session Cookie

**说明**：返回当前卖家发布的全部商品，包含 `pending`（待审核）、`rejected`（已拒绝）等状态，供“我的商品”页展示。不要求商家审核通过即可查询。

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 2,
        "seller_id": 42,
        "title": "待审核商品",
        "price": 300,
        "stock": 5,
        "category": "虚拟",
        "status": "pending",
        "created_at": "2024-01-01T12:00:00+08:00"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

**注意**：路由注册顺序为 `/items` → `/items/mine` → `/items/:id`，`mine` 必须位于参数路由之前，避免被 `:id` 捕获。

### 6.4 创建商品

```
POST /api/marketplace/items
```

**鉴权**：Session Cookie（需是入驻商家）

**请求参数**：
```json
{
  "title": "CampusHub 定制笔记本",
  "description": "精美的 CampusHub 定制笔记本",
  "price": 500,
  "stock": 10,
  "category": "周边",
  "image_url": "https://..."
}
```

### 6.5 申请入驻

```
POST /api/marketplace/apply
```

**鉴权**：Session Cookie

**请求参数**：
```json
{
  "shop_name": "我的店铺",
  "description": "我擅长制作校园周边产品"
}
```


### 6.5a 商家入驻状态

```
GET /api/marketplace/my-status
```

**鉴权**：Bearer Token 或 X-User-ID

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_merchant": false,
    "shop_name": "",
    "status": "none"
  }
}
```

**字段说明**：
- `is_merchant`：是否已入驻（status=approved 时为 true）
- `status`：入驻状态（none/pending/approved/rejected）
- `shop_name`：店铺名称

### 7.1 用户积分流水

```
GET /api/chb/me/transactions
```

**鉴权**：Bearer Token 或 X-User-ID

**查询参数**：
- `page`：页码，默认 1
- `page_size`：每页条数，默认 20
- `type`：流水类型筛选（reward/spend/recover），可选

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 10,
        "tx_type": "reward",
        "discourse_user_id": 1,
        "amount": 10,
        "fee": 0,
        "net_amount": 10,
        "from_type": "pool",
        "to_type": "user",
        "description": null,
        "status": "completed",
        "created_at": "2026-08-13T00:04:56Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 8.1 用户已授权应用列表

```
GET /api/oauth/my-apps
```

**鉴权**：Bearer Token 或 X-User-ID

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "app_name": "CampusHub Test App",
        "client_id": "chb_xxxxxxxx",
        "scopes": "["profile:read", "chb:read", "chb:spend"]",
        "created_at": "2026-08-12T23:34:41Z"
      }
    ]
  }
}
```

### 6.6 创建订单

```
POST /api/marketplace/orders
```

**鉴权**：Session Cookie

**请求参数**：
```json
{
  "item_id": 1,
  "quantity": 2,
  "idempotency_key": "order_20240101001"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_no": "ORD20240101001",
    "total_amount": 1000,
    "fee": 100,
    "net_amount": 900,
    "status": "completed"
  }
}
```

### 6.7 订单列表

```
GET /api/marketplace/orders
```

**鉴权**：Session Cookie

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role | string | 否 | buyer / seller，默认 buyer |
| status | string | 否 | 筛选状态 |
| page | int | 否 | 页码 |

---

## 7. 自研 API - 管理后台

### 7.1 等级日上限配置

```
GET /api/admin/trust-levels
PUT /api/admin/trust-levels
```

**鉴权**：管理后台 API Key

**PUT 请求参数**：
```json
{
  "trust_level": 0,
  "daily_cap": 500,
  "reward_multiplier": 0.50
}
```

### 7.2 应用管理

```
GET /api/admin/apps
```

**鉴权**：管理后台 API Key

**请求参数**：?page=1&page_size=20

```
POST /api/admin/apps
```

**请求参数**：
```json
{
  "app_name": "我的应用",
  "app_description": "应用描述",
  "redirect_uris": ["https://app.example.com/callback"],
  "scopes": ["profile:read", "chb:read", "chb:spend"],
  "min_trust_level": 1,
  "fee_rate": 10.00,
  "bound_user_id": 123
}
```

```
PUT /api/admin/apps/:id
DELETE /api/admin/apps/:id
```

### 7.3 集市审核

```
GET /api/admin/marketplace/applications
PUT /api/admin/marketplace/applications/:id
```

**PUT 请求参数**：
```json
{
  "status": "approved",
  "review_comment": "审核通过"
}
```

```
GET /api/admin/marketplace/items
PUT /api/admin/marketplace/items/:id
```

### 7.4 用户管理

```
GET /api/admin/users?keyword=&page=&page_size=
PUT /api/admin/users/:id/freeze
PUT /api/admin/users/:id/unfreeze
POST /api/admin/users/:id/recover
```

**追回积分 POST 请求参数**：
```json
{
  "amount": 500,
  "reason": "检测到刷积分行为"
}
```

### 7.5 系统配置

```
GET /api/admin/settings
PUT /api/admin/settings
```

**PUT 请求参数**：
```json
{
  "marketplace_fee_rate": 10.00,
  "auto_release_enabled": true,
  "auto_release_threshold": 80,
  "auto_release_ratio": 50,
  "auto_release_monthly_cap": 10000000
}
```

### 7.6 审计日志

```
GET /api/admin/audit-logs
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| action | string | 否 | 筛选操作类型 |
| operator_id | int | 否 | 筛选操作人 |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |
| page | int | 否 | 页码 |

### 7.7 数据统计

```
GET /api/admin/stats
```

**鉴权**：管理后台 API Key

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_users": 1000,
    "active_users_today": 200,
    "total_transactions": 5000,
    "today_transactions": 50,
    "total_marketplace_orders": 100,
    "public_pool_water_level": 0.96,
    "pending_applications": 5,
    "pending_items": 3
  }
}
```

---

## 8. 自研 API - 签到

### 8.1 签到

```
POST /api/chb/checkin
```

**鉴权**：Session Cookie

**请求参数**：无

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "amount": 50,
    "checkin_date": "2024-01-01",
    "streak": 5,
    "status": "completed"
  }
}
```

**错误码**：3003（今日已签到）、3006（账户已冻结）

### 8.2 签到状态查询

```
GET /api/chb/checkin/status
```

**鉴权**：Session Cookie

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "checked_in_today": true,
    "streak": 5,
    "last_checkin_date": "2024-01-01"
  }
}
```

---

## 9. 自研 API - 健康检查

### 9.1 健康检查

```
GET /health
GET /api/health
```

**鉴权**：无（公开接口）

**说明**：用于负载均衡探活、部署健康检查与监控告警，两路由等价。注意：该接口为裸 JSON，不套用统一 `{code, message, data}` 信封。

**响应示例**：
```json
{
  "status": "ok",
  "database": true,
  "version": "1.0.0",
  "service": "chb-backend"
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| status | string | ok（正常）/ degraded（数据库不可用） |
| database | bool | 数据库连通性 |
| version | string | 服务版本 |
| service | string | 服务标识 |

**错误码**：数据库不可用时仍返回 200，`status=degraded`，由监控依据字段值告警。

---

## 10. 文档导航

| 场景 | 引用文档 | 引用内容 |
|------|---------|---------|
| 理解经济模型 | [01-经济模型](01-经济模型.md) | 奖励规则、手续费、池结构 |
| 理解安全规范 | [AGENT.md §7](../AGENT.md#7-安全规范) | 鉴权、幂等、并发控制 |
| 数据库字段参考 | [AGENT.md §6](../AGENT.md#6-数据库设计) | 表字段映射 |
| 开发后端模块 | [03-开发计划](03-开发计划.md) | 模块实现方案 |
| 前端调用 API | [03-开发计划](03-开发计划.md) | 前端页面设计 |


### 9. 争议管理 API

#### 9.1 发起争议

```
POST /api/marketplace/orders/:id/dispute
```

**鉴权**：Bearer Token 或 X-User-ID

**请求参数**：
```json
{
  "reason": "商品与描述不符",
  "images": ["http://xxx/img1.jpg"]
}
```

**约束**：仅买入订单、3天内、未发起过争议

#### 9.2 争议列表

```
GET /api/marketplace/disputes?role=buyer|seller&page=1&page_size=20
```

**鉴权**：Bearer Token 或 X-User-ID

#### 9.3 争议详情

```
GET /api/marketplace/disputes/:id
```

#### 9.4 卖家同意退款

```
PUT /api/marketplace/disputes/:id/accept
```

**鉴权**：Bearer Token 或 X-User-ID（卖家）

#### 9.5 卖家拒绝退款

```
PUT /api/marketplace/disputes/:id/reject
```

**请求参数**：
```json
{
  "reason": "商品已正常交付",
  "images": ["http://xxx/evidence.jpg"]
}
```

#### 9.6 管理员争议列表

```
GET /api/admin/disputes?page=1&page_size=20&status=rejected
```

**鉴权**：X-Admin-Key

#### 9.7 管理员判定

```
PUT /api/admin/disputes/:id/decide
```

**请求参数**：
```json
{
  "decision": "buyer_win",
  "note": "卖家未能提供有效交付证据"
}
```

**decision 枚举**：buyer_win（退款）/ seller_win（卖家胜，48h后转可用）

### 10. 用户统计 API

#### 10.1 首页统计

```
GET /api/chb/me/summary
```

**鉴权**：Bearer Token 或 X-User-ID

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "income_7d": 160,
    "expense_7d": 25,
    "pending_disputes": 0,
    "my_disputes": 0
  }
}
```

**字段说明**：
- `income_7d`：最近7天收入（reward类型交易总和）
- `expense_7d`：最近7天支出（spend类型交易总和）
- `pending_disputes`：待处理争议数（用户作为卖家，status=pending）
- `my_disputes`：我发起的争议数（用户作为买家，所有状态）

#### 10.2 余额查询（更新）

```
GET /api/chb/balance
```

**响应新增字段**：`pending_balance`（未来积分）
