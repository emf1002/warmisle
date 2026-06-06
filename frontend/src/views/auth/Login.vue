<template>
  <div data-testid="login-page">
    <!-- 品牌区域 -->
    <div class="auth-brand">
      <div class="brand-logo">
        <LogoIcon :size="48" />
      </div>
      <h1 class="brand-title">暖屿</h1>
      <p class="brand-subtitle">WarmIsle · 温暖每一刻</p>
    </div>

    <!-- 表单区域 -->
    <a-form
      :model="form"
      :rules="rules"
      layout="vertical"
      @finish="handleLogin"
      class="auth-form"
    >
      <a-form-item name="username" class="form-item">
        <a-input
          v-model:value="form.username"
          placeholder="请输入用户名"
          size="large"
          :min-height="44"
          data-testid="username-input"
          class="auth-input"
        >
          <template #prefix>
            <Icon name="User" :size="18" class="input-icon" />
          </template>
        </a-input>
      </a-form-item>

      <a-form-item name="password" class="form-item">
        <a-input-password
          v-model:value="form.password"
          placeholder="请输入密码"
          size="large"
          :min-height="44"
          data-testid="password-input"
          class="auth-input"
        >
          <template #prefix>
            <Icon name="Lock" :size="18" class="input-icon" />
          </template>
        </a-input-password>
      </a-form-item>

      <div class="form-options">
        <a-checkbox v-model:checked="rememberMe" class="remember-checkbox">
          记住我
        </a-checkbox>
      </div>

      <a-form-item class="form-submit">
        <a-button
          type="primary"
          html-type="submit"
          size="large"
          :loading="loading"
          block
          data-testid="login-btn"
          class="login-btn"
        >
          <template #icon>
            <Icon name="LogIn" :size="18" />
          </template>
          登 录
        </a-button>
      </a-form-item>
    </a-form>

    <!-- 底部提示 -->
    <div class="auth-footer">
      <p class="footer-text">
        <Icon name="Shield" :size="14" />
        <span>私有化部署 · 数据安全有保障</span>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Icon from '@/components/Icon.vue'
import LogoIcon from '@/components/LogoIcon.vue'

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
/* ==================== 品牌区域 ==================== */
.auth-brand {
  text-align: center;
  margin-bottom: var(--space-xl, 32px);
  animation: fadeInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) 100ms both;
}

.brand-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: var(--color-brand-bg);
  border-radius: var(--radius-xl, 20px);
  margin-bottom: var(--space-md, 16px);
  transition: transform var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1);
}

.brand-logo:hover {
  transform: scale(1.05) rotate(5deg);
}

.brand-title {
  font-size: var(--text-3xl, 1.875rem);
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 var(--space-xs, 8px) 0;
  letter-spacing: 2px;
  font-family: var(--font-display);
}

.brand-subtitle {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
  letter-spacing: 1px;
}

/* ==================== 表单样式 ==================== */
.auth-form {
  animation: fadeInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) 200ms both;
}

.form-item {
  margin-bottom: var(--space-md, 16px);
}

:deep(.auth-input) {
  border-radius: var(--radius-md, 12px) !important;
  border: 1px solid var(--color-border) !important;
  transition: all var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1) !important;
}

:deep(.auth-input:hover) {
  border-color: var(--color-brand) !important;
}

:deep(.auth-input:focus),
:deep(.auth-input-focused) {
  border-color: var(--color-brand) !important;
  box-shadow: 0 0 0 3px var(--color-brand-light) !important;
}

:deep(.auth-input .ant-input-prefix) {
  margin-right: var(--space-sm, 12px);
}

.input-icon {
  color: var(--color-text-tertiary, #B8AFA8);
  transition: color var(--duration-fast, 150ms) ease-out;
}

:deep(.auth-input:hover) .input-icon,
:deep(.auth-input-focused) .input-icon {
  color: var(--color-brand);
}

/* ==================== 表单选项 ==================== */
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-lg, 24px);
}

.remember-checkbox {
  color: var(--color-text-secondary);
  font-size: var(--text-sm, 0.875rem);
}

/* ==================== 登录按钮 ==================== */
.form-submit {
  margin-bottom: 0;
}

.login-btn {
  height: 48px !important;
  border-radius: var(--radius-md, 12px) !important;
  font-size: var(--text-base, 1rem) !important;
  font-weight: 600 !important;
  letter-spacing: 4px !important;
  background: var(--gradient-brand, linear-gradient(135deg, #E87461 0%, #F2B84B 100%)) !important;
  border: none !important;
  box-shadow: 0 4px 12px rgba(232, 116, 97, 0.3) !important;
  transition: all var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(232, 116, 97, 0.4) !important;
}

.login-btn:active:not(:disabled) {
  transform: translateY(0) !important;
  box-shadow: 0 2px 8px rgba(232, 116, 97, 0.3) !important;
}

/* ==================== 底部提示 ==================== */
.auth-footer {
  margin-top: var(--space-lg, 24px);
  text-align: center;
  animation: fadeInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) 300ms both;
}

.footer-text {
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-text-tertiary, #B8AFA8);
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs, 8px);
  margin: 0;
}

/* ==================== 动画 ==================== */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ==================== 响应式 ==================== */
@media (max-width: 480px) {
  .brand-logo {
    width: 64px;
    height: 64px;
  }

  .brand-title {
    font-size: var(--text-2xl, 1.5rem);
  }
}
</style>
