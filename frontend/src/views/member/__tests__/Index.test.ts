import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

/** Flush all pending microtasks (Promise resolutions + Vue scheduler). */
const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

// ── Mocks ────────────────────────────────────────────────────────────────────

const {
  mockGetMembers,
  mockCreateMember,
  mockUpdateMember,
  mockDeleteMember,
  mockDisableMember,
  mockEnableMember,
  mockResetPassword,
  mockModalConfirm,
  mockMessageSuccess,
  mockMessageError,
} = vi.hoisted(() => ({
  mockGetMembers: vi.fn(),
  mockCreateMember: vi.fn(),
  mockUpdateMember: vi.fn(),
  mockDeleteMember: vi.fn(),
  mockDisableMember: vi.fn(),
  mockEnableMember: vi.fn(),
  mockResetPassword: vi.fn(),
  mockModalConfirm: vi.fn(),
  mockMessageSuccess: vi.fn(),
  mockMessageError: vi.fn(),
}))

vi.mock('@/api/member', () => ({
  getMembers: mockGetMembers,
  createMember: mockCreateMember,
  updateMember: mockUpdateMember,
  deleteMember: mockDeleteMember,
  disableMember: mockDisableMember,
  enableMember: mockEnableMember,
  resetPassword: mockResetPassword,
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

// Shared auth store mock (parses JWT from localStorage)
import '@/test/auth-mock'

import Member from '../Index.vue'

// ── Stubs ────────────────────────────────────────────────────────────────────

const stubs: Record<string, any> = {
  'a-button': {
    props: ['type', 'size', 'danger', 'htmlType', 'loading', 'block'],
    template:
      '<button v-bind="$attrs" @click="$emit(\'click\', $event)"><slot /></button>',
    emits: ['click'],
  },
  'a-table': {
    props: ['columns', 'dataSource', 'rowKey', 'pagination'],
    template: `<table v-bind="$attrs">
      <tbody>
        <tr v-for="(record, i) in dataSource || []" :key="record[rowKey] || i">
          <td v-for="col in columns || []" :key="col.key || col.dataIndex">
            <slot name="bodyCell" :column="col" :record="record" :text="record[col.dataIndex]" :index="i" />
          </td>
        </tr>
      </tbody>
    </table>`,
  },
  'a-tag': {
    props: ['color', 'size'],
    template: '<span v-bind="$attrs"><slot /></span>',
  },
  'a-space': {
    template: '<span v-bind="$attrs"><slot /></span>',
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
  'a-input-password': {
    props: ['value', 'placeholder', 'maxlength'],
    template:
      '<input type="password" v-bind="$attrs" :value="value" :placeholder="placeholder" @input="$emit(\'update:value\', $event.target.value)" />',
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
  'router-link': {
    template: '<a><slot /></a>',
  },
  'router-view': {
    template: '<div><slot /></div>',
  },
}

const mockMembers = [
  { id: 1, username: 'admin', name: '管理员', avatar: '👨', role: 'admin', status: 'active' },
  { id: 2, username: 'member', name: '普通成员', avatar: '👩', role: 'member', status: 'active' },
  { id: 3, username: 'disabled', name: '禁用成员', avatar: '👶', role: 'member', status: 'disabled' },
]

function createWrapper() {
  return mount(Member, {
    global: { stubs },
  })
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe('Member view', () => {
  beforeEach(() => {
    vi.clearAllMocks()

    // Set admin JWT token so isAdmin is true
    const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
    localStorage.setItem('token', `header.${payload}.sig`)

    // Default API response
    mockGetMembers.mockResolvedValue({
      code: 0,
      message: 'ok',
      data: mockMembers,
    })
  })

  it('renders member list from API', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('管理员')
    expect(text).toContain('👨')
    expect(text).toContain('admin')
    expect(text).toContain('普通成员')
    expect(text).toContain('👩')
    expect(text).toContain('member')
  })

  it('renders role tags', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // admin role renders as "管理员"
    expect(text).toContain('管理员')
    // member role renders as "成员"
    expect(text).toContain('成员')
  })

  it('renders status tags', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    // active status renders as "正常"
    expect(text).toContain('正常')
    // disabled status renders as "已禁用"
    expect(text).toContain('已禁用')
  })

  it('creates member and refreshes list', async () => {
    mockCreateMember.mockResolvedValue({ code: 0, message: 'ok', data: { id: 4 } })

    // After create, return updated list with new member
    const updatedMembers = [
      ...mockMembers,
      { id: 4, username: 'newuser', name: '新成员', avatar: '🌟', role: 'member', status: 'active' },
    ]
    mockGetMembers
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: mockMembers })
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: updatedMembers })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    // Verify initial load
    expect(mockGetMembers).toHaveBeenCalledTimes(1)

    // Click add button to open the create dialog
    await wrapper.find('[data-testid="add-btn"]').trigger('click')
    await nextTick()

    // Fill in the form
    const usernameInput = wrapper.find('[data-testid="username-input"]')
    await usernameInput.setValue('newuser')
    const passwordInput = wrapper.find('[data-testid="password-input"]')
    await passwordInput.setValue('pass123')
    const nameInput = wrapper.find('[data-testid="name-input"]')
    await nameInput.setValue('新成员')
    await nextTick()

    // Click the modal OK button to submit
    await wrapper.find('.modal-ok-btn').trigger('click')
    await flushPromises()
    await nextTick()

    // Verify createMember was called with correct data
    expect(mockCreateMember).toHaveBeenCalledWith(
      expect.objectContaining({
        username: 'newuser',
        password: 'pass123',
        name: '新成员',
      }),
    )

    // Verify getMembers was called again (refresh)
    expect(mockGetMembers).toHaveBeenCalledTimes(2)
  })

  it('disables member and refreshes list', async () => {
    mockDisableMember.mockResolvedValue({ code: 0, message: 'ok' })

    // After disable, return updated list with disabled member
    const updatedMembers = mockMembers.map(m =>
      m.id === 2 ? { ...m, status: 'disabled' } : m,
    )
    mockGetMembers
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: mockMembers })
      .mockResolvedValueOnce({ code: 0, message: 'ok', data: updatedMembers })

    // Make Modal.confirm immediately invoke onOk
    mockModalConfirm.mockImplementation(({ onOk }) => {
      onOk()
    })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    expect(mockGetMembers).toHaveBeenCalledTimes(1)

    // Find disable buttons (for members with status === 'active')
    const disableButtons = wrapper.findAll('[data-testid="disable-btn"]')
    expect(disableButtons.length).toBeGreaterThan(0)

    // Click the first disable button
    await disableButtons[0].trigger('click')
    await flushPromises()
    await nextTick()

    // Verify Modal.confirm was called
    expect(mockModalConfirm).toHaveBeenCalled()

    // Verify disableMember was called
    expect(mockDisableMember).toHaveBeenCalled()

    // Verify getMembers was called again (refresh)
    expect(mockGetMembers).toHaveBeenCalledTimes(2)
  })
})
