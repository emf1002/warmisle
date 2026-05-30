/**
 * Shared vi.mock factory for @/stores/auth.
 *
 * Provides the same JWT-parsing mock that the real auth store uses, so tests
 * can control the current user via localStorage.getItem('token').
 *
 * Usage in a test file (must appear BEFORE the import of the component under test):
 *
 *   import '@/test-utils/auth-mock'
 *
 * Then set the token in beforeEach:
 *
 *   const payload = btoa(JSON.stringify({ member_id: 1, role: 'admin', username: 'admin' }))
 *   localStorage.setItem('token', `header.${payload}.sig`)
 */
import { vi } from 'vitest'

vi.mock('@/stores/auth', () => {
  function parseJwt(raw: string) {
    if (!raw) return null
    try {
      return JSON.parse(atob(raw.split('.')[1]))
    } catch {
      return null
    }
  }
  return {
    useAuthStore: () => {
      const token = localStorage.getItem('token') || ''
      const payload = parseJwt(token)
      return {
        currentUserId: (payload?.member_id as number) || 0,
        currentUserRole: payload?.role || '',
        isAdmin: payload?.role === 'admin',
      }
    },
  }
})
