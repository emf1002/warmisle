// e2e/tests/ledger.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { LedgerPage } from '../pages/ledger.page';

test.describe('记账本', () => {
  test('页面加载显示空状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('新增支出记录', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.selectCategory('餐饮');
    await ledger.fillAmount('35.5');
    await ledger.fillNote('午饭');
    await ledger.submit();
    await ledger.expectRecordCount(1);
  });

  test('新增收入记录', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.selectCategory('工资');
    await ledger.fillAmount('10000');
    await ledger.submit();
    await ledger.expectRecordCount(1);
  });

  test('删除记录', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.selectCategory('餐饮');
    await ledger.fillAmount('10');
    await ledger.submit();
    await ledger.expectRecordCount(1);
    await ledger.clickRecord(0);
    await ledger.deleteRecord();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('记账本页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.screenshot('ledger-empty.png');
  });

  test('记账弹窗视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const ledger = new LedgerPage(page);
    await ledger.goto();
    await ledger.openCreate();
    await ledger.screenshotComponent('[data-testid="ledger-modal"]', 'ledger-modal.png');
  });
});
