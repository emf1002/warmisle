# 目录结构清理设计文档

## 背景

项目经过多轮迭代，积累了一些设计产物和冗余文件。本次清理旨在精简仓库、统一工具入口、整理测试工具目录。

## 改进清单

### 1. 精简 docs/ 目录

**删除文件（设计迭代产物）：**
- `docs/dashboard-preview-v1.html`
- `docs/dashboard-preview-v2.html`
- `docs/dashboard-preview-v3.html`
- `docs/login-page-preview.html`
- `docs/login-page-preview-v2.html`
- `docs/login-page-preview-v3.html`
- `docs/logo-badge-dark.svg`
- `docs/logo-badge-light.svg`
- `docs/logo-concept-1-island-home.svg`
- `docs/logo-concept-2-warm-circle.svg`
- `docs/logo-concept-3-minimal.svg`
- `docs/logo-concept-4-badge.svg`
- `docs/logo-concept-5-modern-stack.svg`
- `docs/warmisle-prototype.html`

**保留文件：**
- `docs/prd.md` — 产品需求文档（CLAUDE.md 引用）
- `docs/superpowers/specs/*.md` — 技术设计文档
- `docs/superpowers/plans/*.md` — 实施计划

**理由：** 设计迭代 HTML 和 logo 概念图是设计阶段产物，最终版 logo 已在 `frontend/public/favicon.svg`。保留 `prd.md` 和 `superpowers/` 作为项目文档。

### 2. 删除冗余构建脚本

**删除文件：**
- `build.ps1`（PowerShell 构建脚本）

**保留文件：**
- `Makefile` — 跨平台构建入口（dev/build/test/lint/e2e）
- `build.bat` — Windows CMD 用户友好入口
- `build.sh` — Bash 用户友好入口

**理由：** Makefile 已覆盖所有构建场景，build.bat 和 build.sh 提供更简单的单命令构建。build.ps1 功能与 build.bat 重复。

### 3. 删除前端占位文件

**删除文件：**
- `frontend/public/vite.svg` — Vite 默认占位图标

**保留文件：**
- `frontend/public/favicon.svg` — 项目实际 logo

**理由：** favicon.svg 已在 `index.html` 中引用，vite.svg 是脚手架默认文件，不再需要。

### 4. 合并前端测试工具目录

**操作：**
- 将 `frontend/src/test-utils/auth-mock.ts` 移动到 `frontend/src/test/auth-mock.ts`
- 删除空的 `frontend/src/test-utils/` 目录
- 更新所有引用 `@/test-utils/auth-mock` 的测试文件为 `@/test/auth-mock`

**涉及文件：**
- `frontend/src/views/*/tests/*.test.ts`（约 8 个文件）

**理由：** 统一测试工具到 `src/test/` 目录，减少目录碎片化。

### 5. 提交未跟踪的源码文件

**提交文件：**
- `backend/internal/model/local_time.go` — LocalTime 类型（Ledger 模型依赖）
- `frontend/src/components/LogoIcon.vue` — Logo 图标组件
- `frontend/public/favicon.svg` — 项目 favicon

**不提交文件（设计产物，已在 #1 中删除）：**
- `docs/logo-*.svg`
- `docs/warmisle-prototype.html`

### 6. 更新 .gitignore

**新增规则：**
```
# IDE
.idea/
.vscode/

# OS
.DS_Store
Thumbs.db
```

## 执行顺序

1. 删除 docs/ 下的设计迭代文件（#1）
2. 删除 build.ps1（#2）
3. 删除 frontend/public/vite.svg（#3）
4. 移动前端测试工具并更新引用（#4）
5. 提交源码文件（#5）
6. 更新 .gitignore（#6）

## 验证

- `go build ./backend/...` 编译通过
- `cd frontend && npm test` 测试通过
- `git status` 确认无遗漏文件
