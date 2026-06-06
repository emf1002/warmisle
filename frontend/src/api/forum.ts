import request from '@/utils/request'

// --- Feed ---

export function getFeed(params: {
  page?: number
  page_size?: number
}) {
  return request.get('/feed', { params })
}

// --- Posts ---

export function createPost(data: {
  content: string
}) {
  return request.post('/posts', data)
}

export function updatePost(id: number, data: {
  content: string
}) {
  return request.put(`/posts/${id}`, data)
}

export function deletePost(id: number) {
  return request.delete(`/posts/${id}`)
}

// --- Topics ---

export function createTopic(data: {
  title: string
  content?: string
  tag_id?: number
}) {
  return request.post('/topics', data)
}

export function updateTopic(id: number, data: {
  title?: string
  content?: string
  tag_id?: number
}) {
  return request.put(`/topics/${id}`, data)
}

export function deleteTopic(id: number) {
  return request.delete(`/topics/${id}`)
}

export function togglePin(id: number) {
  return request.put(`/topics/${id}/pin`)
}

export function getTopic(id: number) {
  return request.get(`/topics/${id}`)
}

// --- Comments ---

export function getComments(params: {
  target_type: string
  target_id: number
}) {
  return request.get('/comments', { params })
}

export function createComment(data: {
  target_type: string
  target_id: number
  parent_id?: number
  content: string
}) {
  return request.post('/comments', data)
}

export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}

// --- Likes ---

export function toggleLike(data: {
  target_type: string
  target_id: number
}) {
  return request.post('/likes', data)
}

// --- Votes ---

export function listVotes(params: {
  page?: number
  page_size?: number
}) {
  return request.get('/votes', { params })
}

export function createVote(data: {
  title: string
  options: string[]
  is_multi?: boolean
  deadline?: string
}) {
  return request.post('/votes', data)
}

export function deleteVote(id: number) {
  return request.delete(`/votes/${id}`)
}

export function vote(id: number, data: {
  option_id: number
}) {
  return request.post(`/votes/${id}/vote`, data)
}

export function getVote(id: number) {
  return request.get(`/votes/${id}`)
}

// --- Tags ---

export function getTags() {
  return request.get('/tags')
}

export function createTag(data: {
  name: string
}) {
  return request.post('/tags', data)
}

export function updateTag(id: number, data: {
  name: string
}) {
  return request.put(`/tags/${id}`, data)
}

export function deleteTag(id: number) {
  return request.delete(`/tags/${id}`)
}
