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
            placeholder="请输入用户名"
            size="large"
            :min-height="44"
            data-testid="username-input"
          >
            <template #prefix><Icon name="User" :size="16" color="var(--color-muted)" /></template>
          </a-input>
        </a-form-item>
        <a-form-item name="password">
          <a-input-password
            v-model:value="form.password"
            placeholder="请输入密码"
            size="large"
            :min-height="44"
            data-testid="password-input"
          >
            <template #prefix><Icon name="Lock" :size="16" color="var(--color-muted)" /></template>
          </a-input-password>
        </a-form-item>
        <div class="login-options">
          <a-checkbox v-model:checked="rememberMe">记住我</a-checkbox>
        </div>
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
import Icon from '@/components/Icon.vue'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const rememberMe = ref(false)
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
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  box-shadow: var(--shadow-level-2);
}

:root[data-theme="dark"] .auth-card {
  background: rgba(45, 42, 38, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.08);
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

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}
</style>
