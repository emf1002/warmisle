import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

/** Flush all pending microtasks (Promise resolutions + Vue scheduler). */
const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

const {
  mockGetSummary,
  mockGetExpenseChart,
  mockGetUpcomingTodos,
  mockGetWishTrends,
  mockGetForumHot,
} = vi.hoisted(() => ({
  mockGetSummary: vi.fn(),
  mockGetExpenseChart: vi.fn(),
  mockGetUpcomingTodos: vi.fn(),
  mockGetWishTrends: vi.fn(),
  mockGetForumHot: vi.fn(),
}))

vi.mock('@/api/dashboard', () => ({
  getSummary: mockGetSummary,
  getExpenseChart: mockGetExpenseChart,
  getUpcomingTodos: mockGetUpcomingTodos,
  getWishTrends: mockGetWishTrends,
  getForumHot: mockGetForumHot,
}))

import Dashboard from '../Index.vue'

/**
 * Override the global `true` stubs (which swallow slot content because
 * renderStubDefaultSlot is false) with stubs that forward their slots.
 * Also add stubs for components not listed in the global setup.
 */
const antStubs: Record<string, any> = {
  // Already globally stubbed but need slot forwarding:
  'a-button': {
    template: '<button><slot /></button>',
  },
  'a-card': {
    props: ['title', 'bordered'],
    template: '<div><slot /></div>',
  },
  'a-tag': {
    props: ['color', 'size'],
    template: '<span><slot /></span>',
  },
  // Not globally stubbed — must be provided:
  'a-statistic': {
    props: ['value', 'title', 'precision'],
    template: '<div><slot name="prefix" /><span>{{ fmt() }}</span></div>',
    methods: {
      fmt() {
        const v = this.value ?? 0
        return this.precision != null
          ? Number(v).toFixed(Number(this.precision))
          : String(v)
      },
    },
  },
  'a-list': {
    props: ['dataSource', 'size'],
    template:
      '<div><template v-for="(item, i) in dataSource || []" :key="i"><slot name="renderItem" :item="item" :index="i" /></template></div>',
  },
  'a-list-item': {
    template: '<div><slot /><slot name="extra" /></div>',
  },
  'a-list-item-meta': {
    template:
      '<div><div><slot name="title" /></div><div><slot name="description" /></div></div>',
  },
}

const emptySummary = {
  code: 0,
  message: 'ok',
  data: { income: 0, expense: 0, balance: 0 },
}
const emptyList = { code: 0, message: 'ok', data: [] }

function createWrapper() {
  return mount(Dashboard, {
    global: {
      stubs: {
        EmptyState: true,
        ...antStubs,
      },
    },
  })
}

describe('Dashboard view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Default every API call to empty data so each test only overrides what it needs.
    mockGetSummary.mockResolvedValue(emptySummary)
    mockGetExpenseChart.mockResolvedValue(emptyList)
    mockGetUpcomingTodos.mockResolvedValue(emptyList)
    mockGetWishTrends.mockResolvedValue(emptyList)
    mockGetForumHot.mockResolvedValue(emptyList)
  })

  it('renders summary with snake_case data', async () => {
    mockGetSummary.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: { income: 10000, expense: 5000, balance: 5000 },
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // 10000 / 100 = 100.00, 5000 / 100 = 50.00
    expect(text).toContain('100.00')
    expect(text).toContain('50.00')
  })

  it('renders expense chart with snake_case keys', async () => {
    mockGetExpenseChart.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: [
        { category_id: 1, category_name: '餐饮', icon: '🍔', amount: 3000 },
        { category_id: 2, category_name: '交通', icon: '🚌', amount: 2000 },
      ],
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('🍔')
    expect(text).toContain('餐饮')
    expect(text).toContain('🚌')
    expect(text).toContain('交通')
  })

  it('renders upcoming todos with nested assignee', async () => {
    mockGetUpcomingTodos.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: [
        {
          id: 1,
          title: '买菜',
          priority: 'important',
          due_date: '2026-05-25',
          assignee: { id: 1, name: '管理员', avatar: '👩' },
        },
        {
          id: 2,
          title: '修水龙头',
          priority: 'urgent',
          due_date: '2026-05-22',
          assignee: { id: 2, name: '小明', avatar: '👦' },
        },
      ],
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('买菜')
    expect(text).toContain('管理员')
    expect(text).toContain('修水龙头')
    expect(text).toContain('小明')
  })

  it('renders wish trends with vote_count', async () => {
    mockGetWishTrends.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: [
        {
          id: 1,
          title: '买iPad',
          creator: { name: '管理员' },
          vote_count: 3,
          status: 'pending',
        },
        {
          id: 2,
          title: '去旅游',
          creator: { name: '小明' },
          vote_count: 7,
          status: 'agreed',
        },
      ],
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('买iPad')
    expect(text).toContain('管理员')
    expect(text).toContain('3 票')
    expect(text).toContain('去旅游')
    expect(text).toContain('小明')
    expect(text).toContain('7 票')
  })

  it('renders forum hot with created_at', async () => {
    mockGetForumHot.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: [
        {
          type: 'post',
          title: '',
          content: '讨论周末安排',
          creator: { name: '管理员' },
          created_at: '2026-05-23T10:00:00Z',
        },
      ],
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // title is empty so the template falls back to content
    expect(text).toContain('讨论周末安排')
    expect(text).toContain('管理员')
  })

  it('fails to render when PascalCase keys used for summary', async () => {
    mockGetSummary.mockResolvedValue({
      code: 0,
      message: 'ok',
      // PascalCase keys do not match the snake_case bindings in the template
      data: { Income: 10000, Expense: 5000, Balance: 5000 },
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // summary.income is undefined (key is "Income"), so undefined / 100 = NaN
    // The correct amounts must NOT appear.
    expect(text).not.toContain('100.00')
    expect(text).not.toContain('50.00')
  })
})
