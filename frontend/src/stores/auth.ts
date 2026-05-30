import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, checkInit as checkInitApi } from '@/api/auth'

function parseJwtPayload(rawToken: string): Record<string, any> | null {
  if (!rawToken) return null
  try {
    return JSON.parse(atob(rawToken.split('.')[1]))
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const memberInfo = ref<any>(null)

  const currentUserId = computed(() => {
    const payload = parseJwtPayload(token.value)
    if (!payload) return 0
    return (payload.member_id as number) || 0
  })

  const currentUserRole = computed(() => {
    const payload = parseJwtPayload(token.value)
    if (!payload) return ''
    return payload.role || ''
  })

  const isAdmin = computed(() => currentUserRole.value === 'admin')

  async function login(username: string, password: string) {
    const res: any = await loginApi(username, password)
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    if (res.data.member) {
      memberInfo.value = res.data.member
    }
    return res
  }

  function logout() {
    token.value = ''
    memberInfo.value = null
    localStorage.removeItem('token')
  }

  async function checkInit() {
    const res: any = await checkInitApi()
    return res.data.need_init
  }

  return { token, memberInfo, currentUserId, currentUserRole, isAdmin, login, logout, checkInit }
})
