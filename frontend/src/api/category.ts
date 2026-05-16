import request from '@/utils/request'

export function getCategories() {
  return request.get('/categories')
}

export function createCategory(data: {
  type: string
  name: string
  icon?: string
  sort_order?: number
}) {
  return request.post('/categories', data)
}

export function updateCategory(id: number, data: {
  type?: string
  name?: string
  icon?: string
  sort_order?: number
}) {
  return request.put(`/categories/${id}`, data)
}

export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}
