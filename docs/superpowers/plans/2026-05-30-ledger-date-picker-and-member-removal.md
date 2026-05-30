# 记账本日期选择器 & 移除关联成员 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the ledger month navigator with a date range picker (RangePicker + month presets), and remove the ledger_members many-to-many association while keeping creator_id for permissions.

**Architecture:** Two independent changes to the ledger module: (1) frontend date range picker replacing month arrows, backend `month` param replaced with `start_date`/`end_date`; (2) full removal of `ledger_members` pivot table and all related code across model/repository/service/handler/frontend.

**Tech Stack:** Go, Gin, GORM, SQLite, Vue 3, Ant Design Vue, TypeScript, Vitest

---

## File Map

| File | Change |
|------|--------|
| `backend/internal/model/ledger.go` | Remove `Members` field, delete `LedgerMember` struct |
| `backend/internal/repository/ledger.go` | Remove member filter/preload/CRUD, change `Month` to `StartDate`/`EndDate` |
| `backend/internal/service/ledger.go` | Remove member params/validation, remove `memberRepo` |
| `backend/internal/service/errors.go` | Remove `ErrNoMembers` |
| `backend/internal/handler/ledger.go` | Remove `member_ids`/`member_id` fields, change `month` to `start_date`/`end_date` |
| `backend/internal/service/ledger_test.go` | Update all tests: remove member assertions, use date range |
| `backend/internal/handler/ledger_test.go` | Update all tests: remove member_ids, use date range |
| `backend/migrations/003_drop_ledger_members.up.sql` | New: `DROP TABLE IF EXISTS ledger_members` |
| `frontend/src/api/ledger.ts` | Remove `member_id`/`member_ids`, add `start_date`/`end_date` |
| `frontend/src/api/__tests__/ledger.test.ts` | Update tests for new params |
| `frontend/src/views/ledger/Index.vue` | Remove member UI, replace month nav with RangePicker |
| `frontend/src/views/ledger/__tests__/Index.test.ts` | Update tests for new UI |

---

### Task 1: Backend model — remove LedgerMember and Members association

**Files:**
- Modify: `backend/internal/model/ledger.go`

- [ ] **Step 1: Remove Members field from Ledger struct**

In `backend/internal/model/ledger.go`, remove line 23 (`Members  []Member ...`) from the `Ledger` struct. Keep `Category` and `Creator` fields.

- [ ] **Step 2: Delete LedgerMember struct**

Remove lines 28-33 (the `LedgerMember` struct and its `TableName` method) entirely.

- [ ] **Step 3: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/...`
Expected: Compilation errors in repository, service, handler (expected — will fix in later tasks). Verify the model package itself compiles: `go build ./backend/internal/model/`

- [ ] **Step 4: Commit**

```bash
git add backend/internal/model/ledger.go
git commit -m "refactor(backend): remove LedgerMember model and Members association from Ledger"
```

---

### Task 2: Backend repository — remove member logic, change month to date range

**Files:**
- Modify: `backend/internal/repository/ledger.go`

- [ ] **Step 1: Update LedgerFilter struct**

Replace the `Month string` field with `StartDate string` and `EndDate string`. Remove the `MemberID *uint` field. The struct should be:

```go
type LedgerFilter struct {
	StartDate  string // "2026-05-01"
	EndDate    string // "2026-06-01" (exclusive upper bound)
	CategoryID *uint
	CreatorID  *uint
	Page       int
	PageSize   int
}
```

- [ ] **Step 2: Update LedgerWithAssoc struct**

Remove the `Members []model.Member` field. The struct should be:

```go
type LedgerWithAssoc struct {
	model.Ledger
	Category model.Category `json:"category"`
	Creator  model.Member   `json:"creator"`
}
```

- [ ] **Step 3: Update List method — summary query**

Change the summary query's WHERE clause from `strftime('%Y-%m', ledgers.occurred_at) = ?` to `ledgers.occurred_at >= ? AND ledgers.occurred_at < ?` using `filter.StartDate` and `filter.EndDate`. The args change from one string to two strings.

- [ ] **Step 4: Update List method — main query**

In the main query:
- Remove `Preload("Members")` (line 78)
- Change the month WHERE clause to date range: `ledgers.occurred_at >= ? AND ledgers.occurred_at < ?` with `filter.StartDate` and `filter.EndDate`
- Remove the `MemberID` filter block (lines 87-90, the sub-query)

- [ ] **Step 5: Update groupByDate method**

In the `groupByDate` method, remove `Members: l.Members` from the `LedgerWithAssoc` construction (line 132).

- [ ] **Step 6: Update Create method**

Change the signature from `Create(ledger *model.Ledger, memberIDs []uint)` to `Create(ledger *model.Ledger)`. Remove the `for _, memberID := range memberIDs` loop that inserts `LedgerMember` rows (lines 165-173). The method should just do `pkg.DB.Create(ledger)`.

- [ ] **Step 7: Update Update method**

Change the signature from `Update(ledger *model.Ledger, memberIDs []uint)` to `Update(ledger *model.Ledger)`. Remove the `LedgerMember` delete-and-reinsert block (lines 182-196). The method should just do `pkg.DB.Save(ledger)`.

- [ ] **Step 8: Update FindByID method**

Remove `Preload("Members")` (line 204). Remove `Members: ledger.Members` from the return value (line 213).

- [ ] **Step 9: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/internal/repository/`
Expected: May still have errors in service/handler that depend on old signatures.

- [ ] **Step 10: Commit**

```bash
git add backend/internal/repository/ledger.go
git commit -m "refactor(backend): remove member filter/CRUD, switch month to date range in ledger repo"
```

---

### Task 3: Backend service — remove member logic

**Files:**
- Modify: `backend/internal/service/ledger.go`
- Modify: `backend/internal/service/errors.go`

- [ ] **Step 1: Remove ErrNoMembers from errors.go**

In `backend/internal/service/errors.go`, remove the `ErrNoMembers` line (line 15).

- [ ] **Step 2: Remove memberRepo from LedgerService**

In `backend/internal/service/ledger.go`, remove the `memberRepo *repository.MemberRepo` field from `LedgerService` (line 22) and its initialization in `NewLedgerService()` (line 30).

- [ ] **Step 3: Update Create method signature**

Change from:
```go
func (s *LedgerService) Create(amount int64, note string, categoryID uint, memberIDs []uint, occurredAt time.Time, creatorID uint) (*repository.LedgerWithAssoc, error) {
```
to:
```go
func (s *LedgerService) Create(amount int64, note string, categoryID uint, occurredAt time.Time, creatorID uint) (*repository.LedgerWithAssoc, error) {
```

Remove the `memberIDs` validation block (lines 56-58: `if len(memberIDs) == 0 { return nil, ErrNoMembers }`).

Change the repo call from `s.repo.Create(ledger, memberIDs)` to `s.repo.Create(ledger)`.

- [ ] **Step 4: Update Update method signature**

Change from:
```go
func (s *LedgerService) Update(id uint, amount *int64, note *string, categoryID *uint, memberIDs []uint, occurredAt *time.Time, currentMemberID uint, currentRole string) (*repository.LedgerWithAssoc, error) {
```
to:
```go
func (s *LedgerService) Update(id uint, amount *int64, note *string, categoryID *uint, occurredAt *time.Time, currentMemberID uint, currentRole string) (*repository.LedgerWithAssoc, error) {
```

Remove the entire member replacement block (lines 109-119):
```go
// Members: if nil, keep existing; if empty, error; otherwise replace
if memberIDs != nil {
    if len(memberIDs) == 0 {
        return nil, ErrNoMembers
    }
} else {
    memberIDs = make([]uint, len(existing.Members))
    for i, m := range existing.Members {
        memberIDs[i] = m.ID
    }
}
```

Change the repo call from `s.repo.Update(&existing.Ledger, memberIDs)` to `s.repo.Update(&existing.Ledger)`.

- [ ] **Step 5: Remove unused imports**

Remove the `"warmisle/internal/repository"` import if no longer used (check — it's still used for `LedgerFilter` and `ListResult`, so keep it). Remove the `"warmisle/internal/model"` import if no longer used (check — it's still used for `model.Ledger` and `model.FromTime`, so keep it).

- [ ] **Step 6: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/internal/service/`
Expected: May still have errors in handler.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/service/ledger.go backend/internal/service/errors.go
git commit -m "refactor(backend): remove member validation and memberRepo from ledger service"
```

---

### Task 4: Backend handler — remove member params, switch to date range

**Files:**
- Modify: `backend/internal/handler/ledger.go`

- [ ] **Step 1: Update List handler — replace month with date range**

In the `List` method, change the request struct from:
```go
var req struct {
    Month      string `form:"month"`
    MemberID   *uint  `form:"member_id"`
    CategoryID *uint  `form:"category_id"`
    CreatorID  *uint  `form:"creator_id"`
    Page       int    `form:"page"`
    PageSize   int    `form:"page_size"`
}
```
to:
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

Replace the `getCurrentMonth()` default logic (lines 42-44) with date range defaults:
```go
if req.StartDate == "" {
    req.StartDate = time.Now().Format("2006-01") + "-01"
}
if req.EndDate == "" {
    // First day of next month
    now := time.Now()
    firstOfNext := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
    req.EndDate = firstOfNext.Format("2006-01-02")
}
```

Update the filter construction to use `StartDate` and `EndDate` instead of `Month`, and remove `MemberID`:
```go
filter := repository.LedgerFilter{
    StartDate:  req.StartDate,
    EndDate:    req.EndDate,
    CategoryID: req.CategoryID,
    CreatorID:  req.CreatorID,
    Page:       req.Page,
    PageSize:   req.PageSize,
}
```

Remove the `getCurrentMonth()` function (lines 22-24) as it's no longer needed.

- [ ] **Step 2: Update Create handler — remove member_ids**

In `createLedgerRequest` (lines 93-99), remove `MemberIDs []uint` field.

In the `Create` method:
- Remove the `member_ids` validation block (lines 113-116: `if len(req.MemberIDs) == 0 { ... }`)
- Change the service call from `h.svc.Create(amountCents, req.Note, req.CategoryID, req.MemberIDs, occurredAt, creatorID)` to `h.svc.Create(amountCents, req.Note, req.CategoryID, occurredAt, creatorID)`
- Remove the `serviceError{service.ErrNoMembers, ...}` from `handleServiceError` (line 144)

- [ ] **Step 3: Update Update handler — remove member_ids**

In `updateLedgerRequest` (lines 152-158), remove `MemberIDs []uint` field.

In the `Update` method, change the service call from `h.svc.Update(uint(id), amountCents, req.Note, req.CategoryID, req.MemberIDs, occurredAt, currentMemberID, currentRole)` to `h.svc.Update(uint(id), amountCents, req.Note, req.CategoryID, occurredAt, currentMemberID, currentRole)`.

Remove the `serviceError{service.ErrNoMembers, ...}` from `handleServiceError` (line 210).

- [ ] **Step 4: Remove unused import**

The `"warmisle/internal/repository"` import is still needed for `repository.LedgerFilter`. Keep it.

- [ ] **Step 5: Verify compilation**

Run: `cd D:/Projects/my_projects/home-center-v1 && go build ./backend/...`
Expected: PASS (all backend packages compile).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/handler/ledger.go
git commit -m "refactor(backend): remove member_ids from ledger handler, switch month to date range"
```

---

### Task 5: Database migration — drop ledger_members table

**Files:**
- Create: `backend/migrations/003_drop_ledger_members.up.sql`

- [ ] **Step 1: Create migration file**

Create `backend/migrations/003_drop_ledger_members.up.sql` with:
```sql
-- +goose Up
-- 移除记账记录关联成员表
DROP TABLE IF EXISTS ledger_members;

-- +goose Down
-- 不支持降级
```

Note: This project uses goose for migrations. The `+goose` annotations are required.

- [ ] **Step 2: Commit**

```bash
git add backend/migrations/003_drop_ledger_members.up.sql
git commit -m "feat(backend): add migration to drop ledger_members table"
```

---

### Task 6: Backend tests — update service tests

**Files:**
- Modify: `backend/internal/service/ledger_test.go`

- [ ] **Step 1: Update seedLedgerFixtures**

Remove `member2` from the return values. The function currently returns `(member1, member2 model.Member, cat model.Category)`. Change to return only `(member1 model.Member, cat model.Category)` since member2 is only used for permission tests (which use a separately created admin/member, not the fixture's member2).

Wait — actually, `member2` IS used in `TestLedgerService_Update_ByNonCreator` and `TestLedgerService_Delete_ByNonCreator`. Keep returning both members for those tests, but they don't need to be associated via `ledger_members` anymore. The permission check uses `CreatorID`, not `Members`.

Actually, let me re-read the tests more carefully. The `seedLedgerFixtures` returns `member1, member2, cat`. `member2` is used in permission tests. These tests don't need `ledger_members` at all — they just need `member2` to exist as a different member. So keep `seedLedgerFixtures` as-is, but remove `member_ids` from all `Create` calls.

- [ ] **Step 2: Update Create calls — remove memberIDs parameter**

All `svc.Create(...)` calls need to remove the `[]uint{...}` memberIDs argument. For example:

`svc.Create(3550, "午餐", cat.ID, []uint{m1.ID, m2.ID}, time.Now(), m1.ID)` becomes `svc.Create(3550, "午餐", cat.ID, time.Now(), m1.ID)`

Update ALL `svc.Create` calls in the file (there are ~10).

- [ ] **Step 3: Remove TestLedgerService_Create_NoMembers test**

Delete the `TestLedgerService_Create_NoMembers` test entirely (lines 64-72) — this validation no longer exists.

- [ ] **Step 4: Update assertion in TestLedgerService_Create_Success**

Remove `assert.Len(t, result.Members, 2)` (line 42) — the `Members` field no longer exists.

- [ ] **Step 5: Update List test — use date range**

In `TestLedgerService_List_ByMonth`, change the filter from `LedgerFilter{Month: "2026-05", ...}` to `LedgerFilter{StartDate: "2026-05-01", EndDate: "2026-06-01", ...}`.

- [ ] **Step 6: Update Update calls — remove memberIDs parameter**

All `svc.Update(...)` calls need to remove the `memberIDs` argument (the `nil` or `[]uint{...}` after `categoryID`). For example:

`svc.Update(created.ID, &newAmount, &newNote, nil, nil, nil, m1.ID, "member")` becomes `svc.Update(created.ID, &newAmount, &newNote, nil, nil, m1.ID, "member")`

Update ALL `svc.Update` calls in the file.

- [ ] **Step 7: Run tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/internal/service/ -run TestLedger -v`
Expected: All tests pass.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/service/ledger_test.go
git commit -m "test(backend): update ledger service tests for member removal and date range"
```

---

### Task 7: Backend tests — update handler tests

**Files:**
- Modify: `backend/internal/handler/ledger_test.go`

- [ ] **Step 1: Update Create test — remove member_ids from request body**

In `TestHandler_Ledger_Create_Success` (line 24), change the body from:
```go
body := fmt.Sprintf(`{"amount":3550,"note":"午餐","category_id":%d,"member_ids":[%d],"occurred_at":"%s"}`, cat.ID, member.ID, ...)
```
to:
```go
body := fmt.Sprintf(`{"amount":3550,"note":"午餐","category_id":%d,"occurred_at":"%s"}`, cat.ID, ...)
```

Do the same for all Create test bodies: `TestHandler_Ledger_Create_ZeroAmount`, `TestHandler_Ledger_Create_NegativeAmount`, `TestHandler_Ledger_Create_CategoryNotFound`.

- [ ] **Step 2: Remove TestHandler_Ledger_Create_NoMembers test**

Delete the `TestHandler_Ledger_Create_NoMembers` test entirely (lines 60-71) — this validation no longer exists.

- [ ] **Step 3: Update List test — use date range params**

In `TestHandler_Ledger_List_ByMonth`, change the query from `?month=%s` to `?start_date=%s&end_date=%s`. Compute the first day of the current month and the first day of next month:

```go
now := time.Now()
startDate := now.Format("2006-01") + "-01"
firstOfNext := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
endDate := firstOfNext.Format("2006-01-02")
w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s", startDate, endDate), nil, memberToken)
```

Also remove `pkg.DB.Create(&model.LedgerMember{...})` lines from tests that create ledger fixtures directly (lines 97, 117, 145, 162).

- [ ] **Step 4: Update remaining tests that create ledger fixtures**

In tests like `TestHandler_Ledger_Update_ByCreator`, `TestHandler_Ledger_Update_ByNonCreator`, `TestHandler_Ledger_Delete_ByAdmin`, remove the `pkg.DB.Create(&model.LedgerMember{...})` line after creating a ledger.

- [ ] **Step 5: Run tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/internal/handler/ -run TestHandler_Ledger -v`
Expected: All tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/handler/ledger_test.go
git commit -m "test(backend): update ledger handler tests for member removal and date range"
```

---

### Task 8: Frontend API — update ledger.ts types

**Files:**
- Modify: `frontend/src/api/ledger.ts`

- [ ] **Step 1: Update getLedgers params**

Change the `getLedgers` function params from:
```ts
params: {
  month?: string
  member_id?: number
  category_id?: number
  creator_id?: number
  page?: number
  page_size?: number
}
```
to:
```ts
params: {
  start_date?: string
  end_date?: string
  category_id?: number
  creator_id?: number
  page?: number
  page_size?: number
}
```

- [ ] **Step 2: Update createLedger data**

Change the `createLedger` function data from:
```ts
data: {
  amount: number
  note?: string
  category_id: number
  member_ids: number[]
  occurred_at?: string
}
```
to:
```ts
data: {
  amount: number
  note?: string
  category_id: number
  occurred_at?: string
}
```

- [ ] **Step 3: Update updateLedger data**

Change the `updateLedger` function data from:
```ts
data: {
  amount?: number
  note?: string
  category_id?: number
  member_ids?: number[]
  occurred_at?: string
}
```
to:
```ts
data: {
  amount?: number
  note?: string
  category_id?: number
  occurred_at?: string
}
```

- [ ] **Step 4: Run frontend tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run src/api/__tests__/ledger.test.ts`
Expected: Tests will fail (expected — will fix in next task).

- [ ] **Step 5: Commit**

```bash
git add frontend/src/api/ledger.ts
git commit -m "refactor(frontend): update ledger API types — remove member fields, add date range"
```

---

### Task 9: Frontend tests — update API tests

**Files:**
- Modify: `frontend/src/api/__tests__/ledger.test.ts`

- [ ] **Step 1: Update getLedgers test — use date range params**

In the `getLedgers` test, change:
```ts
const res = await getLedgers({ month: '2026-05', page: 1, page_size: 20 })
expect(mockRequest.get).toHaveBeenCalledWith('/ledgers', {
  params: { month: '2026-05', page: 1, page_size: 20 },
})
```
to:
```ts
const res = await getLedgers({ start_date: '2026-05-01', end_date: '2026-06-01', page: 1, page_size: 20 })
expect(mockRequest.get).toHaveBeenCalledWith('/ledgers', {
  params: { start_date: '2026-05-01', end_date: '2026-06-01', page: 1, page_size: 20 },
})
```

- [ ] **Step 2: Update getLedgers mock response — remove members**

In the mock response, remove the `members` field from items:
```ts
members: [{ id: 1, name: '管理员' }],
```
And remove the assertion:
```ts
expect(item.members[0].name).toBe('管理员')
```

- [ ] **Step 3: Update createLedger test — remove member_ids**

In the `createLedger` test, change the test data from:
```ts
const newItem = {
  amount: 5000,
  note: '晚餐',
  category_id: 2,
  member_ids: [1],
  occurred_at: '2026-05-23',
}
```
to:
```ts
const newItem = {
  amount: 5000,
  note: '晚餐',
  category_id: 2,
  occurred_at: '2026-05-23',
}
```

- [ ] **Step 4: Run tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run src/api/__tests__/ledger.test.ts`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/api/__tests__/ledger.test.ts
git commit -m "test(frontend): update ledger API tests for member removal and date range"
```

---

### Task 10: Frontend view — remove member UI, add RangePicker

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue`

- [ ] **Step 1: Remove member filter from filter bar**

Remove lines 42-52 (the `a-select` for `filters.member_id`). Keep the category and creator selects.

- [ ] **Step 2: Remove member_ids from form state**

In the `form` reactive (lines 342-348), remove `member_ids: [] as number[]`.

- [ ] **Step 3: Remove member_ids from Filters interface and filters reactive**

In the `Filters` interface (lines 313-317), remove `member_id: number | undefined`. In the `filters` reactive (lines 336-340), remove `member_id: undefined`.

- [ ] **Step 4: Remove member_ids validation from handleSubmit**

In `handleSubmit()` (lines 471-513), remove lines 480-483:
```ts
if (form.member_ids.length === 0) {
  message.error('❌ 请至少选择一位关联成员')
  return
}
```

- [ ] **Step 5: Remove member_ids from payload in handleSubmit**

In the payload construction (lines 487-492), remove `member_ids: form.member_ids`.

- [ ] **Step 6: Remove member_ids from openCreate and openEdit**

In `openCreate()` (lines 442-451), remove `form.member_ids = []`.

In `openEdit()` (lines 453-463), remove `form.member_ids = record.members ? record.members.map((m) => m.id) : []`.

- [ ] **Step 7: Remove member_ids from clearFilters**

In `clearFilters()` (lines 435-440), remove `filters.member_id = undefined`.

- [ ] **Step 8: Remove members from LedgerItem interface**

In the `LedgerItem` interface (lines 289-299), remove `members: Member[]`.

- [ ] **Step 9: Remove member display from template**

Remove lines 135-139 (the `item-members` span in the list):
```html
<span v-if="item.members && item.members.length > 0" class="item-members">
  关联：<span v-for="m in item.members.slice(0, 4)" :key="m.id" class="member-tag">{{ m.avatar }}</span>
  <span v-if="item.members.length > 4" class="member-more">等 {{ item.members.length }} 人</span>
</span>
<span v-else class="item-members"></span>
```

- [ ] **Step 10: Remove member select from modal form**

Remove lines 224-235 (the `a-form-item` for "关联成员"):
```html
<a-form-item label="关联成员" required>
  <a-select v-model:value="form.member_ids" mode="multiple" placeholder="选择成员" data-testid="member-select">
    <a-select-option v-for="m in members" :key="m.id" :value="m.id">
      {{ m.avatar }} {{ m.name }}
    </a-select-option>
  </a-select>
</a-form-item>
```

- [ ] **Step 11: Remove member-related CSS**

Remove the `.item-members`, `.member-tag`, `.member-more` CSS rules (lines 741-755).

- [ ] **Step 12: Remove month navigation template**

Replace the entire "Month Switcher" section (lines 9-20) with a new row containing the RangePicker and the "记一笔" button:
```html
<div class="month-row">
  <a-range-picker
    v-model:value="dateRange"
    format="YYYY-MM-DD"
    :presets="rangePresets"
    @change="onDateRangeChange"
    data-testid="date-range-picker"
  />
  <a-button type="primary" @click="openCreate()" data-testid="add-btn">记一笔</a-button>
</div>
```

- [ ] **Step 13: Add dateRange state and rangePresets**

Replace `selectedMonth` (line 327) with:
```ts
const dateRange = ref<[Dayjs, Dayjs]>([dayjs().startOf('month'), dayjs().endOf('month')])
```

Add `rangePresets` computed or constant:
```ts
const rangePresets = [
  { label: '本月', value: [dayjs().startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '上月', value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] as [Dayjs, Dayjs] },
  { label: '近三个月', value: [dayjs().subtract(2, 'month').startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '近半年', value: [dayjs().subtract(5, 'month').startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '今年', value: [dayjs().startOf('year'), dayjs().endOf('year')] as [Dayjs, Dayjs] },
]
```

- [ ] **Step 14: Replace month functions with date range handler**

Remove `getMonthParam()`, `goPrevMonth()`, `goNextMonth()` (lines 373-385).

Add:
```ts
function onDateRangeChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    dateRange.value = dates
    fetchLedgers()
  }
}
```

- [ ] **Step 15: Update fetchLedgers to use date range**

In `fetchLedgers()` (lines 387-424), change the params from:
```ts
const params: Record<string, unknown> = {
  month: getMonthParam(),
  page: page.value,
  page_size: pageSize,
}
if (filters.member_id) params.member_id = filters.member_id
```
to:
```ts
const params: Record<string, unknown> = {
  start_date: dateRange.value[0].format('YYYY-MM-DD'),
  end_date: dateRange.value[1].add(1, 'day').format('YYYY-MM-DD'),
  page: page.value,
  page_size: pageSize,
}
```

Note: `end_date` is exclusive (the backend uses `<`), so we add 1 day to the selected end date.

- [ ] **Step 16: Remove unused imports**

Remove `getMembers` import from `@/api/member` if it's no longer used. Check — the "创建者" filter still needs the `members` list, so `getMembers()` and the `members` ref should stay.

Remove the `Member` interface if still needed for the creator type — actually, `Member` is still used for `creator: Member` in `LedgerItem` and for the members list in the creator filter. Keep it.

- [ ] **Step 17: Update month-row CSS**

Remove `.month-switcher`, `.month-arrow`, `.month-arrow:disabled`, `.month-text` CSS rules (lines 572-595). The RangePicker from Ant Design handles its own styling.

- [ ] **Step 18: Add RangePicker stub to test stubs**

In the test file `frontend/src/views/ledger/__tests__/Index.test.ts`, add a stub for `a-range-picker`:
```ts
'a-range-picker': { template: '<div><slot /></div>' },
```

- [ ] **Step 19: Run all backend tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/... -v`
Expected: All tests pass.

- [ ] **Step 20: Run frontend tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Expected: View tests may fail — will fix in next task.

- [ ] **Step 21: Commit**

```bash
git add frontend/src/views/ledger/Index.vue
git commit -m "refactor(frontend): replace month nav with RangePicker, remove member UI from ledger"
```

---

### Task 11: Frontend tests — update view tests

**Files:**
- Modify: `frontend/src/views/ledger/__tests__/Index.test.ts`

- [ ] **Step 1: Add RangePicker stub**

In the `stubs` object, add:
```ts
'a-range-picker': { template: '<div><slot /></div>' },
```

- [ ] **Step 2: Remove members from mock ledger data**

In `mockLedgersData` (lines 124-139), remove the `members` field from all items:
```ts
members: [{ id: 1, name: '管理员', avatar: '👨' }],
```

Also remove `members` from the `updatedData` mock in the create and delete tests.

- [ ] **Step 3: Remove member_ids from create test**

In the `creates ledger and refreshes list` test (lines 196-249), remove:
```ts
vm.form.member_ids = [1]
```

- [ ] **Step 4: Update mockLedgersData references**

Ensure all mock data items no longer have `members` fields.

- [ ] **Step 5: Run view tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run src/views/ledger/__tests__/Index.test.ts`
Expected: All tests pass.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/views/ledger/__tests__/Index.test.ts
git commit -m "test(frontend): update ledger view tests for RangePicker and member removal"
```

---

### Task 12: Full verification

- [ ] **Step 1: Run all backend tests**

Run: `cd D:/Projects/my_projects/home-center-v1 && go test ./backend/... -v`
Expected: All tests pass.

- [ ] **Step 2: Run all frontend tests**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Expected: All tests pass.

- [ ] **Step 3: Verify build**

Run: `cd D:/Projects/my_projects/home-center-v1 && make build`
Expected: Build succeeds.

- [ ] **Step 4: Commit (if any fixups needed)**

If any fixups were needed, commit them:
```bash
git add -A
git commit -m "fix: address test/build issues from ledger refactor"
```
