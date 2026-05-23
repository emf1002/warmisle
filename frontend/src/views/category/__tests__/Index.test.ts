import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

/** Flush all pending microtasks (Promise resolutions + Vue scheduler). */
const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

// ── Mocks ────────────────────────────────────────────────────────────────────

const {
  mockGetCategories,
  mockCreateCategory,
  mockUpdateCategory,
  mockDeleteCategory,
  mockModalConfirm,
  mockMessageSuccess,
  mockMessageError,
} = vi.hoisted(() => ({
  mockGetCategories: vi.fn(),
  mockCreateCategory: vi.fn(),
  mockUpdateCategory: vi.fn(),
  mockDeleteCategory: vi.fn(),
  mockModalConfirm: vi.fn(),
  mockMessageSuccess: vi.fn(),
  mockMessageError: vi.fn(),
}))

vi.mock('@/api/category', () => ({
  getCategories: mockGetCategories,
  createCategory: mockCreateCategory,
  updateCategory: mockUpdateCategory,
  deleteCategory: mockDeleteCategory,
}))

// Mock ant-design-vue JS utilities (Modal.confirm, message)
vi.mock('ant-design-vue', () => ({
  message: {
    success: mockMessageSuccess,
    error: mockMessageError,
  },
  Modal: {
    confirm: mockModalConfirm,
  },
}))

import Category from '../Index.vue'

// ── Stubs ────────────────────────────────────────────────────────────────────

const stubs: Record<string, any> = {
  'a-button': {
    props: ['type', 'size', 'danger', 'htmlType', 'loading', 'block'],
    template:
      '<button v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>',
    emits: ['click'],
  },
  'a-modal': {
    props: ['open', 'title', 'confirmLoading'],
    template:
      '<div v-if="open" v-bind="$attrs"><div class="modal-title">{{ title }}</div><slot /><button class="modal-ok-btn" @click="$emit(\'ok\')">OK</button></div>',
    emits: ['ok', 'update:open'],
  },
  'a-form': {
    template: '<form v-bind="$attrs"><slot /></form>',
  },
  'a-form-item': {
    props: ['label', 'required'],
    template: '<div v-bind="$attrs"><slot /></div>',
  },
  'a-input': {
    props: ['value', 'placeholder', 'maxlength'],
    template:
      '<input v-bind="$attrs" :value="value" :placeholder="placeholder" @input="$emit(\'update:value\', $event.target.value)" />',
    emits: ['update:value'],
  },
  'a-select': {
    props: ['value'],
    template:
      '<select v-bind="$attrs" :value="value" @change="$emit(\'update:value\', $event.target.value)"><slot /></select>',
    emits: ['update:value'],
  },
  'a-select-option': {
    props: ['value'],
    template: '<option :value="value"><slot /></option>',
  },
  'a-input-number': {
    props: ['value', 'min', 'max', 'placeholder'],
    template:
      '<input type="number" v-bind="$attrs" :value="value" :min="min" :max="max" :placeholder="placeholder" @input="$emit(\'update:value\', Number($event.target.value))" />',
    emits: ['update:value'],
  },
  'router-link': {
    template: '<a><slot /></a>',
  },
  'router-view': {
    template: '<div><slot /></div>',
  },
}

const mockCategories = [
  { id: 1, type: 'expense', name: '餐饮', icon: '🍱', sort_order: 1, preset: true },
  { id: 2, type: 'expense', name: '交通', icon: '🚗', sort_order: 2, preset: false },
  { id: 3, type: 'income', name: '工资', icon: '💰', sort_order: 1, preset: true },
]

function createWrapper() {
  return mount(Category, {
    global: { stubs },
  })
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe('Category view', () => {
  beforeEach(() => {
    vi.clearAllMocks()

    // Set admin JWT token so isAdmin is true
    const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
    localStorage.setItem('token', `header.${payload}.sig`)

    // Default API response
    mockGetCategories.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: mockCategories,
    })
  })

  it('renders categories from API', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('餐饮')
    expect(text).toContain('🍱')
    expect(text).toContain('交通')
    expect(text).toContain('🚗')
    expect(text).toContain('工资')
    expect(text).toContain('💰')
  })

  it('renders expense and income sections separately', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const expenseSection = wrapper.find('[data-testid="expense-categories"]')
    const incomeSection = wrapper.find('[data-testid="income-categories"]')

    expect(expenseSection.exists()).toBe(true)
    expect(incomeSection.exists()).toBe(true)

    const expenseText = expenseSection.text()
    const incomeText = incomeSection.text()

    // Expense section should contain expense categories
    expect(expenseText).toContain('餐饮')
    expect(expenseText).toContain('交通')
    expect(expenseText).not.toContain('工资')

    // Income section should contain income categories
    expect(incomeText).toContain('工资')
    expect(incomeText).not.toContain('餐饮')
    expect(incomeText).not.toContain('交通')
  })

  it('creates category and refreshes list', async () => {
    mockCreateCategory.mockResolvedValue({ code: 0, message: 'ok', data: { id: 4 } })

    // After create, return updated list with new item
    const updatedCategories = [
      ...mockCategories,
      { id: 4, type: 'expense', name: '娱乐', icon: '🎮', sort_order: 3, preset: false },
    ]
    mockGetCategories
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: mockCategories })
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: updatedCategories })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    // Verify initial load
    expect(mockGetCategories).toHaveBeenCalledTimes(1)

    // Click add button to open the create dialog
    await wrapper.find('[data-testid="add-btn"]').trigger('click')
    await nextTick()

    // Fill in the form: type is already 'expense' by default, set the name
    const nameInput = wrapper.find('[data-testid="name-input"]')
    await nameInput.setValue('娱乐')
    await nextTick()

    // Click the modal OK button to submit
    await wrapper.find('.modal-ok-btn').trigger('click')
    await flushPromises()
    await nextTick()

    // Verify createCategory was called with correct data
    expect(mockCreateCategory).toHaveBeenCalledWith(
      expect.objectContaining({
        type: 'expense',
        name: '娱乐',
      }),
    )

    // Verify getCategories was called again (refresh)
    expect(mockGetCategories).toHaveBeenCalledTimes(2)
  })

  it('deletes category and refreshes list', async () => {
    mockDeleteCategory.mockResolvedValue({ code: 0, message: 'ok' })

    // After delete, return list without the deleted item
    const reducedCategories = mockCategories.filter(c => c.id !== 2)
    mockGetCategories
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: mockCategories })
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: reducedCategories })

    // Make Modal.confirm immediately invoke onOk
    mockModalConfirm.mockImplementation(({ onOk }) => {
      onOk()
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetCategories).toHaveBeenCalledTimes(1)

    // Find all delete buttons and click the first one (for category id=1 "餐饮")
    const deleteButtons = wrapper.findAll('[data-testid="delete-btn"]')
    expect(deleteButtons.length).toBeGreaterThan(0)
    await deleteButtons[0].trigger('click')
    await flushPromises()
    await nextTick()

    // Verify Modal.confirm was called
    expect(mockModalConfirm).toHaveBeenCalled()

    // Verify deleteCategory was called
    expect(mockDeleteCategory).toHaveBeenCalled()

    // Verify getCategories was called again (refresh)
    expect(mockGetCategories).toHaveBeenCalledTimes(2)
  })
})
