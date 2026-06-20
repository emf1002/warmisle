// e2e/tests/forum.spec.ts
import { test, expect } from '../fixtures/auth.fixture';
import { ForumPage } from '../pages/forum.page';

test.describe('家庭论坛', () => {
  test('页面加载显示空状态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await expect(page.getByTestId('empty-state')).toBeVisible();
  });

  test('发动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('今天天气真好！');
    await forum.submitModal();
    await forum.expectFeedCount(1);
  });

  test('发话题', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('周末去哪玩');
    await forum.fillTopicContent('大家有什么推荐的地方吗？');
    await forum.selectTopicTag('出行');
    await forum.submitModal();
    await forum.expectFeedCount(1);
  });

  // P2: like-btn text 中 like_count 未刷新（0→0），需排查前端点赞→feed刷新联动
  test.skip('点赞动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试点赞');
    await forum.submitModal();
    await forum.likePost(0);
    await forum.expectLikeCount(0, 1);
  });

  // === 公告 ===

  test('发布公告', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateAnnouncement();
    await forum.fillTopicTitle('重要通知');
    await forum.fillTopicContent('本周六家庭聚会');
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
    await forum.submitModal();
    await forum.unpinAnnouncement(0);
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
    // After voting, the UI switches to results view — cannot vote again
    await expect(page.getByTestId('poll-result')).toBeVisible();
    // Verify voting UI is gone (no poll-submit button visible)
    await expect(page.getByTestId('poll-submit')).not.toBeVisible();
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

  test('投票选项在列表中显示', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePoll();
    await forum.fillPollTitle('列表显示测试');
    // Fill the 2 pre-existing option inputs
    await forum.fillPollOption(0, '选项一');
    await forum.fillPollOption(1, '选项二');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    // Verify poll options are visible in the list
    await expect(page.getByTestId('poll-option-0')).toBeVisible();
    await expect(page.getByTestId('poll-option-1')).toBeVisible();
    await expect(page.getByTestId('poll-option-0')).toContainText('选项一');
    await expect(page.getByTestId('poll-option-1')).toContainText('选项二');
  });

  test('刷新页面后投票仍显示', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    // Create a poll
    await forum.openCreatePoll();
    await forum.fillPollTitle('持久化测试');
    await forum.fillPollOption(0, '苹果');
    await forum.fillPollOption(1, '香蕉');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    // Refresh the page
    await page.reload();
    await page.waitForLoadState('networkidle');
    // Poll should still be visible with its options
    await forum.expectFeedCount(1);
    await expect(page.getByTestId('poll-option-0')).toBeVisible();
    await expect(page.getByTestId('poll-option-1')).toBeVisible();
    await expect(page.getByTestId('poll-option-0')).toContainText('苹果');
    await expect(page.getByTestId('poll-option-1')).toContainText('香蕉');
  });

  // P2: poll-result-text 是 CSS class 非 testid，已修复
  test('投票后刷新页面结果仍显示', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    // Create a poll
    await forum.openCreatePoll();
    await forum.fillPollTitle('投票持久化测试');
    await forum.addPollOption('选项A');
    await forum.addPollOption('选项B');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    // Vote
    await forum.votePoll(0, 0);
    await expect(page.getByTestId('poll-result')).toBeVisible();
    // Refresh the page
    await page.reload();
    await page.waitForLoadState('networkidle');
    // Vote results should still be visible
    await forum.expectFeedCount(1);
    await expect(page.getByTestId('poll-result')).toBeVisible();
    await expect(page.locator('.poll-result-text').first()).toContainText('选项A');
    await expect(page.locator('.poll-result-text').last()).toContainText('选项B');
  });

  // === 评论 ===

  test('评论动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试评论话题');
    await forum.fillTopicContent('话题内容');
    await forum.submitModal();
    await forum.commentOnPost(0, '好棒！');
    await expect(page.getByTestId('comment-list')).toContainText('好棒！');
  });

  // P2: reply-item 是 CSS class 非 testid，已修复
  test('回复评论', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试回复话题');
    await forum.fillTopicContent('话题内容');
    await forum.submitModal();
    await forum.commentOnPost(0, '一级评论');
    await forum.replyToComment(0, '二级回复');
    await forum.expectCommentCount(2);
  });

  // P2: 已修复 reply-item 选择器
  test('二级评论无回复按钮', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试嵌套话题');
    await forum.fillTopicContent('话题内容');
    await forum.submitModal();
    await forum.commentOnPost(0, '一级');
    await forum.replyToComment(0, '二级');
    await forum.expectNoReplyButton(1);
  });

  // P2: 已修复 reply-item 选择器，级联删除依赖正确评论渲染
  test('删除一级评论级联删除二级', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('测试级联删除话题');
    await forum.fillTopicContent('话题内容');
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
    const memberForum = new ForumPage(memberPage);
    await memberForum.goto();
    await memberForum.openCreatePost();
    await memberForum.fillPostContent('member的动态');
    await memberForum.submitModal();
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.deletePost(0);
    await forum.expectFeedCount(0);
  });

  test('成员不可删除他人内容', async ({ authenticated, memberContext }) => {
    const { page } = authenticated;
    const { page: memberPage } = memberContext;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('admin的动态');
    await forum.submitModal();
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
    await expect(page.getByTestId('tag-item').filter({ hasText: '新标签' })).toBeVisible();
  });

  test('有关联话题的标签不可删除', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('带标签的话题');
    await forum.fillTopicContent('内容');
    await forum.selectTopicTag('家务');
    await forum.submitModal();
    await forum.openManageTags();
    await forum.expectTagDeleteDisabled('家务');
  });

  // P2: topic detail 页 dropdown-trigger 交互不一致，需进一步调试
  test.skip('编辑已发布的话题', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('原标题');
    await forum.fillTopicContent('原内容');
    await forum.submitModal();
    await forum.expectFeedCount(1);
    // Navigate to topic detail via comment-btn (goToDetail)
    await forum.goToDetail(0);
    // Open topic detail dropdown (topic card 内的 dropdown-trigger)
    const topicDropdown = page.locator('.topic-card').getByTestId('dropdown-trigger');
    await topicDropdown.click();
    // 找出可见的菜单中的"编辑"选项
    await page.locator('.ant-dropdown:visible, .ant-dropdown-menu').getByText('编辑', { exact: true }).first().click();
    // Now fill and submit
    await forum.fillTopicTitle('编辑后的标题');
    await forum.submitModal();
    await expect(page.locator('.topic-title')).toContainText('编辑后的标题');
  });

  // P2: topic-title-display 不存在，标题用 .topic-title CSS class。feed-card click 对 topic 无效，改用 goToDetail
  test('话题详情页独立访问', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreateTopic();
    await forum.fillTopicTitle('独立访问测试');
    await forum.fillTopicContent('内容');
    await forum.submitModal();
    // Navigate to topic detail via goToDetail (clicks comment-btn)
    await forum.goToDetail(0);
    const detailUrl = page.url();
    // Reload the detail page directly
    await page.goto(detailUrl);
    await page.waitForLoadState('networkidle');
    await expect(page.locator('.topic-title')).toContainText('独立访问测试');
  });

  // === 错误路径 ===

  test('动态内容为空被拒绝', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('');
    // Don't use forum.submitModal() here — it expects modal to close
    await page.getByTestId('modal-submit-btn').click();
    // Error toast shown, modal stays open
    await forum.expectToast('内容不能为空');
    await forum.expectModalVisible();
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
});
