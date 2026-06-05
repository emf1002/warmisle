# Ledger E2E Load Tests Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add bulk data seeding and 4 new E2E test scenarios for the ledger module covering infinite scroll, summary stats, date range filtering, and category filtering.

**Architecture:** A new test-only backend endpoint (`POST /api/test/seed-ledgers`) bulk-inserts deterministic ledger data. The E2E page object gets fixed and extended with new helper methods. A new Playwright spec file contains 4 test scenarios that use the seed endpoint for data setup.

**Tech Stack:** Go (Gin, GORM), Vue 3, Playwright, TypeScript

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `backend/internal/handler/test_seed_ledgers.go` | Create | Bulk seed endpoint handler |
| `backend/internal/routes/router.go:16` | Modify | Register seed route |
| `frontend/src/views/ledger/Index.vue` | Modify | Add 8 data-testid attributes |
| `e2e/pages/ledger.page.ts` | Modify | Fix selectCategory, add 7 new methods |
| `e2e/fixtures/db.fixture.ts` | Modify | Add seedLedgers function |
| `e2e/fixtures/auth.fixture.ts` | Modify | Extend authenticated fixture with seedLedgers |
| `e2e/tests/ledger-load.spec.ts` | Create | 4 test scenarios |

---

### Task 1: Backend — Create bulk seed endpoint

**Files:**
- Create: `backend/internal/handler/test_seed_ledgers.go`
- Modify: `backend/internal/routes/router.go:16`

- [ ] **Step 1: Create the seed handler file**

Create `backend/internal/handler/test_seed_ledgers.go`:

```go
package handler

import (
	"math"
	"net/http"
	"os"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

type seedLedgersRequest struct {
	Count     int    `json:"count" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type seedSummary struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Balance int64 `json:"balance"`
}

type seedResult struct {
	Count              int         `json:"count"`
	Summary            seedSummary `json:"summary"`
	ExpenseCategoryCnt int         `json:"expense_category_count"`
	IncomeCategoryCnt  int         `json:"income_category_count"`
}

// TestSeedLedgers 批量播种账单数据（仅测试模式可用）
func TestSeedLedgers(c *gin.Context) {
	if os.Getenv("HC_TEST_MODE") != "true" {
		pkg.Error(c, http.StatusNotFound, 404, "not found")
		return
	}

	var req seedLedgersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Count <= 0 {
		pkg.Error(c, 400, 40001, "count 必须大于 0")
		return
	}

	db := pkg.DB

	// 解析日期范围，默认近 7 天
	endDate := time.Now()
	if req.EndDate != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			pkg.Error(c, 400, 40001, "end_date 格式错误，应为 YYYY-MM-DD")
			return
		}
	}

	startDate := endDate.AddDate(0, 0, -6) // 默认 7 天前
	if req.StartDate != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			pkg.Error(c, 400, 40001, "start_date 格式错误，应为 YYYY-MM-DD")
			return
		}
	}

	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days <= 0 {
		days = 1
	}

	// 获取预置分类
	var expenseCats, incomeCats []model.Category
	db.Where("type = ? AND deleted_at IS NULL", "expense").Order("sort_order").Find(&expenseCats)
	db.Where("type = ? AND deleted_at IS NULL", "income").Order("sort_order").Find(&incomeCats)

	if len(expenseCats) == 0 || len(incomeCats) == 0 {
		pkg.Error(c, 400, 40002, "请先调用 /api/test/reset 创建预置分类")
		return
	}

	// 获取第一个成员作为 creator
	var member model.Member
	if err := db.Where("deleted_at IS NULL").Order("id").First(&member).Error; err != nil {
		pkg.Error(c, 400, 40003, "请先创建至少一个成员")
		return
	}

	// 批量创建记录
	ledgers := make([]model.Ledger, 0, req.Count)
	expIdx := 0
	incIdx := 0

	for i := 0; i < req.Count; i++ {
		dayOffset := i * days / req.Count
		occurred := time.Date(
			startDate.Year(), startDate.Month(), startDate.Day()+dayOffset,
			12, 0, 0, 0, time.UTC,
		)

		var cat model.Category
		var amount int64

		if i%5 == 4 {
			// 每 5 条中 1 条收入
			cat = incomeCats[incIdx%len(incomeCats)]
			amount = 50000 // 500 元
			incIdx++
		} else {
			// 4 条支出
			cat = expenseCats[expIdx%len(expenseCats)]
			amount = int64((i%10 + 1) * 100) // 100~1000 分
			expIdx++
		}

		ledgers = append(ledgers, model.Ledger{
			Amount:     amount,
			Note:       "测试账单-" + time.Now().Format("150405") + "-" + string(rune('A'+i%26)),
			CategoryID: cat.ID,
			CreatorID:  member.ID,
			OccurredAt: model.FromTime(occurred),
		})
	}

	// 批量插入
	if err := db.CreateInBatches(&ledgers, 100).Error; err != nil {
		pkg.Error(c, 500, 50001, "播种数据失败: "+err.Error())
		return
	}

	// 计算 summary
	var income, expense int64
	db.Table("ledgers").
		Select("COALESCE(SUM(ledgers.amount), 0)").
		Joins("JOIN categories ON categories.id = ledgers.category_id").
		Where("categories.type = ? AND ledgers.deleted_at IS NULL", "income").
		Scan(&income)

	db.Table("ledgers").
		Select("COALESCE(SUM(ledgers.amount), 0)").
		Joins("JOIN categories ON categories.id = ledgers.category_id").
		Where("categories.type = ? AND ledgers.deleted_at IS NULL", "expense").
		Scan(&expense)

	pkg.Success(c, seedResult{
		Count: req.Count,
		Summary: seedSummary{
			Income:  income,
			Expense: expense,
			Balance: income - expense,
		},
		ExpenseCategoryCnt: len(expenseCats),
		IncomeCategoryCnt:  len(incomeCats),
	})
}
```

- [ ] **Step 2: Register the route**

In `backend/internal/routes/router.go`, add after line 16 (`api.POST("/test/reset", handler.TestReset)`):

```go
		api.POST("/test/seed-ledgers", handler.TestSeedLedgers)
```

- [ ] **Step 3: Verify backend compiles**

Run: `cd D:/Projects/my_projects/home-center-v1/backend && go build ./...`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add backend/internal/handler/test_seed_ledgers.go backend/internal/routes/router.go
git commit -m "feat(backend): add test seed-ledgers endpoint for E2E bulk data"
```

---

### Task 2: Frontend — Add data-testid attributes

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue`

- [ ] **Step 1: Add data-testid to summary bar**

In `frontend/src/views/ledger/Index.vue`, make these edits:

Line 24 — add `data-testid="summary-income"`:
```html
        <span class="income-amount" data-testid="summary-income">+{{ formatYuan(summary.income) }}</span>
```

Line 27 — add `data-testid="summary-expense"`:
```html
        <span class="expense-amount" data-testid="summary-expense">-{{ formatYuan(summary.expense) }}</span>
```

Lines 32-33 — add `data-testid="summary-balance"`:
```html
        <span :class="summary.balance >= 0 ? 'income-amount' : 'expense-amount'" data-testid="summary-balance">
```

- [ ] **Step 2: Add data-testid to filter selects**

Line 40 — add `data-testid="filter-category"` to the category `<a-select>`:
```html
      <a-select
        v-model:value="filters.category_id"
        placeholder="全部分类"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="filter-category"
      >
```

Line 60 — add `data-testid="filter-creator"` to the creator `<a-select>`:
```html
      <a-select
        v-model:value="filters.creator_id"
        placeholder="全部创建者"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="filter-creator"
      >
```

- [ ] **Step 3: Add data-testid to date group and daily total**

Line 87 — add `data-testid="date-group"`:
```html
      <div v-for="group in groups" :key="group.date" class="date-group" data-testid="date-group">
```

Line 91 — add `data-testid="daily-total"`:
```html
          <span :class="group.daily_total >= 0 ? 'income-amount' : 'expense-amount'" class="date-total" data-testid="daily-total">
```

- [ ] **Step 4: Add data-testid to sentinel**

Line 130 — add `data-testid="load-sentinel"`:
```html
    <div ref="sentinelRef" v-if="hasMore" class="load-sentinel" data-testid="load-sentinel">
```

- [ ] **Step 5: Verify frontend builds**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npm run build`
Expected: build succeeds

- [ ] **Step 6: Commit**

```bash
git add frontend/src/views/ledger/Index.vue
git commit -m "feat(frontend): add data-testid attrs to ledger for E2E testing"
```

---

### Task 3: E2E — Fix and extend LedgerPage object

**Files:**
- Modify: `e2e/pages/ledger.page.ts`

- [ ] **Step 1: Rewrite the page object**

Replace the entire content of `e2e/pages/ledger.page.ts` with:

```typescript
// e2e/pages/ledger.page.ts
import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class LedgerPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/ledger');
  }

  async expectOnLedger() {
    await expect(this.page).toHaveURL(/\/#\/ledger/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.getByTestId('ledger-modal')).toBeVisible();
  }

  /** 在弹窗网格选择器中点击分类（按文本匹配） */
  async pickCategory(name: string) {
    await this.page.locator('.category-pick-item', { hasText: name }).click();
  }

  async fillAmount(amount: string) {
    await this.page.getByTestId('amount-input').click();
    await this.page.getByTestId('amount-input').locator('input').fill(amount);
  }

  async fillNote(note: string) {
    await this.page.getByTestId('note-input').fill(note);
  }

  async submit() {
    await this.page.getByTestId('submit-btn').click();
  }

  async expectRecordCount(count: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items).toHaveCount(count);
  }

  async clickRecord(index: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await items.nth(index).click();
  }

  async deleteRecord() {
    await this.page.getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  async clearFilters() {
    await this.page.getByTestId('clear-filters').click();
  }

  /** 通过滤选栏选择分类 */
  async filterByCategory(name: string) {
    await this.page.getByTestId('filter-category').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 通过滤选栏选择创建者 */
  async filterByCreator(name: string) {
    await this.page.getByTestId('filter-creator').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 滚动到 sentinel 触发无限加载 */
  async scrollToLoadMore() {
    const sentinel = this.page.getByTestId('load-sentinel');
    await sentinel.scrollIntoViewIfNeeded();
  }

  /** 断言汇总栏文本 */
  async expectSummary(values: { income?: string; expense?: string; balance?: string }) {
    if (values.income !== undefined) {
      await expect(this.page.getByTestId('summary-income')).toContainText(values.income);
    }
    if (values.expense !== undefined) {
      await expect(this.page.getByTestId('summary-expense')).toContainText(values.expense);
    }
    if (values.balance !== undefined) {
      await expect(this.page.getByTestId('summary-balance')).toContainText(values.balance);
    }
  }

  /** 断言日期分组数量 */
  async expectDateGroupCount(count: number) {
    await expect(this.page.getByTestId('date-group')).toHaveCount(count);
  }

  /** 断言所有记录总数 */
  async expectTotalItemCount(count: number) {
    const items = this.page.getByTestId(/^ledger-item-/);
    await expect(items).toHaveCount(count);
  }

  /** 断言第 n 个日期组的每日小计包含指定文本 */
  async expectDailyTotal(index: number, text: string) {
    await expect(this.page.getByTestId('daily-total').nth(index)).toContainText(text);
  }
}
```

- [ ] **Step 2: Verify TypeScript compiles**

Run: `cd D:/Projects/my_projects/home-center-v1/e2e && npx tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add e2e/pages/ledger.page.ts
git commit -m "fix(e2e): fix selectCategory and extend LedgerPage with new helpers"
```

---

### Task 4: E2E — Extend fixtures with seedLedgers

**Files:**
- Modify: `e2e/fixtures/db.fixture.ts`
- Modify: `e2e/fixtures/auth.fixture.ts`

- [ ] **Step 1: Add seedLedgers to db.fixture.ts**

Add the following function to `e2e/fixtures/db.fixture.ts` after the `initAdmin` function:

```typescript
export interface SeedLedgersOptions {
  count?: number;
  startDate?: string;
  endDate?: string;
}

export interface SeedLedgersResult {
  code: number;
  message: string;
  data: {
    count: number;
    summary: { income: number; expense: number; balance: number };
    expense_category_count: number;
    income_category_count: number;
  };
}

export async function seedLedgers(
  token: string,
  options?: SeedLedgersOptions
): Promise<SeedLedgersResult> {
  const res = await fetch(`${BASE_URL}/api/test/seed-ledgers`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      count: options?.count ?? 35,
      start_date: options?.startDate,
      end_date: options?.endDate,
    }),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Seed ledgers failed: ${res.status} ${text}`);
  }

  return res.json();
}
```

- [ ] **Step 2: Extend auth.fixture.ts**

Replace the entire content of `e2e/fixtures/auth.fixture.ts`:

```typescript
import { test as base, type Page } from '@playwright/test';
import { resetDatabase, initAdmin, seedLedgers, type SeedLedgersOptions, type SeedLedgersResult } from './db.fixture';

type AuthFixtures = {
  authenticated: {
    page: Page;
    token: string;
    seedLedgers: (options?: SeedLedgersOptions) => Promise<SeedLedgersResult>;
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
});

export { expect } from '@playwright/test';
```

- [ ] **Step 3: Update existing ledger.spec.ts to use new pickCategory**

In `e2e/tests/ledger.spec.ts`, replace all `selectCategory` calls with `pickCategory`. There are 3 occurrences on lines 18, 30, 42:

```
await ledger.selectCategory('餐饮');
→ await ledger.pickCategory('餐饮');
```
```
await ledger.selectCategory('工资');
→ await ledger.pickCategory('工资');
```
```
await ledger.selectCategory('餐饮');
→ await ledger.pickCategory('餐饮');
```

- [ ] **Step 4: Verify TypeScript compiles**

Run: `cd D:/Projects/my_projects/home-center-v1/e2e && npx tsc --noEmit`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add e2e/fixtures/db.fixture.ts e2e/fixtures/auth.fixture.ts e2e/tests/ledger.spec.ts
git commit -m "feat(e2e): add seedLedgers fixture and update auth fixture"
```

---

### Task 5: E2E — Write the 4 test scenarios

**Files:**
- Create: `e2e/tests/ledger-load.spec.ts`

- [ ] **Step 1: Create the test file**

Create `e2e/tests/ledger-load.spec.ts`:

```typescript
// e2e/tests/ledger-load.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { LedgerPage } from '../pages/ledger.page';

test.describe('记账本 — 负载场景', () => {
  test('无限滚动：首次加载 20 条，滚动后加载全部 35 条', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 35 });
    expect(result.code).toBe(0);

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 首次加载 20 条（默认 limit=20）
    await ledger.expectTotalItemCount(20);

    // sentinel 可见，表示还有更多数据
    await expect(page.getByTestId('load-sentinel')).toBeVisible();

    // 滚动触发加载更多
    await ledger.scrollToLoadMore();

    // 等待全部 35 条加载完成
    await ledger.expectTotalItemCount(35);

    // sentinel 消失
    await expect(page.getByTestId('load-sentinel')).not.toBeVisible();

    // 7 天 = 7 个日期分组
    await ledger.expectDateGroupCount(7);
  });

  test('统计汇总：收入/支出/结余与 seed 数据一致', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 35 });
    const { summary } = result.data;

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载完成
    await ledger.expectTotalItemCount(20);

    // 断言汇总值（页面首次加载即包含 summary，基于全量数据计算）
    const expectedIncome = (summary.income / 100).toFixed(2);
    const expectedExpense = (summary.expense / 100).toFixed(2);
    const expectedBalance = (summary.balance / 100).toFixed(2);

    await ledger.expectSummary({
      income: expectedIncome,
      expense: expectedExpense,
    });

    // 结余符号取决于正负
    if (summary.balance >= 0) {
      await ledger.expectSummary({ balance: `+${expectedBalance}` });
    } else {
      await ledger.expectSummary({ balance: `-${Math.abs(Number(expectedBalance)).toFixed(2)}` });
    }

    // 滚动加载全部后，汇总值不变
    await ledger.scrollToLoadMore();
    await ledger.expectTotalItemCount(35);
    await ledger.expectSummary({
      income: expectedIncome,
      expense: expectedExpense,
    });
  });

  test('日期范围筛选：选择子范围后只显示匹配记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    await seedLedgers({
      count: 35,
      startDate: '2026-06-01',
      endDate: '2026-06-07',
    });

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待全量数据加载
    await ledger.scrollToLoadMore();
    await ledger.expectTotalItemCount(35);

    // 使用 RangePicker 选择 6 月 3 日 ~ 6 月 5 日（3 天）
    const rangePicker = page.getByTestId('date-range-picker');
    await rangePicker.click();

    // 点击 6 月 3 日
    await page.locator('.ant-picker-cell[title="2026-06-03"]').click();
    // 点击 6 月 5 日
    await page.locator('.ant-picker-cell[title="2026-06-05"]').click();

    // 等待数据重新加载
    await page.waitForLoadState('networkidle');

    // 记录数应少于 35
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeLessThan(35);
    expect(count).toBeGreaterThan(0);

    // sentinel 不应存在（筛选后数据量小于 limit）
    await expect(page.getByTestId('load-sentinel')).not.toBeVisible();
  });

  test('分类筛选：选择分类后只显示该分类记录', async ({ authenticated }) => {
    const { page, seedLedgers } = authenticated;
    const result = await seedLedgers({ count: 35 });
    const fullSummary = result.data.summary;

    const ledger = new LedgerPage(page);
    await ledger.goto();

    // 等待首次加载
    await ledger.expectTotalItemCount(20);

    // 筛选"餐饮"分类
    await ledger.filterByCategory('餐饮');

    // 等待数据重新加载
    await page.waitForLoadState('networkidle');

    // 记录数应少于 20（只显示餐饮分类）
    const items = page.getByTestId(/^ledger-item-/);
    const count = await items.count();
    expect(count).toBeLessThan(20);
    expect(count).toBeGreaterThan(0);

    // 验证所有可见记录都是餐饮分类
    const categoryNames = page.locator('.item-cat-name');
    const catCount = await categoryNames.count();
    for (let i = 0; i < catCount; i++) {
      await expect(categoryNames.nth(i)).toHaveText('餐饮');
    }

    // 清除筛选后恢复全部
    await ledger.clearFilters();
    await page.waitForLoadState('networkidle');
    await ledger.expectTotalItemCount(20);
  });
});
```

- [ ] **Step 2: Verify TypeScript compiles**

Run: `cd D:/Projects/my_projects/home-center-v1/e2e && npx tsc --noEmit`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add e2e/tests/ledger-load.spec.ts
git commit -m "feat(e2e): add ledger load test scenarios (scroll, summary, filters)"
```

---

### Task 6: Run full E2E suite and fix issues

**Files:**
- Potentially modify any files from previous tasks if tests fail

- [ ] **Step 1: Build the binary**

Run: `cd D:/Projects/my_projects/home-center-v1 && make build`
Expected: `dist/warmisle.exe` built successfully

- [ ] **Step 2: Run the full E2E suite**

Run: `cd D:/Projects/my_projects/home-center-v1/e2e && npx playwright test`
Expected: all 50 tests pass (46 existing + 4 new)

- [ ] **Step 3: Fix any failures**

If tests fail:
1. Read the error output carefully
2. Check if data-testid selectors match the actual rendered HTML
3. Check if the seed endpoint returns expected data
4. Fix and re-run

- [ ] **Step 4: Final commit (if fixes needed)**

```bash
git add -A
git commit -m "fix(e2e): fix test failures in ledger load scenarios"
```
