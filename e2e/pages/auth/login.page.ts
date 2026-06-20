import { type Page, expect } from '@playwright/test';
import { BasePage } from '../base.page';

export class LoginPage extends BasePage {
  private usernameInput = this.page.getByTestId('username-input');
  private passwordInput = this.page.getByTestId('password-input');
  private loginBtn = this.page.getByTestId('login-btn');

  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/login');
  }

  async login(username: string, password: string) {
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
    await this.loginBtn.click();
  }

  async expectLoginError(message: string) {
    // Login errors are shown via Ant Design message.toast, not inline form errors
    await expect(this.page.locator('.ant-message').first()).toContainText(message);
  }

  async expectOnLoginPage() {
    await expect(this.page).toHaveURL(/\/#\/login/);
  }
}
