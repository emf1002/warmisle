# Bug Fix & UI Issues Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix 6 reported bugs: double error messages, English modal buttons, ledger refresh after create, category filter disabled, category selector redesign, and date picker locale.

**Architecture:** All fixes are in the Vue 3 frontend. The main change is in `frontend/src/views/ledger/Index.vue` (category tabs redesign + bug fixes), with smaller changes across other view files and `App.vue` for global locale.

**Tech Stack:** Vue 3, Ant Design Vue 4.x, dayjs, TypeScript

---

## Root Cause Analysis

### Issue 1: Double Error Messages
The axios interceptor (`request.ts:18-35`) calls `message.error()` for ALL errors (HTTP 4xx/5xx at line 32, business errors at line 22). But page-level `catch` blocks also call `message.error(e.response.data.message)` for HTTP errors. Since the interceptor already showed the error, users see it twice.

**Affected files (catch blocks with `e?.response?.data?.message`):**
- `member/Index.vue:234-238` (handleSubmit)
- `category/Index.vue:209-212` (handleSubmit)
- `forum/Index.vue:458-459` (handlePostSubmit)
- `forum/Index.vue:526-527` (handleTopicSubmit)
- `forum/TopicDetail.vue` (handleEditSubmit, submitComment, submitReply)
- `wish/Index.vue:350-351` (handleSubmit)

### Issue 2: English Modal Buttons
6 `<a-modal>` components don't set `okText`/`cancelText`, falling back to Ant Design Vue's English defaults ("OK"/"Cancel").

### Issue 3: Ledger Not Refreshing After Create
The `fetchLedgers()` catch block at line 438 is empty. If the refresh request fails (e.g., network glitch), the user sees "记账成功" but the list stays stale. Additionally, the catch block in `handleSubmit` (line 526-529) is dead code for business errors (interceptor rejects with `new Error()` which has no `response` property).

### Issue 4: Category Filter Disabled
The `a-select-opt-group` with `a-select-option` inside may render as disabled when Ant Design Vue doesn't properly recognize nested options. The fix is to flatten the options with visual separators.

### Issue 5: Category Selector Redesign
Replace the `a-select` dropdown with `a-tabs` + clickable category cards (icon + text), matching the emoji picker pattern already used in `category/Index.vue`.

### Issue 6: Date Picker English Locale
No Chinese locale is configured. Fix globally via `<a-config-provider :locale="zhCN">` in `App.vue`.

---

## File Structure

| File | Change |
|------|--------|
| `frontend/src/App.vue` | Add `zhCN` locale to `<a-config-provider>` |
| `frontend/src/views/member/Index.vue` | Add `okText`/`cancelText` to modal, remove duplicate error in catch |
| `frontend/src/views/category/Index.vue` | Add `okText`/`cancelText` to modal, remove duplicate error in catch |
| `frontend/src/views/forum/Index.vue` | Add `okText`/`cancelText` to 2 modals, remove duplicate errors in catch blocks |
| `frontend/src/views/forum/TopicDetail.vue` | Add `okText`/`cancelText` to modal, remove duplicate errors in catch blocks |
| `frontend/src/views/wish/Index.vue` | Add `okText`/`cancelText` to modal, remove duplicate errors in catch blocks |
| `frontend/src/views/ledger/Index.vue` | Fix refresh error handling, flatten filter select, redesign category selector as tabs |
| `frontend/src/views/ledger/__tests__/Index.test.ts` | Add test for refresh error handling, update category selector test |

---

### Task 1: Global Chinese Locale for Date Picker

**Files:**
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Import zhCN locale and add to ConfigProvider**

In `App.vue`, import `zhCN` from `ant-design-vue/es/locale/zh_CN` and add it to the `<a-config-provider>`:

```vue
<template>
  <a-config-provider :theme="themeConfig" :locale="zhCN">
    <component :is="layout">
      <router-view v-slot="{ Component: RouteComponent }">
        <transition name="page-fade" mode="out-in">
          <component :is="RouteComponent" />
        </transition>
      </router-view>
    </component>
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { theme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import MainLayout from '@/layouts/MainLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useThemeStore } from '@/stores/theme'

dayjs.locale('zh-cn')

// ... rest unchanged
</script>
```

- [ ] **Step 2: Verify locale is applied**

Run the dev server and check that all date pickers show Chinese labels (年/月/日, 确定/取消 buttons in calendar popup).

- [ ] **Step 3: Commit**

```bash
git add frontend/src/App.vue
git commit -m "fix: add Chinese locale globally for ant-design-vue and dayjs"
```

---

### Task 2: Fix Double Error Messages (All Pages)

**Files:**
- Modify: `frontend/src/views/member/Index.vue:234-238`
- Modify: `frontend/src/views/category/Index.vue:209-212`
- Modify: `frontend/src/views/forum/Index.vue:458-459,526-527`
- Modify: `frontend/src/views/wish/Index.vue:350-351`

- [ ] **Step 1: Remove duplicate error handling in member/Index.vue**

In `handleSubmit()`, replace the catch block:

```typescript
// Before (lines 234-238):
  } catch (e: any) {
    // error handled by interceptor; show fallback if data has message
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

- [ ] **Step 2: Remove duplicate error handling in category/Index.vue**

In `handleSubmit()`, replace the catch block:

```typescript
// Before (lines 209-212):
  } catch (e: any) {
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

- [ ] **Step 3: Remove duplicate error handling in forum/Index.vue**

In `handlePostSubmit()` (line 458-459):
```typescript
// Before:
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

In `handleTopicSubmit()` (line 526-527):
```typescript
// Before:
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

Also fix `handleTogglePin()` (line 558-559) if it has the same pattern.

- [ ] **Step 4: Remove duplicate error handling in wish/Index.vue**

In `handleSubmit()` (line 350-351):
```typescript
// Before:
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

Also fix `handleStatusChange()` (line 381-382) if it has the same pattern.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/views/member/Index.vue frontend/src/views/category/Index.vue frontend/src/views/forum/Index.vue frontend/src/views/wish/Index.vue
git commit -m "fix: remove duplicate error messages in catch blocks across all pages"
```

---

### Task 3: Add Chinese Text to All Modal Buttons

**Files:**
- Modify: `frontend/src/views/member/Index.vue:65`
- Modify: `frontend/src/views/category/Index.vue:57`
- Modify: `frontend/src/views/forum/Index.vue:153,195`
- Modify: `frontend/src/views/forum/TopicDetail.vue:188`
- Modify: `frontend/src/views/wish/Index.vue:118`

- [ ] **Step 1: Add okText/cancelText to member/Index.vue modal**

```html
<!-- Before (line 65): -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingMember ? '编辑成员' : '添加成员'"
  @ok="handleSubmit"
  :confirm-loading="submitting"
  data-testid="member-modal"
>

<!-- After: -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingMember ? '编辑成员' : '添加成员'"
  ok-text="保存"
  cancel-text="取消"
  @ok="handleSubmit"
  :confirm-loading="submitting"
  data-testid="member-modal"
>
```

- [ ] **Step 2: Add okText/cancelText to category/Index.vue modal**

```html
<!-- Before (line 57): -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingCategory ? '编辑分类' : '添加分类'"
  @ok="handleSubmit"
  :confirm-loading="submitting"
  data-testid="category-modal"
>

<!-- After: -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingCategory ? '编辑分类' : '添加分类'"
  ok-text="保存"
  cancel-text="取消"
  @ok="handleSubmit"
  :confirm-loading="submitting"
  data-testid="category-modal"
>
```

- [ ] **Step 3: Add okText/cancelText to forum/Index.vue post modal**

```html
<!-- Before (line 153): -->
<a-modal
  v-model:open="postDialogOpen"
  :title="editingPostItem ? '编辑动态' : '发动态'"
  :confirm-loading="postSubmitting"
  @ok="handlePostSubmit"
  @cancel="postDialogOpen = false"
  data-testid="post-modal"
>

<!-- After: -->
<a-modal
  v-model:open="postDialogOpen"
  :title="editingPostItem ? '编辑动态' : '发动态'"
  ok-text="发布"
  cancel-text="取消"
  :confirm-loading="postSubmitting"
  @ok="handlePostSubmit"
  @cancel="postDialogOpen = false"
  data-testid="post-modal"
>
```

- [ ] **Step 4: Add okText/cancelText to forum/Index.vue topic modal**

```html
<!-- Before (line 195): -->
<a-modal
  v-model:open="topicDialogOpen"
  :title="editingTopicItem ? '编辑话题' : '发话题'"
  :confirm-loading="topicSubmitting"
  width="520px"
  @ok="handleTopicSubmit"
  @cancel="topicDialogOpen = false"
  data-testid="topic-modal"
>

<!-- After: -->
<a-modal
  v-model:open="topicDialogOpen"
  :title="editingTopicItem ? '编辑话题' : '发话题'"
  ok-text="发布"
  cancel-text="取消"
  :confirm-loading="topicSubmitting"
  width="520px"
  @ok="handleTopicSubmit"
  @cancel="topicDialogOpen = false"
  data-testid="topic-modal"
>
```

- [ ] **Step 5: Add okText/cancelText to forum/TopicDetail.vue edit modal**

```html
<!-- Before (line 188): -->
<a-modal
  v-model:open="editDialogOpen"
  title="编辑话题"
  :confirm-loading="editSubmitting"
  width="520px"
  @ok="handleEditSubmit"
  @cancel="editDialogOpen = false"
>

<!-- After: -->
<a-modal
  v-model:open="editDialogOpen"
  title="编辑话题"
  ok-text="保存"
  cancel-text="取消"
  :confirm-loading="editSubmitting"
  width="520px"
  @ok="handleEditSubmit"
  @cancel="editDialogOpen = false"
>
```

- [ ] **Step 6: Add okText/cancelText to wish/Index.vue modal**

```html
<!-- Before (line 118): -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingWish ? '编辑愿望' : '新建愿望'"
  :confirm-loading="submitting"
  width="480px"
  @ok="handleSubmit"
  @cancel="dialogOpen = false"
  data-testid="wish-modal"
>

<!-- After: -->
<a-modal
  v-model:open="dialogOpen"
  :title="editingWish ? '编辑愿望' : '新建愿望'"
  ok-text="保存"
  cancel-text="取消"
  :confirm-loading="submitting"
  width="480px"
  @ok="handleSubmit"
  @cancel="dialogOpen = false"
  data-testid="wish-modal"
>
```

- [ ] **Step 7: Commit**

```bash
git add frontend/src/views/member/Index.vue frontend/src/views/category/Index.vue frontend/src/views/forum/Index.vue frontend/src/views/forum/TopicDetail.vue frontend/src/views/wish/Index.vue
git commit -m "fix: add Chinese button text to all modal dialogs"
```

---

### Task 4: Fix Ledger Refresh Error Handling

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue:438,526-529`

- [ ] **Step 1: Add error feedback to fetchLedgers catch block**

```typescript
// Before (lines 438-439):
  } catch {
    // error handled by interceptor
  }

// After:
  } catch (e: any) {
    if (!isLoadMore) {
      groups.value = []
      summary.value = { income: 0, expense: 0, balance: 0 }
      total.value = 0
    }
  }
```

This ensures that if `fetchLedgers` fails after a successful create, the UI at least clears stale data rather than showing the old list (which might confuse the user into thinking the create didn't work).

- [ ] **Step 2: Clean up handleSubmit catch block**

```typescript
// Before (lines 526-529):
  } catch (e: any) {
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  }

// After:
  } catch {
    // error handled by interceptor
  }
```

- [ ] **Step 3: Add test for create-and-refresh flow**

In `frontend/src/views/ledger/__tests__/Index.test.ts`, update the existing create test to verify the list actually updates:

```typescript
it('creates ledger and refreshes list with updated data', async () => {
  mockCreateLedger.mockResolvedValue({ code: 0, message: 'ok', data: { id: 99 } })

  const updatedData = {
    code: 0, message: 'ok', data: {
      summary: { income: 10000, expense: 4550, balance: 5450 },
      groups: [{
        date: '2026-05-23', daily_total: -4550,
        items: [{
          id: 99, amount: 1000, note: '新记录', category_id: 1, creator_id: 1,
          occurred_at: '2026-05-23T14:00:00Z',
          category: { id: 1, name: '餐饮', icon: '🍱', type: 'expense' },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          members: [{ id: 1, name: '管理员', avatar: '👨' }],
        }, {
          id: 1, amount: 3550, note: '午餐', category_id: 1, creator_id: 1,
          occurred_at: '2026-05-23T12:00:00Z',
          category: { id: 1, name: '餐饮', icon: '🍱', type: 'expense' },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          members: [{ id: 1, name: '管理员', avatar: '👨' }],
        }],
      }],
      total: 2,
    },
  }

  mockGetLedgers.mockResolvedValueOnce(mockLedgersData)
  mockGetLedgers.mockResolvedValueOnce(updatedData)

  const wrapper = createWrapper()
  await flushPromises()
  await nextTick()

  const vm = wrapper.vm as any
  vm.dialogOpen = true
  vm.form.category_id = 1
  vm.form.amount = 10
  vm.form.member_ids = [1]
  vm.form.note = '新记录'
  vm.form.occurred_at = undefined
  await nextTick()

  await vm.handleSubmit()
  await flushPromises()
  await nextTick()

  expect(mockCreateLedger).toHaveBeenCalled()
  expect(mockGetLedgers).toHaveBeenCalledTimes(2)

  // Verify the new record appears in the rendered list
  const text = wrapper.text()
  expect(text).toContain('新记录')
  expect(text).toContain('10.00')
})
```

- [ ] **Step 4: Run tests**

```bash
cd frontend && npx vitest run src/views/ledger/__tests__/Index.test.ts
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/views/ledger/Index.vue frontend/src/views/ledger/__tests__/Index.test.ts
git commit -m "fix: improve ledger refresh error handling and add test coverage"
```

---

### Task 5: Fix Category Filter Dropdown (Disabled Options)

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue:53-70`

- [ ] **Step 1: Flatten the category filter select**

Replace the `a-select-opt-group` approach with flat options and a divider-style label:

```html
<!-- Before (lines 53-70): -->
<a-select
  v-model:value="filters.category_id"
  placeholder="全部分类"
  allow-clear
  style="width: 140px"
  @change="onFilterChange"
>
  <a-select-opt-group label="支出">
    <a-select-option v-for="c in expenseCategories" :key="c.id" :value="c.id">
      {{ c.icon }} {{ c.name }}
    </a-select-option>
  </a-select-opt-group>
  <a-select-opt-group label="收入">
    <a-select-option v-for="c in incomeCategories" :key="c.id" :value="c.id">
      {{ c.icon }} {{ c.name }}
    </a-select-option>
  </a-select-opt-group>
</a-select>

<!-- After: -->
<a-select
  v-model:value="filters.category_id"
  placeholder="全部分类"
  allow-clear
  style="width: 140px"
  @change="onFilterChange"
>
  <template v-if="expenseCategories.length > 0">
    <a-select-option disabled :value="-1" class="category-group-label">支出</a-select-option>
    <a-select-option v-for="c in expenseCategories" :key="c.id" :value="c.id">
      {{ c.icon }} {{ c.name }}
    </a-select-option>
  </template>
  <template v-if="incomeCategories.length > 0">
    <a-select-option disabled :value="-2" class="category-group-label">收入</a-select-option>
    <a-select-option v-for="c in incomeCategories" :key="c.id" :value="c.id">
      {{ c.icon }} {{ c.name }}
    </a-select-option>
  </template>
</a-select>
```

Add CSS for the group label:
```css
.category-group-label {
  font-weight: 600;
  color: var(--color-text-secondary);
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/ledger/Index.vue
git commit -m "fix: flatten category filter select to fix disabled options issue"
```

---

### Task 6: Redesign Category Selector as Tabs

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue:178-198`

- [ ] **Step 1: Add tab state to script**

Add a reactive variable to track the active tab (expense/income):

```typescript
// Add after form reactive (around line 333):
const categoryTab = ref<'expense' | 'income'>('expense')
```

Update `openCreate()` to reset the tab:
```typescript
function openCreate() {
  editingRecord.value = null
  form.category_id = undefined
  form.amount = null
  form.member_ids = []
  form.occurred_at = dayjs()
  form.note = ''
  categoryTab.value = 'expense'  // Add this line
  dialogOpen.value = true
}
```

Update `openEdit()` to set the tab based on the selected category:
```typescript
function openEdit(record: LedgerItem) {
  editingRecord.value = record
  form.category_id = record.category_id
  form.amount = record.amount / 100
  form.member_ids = record.members ? record.members.map((m) => m.id) : []
  form.occurred_at = dayjs(record.occurred_at)
  form.note = record.note || ''
  // Set tab based on category type
  const cat = categories.value.find(c => c.id === record.category_id)
  categoryTab.value = cat?.type === 'income' ? 'income' : 'expense'
  dialogOpen.value = true
}
```

- [ ] **Step 2: Replace select with tabs and category cards**

Replace lines 178-198:

```html
<!-- Before: -->
<a-form-item label="分类" required>
  <a-select v-model:value="form.category_id" placeholder="选择分类" data-testid="category-select">
    <a-select-opt-group label="支出分类">
      <a-select-option
        v-for="cat in expenseCategories"
        :key="cat.id"
        :value="cat.id"
      >
        {{ cat.icon }} {{ cat.name }}
      </a-select-option>
    </a-select-opt-group>
    <a-select-opt-group label="收入分类">
      <a-select-option
        v-for="cat in incomeCategories"
        :key="cat.id"
        :value="cat.id"
      >
        {{ cat.icon }} {{ cat.name }}
      </a-select-option>
    </a-select-opt-group>
  </a-select>
</a-form-item>

<!-- After: -->
<a-form-item label="分类" required>
  <a-tabs v-model:activeKey="categoryTab" size="small">
    <a-tab-pane key="expense" tab="支出">
      <div class="category-grid-picker">
        <div
          v-for="cat in expenseCategories"
          :key="cat.id"
          :class="['category-pick-item', { active: form.category_id === cat.id }]"
          @click="form.category_id = cat.id"
          :data-testid="'cat-expense-' + cat.id"
        >
          <span class="category-pick-icon">{{ cat.icon }}</span>
          <span class="category-pick-name">{{ cat.name }}</span>
        </div>
      </div>
    </a-tab-pane>
    <a-tab-pane key="income" tab="收入">
      <div class="category-grid-picker">
        <div
          v-for="cat in incomeCategories"
          :key="cat.id"
          :class="['category-pick-item', { active: form.category_id === cat.id }]"
          @click="form.category_id = cat.id"
          :data-testid="'cat-income-' + cat.id"
        >
          <span class="category-pick-icon">{{ cat.icon }}</span>
          <span class="category-pick-name">{{ cat.name }}</span>
        </div>
      </div>
    </a-tab-pane>
  </a-tabs>
</a-form-item>
```

- [ ] **Step 3: Add CSS for the category picker**

Add to the `<style scoped>` section:

```css
/* Category Grid Picker (in modal) */
.category-grid-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.category-pick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 64px;
  background: var(--color-bg-container);
}

.category-pick-item:hover {
  border-color: var(--color-border);
  background: var(--color-bg-layout);
}

.category-pick-item.active {
  border-color: var(--color-brand);
  background: var(--color-brand-light);
}

.category-pick-icon {
  font-size: 24px;
  line-height: 1;
}

.category-pick-name {
  font-size: 12px;
  color: var(--color-text-primary);
  white-space: nowrap;
}
```

- [ ] **Step 4: Add stubs for a-tabs and a-tab-pane in tests**

In `frontend/src/views/ledger/__tests__/Index.test.ts`, add to the stubs:

```typescript
'a-tabs': {
  props: ['activeKey'],
  template: '<div><slot /></div>',
  emits: ['update:activeKey'],
},
'a-tab-pane': {
  props: ['key', 'tab'],
  template: '<div><slot /></div>',
},
```

- [ ] **Step 5: Run tests**

```bash
cd frontend && npx vitest run src/views/ledger/__tests__/Index.test.ts
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/views/ledger/Index.vue frontend/src/views/ledger/__tests__/Index.test.ts
git commit -m "feat: replace category dropdown with tab-based icon picker in ledger form"
```

---

### Task 7: Run All Frontend Tests and Final Verification

- [ ] **Step 1: Run all frontend tests**

```bash
cd frontend && npx vitest run
```

- [ ] **Step 2: Fix any failing tests**

If tests fail due to the modal stub not rendering `okText`/`cancelText`, update the `a-modal` stub in test files to render the button text.

- [ ] **Step 3: Build check**

```bash
cd frontend && npx vue-tsc --noEmit
```

- [ ] **Step 4: Commit any fixes**

```bash
git add -A
git commit -m "fix: update tests for modal and category selector changes"
```
