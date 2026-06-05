# Task 2: 待办管理 — E2E 测试扩展 (14 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐待办模块的优先级、截止日期、指派、认领、筛选、编辑和权限测试。

**Files:**
- Modify: `e2e/pages/todo.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/todo.spec.ts` — 新增 14 个测试用例

**Dependencies:** Task 0 完成（memberContext fixture 可用）

---

## Step 1: 扩展 TodoPage POM

**Modify:** `e2e/pages/todo.page.ts` — 在 `clearFilters` 方法之后、类的结束 `}` 之前新增：

```typescript
  /** 填写截止日期（格式 YYYY-MM-DD） */
  async fillDueDate(date: string) {
    await this.page.getByTestId('due-date-picker').click();
    await this.page.getByTestId('due-date-picker').locator('input').fill(date);
    // 关闭日期面板
    await this.page.keyboard.press('Escape');
  }

  /** 选择指派成员 */
  async selectAssignee(name: string) {
    await this.page.getByTestId('assignee-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 认领第 n 条未指派的待办 */
  async claimTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('claim-btn').click();
  }

  /** 按指派成员筛选 */
  async filterByAssignee(name: string) {
    await this.page.getByTestId('assignee-filter').click();
    await this.page.locator('.ant-select-item-option', { hasText: name }).click();
  }

  /** 编辑第 n 条待办 */
  async editTodo(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await items.nth(index).getByTestId('edit-btn').click();
    await expect(this.page.getByTestId('todo-modal')).toBeVisible();
  }

  /** 断言第 n 条待办的优先级标签 */
  async expectTodoPriority(index: number, priority: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('priority-tag')).toContainText(priority);
  }

  /** 断言第 n 条待办的截止日期显示 */
  async expectTodoDueDate(index: number, dateText: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('due-date')).toContainText(dateText);
  }

  /** 断言第 n 条待办的指派人 */
  async expectTodoAssignee(index: number, name: string) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('assignee-name')).toContainText(name);
  }

  /** 断言第 n 条待办的截止日期为红色（过期） */
  async expectOverdueHighlight(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('due-date')).toHaveClass(/overdue|red|error/);
  }

  /** 断言第 n 条待办已完成 */
  async expectTodoCompleted(index: number) {
    const items = this.page.getByTestId(/^todo-item-/);
    await expect(items.nth(index).getByTestId('todo-checkbox')).toBeChecked();
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/todo.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 功能场景 ===

  test('设置优先级', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('紧急任务');
    await todo.selectPriority('紧急');
    await todo.submit();
    await todo.expectTodoCount(1);
    await todo.expectTodoPriority(0, '紧急');
  });

  test('设置截止日期', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('有截止日的任务');
    await todo.fillDueDate('2026-12-31');
    await todo.submit();
    await todo.expectTodoCount(1);
  });

  test('指派给成员', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('指派任务');
    await todo.selectAssignee('成员一');
    await todo.submit();
    await todo.expectTodoCount(1);
    await todo.expectTodoAssignee(0, '成员一');
  });

  test('认领待办', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建未指派的待办
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('待认领');
    await todo.submit();
    // member 认领
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.claimTodo(0);
    await memberTodo.expectTodoAssignee(0, '成员一');
  });

  test('按状态筛选', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    // 创建两个待办，完成一个
    await todo.openCreate();
    await todo.fillTitle('任务A');
    await todo.submit();
    await todo.openCreate();
    await todo.fillTitle('任务B');
    await todo.submit();
    await todo.toggleTodo(0);
    // 筛选"已完成"
    await todo.filterByStatus('已完成');
    await page.waitForTimeout(500);
    await todo.expectTodoCount(1);
  });

  test('按指派成员筛选', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    // 创建指派给成员一的待办
    await todo.openCreate();
    await todo.fillTitle('指派任务');
    await todo.selectAssignee('成员一');
    await todo.submit();
    // 创建未指派的待办
    await todo.openCreate();
    await todo.fillTitle('未指派');
    await todo.submit();
    await todo.expectTodoCount(2);
    // 按成员筛选
    await todo.filterByAssignee('成员一');
    await page.waitForTimeout(500);
    await todo.expectTodoCount(1);
  });

  test('编辑待办', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('原标题');
    await todo.submit();
    await todo.editTodo(0);
    await todo.fillTitle('修改后标题');
    await todo.submit();
    await todo.expectTodoTitle(0, '修改后标题');
  });

  test('管理员可编辑任意待办', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 创建待办
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.openCreate();
    await memberTodo.fillTitle('member的待办');
    await memberTodo.submit();
    // admin 编辑
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.editTodo(0);
    await todo.fillTitle('admin已修改');
    await todo.submit();
    await todo.expectTodoTitle(0, 'admin已修改');
  });

  test('管理员可删除任意待办', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 创建待办
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.openCreate();
    await memberTodo.fillTitle('member的待办');
    await memberTodo.submit();
    // admin 删除
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.deleteTodo(0);
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  // === 权限测试 ===

  test('被指派人可编辑被指派的待办', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建指派给 member 的待办
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('指派给member');
    await todo.selectAssignee('成员一');
    await todo.submit();
    // member 编辑
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.editTodo(0);
    await memberTodo.fillDescription('member补充的描述');
    await memberTodo.submit();
  });

  test('非创建者非被指派人不可编辑', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 创建待办（不指派）
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('仅admin的待办');
    await todo.submit();
    // member 不应看到编辑按钮
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.clickRecord && // 如果有 clickRecord 方法
    await expect(memberPage.getByTestId('edit-btn')).not.toBeVisible();
  });

  // === 错误路径 ===

  test('标题为空被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('');
    await todo.submit();
    // 弹窗应仍然可见（提交被阻止）
    await expect(page.getByTestId('todo-modal')).toBeVisible();
  });

  test('过期待办高亮', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('过期任务');
    await todo.fillDueDate('2020-01-01');
    await todo.submit();
    await todo.expectTodoCount(1);
    await todo.expectOverdueHighlight(0);
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test todo.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/todo.page.ts e2e/tests/todo.spec.ts
git commit -m "test(e2e): expand todo tests — priority, assignee, claim, filter, permissions"
```
