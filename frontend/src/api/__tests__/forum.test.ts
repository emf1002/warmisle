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
  getFeed,
  createPost,
  updatePost,
  deletePost,
  createTopic,
  updateTopic,
  deleteTopic,
  togglePin,
  getTopic,
  getComments,
  createComment,
  deleteComment,
  toggleLike,
  createVote,
  deleteVote,
  vote,
  getVote,
  getTags,
  createTag,
  updateTag,
  deleteTag,
} from '../forum'

describe('Forum API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  // --- Feed ---

  describe('getFeed', () => {
    it('should call GET /feed with params', async () => {
      const mockResponse = {
        data: {
          pinned: [
            {
              type: 'topic',
              id: 1,
              title: '置顶话题',
              content: '内容',
              creator: { id: 1, name: '管理员', avatar: '👨' },
              tag: { id: 1, name: '讨论', preset: true },
              is_pinned: true,
              created_at: '2026-05-23',
            },
          ],
          items: [
            {
              type: 'post',
              id: 1,
              title: '',
              content: '今天天气好',
              creator: { id: 1, name: '管理员', avatar: '👨' },
              is_pinned: false,
              is_liked: false,
              like_count: 0,
              comment_count: 0,
              created_at: '2026-05-23',
            },
          ],
          total: 2,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getFeed({ page: 1, page_size: 10 })

      expect(mockRequest.get).toHaveBeenCalledWith('/feed', {
        params: { page: 1, page_size: 10 },
      })
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in feed response', async () => {
      const mockResponse = {
        data: {
          pinned: [
            {
              type: 'topic',
              id: 1,
              title: '置顶话题',
              content: '内容',
              creator: { id: 1, name: '管理员', avatar: '👨' },
              tag: { id: 1, name: '讨论', preset: true },
              is_pinned: true,
              created_at: '2026-05-23',
            },
          ],
          items: [
            {
              type: 'post',
              id: 1,
              title: '',
              content: '今天天气好',
              creator: { id: 1, name: '管理员', avatar: '👨' },
              is_pinned: false,
              is_liked: false,
              like_count: 0,
              comment_count: 0,
              created_at: '2026-05-23',
            },
          ],
          total: 2,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getFeed({ page: 1, page_size: 10 })
      const pinned = result.data.pinned[0]
      const item = result.data.items[0]

      // Verify pinned topic snake_case keys
      expect(pinned).toHaveProperty('is_pinned')
      expect(pinned).toHaveProperty('created_at')
      expect(pinned.creator).toHaveProperty('avatar')
      expect(pinned.tag).toHaveProperty('preset')

      // Verify feed item snake_case keys
      expect(item).toHaveProperty('is_pinned')
      expect(item).toHaveProperty('is_liked')
      expect(item).toHaveProperty('like_count')
      expect(item).toHaveProperty('comment_count')
      expect(item).toHaveProperty('created_at')
      expect(item.creator).toHaveProperty('avatar')
    })
  })

  // --- Post CRUD ---

  describe('createPost', () => {
    it('should call POST /posts with data', async () => {
      const postData = { content: '今天天气好' }
      const mockResponse = {
        data: {
          id: 1,
          content: '今天天气好',
          creator_id: 1,
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createPost(postData)

      expect(mockRequest.post).toHaveBeenCalledWith('/posts', postData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created post', async () => {
      const mockResponse = {
        data: {
          id: 1,
          content: '今天天气好',
          creator_id: 1,
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createPost({ content: '今天天气好' })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('content')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('created_at')
    })
  })

  describe('updatePost', () => {
    it('should call PUT /posts/:id with data', async () => {
      const updateData = { content: '修改后的内容' }
      const mockResponse = {
        data: {
          id: 1,
          content: '修改后的内容',
          creator_id: 1,
          created_at: '2026-05-23',
        },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updatePost(1, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/posts/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deletePost', () => {
    it('should call DELETE /posts/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deletePost(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/posts/1')
      expect(result).toEqual(mockResponse)
    })
  })

  // --- Topic CRUD ---

  describe('createTopic', () => {
    it('should call POST /topics with data', async () => {
      const topicData = { title: '新话题', content: '话题内容', tag_id: 1 }
      const mockResponse = {
        data: {
          id: 1,
          title: '新话题',
          content: '话题内容',
          tag_id: 1,
          creator_id: 1,
          is_pinned: false,
          tag: { id: 1, name: '讨论', preset: true },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTopic(topicData)

      expect(mockRequest.post).toHaveBeenCalledWith('/topics', topicData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created topic', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '新话题',
          content: '话题内容',
          tag_id: 1,
          creator_id: 1,
          is_pinned: false,
          tag: { id: 1, name: '讨论', preset: true },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTopic({ title: '新话题' })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('title')
      expect(item).toHaveProperty('content')
      expect(item).toHaveProperty('tag_id')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('is_pinned')
      expect(item.tag).toHaveProperty('name')
      expect(item.creator).toHaveProperty('avatar')
    })
  })

  describe('updateTopic', () => {
    it('should call PUT /topics/:id with data', async () => {
      const updateData = { title: '更新话题' }
      const mockResponse = {
        data: { id: 1, title: '更新话题', tag_id: 1, creator_id: 1, is_pinned: false },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updateTopic(1, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/topics/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deleteTopic', () => {
    it('should call DELETE /topics/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteTopic(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/topics/1')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('togglePin', () => {
    it('should call PUT /topics/:id/pin', async () => {
      const mockResponse = {
        data: { id: 1, is_pinned: true },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await togglePin(1)

      expect(mockRequest.put).toHaveBeenCalledWith('/topics/1/pin')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case is_pinned in response', async () => {
      const mockResponse = {
        data: { id: 1, is_pinned: true },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await togglePin(1)
      const item = result.data

      expect(item).toHaveProperty('is_pinned')
      expect(item.is_pinned).toBe(true)
    })
  })

  describe('getTopic', () => {
    it('should call GET /topics/:id', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '话题详情',
          content: '详细内容',
          tag_id: 1,
          creator_id: 1,
          is_pinned: true,
          tag: { id: 1, name: '讨论', preset: true },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          created_at: '2026-05-23',
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTopic(1)

      expect(mockRequest.get).toHaveBeenCalledWith('/topics/1')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys with nested tag and creator', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '话题详情',
          content: '详细内容',
          tag_id: 1,
          creator_id: 1,
          is_pinned: true,
          tag: { id: 1, name: '讨论', preset: true },
          creator: { id: 1, name: '管理员', avatar: '👨' },
          created_at: '2026-05-23',
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTopic(1)
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('title')
      expect(item).toHaveProperty('content')
      expect(item).toHaveProperty('tag_id')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('is_pinned')
      expect(item).toHaveProperty('created_at')
      expect(item.tag).toHaveProperty('name')
      expect(item.tag).toHaveProperty('preset')
      expect(item.creator).toHaveProperty('avatar')
    })
  })

  // --- Comments ---

  describe('getComments', () => {
    it('should call GET /comments with params', async () => {
      const mockResponse = {
        data: {
          list: [
            {
              id: 1,
              target_type: 'post',
              target_id: 1,
              parent_id: null,
              content: '好帖子',
              creator_id: 1,
              creator: { id: 1, name: '管理员', avatar: '👨' },
              created_at: '2026-05-23',
            },
          ],
          total: 1,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getComments({ target_type: 'post', target_id: 1 })

      expect(mockRequest.get).toHaveBeenCalledWith('/comments', {
        params: { target_type: 'post', target_id: 1 },
      })
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in comments', async () => {
      const mockResponse = {
        data: {
          list: [
            {
              id: 1,
              target_type: 'post',
              target_id: 1,
              parent_id: null,
              content: '好帖子',
              creator_id: 1,
              creator: { id: 1, name: '管理员', avatar: '👨' },
            },
          ],
          total: 1,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getComments({ target_type: 'post', target_id: 1 })
      const comment = result.data.list[0]

      expect(comment).toHaveProperty('id')
      expect(comment).toHaveProperty('target_type')
      expect(comment).toHaveProperty('target_id')
      expect(comment).toHaveProperty('parent_id')
      expect(comment).toHaveProperty('content')
      expect(comment).toHaveProperty('creator_id')
      expect(comment.creator).toHaveProperty('name')
    })
  })

  describe('createComment', () => {
    it('should call POST /comments with data', async () => {
      const commentData = {
        target_type: 'post',
        target_id: 1,
        content: '好帖子',
      }
      const mockResponse = {
        data: {
          id: 1,
          target_type: 'post',
          target_id: 1,
          parent_id: null,
          content: '好帖子',
          creator_id: 1,
          creator: { id: 1, name: '管理员', avatar: '👨' },
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createComment(commentData)

      expect(mockRequest.post).toHaveBeenCalledWith('/comments', commentData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created comment', async () => {
      const mockResponse = {
        data: {
          id: 1,
          target_type: 'post',
          target_id: 1,
          parent_id: null,
          content: '好帖子',
          creator_id: 1,
          creator: { id: 1, name: '管理员' },
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createComment({
        target_type: 'post',
        target_id: 1,
        content: '好帖子',
      })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('target_type')
      expect(item).toHaveProperty('target_id')
      expect(item).toHaveProperty('parent_id')
      expect(item).toHaveProperty('content')
      expect(item).toHaveProperty('creator_id')
      expect(item.creator).toHaveProperty('name')
    })
  })

  describe('deleteComment', () => {
    it('should call DELETE /comments/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteComment(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/comments/1')
      expect(result).toEqual(mockResponse)
    })
  })

  // --- Likes ---

  describe('toggleLike', () => {
    it('should call POST /likes with data', async () => {
      const likeData = { target_type: 'post', target_id: 1 }
      const mockResponse = {
        data: {
          id: 1,
          target_type: 'post',
          target_id: 1,
          member_id: 1,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await toggleLike(likeData)

      expect(mockRequest.post).toHaveBeenCalledWith('/likes', likeData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in like response', async () => {
      const mockResponse = {
        data: {
          id: 1,
          target_type: 'post',
          target_id: 1,
          member_id: 1,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await toggleLike({ target_type: 'post', target_id: 1 })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('target_type')
      expect(item).toHaveProperty('target_id')
      expect(item).toHaveProperty('member_id')
    })
  })

  // --- Votes ---

  describe('createVote', () => {
    it('should call POST /votes with data', async () => {
      const voteData = {
        title: '今晚吃什么',
        options: ['火锅', '烧烤', '外卖'],
        is_multi: false,
        deadline: '2026-05-24T18:00:00Z',
      }
      const mockResponse = {
        data: {
          id: 1,
          title: '今晚吃什么',
          is_multi: false,
          deadline: '2026-05-24T18:00:00Z',
          creator_id: 1,
          options: [
            { id: 1, vote_id: 1, content: '火锅', sort_order: 0 },
            { id: 2, vote_id: 1, content: '烧烤', sort_order: 1 },
            { id: 3, vote_id: 1, content: '外卖', sort_order: 2 },
          ],
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createVote(voteData)

      expect(mockRequest.post).toHaveBeenCalledWith('/votes', voteData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created vote', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '今晚吃什么',
          is_multi: false,
          deadline: '2026-05-24T18:00:00Z',
          creator_id: 1,
          options: [
            { id: 1, vote_id: 1, content: '火锅', sort_order: 0 },
          ],
          created_at: '2026-05-23',
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createVote({
        title: '今晚吃什么',
        options: ['火锅'],
      })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('title')
      expect(item).toHaveProperty('is_multi')
      expect(item).toHaveProperty('deadline')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('created_at')
      expect(item.options[0]).toHaveProperty('vote_id')
      expect(item.options[0]).toHaveProperty('sort_order')
    })
  })

  describe('deleteVote', () => {
    it('should call DELETE /votes/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteVote(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/votes/1')
      expect(result).toEqual(mockResponse)
    })
  })

  describe('vote', () => {
    it('should call POST /votes/:id/vote with option_id', async () => {
      const mockResponse = {
        data: {
          id: 1,
          option_id: 2,
          member_id: 1,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await vote(1, { option_id: 2 })

      expect(mockRequest.post).toHaveBeenCalledWith('/votes/1/vote', { option_id: 2 })
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in vote record', async () => {
      const mockResponse = {
        data: {
          id: 1,
          option_id: 2,
          member_id: 1,
        },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await vote(1, { option_id: 2 })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('option_id')
      expect(item).toHaveProperty('member_id')
    })
  })

  describe('getVote', () => {
    it('should call GET /votes/:id', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '今晚吃什么',
          is_multi: false,
          deadline: '2026-05-24T18:00:00Z',
          creator_id: 1,
          options: [
            { id: 1, vote_id: 1, content: '火锅', sort_order: 0 },
            { id: 2, vote_id: 1, content: '烧烤', sort_order: 1 },
          ],
          created_at: '2026-05-23',
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getVote(1)

      expect(mockRequest.get).toHaveBeenCalledWith('/votes/1')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys with nested options', async () => {
      const mockResponse = {
        data: {
          id: 1,
          title: '今晚吃什么',
          is_multi: false,
          deadline: '2026-05-24T18:00:00Z',
          creator_id: 1,
          options: [
            { id: 1, vote_id: 1, content: '火锅', sort_order: 0 },
          ],
          created_at: '2026-05-23',
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getVote(1)
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('is_multi')
      expect(item).toHaveProperty('deadline')
      expect(item).toHaveProperty('creator_id')
      expect(item).toHaveProperty('created_at')
      expect(item.options[0]).toHaveProperty('vote_id')
      expect(item.options[0]).toHaveProperty('sort_order')
    })
  })

  // --- Tags ---

  describe('getTags', () => {
    it('should call GET /tags', async () => {
      const mockResponse = {
        data: {
          list: [
            { id: 1, name: '讨论', preset: true },
            { id: 2, name: '公告', preset: false },
          ],
          total: 2,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTags()

      expect(mockRequest.get).toHaveBeenCalledWith('/tags')
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in tags', async () => {
      const mockResponse = {
        data: {
          list: [
            { id: 1, name: '讨论', preset: true },
          ],
          total: 1,
        },
      }
      mockRequest.get.mockResolvedValue(mockResponse)

      const result = await getTags()
      const tag = result.data.list[0]

      expect(tag).toHaveProperty('id')
      expect(tag).toHaveProperty('name')
      expect(tag).toHaveProperty('preset')
    })
  })

  describe('createTag', () => {
    it('should call POST /tags with data', async () => {
      const tagData = { name: '新标签' }
      const mockResponse = {
        data: { id: 3, name: '新标签', preset: false },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTag(tagData)

      expect(mockRequest.post).toHaveBeenCalledWith('/tags', tagData)
      expect(result).toEqual(mockResponse)
    })

    it('should return snake_case keys in created tag', async () => {
      const mockResponse = {
        data: { id: 3, name: '新标签', preset: false },
      }
      mockRequest.post.mockResolvedValue(mockResponse)

      const result = await createTag({ name: '新标签' })
      const item = result.data

      expect(item).toHaveProperty('id')
      expect(item).toHaveProperty('name')
      expect(item).toHaveProperty('preset')
    })
  })

  describe('updateTag', () => {
    it('should call PUT /tags/:id with data', async () => {
      const updateData = { name: '更新标签' }
      const mockResponse = {
        data: { id: 1, name: '更新标签', preset: true },
      }
      mockRequest.put.mockResolvedValue(mockResponse)

      const result = await updateTag(1, updateData)

      expect(mockRequest.put).toHaveBeenCalledWith('/tags/1', updateData)
      expect(result).toEqual(mockResponse)
    })
  })

  describe('deleteTag', () => {
    it('should call DELETE /tags/:id', async () => {
      const mockResponse = { data: null }
      mockRequest.delete.mockResolvedValue(mockResponse)

      const result = await deleteTag(1)

      expect(mockRequest.delete).toHaveBeenCalledWith('/tags/1')
      expect(result).toEqual(mockResponse)
    })
  })

  // --- Negative Tests: PascalCase keys would be undefined ---

  describe('Negative tests: PascalCase keys would be undefined', () => {
    it('should NOT match PascalCase keys for feed items (IsPinned, CreatorID, LikeCount)', async () => {
      const pascalCaseResponse = {
        data: {
          pinned: [],
          items: [
            {
              Type: 'post',
              ID: 1,
              Title: '',
              Content: '今天天气好',
              Creator: { ID: 1, Name: '管理员', Avatar: '👨' },
              IsPinned: false,
              IsLiked: false,
              LikeCount: 0,
              CommentCount: 0,
              CreatedAt: '2026-05-23',
            },
          ],
          Total: 1,
        },
      }
      mockRequest.get.mockResolvedValue(pascalCaseResponse)

      const result = await getFeed({ page: 1, page_size: 10 })
      const item = result.data.items[0]

      // PascalCase keys should NOT match expected snake_case access patterns
      expect(item).not.toHaveProperty('is_pinned')
      expect(item).not.toHaveProperty('is_liked')
      expect(item).not.toHaveProperty('like_count')
      expect(item).not.toHaveProperty('comment_count')
      expect(item).not.toHaveProperty('created_at')

      // PascalCase keys exist in wrongly formatted response
      expect(item).toHaveProperty('IsPinned')
      expect(item).toHaveProperty('IsLiked')
      expect(item).toHaveProperty('LikeCount')
      expect(item).toHaveProperty('CommentCount')
    })

    it('should NOT match PascalCase keys for comments (TargetType, CreatorID)', async () => {
      const pascalCaseResponse = {
        data: {
          list: [
            {
              ID: 1,
              TargetType: 'post',
              TargetID: 1,
              ParentID: null,
              Content: '好帖子',
              CreatorID: 1,
              Creator: { ID: 1, Name: '管理员' },
            },
          ],
          Total: 1,
        },
      }
      mockRequest.get.mockResolvedValue(pascalCaseResponse)

      const result = await getComments({ target_type: 'post', target_id: 1 })
      const comment = result.data.list[0]

      // PascalCase keys should NOT match expected snake_case access patterns
      expect(comment).not.toHaveProperty('target_type')
      expect(comment).not.toHaveProperty('target_id')
      expect(comment).not.toHaveProperty('parent_id')
      expect(comment).not.toHaveProperty('creator_id')

      // PascalCase keys exist in wrongly formatted response
      expect(comment).toHaveProperty('TargetType')
      expect(comment).toHaveProperty('TargetID')
      expect(comment).toHaveProperty('ParentID')
      expect(comment).toHaveProperty('CreatorID')
    })
  })
})
