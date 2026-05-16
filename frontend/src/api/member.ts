import request from '@/utils/request'

export function getMembers() {
  return request.get('/members')
}

export function createMember(data: {
  username: string
  password: string
  name?: string
  avatar?: string
  role?: string
}) {
  return request.post('/members', data)
}

export function updateMember(id: number, data: {
  name?: string
  avatar?: string
  role?: string
}) {
  return request.put(`/members/${id}`, data)
}

export function deleteMember(id: number) {
  return request.delete(`/members/${id}`)
}

export function disableMember(id: number) {
  return request.put(`/members/${id}/disable`)
}

export function enableMember(id: number) {
  return request.put(`/members/${id}/enable`)
}

export function resetPassword(id: number) {
  return request.put(`/members/${id}/reset-pwd`)
}

export function getProfile() {
  return request.get('/profile')
}

export function updateProfile(data: { name?: string; avatar?: string }) {
  return request.put('/profile', data)
}
