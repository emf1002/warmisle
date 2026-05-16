import request from '@/utils/request'

export function getWishList(params: {
  type?: string
  status?: string
  creator_id?: number
  page?: number
  page_size?: number
}) {
  return request.get('/wishes', { params })
}

export function createWish(data: {
  title: string
  description?: string
  category?: string
  priority?: string
  amount?: number
}) {
  return request.post('/wishes', data)
}

export function updateWish(id: number, data: {
  title?: string
  description?: string
  category?: string
  priority?: string
  amount?: number
}) {
  return request.put(`/wishes/${id}`, data)
}

export function deleteWish(id: number) {
  return request.delete(`/wishes/${id}`)
}

export function promoteWish(id: number) {
  return request.post(`/wishes/${id}/promote`)
}

export function updateWishStatus(id: number, status: string) {
  return request.put(`/wishes/${id}/status`, { status })
}

export function voteWish(id: number) {
  return request.post(`/wishes/${id}/vote`)
}

export function unvoteWish(id: number) {
  return request.delete(`/wishes/${id}/vote`)
}
