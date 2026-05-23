import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => ({
  mockRequest: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import { getSummary, getExpenseChart, getUpcomingTodos, getWishTrends, getForumHot } from '../dashboard'

describe('Dashboard API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getSummary', () => {
    it('should return snake_case keys for income, expense, balance', async () => {
      mockRequest.get.mockResolvedValue({
        data: { income: 10000, expense: 5000, balance: 5000 },
      })

      const res = await getSummary('2026-05')

      expect(mockRequest.get).toHaveBeenCalledWith('/dashboard/summary', { params: { month: '2026-05' } })
      expect(res.data.income).toBe(10000)
      expect(res.data.expense).toBe(5000)
      expect(res.data.balance).toBe(5000)
    })

    it('should fail if PascalCase keys are used instead of snake_case', async () => {
      mockRequest.get.mockResolvedValue({
        data: { Income: 10000, Expense: 5000, Balance: 5000 },
      })

      const res = await getSummary()

      // PascalCase keys should NOT exist on the response
      expect(res.data.income).toBeUndefined()
      expect(res.data.expense).toBeUndefined()
      expect(res.data.balance).toBeUndefined()
    })
  })

  describe('getExpenseChart', () => {
    it('should return snake_case keys: category_id, category_name, icon, amount', async () => {
      mockRequest.get.mockResolvedValue({
        data: [{ category_id: 1, category_name: '餐饮', icon: '🍱', amount: 3000 }],
      })

      const res = await getExpenseChart('2026-05')

      expect(mockRequest.get).toHaveBeenCalledWith('/dashboard/expense-chart', { params: { month: '2026-05' } })
      expect(res.data[0].category_id).toBe(1)
      expect(res.data[0].category_name).toBe('餐饮')
      expect(res.data[0].icon).toBe('🍱')
      expect(res.data[0].amount).toBe(3000)
    })

    it('should fail if PascalCase keys are used instead of snake_case', async () => {
      mockRequest.get.mockResolvedValue({
        data: [{ CategoryId: 1, CategoryName: '餐饮', Icon: '🍱', Amount: 3000 }],
      })

      const res = await getExpenseChart()

      // PascalCase keys should NOT exist on the response
      expect(res.data[0].category_id).toBeUndefined()
      expect(res.data[0].category_name).toBeUndefined()
    })
  })

  describe('getUpcomingTodos', () => {
    it('should return snake_case keys: title, priority, due_date, assignee.name', async () => {
      mockRequest.get.mockResolvedValue({
        data: [
          {
            id: 1,
            title: '买菜',
            priority: 'important',
            due_date: '2026-05-25',
            assignee: { id: 1, name: '管理员', avatar: '👨' },
          },
        ],
      })

      const res = await getUpcomingTodos()

      expect(mockRequest.get).toHaveBeenCalledWith('/dashboard/upcoming-todos')
      expect(res.data[0].title).toBe('买菜')
      expect(res.data[0].priority).toBe('important')
      expect(res.data[0].due_date).toBe('2026-05-25')
      expect(res.data[0].assignee.name).toBe('管理员')
    })
  })

  describe('getWishTrends', () => {
    it('should return snake_case keys: vote_count, creator.name, status', async () => {
      mockRequest.get.mockResolvedValue({
        data: [
          {
            id: 1,
            title: '买iPad',
            creator: { name: '管理员' },
            vote_count: 3,
            status: 'pending',
          },
        ],
      })

      const res = await getWishTrends()

      expect(mockRequest.get).toHaveBeenCalledWith('/dashboard/wish-trends')
      expect(res.data[0].vote_count).toBe(3)
      expect(res.data[0].creator.name).toBe('管理员')
      expect(res.data[0].status).toBe('pending')
    })
  })

  describe('getForumHot', () => {
    it('should return snake_case keys: type, content, creator.name, created_at', async () => {
      mockRequest.get.mockResolvedValue({
        data: [
          {
            type: 'post',
            title: '',
            content: '今天天气好',
            creator: { name: '管理员' },
            created_at: '2026-05-23T10:00:00Z',
          },
        ],
      })

      const res = await getForumHot()

      expect(mockRequest.get).toHaveBeenCalledWith('/dashboard/forum-hot')
      expect(res.data[0].type).toBe('post')
      expect(res.data[0].content).toBe('今天天气好')
      expect(res.data[0].creator.name).toBe('管理员')
      expect(res.data[0].created_at).toBe('2026-05-23T10:00:00Z')
    })
  })
})
