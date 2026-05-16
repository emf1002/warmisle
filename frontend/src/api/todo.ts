import request from '@/utils/request'

export function getTodoList(params: {
  status?: string
  assignee_id?: number
  page?: number
  page_size?: number
}) {
  return request.get('/todos', { params })
}

export function createTodo(data: {
  title: string
  description?: string
  priority?: string
  assignee_id?: number
  due_date?: string
}) {
  return request.post('/todos', data)
}

export function updateTodo(id: number, data: {
  title?: string
  description?: string
  priority?: string
  assignee_id?: number
  due_date?: string
}) {
  return request.put(`/todos/${id}`, data)
}

export function deleteTodo(id: number) {
  return request.delete(`/todos/${id}`)
}

export function toggleTodo(id: number) {
  return request.put(`/todos/${id}/toggle`)
}

export function claimTodo(id: number) {
  return request.put(`/todos/${id}/claim`)
}
