import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => ({
  mockRequest: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import { getMembers, createMember, updateMember, getProfile, updateProfile } from '../member'

describe('Member API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMembers', () => {
    it('sends GET to /members', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: [] })

      await getMembers()
      expect(mockRequest.get).toHaveBeenCalledWith('/members')
    })

    it('returns snake_case keys on member items', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: [
          { id: 1, username: 'admin', name: 'Admin', avatar: 'cat', role: 'admin', status: 'active' },
          { id: 2, username: 'user1', name: 'User One', avatar: 'dog', role: 'member', status: 'active' },
        ],
      })

      const res = await getMembers()
      const item = res.data[0]
      expect(item.id).toBe(1)
      expect(item.username).toBe('admin')
      expect(item.name).toBe('Admin')
      expect(item.avatar).toBe('cat')
      expect(item.role).toBe('admin')
      expect(item.status).toBe('active')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: [
          { Id: 1, Username: 'admin', Name: 'Admin', Avatar: 'cat', Role: 'admin', Status: 'active' },
        ],
      })

      const res = await getMembers()
      const item = res.data[0]
      expect(item.id).toBeUndefined()
      expect(item.username).toBeUndefined()
      expect(item.name).toBeUndefined()
      expect(item.avatar).toBeUndefined()
      expect(item.role).toBeUndefined()
      expect(item.status).toBeUndefined()
    })
  })

  describe('createMember', () => {
    it('sends member data via POST', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 3, username: 'newuser', name: 'New User', avatar: 'rabbit', role: 'member', status: 'active' },
      })

      await createMember({ username: 'newuser', password: 'pass', name: 'New User', avatar: 'rabbit' })
      expect(mockRequest.post).toHaveBeenCalledWith('/members', {
        username: 'newuser',
        password: 'pass',
        name: 'New User',
        avatar: 'rabbit',
      })
    })

    it('returns snake_case keys on created member', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 3, username: 'newuser', name: 'New User', avatar: 'rabbit', role: 'member', status: 'active' },
      })

      const res = await createMember({ username: 'newuser', password: 'pass', name: 'New User', avatar: 'rabbit' })
      expect(res.data.id).toBe(3)
      expect(res.data.username).toBe('newuser')
      expect(res.data.name).toBe('New User')
      expect(res.data.avatar).toBe('rabbit')
      expect(res.data.role).toBe('member')
      expect(res.data.status).toBe('active')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { Id: 3, Username: 'newuser', Name: 'New User', Avatar: 'rabbit', Role: 'member', Status: 'active' },
      })

      const res = await createMember({ username: 'newuser', password: 'pass', name: 'New User', avatar: 'rabbit' })
      expect(res.data.id).toBeUndefined()
      expect(res.data.username).toBeUndefined()
      expect(res.data.status).toBeUndefined()
    })
  })

  describe('updateMember', () => {
    it('sends updated data via PUT', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 2, username: 'user1', name: 'Updated Name', avatar: 'bear', role: 'member', status: 'active' },
      })

      await updateMember(2, { name: 'Updated Name', avatar: 'bear' })
      expect(mockRequest.put).toHaveBeenCalledWith('/members/2', { name: 'Updated Name', avatar: 'bear' })
    })

    it('returns snake_case keys on updated member', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 2, username: 'user1', name: 'Updated Name', avatar: 'bear', role: 'member', status: 'active' },
      })

      const res = await updateMember(2, { name: 'Updated Name', avatar: 'bear' })
      expect(res.data.id).toBe(2)
      expect(res.data.name).toBe('Updated Name')
      expect(res.data.avatar).toBe('bear')
      expect(res.data.role).toBe('member')
      expect(res.data.status).toBe('active')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { Id: 2, Username: 'user1', Name: 'Updated Name', Avatar: 'bear', Role: 'member', Status: 'active' },
      })

      const res = await updateMember(2, { name: 'Updated Name', avatar: 'bear' })
      expect(res.data.id).toBeUndefined()
      expect(res.data.status).toBeUndefined()
    })
  })

  describe('getProfile', () => {
    it('sends GET to /profile', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: {} })

      await getProfile()
      expect(mockRequest.get).toHaveBeenCalledWith('/profile')
    })

    it('returns snake_case keys including last_login', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          id: 1,
          username: 'admin',
          name: 'Admin',
          avatar: 'cat',
          role: 'admin',
          status: 'active',
          last_login: '2026-05-23T10:00:00Z',
        },
      })

      const res = await getProfile()
      expect(res.data.id).toBe(1)
      expect(res.data.username).toBe('admin')
      expect(res.data.name).toBe('Admin')
      expect(res.data.avatar).toBe('cat')
      expect(res.data.role).toBe('admin')
      expect(res.data.status).toBe('active')
      expect(res.data.last_login).toBe('2026-05-23T10:00:00Z')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: {
          Id: 1,
          Username: 'admin',
          Name: 'Admin',
          Avatar: 'cat',
          Role: 'admin',
          Status: 'active',
          LastLogin: '2026-05-23T10:00:00Z',
        },
      })

      const res = await getProfile()
      expect(res.data.id).toBeUndefined()
      expect(res.data.last_login).toBeUndefined()
    })
  })

  describe('updateProfile', () => {
    it('sends profile data via PUT', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 1, username: 'admin', name: 'New Name', avatar: 'star', role: 'admin', status: 'active', last_login: '2026-05-23T10:00:00Z' },
      })

      await updateProfile({ name: 'New Name', avatar: 'star' })
      expect(mockRequest.put).toHaveBeenCalledWith('/profile', { name: 'New Name', avatar: 'star' })
    })

    it('returns snake_case keys on updated profile', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 1, username: 'admin', name: 'New Name', avatar: 'star', role: 'admin', status: 'active', last_login: '2026-05-23T10:00:00Z' },
      })

      const res = await updateProfile({ name: 'New Name', avatar: 'star' })
      expect(res.data.name).toBe('New Name')
      expect(res.data.avatar).toBe('star')
      expect(res.data.last_login).toBe('2026-05-23T10:00:00Z')
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { Id: 1, Username: 'admin', Name: 'New Name', Avatar: 'star', Role: 'admin', Status: 'active', LastLogin: '2026-05-23T10:00:00Z' },
      })

      const res = await updateProfile({ name: 'New Name', avatar: 'star' })
      expect(res.data.last_login).toBeUndefined()
    })
  })
})
