# 记账本日期选择器 & 移除关联成员设计

## 概述

对记账本模块进行两项改动：
1. 将月份导航替换为日期范围选择器（RangePicker + 快捷月份预设）
2. 移除记录的关联成员（ledger_members 多对多关系），保留 creator_id 用于权限控制

## 动机

- 月份导航只能按月查看，不够灵活，用户希望能按日期范围筛选
- 关联成员功能增加了记账操作的复杂度，实际使用中价值不大，简化为仅保留创建者

## 一、日期范围选择器

### 1.1 前端变更

**替换月份导航组件：**

将当前的 `◀ 2026年5月 ▶` 箭头导航替换为 Ant Design Vue 的 `<a-range-picker>`：

- 默认值：当月第一天 ~ 当月最后一天（`[dayjs().startOf('month'), dayjs().endOf('month')]`）
- `presets` 快捷选项：本月、上月、近三个月、近半年、今年
- 日期格式：`YYYY-MM-DD`
- 选择后自动触发 `fetchLedgers()`

**状态变更：**

- 移除：`selectedMonth: ref<Dayjs>(dayjs())`
- 新增：`dateRange: ref<[Dayjs, Dayjs]>`，初始值为当月范围
- 移除函数：`goPrevMonth()`、`goNextMonth()`、`getMonthParam()`

**API 请求参数变更：**

- 移除 `month` 参数
- 新增 `start_date` 和 `end_date` 参数（格式 `YYYY-MM-DD`）

**摘要数据：**

当前的月度摘要（收入/支出/结余）改为基于所选日期范围的摘要，文案从"X月"改为日期范围显示。

### 1.2 后端变更

**Handler (`handler/ledger.go`)：**

`List` 接口参数变更：
- 移除 `month` 查询参数
- 新增 `start_date`（格式 `YYYY-MM-DD`，默认当月第一天）
- 新增 `end_date`（格式 `YYYY-MM-DD`，默认当月最后一天）
- 验证：`start_date` 不能晚于 `end_date`

**Repository (`repository/ledger.go`)：**

`LedgerFilter` 结构体：
- 移除 `Month string`
- 新增 `StartDate string`、`EndDate string`

查询条件变更：
- 月度过滤：`strftime('%Y-%m', ledgers.occurred_at) = ?`
- 改为日期范围：`ledgers.occurred_at >= ? AND ledgers.occurred_at < ?`
- `StartDate` 对应当天 00:00:00，`EndDate` 对应次日 00:00:00（包含结束日期当天的所有记录）

Summary 计算同样使用日期范围过滤。

## 二、移除关联成员

### 2.1 前端变更

**移除成员筛选器：**

从筛选栏移除"关联成员"下拉框（`filters.member_id`）。保留"创建者"筛选和"分类"筛选。

**移除表单中的成员多选：**

创建/编辑模态框中移除"关联成员" `<a-form-item>` 及其验证逻辑（`member_ids` 非空校验）。

**移除列表中的成员标签：**

每条记录下方不再显示 `关联：😀 😊 ...` 成员头像标签。

**清理接口类型：**

- `LedgerItem` 移除 `members: Member[]` 字段
- `LedgerItem` 保留 `creator_id` 和 `creator`
- 移除 `member_ids` 相关的表单状态和 payload 构造

**成员列表获取：**

如果"创建者"筛选器仍需要成员列表，则保留 `getMembers()` 调用和 `members` ref。仅移除关联成员相关的使用。

### 2.2 后端变更

**Model (`model/ledger.go`)：**

- `Ledger` 结构体移除 `Members []Member` 关联字段和对应的 GORM tag
- 删除 `LedgerMember` 结构体定义

**Repository (`repository/ledger.go`)：**

- `LedgerFilter` 移除 `MemberID *uint` 字段
- `List` 方法：移除 `Preload("Members")`、移除 member_id 子查询过滤
- `Create` 方法：移除 `memberIDs` 参数、移除 `ledger_members` 批量插入
- `Update` 方法：移除 `memberIDs` 参数、移除 `ledger_members` 删除+重插逻辑
- `FindByID` 方法：移除 `Preload("Members")`
- `LedgerWithAssoc` 响应结构体移除 `Members` 字段（如有）

**Service (`service/ledger.go`)：**

- `Create` 方法：移除 `memberIDs` 参数、移除成员验证（`ErrNoMembers` 检查）
- `Update` 方法：移除 `memberIDs` 参数、移除成员替换逻辑（保留/替换判断）
- 移除 `ErrNoMembers` 错误哨兵
- `LedgerService` 如果 `MemberRepo` 仅用于关联成员操作，可移除该依赖

**Handler (`handler/ledger.go`)：**

- `Create` 请求体移除 `MemberIDs []uint` 字段
- `Update` 请求体移除 `MemberIDs []uint` 字段
- `List` 查询参数移除 `member_id`

**错误码 (`service/errors.go`)：**

- 移除 `ErrNoMembers`

### 2.3 数据库迁移

新增迁移文件 `backend/migrations/003_drop_ledger_members.up.sql`：

```sql
-- 移除记账记录关联成员表
DROP TABLE IF EXISTS ledger_members;
```

## 三、保留不变的部分

- `creator_id` 字段和 `Creator` 关联保留
- `canEdit()` 权限逻辑不变（创建者或管理员可编辑/删除）
- "创建者"筛选器保留
- 分类筛选器保留
- 记录的创建、编辑、删除流程（除成员部分）不变

## 四、测试策略

- 更新前端单元测试：移除成员相关的测试用例，新增日期范围选择器的测试
- 更新后端单元测试：移除成员相关的测试用例，新增日期范围过滤的测试
- 更新 handler 测试：验证 API 参数变更
- 手动测试：验证日期范围选择、月份快捷切换、摘要数据正确性

## 五、影响范围

| 文件 | 变更类型 |
|------|---------|
| `frontend/src/views/ledger/Index.vue` | 大量修改（UI + 逻辑） |
| `frontend/src/api/ledger.ts` | 修改接口参数类型 |
| `frontend/src/views/ledger/__tests__/Index.test.ts` | 更新测试 |
| `frontend/src/api/__tests__/ledger.test.ts` | 更新测试 |
| `backend/internal/model/ledger.go` | 移除 LedgerMember、Members 字段 |
| `backend/internal/repository/ledger.go` | 修改 filter、查询逻辑 |
| `backend/internal/service/ledger.go` | 移除成员验证逻辑 |
| `backend/internal/service/errors.go` | 移除 ErrNoMembers |
| `backend/internal/handler/ledger.go` | 修改请求/响应结构 |
| `backend/internal/handler/ledger_test.go` | 更新测试 |
| `backend/internal/service/ledger_test.go` | 更新测试 |
| `backend/migrations/003_drop_ledger_members.up.sql` | 新增迁移 |
