import request from '@/utils/request'

export function getConfig() {
  return request.get('/backup/config')
}

export function saveConfig(data: { app_id: string; app_secret: string; redirect_uri: string; backup_dir: string }) {
  return request.put('/backup/config', data)
}

export function getAuthUrl() {
  return request.get('/backup/auth-url')
}

export function callback(data: { code: string; state: string }) {
  return request.post('/backup/callback', data)
}

export function triggerBackup() {
  return request.post('/backup/trigger')
}

export function listCloudFiles() {
  return request.get('/backup/cloud-files')
}

export function restoreBackup(cloudFileId: string, confirmText: string) {
  return request.post(`/backup/restore/${cloudFileId}`, { confirm_text: confirmText })
}

export function listHistory(params: { page: number; page_size: number }) {
  return request.get('/backup/history', { params })
}

export function deleteHistory(id: number) {
  return request.delete(`/backup/history/${id}`)
}

export function getSchedule() {
  return request.get('/backup/schedule')
}

export function saveSchedule(data: { schedule_enabled: boolean; schedule_time: string; retention_days: number }) {
  return request.put('/backup/schedule', data)
}
