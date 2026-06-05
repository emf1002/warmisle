# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## 项目概述

暖屿 (WarmIsle) V1 — 面向家庭的私有化协作平台。包含记账本、待办管理、愿望清单、家庭论坛、仪表盘 5 大模块。Go 后端 + Vue 3 前端，编译为单二进制部署。

## 构建与开发命令

```bash
make dev       # 启动开发环境（后端 air 热重载 + 前端 vite dev）
make build     # 构建单二进制（前端 dist → embed 到 Go binary）
make clean     # 清理构建产物
```

前端开发代理：Vite dev server (port 3000) 代理 `/api` 到后端 `localhost:8080`。

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.22+, Gin, GORM, SQLite (WAL 模式), goose (数据库迁移), JWT (golang-jwt/v5) |
| 前端 | Vue 3 (Composition API), Ant Design Vue, Vite, Pinia, Axios, vue-router 4 |

## 项目结构

```
warmisle/
├── backend/
│   ├── cmd/server/main.go       # Web 服务入口，embed 前端 dist
│   ├── cmd/cli/main.go          # CLI 工具（reset-password）
│   ├── internal/
│   │   ├── handler/             # Gin 处理器（HTTP 层）
│   │   ├── service/             # 业务逻辑层
│   │   ├── repository/          # GORM 数据访问层
│   │   ├── model/               # 数据模型
│   │   ├── middleware/          # JWT 认证、CORS 中间件
│   │   ├── routes/              # 路由注册
│   │   └── pkg/                 # 工具包（jwt/password/response/database）
│   └── migrations/              # goose SQL 迁移文件
├── frontend/
│   └── src/
│       ├── views/               # 页面组件（按模块：auth/dashboard/ledger/todo/wish/forum/member/category/profile）
│       ├── components/          # 通用组件（EmptyState 等）
│       ├── layouts/             # 布局组件（MainLayout: 桌面侧边栏+移动端底栏, AuthLayout）
│       ├── stores/              # Pinia 状态管理
│       ├── router/              # 路由配置（含守卫：auth + init 检测）
│       ├── api/                 # Axios 请求封装
│       └── utils/               # 工具函数（request.ts: Axios 拦截器）
├── docs/prd.md                  # 产品需求文档（权威需求来源）
├── docs/superpowers/specs/      # 技术设计 + UI 风格指南
├── docs/superpowers/plans/      # 实施计划（包含 task 依赖关系）
└── Makefile
```

前端构建输出到 `backend/frontend/dist/`，通过 `//go:embed` 嵌入 Go 二进制。

## 架构分层

**后端三层架构**：`handler（HTTP） → service（业务逻辑） → repository（数据访问）`

- Handler：参数绑定与校验，调用 service，返回统一响应格式 `{code, message, data}`
- Service：业务规则校验、权限检查、调用 repository
- Repository：GORM 数据库操作，不含业务逻辑

**统一响应工具** (`backend/internal/pkg/response.go`)：
- `Success(c, data)` → `{code: 0, message: "ok", data: ...}`
- `Error(c, httpCode, bizCode, msg)` → 自定义错误码
- 分页响应含 `list`, `total`, `page`, `page_size`

**API 约定**：基础路径 `/api`，认证 `Authorization: Bearer <token>`

## 核心设计决策（必须遵守）

1. **金额存储**：数据库以"分"(int64) 存储，前端以"元"展示（除以 100）。金额始终正整数，收支方向由分类 type 实时决定。
2. **收支类型**：记账记录不存储 type 字段，方向由关联分类 `category.type`（income/expense）实时引用。修改分类 type 时关联记录同步变化。
3. **删除策略**：全系统统一软删除（`deleted_at` 字段），无硬删除。GORM 查询自动过滤 `deleted_at IS NULL` 的记录。
4. **并发策略**：V1 采用 Last Write Wins，不引入乐观锁。SQLite WAL 模式保证写事务串行化。
5. **JWT**：7 天有效期，V1 不做 refresh token。修改密码后当前 Token 立即失效。`HC_JWT_SECRET` 不设置时自动生成随机 32 位 hex。
6. **登录锁定**：5 次失败锁定 15 分钟（内存记录，服务重启清空）。CLI `reset-password` 同步清除锁定。
7. **评论嵌套**：最多 2 层（一级评论 + 二级回复）。删除一级评论时二级同步软删除。
8. **预置数据**：由 goose migration 创建（支出 15 个 + 收入 5 个分类，10 个标签），迁移时用 `INSERT OR IGNORE` 保证幂等。用户已删除的预置条目不会因迁移恢复。
9. **数据库迁移**：启动时自动执行 goose Up 迁移（不支持降级），迁移前自动备份到 `backups/`，迁移失败则停止启动。
10. **角色**：admin（管理员）和 member（普通成员）。管理员至少一个，禁止删除/禁用最后一个管理员。

## 配置（环境变量）

| 变量 | 说明 | 默认值 |
|---|---|---|
| `HC_JWT_SECRET` | JWT 签名密钥 | 随机生成 32 位 hex（持久化到 `data/secret.key`） |
| `HC_DB_PATH` | SQLite 文件路径 | `./data/warmisle.db` |
| `HC_PORT` | 监听端口 | `8080` |
| `HC_CORS_ORIGINS` | CORS 白名单（逗号分隔） | 开发：`*`，生产（嵌入前端）：仅同源 |

## 前端路由与布局

```
/login       → 登录页（AuthLayout）
/init        → 初始化页（AuthLayout）
/            → 仪表盘
/ledger      → 记账本
/todo        → 待办管理
/wish        → 愿望清单
/forum       → 家庭论坛（信息流 + 话题详情 /forum/topic/:id）
/members     → 成员管理（admin only）
/categories  → 分类管理（admin only）
/profile     → 个人中心
```

桌面端（≥768px）：侧边栏导航 + 顶部栏（头像下拉菜单）+ 主内容区
移动端（<768px）：底部 TabBar（仪表盘/记账/待办/论坛/我的）+ 顶部栏

## UI 设计关键规范

- 金额：收入绿色 `+¥` 前缀，支出红色 `-¥` 前缀
- 所有可点击元素最小触控高度 44px（移动端）
- 删除操作统一使用 `<Modal>` 二次确认
- 空状态使用通用 `<EmptyState>` 组件，区分 `no-data`（有创建引导）和 `no-result`（有清除筛选链接）
- 日期格式：今天/"昨天"/"M月D日 周X"
- 头像使用预置 emoji 列表（约 30 个），不支持自定义图片上传

## 实施顺序

按 RICE 优先级：认证与权限 → 成员管理 → 分类管理 → 记账本 → 待办管理 → 仪表盘 → 愿望清单 → 家庭论坛 → 个人中心 → 布局与响应式。

完整实施计划（含 task 依赖关系 DAG）见 `docs/superpowers/plans/2026-05-16-warmisle-implementation.md`。
详细验收标准、业务规则、字段约束见 `docs/prd.md`。
