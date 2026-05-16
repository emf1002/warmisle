# 家庭数字中心 V1 — 技术设计文档

| 项目 | 内容 |
|------|------|
| 产品名称 | 家庭数字中心 V1 |
| 版本 | V1.0.1 |
| 创建日期 | 2026-05-16 |
| 状态 | 待实施 |

**修订记录**：

| 版本 | 日期 | 修改内容 |
|------|------|----------|
| V1.0.1 | 2026-05-16 | Grill-me 审查修正：1) CLI reset-password 同步清除登录锁定；2) 预置数据改由 migration 创建，init/setup 仅创建管理员；3) init/setup 补充请求体定义；4) JWT 中间件增加 status=disabled 校验；5) 仪表盘 API 补充 month 查询参数；6) 信息流 API 补充分页参数+响应结构；7) 待办/愿望 API 补充查询参数；8) 评论/点赞 API 补充请求体定义；9) 新增数据库索引策略（3.4 节）|

---

## 1. 技术栈

| 层级 | 技术 |
|------|------|
| 后端语言 | Go |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | SQLite (WAL 模式) |
| 数据库迁移 | goose |
| 认证 | JWT (bcrypt 密码存储) |
| 前端框架 | Vue 3 (Composition API) |
| UI 组件库 | Ant Design Vue |
| 构建工具 | Vite |
| 状态管理 | Pinia |
| HTTP 客户端 | Axios |

---

## 2. 项目目录结构

```
home-center-v1/
├── backend/
│   ├── cmd/
│   │   ├── server/main.go          # Web 服务入口
│   │   └── cli/main.go             # CLI 工具入口（重置管理员密码等）
│   ├── internal/
│   │   ├── handler/                # Gin 处理器
│   │   ├── service/                # 业务逻辑
│   │   ├── repository/             # GORM 数据访问
│   │   ├── model/                  # 数据模型
│   │   ├── middleware/             # JWT、CORS 中间件
│   │   ├── routes/                 # 路由注册
│   │   └── pkg/                    # 工具包(jwt/password/response)
│   ├── migrations/                 # goose SQL 迁移文件
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── views/                  # 页面(按模块)
│   │   ├── components/             # 通用组件
│   │   ├── stores/                 # Pinia
│   │   ├── router/                 # 路由
│   │   ├── api/                    # Axios 请求封装
│   │   └── utils/                  # 工具函数
│   ├── package.json
│   └── vite.config.ts
├── Makefile                        # 统一构建入口
└── README.md
```

### 构建策略

- `make dev`: 同时启动后端 (air 热重载) + 前端 (vite dev)
- `make build`: 构建前端 → `embed.FS` 嵌入 Go 二进制 → 编译单文件
- 前端产物嵌入 `backend/cmd/server/main.go`
- 最终产出单二进制 `home-center`，一条命令启动

### CLI 工具

`./home-center cli reset-password --username admin --password new123`

**功能**：管理员忘记密码或被登录锁定时的线下恢复工具。绕过认证直接操作数据库。

**使用场景**：

| 命令 | 说明 |
|------|------|
| `home-center cli reset-password` | 重置指定成员密码（默认 `home123`） |

**实现要点**：
- 读取相同的 SQLite 数据库路径（与 Web 服务共享配置）
- 直接调用 `service.MemberService.ResetPassword()` 复用业务逻辑
- 不受登录锁定状态影响，可强制重置被锁定账号的密码
- 执行 reset-password 时同步清除该用户名的登录失败计数（解锁账号）
- 命令执行完后退出，不启动 Web 服务

### 2.3 配置管理

**配置方式**：环境变量优先，低于命令行参数。应用启动时从环境变量读取，不存在则使用默认值。

**配置项清单**：

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `HC_JWT_SECRET` | JWT 签名密钥，不设置时启动自动生成 32 位随机 hex | 随机生成 |
| `HC_DB_PATH` | SQLite 数据库文件路径 | `./data/home-center.db` |
| `HC_PORT` | 监听端口 | `8080` |
| `HC_CORS_ORIGINS` | 允许的 CORS Origin，逗号分隔 | `*`（开发）/ 空=仅同源（生产） |

**JWT Secret 说明**：
- 不设置 `HC_JWT_SECRET` 时，首次启动自动生成随机 32 位 hex 密钥并持久化到 `data/secret.key`
- 可通过环境变量显式指定，推荐生产环境使用 `openssl rand -hex 32` 生成
- Secret 变更会导致已有 Token 全部失效

**CORS 说明**：
- 开发环境（未嵌入前端）：默认允许所有 Origin
- 生产环境（嵌入前端）：默认仅允许同源访问，可通过 `HC_CORS_ORIGINS` 配置白名单

---

## 3. 数据库设计

### 3.1 ER 概览

```
members ──1:N── ledgers ──N:M── members (通过 ledger_members)
members ──1:N── todos
members ──1:N── wishes ──N:M── members (通过 wish_votes)
members ──1:N── posts
members ──1:N── topics
members ──1:N── votes ──1:N── vote_options ──1:N── vote_records
members ──1:N── comments (polymorphic: post/topic/wish)
members ──1:N── likes (polymorphic: post/topic/comment)
categories ──1:N── ledgers
tags ──1:N── topics
```

### 3.2 表结构

#### members

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| username | string(20) | 唯一，字母数字下划线 |
| password | string | bcrypt hash |
| name | string(20) | 显示姓名 |
| avatar | string | emoji，默认 👨 |
| role | string | admin / member |
| status | string | active / disabled |
| last_login | datetime | 可选 |
| created_at | datetime | |
| updated_at | datetime | |
| deleted_at | datetime | 软删除 |

#### categories

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| type | string | income / expense |
| name | string(20) | 同 type 下唯一 |
| icon | string | emoji |
| sort_order | int | 排序权重 |
| preset | bool | 是否预置 |
| created_at | datetime | |
| updated_at | datetime | |
| deleted_at | datetime | 软删除 |

**预置数据**：支出 15 个 + 收入 5 个，见 PRD 第 6.3 节。

#### ledgers

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| amount | int64 | 存储单位：分（1元=100分） |
| note | string(200) | 备注 |
| category_id | uint | FK → categories |
| creator_id | uint | FK → members，记录者 |
| occurred_at | datetime | 账单发生时间 |
| created_at | datetime | |
| updated_at | datetime | |
| deleted_at | datetime | 软删除 |

#### ledger_members

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| ledger_id | uint | FK → ledgers |
| member_id | uint | FK → members |

多对多关联，至少选一个成员。

#### todos

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| title | string(100) | |
| description | string(500) | 可选 |
| priority | string | normal / important / urgent |
| status | string | pending / completed |
| creator_id | uint | FK → members |
| assignee_id | uint | FK → members，可空 |
| due_date | date | 仅日期，可选；已过截止日本完成时前端红色高亮 |
| completed_at | datetime | 可选 |
| created_at | datetime | |
| updated_at | datetime | |
| deleted_at | datetime | 软删除 |

#### todo_logs

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| todo_id | uint | FK → todos |
| field | string | 变更字段 |
| old_value | string | |
| new_value | string | |
| operator_id | uint | FK → members |
| created_at | datetime | |

#### wishes

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| title | string(100) | |
| description | string(500) | 可选 |
| category | string | item / travel / experience / other |
| amount | int64 | 预估金额(分)，可选 |
| priority | string | normal / important / urgent |
| type | string | personal / family |
| status | string | pending / agreed / achieved / abandoned |
| creator_id | uint | FK → members |
| promoted_at | datetime | 提升时间，可选 |
| created_at | datetime | |
| updated_at | datetime | |
| deleted_at | datetime | 软删除 |

#### wish_votes

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| wish_id | uint | FK → wishes |
| member_id | uint | FK → members |
| created_at | datetime | |

唯一约束：(wish_id, member_id)

#### posts

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| content | string(1000) | |
| creator_id | uint | FK → members |
| created_at | datetime | |
| deleted_at | datetime | 软删除 |

#### topics

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| title | string(100) | |
| content | string(2000) | |
| tag_id | uint | FK → tags，可选 |
| is_pinned | bool | 是否置顶(公告) |
| creator_id | uint | FK → members |
| created_at | datetime | |
| deleted_at | datetime | 软删除 |

#### votes

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| title | string(100) | |
| type | string | single / multiple |
| deadline_at | date | 可选，不设则永久 |
| creator_id | uint | FK → members |
| created_at | datetime | |
| deleted_at | datetime | 软删除 |

#### vote_options

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| vote_id | uint | FK → votes |
| content | string(50) | |
| sort_order | int | |

#### vote_records

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| vote_id | uint | FK → votes |
| option_id | uint | FK → vote_options |
| member_id | uint | FK → members |
| created_at | datetime | |

唯一约束：(vote_id, member_id, option_id)

#### comments

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| target_type | string | post / topic / wish |
| target_id | uint | |
| parent_id | uint | 回复的一级评论ID，空=一级 |
| content | string(500) | |
| creator_id | uint | FK → members |
| created_at | datetime | |
| deleted_at | datetime | 软删除 |

#### likes

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| target_type | string | post / topic / comment |
| target_id | uint | |
| member_id | uint | FK → members |
| created_at | datetime | |

唯一约束：(target_type, target_id, member_id)

#### tags

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint PK | |
| name | string(20) | 唯一 |
| preset | bool | 是否预置 |
| created_at | datetime | |
| deleted_at | datetime | 软删除 |

**预置标签（10个）**：家务、育儿、出行、饮食、健康、教育、财务、购物、装修、宠物

### 3.3 数据库迁移与备份

**迁移工具**：goose（SQL 迁移文件，目录 `backend/migrations/`）

**迁移策略**：
- 应用启动时自动执行增量迁移（仅 Up 方向，不支持降级）
- 迁移文件命名规范：`{版本号}_{描述}.sql`，如 `2026051601_init.sql`
- 每次启动时 goose 自动检测未执行的迁移并按版本号顺序执行

**备份策略**：
- 迁移执行前自动备份当前 SQLite 文件到 `backups/backup-{timestamp}.db`
- 迁移失败时启动中止，日志输出错误信息，提示用户手动用备份文件恢复
- 迁移成功则保留最近 7 份备份，超出后自动清理最旧备份

**预置数据初始化**：
- 预置数据由 goose migration 在数据库迁移阶段创建，不在 `POST /api/init/setup` 中创建：
  - 预置分类 20 个（支出 15 个 + 收入 5 个）
  - 预置标签 10 个
- `POST /api/init/setup` 仅负责创建管理员成员 + 管理员账号（用户名、密码、姓名），头像默认 👨，角色固定 admin
- 后续迁移中如需补充新的预置数据，在迁移文件中使用 `INSERT OR IGNORE` 确保幂等
- 用户已手动删除的预置条目不会因迁移恢复

### 3.4 数据库索引

| 表 | 索引 | 理由 |
|------|------|------|
| ledgers | `(deleted_at, occurred_at)` | 按月查询+排序高频 |
| ledgers | `(creator_id, deleted_at)` | 按记录者筛选 |
| ledgers | `(category_id, deleted_at)` | 按分类筛选+删除校验 |
| ledger_members | `(member_id)` | 按成员筛选 |
| todos | `(deleted_at, status, priority, due_date)` | 列表排序 |
| todos | `(assignee_id, deleted_at)` | 按指派人筛选 |
| posts | `(deleted_at, created_at)` | 信息流排序 |
| topics | `(deleted_at, created_at)` | 信息流排序 |
| topics | `(is_pinned, deleted_at, created_at)` | 公告置顶查询 |
| comments | `(target_type, target_id, deleted_at)` | 评论列表 |
| likes | `(target_type, target_id)` | 点赞数统计 |
| wish_votes | `(wish_id)` | 愿望投票数 |
| vote_records | `(vote_id, member_id)` | 投票去重+结果 |
| members | `(username, deleted_at)` | 登录查询+用户名唯一校验 |

---

## 4. API 设计

### 4.1 统一约定

- 基础路径：`/api`
- 响应格式：
  ```json
  { "code": 0, "message": "ok", "data": {} }
  ```
- 列表分页：`?page=1&page_size=20`（上限 100），响应带 `total`
- 认证方式：`Authorization: Bearer <token>`
- 业务错误码：
  - 40001: 参数错误
  - 40004: 资源不存在
  - 40101: Token 过期/无效
  - 40301: 权限不足
  - 40302: 只能操作自己的内容
  - 42901: 登录频繁（锁定）

### 4.2 接口列表

#### 认证与初始化

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/init/check | 检查初始化状态 | 无 |
| POST | /api/init/setup | 系统初始化（创建管理员） | 无 |
| POST | /api/auth/login | 登录 | 无 |
| POST | /api/auth/logout | 退出 | 登录 |
| PUT | /api/auth/password | 修改密码 | 登录 |

**`POST /api/init/setup` 请求体**：

```json
{
  "username": "admin",
  "password": "secure123",
  "name": "张三"
}
```

头像默认 `👨`，角色固定 `admin`，后端自动设置，无需前端传入。

#### 成员管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/members | 成员列表 | 登录 |
| POST | /api/members | 添加成员 | admin |
| PUT | /api/members/:id | 编辑成员 | admin |
| DELETE | /api/members/:id | 删除成员 | admin |
| PUT | /api/members/:id/disable | 禁用成员 | admin |
| PUT | /api/members/:id/enable | 启用成员 | admin |
| PUT | /api/members/:id/reset-pwd | 重置密码 | admin |
| PUT | /api/profile | 修改个人信息 | 登录 |

#### 分类管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/categories | 分类列表 | 登录 |
| POST | /api/categories | 添加分类 | admin |
| PUT | /api/categories/:id | 编辑分类 | admin |
| DELETE | /api/categories/:id | 删除分类 | admin |

#### 记账本

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/ledgers | 记账列表(筛选+分页) | 登录 |
| POST | /api/ledgers | 记一笔 | 登录 |
| PUT | /api/ledgers/:id | 编辑 | 本人/admin |
| DELETE | /api/ledgers/:id | 删除 | 本人/admin |

查询参数：`?month=2026-05&member_id=&category_id=&creator_id=&page=1&page_size=20`

#### 待办

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/todos | 待办列表 | 登录 |
| POST | /api/todos | 添加待办 | 登录 |
| PUT | /api/todos/:id | 编辑 | 创建者/被指派人/admin |
| DELETE | /api/todos/:id | 删除 | 创建者/admin |
| PUT | /api/todos/:id/toggle | 完成/恢复 | 创建者/被指派人/admin |
| PUT | /api/todos/:id/claim | 认领 | 登录(仅未指派) |

查询参数：`?status=pending|completed&assignee_id=&page=1&page_size=20`

#### 愿望清单

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/wishes | 愿望列表 | 登录 |
| POST | /api/wishes | 添加愿望 | 登录 |
| PUT | /api/wishes/:id | 编辑 | 创建者/admin |
| DELETE | /api/wishes/:id | 删除 | 创建者/admin |
| POST | /api/wishes/:id/promote | 提升为家庭愿望 | 创建者 |
| PUT | /api/wishes/:id/status | 变更状态 | admin/创建者(仅放弃) |
| POST | /api/wishes/:id/vote | 投票 | 登录 |
| DELETE | /api/wishes/:id/vote | 取消投票 | 登录 |

查询参数：`?type=personal|family&status=pending|agreed|achieved|abandoned&creator_id=&page=1&page_size=20`

#### 论坛

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/feed | 信息流 | 登录 |

查询参数：`?page=1&page_size=20`

响应结构：
```json
{
  "pinned": [...],
  "items": [...],
  "total": 150
}
```
- `pinned`：所有有效公告（置顶且未删除），按创建时间倒序；前端若为空数组则隐藏公告区
- `items`：动态+话题混合列表，按创建时间倒序，分页
| POST | /api/posts | 发动态 | 登录 |
| DELETE | /api/posts/:id | 删动态 | 创建者/admin |
| POST | /api/topics | 发话题 | 登录 |
| GET | /api/topics/:id | 话题详情(含评论) | 登录 |
| DELETE | /api/topics/:id | 删话题 | 创建者/admin |
| PUT | /api/topics/:id/pin | 置顶/取消(公告) | admin |
| POST | /api/votes | 发起投票 | 登录 |
| DELETE | /api/votes/:id | 删投票(截止前) | 创建者/admin |
| POST | /api/votes/:id/vote | 参与投票 | 登录 |
| GET | /api/votes/:id/result | 投票结果 | 登录 |
| POST | /api/comments | 发表评论 | 登录 |
| DELETE | /api/comments/:id | 删评论 | 创建者/admin |
| POST | /api/likes | 点赞 | 登录 |
| DELETE | /api/likes | 取消赞 | 登录 |

**`POST /api/comments` 请求体**：

```json
{
  "target_type": "post",
  "target_id": 1,
  "parent_id": null,
  "content": "说得好"
}
```

- `target_type` 枚举：`post` / `topic` / `wish`
- `parent_id` 为 null → 一级评论；指向一级评论 ID → 二级评论
- 后端校验：`parent_id` 指向的评论的 `parent_id` 必须为 null（禁止三级嵌套）
- `creator_id` 从 Token 中取

**`POST /api/likes` 请求体**：

```json
{
  "target_type": "post",
  "target_id": 1
}
```

**`DELETE /api/likes` 请求体**：同上，通过 `(target_type, target_id, member_id)` 唯一约束定位。
| GET | /api/tags | 标签列表 | 登录 |
| POST | /api/tags | 添加标签 | admin |
| PUT | /api/tags/:id | 编辑标签 | admin |
| DELETE | /api/tags/:id | 删标签(无关联) | admin |

#### 仪表盘

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/dashboard/summary | 月度收支概览 | 登录 |
| GET | /api/dashboard/expense-chart | 支出分类饼图 | 登录 |
| GET | /api/dashboard/upcoming-todos | 近期待办(5条) | 登录 |
| GET | /api/dashboard/wish-trends | 愿望动态(5条) | 登录 |
| GET | /api/dashboard/forum-hot | 论坛热点 | 登录 |

查询参数：
- `summary?month=2026-05`（可选，默认当前月份）
- `expense-chart?month=2026-05`（可选，默认当前月份）
- `upcoming-todos`、`wish-trends`、`forum-hot` 无需查询参数

#### 个人信息

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | /api/profile | 查看个人信息 | 登录 |
| PUT | /api/profile | 修改姓名/头像 | 登录 |

---

## 5. 前端路由设计

```
/login                → 登录页
/init                 → 初始化页（首次访问）
/                     → 仪表盘（首页）
/ledger               → 记账本
/todo                 → 待办管理
/wish                 → 愿望清单
/forum                → 家庭论坛
/members              → 成员管理（admin）
/categories           → 分类管理（admin）
/profile              → 个人中心
```

### 布局

**桌面端**：侧边栏导航（Logo + 模块菜单） + 顶部栏（当前用户头像+下拉菜单） + 主内容区

顶部栏下拉菜单（点击头像展开）：
- 个人信息 → 跳转 `/profile`
- 修改密码 → 弹出修改密码表单
- 退出登录 → 二次确认后退出

分类管理（`/categories`）和成员管理（`/members`）通过侧边栏直接访问，不出现在下拉菜单中。

**移动端**：底部 TabBar（仪表盘、记账、待办、论坛、我的）+ 顶栏 + 主内容区

"我的" Tab 进入个人中心页面（`/profile`），包含：
- 顶部显示头像、姓名、角色标签、用户名
- 列表项：账单分类（→ `/categories`）、家庭成员（→ `/members`）
- 操作项：修改个人信息、修改密码、退出登录

**登录/初始化页**：独立布局，不含侧边栏或 TabBar。

### 5.1 空状态组件

各模块无数据时遵循统一规范和复用组件，定义如下通用 `<EmptyState>` 组件：

| 参数 | 类型 | 说明 |
|------|------|------|
| message | string | 提示文案 |
| actionText | string | 引导按钮文字，为空则不显示按钮 |
| onAction | function | 按钮点击回调 |
| type | string | no-data / no-result，默认 no-data |

**各页面调用规范**（出自 PRD 第 7.7 节）：

| 模块 | 场景 | message | actionText | 类型 |
|------|------|---------|------------|------|
| 记账本 | 当月无记录 | "本月还没有记录，记一笔吧" | "记一笔" | no-data |
| 记账本 | 筛选无结果 | "没有找到符合条件的记录" | "清除筛选" | no-result |
| 待办 | 无待办 | "暂无待办，添加一个吧" | "添加待办" | no-data |
| 待办 | 无已完成 | "还没有完成的事项" | 无 | no-data |
| 愿望清单 | 无个人愿望 | "还没有愿望，许一个吧" | "添加愿望" | no-data |
| 愿望清单 | 无家庭愿望 | "家庭愿望池是空的，从个人愿望提升或直接创建" | 无 | no-data |
| 论坛 | 无动态/话题 | "还没有内容，发条动态或话题吧" | "发动态""发话题" | no-data |
| 论坛 | 无公告 | 不展示（隐藏公告区域） | 无 | — |
| 仪表盘饼图 | 无支出数据 | "暂无支出数据" | 无 | no-data |

**通用原则**：
1. `no-data` 类型：居中图标 + 文案 + 引导按钮，引导用户创建内容
2. `no-result` 类型：居中图标 + 文案 + "清除筛选"链接，不提供创建入口（筛选条件导致的无结果不应引导创建）
3. 仪表盘数字为 ¥0.00 是正常数据状态，不展示空状态组件
4. 公告区无内容直接在 DOM 中隐藏，不占空间

---

## 6. 认证与授权流程

1. **用户登录** → 服务端校验用户名密码 → 校验成员 status（若为 disabled 则拒绝登录 → 签发 JWT Token（7天有效期）
2. **请求拦截** → 除登录/初始化接口外，Gin 中间件解析 Token → **校验成员 status，若为 disabled 则返回 40301** → 注入当前用户信息到 context
3. **角色校验** → admin 接口二次校验角色
4. **资源归属** → 编辑/删除操作校验当前用户是否为创建者或 admin
5. **修改密码** → 旧密码验证 → 新密码存储 → 清除 Token → 跳转登录
6. **登录锁定** → 5次失败后锁定用户名 15分钟（内存记录，服务重启重置）
7. **默认密码** → `home123`，使用默认密码时提示修改

---

## 7. 并发与一致性

| 策略 | 说明 |
|------|------|
| 写入策略 | Last Write Wins（BR-CC01） |
| SQLite 模式 | WAL 模式保证写事务串行化（BR-CC02） |
| 唯一约束 | 数据库层保证（vote、like 等） |
| 写入失败 | 返回通用错误"操作失败，请重试"（BR-CC03） |

---

## 8. 删除策略

所有删除统一软删除（BR-DL01），设置 `deleted_at` 字段：

| 场景 | 行为 |
|------|------|
| 删除分类 | 有未删除记账记录 → 拒绝；否则软删除 |
| 删除标签 | 有关联话题 → 拒绝；否则软删除 |
| 删除成员 | 有活动记录 → 仅可禁用；否则软删除+停用账号 |
| 删除其他数据 | 直接软删除，关联数据同步软删除 |

---

## 9. 错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 40001 | 参数错误 |
| 40002 | 用户名已存在 |
| 40003 | 资源冲突（如重复投票） |
| 40004 | 资源不存在 |
| 40005 | 操作被拒绝（如有引用依赖） |
| 40101 | Token 过期/无效 |
| 40301 | 权限不足 |
| 40302 | 只能操作自己的内容 |
| 42901 | 登录频次限制 |

---

## 10. 实施顺序

按 RICE 优先级排序：

1. **认证与权限** — 初始化、登录、JWT 中间件
2. **成员管理** — CRUD、禁用/启用、密码重置
3. **分类管理** — CRUD、预置数据初始化
4. **记账本** — CRUD、筛选、日期分组
5. **待办管理** — CRUD、指派、认领、完成切换
6. **仪表盘** — 聚合查询、饼图
7. **愿望清单** — CRUD、提升、投票、状态流转
8. **家庭论坛** — 动态、话题、公告、投票、评论、点赞
9. **个人中心** — 信息查看/编辑、密码修改
10. **移动端适配** — 响应式布局打磨
