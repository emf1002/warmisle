# E2E 浏览器自动化测试设计

## 概述

为 WarmIsle 项目引入 Playwright E2E 测试，覆盖全部 11 个路由页面的功能测试和视觉回归测试。

## 目标

- **全模块功能覆盖**：登录、初始化、仪表盘、记账本、待办、愿望清单、论坛、成员管理、分类管理、个人中心
- **视觉回归测试**：页面级 + 组件级截图对比，防止 UI 样式被意外改动

## 技术选型

| 决策 | 选择 | 理由 |
|---|---|---|
| 测试框架 | Playwright | 内置截图对比、多浏览器支持、速度快、API 现代 |
| 测试模式 | Page Object Model (POM) | 11 个页面需要高可维护性，选择器集中管理 |
| DB 重置 | 后端测试端点 + 环境变量保护 | 不重启进程，速度快，受 HC_TEST_MODE 保护 |
| 认证 | fixture 注入 localStorage token | 复用现有 JWT 机制 |
| 视觉回归容差 | 1% 像素差异 | 平衡敏感度和误报 |

## 项目结构

```
e2e/
├── playwright.config.ts
├── package.json
├── fixtures/
│   ├── auth.fixture.ts           # 登录 fixture
│   └── db.fixture.ts             # 数据库重置 fixture
├── pages/
│   ├── base.page.ts              # BasePage 基类
│   ├── auth/
│   │   ├── login.page.ts
│   │   └── init.page.ts
│   ├── dashboard.page.ts
│   ├── ledger.page.ts
│   ├── todo.page.ts
│   ├── wish.page.ts
│   ├── forum.page.ts
│   ├── members.page.ts
│   ├── categories.page.ts
│   └── profile.page.ts
├── tests/
│   ├── auth.spec.ts
│   ├── dashboard.spec.ts
│   ├── ledger.spec.ts
│   ├── todo.spec.ts
│   ├── wish.spec.ts
│   ├── forum.spec.ts
│   ├── members.spec.ts
│   ├── categories.spec.ts
│   └── profile.spec.ts
└── screenshots/                  # 视觉回归基准图（自动生成）
```

`e2e/` 目录在项目根目录，与 `frontend/` 和 `backend/` 平级。

## Playwright 配置

`e2e/playwright.config.ts` 关键配置：

- Base URL: `http://localhost:8080`
- 浏览器: 默认 Chromium
- 测试目录: `./tests/`
- 截图对比: `toHaveScreenshot()` + `maxDiffPixelRatio: 0.01`
- 全局超时: 30s per test
- Web Server: 自动启动编译后的 Go 二进制，`HC_DB_PATH=./e2e-data/test.db` + `HC_TEST_MODE=true`

## 数据库重置机制

**策略**: 后端新增测试专用 reset 端点 + 环境变量保护

Playwright 的 `webServer` 是整个测试运行期间启动一次（不会每个测试文件重启），因此纯文件级删除方案不可行。改为在后端新增一个受环境变量保护的 reset 端点。

### 后端新增端点

`POST /api/test/reset` — 仅在 `HC_TEST_MODE=true` 环境变量下可用，否则返回 404。

功能：
- 清空所有业务数据表（保留 goose migration 记录）
- 重新插入预置种子数据（分类、标签）
- 返回 200 OK

路由注册：在 `routes/` 中条件注册，`os.Getenv("HC_TEST_MODE") == "true"` 时才挂载。

### 测试流程

1. Playwright `webServer` 配置启动后端（`HC_DB_PATH=./e2e-data/test.db`, `HC_TEST_MODE=true`）
2. 每个测试文件的 `beforeAll` 中：
   - 调用 `POST /api/test/reset` 清空数据并重建种子
   - 调用 `POST /api/init` 创建管理员账号
3. 测试运行期间使用该 DB
4. 测试结束后清理临时文件

**为什么不直接删除 DB 文件**：Playwright 的 `webServer` 只启动一次 server 进程，删除 DB 文件后 server 会崩溃，无法自动重启。

**为什么新增端点而不是扩展现有接口**：`/api/init` 的语义是首次初始化，重置是不同操作，混在一起会污染其职责。测试端点受环境变量保护，生产环境不可访问。

## 认证 Fixture

`e2e/fixtures/auth.fixture.ts` 提供 `authenticated` fixture：

```
authenticated fixture 流程：
1. 调用 POST /api/init 创建管理员（admin / test123）
2. 调用 POST /api/auth/login 获取 JWT token
3. 在浏览器 context 中注入 localStorage: { token: "<jwt>" }
4. 返回 { page, token }
```

使用方式：
```typescript
import { test, expect } from './fixtures/auth.fixture';

test('记账本页面', async ({ authenticated }) => {
  const { page } = authenticated;
  await page.goto('/#/ledger');
  // ...
});
```

## Page Object Model

### BasePage 基类

`e2e/pages/base.page.ts` 提供通用方法：

- `navigate(path)` — 导航到指定路由（带 hash history）
- `waitForPage()` — 等待页面加载完成
- `screenshot(name)` — 全页截图（视觉回归）
- `screenshotComponent(selector, name)` — 组件级截图
- `expectToast(message)` — 断言 Ant Design 提示消息

### 各页面 POM

每个 POM 类继承 BasePage，封装：

- **选择器**：页面关键元素的 data-testid 或语义化选择器
- **操作方法**：如 `ledger.addRecord(data)`, `todo.toggleItem(index)`
- **断言方法**：如 `ledger.expectRecordsCount(n)`, `dashboard.expectStatValue(key, value)`

### 测试文件中的 data-testid 约定

为确保选择器稳定，在前端组件中添加 `data-testid` 属性：

- `data-testid="page-title"` — 页面标题
- `data-testid="add-btn"` — 新增按钮
- `data-testid="list-item-{id}"` — 列表项
- `data-testid="modal-confirm"` — 弹窗确认按钮
- `data-testid="empty-state"` — 空状态组件

## 视觉回归测试

### 截图策略

| 级别 | 截图内容 | 示例 |
|---|---|---|
| 页面级 | 每个路由页面全页截图 | 仪表盘、记账列表、待办列表 |
| 组件级 | 关键交互组件截图 | 新增弹窗、表单验证错误、确认删除弹窗 |

### 基准图管理

- 存储路径: `e2e/screenshots/<module>/<test-name>.png`
- 首次运行: 自动生成基准图
- 更新基准: `npx playwright test --update-snapshots`
- 对比容差: `maxDiffPixelRatio: 0.01`（1% 像素差异）

### 截图对比示例

```typescript
test('仪表盘页面视觉回归', async ({ authenticated }) => {
  const { page } = authenticated;
  await page.goto('/#/');
  await page.waitForLoadState('networkidle');
  await expect(page).toHaveScreenshot('dashboard-full.png', {
    fullPage: true,
    maxDiffPixelRatio: 0.01,
  });
});

test('新增记账弹窗视觉回归', async ({ authenticated }) => {
  const { page } = authenticated;
  await page.goto('/#/ledger');
  await page.click('[data-testid="add-btn"]');
  await expect(page.locator('.ant-modal')).toHaveScreenshot('ledger-add-modal.png');
});
```

## 测试覆盖范围

### auth.spec.ts
- 首次访问重定向到 /init
- 初始化系统创建管理员
- 登录成功跳转仪表盘
- 登录失败显示错误提示
- 未登录访问受保护页面重定向到 /login

### dashboard.spec.ts
- 页面加载显示统计数据
- 快速记账入口功能
- 视觉回归截图

### ledger.spec.ts
- 记账列表展示
- 新增收入/支出记录
- 编辑记录
- 删除记录（二次确认）
- 筛选和搜索
- 分页
- 视觉回归截图

### todo.spec.ts
- 待办列表展示
- 新增待办
- 切换完成状态
- 编辑待办
- 删除待办
- 筛选（全部/进行中/已完成）
- 视觉回归截图

### wish.spec.ts
- 愿望列表展示
- 新增愿望
- 实现愿望
- 删除愿望
- 视觉回归截图

### forum.spec.ts
- 信息流列表
- 创建话题
- 话题详情
- 发表评论
- 回复评论（最多 2 层）
- 删除话题/评论
- 视觉回归截图

### members.spec.ts (admin only)
- 成员列表展示
- 添加成员
- 编辑成员
- 禁用/启用成员
- 角色切换
- 权限校验（非 admin 不可访问）
- 视觉回归截图

### categories.spec.ts (admin only)
- 分类列表展示
- 新增收入/支出分类
- 编辑分类
- 删除分类
- 预置分类展示
- 视觉回归截图

### profile.spec.ts
- 个人信息展示
- 修改昵称
- 修改密码
- 修改头像（emoji 选择）
- 视觉回归截图

## Makefile 集成

```makefile
e2e-install:
	cd e2e && npx playwright install --with-deps chromium

e2e:
	make build
	cd e2e && npx playwright test

e2e-ui:
	cd e2e && npx playwright test --ui

e2e-update:
	cd e2e && npx playwright test --update-snapshots

e2e-report:
	cd e2e && npx playwright show-report
```

## 依赖

`e2e/package.json`：
```json
{
  "name": "warmisle-e2e",
  "private": true,
  "devDependencies": {
    "@playwright/test": "^1.50.0"
  }
}
```

仅需 `@playwright/test` 一个依赖。

## Windows 兼容性

- Playwright 完全支持 Windows，自动下载 Chromium
- Makefile 中的 `rm`/`cd` 命令在 bash（Git Bash/MSYS2）下可用
- 也可直接用 `npx playwright test` 绕过 Makefile
