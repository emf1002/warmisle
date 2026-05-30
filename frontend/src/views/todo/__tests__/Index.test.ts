import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

const {
  mockGetTodoList,
  mockCreateTodo,
  mockUpdateTodo,
  mockDeleteTodo,
  mockToggleTodo,
  mockClaimTodo,
  mockGetMembers,
} = vi.hoisted(() => ({
  mockGetTodoList: vi.fn(),
  mockCreateTodo: vi.fn(),
  mockUpdateTodo: vi.fn(),
  mockDeleteTodo: vi.fn(),
  mockToggleTodo: vi.fn(),
  mockClaimTodo: vi.fn(),
  mockGetMembers: vi.fn(),
}))

vi.mock('@/api/todo', () => ({
  getTodoList: mockGetTodoList,
  createTodo: mockCreateTodo,
  updateTodo: mockUpdateTodo,
  deleteTodo: mockDeleteTodo,
  toggleTodo: mockToggleTodo,
  claimTodo: mockClaimTodo,
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
import '@/test-utils/auth-mock'

import TodoIndex from '../Index.vue'

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

const mockMembersList = {
  code: 0, message: 'ok', data: [
    { id: 1, username: 'admin', name: '管理员', avatar: '\uD83D\uDC68', role: 'admin', status: 'active' },
    { id: 2, username: 'member', name: '成员', avatar: '\uD83D\uDC69', role: 'member', status: 'active' },
  ],
}

const mockTodoData = {
  code: 0, message: 'ok', data: {
    list: [
      {
        id: 1, title: '买菜', description: '去超市买菜', priority: 'important', status: 'pending',
        assignee_id: 2, creator_id: 1, due_date: '2026-05-25', completed_at: null, created_at: '2026-05-23',
        assignee: { id: 2, name: '成员', avatar: '\uD83D\uDC69' },
        creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
      },
      {
        id: 2, title: '做饭', description: '', priority: 'normal', status: 'completed',
        assignee_id: null, creator_id: 1, due_date: null, completed_at: '2026-05-23', created_at: '2026-05-22',
        assignee: null,
        creator: { id: 1, name: '管理员', avatar: '\uD83D\uDC68' },
      },
    ],
    total: 2, page: 1, page_size: 20,
  },
}

function createWrapper() {
  return mount(TodoIndex, {
    global: { stubs },
  })
}

describe('Todo view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
    localStorage.setItem('token', `header.${payload}.sig`)
    mockGetMembers.mockResolvedValue(mockMembersList)
    mockGetTodoList.mockResolvedValue(mockTodoData)
  })

  it('renders todo items with snake_case keys', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('买菜')
    expect(text).toContain('去超市买菜')
    expect(text).toContain('重要')
    expect(text).toContain('成员')
    expect(text).toContain('\uD83D\uDC69') // assignee avatar
  })

  it('renders completed state with title-done class', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    // The completed todo (id=2, title="做饭") should have the title-done class
    const doneTitle = wrapper.find('.title-done')
    expect(doneTitle.exists()).toBe(true)
    expect(doneTitle.text()).toContain('做饭')
  })

  it('renders assignee and creator avatars', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // Creator avatar for both items
    expect(text).toContain('\uD83D\uDC68') // creator avatar
    // Assignee avatar for the first item
    expect(text).toContain('\uD83D\uDC69') // assignee avatar
    // Assignee name
    expect(text).toContain('成员')
  })

  it('renders unassigned todo', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // The second todo has no assignee, should show "未指派"
    expect(text).toContain('未指派')
  })

  it('toggles todo and refreshes list', async () => {
    mockToggleTodo.mockResolvedValue({ code: 0, message: 'ok' })

    const updatedData = {
      ...mockTodoData,
      data: {
        ...mockTodoData.data,
        list: mockTodoData.data.list.map(t =>
          t.id === 1 ? { ...t, status: 'completed', completed_at: '2026-05-23' } : t
        ),
      },
    }
    mockGetTodoList.mockResolvedValueOnce(mockTodoData)
    mockGetTodoList.mockResolvedValueOnce(updatedData)

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetTodoList).toHaveBeenCalledTimes(1)

    // Call handleToggle on the first todo
    const vm = wrapper.vm as any
    await vm.handleToggle(mockTodoData.data.list[0])
    await flushPromises()
    await nextTick()

    expect(mockToggleTodo).toHaveBeenCalledWith(1)
    expect(mockGetTodoList).toHaveBeenCalledTimes(2)
  })

  it('creates todo and refreshes list', async () => {
    mockCreateTodo.mockResolvedValue({ code: 0, message: 'ok', data: { id: 99 } })

    const updatedData = {
      ...mockTodoData,
      data: { ...mockTodoData.data, total: 3 },
    }
    mockGetTodoList.mockResolvedValueOnce(mockTodoData)
    mockGetTodoList.mockResolvedValueOnce(updatedData)

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetTodoList).toHaveBeenCalledTimes(1)

    // Fill form and submit
    const vm = wrapper.vm as any
    vm.dialogOpen = true
    vm.form.title = '新待办'
    vm.form.description = '描述'
    vm.form.priority = 'normal'
    await nextTick()

    await vm.handleSubmit()
    await flushPromises()
    await nextTick()

    expect(mockCreateTodo).toHaveBeenCalled()
    expect(mockGetTodoList).toHaveBeenCalledTimes(2)
  })
})
