// e2e/tests/todo.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { TodoPage } from '../pages/todo.page';

test.describe('待办管理', () => {
  test('页面加载显示空状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('新建待办', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('买菜');
    await todo.fillDescription('西红柿、鸡蛋、青菜');
    await todo.submit();
    await todo.expectTodoCount(1);
    await todo.expectTodoTitle(0, '买菜');
  });

  test('切换待办完成状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('洗碗');
    await todo.submit();
    await todo.toggleTodo(0);
    // 待办应标记为已完成
  });

  test('删除待办', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('临时任务');
    await todo.submit();
    await todo.expectTodoCount(1);
    await todo.deleteTodo(0);
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

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
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('待认领');
    await todo.submit();
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.claimTodo(0);
    await memberTodo.expectTodoAssignee(0, '成员一');
  });

  test('按状态筛选', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('任务A');
    await todo.submit();
    await todo.openCreate();
    await todo.fillTitle('任务B');
    await todo.submit();
    await todo.toggleTodo(0);
    await todo.filterByStatus('已完成');
    await page.waitForTimeout(500);
    await todo.expectTodoCount(1);
  });

  test('按指派成员筛选', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('指派任务');
    await todo.selectAssignee('成员一');
    await todo.submit();
    await todo.openCreate();
    await todo.fillTitle('未指派');
    await todo.submit();
    await todo.expectTodoCount(2);
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
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.openCreate();
    await memberTodo.fillTitle('member的待办');
    await memberTodo.submit();
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
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.openCreate();
    await memberTodo.fillTitle('member的待办');
    await memberTodo.submit();
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.deleteTodo(0);
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  // === 权限测试 ===

  test('被指派人可编辑被指派的待办', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('指派给member');
    await todo.selectAssignee('成员一');
    await todo.submit();
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
    await memberTodo.editTodo(0);
    await memberTodo.fillDescription('member补充的描述');
    await memberTodo.submit();
  });

  test('非创建者非被指派人不可编辑', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.openCreate();
    await todo.fillTitle('仅admin的待办');
    await todo.submit();
    const memberTodo = new TodoPage(memberPage);
    await memberTodo.goto();
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
    await expect(page.locator('.ant-modal-wrap:visible')).toBeVisible();
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
});
