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
    await this.expectModalVisible();
  }

  async openChangePassword() {
    await this.page.getByTestId('change-pwd-btn').click();
    await this.expectModalVisible();
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
    await this.submitModal();
  }

  async submitPassword() {
    await this.submitModal();
  }

  async logout() {
    await this.page.getByTestId('logout-btn').click();
    await this.confirmModal();
  }

  async expectProfileName(name: string) {
    await expect(this.page.getByTestId('profile-name')).toContainText(name);
  }

  async selectAvatar(avatar: string) {
    await this.page.getByTestId('avatar-picker').click();
    await this.page.getByTestId(`icon-item-${avatar}`).click();
  }

  async fillConfirmPassword(password: string) {
    await this.page.getByTestId('confirm-pwd-input').fill(password);
  }

  async expectOldPasswordError() {
    await expect(this.page.locator('.ant-form-item-explain-error').filter({ hasText: '原密码错误' })).toBeVisible();
  }

  async expectPasswordMismatchError() {
    await expect(this.page.locator('.ant-form-item-explain-error').filter({ hasText: '两次输入的密码不一致' })).toBeVisible();
  }

  async expectSamePasswordError() {
    await expect(this.page.locator('.ant-form-item-explain-error').filter({ hasText: '新密码不能与旧密码相同' })).toBeVisible();
  }

  async expectNameTooLongError() {
    await expect(this.page.locator('.ant-form-item-explain-error').filter({ hasText: '请输入1-20字符' })).toBeVisible();
  }
}
