# 🏝️ 暖屿 · WarmIsle

> 面向家庭的私有化协作平台 — 记账、待办、愿望、论坛，一个地方搞定。

**暖屿（WarmIsle）** 是一座属于你们家的温暖岛屿。在这里，每一笔收支有人记，每一件待办有人理，每一个愿望被看见，每一次讨论有回音。单二进制一键部署，数据完全私有。

---

## ✨ 功能

| 模块 | 说明 |
|------|------|
| 📒 **记账本** | 全员协作记账，按日分组、月度小计、分类筛选、支出饼图 |
| ✅ **待办管理** | 创建指派、主动认领、优先级、截止日期、过期高亮 |
| ⭐ **愿望清单** | 个人愿望 → 提升为家庭愿望 → 投票讨论 → 落实 |
| 💬 **家庭论坛** | 动态墙、话题讨论、家庭公告、投票决策、评论点赞 |
| 📊 **仪表盘** | 收支概览、支出分类占比、近期待办、愿望动态、论坛热点 |

## 🚀 快速开始

```bash
# 下载对应平台的二进制，单文件运行
./warmisle

# 浏览器打开 http://localhost:8080
# 首次访问进入初始化页面，创建管理员账号
```

### 从源码构建

```bash
git clone https://github.com/emf1002/warmisle.git
cd warmisle
make build    # 构建单二进制
./warmisle    # 启动
```

### 开发模式

```bash
make dev      # 后端 air 热重载 :8080 + 前端 Vite dev :3000
```

## 🛠️ 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.22+, Gin, GORM, SQLite (WAL), goose |
| 前端 | Vue 3 (Composition API), Ant Design Vue, Vite, Pinia |
| 部署 | 单二进制（前端 dist embed 到 Go binary） |

## ⚙️ 配置

通过环境变量配置，开箱即用：

| 变量 | 说明 | 默认值 |
|---|---|---|
| `HC_JWT_SECRET` | JWT 签名密钥 | 自动生成 32 位 hex |
| `HC_DB_PATH` | SQLite 文件路径 | `./data/warmisle.db` |
| `HC_PORT` | 监听端口 | `8080` |
| `HC_CORS_ORIGINS` | CORS 白名单 | 开发 `*`，生产仅同源 |

## 📁 项目结构

```
warmisle/
├── backend/
│   ├── cmd/server/        # Web 服务入口
│   ├── cmd/cli/           # CLI 工具（reset-password）
│   ├── internal/
│   │   ├── handler/       # HTTP 处理器
│   │   ├── service/       # 业务逻辑
│   │   ├── repository/    # 数据访问
│   │   ├── model/         # 数据模型
│   │   ├── middleware/    # JWT 认证、CORS
│   │   ├── routes/        # 路由注册
│   │   └── pkg/           # 工具包
│   └── migrations/        # 数据库迁移
├── frontend/
│   └── src/
│       ├── views/         # 页面组件
│       ├── components/    # 通用组件
│       ├── layouts/       # 布局组件
│       ├── stores/        # Pinia 状态管理
│       ├── router/        # 路由配置
│       └── api/           # API 请求封装
└── Makefile
```

## 🏗️ 架构

后端三层架构：**Handler → Service → Repository**

- **Handler**：参数绑定与校验，调用 service，返回统一响应
- **Service**：业务规则校验、权限检查，调用 repository
- **Repository**：GORM 数据库操作，不含业务逻辑

## 📜 许可

MIT
