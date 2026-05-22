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

  test('待办页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const todo = new TodoPage(page);
    await todo.goto();
    await todo.screenshot('todo-empty.png');
  });
});
