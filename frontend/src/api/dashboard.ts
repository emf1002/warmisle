import request from '@/utils/request'

export function getSummary(month?: string) {
  return request.get('/dashboard/summary', { params: { month } })
}

export function getExpenseChart(month?: string) {
  return request.get('/dashboard/expense-chart', { params: { month } })
}

export function getUpcomingTodos() {
  return request.get('/dashboard/upcoming-todos')
}

export function getWishTrends() {
  return request.get('/dashboard/wish-trends')
}

export function getForumHot() {
  return request.get('/dashboard/forum-hot')
}
