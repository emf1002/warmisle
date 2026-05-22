import { type Page, expect } from '@playwright/test';
import { BasePage } from './base.page';

export class ForumPage extends BasePage {
  constructor(page: Page) {
    super(page);
  }

  async goto() {
    await this.navigate('/forum');
  }

  async expectOnForum() {
    await expect(this.page).toHaveURL(/\/#\/forum/);
  }

  async openCreatePost() {
    await this.page.getByTestId('create-post-btn').click();
    await expect(this.page.getByTestId('post-modal')).toBeVisible();
  }

  async openCreateTopic() {
    await this.page.getByTestId('create-topic-btn').click();
    await expect(this.page.getByTestId('topic-modal')).toBeVisible();
  }

  async fillPostContent(content: string) {
    await this.page.getByTestId('post-content').fill(content);
  }

  async fillTopicTitle(title: string) {
    await this.page.getByTestId('topic-title').fill(title);
  }

  async fillTopicContent(content: string) {
    await this.page.getByTestId('topic-content').fill(content);
  }

  async selectTopicTag(tag: string) {
    await this.page.getByTestId('topic-tag').click();
    await this.page.locator('.ant-select-item-option', { hasText: tag }).click();
  }

  async submitModal() {
    await this.page.locator('.ant-modal-footer .ant-btn-primary').click();
  }

  async expectFeedCount(count: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await expect(items).toHaveCount(count);
  }

  async likePost(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await items.nth(index).getByTestId('like-btn').click();
  }

  async goToDetail(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await items.nth(index).getByTestId('comment-btn').click();
  }
}
