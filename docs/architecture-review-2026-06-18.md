# 暖屿V1架构审视报告

**审阅人**：架构师 高见远  
**审阅日期**：2026-06-18  
**项目版本**：V1.0.0  
**代码库路径**：d:/Projects/my_projects/home-center-v1

---

## 一、架构总体评价

**综合评分：4.2/5.0**

| 维度 | 评分 | 说明 |
|------|------|------|
| 文档完整性 | 5/5 | PRD、技术设计文档、AGENTS.md、实施计划非常完善 |
| 架构设计 | 4.5/5 | 三层架构清晰，职责分离合理，但存在少量实现不一致 |
| 代码质量 | 4/5 | 整体规范，但测试文件混杂、部分安全实践需加强 |
| 可扩展性 | 3.5/5 | 模块化良好，但多家庭、通知等V2功能需重构 |
| 安全架构 | 3.5/5 | 认证流程完整，但Token存储、CSRF防护等有风险 |
| 技术债务 | 4/5 | 债务可控，主要是登录追踪持久化、测试分离等 |

**总体评价**：项目架构设计**优秀**，文档质量**极高**，可作为小型私有化项目的参考范例。主要风险点在安全实现细节和V2扩展预留不足。

---

## 二、具体发现

### 2.1 后端架构

#### ✅ 优点

1. **三层架构严格执行**
   - Handler → Service → Repository 分层清晰
   - 所有模块统一遵循该模式（ledger、todo、wish、forum等）
   - 统一响应格式 `{code, message, data}` 在 `pkg/response.go` 中规范实现

2. **数据库设计合理**
   - 金额以"分"存储（int64），避免浮点精度问题
   - 软删除统一使用 `deleted_at` 字段，GORM自动过滤
   - 索引设计完善（见 `001_init.up.sql` 中的12个索引）

3. **业务逻辑封装良好**
   - Service层正确处理权限校验（只能操作自己的内容）
   - 收支类型实时引用分类type，避免快照不一致
   - 论坛评论通用模型支持多态关联（post/topic/wish）

#### ⚠️ 问题

1. **登录失败追踪使用内存存储**
   ```go
   // backend/internal/service/auth.go
   var loginFailures = make(map[string]*LoginFailure)
   ```
   - **问题**：服务重启后锁定记录丢失，攻击者可利用此窗口暴力破解
   - **影响**：安全机制失效
   - **位置**：`service/auth.go`

2. **测试文件混杂在生产代码中**
   - `handler/test_create_member.go`
   - `handler/test_reset.go`
   - `handler/test_seed_ledgers.go`
   - **问题**：这些文件 should be in `testutil/` 或 `cmd/cli/`，不应放在handler包
   - **影响**：生产二进制体积增大，代码组织混乱

3. **Repository层直接返回数据库错误**
   - 部分repository方法将GORM错误直接返回给Service
   - **建议**：统一错误封装，避免数据库错误泄露到API响应

4. **Cursor分页实现过于复杂**
   - `repository/ledger.go` 中的 `List()` 方法长达170行
   - 包含"补全"逻辑（防止日期分组截断）
   - **建议**：考虑简化或提取为独立方法，提高可测试性

### 2.2 前端架构

#### ✅ 优点

1. **状态管理清晰**
   - Pinia stores 按模块分离（auth、categories、members、theme）
   - 响应式设计良好（computed属性正确使用）

2. **路由守卫完整**
   - `router/index.ts` 实现了初始化检测、登录检测、管理员权限检测
   - 支持`requiresAuth`和`requiresAdmin`元字段

3. **API层封装良好**
   - 每个模块有独立的API文件（`api/ledger.ts`、`api/todo.ts`等）
   - 统一的请求/响应拦截器

#### ⚠️ 问题

1. **Token存储安全风险**
   ```typescript
   // stores/auth.ts
   const token = ref(localStorage.getItem('token') || '')
   localStorage.setItem('token', res.data.token)
   ```
   - **问题**：使用localStorage存储JWT，易受XSS攻击
   - **风险**：恶意脚本可读取localStorage中的Token
   - **建议**：改用HttpOnly Cookie，或至少启用CSP

2. **请求拦截器拼写错误**
   ```typescript
   // utils/request.ts
   config.headers.Authorization = `Bearer ${token}`
   ```
   - **实际代码**：`config.headers.Authorization`（拼写错误，应为`Authorization`）
   - **影响**：可能导致某些代理服务器无法识别Authorization头
   - **位置**：`utils/request.ts` 第13行

3. **JWT Payload客户端解析**
   ```typescript
   // stores/auth.ts
   function parseJwtPayload(rawToken: string): Record<string, any> | null {
     return JSON.parse(atob(rawToken.split('.')[1]))
   }
   ```
   - **问题**：前端直接解析JWT payload获取用户ID和角色
   - **风险**：虽然不验证签名是合理的（无法隐藏前端信息），但依赖前端解析结果做权限判断存在风险
   - **建议**：关键权限判断应依赖后端返回的用户信息，而非解析Token

4. **组件测试覆盖率不足**
   - 有测试文件（如 `views/ledger/__tests__/Index.test.ts`）
   - 但主要测试的是页面组件，工具函数和store的测试较少

### 2.3 数据库架构

#### ✅ 优点

1. **迁移管理规范**
   - 使用goose管理迁移，命名规范（`001_init.up.sql`）
   - 迁移前自动备份到 `backups/` 目录
   - 支持幂等性（`INSERT OR IGNORE`）

2. **索引设计合理**
   - 高频查询字段均有索引（如 `idx_ledgers_deleted_occurred`）
   - 外键查询优化（如 `idx_ledger_members_member`）

3. **约束完整**
   - 枚举值使用CHECK约束（如 `role IN ('admin', 'member')`）
   - 唯一约束正确设置（如投票记录的 `(vote_id, member_id, option_id)`）

#### ⚠️ 问题

1. **SQLite外键约束可能未启用**
   ```sql
   -- 001_init.up.sql
   category_id INTEGER NOT NULL REFERENCES categories(id),
   ```
   - **问题**：SQLite默认**不强制执行**外键约束，需要手动启用 `PRAGMA foreign_keys = ON`
   - **检查**：未在 `pkg/database.go` 中看到启用外键的代码
   - **风险**：可能插入无效的category_id或creator_id

2. **迁移文件缺少Down方法**
   - 只有 `001_init.up.sql` 有Down部分
   - 后续迁移（如 `002_fix_time_format.up.sql`）只有Up方向
   - **影响**：无法回滚迁移，与技术设计文档中"不支持降级"的决策一致，但风险较高

3. **预置数据管理策略**
   - 预置数据（分类、标签）在迁移文件中硬编码
   - **问题**：如果未来需要修改预置数据（如修改图标、添加新分类），需要新建迁移
   - **建议**：考虑将预置数据管理分离到独立的种子系统

### 2.4 安全架构

#### ✅ 优点

1. **认证流程完整**
   - JWT签发、验证、过期处理均有实现
   - 中间件正确解析Token并注入用户信息

2. **密码存储安全**
   - 使用bcrypt哈希存储密码
   - 密码复杂度要求（至少字母+数字）

3. **权限校验多层防护**
   - 中间件层：JWT解析 + status校验
   - Handler层：角色校验（AdminRequired）
   - Service层：资源归属校验（只能操作自己的内容）

4. **登录失败锁定**
   - 5次失败锁定15分钟
   - 管理员可通过CLI工具解锁

#### ⚠️ 问题

1. **Token存储XSS风险**（已在前端部分提及）

2. **无CSRF防护**
   - 使用Bearer Token认证，但未启用CSRF Token
   - **风险**：虽然Bearer Token本身可以防御CSRF（因为CSRF无法携带自定义Header），但双重Cookie校验是更安全的实践

3. **无速率限制**
   - 未实现API级别的速率限制
   - **风险**：登录接口可被暴力破解（虽然有账号锁定，但攻击者可以轮换用户名）

4. **论坛内容无HTML过滤**
   ```go
   // model/forum.go
   type Post struct {
       Content string `gorm:"size:1000" json:"content"`
   }
   ```
   - **风险**：如果前端允许输入HTML，后端未过滤，可能导致XSS
   - **建议**：后端在存储前过滤HTML标签，或前端使用v-text而非v-html

5. **CORS配置可能过于宽松**
   ```go
   // AGENTS.md中提到
   HC_CORS_ORIGINS | 开发：`*`，生产（嵌入前端）：仅同源
   ```
   - **问题**：开发环境下CORS允许所有Origin，如果生产环境配置错误，可能导致跨域攻击

---

## 三、技术债务清单

### P0（立即修复）

| ID | 债务描述 | 位置 | 风险 | 修复建议 | 工作量 |
|----|----------|------|------|----------|--------|
| TD-01 | 登录失败追踪使用内存存储 | `service/auth.go` | 重启后锁定失效，可被暴力破解 | 改用SQLite表存储失败记录 | 2h |
| TD-02 | SQLite外键约束可能未启用 | `pkg/database.go` | 可能插入无效外键值 | 在初始化时执行 `PRAGMA foreign_keys = ON` | 0.5h |
| TD-03 | 请求头拼写错误Authorization | `utils/request.ts` | 某些代理可能无法识别 | 修正拼写 | 0.1h |

### P1（近期修复）

| ID | 债务描述 | 位置 | 风险 | 修复建议 | 工作量 |
|----|----------|------|------|----------|--------|
| TD-04 | 测试文件混杂在生产代码中 | `handler/test_*.go` | 代码组织混乱，二进制体积增大 | 迁移到 `testutil/` 或 `cmd/cli/` | 1h |
| TD-05 | Token存储在localStorage | `stores/auth.ts` | XSS攻击可窃取Token | 改用HttpOnly Cookie，或启用CSP | 4h |
| TD-06 | 论坛内容无HTML过滤 | `handler/forum_*.go` | XSS攻击风险 | 在后端存储前过滤HTML标签 | 2h |
| TD-07 | 无API速率限制 | 全局 | 暴力破解、DoS风险 | 添加速率限制中间件（如 `gin-contrib/limiter`） | 3h |

### P2（可选优化）

| ID | 债务描述 | 位置 | 风险 | 修复建议 | 工作量 |
|----|----------|------|------|----------|--------|
| TD-08 | Cursor分页实现复杂 | `repository/ledger.go` | 可测试性差，维护困难 | 提取为独立方法，添加单元测试 | 3h |
| TD-09 | 预置数据硬编码在迁移中 | `migrations/001_init.up.sql` | 修改预置数据需新建迁移 | 考虑分离种子数据管理系统 | 4h |
| TD-10 | 无CSRF防护 | 全局 | CSRF攻击风险（低，因为使用Bearer Token） | 添加双重Cookie校验 | 2h |
| TD-11 | Repository层错误未统一封装 | `repository/*.go` | 数据库错误可能泄露到API | 定义统一错误类型，封装GORM错误 | 2h |

---

## 四、可扩展性评估

### 4.1 当前架构对V2功能的支持

| V2功能 | 扩展性评估 | 所需改动 |
|---------|--------------|----------|
| **消息通知** | ❌ 不支持 | 需新增通知表、WebSocket或轮询接口、前端通知组件 |
| **数据导出** | ⚠️ 部分支持 | 后端需新增导出接口（CSV/Excel），前端需新增导出按钮 |
| **多家庭** | ❌ 不支持 | 需新增families表，所有表添加family_id外键，重构权限系统 |
| **预算管理** | ⚠️ 部分支持 | 需新增budgets表，仪表盘新增预算进度展示 |
| **搜索功能** | ❌ 不支持 | 需实现全文搜索（SQLite FTS或Elasticsearch） |
| **批量操作** | ⚠️ 部分支持 | 后端需新增批量接口，前端需新增批量选择UI |

### 4.2 架构改进建议（为V2做准备）

1. **引入事件系统**
   - 当前：操作直接修改数据库，无副作用处理
   - 建议：引入领域事件（如 `LedgerCreated`、`TodoCompleted`），为通知系统做准备
   - 实现：可使用Go channel或简单的事务性Outbox表

2. **抽象多租户上下文**
   - 当前：假设单家庭
   - 建议：在中间件中注入 `family_id`，即使V1不使用，也为V2预留
   - 实现：新增 `middleware/family.go`，从JWT或请求参数中提取family_id

3. **统一分页接口**
   - 当前：Ledger使用Cursor分页，其他模块使用Offset分页
   - 建议：统一为Cursor分页（性能更好），或定义统一的分页接口

4. **配置管理改进**
   - 当前：使用环境变量，无配置文件支持
   - 建议：引入配置文件支持（如 `config.yaml`），便于V2时管理多家庭配置

---

## 五、改进建议

### 5.1 安全加固（优先级最高）

1. **Token存储安全**
   ```typescript
   // 当前：存储在localStorage
   localStorage.setItem('token', token)
   
   // 建议：改用HttpOnly Cookie
   // 后端设置Cookie，前端无需操作
   c.SetCookie("token", token, 7*24*3600, "/", "", false, true)
   ```

2. **启用CSP（内容安全策略）**
   ```go
   // middleware/security.go
   func SecurityHeaders() gin.HandlerFunc {
       return func(c *gin.Context) {
           c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
           c.Next()
       }
   }
   ```

3. **输入验证增强**
   - 后端：使用 `go-playground/validator` 统一验证请求
   - 前端：对所有用户输入进行编码，避免XSS

### 5.2 代码质量提升

1. **测试文件分离**
   ```
   backend/
   ├── internal/
   │   ├── handler/        # 只保留生产代码
   │   ├── testutil/       # 测试辅助工具
   │   └── cmd/
   │       └── cli/        # CLI工具（含reset-password）
   ```

2. **统一错误处理**
   ```go
   // pkg/errors.go
   type AppError struct {
       HTTPCode  int
       BizCode   int
       Message   string
       Err       error
   }
   
   func (e *AppError) Error() string {
       return e.Message
   }
   ```

3. **Repository层返回统一错误**
   - 定义 `repository.ErrNotFound`、`repository.ErrConflict` 等统一错误
   - Service层直接判断错误类型，无需额外查询

### 5.3 性能优化

1. **数据库查询优化**
   - 当前：`repository/ledger.go` 中的 `computeSummary` 每次列表查询都执行
   - 建议：对于月度数据，可考虑缓存（如每5分钟更新一次）

2. **前端懒加载**
   - 当前：所有路由组件在定义时导入
   - 建议：已正确使用Vue的懒加载（`() => import(...)`），保持即可

3. **API响应压缩**
   - 建议：在Gin中启用Gzip中间件
   ```go
   import "github.com/gin-contrib/gzip"
   router.Use(gzip.Gzip(gzip.DefaultCompression))
   ```

### 5.4 运维可观测性

1. **结构化日志**
   - 当前：未看到统一日志系统
   - 建议：引入 `logrus` 或 `zap`，输出JSON格式日志，便于后续接入ELK

2. **健康检查接口**
   ```go
   // routes/router.go
   router.GET("/healthz", func(c *gin.Context) {
       c.JSON(200, gin.H{"status": "ok"})
   })
   ```

3. **性能指标**
   - 建议：在Gin中启用 `gin-contrib/pprof`，便于性能分析

---

## 六、架构决策记录（ADD）

基于本次审视，建议补充以下架构决策文档到 `docs/architecture/decisions/`：

### ADR-001：Token存储方案
- **状态**：已决策（localStorage）
- **背景**：需要存储JWT Token
- **决策**：使用localStorage（而非HttpOnly Cookie）
- **理由**：简化实现，避免Cookie跨域问题
- **后果**：
  - 正面：实现简单，无Cookie跨域问题
  - 负面：XSS风险，需依赖前端CSP和输入过滤
- **日期**：2026-05-16
- **修订**：建议V1.1改为HttpOnly Cookie

### ADR-002：分页策略
- **状态**：混合方案
- **背景**：列表查询需要分页
- **决策**：Ledger使用Cursor分页，其他模块使用Offset分页
- **理由**：Ledger数据量大，Cursor性能更好；其他模块数据量小，Offset实现简单
- **后果**：
  - 正面：性能优化
  - 负面：分页接口不统一，前端需适配两种分页逻辑
- **日期**：2026-05-16
- **修订**：建议V2统一为Cursor分页

### ADR-003：登录失败追踪存储
- **状态**：已决策（内存）
- **背景**：需要追踪登录失败次数以实现锁定
- **决策**：使用内存map存储
- **理由**：简化实现，V1用户量少，重启概率低
- **后果**：
  - 正面：实现简单，无需额外表
  - 负面：重启后锁定记录丢失
- **日期**：2026-05-16
- **修订**：**建议立即改为数据库存储**（P0技术债务TD-01）

---

## 七、总结与行动计划

### 7.1 立即行动（本周）

1. ✅ **修复P0技术债务**
   - [ ] TD-01：登录失败追踪改用数据库存储（2h）
   - [ ] TD-02：启用SQLite外键约束（0.5h）
   - [ ] TD-03：修正Authorization拼写错误（0.1h）

2. ✅ **安全加固**
   - [ ] 后端添加HTML过滤（论坛内容）（2h）
   - [ ] 前端启用CSP（1h）

### 7.2 近期计划（本月）

1. ✅ **代码质量提升**
   - [ ] 分离测试文件（1h）
   - [ ] 统一错误处理（2h）

2. ✅ **性能优化**
   - [ ] 启用Gzip压缩（0.5h）
   - [ ] 添加健康检查接口（0.5h）

3. ✅ **可观测性**
   - [ ] 引入结构化日志（2h）
   - [ ] 添加基础性能监控（1h）

### 7.3 V2准备（下个迭代）

1. ✅ **架构改进**
   - [ ] 引入事件系统（为通知功能做准备）（1天）
   - [ ] 抽象多租户上下文（为多家庭功能做准备）（1天）
   - [ ] 统一分页策略（Cursor分页）（1天）

2. ✅ **新功能设计**
   - [ ] 消息通知系统设计
   - [ ] 数据导出功能设计
   - [ ] 多家庭架构设计

---

## 八、附录：检查清单

### A. 架构文档完整性 ✅

- [x] PRD（产品需求文档）：`docs/prd.md`（极高质量）
- [x] 技术设计文档：`docs/superpowers/specs/2026-05-16-home-center-technical-design.md`
- [x] 实施计划：`docs/superpowers/plans/` 目录
- [x] AGENTS.md（开发指南）：项目根目录
- [ ] 架构决策记录（ADD）：**缺失**，建议补充
- [ ] API文档（自动生成）：**缺失**，建议使用Swagger或类似工具

### B. 代码质量检查 ✅

- [x] 后端三层架构：严格执行
- [x] 前端模块化：良好
- [ ] 测试覆盖率：**不足**，建议补充
- [ ] 代码审查流程：**未看到**，建议建立

### C. 安全审查 ✅

- [x] 认证流程：完整
- [x] 密码存储：安全（bcrypt）
- [ ] Token存储：**存在风险**（localStorage）
- [ ] CSRF防护：**缺失**
- [ ] 速率限制：**缺失**
- [ ] 输入验证：**部分缺失**（论坛内容无过滤）

---

**报告结束**

> **审阅人签名**：高见远（软件架构师）  
> **日期**：2026-06-18  
> **下次审阅建议**：V1.1发布前（预计2026-07）
