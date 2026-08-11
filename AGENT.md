# CampusHub 规则文档

> 项目最高规则文档，所有开发者和贡献者必须遵循。

## 文档作用

本文档是项目的单一信息源，涵盖：项目简介、技术栈、架构、数据库、安全、编码规范、开发命令、开发规则、测试标准、插件设计、部署运维。所有功能开发和架构决策必须遵循本文档。当其他文档与本文档冲突时，以本文档为准。

> 引用场景：项目启动、新成员入职、代码审查、制定开发计划时阅读。

## 1. 项目简介

CampusHub 是一个以知识共享为核心的论坛生态系统，参考 linux.do 的社区运营模式，构建包含论坛、积分体系（CHB）、积分集市、OAuth2 应用接入的完整生态。

- **CHB**：积分单位，总量恒定 5 亿，不支持充值或兑换法币
- **公共积分池**：5 亿 CHB，用户通过发帖、回复、签到、被点赞获取积分
- **官方积分池**：通过收取手续费获得收益，与公共积分池逻辑隔绝
- **积分集市**：用户申请入驻后可发布积分商品，实现积分流转
- **OAuth2 应用接入**：内部应用通过 OAuth2 授权使用论坛积分，后期开放平台
- **论坛定位**：知识共享，无积分付费设计，无悬赏提问

## 2. 技术栈

| 层级 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 论坛引擎 | Discourse | 最新稳定版 | 论坛主站，仅通过插件扩展 |
| 数据库 | PostgreSQL | >= 15 | 核心后端数据存储 |
| 后端语言 | Go | >= 1.21 | 积分账本、奖励引擎、IdP、集市、管理后台 API |
| 前端框架 | Next.js | >= 14 | 集市前端、管理后台前端 |
| 前端动画 | framer-motion | 最新 | 入场、交错、滚动揭示等动效，全站唯一动效方案 |
| 前端图标 | lucide-react | 最新 | 全部图标来源，禁止 emoji 及其他图标库 |
| 容器 | Docker | 最新 | Discourse 部署 |
| 反向代理 | Nginx | 最新 | 负载均衡、SSL 终止 |

## 3. 开源与自研边界

| 程序 | 引用方式 | 修改边界 |
|------|---------|---------|
| Discourse | Docker 部署 + 插件 | 仅通过插件扩展，不修改核心代码 |
| PostgreSQL | 独立部署 | 直接使用，不自研 |
| Nginx | 反向代理 | 配置文件自定义 |
| 核心后端服务 (Go) | 自研 | 完全自主开发 |
| 奖励插件 (Ruby) | 自研 | 完全自主开发 |
| 集市/管理后台 (Next.js) | 自研 | 完全自主开发 |

**边界原则**：Discourse 不做积分存储（只触发事件）；后端服务不做前端渲染；集市和管理后台不直接操作数据库（通过 API 访问）；插件不做业务逻辑（只负责事件转发）。

**开发环境原则**：Discourse 开源代码不下载到本地，仅通过 Docker 官方镜像运行。自研的 Go 后端、Next.js 前端、Ruby 插件全部在本地开发。Discourse 的 Docker 容器支持挂载本地插件目录，无需每次重新构建。

## 4. 项目目录结构

```
xiaoyuan/
├── AGENT.md               # 规则文档（本文档）
├── README.md              # 文档导航
├── docs/
│   ├── 01-经济模型.md         # 经济模型
│   ├── 02-API文档.md         # API 参考
│   └── 03-开发计划.md         # 开发计划
├── chb-backend/           # Go 后端服务（已实现：账本/奖励/IdP/集市/管理后台 API）
│   ├── config.yaml        # 开发配置
│   ├── config.test.yaml   # 集成测试配置
│   ├── integration_test.go# 集成测试（26 个用例）
│   └── migrations/        # 8 组迁移文件
├── chb-frontend/          # Next.js 前端（已实现：首页/集市/管理后台/授权页，14 条路由）
│   ├── app/
│   │   ├── page.tsx       # 首页（沉浸式 Awwwards 级着陆页）
│   │   ├── marketplace/   # 集市：列表、详情、入驻申请、我的商品、订单
│   │   ├── admin/         # 管理后台：仪表盘、审核、用户、应用、配置、审计
│   │   └── oauth/         # OAuth2 授权同意页
│   ├── src/components/    # Navbar、Toast、AuroraBackground、Reveal（滚动揭示）
│   └── src/lib/api.ts     # API 客户端（X-User-ID 模拟）
├── chb-reward-plugin/     # Discourse Ruby 插件（已实现：事件监听+签到）
├── docker-compose.yml     # Docker 编排配置
└── .gitignore
```

## 5. 系统架构

~~~mermaid
graph TB
    B["Browser"] --> N["Nginx"]
    N --> D["Discourse Docker"]
    N --> F["Next.js 前端"]
    N --> API["Go 后端 API"]
    API --> PG[("PostgreSQL")]
    D -->|"事件"| P["chb_reward_plugin"]
    P -->|"HTTP"| API
    F -->|"OAuth2"| API
~~~

| 模块 | 职责 |
|------|------|
| Discourse | 论坛核心、用户认证、信任等级、内容管理 |
| chb_reward_plugin | 监听 Discourse 事件，转发到后端奖励引擎 |
| Go 后端 | IdP + 积分账本 + 奖励引擎 + 集市 API + 管理后台 API |
| Next.js 前端 | 集市前端 + 管理后台前端 + OAuth2 授权页 |

**完整功能清单**参见 [03-开发计划](docs/03-开发计划.md)。

## 6. 数据库设计

16 张表，金额字段统一 BIGINT（单位 0.01 CHB），时区 Asia/Shanghai。

| 表名 | 用途 | 关键字段 |
|------|------|---------|
| pools | 公共池 + 官方池余额 | pool_type, total_supply, balance |
| user_balances | 用户余额和等级 | discourse_user_id, balance, version, trust_level |
| transactions | 积分变动流水（含幂等键） | tx_type, amount, fee, idempotency_key (UNIQUE) |
| reward_rules | 奖励规则配置 | action, amount, cooldown_seconds, daily_cap |
| trust_level_caps | 等级日上限 | trust_level, daily_cap, reward_multiplier |
| daily_reward_quotas | 用户每日已获取量 | discourse_user_id, reward_date, earned_today |
| reward_logs | 奖励发放日志 | action, ref_type, ref_id, discourse_user_id, amount |
| apps | OAuth2 应用 | client_id, scopes, min_trust_level, fee_rate, bound_user_id |
| auth_codes | 授权码 | code, client_id, expires_at, used |
| access_tokens | OAuth2 Token | access_token, refresh_token, scopes, expires_at |
| marketplace_items | 集市商品 | seller_id, title, price, stock, status |
| marketplace_orders | 集市订单 | order_no, item_id, buyer_id, seller_id, total_amount, fee |
| merchant_applications | 入驻申请 | discourse_user_id, shop_name, status |
| release_logs | 释放阀操作 | amount, from_pool, to_pool, operator_id |
| system_configs | 系统配置 | config_key, config_value |
| audit_logs | 操作审计 | operator_id, action, target_type, target_id, detail |

**核心约束**：transactions 表的 idempotency_key 唯一索引保证幂等；reward_logs 的 (ref_type, ref_id, discourse_user_id) 唯一索引保证不重复发奖；user_balances 的 version 字段用于乐观锁。

**迁移策略**：使用 cmd/migrate/main.go 管理，8 组迁移文件（000001-000008，含 up + down），生产环境迁移前必须备份。

## 7. 安全规范

### 7.1 OAuth2 Scope

| Scope | 权限 | 敏感级别 |
|-------|------|---------|
| profile:read | 读取用户资料 | 低 |
| chb:read | 查询 CHB 余额 | 低 |
| chb:spend | 消费 CHB 积分 | 高 |

### 7.2 鉴权方式

| 鉴权方式 | 请求头 | 说明 |
|---------|--------|------|
| Bearer Token | Authorization: Bearer xxx | OAuth2 access_token |
| API Key | X-API-Key: xxx | 内部服务调用（奖励引擎） |
| 管理后台 Key | X-Admin-Key: xxx | 管理后台 API |
| 用户 ID | X-User-ID: xxx | 开发/测试环境用户模拟 |

### 7.3 幂等与并发控制

- 所有积分扣减操作必须携带 idempotency_key，数据库唯一索引保证幂等
- 用户余额更新使用行级锁（SELECT FOR UPDATE）+ 乐观锁（version 字段）
- 积分池更新使用行级锁，防止并发超发
- 集市下单使用库存行级锁，防止超卖

### 7.4 审计与合规

- 管理后台所有敏感操作写入 audit_logs
- 积分变动全部写入 transactions（含手续费、流向）
- 奖励发放全部写入 reward_logs（含等级、倍率）

## 8. 编码规范

### Go 后端
- 标准 Go project layout；API 响应统一 {code, message, data} 格式
- 错误码：0=成功，1000-1999=参数，2000-2999=权限，3000-3999=积分，4000-4999=集市，5000+=系统
- 参数化查询，禁止 SQL 拼接；金额 int64（0.01 CHB 单位）

### Next.js 前端
- 组件 PascalCase，文件 kebab-case；API 请求统一封装自动携带 token
- 环境变量 NEXT_PUBLIC_ 前缀才暴露给浏览器
- 图标统一使用 lucide-react，界面全程禁止使用 emoji 及 emoji 类符号（如 →、❤️），一律用 Lucide 图标表达
- 设计标准对标 Awwwards / FWA / CSS Design Awards 顶级网站：将浏览器视作交互式艺术画布，追求先锋视觉风格、实验性排版、流畅物理动效与冲击力文字版式
- 动效统一使用 framer-motion（入场、交错、悬浮、滚动揭示），禁止散落手写动画
- 沉浸式体验：结合极光背景、玻璃拟态、渐变光晕与微交互，保证统一完整的精品页面体验
- 视觉规范（CSS 变量）以 app/globals.css 为单一来源，禁止页面内硬编码色值

### Discourse 插件
- 前缀 chb_；使用 DiscourseEvent.on 回调；不修改核心代码

### 中文编码规范
- 所有源文件统一 UTF-8 编码（无 BOM），禁止使用 GBK、GB2312、ISO-8859-1 等非 UTF-8 编码
- 数据库字符集 utf8mb4（PostgreSQL 使用 UTF-8 编码）
- HTML 页面声明 <meta charset="utf-8">
- JSON API 响应头 Content-Type: application/json; charset=utf-8
- 用户可见的错误消息使用中文，日志记录使用英文
- 所有用户输入在写入数据库前统一转码为 UTF-8，防止特殊字符导致乱码
- 迁移 SQL 文件必须保存为 UTF-8 无 BOM 格式（BOM 会导致 PostgreSQL 解析失败）

### 通用
- 提交格式 type(scope): description（feat/fix/docs/refactor/test/chore）
- 分支策略：main 受保护，feature/ 功能分支，hotfix/ 紧急修复
- 格式化：Go 用 gofmt，前端用 Prettier；时区 Asia/Shanghai

## 9. 开发命令

| 模块 | 命令 | 说明 |
|------|------|------|
| Go 后端 | GOPROXY=https://goproxy.cn,direct go build ./cmd/server | 构建（WSL 需国内代理） |
| Go 后端 | go test ./... -v -race | 单元测试（含竞态检测） |
| Go 集成测试 | go test -v -timeout 180s ./... | 集成测试（需先建立 SSH 隧道 ssh -L 5433:localhost:5432 devbox） |
| Next.js | npm run dev / npm run build | 本地开发 / 生产构建 |
| 数据库迁移 | go run ./cmd/migrate | 应用迁移 + 种子数据 |
| Discourse 插件 | ./launcher rebuild app | 容器内部署插件 |

## 10. 开发规则

- main 分支受保护，禁止直接推送；合并需 PR + 代码审查 + CI 通过
- 依赖管理：Go 锁定 go.sum，前端锁定 package-lock.json，禁止未审查的第三方库
- 每一步开发完成后必须同步更新相关文档（API 文档、开发计划）
- 文档与代码不一致时，以代码为准更新文档

## 11. 代码审查与问题修复

**审查清单**：遵循编码规范；积分操作使用幂等键；参数化查询；无硬编码密钥；完整错误处理；足够日志。

**修复流程**：复现 → 定位根因 → 编写测试 → 最小改动修复 → 验证 → 审查合并。

**升级**：P0（积分丢失/多扣）2 小时 hotfix；P1（功能不可用）当天修复；P2（体验问题）排入迭代。

## 12. 测试验收标准

- 后端覆盖率 >= 80%，积分账本 100%，奖励引擎 >= 90%
- 集成测试（P13 已完成）：26 个用例覆盖账本(7)/奖励(4)/集市(4)/OAuth2(4)/管理后台(4)/响应格式与错误码(3)
- 集成测试：OAuth2 全链路、100 并发扣减、防刷、集市交易
- 验收：P0/P1 修复；对账平衡；OAuth2 符合 RFC 6749；API P99 < 200ms

## 13. Discourse 插件

chb_reward_plugin 监听 Discourse 事件，通过 HTTP 转发到后端奖励引擎：

| 事件 | 触发 | 转发目标 |
|------|------|---------|
| topic_created | 发帖 | POST /api/chb/reward (action=post) |
| post_created | 回复 | POST /api/chb/reward (action=reply) |
| like_added | 被点赞 | POST /api/chb/reward (action=liked) |
| user_trust_level_change | 等级变更 | POST /api/chb/sync/trust-level |
| 签到按钮 | 用户点击 | POST /api/chb/checkin（插件路由代理） |

插件目录结构：plugin.rb（入口）+ config/settings.yml（配置 API 地址和密钥）+ lib/（事件处理、API 客户端、序列化、等级同步）+ assets/（签到前端脚本）+ spec/（单元测试）。

**签到流程**：论坛页面按钮 → AJAX /chb/checkin（插件路由）→ 插件转发 /api/chb/checkin → 后端发放奖励。

## 14. 部署运维

**部署架构**：Nginx 反向代理 → Discourse Docker + Go 后端 x2 + Next.js 前端 x2 → PostgreSQL 主从 + Redis。

**服务器规格**：Discourse 4C8G，Go 后端/前端 2C4G x2，PostgreSQL 4C8G。

**备份策略**：数据库每日 03:00 pg_dump 全量保留 30 天 + WAL 连续归档支持 PITR。恢复 RTO：主节点故障 30 分钟，完全损坏 2 小时。

**监控指标**：公共池水位 < 80% 告警；API P99 > 200ms 告警；错误率 > 1% 告警；服务器 CPU/内存/磁盘 > 80% 告警。

## 15. 文档导航

| 场景 | 引用文档 |
|------|---------|
| 理解积分经济规则 | [01-经济模型](docs/01-经济模型.md) |
| 查看 API 路由和参数 | [02-API文档](docs/02-API文档.md) |
| 查看开发里程碑 | [03-开发计划](docs/03-开发计划.md) |
| 了解项目全貌 | 本文档 §1-§5 |
| 数据库表结构 | 本文档 §6 |
| 安全设计 | 本文档 §7 |
| 编码规范 | 本文档 §8 |
| 测试标准 | 本文档 §12 |
| 部署配置 | 本文档 §14 |
| 环境边界和操作红线 | 本文档 §16 |

## 16. 环境边界

### 16.1 开发环境（WSL2）

| 属性 | 说明 |
|------|------|
| 用途 | 日常编码、本地编译、单元测试、前端开发调试 |
| 访问方式 | VS Code Remote - WSL 连接 |
| 代码位置 | 项目文件位于 Windows 文件系统，WSL 通过 /mnt/c/ 访问 |
| 版本控制 | 所有 Git 操作在此环境执行 |
| 运行服务 | PostgreSQL（Docker 容器）、Go 后端（开发模式热重载）、Next.js 前端（dev server） |
| 不运行 | Discourse Docker（内存不足，仅通过服务器测试） |
| 修改边界 | 仅修改自研代码（chb-backend/、chb-frontend/、chb-reward-plugin/），不修改 Discourse 核心代码 |
| 安装的运行时 | Go、Node.js、npm、Docker CLI（通过 Docker Desktop WSL2 后端） |
| 数据 | 测试数据库数据可随时重建，不做持久化存储 |

### 16.2 部署测试环境（devbox 服务器）

| 属性 | 说明 |
|------|------|
| 用途 | 集成测试、预发布验证、Discourse 插件测试、最终部署 |
| 访问方式 | SSH（已配置别名 devbox） |
| IP | 159.75.116.207 |
| 用户名 | ubuntu |
| 操作系统 | Ubuntu 24.04 LTS |
| 代码来源 | 仅从 Git 仓库拉取，禁止直接修改代码 |
| 配置管理 | 配置文件通过部署脚本管理，不手动编辑 |
| 运行服务 | Discourse Docker、PostgreSQL、Go 后端（生产二进制）、Next.js 前端（生产构建）、Nginx |
| 不运行 | 开发服务器（npm run dev、go run），不安装编译工具链 |
| 修改边界 | 仅修改部署配置（nginx.conf、docker-compose.yml、.env）、系统配置（systemd service），不修改业务代码 |
| 数据持久化 | PostgreSQL 数据持久化，定期备份 |
| 数据安全 | 禁止在服务器上运行测试数据生成脚本，防止污染生产数据 |

### 16.3 环境协作流程

代码仓库：https://github.com/uupup-cn/CampusHub（公开，默认分支 main，已启用分支保护：禁止直推、强制 PR + 1 人审查、禁止强推与删除）

WSL2 开发 → Git 提交 → GitHub 仓库（PR 审查）→ devbox 服务器拉取部署 → 测试验证 → 生产运行

### 16.4 操作红线

| 操作 | 禁止原因 |
|------|---------|
| 在服务器上直接编辑 Go 或前端代码 | 版本失控，无法追溯 |
| 在服务器上运行 npm install / go build | 消耗服务器资源，应该用 CI/CD 构建 |
| 在 WSL 中运行 Discourse Docker | 内存不足，影响开发体验 |
| 在服务器上执行 git reset --hard | 可能导致部署状态与代码库不一致 |
| 手动修改服务器数据库 | 应通过迁移脚本或 API 操作 |
| 将服务器数据库密码硬编码到代码中 | 安全风险，应使用环境变量或密钥管理 |
