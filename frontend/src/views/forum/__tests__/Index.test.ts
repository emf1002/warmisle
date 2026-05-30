import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { nextTick } from 'vue'

const {
  mockGetFeed,
  mockCreatePost,
  mockUpdatePost,
  mockDeletePost,
  mockCreateTopic,
  mockUpdateTopic,
  mockDeleteTopic,
  mockTogglePin,
  mockGetTags,
} = vi.hoisted(() => ({
  mockGetFeed: vi.fn(),
  mockCreatePost: vi.fn(),
  mockUpdatePost: vi.fn(),
  mockDeletePost: vi.fn(),
  mockCreateTopic: vi.fn(),
  mockUpdateTopic: vi.fn(),
  mockDeleteTopic: vi.fn(),
  mockTogglePin: vi.fn(),
  mockGetTags: vi.fn(),
}))

vi.mock('@/api/forum', () => ({
  getFeed: mockGetFeed,
  createPost: mockCreatePost,
  updatePost: mockUpdatePost,
  deletePost: mockDeletePost,
  createTopic: mockCreateTopic,
  updateTopic: mockUpdateTopic,
  deleteTopic: mockDeleteTopic,
  togglePin: mockTogglePin,
  getTags: mockGetTags,
}))

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => ({ params: {} }),
  createRouter: vi.fn(),
  createWebHashHistory: vi.fn(),
}))

vi.mock('ant-design-vue', () => ({
  message: { success: vi.fn(), error: vi.fn(), info: vi.fn() },
  Modal: { confirm: vi.fn() },
}))

vi.mock('@/utils/request', () => ({
  default: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn(), interceptors: { request: { use: vi.fn() }, response: { use: vi.fn() } } },
}))

// Shared auth store mock (parses JWT from localStorage)
import '@/test-utils/auth-mock'

import Forum from '../Index.vue'

const stubs = {
  'a-button': { template: '<button><slot /></button>' },
  'a-tag': { template: '<span><slot /></span>' },
  'a-list': { template: '<div><slot /></div>' },
  'a-list-item': { template: '<div><slot /></div>' },
  'a-spin': { template: '<div>loading</div>' },
  'a-modal': { template: '<div v-bind="$attrs"><slot /></div>' },
  'a-form': { template: '<form><slot /></form>' },
  'a-form-item': { template: '<div><slot /></div>' },
  'a-input': { template: '<input />' },
  'a-textarea': { template: '<textarea />' },
  'a-select': { template: '<div><slot /></div>' },
  'a-select-option': { template: '<option><slot /></option>' },
  'a-dropdown': { template: '<div><slot /></div>' },
  'a-menu': { template: '<ul><slot /></ul>' },
  'a-menu-item': { template: '<li><slot /></li>' },
  'a-menu-divider': { template: '<hr />' },
  'a-drawer': { template: '<div><slot /></div>' },
  'a-pagination': { template: '<div />' },
  'router-link': { template: '<a><slot /></a>' },
  'router-view': { template: '<div><slot /></div>' },
}

const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))

function makeWrapper() {
  return mount(Forum, {
    global: { stubs },
  })
}

describe('Forum Index', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.setItem('token', `header.${payload}.sig`)

    mockGetTags.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: [
        { id: 1, name: '讨论', preset: true },
        { id: 2, name: '公告', preset: true },
      ],
    })

    mockGetFeed.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: {
        pinned: [
          {
            type: 'topic',
            id: 10,
            title: '置顶公告',
            content: '这是置顶内容',
            creator: { id: 1, name: '管理员', avatar: '👨' },
            tag: { id: 2, name: '公告', preset: true },
            is_pinned: true,
            like_count: 5,
            comment_count: 3,
            created_at: '2026-05-22T10:00:00Z',
          },
        ],
        items: [
          {
            type: 'post',
            id: 1,
            title: '',
            content: '今天天气真好！',
            creator: { id: 1, name: '管理员', avatar: '👨' },
            tag: null,
            is_pinned: false,
            is_liked: false,
            like_count: 2,
            comment_count: 1,
            created_at: '2026-05-23T10:00:00Z',
          },
          {
            type: 'topic',
            id: 2,
            title: '周末计划',
            content: '大家有什么计划？',
            creator: { id: 2, name: '成员', avatar: '👩' },
            tag: { id: 1, name: '讨论', preset: true },
            is_pinned: false,
            like_count: 0,
            comment_count: 0,
            created_at: '2026-05-23T09:00:00Z',
          },
        ],
        total: 3,
      },
    })
  })

  it('renders pinned items', async () => {
    const wrapper = makeWrapper()
    await flushPromises()

    const pinnedSection = wrapper.find('.pinned-section')
    expect(pinnedSection.exists()).toBe(true)
    expect(pinnedSection.text()).toContain('📌')
    expect(pinnedSection.text()).toContain('置顶公告')
    expect(pinnedSection.text()).toContain('这是置顶内容')
    expect(pinnedSection.text()).toContain('#公告')
  })

  it('renders feed posts', async () => {
    const wrapper = makeWrapper()
    await flushPromises()

    const feedList = wrapper.find('[data-testid="feed-list"]')
    expect(feedList.exists()).toBe(true)
    expect(feedList.text()).toContain('今天天气真好！')
    expect(feedList.text()).toContain('管理员')
    expect(feedList.text()).toContain('👨')
  })

  it('renders feed topics', async () => {
    const wrapper = makeWrapper()
    await flushPromises()

    const feedList = wrapper.find('[data-testid="feed-list"]')
    expect(feedList.text()).toContain('周末计划')
    expect(feedList.text()).toContain('大家有什么计划？')
    expect(feedList.text()).toContain('#讨论')
    expect(feedList.text()).toContain('👩')
    expect(feedList.text()).toContain('成员')
  })

  it('renders like and comment counts', async () => {
    const wrapper = makeWrapper()
    await flushPromises()

    const html = wrapper.text()
    // Pinned item: like_count=5, comment_count=3
    expect(html).toContain('5')
    expect(html).toContain('3')
    // Post item: like_count=2, comment_count=1
    expect(html).toContain('2')
    expect(html).toContain('1')
  })

  it('renders creator info', async () => {
    const wrapper = makeWrapper()
    await flushPromises()

    const html = wrapper.text()
    // Post creator
    expect(html).toContain('👨')
    expect(html).toContain('管理员')
    // Topic creator
    expect(html).toContain('👩')
    expect(html).toContain('成员')
  })

  it('creates post and refreshes feed', async () => {
    mockCreatePost.mockResolvedValue({ code: 0, message: 'ok', data: { id: 99 } })

    const updatedFeedData = {
      code: 0,
      message: 'ok',
      data: {
        pinned: [],
        items: [
          {
            type: 'post',
            id: 99,
            title: '',
            content: '新发布的动态',
            creator: { id: 1, name: '管理员', avatar: '👨' },
            tag: null,
            is_pinned: false,
            like_count: 0,
            comment_count: 0,
            created_at: '2026-05-23T12:00:00Z',
          },
        ],
        total: 1,
      },
    }
    mockGetFeed.mockResolvedValueOnce({
      code: 0,
      message: 'ok',
      data: {
        pinned: [{ type: 'topic', id: 10, title: '置顶公告', content: '这是置顶内容', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: { id: 2, name: '公告', preset: true }, is_pinned: true, like_count: 5, comment_count: 3, created_at: '2026-05-22T10:00:00Z' }],
        items: [
          { type: 'post', id: 1, title: '', content: '今天天气真好！', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: null, is_pinned: false, is_liked: false, like_count: 2, comment_count: 1, created_at: '2026-05-23T10:00:00Z' },
          { type: 'topic', id: 2, title: '周末计划', content: '大家有什么计划？', creator: { id: 2, name: '成员', avatar: '👩' }, tag: { id: 1, name: '讨论', preset: true }, is_pinned: false, like_count: 0, comment_count: 0, created_at: '2026-05-23T09:00:00Z' },
        ],
        total: 3,
      },
    }).mockResolvedValueOnce(updatedFeedData)

    const wrapper = makeWrapper()
    await flushPromises()

    const initialCallCount = mockGetFeed.mock.calls.length

    // Open create post dialog
    await wrapper.find('[data-testid="create-post-btn"]').trigger('click')
    await nextTick()

    // Set post content via component data
    const vm = wrapper.vm as any
    vm.postForm.content = '新发布的动态'
    vm.postDialogOpen = true
    await nextTick()

    // Submit the post
    await vm.handlePostSubmit()
    await flushPromises()

    expect(mockCreatePost).toHaveBeenCalledWith({ content: '新发布的动态' })
    expect(mockGetFeed).toHaveBeenCalledTimes(initialCallCount + 1)
  })

  it('creates topic and refreshes feed', async () => {
    mockCreateTopic.mockResolvedValue({ code: 0, message: 'ok', data: { id: 88 } })

    mockGetFeed.mockResolvedValueOnce({
      code: 0,
      message: 'ok',
      data: {
        pinned: [{ type: 'topic', id: 10, title: '置顶公告', content: '这是置顶内容', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: { id: 2, name: '公告', preset: true }, is_pinned: true, like_count: 5, comment_count: 3, created_at: '2026-05-22T10:00:00Z' }],
        items: [
          { type: 'post', id: 1, title: '', content: '今天天气真好！', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: null, is_pinned: false, is_liked: false, like_count: 2, comment_count: 1, created_at: '2026-05-23T10:00:00Z' },
          { type: 'topic', id: 2, title: '周末计划', content: '大家有什么计划？', creator: { id: 2, name: '成员', avatar: '👩' }, tag: { id: 1, name: '讨论', preset: true }, is_pinned: false, like_count: 0, comment_count: 0, created_at: '2026-05-23T09:00:00Z' },
        ],
        total: 3,
      },
    }).mockResolvedValueOnce({
      code: 0,
      message: 'ok',
      data: {
        pinned: [],
        items: [
          {
            type: 'topic',
            id: 88,
            title: '新话题',
            content: '话题内容',
            creator: { id: 1, name: '管理员', avatar: '👨' },
            tag: { id: 1, name: '讨论', preset: true },
            is_pinned: false,
            like_count: 0,
            comment_count: 0,
            created_at: '2026-05-23T12:00:00Z',
          },
        ],
        total: 1,
      },
    })

    const wrapper = makeWrapper()
    await flushPromises()

    const initialCallCount = mockGetFeed.mock.calls.length

    // Open create topic dialog
    await wrapper.find('[data-testid="create-topic-btn"]').trigger('click')
    await nextTick()

    // Set topic form via component data
    const vm = wrapper.vm as any
    vm.topicForm.title = '新话题'
    vm.topicForm.content = '话题内容'
    vm.topicForm.tag_id = 1
    vm.topicDialogOpen = true
    await nextTick()

    // Submit the topic
    await vm.handleTopicSubmit()
    await flushPromises()

    expect(mockCreateTopic).toHaveBeenCalledWith({
      title: '新话题',
      content: '话题内容',
      tag_id: 1,
    })
    expect(mockGetFeed).toHaveBeenCalledTimes(initialCallCount + 1)
  })

  it('toggles pin and refreshes feed', async () => {
    mockTogglePin.mockResolvedValue({ code: 0, message: 'ok', data: null })

    mockGetFeed.mockResolvedValueOnce({
      code: 0,
      message: 'ok',
      data: {
        pinned: [{ type: 'topic', id: 10, title: '置顶公告', content: '这是置顶内容', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: { id: 2, name: '公告', preset: true }, is_pinned: true, like_count: 5, comment_count: 3, created_at: '2026-05-22T10:00:00Z' }],
        items: [
          { type: 'post', id: 1, title: '', content: '今天天气真好！', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: null, is_pinned: false, is_liked: false, like_count: 2, comment_count: 1, created_at: '2026-05-23T10:00:00Z' },
          { type: 'topic', id: 2, title: '周末计划', content: '大家有什么计划？', creator: { id: 2, name: '成员', avatar: '👩' }, tag: { id: 1, name: '讨论', preset: true }, is_pinned: false, like_count: 0, comment_count: 0, created_at: '2026-05-23T09:00:00Z' },
        ],
        total: 3,
      },
    }).mockResolvedValueOnce({
      code: 0,
      message: 'ok',
      data: {
        pinned: [],
        items: [
          { type: 'post', id: 1, title: '', content: '今天天气真好！', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: null, is_pinned: false, is_liked: false, like_count: 2, comment_count: 1, created_at: '2026-05-23T10:00:00Z' },
          { type: 'topic', id: 10, title: '置顶公告', content: '这是置顶内容', creator: { id: 1, name: '管理员', avatar: '👨' }, tag: { id: 2, name: '公告', preset: true }, is_pinned: false, like_count: 5, comment_count: 3, created_at: '2026-05-22T10:00:00Z' },
          { type: 'topic', id: 2, title: '周末计划', content: '大家有什么计划？', creator: { id: 2, name: '成员', avatar: '👩' }, tag: { id: 1, name: '讨论', preset: true }, is_pinned: false, like_count: 0, comment_count: 0, created_at: '2026-05-23T09:00:00Z' },
        ],
        total: 3,
      },
    })

    const wrapper = makeWrapper()
    await flushPromises()

    const initialCallCount = mockGetFeed.mock.calls.length

    // Call handleTogglePin directly on the pinned topic item
    const vm = wrapper.vm as any
    const pinnedItem = {
      type: 'topic',
      id: 10,
      title: '置顶公告',
      content: '这是置顶内容',
      creator: { id: 1, name: '管理员', avatar: '👨' },
      tag: { id: 2, name: '公告', preset: true },
      is_pinned: true,
      like_count: 5,
      comment_count: 3,
      created_at: '2026-05-22T10:00:00Z',
    }

    await vm.handleTogglePin(pinnedItem)
    await flushPromises()

    expect(mockTogglePin).toHaveBeenCalledWith(10)
    expect(mockGetFeed).toHaveBeenCalledTimes(initialCallCount + 1)
  })
})
