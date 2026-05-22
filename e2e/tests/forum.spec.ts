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

  test('点赞动态', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.openCreatePost();
    await forum.fillPostContent('测试点赞');
    await forum.submitModal();
    await forum.likePost(0);
    // 点赞数应增加
  });

  test('论坛页面视觉回归', async ({ authenticated }) => {
    const { page } = authenticated;
    const forum = new ForumPage(page);
    await forum.goto();
    await forum.screenshot('forum-empty.png');
  });
});
