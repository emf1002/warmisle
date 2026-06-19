# 暖屿 V1 安全审计报告

**审计日期**：2026-06-19  
**审计范围**：完整代码库（后端 Go + 前端 Vue 3）  
**审计标准**：OWASP Top 10:2021 + CWE Top 25  
**审计人**：安全工程师（Security Engineer Agent）

---

## 一、系统概览

| 属性 | 值 |
|------|-----|
| 架构 | 单体应用（Go Gin + Vue 3 SPA，编译为单二进制） |
| 数据分类 | PII（用户名、姓名）+ 财务数据（记账金额）+ 家庭隐私 |
| 信任边界 | 浏览器 → API（JWT Bearer）→ 服务层 → SQLite |
| 认证方式 | 用户名+密码 → JWT Token（HS256, 7天有效期） |
| 授权模型 | RBAC（admin / member）+ 资源级权限检查 |

---

## 二、攻击面分析

```
[浏览器] ──/api/*──▶ [Gin HTTP Server] ──▶ [Service Layer] ──▶ [SQLite WAL]
                           │
                    ┌──────┴──────┐
                    │  JWT中间件   │
                    │  静态文件服务 │
                    └─────────────┘
```

- **外部攻击面**：`/api/auth/login`、`/api/init/setup`、`/api/backup/callback`（无需认证）
- **内部攻击面**：所有 `/api/*` 端点（认证后）、SQLite 数据库文件
- **数据攻击面**：localStorage Token、响应中的 PII、数据库备份文件

---

## 三、漏洞发现（按严重度排序）

### 🔴 CRITICAL（P0 — 必须立即修复）

#### C-1：全站缺少速率限制（CWE-307）
| 项目 | 详情 |
|------|------|
| **影响端点** | 所有 API 端点，尤其是 `/api/auth/login` |
| **风险** | 分布式暴力破解、资源耗尽 DoS |
| **现状** | 仅有应用层账号锁定（5次/15分钟），无法防御 IP 轮换攻击 |
| **CVSS 评分** | 7.5 (AV:N/AC:L/PR:N/UI:N/S:U/C:L/I:L/A:H) |

#### C-2：缺少安全响应头（CWE-693）
| 项目 | 详情 |
|------|------|
| **缺失头** | CSP、HSTS、X-Frame-Options、X-Content-Type-Options、Referrer-Policy、Permissions-Policy |
| **风险** | XSS 利用、点击劫持、MIME 嗅探攻击 |
| **现状** | Gin 框架未设置任何安全头，静态文件有 ETag 和 Cache-Control 但无安全相关 |

---

### 🟠 HIGH（P1 — 本周内修复）

#### H-1：论坛动态（Post）内容缺少 XSS 过滤（CWE-79）
| 项目 | 详情 |
|------|------|
| **影响文件** | `backend/internal/service/forum_post.go:10` |
| **风险** | 存储型 XSS — 恶意脚本嵌入动态内容，所有查看者受影响 |
| **根源** | `CreatePost` 和 `UpdatePost` 未调用 `pkg.SanitizeHTML()`，而 Topic 和 Comment 已正确调用 |
| **CVSS 评分** | 6.1 (AV:N/AC:L/PR:N/UI:R/S:C/C:L/I:L/A:N) |

#### H-2：前端 Token 存储于 localStorage（CWE-312）
| 项目 | 详情 |
|------|------|
| **影响文件** | `frontend/src/stores/auth.ts:15`, `frontend/src/utils/request.ts:11` |
| **风险** | 任何 XSS 漏洞可导致 Token 被盗，实现账户接管 |
| **现状** | JWT 存储在 `localStorage`，通过 `Authorization: Bearer` 头发送 |

#### H-3：缺少 CORS 配置（CWE-942）
| 项目 | 详情 |
|------|------|
| **影响** | 生产环境未显式配置 CORS，依赖浏览器默认同源策略 |
| **风险** | 如果反向代理配置不当，可能导致跨域请求伪造 |
| **现状** | 后端无任何 CORS 中间件 |

---

### 🟡 MEDIUM（P2 — 本月内修复）

#### M-1：登录接口用户枚举风险（CWE-204）
| 项目 | 详情 |
|------|------|
| **影响文件** | `backend/internal/handler/auth.go:36-47` |
| **风险** | 攻击者可通过不同错误响应区分"用户不存在"和"密码错误" |
| **现状** | 虽然已经统一错误消息，但锁定状态单独返回了 HTTP 429 和不同消息 |

#### M-2：LoginFailure 存储存在 TOCTOU 竞态条件（CWE-367）
| 项目 | 详情 |
|------|------|
| **影响文件** | `backend/internal/repository/login_failure.go:22-32` |
| **风险** | 并发登录请求可能绕过锁定检查 |
| **现状** | `Save` 方法先查询后操作，中间存在窗口期 |

#### M-3：敏感信息在 API 响应中暴露（CWE-200）
| 项目 | 详情 |
|------|------|
| **影响文件** | `backend/internal/model/member.go:21` |
| **风险** | `deleted_at` 时间戳暴露给前端 |
| **现状** | `Member.DeletedAt` 字段有 `json:"deleted_at"` 标签 |

#### M-4：缺少请求体大小限制（CWE-770）
| 项目 | 详情 |
|------|------|
| **风险** | 无限制的请求体可导致内存耗尽 |
| **现状** | Gin 默认无请求体大小限制 |

---

### 🟢 LOW（P3 — 后续迭代）

#### L-1：JWT 使用对称加密算法 HS256
- 建议未来版本迁移到 RS256/ES256 以支持多服务架构和密钥轮换

#### L-2：缺少结构化安全审计日志
- 关键操作（登录成功/失败、权限变更、数据删除）无结构化日志

#### L-3：密码策略较弱
- 密码长度要求 6-32 位，无复杂度要求（大小写+数字+特殊字符）

---

## 四、正面发现

以下安全措施值得肯定：

1. **bcrypt 密码哈希**：使用 `golang.org/x/crypto/bcrypt` 和 `DefaultCost`，业界标准
2. **参数化查询**：全项目使用 GORM 参数化查询，无 SQL 注入风险
3. **JWT 中间件**：正确验证并提取 claims，禁用用户实时检查
4. **资源级权限**：所有 CRUD 操作检查 `creator_id == currentMemberID || role == admin`
5. **HTML 过滤**：Topic/Comment 使用 bluemonday UGCPolicy 防 XSS（但 Post 遗漏了）
6. **密码字段 JSON 隐藏**：`Member.Password` 使用 `json:"-"` 标签，API 从不返回
7. **统一错误处理**：`handleServiceError` 模式避免原始错误泄露
8. **数据库加密备份**：AES-256-GCM 加密备份，密钥安全持久化
9. **LoginFailure 持久化**：登录失败记录写入 SQLite，优于设计文档中描述的"内存存储"
10. **软删除统一策略**：全系统使用 GORM 软删除，无数据丢失风险

---

## 五、修复优先级路线图

| 优先级 | 漏洞编号 | 修复内容 | 预计工时 |
|--------|---------|---------|---------|
| **P0（今天）** | C-1, C-2 | 添加速率限制中间件 + 安全响应头 | 2h |
| **P1（本周）** | H-1 | 论坛 Post 内容 XSS 过滤 | 15min |
| **P1（本周）** | H-3 | CORS 显式配置 | 30min |
| **P2（本月）** | M-1, M-2, M-3, M-4 | 用户枚举修复、竞态修复、敏感信息隐藏、请求大小限制 | 4h |
| **P3（下个迭代）** | L-1, L-2, L-3 | JWT 算法升级、审计日志、密码策略 | 待评估 |

---

## 六、安全加固检查表

- [ ] 速率限制中间件：100 req/min/IP（登录端点 5 req/min）
- [ ] 安全响应头：CSP, HSTS, X-Frame-Options, X-Content-Type-Options
- [ ] Post 内容 XSS 过滤
- [ ] CORS 白名单配置
- [ ] 登录错误消息统一化
- [ ] LoginFailure upsert 原子化
- [ ] Member.DeletedAt 从 JSON 响应中隐藏
- [ ] 请求体大小限制（如 1MB）
- [ ] 反向代理层 TLS 终结配置建议
