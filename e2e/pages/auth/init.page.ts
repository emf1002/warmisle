import { type Page, expect } from '@playwright/test';
import { BasePage } from '../base.page';

export class InitPage extends BasePage {
  private nameInput = this.page.getByTestId('name-input');
  private usernameInput = this.page.getByTestId('username-input');
  private passwordInput = this.page.getByTestId('password-input');
  private initBtn = this.page.getByTestId('init-btn');

  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/init');
  }

  async setup(name: string, username: string, password: string) {
    await this.nameInput.fill(name);
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
    await this.initBtn.click();
  }

  async expectOnInitPage() {
    await expect(this.page).toHaveURL(/\/#\/init/);
  }
}
