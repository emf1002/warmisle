import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class ProfilePage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/profile');
  }

  async expectOnProfile() {
    await expect(this.page).toHaveURL(/\/#\/profile/);
  }

  async openEditProfile() {
    await this.page.getByTestId('edit-profile-btn').click();
    await expect(this.page.getByTestId('profile-modal')).toBeVisible();
  }

  async openChangePassword() {
    await this.page.getByTestId('change-pwd-btn').click();
    await expect(this.page.getByTestId('pwd-modal')).toBeVisible();
  }

  async fillName(name: string) {
    await this.page.getByTestId('name-input').fill(name);
  }

  async fillOldPassword(password: string) {
    await this.page.getByTestId('old-pwd-input').fill(password);
  }

  async fillNewPassword(password: string) {
    await this.page.getByTestId('new-pwd-input').fill(password);
  }

  async submitProfile() {
    await this.page.getByTestId('profile-modal').locator('.ant-modal-footer .ant-btn-primary').click();
  }

  async submitPassword() {
    await this.page.getByTestId('pwd-modal').locator('.ant-modal-footer .ant-btn-primary').click();
  }

  async logout() {
    await this.page.getByTestId('logout-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
  }

  async expectProfileName(name: string) {
    await expect(this.page.locator('.profile-header')).toContainText(name);
  }
}
