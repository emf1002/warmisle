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
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async openChangePassword() {
    await this.page.getByTestId('change-pwd-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
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
    await this.page.locator('.ant-modal:visible .ant-modal-footer .ant-btn-primary').click();
  }

  async submitPassword() {
    await this.page.locator('.ant-modal:visible .ant-modal-footer .ant-btn-primary').click();
  }

  async logout() {
    await this.page.getByTestId('logout-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
  }

  async expectProfileName(name: string) {
    await expect(this.page.locator('.profile-header')).toContainText(name);
  }

  async selectAvatar(avatar: string) {
    await this.page.getByTestId('avatar-picker').click();
    await this.page.getByTestId(`icon-item-${avatar}`).click();
  }

  async fillConfirmPassword(password: string) {
    await this.page.getByTestId('confirm-pwd-input').fill(password);
  }

  async expectOldPasswordError() {
    await expect(this.page.locator('.ant-modal:visible .ant-form-item:has-text("当前密码") .ant-form-item-explain')).toContainText('原密码错误');
  }

  async expectPasswordMismatchError() {
    await expect(this.page.locator('.ant-modal:visible .ant-form-item:has-text("确认新密码") .ant-form-item-explain')).toContainText('两次输入的密码不一致');
  }

  async expectSamePasswordError() {
    await expect(this.page.locator('.ant-modal:visible .ant-form-item:has-text("新密码") .ant-form-item-explain')).toContainText('新密码不能与旧密码相同');
  }

  async expectNameTooLongError() {
    await expect(this.page.locator('.ant-modal:visible .ant-form-item:has-text("姓名") .ant-form-item-explain')).toContainText('请输入1-20字符');
  }
}
