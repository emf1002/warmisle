import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => ({
  mockRequest: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import { login, checkInit, setupInit } from '../auth'

describe('Auth API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('login', () => {
    it('sends username and password via POST', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { token: 'abc123', member: { id: 1, username: 'admin', name: 'Admin', avatar: 'cat', role: 'admin' } },
      })

      await login('admin', 'pass123')
      expect(mockRequest.post).toHaveBeenCalledWith('/auth/login', { username: 'admin', password: 'pass123' })
    })

    it('returns snake_case keys in data.token and data.member', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          token: 'abc123',
          member: { id: 1, username: 'admin', name: 'Admin', avatar: 'cat', role: 'admin' },
        },
      })

      const res = await login('admin', 'pass123')
      expect(res.data.token).toBe('abc123')
      expect(res.data.member.id).toBe(1)
      expect(res.data.member.username).toBe('admin')
      expect(res.data.member.name).toBe('Admin')
      expect(res.data.member.avatar).toBe('cat')
      expect(res.data.member.role).toBe('admin')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          Token: 'abc123',
          Member: { Id: 1, Username: 'admin', Name: 'Admin', Avatar: 'cat', Role: 'admin' },
        },
      })

      const res = await login('admin', 'pass123')
      expect(res.data.token).toBeUndefined()
      expect(res.data.member).toBeUndefined()
    })
  })

  describe('checkInit', () => {
    it('sends GET to /init/check', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: { need_init: true } })

      await checkInit()
      expect(mockRequest.get).toHaveBeenCalledWith('/init/check')
    })

    it('returns snake_case keys with need_init boolean', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: { need_init: true } })

      const res = await checkInit()
      expect(typeof res.data.need_init).toBe('boolean')
      expect(res.data.need_init).toBe(true)
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: { NeedInit: true } })

      const res = await checkInit()
      expect(res.data.need_init).toBeUndefined()
    })
  })

  describe('setupInit', () => {
    it('sends admin_name, username, and password via POST', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { token: 'xyz789', member: { id: 1, username: 'admin', name: 'Boss', avatar: 'dog', role: 'admin' } },
      })

      await setupInit('Boss', 'admin', 'secret')
      expect(mockRequest.post).toHaveBeenCalledWith('/init/setup', { admin_name: 'Boss', username: 'admin', password: 'secret' })
    })

    it('returns snake_case keys in data.token and data.member', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          token: 'xyz789',
          member: { id: 1, username: 'admin', name: 'Boss', avatar: 'dog', role: 'admin' },
        },
      })

      const res = await setupInit('Boss', 'admin', 'secret')
      expect(res.data.token).toBe('xyz789')
      expect(res.data.member.id).toBe(1)
      expect(res.data.member.username).toBe('admin')
      expect(res.data.member.name).toBe('Boss')
      expect(res.data.member.avatar).toBe('dog')
      expect(res.data.member.role).toBe('admin')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          Token: 'xyz789',
          Member: { Id: 1, Username: 'admin', Name: 'Boss', Avatar: 'dog', Role: 'admin' },
        },
      })

      const res = await setupInit('Boss', 'admin', 'secret')
      expect(res.data.token).toBeUndefined()
      expect(res.data.member).toBeUndefined()
    })
  })
})
