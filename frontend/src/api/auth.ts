import request from '@/utils/request'

export function login(username: string, password: string) {
  return request.post('/auth/login', { username, password })
}

export function checkInit() {
  return request.get('/init/check')
}

export function setupInit(admin_name: string, username: string, password: string) {
  return request.post('/init/setup', { admin_name, username, password })
}
