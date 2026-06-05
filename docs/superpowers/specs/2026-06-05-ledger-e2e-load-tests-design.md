# 账本 E2E 负载测试设计

日期：2026-06-05

## 背景

当前账本 E2E 测试仅覆盖基本 CRUD 操作（每种操作创建 1 条记录）。需要扩展测试覆盖范围，模拟用户连续多天大量记账后的页面行为，包括无限滚动分页、统计汇总准确性、日期范围筛选、分类/创建者筛选等场景。

## 设计目标

1. 高效批量播种账单测试数据（30+ 条跨多天记录）
2. 验证无限滚动/游标分页在大量数据下的正确性
3. 验证汇总统计（收入/支出/结余/每日小计）的准确性
4. 验证日期范围筛选和分类/创建者筛选的联动

## 一、后端：批量播种端点

### 端点定义

```
POST /api/test/seed-ledgers
Authorization: Bearer <token>
```

仅在 `HC_TEST_MODE=true` 时可用（与 `/api/test/reset` 同模式）。

### 请求体

```json
{
  "count": 35,
  "start_date": "2026-06-01",
  "end_date": "2026-06-07"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `count` | int | 是 | 要创建的记录总数 |
| `start_date` | string | 否 | 起始日期，默认 7 天前 |
| `end_date` | string | 否 | 结束日期，默认今天 |

### 行为

1. 从数据库获取已有预置分类（由 `/api/test/reset` 创建的 15 支出 + 5 收入）
2. 获取数据库中第一个成员作为 `creator_id`（由 `/api/init/setup` 创建的管理员）
3. 将 `count` 条记录均匀分布在日期范围内：
   - 每天分配 `ceil(count / days)` 条记录
   - 最后一天补齐剩余记录
4. 记录的属性：
   - **金额**：使用固定序列，确保测试可断言。支出金额 = `(index % 10 + 1) * 100` 分（100~1000 分），收入金额 = `50000` 分（500 元）
   - **收支比例**：每 5 条记录中 4 条支出、1 条收入（index % 5 === 4 时为收入）
   - **分类**：支出按顺序轮询 expense 分类，收入按顺序轮询 income 分类
   - **备注**：`"测试账单-{index}"`
   - **时间**：该天的 `12:00:00`

### 响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "count": 35,
    "summary": {
      "income": 350000,
      "expense": 6300,
      "balance": 343700
    },
    "expense_category_count": 15,
    "income_category_count": 5
  }
}
```

后端在播种完成后立即查询计算 summary 并返回，测试可直接用于断言。

### 实现

- 新增文件：`backend/internal/handler/test_seed_ledgers.go`
- 路由注册：`backend/internal/routes/router.go` 中 test reset 路由旁添加

## 二、前端组件：添加 data-testid

修改 `frontend/src/views/ledger/Index.vue`，为筛选栏和汇总栏添加测试标识：

| 元素 | 添加属性 | 位置 |
|------|----------|------|
| 分类筛选 `<a-select>` | `data-testid="filter-category"` | 模板第 40 行附近 |
| 创建者筛选 `<a-select>` | `data-testid="filter-creator"` | 模板第 60 行附近 |
| 收入汇总 `<span>` | `data-testid="summary-income"` | 模板第 24 行 |
| 支出汇总 `<span>` | `data-testid="summary-expense"` | 模板第 27 行 |
| 结余汇总 `<span>` | `data-testid="summary-balance"` | 模板第 32 行 |
| 日期分组 `<div>` | `data-testid="date-group"` | 模板第 87 行 |
| 每日小计 `<span>` | `data-testid="daily-total"` | 模板第 91 行 |
| 加载哨兵 `<div>` | `data-testid="load-sentinel"` | 模板第 130 行 |

## 三、E2E 页面对象修复与扩展

修改 `e2e/pages/ledger.page.ts`：

### 修复

- `selectCategory(name)` → 重命名为 `pickCategory(name)`，改为点击网格选择器 `.category-pick-item`（匹配文本），替代已失效的 `data-testid="category-select"` 定位

### 新增方法

| 方法 | 说明 |
|------|------|
| `filterByCategory(name)` | 操作滤选栏的分类 select（`data-testid="filter-category"`） |
| `filterByCreator(name)` | 操作滤选栏的创建者 select（`data-testid="filter-creator"`） |
| `scrollToLoadMore()` | 滚动到 sentinel 元素（`data-testid="load-sentinel"`），触发 IntersectionObserver |
| `expectSummary({income?, expense?, balance?})` | 断言汇总栏文本 |
| `expectDateGroupCount(n)` | 断言 `[data-testid="date-group"]` 数量 |
| `expectTotalItemCount(n)` | 断言所有 `[data-testid^="ledger-item-"]` 总数 |
| `expectDailyTotal(index, text)` | 断言第 n 个日期组的每日小计文本 |

## 四、E2E Fixture 扩展

### `e2e/fixtures/db.fixture.ts`

新增 `seedLedgers(request, token, options?)` 函数：

```typescript
export async function seedLedgers(
  request: APIRequestContext,
  token: string,
  options?: { count?: number; startDate?: string; endDate?: string }
) {
  const res = await request.post('/api/test/seed-ledgers', {
    data: {
      count: options?.count ?? 35,
      start_date: options?.startDate,
      end_date: options?.endDate,
    },
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.json();
}
```

### `e2e/fixtures/auth.fixture.ts`

扩展 `authenticated` fixture，额外返回 `seedLedgers(options?)` 便捷方法（已绑定 token 和 request）。

## 五、新增测试文件

新增 `e2e/tests/ledger-load.spec.ts`，包含以下场景：

### 场景 1：无限滚动/分页

**前置**：播种 35 条记录跨 7 天（默认 limit=20，需要 2 页）

**步骤**：
1. 打开账本页面
2. 断言首次加载显示 20 条记录（`expectTotalItemCount(20)`）
3. 断言 sentinel 可见（`v-if="hasMore"` 渲染了 `data-testid="load-sentinel"`）
4. 滚动到 sentinel 触发加载更多
5. 等待加载完成（sentinel 消失或 spinner 消失）
6. 断言总记录数 = 35
7. 断言日期分组数量 = 7
8. 断言 sentinel 不再渲染

### 场景 2：统计汇总准确性

**前置**：播种 35 条记录

**步骤**：
1. 打开账本页面
2. 等待数据加载完成
3. 从 seed 响应中获取预期 summary
4. 断言汇总栏收入显示 `+¥{income/100}.00`
5. 断言汇总栏支出显示 `-¥{expense/100}.00`
6. 断言汇总栏结余显示正确符号和金额
7. 滚动加载全部数据后，断言汇总值不变（summary 基于全量数据计算）

### 场景 3：日期范围筛选

**前置**：播种 35 条记录跨 6 月 1 日 ~ 6 月 7 日

**步骤**：
1. 打开账本页面
2. 通过 RangePicker 选择 6 月 3 日 ~ 6 月 5 日
3. 等待数据重新加载
4. 断言记录数 < 35（仅 3 天的数据）
5. 断言 summary 更新（与全量不同）
6. 断言所有可见记录的日期在筛选范围内

### 场景 4：分类/创建者筛选

**前置**：播种 35 条记录（包含多种分类）

**步骤**：
1. 打开账本页面
2. 通过 `filterByCategory('餐饮')` 筛选
3. 断言只显示餐饮分类的记录
4. 断言 summary 更新为仅餐饮的汇总
5. 点击"清除筛选"
6. 断言恢复全部记录

## 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `backend/internal/handler/test_seed_ledgers.go` | 新增 | 批量播种端点 |
| `backend/internal/routes/router.go` | 修改 | 添加 1 行路由 |
| `frontend/src/views/ledger/Index.vue` | 修改 | 添加 8 个 data-testid |
| `e2e/pages/ledger.page.ts` | 修改 | 修复 selectCategory + 新增 7 个方法 |
| `e2e/fixtures/db.fixture.ts` | 修改 | 新增 seedLedgers 函数 |
| `e2e/fixtures/auth.fixture.ts` | 修改 | 扩展 authenticated fixture |
| `e2e/tests/ledger-load.spec.ts` | 新增 | 4 个测试场景 |

## 验收标准

1. `make e2e` 全部通过（含现有 46 个 + 新增 4 个测试）
2. 新增测试可重复运行（每次测试前 reset + seed 保证隔离）
3. 批量播种 35 条记录耗时 < 1 秒
4. 无限滚动测试验证分页完整性（所有记录最终可见）
5. 汇总统计与 seed 端点返回值一致
