import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockRequest } = vi.hoisted(() => {
  const mockRequest = {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  }
  return { mockRequest }
})

vi.mock('@/utils/request', () => ({ default: mockRequest }))

import {
  getWishList,
  createWish,
  updateWish,
  deleteWish,
  promoteWish,
  updateWishStatus,
  voteWish,
  unvoteWish,
} from '../wish'

describe('Wish API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getWishList', () => {
    it('should call GET /wishes with params', async () => {
      const mockResponse = {
        data: {
          list: [
            {
              id: 1,
              title: '买iPad',
              description: '学习用',
              category: 'item',
              amount: 500000,
              priority: 'normal',
              type: 'personal',
              status: 'pending',
              creator_id: 1,
              vote_count: 3,
              creator: { id: 1, name: '管理员', avatar: '👨' },
            },
          ],
          total: 1,
          page: 1,
          page_size: 20,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getWishList({ page: 1, page_size: 20 })

      expect(mockRequest.get).toHaveBeenCalledWith('/wishes', {
        params: { page: 1, page_size: 20 },
      })
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in response data', async () => {
      const mockResponse = {
        data: {
          list: [
            {
              id: 1,
              title: '买iPad',
              description: '学习用',
              category: 'item',
              amount: 500000,
              priority: 'normal',
              type: 'personal',
              status: 'pending',
              creator_id: 1,
              vote_count: 3,
              creator: { id: 1, name: '管理员', avatar: '👨' },
            },
          ],
          total: 1,
          page: 1,
          page_size: 20,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getWishList({})
      const item = result.data.list[0]

      // Verify all snake_case keys exist
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('vote_count')

      // Verify nested creator object
      expect(item.creator).toHaveProperty('name')
      expect(item.creator).toHaveProperty('avatar')
    })

    it('should NOT have PascalCase keys (CreatorID would not match creator_id)', async () => {
      const pascalCaseResponse = {
        data: {
          list: [
            {
              ID: 1,
              Title: '买iPad',
              CreatorID: 1,
              VoteCount: 3,
              CreatedAt: '2026-05-23',
            },
          ],
          Total: 1,
          Page: 1,
          PageSize: 20,
        },
      }
      mockRequest.get.mockResolvedValue(pascalCaseResponse)

      const result = await getWishList({})
      const item = result.data.list[0]

      // PascalCase keys should NOT match the expected snake_case keys
      expect(item).not.toHaveProperty('creator_id')
      expect(item).not.toHaveProperty('vote_count')

      // These PascalCase keys would exist in a wrongly formatted response
      expect(item).toHaveProperty('CreatorID')
      expect(item).toHaveProperty('VoteCount')
    })
  })

  describe('createWish', () => {
    it('should call POST /wishes with data', async () => {
      const wishData = {
        title: '买iPad',
        description: '学习用',
        category: 'item',
        priority: 'normal',
        amount: 500000,
      }
      const mockResponse = {
        data: {
          id: 2,
          title: '买iPad',
          description: '学习用',
          category: 'item',
          amount: 500000,
          priority: 'normal',
          type: 'personal',
          status: 'pending',
          creator_id: 1,
          vote_count: 0,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createWish(wishData)

      expect(mockRequest.post).toHaveBeenCalledWith('/wishes', wishData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created wish', async () => {
      const mockResponse = {
        data: {
          id: 2,
          title: '买iPad',
          creator_id: 1,
          vote_count: 0,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createWish({ title: '买iPad' })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('title')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('vote_count')
    })
  })

  describe('updateWish', () => {
    it('should call PUT /wishes/:id with data', async () => {
      const updateData = { title: '买iPad Pro' }
      const mockResponse = { data: { id: 1, title: '买iPad Pro' } }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updateWish(1, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/wishes/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deleteWish', () => {
    it('should call DELETE /wishes/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteWish(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/wishes/1')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('promoteWish', () => {
    it('should call POST /wishes/:id/promote', async () => {
      const mockResponse = {
        data: {
          id: 1,
          type: 'family',
          status: 'pending',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await promoteWish(1)

      expect(mockRequest.post).toHaveBeenCalledWith('/wishes/1/promote')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys with type changed to family', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '买iPad',
          type: 'family',
          status: 'pending',
          creator_id: 1,
          vote_count: 3,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await promoteWish(1)
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('type')
      expect(item).toHaveProperty('status')
      expect(item.type).toBe('family')
    })
  })

  describe('updateWishStatus', () => {
    it('should call PUT /wishes/:id/status with status', async () => {
      const mockResponse = { data: { id: 1, status: 'approved' } }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updateWishStatus(1, 'approved')

      expect(mockRequest.put).toHaveBeenCalledWith('/wishes/1/status', { status: 'approved' })
      expect(result).toEqual(mockResponse)
    })
  })

  describe('voteWish', () => {
    it('should call POST /wishes/:id/vote', async () => {
      const mockResponse = {
        data: {
          id: 1,
          vote_count: 4,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await voteWish(1)

      expect(mockRequest.post).toHaveBeenCalledWith('/wishes/1/vote')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case vote_count in response', async () => {
      const mockResponse = {
        data: {
          id: 1,
          vote_count: 4,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await voteWish(1)
      const item = result.data

      expect(item).toHaveProperty('vote_count')
      expect(item.vote_count).toBe(4)

      // Verify it's snake_case, not PascalCase
      expect(item).not.toHaveProperty('VoteCount')
    })
  })

  describe('unvoteWish', () => {
    it('should call DELETE /wishes/:id/vote', async () => {
      const mockResponse = {
        data: {
          id: 1,
          vote_count: 2,
        },
      }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await unvoteWish(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/wishes/1/vote')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case vote_count in response', async () => {
      const mockResponse = {
        data: {
          id: 1,
          vote_count: 2,
        },
      }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await unvoteWish(1)
      const item = result.data

      expect(item).toHaveProperty('vote_count')
      expect(item.vote_count).toBe(2)
    })
  })
})
