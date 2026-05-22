import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class MembersPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/members');
  }

  async expectOnMembers() {
    await expect(this.page).toHaveURL(/\/#\/members/);
  }

  async openCreate() {
    await this.page.getByTestId('add-btn').click();
    await expect(this.page.getByTestId('member-modal')).toBeVisible();
  }

  async fillUsername(username: string) {
    await this.page.getByTestId('username-input').fill(username);
  }

  async fillPassword(password: string) {
    await this.page.getByTestId('password-input').fill(password);
  }

  async fillName(name: string) {
    await this.page.getByTestId('name-input').fill(name);
  }

  async selectRole(role: string) {
    await this.page.getByTestId('role-select').click();
    await this.page.locator('.ant-select-item-option', { hasText: role }).click();
  }

  async submit() {
    await this.page.locator('.ant-modal-footer .ant-btn-primary').click();
  }

  async expectMemberCount(count: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await expect(rows).toHaveCount(count);
  }

  async editMember(index: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await rows.nth(index).getByTestId('edit-btn').click();
  }

  async disableMember(index: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await rows.nth(index).getByTestId('disable-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  async enableMember(index: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await rows.nth(index).getByTestId('enable-btn').click();
  }

  async deleteMember(index: number) {
    const rows = this.page.getByTestId('member-table').locator('tbody tr');
    await rows.nth(index).getByTestId('delete-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }
}
