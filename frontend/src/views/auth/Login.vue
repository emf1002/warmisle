<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>暖屿 · WarmIsle</h1>
        <p>登录您的账号</p>
      </div>
      <a-form
        :model="form"
        :rules="rules"
        layout="vertical"
        @finish="handleLogin"
      >
        <a-form-item name="username">
          <a-input
            v-model:value="form.username"
            placeholder="用户名"
            size="large"
            :min-height="44"
            data-testid="username-input"
          />
        </a-form-item>
        <a-form-item name="password">
          <a-input-password
            v-model:value="form.password"
            placeholder="密码"
            size="large"
            :min-height="44"
            data-testid="password-input"
          />
        </a-form-item>
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            :loading="loading"
            block
            data-testid="login-btn"
          >
            登 录
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

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }]
}

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login(form.username, form.password)
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
  width: 100%;
}

.auth-card {
  width: 100%;
  background: transparent;
  border: none;
  box-shadow: none;
  padding: 0;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text-primary);
}

.auth-header p {
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
