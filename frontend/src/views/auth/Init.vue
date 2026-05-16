<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>初始化家庭数字中心</h1>
        <p>创建管理员账号</p>
      </div>
      <a-form
        :model="form"
        :rules="rules"
        layout="vertical"
        @finish="handleSetup"
      >
        <a-form-item name="name">
          <a-input
            v-model:value="form.name"
            placeholder="姓名"
            size="large"
            :min-height="44"
          />
        </a-form-item>
        <a-form-item name="username">
          <a-input
            v-model:value="form.username"
            placeholder="用户名，3-20位字母数字下划线"
            size="large"
            :min-height="44"
          />
        </a-form-item>
        <a-form-item name="password">
          <a-input-password
            v-model:value="form.password"
            placeholder="密码，6-32位"
            size="large"
            :min-height="44"
          />
        </a-form-item>
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            :loading="loading"
            block
          >
            开始使用
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { setupInit } from '@/api/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  name: '',
  username: '',
  password: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名' }],
  username: [
    { required: true, message: '请输入用户名' },
    { pattern: /^[a-zA-Z0-9_]{3,20}$/, message: '用户名需3-20位字母数字下划线' }
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, max: 32, message: '密码需6-32位' }
  ]
}

async function handleSetup() {
  loading.value = true
  try {
    const res: any = await setupInit(form.name, form.username, form.password)
    authStore.token = res.data.token
    localStorage.setItem('token', res.data.token)
    if (res.data.member) {
      authStore.memberInfo = res.data.member
    }
    router.push('/')
  } catch {
    // Error already handled by request interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  padding: 16px;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 8px;
  padding: 40px 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 24px;
  margin: 0 0 8px;
  color: #1a1a1a;
}

.auth-header p {
  color: #666;
  margin: 0;
}
</style>
