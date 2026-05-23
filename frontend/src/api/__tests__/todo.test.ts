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

import { getTodoList, createTodo, updateTodo, deleteTodo, toggleTodo, claimTodo } from '../todo'

describe('Todo API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getTodoList', () => {
    it('should call GET /todos with params', async () => {
      const mockResponse = {
        data: {
          list: [
            {
              id: 1,
              title: '买菜',
              description: '去超市',
              priority: 'important',
              status: 'pending',
              assignee_id: 1,
              creator_id: 1,
              due_date: '2026-05-25',
              completed_at: null,
              created_at: '2026-05-23',
              assignee: { id: 1, name: '管理员', avatar: '👨' },
              creator: { id: 1, name: '管理员', avatar: '👨' },
            },
          ],
          total: 1,
          page: 1,
          page_size: 20,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTodoList({ page: 1, page_size: 20 })

      expect(mockRequest.get).toHaveBeenCalledWith('/todos', {
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
              title: '买菜',
              description: '去超市',
              priority: 'important',
              status: 'pending',
              assignee_id: 1,
              creator_id: 1,
              due_date: '2026-05-25',
              completed_at: null,
              created_at: '2026-05-23',
              assignee: { id: 1, name: '管理员', avatar: '👨' },
              creator: { id: 1, name: '管理员', avatar: '👨' },
            },
          ],
          total: 1,
          page: 1,
          page_size: 20,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTodoList({})
      const item = result.data.list[0]

      // Verify all snake_case keys exist
      expect(item).toHaveProperty('assignee_id')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('due_date')
      expect(item).toHaveProperty('completed_at')
      expect(item).toHaveProperty('created_at')

      // Verify nested objects have snake_case-compatible fields
      expect(item.assignee).toHaveProperty('avatar')
      expect(item.creator).toHaveProperty('avatar')
    })

    it('should NOT have PascalCase keys (AssigneeID would not match assignee_id)', async () => {
      const pascalCaseResponse = {
        data: {
          list: [
            {
              ID: 1,
              Title: '买菜',
              AssigneeID: 1,
              CreatorID: 1,
              DueDate: '2026-05-25',
              CompletedAt: null,
              CreatedAt: '2026-05-23',
            },
          ],
          Total: 1,
          Page: 1,
          PageSize: 20,
        },
      }
      mockRequest.get.mockResolvedValue(pascalCaseResponse)

      const result = await getTodoList({})
      const item = result.data.list[0]

      // PascalCase keys should NOT match the expected snake_case keys
      expect(item).not.toHaveProperty('assignee_id')
      expect(item).not.toHaveProperty('creator_id')
      expect(item).not.toHaveProperty('due_date')
      expect(item).not.toHaveProperty('completed_at')
      expect(item).not.toHaveProperty('created_at')

      // These PascalCase keys would exist in a wrongly formatted response
      expect(item).toHaveProperty('AssigneeID')
      expect(item).toHaveProperty('CreatorID')
      expect(item).toHaveProperty('DueDate')
    })
  })

  describe('createTodo', () => {
    it('should call POST /todos with data', async () => {
      const todoData = {
        title: '买菜',
        description: '去超市',
        priority: 'important',
        assignee_id: 1,
        due_date: '2026-05-25',
      }
      const mockResponse = {
        data: {
          id: 2,
          title: '买菜',
          priority: 'important',
          status: 'pending',
          assignee_id: 1,
          due_date: '2026-05-25',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTodo(todoData)

      expect(mockRequest.post).toHaveBeenCalledWith('/todos', todoData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created todo', async () => {
      const mockResponse = {
        data: {
          id: 2,
          title: '买菜',
          priority: 'important',
          status: 'pending',
          assignee_id: 1,
          due_date: '2026-05-25',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTodo({ title: '买菜' })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('title')
      expect(item).toHaveProperty('priority')
      expect(item).toHaveProperty('status')
      expect(item).toHaveProperty('assignee_id')
      expect(item).toHaveProperty('due_date')
    })
  })

  describe('updateTodo', () => {
    it('should call PUT /todos/:id with data', async () => {
      const updateData = { title: '买菜（已更新）' }
      const mockResponse = { data: { id: 1, title: '买菜（已更新）' } }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updateTodo(1, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/todos/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deleteTodo', () => {
    it('should call DELETE /todos/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteTodo(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/todos/1')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('toggleTodo', () => {
    it('should call PUT /todos/:id/toggle', async () => {
      const mockResponse = {
        data: {
          id: 1,
          status: 'completed',
          completed_at: '2026-05-23T10:00:00Z',
        },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await toggleTodo(1)

      expect(mockRequest.put).toHaveBeenCalledWith('/todos/1/toggle')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in toggled todo', async () => {
      const mockResponse = {
        data: {
          id: 1,
          status: 'completed',
          completed_at: '2026-05-23T10:00:00Z',
        },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await toggleTodo(1)
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('status')
      expect(item).toHaveProperty('completed_at')

      // Verify it's snake_case, not PascalCase
      expect(item).not.toHaveProperty('CompletedAt')
    })
  })

  describe('claimTodo', () => {
    it('should call PUT /todos/:id/claim', async () => {
      const mockResponse = { data: { id: 1, assignee_id: 2 } }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await claimTodo(1)

      expect(mockRequest.put).toHaveBeenCalledWith('/todos/1/claim')
      expect(result).toEqual(mockResponse)
    })
  })
})
