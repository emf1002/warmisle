import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => ({
  mockRequest: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import { getLedgers, getLedgerById, createLedger, updateLedger, deleteLedger } from '../ledger'

describe('Ledger API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getLedgers', () => {
    it('should return snake_case keys: summary, groups, total, page, page_size', async () => {
      mockRequest.get.mockResolvedValue({
        data: {
          summary: { income: 10000, expense: 5000, balance: 5000 },
          groups: [
            {
              date: '2026-05-23',
              daily_total: 3000,
              items: [
                {
                  id: 1,
                  amount: 3000,
                  note: '午餐',
                  category_id: 1,
                  creator_id: 1,
                  occurred_at: '2026-05-23',
                  category: { id: 1, name: '餐饮', icon: '🍱', type: 'expense' },
                  creator: { id: 1, name: '管理员', avatar: '👨' },
                  members: [{ id: 1, name: '管理员' }],
                },
              ],
            },
          ],
          total: 1,
          page: 1,
          page_size: 20,
        },
      })

      const res = await getLedgers({ month: '2026-05', page: 1, page_size: 20 })

      expect(mockRequest.get).toHaveBeenCalledWith('/ledgers', {
        params: { month: '2026-05', page: 1, page_size: 20 },
      })

      // Top-level response keys
      expect(res.data.total).toBe(1)
      expect(res.data.page).toBe(1)
      expect(res.data.page_size).toBe(20)

      // Summary keys
      expect(res.data.summary.income).toBe(10000)
      expect(res.data.summary.expense).toBe(5000)
      expect(res.data.summary.balance).toBe(5000)

      // Group-level keys
      expect(res.data.groups[0].daily_total).toBe(3000)

      // Item-level snake_case keys
      const item = res.data.groups[0].items[0]
      expect(item.category_id).toBe(1)
      expect(item.creator_id).toBe(1)
      expect(item.occurred_at).toBe('2026-05-23')

      // Nested category
      expect(item.category.type).toBe('expense')
      expect(item.category.icon).toBe('🍱')

      // Nested creator
      expect(item.creator.avatar).toBe('👨')

      // Nested members
      expect(item.members[0].name).toBe('管理员')
    })
  })

  describe('getLedgerById', () => {
    it('should fetch a single ledger by id', async () => {
      mockRequest.get.mockResolvedValue({
        data: {
          id: 1,
          amount: 3000,
          note: '午餐',
          category_id: 1,
          creator_id: 1,
          occurred_at: '2026-05-23',
        },
      })

      const res = await getLedgerById(1)

      expect(mockRequest.get).toHaveBeenCalledWith('/ledgers/1')
      expect(res.data.id).toBe(1)
      expect(res.data.amount).toBe(3000)
      expect(res.data.category_id).toBe(1)
      expect(res.data.creator_id).toBe(1)
      expect(res.data.occurred_at).toBe('2026-05-23')
    })
  })

  describe('createLedger', () => {
    it('should return created item with snake_case keys: id, amount, note, category_id, creator_id, occurred_at', async () => {
      const newItem = {
        amount: 5000,
        note: '晚餐',
        category_id: 2,
        member_ids: [1],
        occurred_at: '2026-05-23',
      }

      mockRequest.post.mockResolvedValue({
        data: {
          id: 2,
          amount: 5000,
          note: '晚餐',
          category_id: 2,
          creator_id: 1,
          occurred_at: '2026-05-23',
        },
      })

      const res = await createLedger(newItem)

      expect(mockRequest.post).toHaveBeenCalledWith('/ledgers', newItem)
      expect(res.data.id).toBe(2)
      expect(res.data.amount).toBe(5000)
      expect(res.data.note).toBe('晚餐')
      expect(res.data.category_id).toBe(2)
      expect(res.data.creator_id).toBe(1)
      expect(res.data.occurred_at).toBe('2026-05-23')
    })
  })

  describe('updateLedger', () => {
    it('should send PUT request with correct id and data', async () => {
      const updateData = { amount: 6000, note: '更新晚餐' }

      mockRequest.put.mockResolvedValue({
        data: {
          id: 2,
          amount: 6000,
          note: '更新晚餐',
          category_id: 2,
          creator_id: 1,
          occurred_at: '2026-05-23',
        },
      })

      const res = await updateLedger(2, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/ledgers/2', updateData)
      expect(res.data.id).toBe(2)
      expect(res.data.amount).toBe(6000)
      expect(res.data.note).toBe('更新晚餐')
      expect(res.data.category_id).toBe(2)
      expect(res.data.creator_id).toBe(1)
      expect(res.data.occurred_at).toBe('2026-05-23')
    })
  })

  describe('deleteLedger', () => {
    it('should send DELETE request with correct id', async () => {
      mockRequest.delete.mockResolvedValue({ data: null })

      const res = await deleteLedger(2)

      expect(mockRequest.delete).toHaveBeenCalledWith('/ledgers/2')
      expect(res.data).toBeNull()
    })
  })
})
