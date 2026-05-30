# 代码库优化设计

日期：2026-05-30
状态：已批准
实施策略：分阶段（方案 B：影响面优先）

## 背景

项目经过 2 周的快速开发，积累了大量代码重复和架构不一致问题。分析发现：
- 前端 JWT 解析在 7 个文件中重复，CSS 公共样式在 6+ 个视图重复
- 后端 handler 错误处理样板代码 ~40+ 处，service 层 not-found 映射 ~25 处
- forum 模块三层共 1223 行，合并了 6 个子资源
- 基础设施存在 Makefile 缺少 test target、迁移缺少 Down 块等问题

## 阶段 1：前端代码重复提取

### 1.1 JWT 解析提取到 Auth Store

**目标**：消除 7 个文件中的 JWT 解析重复

**修改文件**：
- `frontend/src/stores/auth.ts` — 新增 `currentUserId`、`currentUserRole`、`isAdmin` 计算属性
- `frontend/src/views/ledger/Index.vue` — 删除本地 JWT 解析，改用 auth store
- `frontend/src/views/todo/Index.vue` — 同上
- `frontend/src/views/wish/Index.vue` — 同上
- `frontend/src/views/forum/Index.vue` — 同上
- `frontend/src/views/forum/TopicDetail.vue` — 同上
- `frontend/src/views/member/Index.vue` — 同上
- `frontend/src/views/category/Index.vue` — 同上

**实现细节**：
```ts
// stores/auth.ts 扩展
const currentUserId = computed(() => {
  if (!token.value) return null
  const payload = JSON.parse(atob(token.value.split('.')[1]))
  return payload.member_id
})
const currentUserRole = computed(() => {
  if (!token.value) return null
  const payload = JSON.parse(atob(token.value.split('.')[1]))
  return payload.role
})
const isAdmin = computed(() => currentUserRole.value === 'admin')
```

各视图改为 `const authStore = useAuthStore()` 后直接使用 `authStore.currentUserId` 等。

### 1.2 工具函数提取

**目标**：消除 5 个工具函数的重复定义，统一实现

**新建文件**：`frontend/src/utils/format.ts`

| 函数 | 原重复位置 | 统一行为 |
|------|-----------|---------|
| `truncate(str, len)` | ledger, forum (2处) | 截断并加 `...` |
| `timeAgo(date)` | dashboard, forum, TopicDetail (3处) | 统一包含"刚刚"的版本 |
| `formatDate(date)` | dashboard(dayjs), ledger(Date), todo(Date) | 统一用 dayjs，格式 `M月D日 周X` |
| `priorityColor(p)` | todo, wish (2处) | 返回颜色值 |
| `priorityLabel(p)` | todo, wish (2处) | 返回中文标签 |

**修改文件**：删除各视图中的本地定义，改为 `import { truncate, timeAgo, ... } from '@/utils/format'`

### 1.3 CSS 公共类提取

**目标**：消除 6+ 个视图中的重复 CSS

**新建文件**：`frontend/src/styles/components.css`

| 类名 | 重复文件数 | 来源 |
|------|-----------|------|
| `.page-header` | 6 | todo, wish, forum, member, category, profile |
| `.loading-state` | 5 | ledger, todo, wish, forum, TopicDetail |
| `.empty-state` | 4 | ledger, todo, wish, forum |
| `.filter-row` | 2 | ledger, todo |
| `.emoji-picker` | 3 | member, category, profile |
| `.pagination-row` | 3 | todo, wish, forum |

**修改文件**：
- 各视图删除对应的 scoped CSS
- `main.ts` 中 import `components.css`
- `global.css` 中删除未使用的 `.income-color`、`.expense-color`

### 1.4 EmojiPicker 组件提取

**目标**：消除 3 个视图中 emoji 选择器的重复

**新建文件**：`frontend/src/components/EmojiPicker.vue`

```vue
<script setup lang="ts">
defineProps<{ modelValue: string }>()
defineEmits<{ 'update:modelValue': [val: string] }>()
const emojiList = ['😀', '😂', '🥰', '😎', '🤔', '😴', '🥳', '😇', '🤗', '😋', '🤩', '😏', '😤', '🤯', '🫡', '🤭', '😌', '🤓', '🐶', '🐱', '🐼', '🦊', '🐸', '🐵', '🏠', '🌟', '🔥', '❤️', '💪', '🎉', '🌈', '☕']
</script>
```

**修改文件**：member、category、profile 视图改用 `<EmojiPicker v-model="form.avatar" />`

## 阶段 2：后端错误处理统一

### 2.1 Handler 错误映射 Helper

**目标**：消除 ~40 处 handler 错误处理样板代码

**新建文件**：`backend/internal/handler/errors.go`

```go
type serviceError struct {
    err      error
    httpCode int
    bizCode  int
    msg      string
}

func handleServiceError(c *gin.Context, err error, mappings ...serviceError) {
    for _, m := range mappings {
        if errors.Is(err, m.err) {
            pkg.Error(c, m.httpCode, m.bizCode, m.msg)
            return
        }
    }
    pkg.Error(c, 500, 50001, "服务器内部错误")
}

func getMemberID(c *gin.Context) uint { ... }
func getRole(c *gin.Context) string { ... }
```

**修改文件**：所有 handler 文件中的 error handling 改用 `handleServiceError`

### 2.2 Service 层通用 Not-Found 映射

**目标**：消除 ~25 处 `gorm.ErrRecordNotFound` 重复检查

**新建文件**：`backend/internal/service/errors.go`

```go
func wrapNotFound(err error, domainErr error) error {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return domainErr
    }
    return err
}
```

**修改文件**：所有 service 文件中的 not-found 检查改用 `wrapNotFound`

### 2.3 错误码规范

| bizCode | 含义 | HTTP Status |
|---------|------|-------------|
| 40001 | 参数校验失败 | 400 |
| 40002 | 资源已存在（唯一约束） | 409 |
| 40401 | 资源不存在 | 404（新增） |
| 40301 | 权限不足 | 403 |
| 50001 | 服务器内部错误 | 500 |

**注意**：将现有 `40001`（资源不存在）迁移到 `40401`，与参数校验错误区分。

### 2.4 辅助清理

- `getMemberID()`/`getRole()` 从 `handler/ledger.go` 移到 `handler/errors.go`
- `test_reset.go` 改用 `pkg.Success()`/`pkg.Error()`

## 阶段 3：Forum 模块拆分

### 拆分方案

| 新文件 | 包含的端点 |
|--------|-----------|
| `handler/forum_post.go` | Post CRUD |
| `handler/forum_topic.go` | Topic CRUD + pin |
| `handler/forum_comment.go` | Comment CRUD |
| `handler/forum_vote.go` | Vote + Like |
| `handler/forum_feed.go` | Feed list |
| `service/forum_post.go` | Post 业务逻辑 |
| `service/forum_topic.go` | Topic 业务逻辑 |
| `service/forum_comment.go` | Comment 业务逻辑 |
| `service/forum_vote.go` | Vote + Like 业务逻辑 |
| `repository/forum_post.go` | Post 数据访问 |
| `repository/forum_topic.go` | Topic 数据访问 |
| `repository/forum_comment.go` | Comment 数据访问 |
| `repository/forum_vote.go` | Vote + Like 数据访问 |

### 共享逻辑

Forum 子资源间共享的 `FeedItem` 结构体统一定义在 `repository/forum.go`（保留原文件作为公共定义），dashboard service 引用此定义。

### 路由注册

`routes/router.go` 中的 forum 路由组不需要改动，handler 函数名保持一致。

## 阶段 4：基础设施清理

| 任务 | 文件 | 说明 |
|------|------|------|
| Makefile 添加 test/lint/fmt | `Makefile` | `test: go test ./...`，`lint: golangci-lint run`，`fmt: gofmt -w .` |
| 迁移 Down 块 | `migrations/002_fix_time_format.up.sql` | 添加 `+goose Down` 空块 |
| 删除死代码 | `internal/pkg/response.go` | 删除未使用的 `PageData` 结构体 |
| 路由测试统一 | `internal/handler/test_helpers_test.go` | 改用 `routes.Register()` |
| GORM 日志级别 | `internal/pkg/database.go` | `logger.Info` → `logger.Warn` |
| 合并重复函数 | `service/common.go` | `validPriority`/`validWishPriority` 合并 |
| FeedItem 统一 | `repository/forum.go`, `service/dashboard.go` | dashboard 引用 repo 的定义 |

## 验收标准

- [ ] 前端：JWT 解析仅在 `stores/auth.ts` 中存在
- [ ] 前端：`utils/format.ts` 包含 5 个工具函数，无视图中有重复定义
- [ ] 前端：`styles/components.css` 包含 6 类公共样式，无视图中有重复定义
- [ ] 前端：`EmojiPicker.vue` 组件被 3 个视图引用
- [ ] 后端：所有 handler 使用 `handleServiceError`，无手动 error checking 样板
- [ ] 后端：所有 service 使用 `wrapNotFound`，无手动 `gorm.ErrRecordNotFound` 检查
- [ ] 后端：`40401` 用于资源不存在，`40001` 仅用于参数校验
- [ ] 后端：forum 模块拆分为 5 个子资源文件（每层）
- [ ] Makefile 包含 test、lint、fmt target
- [ ] 所有现有测试通过
