import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

const {
  mockGetWishList,
  mockCreateWish,
  mockUpdateWish,
  mockDeleteWish,
  mockPromoteWish,
  mockUpdateWishStatus,
  mockVoteWish,
  mockUnvoteWish,
} = vi.hoisted(() => ({
  mockGetWishList: vi.fn(),
  mockCreateWish: vi.fn(),
  mockUpdateWish: vi.fn(),
  mockDeleteWish: vi.fn(),
  mockPromoteWish: vi.fn(),
  mockUpdateWishStatus: vi.fn(),
  mockVoteWish: vi.fn(),
  mockUnvoteWish: vi.fn(),
}))

vi.mock('@/api/wish', () => ({
  getWishList: mockGetWishList,
  createWish: mockCreateWish,
  updateWish: mockUpdateWish,
  deleteWish: mockDeleteWish,
  promoteWish: mockPromoteWish,
  updateWishStatus: mockUpdateWishStatus,
  voteWish: mockVoteWish,
  unvoteWish: mockUnvoteWish,
}))

import WishIndex from '../Index.vue'

const stubs: Record<string, any> = {
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
  'a-input-number': { template: '<input />' },
  'a-select': { template: '<div><slot /></div>' },
  'a-select-option': { template: '<option><slot /></option>' },
  'a-select-opt-group': { template: '<optgroup><slot /></optgroup>' },
  'a-checkbox': { template: '<input type="checkbox" />' },
  'a-date-picker': { template: '<input />' },
  'a-pagination': { template: '<div />' },
  'a-dropdown': { template: '<div><slot /></div>' },
  'a-menu': { template: '<ul><slot /></ul>' },
  'a-menu-item': { template: '<li><slot /></li>' },
  'a-menu-divider': { template: '<hr />' },
  'a-segmented': { template: '<div />' },
  'router-link': { template: '<a><slot /></a>' },
  'router-view': { template: '<div><slot /></div>' },
  'EmptyState': { template: '<div><slot /></div>' },
}

const mockWishData = {
  code: 0, message: 'ok', data: {
    list: [
      {
        id: 1, title: '买iPad', description: '学习用', category: 'item', amount: 500000,
        priority: 'normal', type: 'personal', status: 'pending', creator_id: 1, vote_count: 3,
        creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
      },
      {
        id: 2, title: '去旅行', description: '全家旅行', category: 'travel', amount: 1000000,
        priority: 'important', type: 'family', status: 'agreed', creator_id: 1, vote_count: 5,
        creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
      },
    ],
    total: 2, page: 1, page_size: 12,
  },
}

function createWrapper() {
  return mount(WishIndex, {
    global: { stubs },
  })
}

describe('Wish view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
    localStorage.setItem('token', `header.${payload}.sig`)
    mockGetWishList.mockResolvedValue(mockWishData)
  })

  it('renders wish cards with snake_case keys', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('买iPad')
    expect(text).toContain('学习用')
    // 500000 / 100 = 5000.00
    expect(text).toContain('5000.00')
    expect(text).toContain('\uD83D\uDC68 管理员')
  })

  it('renders priority and status tags', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // First wish: priority=normal -> "普通", status=pending -> "待定"
    expect(text).toContain('普通')
    expect(text).toContain('待定')
    // Second wish: priority=important -> "重要", status=agreed -> "已同意"
    expect(text).toContain('重要')
    expect(text).toContain('已同意')
  })

  it('renders category labels', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // category=item -> "🛍️ 物品", category=travel -> "✈️ 旅行"
    expect(text).toContain('\uD83D\uDECD\uFE0F 物品')
    expect(text).toContain('\u2708\uFE0F 旅行')
  })

  it('renders vote counts in family mode', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    // Vote buttons only render when activeType === 'family'
    const vm = wrapper.vm as any
    vm.activeType = 'family'
    await nextTick()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('3')
    expect(text).toContain('5')
  })

  it('creates wish and refreshes list', async () => {
    mockCreateWish.mockResolvedValue({ code: 0, message: 'ok', data: { id: 99 } })

    const updatedData = {
      ...mockWishData,
      data: { ...mockWishData.data, total: 3 },
    }
    mockGetWishList.mockResolvedValueOnce(mockWishData)
    mockGetWishList.mockResolvedValueOnce(updatedData)

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetWishList).toHaveBeenCalledTimes(1)

    // Fill form and submit
    const vm = wrapper.vm as any
    vm.dialogOpen = true
    vm.form.title = '新愿望'
    vm.form.description = '描述'
    vm.form.category = 'other'
    vm.form.priority = 'normal'
    await nextTick()

    await vm.handleSubmit()
    await flushPromises()
    await nextTick()

    expect(mockCreateWish).toHaveBeenCalled()
    expect(mockGetWishList).toHaveBeenCalledTimes(2)
  })

  it('votes wish and refreshes list', async () => {
    mockVoteWish.mockResolvedValue({ code: 0, message: 'ok' })

    const updatedData = {
      ...mockWishData,
      data: {
        ...mockWishData.data,
        list: mockWishData.data.list.map(w =>
          w.id === 2 ? { ...w, vote_count: 6 } : w
        ),
      },
    }
    mockGetWishList.mockResolvedValueOnce(mockWishData)
    mockGetWishList.mockResolvedValueOnce(updatedData)

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetWishList).toHaveBeenCalledTimes(1)

    // Call handleVote on the second wish (family type)
    const vm = wrapper.vm as any
    await vm.handleVote(mockWishData.data.list[1])
    await flushPromises()
    await nextTick()

    expect(mockVoteWish).toHaveBeenCalledWith(2)
    expect(mockGetWishList).toHaveBeenCalledTimes(2)
  })
})
