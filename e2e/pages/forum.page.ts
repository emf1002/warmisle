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
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async openCreateTopic() {
    await this.page.getByTestId('create-topic-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
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
    // Wait for modal close animation to finish
    await this.page.locator('.ant-modal-wrap:visible').waitFor({ state: 'hidden', timeout: 5000 }).catch(() => {});
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
    await this.page.waitForURL(/\/#\/forum\/topic\//);
    await this.page.waitForLoadState('networkidle');
  }

  // === 公告 ===

  async openCreateAnnouncement() {
    await this.page.getByTestId('create-announcement-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async unpinAnnouncement(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    // Open the dropdown menu first by clicking the ⋯ trigger
    await items.nth(index).locator('.ant-dropdown-trigger').click();
    await this.page.locator('.ant-dropdown:visible').getByTestId('unpin-btn').click();
    await this.page.waitForTimeout(300);
  }

  async expectAnnouncementPinned(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await expect(items.nth(index).getByTestId('pinned-tag')).toBeVisible();
  }

  // === 投票 ===

  async openCreatePoll() {
    await this.page.getByTestId('create-poll-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async fillPollTitle(title: string) {
    await this.page.getByTestId('poll-title').fill(title);
  }

  async addPollOption(option: string) {
    await this.page.getByTestId('add-option-btn').click();
    const inputs = this.page.getByTestId(/^option-input-/);
    await inputs.last().fill(option);
  }

  async setPollMultiSelect(enabled: boolean) {
    const checkbox = this.page.getByTestId('poll-multi-select');
    const isChecked = await checkbox.isChecked();
    if (enabled !== isChecked) {
      await checkbox.click();
    }
  }

  async votePoll(feedIndex: number, optionIndex: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    const options = items.nth(feedIndex).getByTestId(/^poll-option-/);
    await options.nth(optionIndex).click();
    await items.nth(feedIndex).getByTestId('poll-submit').click();
  }

  // === 评论 ===

  async commentOnPost(feedIndex: number, text: string) {
    await this.goToDetail(feedIndex);
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await this.page.waitForTimeout(300);
  }

  async replyToComment(commentIndex: number, text: string) {
    const comments = this.page.locator('.comment-item, .reply-item');
    await comments.nth(commentIndex).getByTestId('reply-btn').click();
    await this.page.getByTestId('reply-input').fill(text);
    await this.page.getByTestId('reply-submit').click();
    await this.page.waitForTimeout(300);
  }

  async expectNoReplyButton(commentIndex: number) {
    const comments = this.page.locator('.comment-item, .reply-item');
    await expect(comments.nth(commentIndex).getByTestId('reply-btn')).not.toBeVisible();
  }

  async deleteComment(commentIndex: number) {
    const comments = this.page.locator('.comment-item, .reply-item');
    // Use .first() because a parent comment may contain a nested reply's delete button
    await comments.nth(commentIndex).getByTestId('delete-comment-btn').first().click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
    await this.page.waitForTimeout(300);
  }

  async expectCommentCount(count: number) {
    await expect(this.page.locator('.comment-item, .reply-item')).toHaveCount(count);
  }

  // === 内容管理 ===

  async deletePost(feedIndex: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    // Open dropdown menu first
    await items.nth(feedIndex).locator('.ant-dropdown-trigger').click();
    await this.page.locator('.ant-dropdown:visible').getByTestId('delete-feed-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
    await this.page.waitForTimeout(300);
  }

  // === 标签管理 ===

  async openManageTags() {
    await this.page.getByTestId('manage-tags-btn').click();
    await expect(this.page.locator('.ant-modal-wrap:visible')).toBeVisible();
  }

  async addTag(name: string) {
    await this.page.getByTestId('add-tag-btn').click();
    await this.page.getByTestId('tag-name-input').fill(name);
    await this.page.getByTestId('tag-submit-btn').click();
    await this.page.waitForTimeout(300);
  }

  async deleteTag(name: string) {
    await this.page.locator('.ant-modal-wrap:visible')
      .locator(`[data-testid="tag-item"]`, { hasText: name })
      .getByTestId('delete-tag-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn:last-child').click();
    await this.page.waitForTimeout(300);
  }

  async expectTagDeleteDisabled(name: string) {
    const btn = this.page.locator('.ant-modal-wrap:visible')
      .locator(`[data-testid="tag-item"]`, { hasText: name })
      .getByTestId('delete-tag-btn');
    await expect(btn).toBeDisabled();
  }
}
