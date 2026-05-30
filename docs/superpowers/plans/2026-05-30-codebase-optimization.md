# 代码库优化实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 消除前后端代码重复，统一错误处理模式，拆分过大模块，清理基础设施

**Architecture:** 分 4 阶段实施：前端代码提取 → 后端错误处理统一 → Forum 模块拆分 → 基础设施清理。每阶段可独立验证，每任务产出可独立提交。

**Tech Stack:** Vue 3 Composition API, Pinia, Go, Gin, GORM, SQLite

---

## 阶段 1：前端代码重复提取（Task 1-4）

### Task 1: JWT 解析提取到 Auth Store

**Files:**
- Modify: `frontend/src/stores/auth.ts`
- Modify: `frontend/src/views/ledger/Index.vue`
- Modify: `frontend/src/views/todo/Index.vue`
- Modify: `frontend/src/views/wish/Index.vue`
- Modify: `frontend/src/views/forum/Index.vue`
- Modify: `frontend/src/views/forum/TopicDetail.vue`
- Modify: `frontend/src/views/member/Index.vue`
- Modify: `frontend/src/views/category/Index.vue`

- [ ] **Step 1: 扩展 auth store，添加 JWT 解析计算属性**

```ts
// frontend/src/stores/auth.ts — 在现有 token ref 之后添加
import { computed } from 'vue'

// 在 defineStore 内部，token 声明之后
const currentUserId = computed(() => {
  if (!token.value) return null
  try {
    const payload = JSON.parse(atob(token.value.split('.')[1]))
    return payload.member_id
  } catch {
    return null
  }
})

const currentUserRole = computed(() => {
  if (!token.value) return null
  try {
    const payload = JSON.parse(atob(token.value.split('.')[1]))
    return payload.role
  } catch {
    return null
  }
})

const isAdmin = computed(() => currentUserRole.value === 'admin')
```

确保在 `return` 中导出 `currentUserId`, `currentUserRole`, `isAdmin`。

- [ ] **Step 2: 修改 ledger/Index.vue，删除本地 JWT 解析**

删除文件中的 `currentUserRole` 和 `currentUserId` 函数/计算属性（约 lines 359-379），改为：

```ts
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
// 使用 authStore.currentUserId 和 authStore.currentUserRole 替换原有变量
```

搜索文件中所有 `currentUserRole()` 和 `currentUserId()` 调用，替换为 `authStore.currentUserRole` 和 `authStore.currentUserId`。

- [ ] **Step 3: 修改 todo/Index.vue**

删除 `getCurrentUserId` 和 `getCurrentUserRole` 函数（约 lines 205-221），改为使用 auth store。

- [ ] **Step 4: 修改 wish/Index.vue**

删除 `currentUserId` 和 `isAdmin` 计算属性（约 lines 203-219），改为使用 auth store。

- [ ] **Step 5: 修改 forum/Index.vue**

删除 `currentUserId` 和 `isAdmin` 计算属性（约 lines 326-346），改为使用 auth store。

- [ ] **Step 6: 修改 forum/TopicDetail.vue**

删除 `currentUserId` 和 `isAdmin` 计算属性（约 lines 318-338），改为使用 auth store。

- [ ] **Step 7: 修改 member/Index.vue**

删除 `isAdmin` 计算属性（约 lines 152-161），改为使用 auth store。

- [ ] **Step 8: 修改 category/Index.vue**

删除 `isAdmin` 计算属性（约 lines 134-143），改为使用 auth store。

- [ ] **Step 9: 运行前端测试**

```bash
cd frontend && npm test
```

Expected: 所有测试通过

- [ ] **Step 10: 提交**

```bash
git add frontend/src/stores/auth.ts frontend/src/views/
git commit -m "refactor: extract JWT parsing to auth store, remove duplication in 7 views"
```

---

### Task 2: 工具函数提取到 utils/format.ts

**Files:**
- Create: `frontend/src/utils/format.ts`
- Modify: `frontend/src/views/ledger/Index.vue`
- Modify: `frontend/src/views/todo/Index.vue`
- Modify: `frontend/src/views/wish/Index.vue`
- Modify: `frontend/src/views/forum/Index.vue`
- Modify: `frontend/src/views/forum/TopicDetail.vue`
- Modify: `frontend/src/views/dashboard/Index.vue`

- [ ] **Step 1: 创建 utils/format.ts**

```ts
// frontend/src/utils/format.ts
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

/**
 * 截断字符串
 */
export function truncate(str: string, len: number): string {
  if (!str) return ''
  return str.length > len ? str.slice(0, len) + '...' : str
}

/**
 * 相对时间（如"3分钟前"、"刚刚"）
 */
export function timeAgo(date: string | Date): string {
  if (!date) return ''
  const d = dayjs(date)
  const now = dayjs()
  const diffMinutes = now.diff(d, 'minute')
  if (diffMinutes < 1) return '刚刚'
  return d.fromNow()
}

/**
 * 格式化日期（今天/昨天/M月D日 周X）
 */
export function formatDate(date: string | Date): string {
  if (!date) return ''
  const d = dayjs(date)
  const now = dayjs()
  if (d.isSame(now, 'day')) return '今天'
  if (d.isSame(now.subtract(1, 'day'), 'day')) return '昨天'
  return d.format('M月D日') + ' 周' + ['日', '一', '二', '三', '四', '五', '六'][d.day()]
}

/**
 * 优先级颜色
 */
export function priorityColor(priority: string): string {
  const map: Record<string, string> = {
    normal: '#52c41a',
    important: '#faad14',
    urgent: '#ff4d4f',
  }
  return map[priority] || '#999'
}

/**
 * 优先级中文标签
 */
export function priorityLabel(priority: string): string {
  const map: Record<string, string> = {
    normal: '普通',
    important: '重要',
    urgent: '紧急',
  }
  return map[priority] || priority
}
```

- [ ] **Step 2: 修改 dashboard/Index.vue**

删除本地 `timeAgo` 和 `formatDate` 函数（约 lines 183-195），改为：

```ts
import { timeAgo, formatDate } from '@/utils/format'
```

- [ ] **Step 3: 修改 ledger/Index.vue**

删除本地 `truncate` 和 `formatDate` 函数（约 lines 386-400），改为：

```ts
import { truncate, formatDate } from '@/utils/format'
```

- [ ] **Step 4: 修改 todo/Index.vue**

删除本地 `priorityColor`、`priorityLabel`、`formatDate` 函数（约 lines 230-250），改为：

```ts
import { priorityColor, priorityLabel, formatDate } from '@/utils/format'
```

- [ ] **Step 5: 修改 wish/Index.vue**

删除本地 `priorityColor`、`priorityLabel` 函数（约 lines 231-241），改为：

```ts
import { priorityColor, priorityLabel } from '@/utils/format'
```

- [ ] **Step 6: 修改 forum/Index.vue**

删除本地 `truncate` 和 `timeAgo` 函数（约 lines 348-360），改为：

```ts
import { truncate, timeAgo } from '@/utils/format'
```

- [ ] **Step 7: 修改 forum/TopicDetail.vue**

删除本地 `timeAgo` 函数（约 lines 351-358），改为：

```ts
import { timeAgo } from '@/utils/format'
```

- [ ] **Step 8: 运行前端测试**

```bash
cd frontend && npm test
```

Expected: 所有测试通过

- [ ] **Step 9: 提交**

```bash
git add frontend/src/utils/format.ts frontend/src/views/
git commit -m "refactor: extract shared utility functions to utils/format.ts"
```

---

### Task 3: CSS 公共类提取到 components.css

**Files:**
- Create: `frontend/src/styles/components.css`
- Modify: `frontend/src/main.ts`
- Modify: `frontend/src/styles/global.css`
- Modify: `frontend/src/views/ledger/Index.vue`
- Modify: `frontend/src/views/todo/Index.vue`
- Modify: `frontend/src/views/wish/Index.vue`
- Modify: `frontend/src/views/forum/Index.vue`
- Modify: `frontend/src/views/forum/TopicDetail.vue`
- Modify: `frontend/src/views/member/Index.vue`
- Modify: `frontend/src/views/category/Index.vue`
- Modify: `frontend/src/views/profile/Index.vue`

- [ ] **Step 1: 创建 styles/components.css**

从各视图中提取重复的 CSS 到公共文件。以下为合并后的样式：

```css
/* frontend/src/styles/components.css */

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

/* 加载状态 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
  color: var(--text-secondary);
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--text-secondary);
}

/* 筛选行 */
.filter-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}

/* Emoji 选择器 */
.emoji-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.emoji-picker .emoji-item {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  border-radius: var(--radius-sm);
  border: 2px solid transparent;
  transition: all 0.2s;
}

.emoji-picker .emoji-item:hover {
  background: var(--bg-hover);
}

.emoji-picker .emoji-item.active {
  border-color: var(--brand-primary);
  background: var(--bg-active);
}

/* 分页行 */
.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
```

- [ ] **Step 2: 在 main.ts 中导入 components.css**

在 `main.ts` 中添加：

```ts
import './styles/components.css'
```

- [ ] **Step 3: 清理 global.css 中未使用的类**

删除 `global.css` 中的 `.income-color` 和 `.expense-color` 定义（如果存在且未被使用）。

- [ ] **Step 4: 从各视图中删除重复的 scoped CSS**

从以下视图中删除对应的 scoped CSS 块：
- `ledger/Index.vue` — 删除 `.loading-state`, `.empty-state`, `.filter-row`
- `todo/Index.vue` — 删除 `.page-header`, `.loading-state`, `.filter-row`, `.pagination-row`
- `wish/Index.vue` — 删除 `.page-header`, `.loading-state`, `.pagination-row`
- `forum/Index.vue` — 删除 `.page-header`, `.loading-state`, `.pagination-row`
- `forum/TopicDetail.vue` — 删除 `.loading-state`
- `member/Index.vue` — 删除 `.page-header`, `.emoji-picker`, `.emoji-item`
- `category/Index.vue` — 删除 `.page-header`, `.emoji-picker`, `.emoji-item`
- `profile/Index.vue` — 删除 `.page-header`, `.emoji-picker`, `.emoji-item`

**注意**：如果某个视图的样式与公共版本有差异，在该视图中保留差异部分的 scoped CSS。

- [ ] **Step 5: 运行前端测试**

```bash
cd frontend && npm test
```

Expected: 所有测试通过

- [ ] **Step 6: 运行前端构建验证 CSS 无冲突**

```bash
cd frontend && npm run build
```

Expected: 构建成功，无 CSS 错误

- [ ] **Step 7: 提交**

```bash
git add frontend/src/styles/components.css frontend/src/main.ts frontend/src/styles/global.css frontend/src/views/
git commit -m "refactor: extract shared CSS classes to styles/components.css"
```

---

### Task 4: EmojiPicker 组件提取

**Files:**
- Create: `frontend/src/components/EmojiPicker.vue`
- Modify: `frontend/src/views/member/Index.vue`
- Modify: `frontend/src/views/category/Index.vue`
- Modify: `frontend/src/views/profile/Index.vue`

- [ ] **Step 1: 创建 EmojiPicker.vue 组件**

```vue
<!-- frontend/src/components/EmojiPicker.vue -->
<template>
  <div class="emoji-picker">
    <div
      v-for="emoji in emojiList"
      :key="emoji"
      class="emoji-item"
      :class="{ active: modelValue === emoji }"
      @click="$emit('update:modelValue', emoji)"
    >
      {{ emoji }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const emojiList = [
  '😀', '😂', '🥰', '😎', '🤔', '😴', '🥳', '😇',
  '🤗', '😋', '🤩', '😏', '😤', '🤯', '🫡', '🤭',
  '😌', '🤓', '🐶', '🐱', '🐼', '🦊', '🐸', '🐵',
  '🏠', '🌟', '🔥', '❤️', '💪', '🎉', '🌈', '☕'
]
</script>
```

- [ ] **Step 2: 修改 member/Index.vue**

删除本地 `emojiList` 数组和 emoji picker 模板，替换为：

```vue
<EmojiPicker v-model="form.avatar" />
```

添加导入：

```ts
import EmojiPicker from '@/components/EmojiPicker.vue'
```

- [ ] **Step 3: 修改 category/Index.vue**

同上，删除本地 emoji picker，使用组件。

- [ ] **Step 4: 修改 profile/Index.vue**

同上，删除本地 emoji picker，使用组件。

- [ ] **Step 5: 运行前端测试**

```bash
cd frontend && npm test
```

Expected: 所有测试通过

- [ ] **Step 6: 提交**

```bash
git add frontend/src/components/EmojiPicker.vue frontend/src/views/member/ frontend/src/views/category/ frontend/src/views/profile/
git commit -m "refactor: extract EmojiPicker component from 3 views"
```

---

## 阶段 2：后端错误处理统一（Task 5-8）

### Task 5: Handler 错误映射 Helper

**Files:**
- Create: `backend/internal/handler/errors.go`
- Modify: `backend/internal/handler/ledger.go`
- Modify: `backend/internal/handler/todo.go`
- Modify: `backend/internal/handler/wish.go`
- Modify: `backend/internal/handler/forum.go`
- Modify: `backend/internal/handler/member.go`
- Modify: `backend/internal/handler/category.go`
- Modify: `backend/internal/handler/tag.go`
- Modify: `backend/internal/handler/init.go`
- Modify: `backend/internal/handler/dashboard.go`

- [ ] **Step 1: 创建 handler/errors.go**

```go
// backend/internal/handler/errors.go
package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"warmisle/internal/pkg"
)

// serviceError maps a service error to an HTTP response
type serviceError struct {
	err      error
	httpCode int
	bizCode  int
	msg      string
}

// handleServiceError checks err against mappings and returns the first match.
// If no match, returns 500.
func handleServiceError(c *gin.Context, err error, mappings ...serviceError) {
	for _, m := range mappings {
		if errors.Is(err, m.err) {
			pkg.Error(c, m.httpCode, m.bizCode, m.msg)
			return
		}
	}
	pkg.Error(c, 500, 50001, "服务器内部错误")
}

// getMemberID extracts member_id from Gin context (set by JWT middleware)
func getMemberID(c *gin.Context) uint {
	if id, exists := c.Get("member_id"); exists {
		return id.(uint)
	}
	return 0
}

// getRole extracts role from Gin context (set by JWT middleware)
func getRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}
```

- [ ] **Step 2: 从 ledger.go 中删除 getMemberID/getRole**

删除 `backend/internal/handler/ledger.go` 中的 `getMemberID()` 和 `getRole()` 函数定义（约 lines 27-35）。它们现在在 `errors.go` 中定义。

- [ ] **Step 3: 修改 ledger.go 中的错误处理**

将所有手动错误检查模式改为使用 `handleServiceError`。示例：

```go
// 之前
if err != nil {
    if errors.Is(err, service.ErrLedgerNotFound) {
        pkg.Error(c, 404, 40001, "记录不存在")
        return
    }
    if errors.Is(err, service.ErrPermissionDenied) {
        pkg.Error(c, 403, 40301, "只能修改自己创建的记录")
        return
    }
    pkg.Error(c, 500, 50001, "服务器内部错误")
    return
}

// 之后
if err != nil {
    handleServiceError(c, err,
        serviceError{service.ErrLedgerNotFound, 404, 40001, "记录不存在"},
        serviceError{service.ErrPermissionDenied, 403, 40301, "只能修改自己创建的记录"},
    )
    return
}
```

对 `ledger.go` 中的每个 handler 方法重复此操作。

- [ ] **Step 4: 修改 todo.go 中的错误处理**

同 Step 3 模式，处理 `todo.go` 中的所有错误检查。

- [ ] **Step 5: 修改 wish.go 中的错误处理**

同 Step 3 模式。

- [ ] **Step 6: 修改 forum.go 中的错误处理**

同 Step 3 模式。forum.go 有最多的错误检查点。

- [ ] **Step 7: 修改 member.go 中的错误处理**

同 Step 3 模式。

- [ ] **Step 8: 修改 category.go 中的错误处理**

同 Step 3 模式。

- [ ] **Step 9: 修改 tag.go 中的错误处理**

同 Step 3 模式。

- [ ] **Step 10: 检查 init.go 和 dashboard.go**

这两个文件可能没有复杂的错误处理，检查并统一。

- [ ] **Step 11: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 12: 提交**

```bash
git add backend/internal/handler/
git commit -m "refactor: add handleServiceError helper, remove error handling boilerplate in handlers"
```

---

### Task 6: Service 层 wrapNotFound Helper

**Files:**
- Create: `backend/internal/service/errors.go`
- Modify: `backend/internal/service/ledger.go`
- Modify: `backend/internal/service/todo.go`
- Modify: `backend/internal/service/wish.go`
- Modify: `backend/internal/service/forum.go`
- Modify: `backend/internal/service/member.go`
- Modify: `backend/internal/service/category.go`

- [ ] **Step 1: 创建 service/errors.go**

```go
// backend/internal/service/errors.go
package service

import (
	"errors"

	"gorm.io/gorm"
)

// wrapNotFound converts gorm.ErrRecordNotFound to the given domain error.
// Returns the original error if it's not a "not found" error.
func wrapNotFound(err error, domainErr error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainErr
	}
	return err
}
```

- [ ] **Step 2: 修改 ledger.go 中的 not-found 检查**

```go
// 之前
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrLedgerNotFound
}

// 之后
if err != nil {
    return nil, wrapNotFound(err, ErrLedgerNotFound)
}
```

对 `ledger.go` 中的所有 not-found 检查重复此操作。

- [ ] **Step 3: 修改 todo.go 中的 not-found 检查**

同 Step 2 模式。

- [ ] **Step 4: 修改 wish.go 中的 not-found 检查**

同 Step 2 模式。

- [ ] **Step 5: 修改 forum.go 中的 not-found 检查**

同 Step 2 模式。

- [ ] **Step 6: 修改 member.go 中的 not-found 检查**

同 Step 2 模式。

- [ ] **Step 7: 修改 category.go 中的 not-found 检查**

同 Step 2 模式。

- [ ] **Step 8: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 9: 提交**

```bash
git add backend/internal/service/
git commit -m "refactor: add wrapNotFound helper, remove gorm.ErrRecordNotFound boilerplate in services"
```

---

### Task 7: 错误码规范化

**Files:**
- Modify: `backend/internal/handler/errors.go` (update comments/docs)
- Modify: `backend/internal/handler/` (all handler files)
- Modify: `backend/internal/handler/test_helpers_test.go` (if assertions check biz codes)

- [ ] **Step 1: 统一错误码使用**

将所有资源不存在的 bizCode 从 `40001` 改为 `40401`：

| 位置 | 原 bizCode | 新 bizCode |
|------|-----------|-----------|
| handler/ledger.go | 40001 | 40401 |
| handler/todo.go | 40001 | 40401 |
| handler/wish.go | 40001 | 40401 |
| handler/forum.go | 40001 | 40401 |
| handler/member.go | 40001 | 40401 |
| handler/category.go | 40001 | 40401 |
| handler/tag.go | 40001 | 40401 |

**注意**：仅修改"资源不存在"的错误码，参数校验错误保持 `40001`。

- [ ] **Step 2: 更新测试中的 bizCode 断言**

检查所有 `*_test.go` 文件中对 `40001` 的断言，如果是资源不存在的场景，改为 `40401`。

- [ ] **Step 3: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 4: 提交**

```bash
git add backend/internal/handler/
git commit -m "refactor: standardize error codes - 40401 for not-found, 40001 for validation"
```

---

### Task 8: 辅助清理

**Files:**
- Modify: `backend/internal/handler/test_reset.go`
- Test: `backend/internal/handler/init_test.go` (verify test_reset endpoint still works)

- [ ] **Step 1: 修改 test_reset.go 使用 pkg.Success/pkg.Error**

将 `test_reset.go` 中的 `c.JSON(http.StatusOK, gin.H{...})` 替换为 `pkg.Success(c, nil)` 和 `pkg.Error(c, ...)`。

- [ ] **Step 2: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 3: 提交**

```bash
git add backend/internal/handler/test_reset.go
git commit -m "refactor: use pkg.Success/Error in test_reset.go"
```

---

## 阶段 3：Forum 模块拆分（Task 9-11）

### Task 9: Forum Handler 层拆分

**Files:**
- Create: `backend/internal/handler/forum_post.go`
- Create: `backend/internal/handler/forum_topic.go`
- Create: `backend/internal/handler/forum_comment.go`
- Create: `backend/internal/handler/forum_vote.go`
- Create: `backend/internal/handler/forum_feed.go`
- Modify: `backend/internal/handler/forum.go` (删除已拆分的函数，仅保留共享逻辑)

- [ ] **Step 1: 创建 forum_feed.go**

从 `forum.go` 中提取 feed 相关的 handler 函数到 `forum_feed.go`。

```go
// backend/internal/handler/forum_feed.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"warmisle/internal/pkg"
)

// GetFeed handles GET /api/forum/feed
func GetFeed(c *gin.Context) {
	// ... 从 forum.go 移动过来的代码
}
```

- [ ] **Step 2: 创建 forum_post.go**

从 `forum.go` 中提取 post CRUD 函数。

```go
// backend/internal/handler/forum_post.go
package handler

// CreatePost, UpdatePost, DeletePost, GetPost
```

- [ ] **Step 3: 创建 forum_topic.go**

从 `forum.go` 中提取 topic CRUD + pin 函数。

```go
// backend/internal/handler/forum_topic.go
package handler

// CreateTopic, UpdateTopic, DeleteTopic, GetTopic, PinTopic
```

- [ ] **Step 4: 创建 forum_comment.go**

从 `forum.go` 中提取 comment CRUD 函数。

```go
// backend/internal/handler/forum_comment.go
package handler

// CreateComment, UpdateComment, DeleteComment
```

- [ ] **Step 5: 创建 forum_vote.go**

从 `forum.go` 中提取 vote + like 函数。

```go
// backend/internal/handler/forum_vote.go
package handler

// VotePost, LikePost
```

- [ ] **Step 6: 清理原 forum.go**

删除已移动到新文件中的函数。保留任何共享的 helper 函数或类型定义（如果有）。

- [ ] **Step 7: 运行后端测试**

```bash
cd backend && go test ./internal/handler/...
```

Expected: 所有测试通过

- [ ] **Step 8: 提交**

```bash
git add backend/internal/handler/forum*.go
git commit -m "refactor: split forum handler into 5 sub-resource files"
```

---

### Task 10: Forum Service 层拆分

**Files:**
- Create: `backend/internal/service/forum_post.go`
- Create: `backend/internal/service/forum_topic.go`
- Create: `backend/internal/service/forum_comment.go`
- Create: `backend/internal/service/forum_vote.go`
- Modify: `backend/internal/service/forum.go` (删除已拆分的函数)

- [ ] **Step 1: 创建 forum_post.go**

从 `service/forum.go` 中提取 post 业务逻辑。

- [ ] **Step 2: 创建 forum_topic.go**

从 `service/forum.go` 中提取 topic 业务逻辑。

- [ ] **Step 3: 创建 forum_comment.go**

从 `service/forum.go` 中提取 comment 业务逻辑。

- [ ] **Step 4: 创建 forum_vote.go**

从 `service/forum.go` 中提取 vote + like 业务逻辑。

- [ ] **Step 5: 清理原 forum.go**

保留共享的错误变量定义（如 `ErrPostNotFound` 等）和 helper 函数。

- [ ] **Step 6: 运行后端测试**

```bash
cd backend && go test ./internal/service/...
```

Expected: 所有测试通过

- [ ] **Step 7: 提交**

```bash
git add backend/internal/service/forum*.go
git commit -m "refactor: split forum service into 4 sub-resource files"
```

---

### Task 11: Forum Repository 层拆分 + FeedItem 统一

**Files:**
- Create: `backend/internal/repository/forum_post.go`
- Create: `backend/internal/repository/forum_topic.go`
- Create: `backend/internal/repository/forum_comment.go`
- Create: `backend/internal/repository/forum_vote.go`
- Modify: `backend/internal/repository/forum.go` (保留共享类型定义)
- Modify: `backend/internal/service/dashboard.go` (引用 repo.FeedItem)

- [ ] **Step 1: 创建 forum_post.go**

从 `repository/forum.go` 中提取 post 数据访问函数。

- [ ] **Step 2: 创建 forum_topic.go**

从 `repository/forum.go` 中提取 topic 数据访问函数。

- [ ] **Step 3: 创建 forum_comment.go**

从 `repository/forum.go` 中提取 comment 数据访问函数。

- [ ] **Step 4: 创建 forum_vote.go**

从 `repository/forum.go` 中提取 vote + like 数据访问函数。

- [ ] **Step 5: 统一 FeedItem 定义**

保留 `repository/forum.go` 中的 `FeedItem` 结构体定义。删除 `service/dashboard.go` 中的重复定义，改为引用 `repository.FeedItem`。

```go
// service/dashboard.go
// 删除本地 FeedItem 定义，改为：
import "warmisle/internal/repository"

// 使用 repository.FeedItem
```

- [ ] **Step 6: 清理原 repository/forum.go**

保留 `FeedItem` 结构体和共享的 helper 函数。

- [ ] **Step 7: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 8: 提交**

```bash
git add backend/internal/repository/forum*.go backend/internal/service/dashboard.go
git commit -m "refactor: split forum repository, unify FeedItem struct"
```

---

## 阶段 4：基础设施清理（Task 12-17）

### Task 12: Makefile 添加 test/lint/fmt

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: 添加 test、lint、fmt target**

在 `Makefile` 中添加：

```makefile
.PHONY: test lint fmt

test:
	cd backend && go test ./...
	cd frontend && npm test

lint:
	cd backend && golangci-lint run ./...
	cd frontend && npx eslint src/

fmt:
	cd backend && gofmt -w .
	cd frontend && npx eslint --fix src/
```

- [ ] **Step 2: 验证 test target**

```bash
make test
```

Expected: 前后端测试都运行

- [ ] **Step 3: 提交**

```bash
git add Makefile
git commit -m "build: add test, lint, fmt targets to Makefile"
```

---

### Task 13: 迁移 Down 块

**Files:**
- Modify: `backend/migrations/002_fix_time_format.up.sql`

- [ ] **Step 1: 添加 +goose Down 块**

在文件末尾添加：

```sql
-- +goose Down
-- 数据迁移无法回滚，此迁移为单向操作
```

- [ ] **Step 2: 提交**

```bash
git add backend/migrations/002_fix_time_format.up.sql
git commit -m "fix: add empty +goose Down block to migration 002"
```

---

### Task 14: 死代码清理

**Files:**
- Modify: `backend/internal/pkg/response.go`

- [ ] **Step 1: 删除未使用的 PageData 结构体**

从 `response.go` 中删除 `PageData` 结构体定义（约 lines 11-16）。

- [ ] **Step 2: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过（确认 PageData 确实未被使用）

- [ ] **Step 3: 提交**

```bash
git add backend/internal/pkg/response.go
git commit -m "chore: remove unused PageData struct from response.go"
```

---

### Task 15: 路由测试统一

**Files:**
- Modify: `backend/internal/handler/test_helpers_test.go`

- [ ] **Step 1: 修改 setupTestRouter 使用 routes.Register**

将 `test_helpers_test.go` 中的重复路由注册改为调用 `routes.Register()`。

```go
// 之前：手动注册所有路由
// 之后
import "warmisle/internal/routes"

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    routes.Register(r)
    return r
}
```

- [ ] **Step 2: 运行后端测试**

```bash
cd backend && go test ./internal/handler/...
```

Expected: 所有测试通过

- [ ] **Step 3: 提交**

```bash
git add backend/internal/handler/test_helpers_test.go
git commit -m "refactor: use routes.Register() in test helpers instead of duplicate registration"
```

---

### Task 16: GORM 日志级别 + 合并重复函数

**Files:**
- Modify: `backend/internal/pkg/database.go`
- Create: `backend/internal/service/common.go`
- Modify: `backend/internal/service/todo.go`
- Modify: `backend/internal/service/wish.go`

- [ ] **Step 1: 修改 GORM 日志级别**

将 `database.go` 中的 `logger.Info` 改为 `logger.Warn`。

```go
// 之前
logger.Default.LogMode(logger.Info)

// 之后
logger.Default.LogMode(logger.Warn)
```

- [ ] **Step 2: 创建 service/common.go**

```go
// backend/internal/service/common.go
package service

// validPriority checks if the priority value is valid
func validPriority(p string) bool {
	return p == "normal" || p == "important" || p == "urgent"
}
```

- [ ] **Step 3: 修改 todo.go**

删除 `validPriority` 函数，使用 `common.go` 中的版本。

- [ ] **Step 4: 修改 wish.go**

删除 `validWishPriority` 函数，改为调用 `validPriority`。

- [ ] **Step 5: 运行后端测试**

```bash
cd backend && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 6: 提交**

```bash
git add backend/internal/pkg/database.go backend/internal/service/
git commit -m "chore: set GORM log to Warn, merge duplicate priority validation"
```

---

## 总览

| 阶段 | Task | 主要产出 | 影响文件数 |
|------|------|---------|-----------|
| 1 | 1-4 | 前端去重 | ~15 文件 |
| 2 | 5-8 | 后端错误处理统一 | ~12 文件 |
| 3 | 9-11 | Forum 模块拆分 | ~15 文件（新建+修改） |
| 4 | 12-17 | 基础设施清理 | ~8 文件 |

**总计**：17 个 Task，约 50 个文件变更，全部可独立验证和提交。
