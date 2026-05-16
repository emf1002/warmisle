import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as loginApi, checkInit as checkInitApi } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const memberInfo = ref<any>(null)

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

  return { token, memberInfo, login, logout, checkInit }
})
