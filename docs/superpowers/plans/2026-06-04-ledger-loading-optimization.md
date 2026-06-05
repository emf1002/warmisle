# 记账本加载方式优化 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace offset/limit pagination with cursor-based pagination for the ledger list, fix summary-filter inconsistency, add Pinia caches for categories/members, and improve loading UX (skeleton, infinite scroll, debounce, abort).

**Architecture:** Backend: cursor pagination with date-group补全 strategy, summary query now applies category/creator filters. Frontend: Pinia stores for categories/members, skeleton screen replacing spinner, IntersectionObserver infinite scroll, 300ms debounce on filter changes, AbortController for request cancellation.

**Tech Stack:** Go, Gin, GORM, SQLite, Vue 3, Ant Design Vue, TypeScript, Pinia, Vitest

**Spec:** `docs/superpowers/specs/2026-06-04-ledger-loading-optimization-design.md`

---

## File Map

| File | Change |
|------|--------|
| `backend/internal/repository/ledger.go` | Add `CursorData`, `encodeCursor`/`DecodeCursor`; update `LedgerFilter` (add `Limit`/`Cursor`, remove `Page`/`PageSize`); update `ListResult` (add `NextCursor`/`HasMore`, remove `Total`/`Page`/`PageSize`); extract `computeSummary` (now applies filters), `applyOptionalFilters`; rewrite `List` with cursor query + 补全 strategy |
| `backend/internal/handler/ledger.go` | Update `List` handler: bind `limit`/`cursor` instead of `page`/`page_size`; decode cursor via `repository.DecodeCursor`; remove page default logic |
| `backend/internal/service/ledger_test.go` | Update `List` test: use `Limit` instead of `Page`/`PageSize`, assert `HasMore`/`NextCursor` instead of `Total` |
| `backend/internal/handler/ledger_test.go` | Update `List` test: use `limit` query param, assert `has_more`/`groups` in response |
| `frontend/src/stores/categories.ts` | **New** — Pinia store with `categories` ref, `loaded` flag, `fetchCategories()`, `reset()` |
| `frontend/src/stores/members.ts` | **New** — Pinia store with `members` ref, `loaded` flag, `fetchMembers()`, `reset()` |
| `frontend/src/api/ledger.ts` | Update `getLedgers` params: `limit`/`cursor` replace `page`/`page_size`; add optional `signal` param |
| `frontend/src/utils/request.ts` | Add early return for `ERR_CANCELED` in response error interceptor |
| `frontend/src/views/ledger/Index.vue` | Skeleton screen replacing spinner; cursor state (`nextCursor`/`hasMore`); IntersectionObserver infinite scroll (remove "加载更多" button); 300ms debounce on filter changes; AbortController for request cancellation; use Pinia stores for categories/members |

### No changes needed
- `backend/internal/service/ledger.go` — pure passthrough for `List`, no code changes
- `backend/internal/model/` — no struct changes
- Database migrations — existing index `(deleted_at, occurred_at)` covers cursor query
- Routes — API path `/api/ledgers` unchanged

---

### Task 1: Backend repo — add cursor types, update structs, add helpers

**Files:**
- Modify: `backend/internal/repository/ledger.go`

- [ ] **Step 1: Add CursorData struct and encode/decode helpers**

Add imports `"encoding/base64"` and `"encoding/json"` to the import block.

Add above the `LedgerRepo` struct:

```go
type CursorData struct {
	OccurredAt string `json:"occurred_at"` // "2006-01-02 15:04:05"
	ID         uint   `json:"id"`
}

func EncodeCursor(c CursorData) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeCursor(s string) (*CursorData, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var c CursorData
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
```

- [ ] **Step 2: Update LedgerFilter struct**

Replace the current `LedgerFilter`:
```go
type LedgerFilter struct {
	StartDate  string
	EndDate    string
	CategoryID *uint
	CreatorID  *uint
	Page       int
	PageSize   int
}
```
With:
```go
type LedgerFilter struct {
	StartDate  string
	EndDate    string
	CategoryID *uint
	CreatorID  *uint
	Limit      int        // default 20
	Cursor     *CursorData // nil = first page
}
```

- [ ] **Step 3: Update ListResult struct**

Replace:
```go
type ListResult struct {
	Summary  LedgerSummary `json:"summary"`
	Groups   []LedgerGroup `json:"groups"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}
```
With:
```go
type ListResult struct {
	Summary    LedgerSummary `json:"summary"`
	Groups     []LedgerGroup `json:"groups"`
	NextCursor *string       `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}
```

- [ ] **Step 4: Add helper methods**

Add `applyOptionalFilters`, `computeSummary`, and `calcDailyTotal` methods:

```go
func (r *LedgerRepo) applyOptionalFilters(query *gorm.DB, filter LedgerFilter) *gorm.DB {
	if filter.CategoryID != nil {
		query = query.Where("ledgers.category_id = ?", *filter.CategoryID)
	}
	if filter.CreatorID != nil {
		query = query.Where("ledgers.creator_id = ?", *filter.CreatorID)
	}
	return query
}

func (r *LedgerRepo) computeSummary(filter LedgerFilter) (LedgerSummary, error) {
	var summary LedgerSummary
	type summaryRow struct {
		Type   string
		Amount int64
	}
	var rows []summaryRow

	query := pkg.DB.Table("ledgers").
		Select("categories.type, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", filter.StartDate, filter.EndDate)

	query = r.applyOptionalFilters(query, filter)

	if err := query.Group("categories.type").Scan(&rows).Error; err != nil {
		return summary, err
	}
	for _, row := range rows {
		if row.Type == "income" {
			summary.Income = row.Amount
		} else if row.Type == "expense" {
			summary.Expense = row.Amount
		}
	}
	summary.Balance = summary.Income - summary.Expense
	return summary, nil
}

func (r *LedgerRepo) calcDailyTotal(items []LedgerWithAssoc) int64 {
	var total int64
	for _, item := range items {
		if item.Category.Type == "income" {
			total += item.Amount
		} else {
			total -= item.Amount
		}
	}
	return total
}
```

Note: `applyOptionalFilters` uses `*gorm.DB` — add `"gorm.io/gorm"` to imports if not already present. Check current imports; GORM is already used indirectly via `pkg.DB`, so you may need to import it explicitly for the type annotation. Alternatively, use the concrete type from `pkg.DB` — since `pkg.DB` is `*gorm.DB`, the import is needed.

- [ ] **Step 5: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/internal/repository/`
Expected: Compilation errors in `List` method (expected — will rewrite in Task 2). The new types and helpers should compile.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/repository/ledger.go
git commit -m "refactor(backend): add cursor types, update LedgerFilter/ListResult structs"
```

---

### Task 2: Backend repo — rewrite List method with cursor pagination + 补全

**Files:**
- Modify: `backend/internal/repository/ledger.go`

- [ ] **Step 1: Rewrite the List method**

Replace the entire `List` method with:

```go
func (r *LedgerRepo) List(filter LedgerFilter) (*ListResult, error) {
	// 1. Compute summary (now applies category/creator filters)
	summary, err := r.computeSummary(filter)
	if err != nil {
		return nil, err
	}

	// 2. Determine limit
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	// 3. Build base query
	query := pkg.DB.Model(&model.Ledger{}).
		Preload("Category").
		Preload("Creator").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", filter.StartDate, filter.EndDate)
	query = r.applyOptionalFilters(query, filter)

	// 4. Apply cursor condition
	if filter.Cursor != nil {
		query = query.Where(
			"(ledgers.occurred_at < ? OR (ledgers.occurred_at = ? AND ledgers.id < ?))",
			filter.Cursor.OccurredAt, filter.Cursor.OccurredAt, filter.Cursor.ID,
		)
	}

	// 5. Fetch limit + 1 records
	var records []model.Ledger
	if err := query.
		Order("ledgers.occurred_at DESC, ledgers.id DESC").
		Limit(limit + 1).
		Find(&records).Error; err != nil {
		return nil, err
	}

	// 6. Determine hasMore
	hasMore := len(records) > limit
	var extraRecord *model.Ledger
	if hasMore {
		extraRecord = &records[limit]
		records = records[:limit]
	}

	// 7. Group by date
	groups := r.groupByDate(records)

	// 8. Check if last date group is incomplete (补全 strategy)
	if hasMore && len(groups) > 0 && extraRecord != nil {
		lastGroup := &groups[len(groups)-1]
		extraDate := time.Time(extraRecord.OccurredAt).Format("2006-01-02")
		if extraDate == lastGroup.Date {
			// Last group was cut off — fetch remaining records for this date
			minId := lastGroup.Items[len(lastGroup.Items)-1].ID
			nextDay := time.Time(extraRecord.OccurredAt).AddDate(0, 0, 1).Format("2006-01-02")

			补全Query := pkg.DB.Model(&model.Ledger{}).
				Preload("Category").
				Preload("Creator").
				Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", lastGroup.Date, nextDay).
				Where("ledgers.id < ?", minId)
			补全Query = r.applyOptionalFilters(补全Query, filter)

			var extraRecords []model.Ledger
			if err := 补全Query.
				Order("ledgers.occurred_at DESC, ledgers.id DESC").
				Find(&extraRecords).Error; err != nil {
				return nil, err
			}

			// Merge into last group
			for _, er := range extraRecords {
				lastGroup.Items = append(lastGroup.Items, LedgerWithAssoc{
					Ledger:   er,
					Category: er.Category,
					Creator:  er.Creator,
				})
			}
			lastGroup.DailyTotal = r.calcDailyTotal(lastGroup.Items)
		}
	}

	// 9. Compute nextCursor
	var nextCursor *string
	if hasMore {
		var lastItem LedgerWithAssoc
		if len(groups) > 0 {
			lastGroup := groups[len(groups)-1]
			lastItem = lastGroup.Items[len(lastGroup.Items)-1]
		}
		cursor := EncodeCursor(CursorData{
			OccurredAt: time.Time(lastItem.OccurredAt).Format("2006-01-02 15:04:05"),
			ID:         lastItem.ID,
		})
		nextCursor = &cursor
	}

	return &ListResult{
		Summary:    summary,
		Groups:     groups,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
```

Note: The variable name uses a Chinese identifier `补全Query` for readability. Go supports Unicode identifiers. If the linter rejects it, rename to `completionQuery`.

- [ ] **Step 2: Remove the now-unused standalone summary block**

The old `List` method had an inline summary computation at the top. That is now handled by `computeSummary`. Make sure no dead code remains.

- [ ] **Step 3: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/internal/repository/`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/repository/ledger.go
git commit -m "feat(backend): rewrite List with cursor pagination and date-group补全 strategy"
```

---

### Task 3: Backend handler — update List params

**Files:**
- Modify: `backend/internal/handler/ledger.go`

- [ ] **Step 1: Update List handler request struct**

Replace the `List` method's request struct from:
```go
var req struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	CategoryID *uint  `form:"category_id"`
	CreatorID  *uint  `form:"creator_id"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}
```
To:
```go
var req struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	CategoryID *uint  `form:"category_id"`
	CreatorID  *uint  `form:"creator_id"`
	Limit      int    `form:"limit"`
	Cursor     string `form:"cursor"`
}
```

- [ ] **Step 2: Replace default logic**

Remove the `req.Page` and `req.PageSize` default blocks. Replace with:
```go
if req.Limit <= 0 {
	req.Limit = 20
}
```

Keep the existing `StartDate`/`EndDate` defaults unchanged.

- [ ] **Step 3: Decode cursor and build filter**

Replace the filter construction block with:
```go
filter := repository.LedgerFilter{
	StartDate:  req.StartDate,
	EndDate:    req.EndDate,
	CategoryID: req.CategoryID,
	CreatorID:  req.CreatorID,
	Limit:      req.Limit,
}

if req.Cursor != "" {
	cursor, err := repository.DecodeCursor(req.Cursor)
	if err != nil {
		pkg.Error(c, 400, 40001, "游标格式错误")
		return
	}
	filter.Cursor = cursor
}
```

- [ ] **Step 4: Update response nil-guard**

Change:
```go
if result.Groups == nil {
	result.Groups = []repository.LedgerGroup{}
}
```
Keep this unchanged — it still applies.

- [ ] **Step 5: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/...`
Expected: PASS (all backend packages compile).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/handler/ledger.go
git commit -m "refactor(backend): update ledger List handler for cursor pagination"
```

---

### Task 4: Backend tests — update service and handler tests

**Files:**
- Modify: `backend/internal/service/ledger_test.go`
- Modify: `backend/internal/handler/ledger_test.go`

#### Service tests

- [ ] **Step 1: Update TestLedgerService_List_ByMonth**

Change the filter from:
```go
result, err := svc.List(repository.LedgerFilter{StartDate: "2026-05-01", EndDate: "2026-06-01", Page: 1, PageSize: 20})
```
To:
```go
result, err := svc.List(repository.LedgerFilter{StartDate: "2026-05-01", EndDate: "2026-06-01", Limit: 20})
```

Change the assertion from:
```go
assert.Equal(t, int64(1), result.Total)
```
To:
```go
assert.Len(t, result.Groups, 1)
assert.False(t, result.HasMore)
assert.Nil(t, result.NextCursor)
```

- [ ] **Step 2: Add cursor pagination test**

Add a new test `TestLedgerService_List_CursorPagination`:

```go
func TestLedgerService_List_CursorPagination(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	// Create 3 records on the same date
	base := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	svc.Create(1000, "A", cat.ID, base, m1.ID)
	svc.Create(2000, "B", cat.ID, base.Add(time.Hour), m1.ID)
	svc.Create(3000, "C", cat.ID, base.Add(2*time.Hour), m1.ID)

	// First page: limit 2
	page1, err := svc.List(repository.LedgerFilter{
		StartDate: "2026-06-01", EndDate: "2026-06-02", Limit: 2,
	})
	require.NoError(t, err)
	assert.True(t, page1.HasMore)
	assert.NotNil(t, page1.NextCursor)
	// 补全: all 3 records are on the same date, so the补全 fetches the 3rd
	totalItems := 0
	for _, g := range page1.Groups {
		totalItems += len(g.Items)
	}
	assert.Equal(t, 3, totalItems) // 补全 ensures complete date group

	// Second page using cursor
	cursor, err := repository.DecodeCursor(*page1.NextCursor)
	require.NoError(t, err)
	page2, err := svc.List(repository.LedgerFilter{
		StartDate: "2026-06-01", EndDate: "2026-06-02", Limit: 2, Cursor: cursor,
	})
	require.NoError(t, err)
	assert.False(t, page2.HasMore)
	assert.Len(t, page2.Groups, 0) // all records returned in page 1 via 补全
}
```

- [ ] **Step 3: Add summary filter alignment test**

Add `TestLedgerService_List_SummaryWithFilters`:

```go
func TestLedgerService_List_SummaryWithFilters(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	incomeCat := testutil.CreateTestCategory(pkg.DB, "income", "工资", "💰", 1)

	svc.Create(3000, "餐饮", cat.ID, time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC), m1.ID)
	svc.Create(10000, "工资", incomeCat.ID, time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), m1.ID)

	// Filter by expense category — summary should only reflect that category
	catID := cat.ID
	result, err := svc.List(repository.LedgerFilter{
		StartDate: "2026-06-01", EndDate: "2026-07-01", Limit: 20, CategoryID: &catID,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.Summary.Income)   // income filtered out
	assert.Equal(t, int64(3000), result.Summary.Expense) // only餐饮
}
```

#### Handler tests

- [ ] **Step 4: Update TestHandler_Ledger_List_ByMonth**

Change the query string from the current date range params to include `limit`:
```go
w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s&limit=20", startDate, endDate), nil, memberToken)
```

Update assertions to check `has_more` and `groups` instead of `total`/`page`/`page_size`:
```go
data := testutil.ParseDataMap(resp)
assert.NotNil(t, data["summary"])
assert.NotNil(t, data["groups"])
assert.NotNil(t, data["has_more"])
```

- [ ] **Step 5: Run all backend tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/... -v`
Expected: All tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/service/ledger_test.go backend/internal/handler/ledger_test.go
git commit -m "test(backend): update ledger tests for cursor pagination and summary filters"
```

---

### Task 5: Frontend — create Pinia stores for categories and members

**Files:**
- Create: `frontend/src/stores/categories.ts`
- Create: `frontend/src/stores/members.ts`

- [ ] **Step 1: Create categories store**

Create `frontend/src/stores/categories.ts`:

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories } from '@/api/category'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<any[]>([])
  const loaded = ref(false)

  async function fetchCategories() {
    if (loaded.value) return categories.value
    const res: any = await getCategories()
    categories.value = res.data || []
    loaded.value = true
    return categories.value
  }

  function reset() {
    categories.value = []
    loaded.value = false
  }

  return { categories, loaded, fetchCategories, reset }
})
```

- [ ] **Step 2: Create members store**

Create `frontend/src/stores/members.ts`:

```typescript
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMembers } from '@/api/member'

export const useMembersStore = defineStore('members', () => {
  const members = ref<any[]>([])
  const loaded = ref(false)

  async function fetchMembers() {
    if (loaded.value) return members.value
    const res: any = await getMembers()
    members.value = res.data || []
    loaded.value = true
    return members.value
  }

  function reset() {
    members.value = []
    loaded.value = false
  }

  return { members, loaded, fetchMembers, reset }
})
```

- [ ] **Step 3: Verify frontend compilation**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit`
Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/stores/categories.ts frontend/src/stores/members.ts
git commit -m "feat(frontend): add Pinia stores for categories and members caching"
```

---

### Task 6: Frontend — update API and request interceptor

**Files:**
- Modify: `frontend/src/api/ledger.ts`
- Modify: `frontend/src/utils/request.ts`

- [ ] **Step 1: Update getLedgers params**

Replace the `getLedgers` function:
```typescript
export function getLedgers(params: {
  start_date?: string
  end_date?: string
  category_id?: number
  creator_id?: number
  page?: number
  page_size?: number
}) {
  return request.get('/ledgers', { params })
}
```
With:
```typescript
export function getLedgers(
  params: {
    start_date?: string
    end_date?: string
    category_id?: number
    creator_id?: number
    limit?: number
    cursor?: string
  },
  signal?: AbortSignal,
) {
  return request.get('/ledgers', { params, signal })
}
```

- [ ] **Step 2: Fix request.ts error interceptor for AbortController**

In `frontend/src/utils/request.ts`, add an early return for canceled requests in the error interceptor. Change:
```typescript
error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
    }
    message.error('❌ ' + (error.message || '网络错误'))
    return Promise.reject(error)
  }
```
To:
```typescript
error => {
    if (error.code === 'ERR_CANCELED') {
      return Promise.reject(error)
    }
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
    }
    message.error('❌ ' + (error.message || '网络错误'))
    return Promise.reject(error)
  }
```

- [ ] **Step 3: Run frontend tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Expected: All tests pass (API tests may need updating in the next task if they reference old params).

- [ ] **Step 4: Commit**

```bash
git add frontend/src/api/ledger.ts frontend/src/utils/request.ts
git commit -m "refactor(frontend): update ledger API for cursor params, fix abort handling"
```

---

### Task 7: Frontend — rewrite Index.vue (skeleton, cursor, infinite scroll, debounce, abort)

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue`

#### Template changes

- [ ] **Step 1: Replace loading spinner with skeleton**

Replace:
```html
<div v-if="loading" class="loading-state">
  <a-spin />
  <span style="margin-left: 8px">加载中...</span>
</div>
```
With:
```html
<div v-if="loading" class="skeleton-state">
  <a-skeleton active :paragraph="{ rows: 6 }" />
</div>
```

- [ ] **Step 2: Remove "加载更多" button, add sentinel div**

Replace:
```html
<div v-if="hasMore" class="load-more">
  <a-button @click="loadMore" :loading="loadingMore">加载更多</a-button>
</div>
```
With:
```html
<div ref="sentinelRef" v-if="hasMore && !loading" class="load-sentinel">
  <a-spin size="small" />
</div>
```

#### Script changes

- [ ] **Step 3: Update imports**

Add to imports:
```typescript
import { ref, reactive, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { useCategoriesStore } from '@/stores/categories'
import { useMembersStore } from '@/stores/members'
```

Remove the direct API imports for categories and members:
```typescript
// Remove:
import { getCategories } from '@/api/category'
import { getMembers } from '@/api/member'
```

Remove the `getLedgers` import's old signature — it now accepts `signal`.

- [ ] **Step 4: Replace pagination state with cursor state**

Remove:
```typescript
const total = ref(0)
const page = ref(1)
const pageSize = 20
```

Add:
```typescript
const nextCursor = ref<string | null>(null)
const hasMore = ref(false)
const pageSize = 20
```

Remove the `hasMore` computed property (it's now a ref).

- [ ] **Step 5: Add AbortController and debounce utility**

Add:
```typescript
let abortController: AbortController | null = null

function debounce<F extends (...args: any[]) => any>(fn: F, delay: number): F {
  let timer: ReturnType<typeof setTimeout>
  return ((...args: any[]) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }) as any
}

const debouncedFetchLedgers = debounce(() => fetchLedgers(), 300)
```

- [ ] **Step 6: Add IntersectionObserver setup**

Add:
```typescript
const sentinelRef = ref<HTMLDivElement | null>(null)
let observer: IntersectionObserver | null = null

watch(sentinelRef, (el, oldEl) => {
  if (oldEl && observer) observer.unobserve(oldEl)
  if (el && observer) observer.observe(el)
})
```

- [ ] **Step 7: Rewrite fetchLedgers**

Replace the existing `fetchLedgers` function with:
```typescript
async function fetchLedgers(isLoadMore = false) {
  if (!isLoadMore) {
    abortController?.abort()
    abortController = new AbortController()
    loading.value = true
    nextCursor.value = null
    groups.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const params: Record<string, unknown> = {
      start_date: dateRange.value[0].format('YYYY-MM-DD'),
      end_date: dateRange.value[1].add(1, 'day').format('YYYY-MM-DD'),
      limit: pageSize,
    }
    if (nextCursor.value && isLoadMore) {
      params.cursor = nextCursor.value
    }
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.creator_id) params.creator_id = filters.creator_id

    const signal = isLoadMore ? undefined : abortController?.signal
    const res: any = await getLedgers(params as any, signal)
    const data = res.data

    if (isLoadMore) {
      groups.value = [...groups.value, ...(data.groups || [])]
    } else {
      groups.value = data.groups || []
    }

    summary.value = data.summary || { income: 0, expense: 0, balance: 0 }
    nextCursor.value = data.next_cursor || null
    hasMore.value = data.has_more || false
  } catch (e: any) {
    if (e?.code === 'ERR_CANCELED') return
    if (!isLoadMore) {
      groups.value = []
      summary.value = { income: 0, expense: 0, balance: 0 }
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}
```

- [ ] **Step 8: Rewrite loadMore**

Replace:
```typescript
async function loadMore() {
  page.value++
  await fetchLedgers(true)
}
```
With:
```typescript
async function loadMore() {
  if (!hasMore.value || loadingMore.value) return
  await fetchLedgers(true)
}
```

- [ ] **Step 9: Update onFilterChange to use debounce**

Replace:
```typescript
function onFilterChange() {
  fetchLedgers()
}
```
With:
```typescript
function onFilterChange() {
  nextCursor.value = null
  hasMore.value = false
  debouncedFetchLedgers()
}
```

- [ ] **Step 10: Update clearFilters**

Replace:
```typescript
function clearFilters() {
  filters.category_id = undefined
  filters.creator_id = undefined
  fetchLedgers()
}
```
With:
```typescript
function clearFilters() {
  filters.category_id = undefined
  filters.creator_id = undefined
  nextCursor.value = null
  hasMore.value = false
  fetchLedgers()
}
```

- [ ] **Step 11: Update onMounted to use Pinia stores and set up observer**

Replace the `onMounted` block:
```typescript
onMounted(async () => {
  try {
    const [catRes, memRes]: any[] = await Promise.all([
      getCategories(),
      getMembers(),
    ])
    categories.value = catRes.data || []
    members.value = memRes.data || []
  } catch {
    // error handled by interceptor
  }
  fetchLedgers()
})
```
With:
```typescript
const categoriesStore = useCategoriesStore()
const membersStore = useMembersStore()

onMounted(async () => {
  try {
    const [cats, mems] = await Promise.all([
      categoriesStore.fetchCategories(),
      membersStore.fetchMembers(),
    ])
    categories.value = cats
    members.value = mems
  } catch {
    // error handled by interceptor
  }
  fetchLedgers()

  // Set up IntersectionObserver after first render
  await nextTick()
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting && hasMore.value && !loading.value && !loadingMore.value) {
        loadMore()
      }
    },
    { threshold: 0.1 },
  )
  if (sentinelRef.value) {
    observer.observe(sentinelRef.value)
  }
})

onUnmounted(() => {
  observer?.disconnect()
  abortController?.abort()
})
```

#### CSS changes

- [ ] **Step 12: Add skeleton and sentinel styles**

Replace the `.load-more` CSS:
```css
.load-more {
  text-align: center;
  padding: 16px 0;
}
```
With:
```css
.skeleton-state {
  padding: 16px;
}

.load-sentinel {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}
```

- [ ] **Step 13: Run frontend tests and fix**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Update any failing tests (view tests may need `a-skeleton` stub, remove `load-more` button references, update mock data to use cursor response format).

- [ ] **Step 14: Commit**

```bash
git add frontend/src/views/ledger/Index.vue
git commit -m "feat(frontend): skeleton, infinite scroll, debounce, abort — ledger loading overhaul"
```

---

### Task 8: Full verification

- [ ] **Step 1: Run all backend tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/... -v`
Expected: All tests pass.

- [ ] **Step 2: Run all frontend tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Expected: All tests pass.

- [ ] **Step 3: Verify build**

Run: `cd D:/Projects/my_projects/home-center-v1 && make build`
Expected: Build succeeds.

- [ ] **Step 4: Smoke test**

Start dev server: `cd D:/Projects/my_projects/home-center-v1 && make dev`
Open the ledger page and verify:
1. Skeleton appears while loading
2. Records display in date groups
3. Scrolling to bottom triggers auto-load (no button)
4. Changing category/creator filter shows correct filtered summary
5. Rapid filter switching does not cause flickering or stale data

- [ ] **Step 5: Commit (if any fixups needed)**

```bash
git add -A
git commit -m "fix: address test/build issues from ledger loading optimization"
```
