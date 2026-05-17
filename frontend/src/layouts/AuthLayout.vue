<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-header">
        <span class="auth-logo">🏠</span>
        <h1 class="auth-title">家庭数字中心</h1>
        <p class="auth-subtitle">全家人的记账、待办、愿望、交流</p>
      </div>
      <div class="auth-content">
        <slot />
      </div>
      <div class="auth-footer">
        <p>&copy; 2026 家庭数字中心</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.auth-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  min-height: 100dvh;
  background: linear-gradient(135deg, #E8734A 0%, #F0884A 50%, #F5A67A 100%);
  padding: 16px;
}

.auth-container {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-logo {
  font-size: 64px;
  display: block;
  margin-bottom: 12px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.auth-title {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
  letter-spacing: 2px;
}

.auth-subtitle {
  font-size: 14px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.75);
  margin: 0;
}

.auth-content {
  width: 100%;
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-level-4);
  padding: 32px;
}

.auth-footer {
  margin-top: 24px;
}

.auth-footer p {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}
</style>
