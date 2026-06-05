import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

const {
  mockGetLedgers,
  mockCreateLedger,
  mockUpdateLedger,
  mockDeleteLedger,
  mockGetCategories,
  mockGetMembers,
  mockModalConfirm,
} = vi.hoisted(() => ({
  mockGetLedgers: vi.fn(),
  mockCreateLedger: vi.fn(),
  mockUpdateLedger: vi.fn(),
  mockDeleteLedger: vi.fn(),
  mockGetCategories: vi.fn(),
  mockGetMembers: vi.fn(),
  mockModalConfirm: vi.fn(),
}))

vi.mock('@/api/ledger', () => ({
  getLedgers: mockGetLedgers,
  createLedger: mockCreateLedger,
  updateLedger: mockUpdateLedger,
  deleteLedger: mockDeleteLedger,
  getLedgerById: vi.fn(),
}))

vi.mock('@/api/category', () => ({
  getCategories: mockGetCategories,
  createCategory: vi.fn(),
  updateCategory: vi.fn(),
  deleteCategory: vi.fn(),
}))

vi.mock('@/api/member', () => ({
  getMembers: mockGetMembers,
  createMember: vi.fn(),
  updateMember: vi.fn(),
  deleteMember: vi.fn(),
  disableMember: vi.fn(),
  enableMember: vi.fn(),
  resetPassword: vi.fn(),
  getProfile: vi.fn(),
  updateProfile: vi.fn(),
  changePassword: vi.fn(),
}))

// Shared auth store mock (parses JWT from localStorage)
import '@/test/auth-mock'

// Mock Pinia stores used by the component
vi.mock('@/stores/categories', () => ({
  useCategoriesStore: () => ({
    categories: [
      { id: 1, type: 'expense', name: '餐饮', icon: '\uD83C\uDF71', sort_order: 1, preset: true },
      { id: 2, type: 'income', name: '工资', icon: '\uD83D\uDCB0', sort_order: 1, preset: true },
    ],
    loaded: true,
    fetchCategories: vi.fn().mockResolvedValue([]),
    reset: vi.fn(),
  }),
}))

vi.mock('@/stores/members', () => ({
  useMembersStore: () => ({
    members: [
      { id: 1, username: 'admin', name: '管理员', avatar: '\uD83D\uDC68', role: 'admin', status: 'active' },
    ],
    loaded: true,
    fetchMembers: vi.fn().mockResolvedValue([]),
    reset: vi.fn(),
  }),
}))

vi.mock('ant-design-vue', async (importOriginal) => {
  const actual = await importOriginal<any>()
  return {
    ...actual,
    Modal: {
      ...actual.Modal,
      confirm: mockModalConfirm,
    },
    message: {
      success: vi.fn(),
      error: vi.fn(),
      warning: vi.fn(),
      info: vi.fn(),
    },
  }
})

import LedgerIndex from '../Index.vue'

const stubs: Record<string, any> = {
  'a-button': { template: '<button><slot /></button>' },
  'a-tag': { template: '<span><slot /></span>' },
  'a-list': { template: '<div><slot /></div>' },
  'a-list-item': { template: '<div><slot /></div>' },
  'a-spin': { template: '<div>loading</div>' },
  'a-skeleton': { template: '<div>skeleton</div>' },
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
  'a-tabs': {
    props: ['activeKey'],
    template: '<div><slot /></div>',
    emits: ['update:activeKey'],
  },
  'a-tab-pane': {
    props: ['key', 'tab'],
    template: '<div><slot /></div>',
  },
  'a-pagination': { template: '<div />' },
  'a-dropdown': { template: '<div><slot /></div>' },
  'a-menu': { template: '<ul><slot /></ul>' },
  'a-menu-item': { template: '<li><slot /></li>' },
  'a-menu-divider': { template: '<hr />' },
  'a-range-picker': { template: '<div><slot /></div>' },
  'a-segmented': { template: '<div />' },
  'router-link': { template: '<a><slot /></a>' },
  'router-view': { template: '<div><slot /></div>' },
}

const mockCategories = {
  code: 0, message: 'ok', data: [
    { id: 1, type: 'expense', name: '餐饮', icon: '\uD83C\uDF71', sort_order: 1, preset: true },
    { id: 2, type: 'income', name: '工资', icon: '\uD83D\uDCB0', sort_order: 1, preset: true },
  ],
}

const mockMembersList = {
  code: 0, message: 'ok', data: [
    { id: 1, username: 'admin', name: '管理员', avatar: '\uD83D\uDC68', role: 'admin', status: 'active' },
  ],
}

const mockLedgersData = {
  code: 0, message: 'ok', data: {
    summary: { income: 10000, expense: 3550, balance: 6450 },
    groups: [{
      date: '2026-05-23', daily_total: -3550,
      items: [{
        id: 1, amount: 3550, note: '午餐', category_id: 1, creator_id: 1,
        occurred_at: '2026-05-23T12:00:00Z',
        category: { id: 1, name: '餐饮', icon: '\uD83C\uDF71', type: 'expense' },
        creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
      }],
    }],
    has_more: false, next_cursor: null,
  },
}

function createWrapper() {
  return mount(LedgerIndex, {
    global: { stubs },
  })
}

describe('Ledger view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
    localStorage.setItem('token', `header.${payload}.sig`)
    mockGetCategories.mockResolvedValue(mockCategories)
    mockGetMembers.mockResolvedValue(mockMembersList)
    mockGetLedgers.mockResolvedValue(mockLedgersData)
  })

  it('renders summary bar with formatted amounts', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('收入')
    expect(text).toContain('支出')
    expect(text).toContain('结余')
    // 10000/100 = 100.00, 3550/100 = 35.50, 6450/100 = 64.50
    expect(text).toContain('100.00')
    expect(text).toContain('35.50')
    expect(text).toContain('64.50')
  })

  it('renders ledger items with snake_case keys', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('\uD83C\uDF71') // bento icon
    expect(text).toContain('餐饮')
    expect(text).toContain('午餐')
    expect(text).toContain('35.50')
    expect(text).toContain('管理员')
    expect(text).toContain('\uD83D\uDC68') // man avatar
  })

  it('renders date group header', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // The date 2026-05-23 renders as "5月23日 周六"
    expect(text).toContain('5月23日 周六')
  })

  it('creates ledger and refreshes list with updated data', async () => {
    mockCreateLedger.mockResolvedValue({ code: 0, message: 'ok', data: { id: 99 } })

    const updatedData = {
      code: 0, message: 'ok', data: {
        summary: { income: 10000, expense: 4550, balance: 5450 },
        groups: [{
          date: '2026-05-23', daily_total: -4550,
          items: [{
            id: 99, amount: 1000, note: '新记录', category_id: 1, creator_id: 1,
            occurred_at: '2026-05-23T14:00:00Z',
            category: { id: 1, name: '餐饮', icon: '\uD83C\uDF71', type: 'expense' },
            creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
          }, {
            id: 1, amount: 3550, note: '午餐', category_id: 1, creator_id: 1,
            occurred_at: '2026-05-23T12:00:00Z',
            category: { id: 1, name: '餐饮', icon: '\uD83C\uDF71', type: 'expense' },
            creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
          }],
        }],
        has_more: false, next_cursor: null,
      },
    }

    mockGetLedgers.mockResolvedValueOnce(mockLedgersData)
    mockGetLedgers.mockResolvedValueOnce(updatedData)

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const vm = wrapper.vm as any
    vm.dialogOpen = true
    vm.form.category_id = 1
    vm.form.amount = 10
    vm.form.note = '新记录'
    vm.form.occurred_at = undefined
    await nextTick()

    await vm.handleSubmit()
    await flushPromises()
    await nextTick()

    expect(mockCreateLedger).toHaveBeenCalled()
    expect(mockGetLedgers).toHaveBeenCalledTimes(2)

    // Verify the new record appears in the rendered list
    const text = wrapper.text()
    expect(text).toContain('新记录')
    expect(text).toContain('10.00')
  })

  it('deletes ledger and refreshes list', async () => {
    mockDeleteLedger.mockResolvedValue({ code: 0, message: 'ok' })

    const updatedData = {
      ...mockLedgersData,
      data: {
        ...mockLedgersData.data,
        summary: { income: 10000, expense: 0, balance: 10000 },
        groups: [],
        has_more: false, next_cursor: null,
      },
    }
    mockGetLedgers.mockResolvedValueOnce(mockLedgersData)
    mockGetLedgers.mockResolvedValueOnce(updatedData)

    // Make Modal.confirm auto-call onOk
    mockModalConfirm.mockImplementation((opts: any) => {
      if (opts.onOk) opts.onOk()
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetLedgers).toHaveBeenCalledTimes(1)

    // Open edit dialog for existing record, then delete
    const vm = wrapper.vm as any
    vm.editingRecord = mockLedgersData.data.groups[0].items[0]
    vm.dialogOpen = true
    await nextTick()

    vm.confirmDelete()
    await flushPromises()
    await nextTick()

    // deleteLedger should be called, and getLedgers should refresh
    expect(mockDeleteLedger).toHaveBeenCalledWith(1)
    expect(mockGetLedgers).toHaveBeenCalledTimes(2)
  })
})
