# E2E 测试覆盖率全量补充 — 实施计划索引

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 按 RICE 优先级逐模块补齐 95 个 e2e 测试用例，覆盖功能场景、权限测试、错误路径和移动端响应式。

**Architecture:** 扩展现有 Playwright POM 架构，新增 member 角色 fixture 支持多用户权限测试，新增后端测试端点快速创建 member 用户，Playwright 配置新增 mobile project 支持响应式验证。

**Tech Stack:** Playwright Test, TypeScript, Go (Gin), Vue 3, Ant Design Vue

**Design Spec:** `docs/superpowers/specs/2026-06-05-e2e-coverage-expansion-design.md`

---

## 任务依赖关系

```
Task 0: 基础设施
  ├── Task 1: 记账本 (14 tests)
  ├── Task 2: 待办管理 (14 tests)
  ├── Task 3: 仪表盘 (7 tests)
  ├── Task 4: 愿望清单 (13 tests)
  ├── Task 5: 家庭论坛 (18 tests)
  ├── Task 6: 成员管理 (11 tests)
  ├── Task 7: 分类管理 (6 tests)
  ├── Task 8: 个人中心 (6 tests)
  └── Task 9: 认证补充 (2 tests)
        └── Task 10: 移动端视觉回归 (4 tests)
```

Task 0 必须最先完成。Task 1-9 相互独立，可并行或按任意顺序执行。Task 10 依赖 Task 0 中的 Playwright mobile project 配置。

## 任务列表

| 任务 | 模块 | 新增用例 | 详细计划 |
|------|------|----------|----------|
| 0 | 基础设施 | — | [e2e-expansion-infrastructure.md](2026-06-05-e2e-expansion-infrastructure.md) |
| 1 | 记账本 | 14 | [e2e-expansion-ledger.md](2026-06-05-e2e-expansion-ledger.md) |
| 2 | 待办管理 | 14 | [e2e-expansion-todo.md](2026-06-05-e2e-expansion-todo.md) |
| 3 | 仪表盘 | 7 | [e2e-expansion-dashboard.md](2026-06-05-e2e-expansion-dashboard.md) |
| 4 | 愿望清单 | 13 | [e2e-expansion-wish.md](2026-06-05-e2e-expansion-wish.md) |
| 5 | 家庭论坛 | 18 | [e2e-expansion-forum.md](2026-06-05-e2e-expansion-forum.md) |
| 6 | 成员管理 | 11 | [e2e-expansion-members.md](2026-06-05-e2e-expansion-members.md) |
| 7 | 分类管理 | 6 | [e2e-expansion-categories.md](2026-06-05-e2e-expansion-categories.md) |
| 8 | 个人中心 | 6 | [e2e-expansion-profile.md](2026-06-05-e2e-expansion-profile.md) |
| 9 | 认证补充 | 2 | [e2e-expansion-auth.md](2026-06-05-e2e-expansion-auth.md) |
| 10 | 移动端视觉回归 | 4 | [e2e-expansion-mobile.md](2026-06-05-e2e-expansion-mobile.md) |
| **合计** | | **95** | |

## 预期结果

- 现有 50 个用例 + 新增 95 个 = **145 个用例**
- 所有模块覆盖 功能场景 + 权限测试 + 错误路径
- 移动端视觉回归覆盖关键页面
- 31 个已有 POM 未使用方法中 20+ 个被激活
