# Task 0: 基础设施 — 后端测试端点 + Fixture 扩展 + Playwright 配置

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为多用户权限测试和移动端测试建立基础设施。

**Files:**
- Create: `backend/internal/handler/test_create_member.go`
- Modify: `backend/internal/routes/router.go:15-18`
- Modify: `e2e/fixtures/db.fixture.ts`
- Modify: `e2e/fixtures/auth.fixture.ts`
- Modify: `e2e/playwright.config.ts`

---

## Step 1: 创建后端测试端点 `POST /api/test/create-member`

遵循 `test_reset.go` 的模式：HC_TEST_MODE 守卫 + `pkg.DB` + `pkg.Success`。

**Create:** `backend/internal/handler/test_create_member.go`

```go
package handler

import (
	"net/http"
	"os"
	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

type createTestMemberRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// TestCreateMember 创建测试用 member 用户并返回 JWT token
// 仅在 HC_TEST_MODE=true 时可用
func TestCreateMember(c *gin.Context) {
	if os.Getenv("HC_TEST_MODE") != "true" {
		pkg.Error(c, http.StatusNotFound, 404, "not found")
		return
	}

	var req createTestMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}

	db := pkg.DB

	// 检查用户名是否已存在
	var count int64
	db.Model(&model.Member{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		pkg.Error(c, http.StatusConflict, 40901, "用户名已被使用")
		return
	}

	// 加密密码
	hashedPwd, err := pkg.HashPassword(req.Password)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "密码加密失败")
		return
	}

	// 创建 member
	member := model.Member{
		Username: req.Username,
		Password: hashedPwd,
		Name:     req.Name,
		Role:     "member",
		Avatar:   "👤",
	}
	if err := db.Create(&member).Error; err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "创建失败")
		return
	}

	// 签发 token
	token, err := pkg.GenerateToken(member.ID, member.Username, member.Role)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "token 签发失败")
		return
	}

	pkg.Success(c, gin.H{
		"token":     token,
		"member_id": member.ID,
	})
}
```

## Step 2: 注册路由

**Modify:** `backend/internal/routes/router.go` — 在测试端点 block 中新增一行

```go
if os.Getenv("HC_TEST_MODE") == "true" {
    api.POST("/test/reset", handler.TestReset)
    api.POST("/test/seed-ledgers", handler.TestSeedLedgers)
    api.POST("/test/create-member", handler.TestCreateMember) // 新增
}
```

## Step 3: 扩展 db.fixture.ts — 新增 `createMember` 函数

**Modify:** `e2e/fixtures/db.fixture.ts` — 在文件末尾新增

```typescript
export interface CreateMemberResult {
  code: number;
  message: string;
  data: {
    token: string;
    member_id: number;
  };
}

export async function createMember(
  options: { username?: string; password?: string; name?: string } = {}
): Promise<CreateMemberResult> {
  const res = await fetch(`${BASE_URL}/api/test/create-member`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: options.username ?? 'member1',
      password: options.password ?? 'test123',
      name: options.name ?? '成员一',
    }),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Create member failed: ${res.status} ${text}`);
  }

  return res.json();
}
```

## Step 4: 扩展 auth.fixture.ts — 新增 `memberContext` fixture

**Modify:** `e2e/fixtures/auth.fixture.ts`

```typescript
import { test as base, type Page, type BrowserContext } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, createMember, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';

type AuthFixtures = {
  authenticated: {
    page: Page;
    token: string;
    seedLedgers: (options?: SeedLedgersOptions) => Promise<SeedLedgersResult>;
  };
  // 新增：member 角色的已认证上下文
  memberContext: {
    page: Page;
    token: string;
  };
};

export const test = base.extend<AuthFixtures>({
  authenticated: async ({ page, browser }, use) => {
    await resetDatabase();
    const { token } = await initAdmin();

    await page.addInitScript((token) => {
      localStorage.setItem('token', token);
    }, token);

    const boundSeedLedgers = (options?: SeedLedgersOptions) =>
      seedLedgers(token, options);

    await use({ page, token, seedLedgers: boundSeedLedgers });
  },

  memberContext: async ({ browser }, use) => {
    // 创建 member 用户
    const memberResult = await createMember();
    const memberToken = memberResult.data.token;

    // 创建新的 browser context 和 page
    const context = await browser.newContext();
    const page = await context.newPage();

    await page.addInitScript((token) => {
      localStorage.setItem('token', token);
    }, memberToken);

    await use({ page, token: memberToken });

    // 清理
    await context.close();
  },
});

export { expect } from '@playwright/test';
```

## Step 5: 扩展 Playwright 配置 — 新增 mobile project

**Modify:** `e2e/playwright.config.ts`

```typescript
import path from 'path';
import { defineConfig, devices } from '@playwright/test';

const rootDir = path.resolve(__dirname, '..');
const binaryPath = path.join(rootDir, 'dist', 'warmisle.exe');

export default defineConfig({
  testDir: './tests',
  snapshotDir: './__snapshots__',
  fullyParallel: false,
  forbidOnly: true,
  retries: 0,
  workers: 1,
  reporter: [['html', { open: 'never' }]],
  use: {
    baseURL: 'http://localhost:8080',
    screenshot: 'only-on-failure',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: {
        viewport: { width: 1280, height: 720 },
      },
    },
    {
      name: 'mobile',
      use: {
        ...devices['iPhone 13'],
      },
      grep: /@mobile/,
    },
  ],
  webServer: {
    command: binaryPath,
    cwd: rootDir,
    port: 8080,
    timeout: 30000,
    reuseExistingServer: true,
    env: {
      HC_DB_PATH: './e2e-data/test.db',
      HC_TEST_MODE: 'true',
      HC_PORT: '8080',
    },
  },
});
```

## Step 6: 验证基础设施

1. `cd backend && go build ./cmd/server/` — 确认 Go 编译通过
2. `cd e2e && npx playwright test auth.spec.ts --project=chromium` — 确认现有测试不受影响
3. 手动验证：启动服务后 `curl -X POST http://localhost:8080/api/test/create-member -H 'Content-Type: application/json' -d '{"username":"test","password":"test123","name":"测试"}'` 应返回 token

## Step 7: Commit

```bash
git add backend/internal/handler/test_create_member.go backend/internal/routes/router.go e2e/fixtures/db.fixture.ts e2e/fixtures/auth.fixture.ts e2e/playwright.config.ts
git commit -m "feat(e2e): add member fixture, test endpoint, and mobile project config"
```
