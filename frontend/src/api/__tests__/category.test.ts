import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => ({
  mockRequest: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import { getCategories, createCategory, updateCategory, deleteCategory } from '../category'

describe('Category API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getCategories', () => {
    it('sends GET to /categories', async () => {
      mockRequest.get.mockResolvedValue({ code: 0, message: 'ok', data: [] })

      await getCategories()
      expect(mockRequest.get).toHaveBeenCalledWith('/categories')
    })

    it('returns snake_case keys on category items', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: [
          { id: 1, type: 'expense', name: 'Food', icon: 'fork', sort_order: 1, preset: true },
          { id: 2, type: 'income', name: 'Salary', icon: 'money', sort_order: 1, preset: true },
        ],
      })

      const res = await getCategories()
      const item = res.data[0]
      expect(item.id).toBe(1)
      expect(item.type).toBe('expense')
      expect(item.name).toBe('Food')
      expect(item.icon).toBe('fork')
      expect(item.sort_order).toBe(1)
      expect(item.preset).toBe(true)
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.get.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: [
          { Id: 1, Type: 'expense', Name: 'Food', Icon: 'fork', SortOrder: 1, Preset: true },
        ],
      })

      const res = await getCategories()
      const item = res.data[0]
      expect(item.id).toBeUndefined()
      expect(item.type).toBeUndefined()
      expect(item.name).toBeUndefined()
      expect(item.sort_order).toBeUndefined()
      expect(item.preset).toBeUndefined()
    })
  })

  describe('createCategory', () => {
    it('sends category data via POST', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 3, type: 'expense', name: 'Transport', icon: 'car', sort_order: 2 },
      })

      await createCategory({ type: 'expense', name: 'Transport', icon: 'car', sort_order: 2 })
      expect(mockRequest.post).toHaveBeenCalledWith('/categories', {
        type: 'expense',
        name: 'Transport',
        icon: 'car',
        sort_order: 2,
      })
    })

    it('returns snake_case keys on created category', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 3, type: 'expense', name: 'Transport', icon: 'car', sort_order: 2 },
      })

      const res = await createCategory({ type: 'expense', name: 'Transport', icon: 'car', sort_order: 2 })
      expect(res.data.id).toBe(3)
      expect(res.data.type).toBe('expense')
      expect(res.data.name).toBe('Transport')
      expect(res.data.icon).toBe('car')
      expect(res.data.sort_order).toBe(2)
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.post.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { Id: 3, Type: 'expense', Name: 'Transport', Icon: 'car', SortOrder: 2 },
      })

      const res = await createCategory({ type: 'expense', name: 'Transport', icon: 'car', sort_order: 2 })
      expect(res.data.id).toBeUndefined()
      expect(res.data.sort_order).toBeUndefined()
    })
  })

  describe('updateCategory', () => {
    it('sends updated data via PUT', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 1, type: 'expense', name: 'Dining', icon: 'fork', sort_order: 1 },
      })

      await updateCategory(1, { name: 'Dining' })
      expect(mockRequest.put).toHaveBeenCalledWith('/categories/1', { name: 'Dining' })
    })

    it('returns snake_case keys on updated category', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { id: 1, type: 'expense', name: 'Dining', icon: 'fork', sort_order: 1 },
      })

      const res = await updateCategory(1, { name: 'Dining' })
      expect(res.data.id).toBe(1)
      expect(res.data.name).toBe('Dining')
      expect(res.data.sort_order).toBe(1)
    })

    it('would fail with PascalCase (missing json tags)', async () => {
      mockRequest.put.mockResolvedValue({
        code: 0,
        message: 'ok',
        data: { Id: 1, Type: 'expense', Name: 'Dining', Icon: 'fork', SortOrder: 1 },
      })

      const res = await updateCategory(1, { name: 'Dining' })
      expect(res.data.id).toBeUndefined()
      expect(res.data.sort_order).toBeUndefined()
    })
  })

  describe('deleteCategory', () => {
    it('sends DELETE to /categories/:id', async () => {
      mockRequest.delete.mockResolvedValue({ code: 0, message: 'ok', data: null })

      await deleteCategory(5)
      expect(mockRequest.delete).toHaveBeenCalledWith('/categories/5')
    })

    it('returns success response', async () => {
      mockRequest.delete.mockResolvedValue({ code: 0, message: 'ok', data: null })

      const res = await deleteCategory(5)
      expect(res.data.code).toBe(0)
    })
  })
})
