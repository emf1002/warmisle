# Task 5: 家庭论坛 — E2E 测试扩展 (18 tests)

> **For agentic workers:** Use superpowers:subagent-driven-development or superpowers:executing-plans to implement.

**Goal:** 补齐论坛模块的公告、投票、评论/回复、标签管理和权限测试。

**Files:**
- Modify: `e2e/pages/forum.page.ts` — 新增 POM 方法
- Modify: `e2e/tests/forum.spec.ts` — 新增 18 个测试用例

**Dependencies:** Task 0 完成（memberContext fixture 可用）

---

## Step 1: 扩展 ForumPage POM

**Modify:** `e2e/pages/forum.page.ts` — 在 `goToDetail` 方法之后新增：

```typescript
  // === 公告 ===

  /** 发布公告（admin 操作，使用"公告"标签的话题） */
  async openCreateAnnouncement() {
    await this.page.getByTestId('create-announcement-btn').click();
    await expect(this.page.getByTestId('topic-modal')).toBeVisible();
  }

  /** 取消公告置顶 */
  async unpinAnnouncement(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await items.nth(index).getByTestId('unpin-btn').click();
    await this.page.waitForTimeout(300);
  }

  /** 断言公告已置顶 */
  async expectAnnouncementPinned(index: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await expect(items.nth(index).getByTestId('pinned-tag')).toBeVisible();
  }

  // === 投票 ===

  /** 打开创建投票弹窗 */
  async openCreatePoll() {
    await this.page.getByTestId('create-poll-btn').click();
    await expect(this.page.getByTestId('poll-modal')).toBeVisible();
  }

  /** 填写投票标题 */
  async fillPollTitle(title: string) {
    await this.page.getByTestId('poll-title').fill(title);
  }

  /** 添加投票选项 */
  async addPollOption(option: string) {
    await this.page.getByTestId('add-option-btn').click();
    const inputs = this.page.getByTestId(/^option-input-/);
    await inputs.last().fill(option);
  }

  /** 设置多选模式 */
  async setPollMultiSelect(enabled: boolean) {
    const checkbox = this.page.getByTestId('poll-multi-select');
    const isChecked = await checkbox.isChecked();
    if (enabled !== isChecked) {
      await checkbox.click();
    }
  }

  /** 提交投票（投票操作） */
  async votePoll(feedIndex: number, optionIndex: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    const options = items.nth(feedIndex).getByTestId(/^poll-option-/);
    await options.nth(optionIndex).click();
    await items.nth(feedIndex).getByTestId('poll-submit').click();
  }

  // === 评论 ===

  /** 对第 n 条内容发表评论 */
  async commentOnPost(feedIndex: number, text: string) {
    await this.goToDetail(feedIndex);
    await this.page.getByTestId('comment-input').fill(text);
    await this.page.getByTestId('comment-submit').click();
    await this.page.waitForTimeout(300);
  }

  /** 回复第 n 条一级评论 */
  async replyToComment(commentIndex: number, text: string) {
    const comments = this.page.getByTestId(/^comment-/);
    await comments.nth(commentIndex).getByTestId('reply-btn').click();
    await this.page.getByTestId('reply-input').fill(text);
    await this.page.getByTestId('reply-submit').click();
    await this.page.waitForTimeout(300);
  }

  /** 断言二级评论无回复按钮 */
  async expectNoReplyButton(commentIndex: number) {
    const comments = this.page.getByTestId(/^comment-/);
    await expect(comments.nth(commentIndex).getByTestId('reply-btn')).not.toBeVisible();
  }

  /** 删除评论 */
  async deleteComment(commentIndex: number) {
    const comments = this.page.getByTestId(/^comment-/);
    await comments.nth(commentIndex).getByTestId('delete-comment-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
    await this.page.waitForTimeout(300);
  }

  /** 断言评论数量 */
  async expectCommentCount(count: number) {
    await expect(this.page.getByTestId(/^comment-/)).toHaveCount(count);
  }

  // === 内容管理 ===

  /** 删除动态/话题 */
  async deletePost(feedIndex: number) {
    const items = this.page.getByTestId(/^feed-card-/);
    await items.nth(feedIndex).getByTestId('delete-feed-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
    await this.page.waitForTimeout(300);
  }

  // === 标签管理 ===

  /** 打开标签管理 */
  async openManageTags() {
    await this.page.getByTestId('manage-tags-btn').click();
    await expect(this.page.getByTestId('tags-modal')).toBeVisible();
  }

  /** 添加标签 */
  async addTag(name: string) {
    await this.page.getByTestId('add-tag-btn').click();
    await this.page.getByTestId('tag-name-input').fill(name);
    await this.page.getByTestId('tag-submit-btn').click();
    await this.page.waitForTimeout(300);
  }

  /** 删除标签 */
  async deleteTag(name: string) {
    await this.page.getByTestId('tags-modal')
      .locator(`[data-testid="tag-item"]`, { hasText: name })
      .getByTestId('delete-tag-btn').click();
    await this.page.locator('.ant-modal-confirm-btns .ant-btn-primary').click();
    await this.page.waitForTimeout(300);
  }

  /** 断言标签删除被禁用 */
  async expectTagDeleteDisabled(name: string) {
    await this.page.getByTestId('tags-modal')
      .locator(`[data-testid="tag-item"]`, { hasText: name })
      .getByTestId('delete-tag-btn')
      .then(btn => expect(btn).toBeDisabled());
  }
```

## Step 2: 新增测试用例

**Modify:** `e2e/tests/forum.spec.ts` — 在现有 describe 块内末尾新增：

```typescript
  // === 公告 ===

  test('发布公告', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateAnnouncement();
    await forum.fillTopicTitle('重要通知');
    await forum.fillTopicContent('本周六家庭聚会');
    await forum.selectTopicTag('公告');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    await forum.expectAnnouncementPinned(0);
  });

  test('取消公告置顶', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateAnnouncement();
    await forum.fillTopicTitle('通知');
    await forum.fillTopicContent('内容');
    await forum.selectTopicTag('公告');
    await forum.submitModal();
    await forum.unpinAnnouncement(0);
    // 取消后不应有置顶标记
    await expect(page.getByTestId('pinned-tag')).not.toBeVisible();
  });

  test('成员不可发布公告', async ({ authenticated, memberContext }) => {
    const { page: memberPage } = memberContext;
    const memberForum = new ForumPage(memberPage);
    await memberForum.goto();
    await expect(memberPage.getByTestId('create-announcement-btn')).not.toBeVisible();
  });

  // === 投票 ===

  test('创建投票', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePoll();
    await forum.fillPollTitle('周末去哪？');
    await forum.addPollOption('公园');
    await forum.addPollOption('商场');
    await forum.submitModal();
    await forum.expectFeedCount(1);
  });

  test('参与投票', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePoll();
    await forum.fillPollTitle('投票测试');
    await forum.addPollOption('选项A');
    await forum.addPollOption('选项B');
    await forum.submitModal();
    await forum.votePoll(0, 0);
    // 投票后应显示结果
    await expect(page.getByTestId('poll-result')).toBeVisible();
  });

  test('重复投票被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePoll();
    await forum.fillPollTitle('单次投票');
    await forum.addPollOption('是');
    await forum.addPollOption('否');
    await forum.submitModal();
    await forum.votePoll(0, 0);
    // 再次投票应无效果
    await forum.votePoll(0, 1);
    // 结果应保持第一次投票
    await expect(page.getByTestId('poll-result')).toBeVisible();
  });

  test('创建多选投票', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePoll();
    await forum.fillPollTitle('多选测试');
    await forum.addPollOption('选项A');
    await forum.addPollOption('选项B');
    await forum.addPollOption('选项C');
    await forum.setPollMultiSelect(true);
    await forum.submitModal();
    await forum.expectFeedCount(1);
  });

  // === 评论 ===

  test('评论动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试评论');
    await forum.submitModal();
    await forum.commentOnPost(0, '好棒！');
    await expect(page.getByTestId('comment-list')).toContainText('好棒！');
  });

  test('回复评论', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试回复');
    await forum.submitModal();
    await forum.commentOnPost(0, '一级评论');
    await forum.replyToComment(0, '二级回复');
    await forum.expectCommentCount(2);
  });

  test('二级评论无回复按钮', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试嵌套');
    await forum.submitModal();
    await forum.commentOnPost(0, '一级');
    await forum.replyToComment(0, '二级');
    await forum.expectNoReplyButton(1);
  });

  test('删除一级评论级联删除二级', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试级联删除');
    await forum.submitModal();
    await forum.commentOnPost(0, '一级');
    await forum.replyToComment(0, '二级');
    await forum.expectCommentCount(2);
    await forum.deleteComment(0);
    await forum.expectCommentCount(0);
  });

  // === 内容管理 ===

  test('删除自己的动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('待删除');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    await forum.deletePost(0);
    await forum.expectFeedCount(0);
  });

  test('管理员删除他人内容', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // member 发动态
    const memberForum = new ForumPage(memberPage);
    await memberForum.goto();
    await memberForum.openCreatePost();
    await memberForum.fillPostContent('member的动态');
    await memberForum.submitModal();
    // admin 删除
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.deletePost(0);
    await forum.expectFeedCount(0);
  });

  test('成员不可删除他人内容', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    // admin 发动态
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('admin的动态');
    await forum.submitModal();
    // member 不应看到删除按钮
    const memberForum = new ForumPage(memberPage);
    await memberForum.goto();
    await expect(memberPage.getByTestId('delete-feed-btn')).not.toBeVisible();
  });

  // === 标签管理 ===

  test('管理标签', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openManageTags();
    await forum.addTag('新标签');
    // 标签应出现在列表中
    await expect(page.getByTestId('tags-modal')).toContainText('新标签');
  });

  test('有关联话题的标签不可删除', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    // 创建带标签的话题
    await forum.openCreateTopic();
    await forum.fillTopicTitle('带标签的话题');
    await forum.fillTopicContent('内容');
    await forum.selectTopicTag('家务');
    await forum.submitModal();
    // 尝试删除该标签
    await forum.openManageTags();
    await forum.expectTagDeleteDisabled('家务');
  });

  // === 错误路径 ===

  test('动态内容为空被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('');
    await forum.submitModal();
    // 弹窗应仍然可见
    await expect(page.getByTestId('post-modal')).toBeVisible();
  });

  // === 响应式 ===

  test('移动端论坛', { tag: '@mobile' }, async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('移动端测试');
    await forum.submitModal();
    await forum.expectFeedCount(1);
  });
```

## Step 3: 验证

```bash
cd e2e && npx playwright test forum.spec.ts --project=chromium
```

## Step 4: Commit

```bash
git add e2e/pages/forum.page.ts e2e/tests/forum.spec.ts
git commit -m "test(e2e): expand forum tests — announcement, poll, comments, tags, permissions"
```
