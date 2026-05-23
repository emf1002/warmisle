import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

/** Flush all pending microtasks (Promise resolutions + Vue scheduler). */
const flushPromises = () => new Promise<void>(resolve => setTimeout(resolve, 0))

// ── Mocks ────────────────────────────────────────────────────────────────────

const { mockLogin } = vi.hoisted(() => ({
  mockLogin: vi.fn(),
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({ login: mockLogin }),
}))

const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush }),
  useRoute: () => ({}),
}))

import Login from '../Login.vue'

// ── Stubs ────────────────────────────────────────────────────────────────────

// Override global `true` stubs with slot-forwarding versions so rendered text
// is visible to assertions.
const stubs: Record<string, any> = {
  'a-form': {
    template: '<form @submit.prevent="$emit(\'finish\')"><slot /></form>',
    emits: ['finish'],
  },
  'a-form-item': {
    props: ['name'],
    template: '<div><slot /></div>',
  },
  'a-input': {
    props: ['value', 'placeholder', 'size', 'data-testid'],
    template:
      '<input :value="value" :placeholder="placeholder" :data-testid="data-testid" @input="$emit(\'update:value\', $event.target.value)" />',
    emits: ['update:value'],
  },
  'a-input-password': {
    props: ['value', 'placeholder', 'size', 'data-testid'],
    template:
      '<input type="password" :value="value" :placeholder="placeholder" :data-testid="data-testid" @input="$emit(\'update:value\', $event.target.value)" />',
    emits: ['update:value'],
  },
  'a-button': {
    props: ['type', 'htmlType', 'size', 'loading', 'block', 'data-testid'],
    template:
      '<button :type="htmlType || \'button\'" :data-testid="data-testid"><slot /></button>',
  },
}

function createWrapper() {
  return mount(Login, {
    global: { stubs },
  })
}

// ── Tests ────────────────────────────────────────────────────────────────────

describe('Login view', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renders login form', async () => {
    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('暖屿')
    expect(text).toContain('登 录')
  })

  it('calls auth store login with form values', async () => {
    mockLogin.mockResolvedValue({ code: 0, data: { token: 'tok' } })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    // Fill in username and password via the stubbed inputs
    const inputs = wrapper.findAll('input')
    // First input is username, second is password
    await inputs[0].setValue('testuser')
    await inputs[1].setValue('testpass')
    await nextTick()

    // Submit the form (triggers the 'finish' event on the a-form stub)
    await wrapper.find('form').trigger('submit')
    await flushPromises()
    await nextTick()

    expect(mockLogin).toHaveBeenCalledWith('testuser', 'testpass')
  })

  it('redirects to / on success', async () => {
    mockLogin.mockResolvedValue({ code: 0, data: { token: 'tok' } })

    const wrapper = createWrapper()
    await flushPromises()
    await nextTick()

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('admin')
    await inputs[1].setValue('password123')
    await nextTick()

    await wrapper.find('form').trigger('submit')
    await flushPromises()
    await nextTick()

    expect(mockPush).toHaveBeenCalledWith('/')
  })
})
