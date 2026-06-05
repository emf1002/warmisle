# 记账本加载方式优化设计

## 问题描述

当前记账本页面的加载机制存在以下问题：

1. **分页与分组边界冲突**：后端使用 `offset/limit` 分页，再在内存中按日期分组（`groupByDate`）。当某天的记录跨两页时，该天会出现两次，每次都是不完整的分组，`daily_total` 也不正确。
2. **汇总与筛选不一致**：顶部收入/支出/结余的汇总查询不带 `category_id` 和 `creator_id` 筛选条件。用户筛选"本月餐饮"时，汇总显示的是所有分类的总额，而非餐饮的总额。
3. **无前端缓存**：分类和成员列表在每次页面挂载时重新请求，离开再进入页面会重复加载。
4. **加载体验差**：使用居中 spinner，无骨架屏，存在布局抖动。
5. **无请求去重**：快速切换筛选条件会发出多次并发请求，最后一次响应不一定最后到达。

## 设计方案

### 1. 后端：游标分页替代 offset/limit

**API 变更**

`GET /api/ledgers` 请求参数变更：

| 参数 | 类型 | 说明 |
|------|------|------|
| start_date | string | 起始日期（含），默认当月 1 号 |
| end_date | string | 结束日期（不含），默认下月 1 号 |
| category_id | uint | 可选，分类筛选 |
| creator_id | uint | 可选，创建者筛选 |
| limit | int | 每页基准条数，默认 20 |
| cursor | string | 可选，游标（上一页返回的 `next_cursor`） |

移除 `page` 和 `page_size` 参数。

**游标格式**

游标为 JSON 的 base64 编码：`base64({"occurred_at":"2026-06-03T15:30:00","id":42})`。
使用 `(occurred_at, id)` 二元组保证在相同时间戳下的确定性排序。

**查询逻辑**

1. 基础查询：`WHERE occurred_at >= start_date AND occurred_at < end_date` + 可选筛选条件
2. 有游标时追加：`AND (occurred_at < ? OR (occurred_at = ? AND id < ?))`
3. 排序：`ORDER BY occurred_at DESC, id DESC`
4. 取 `limit + 1` 条
5. 对前 `limit` 条做内存分组（`groupByDate`）
6. 若第 `limit` 条所在日期被截断（该日期还有更多记录未取到），追加查询把该日剩余记录一并补全
7. 第 `limit + 1` 条存在则 `has_more = true`，`next_cursor` 取返回结果中最后一条的 `(occurred_at, id)`

**补全策略细节**

当分组后最后一天的记录被截断（即 limit 切到了某天的中间），执行一次补充查询：
```sql
SELECT * FROM ledgers
WHERE occurred_at = :last_date AND id < :last_id_of_that_date
  AND deleted_at IS NULL
ORDER BY occurred_at DESC, id DESC
```
将补充结果合并到最后一个分组。这保证每个分页返回的日期分组都是完整的。

**响应结构变更**

```json
{
  "summary": { "income": 12345, "expense": 6789, "balance": 5556 },
  "groups": [
    { "date": "2026-06-03", "daily_total": -500, "items": [...] },
    { "date": "2026-06-01", "daily_total": 2000, "items": [...] }
  ],
  "next_cursor": "eyJvY2N1cnJlZF9hdCI6IjIwMjYtMDYtMDFUMTA6MDA6MDAiLCJpZCI6Mzh9",
  "has_more": true
}
```

移除 `total`、`page`、`page_size` 字段。

### 2. 后端：汇总查询与筛选条件对齐

当前汇总查询只按日期范围过滤，不带 `category_id` 和 `creator_id`。变更：当筛选条件存在时追加对应的 WHERE 子句，使汇总与列表使用相同的筛选条件。

### 3. 前端：Pinia 缓存分类和成员

新增 `frontend/src/stores/categories.ts` 和 `frontend/src/stores/members.ts`。

**categories store**：
- `categories` ref，`loaded` flag
- `fetchCategories()` — 若 `loaded` 为 true 直接返回缓存，否则调用 API 并缓存
- `reset()` — 清空缓存（退出登录时调用）

**members store**：同上结构。

记账本页面的 `onMounted` 中从 store 获取，不再直接调用 API。

### 4. 前端：加载体验优化

**骨架屏替代 spinner**

用 Ant Design Vue 的 `<a-skeleton>` 替代居中 `<a-spin>`。骨架屏模拟日期头 + 3 行记录卡片的形状，避免加载前后的布局抖动。

**IntersectionObserver 无限滚动**

移除"加载更多"按钮。使用 IntersectionObserver 监听列表底部 sentinel 元素进入视口时自动触发加载下一页。滚动加载时底部显示小 spinner。

**筛选防抖**

筛选条件变更（分类、创建者）加 300ms debounce，避免快速切换产生多次请求。

**请求取消**

每次发起新的列表请求时，通过 `AbortController` 取消上一次未完成的请求，确保只有最后一次请求的响应被处理。

### 5. 前端：游标状态管理

替换现有的 `page`/`pageSize`/`total` 状态为：

```typescript
const nextCursor = ref<string | null>(null)
const hasMore = ref(false)
```

筛选条件或日期范围变更时重置 cursor 为 null（触发重新从头加载）。加载更多时携带当前 `nextCursor`，新返回的 groups 追加到现有列表。

## 涉及文件

### 后端
- `backend/internal/repository/ledger.go` — 重写 `List` 方法：游标解析、分页查询、补全策略；重写汇总查询使其应用筛选条件
- `backend/internal/handler/ledger.go` — 更新 `List` handler 的参数绑定（`limit`/`cursor` 替代 `page`/`page_size`）
- `backend/internal/service/ledger.go` — 更新 `LedgerFilter` 结构体（`Limit`/`Cursor` 替代 `Page`/`PageSize`）
- `backend/internal/repository/ledger.go` 中的 `LedgerFilter` 和 `ListResult` 结构体同步更新
- `backend/internal/service/ledger_test.go` — 更新列表相关测试用例
- `backend/internal/handler/ledger_test.go` — 更新 handler 列表测试用例

### 前端
- `frontend/src/views/ledger/Index.vue` — 骨架屏、无限滚动、防抖、请求取消、游标状态
- `frontend/src/api/ledger.ts` — 更新 `getLedgers` 参数类型（`limit`/`cursor` 替代 `page`/`page_size`）
- `frontend/src/stores/categories.ts` — 新增，分类缓存 store
- `frontend/src/stores/members.ts` — 新增，成员缓存 store

### 不需要变更的
- 数据库迁移：现有索引 `(deleted_at, occurred_at)` 已覆盖游标查询模式，无需新增索引
- 路由注册：API 路径 `/api/ledgers` 不变
- 其他页面：不影响其他模块

## 边界与约束

- 游标 base64 编码/解码在后端完成，前端只传递 opaque string
- 第一页不传 cursor，后端返回完整日期分组
- 每页实际返回条数可能略多于 limit（因补全策略），但不会超过 limit + 单日最大记录数
- 筛选条件变更时前端重置列表和 cursor，从头开始加载
- `AbortController` 只在筛选/日期范围变更时取消前序请求，追加加载（loadMore）不取消

## 测试策略

### 后端
- 游标分页：第一页无 cursor 返回正确；第二页带 cursor 从正确位置继续
- 日期补全：limit 切断某天时自动补全该天所有记录
- 筛选一致性：带 category_id 时，summary 的收入/支出与 groups 中的数据一致
- 边界：空结果返回空 groups 数组和零 summary

### 前端
- 骨架屏：loading 状态下显示 skeleton，加载完成后显示列表
- 无限滚动：滚动到底部自动触发加载，has_more 为 false 时停止
- 防抖：快速切换筛选只触发一次请求
- 请求取消：切换筛选时前序请求被取消
