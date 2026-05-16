<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-header">
        <span class="auth-logo">🏠</span>
        <h1 class="auth-title">家庭数字中心</h1>
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
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
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
  font-size: 48px;
  display: block;
  margin-bottom: 8px;
}

.auth-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
  letter-spacing: 1px;
}

.auth-content {
  width: 100%;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  padding: 32px;
}

.auth-footer {
  margin-top: 24px;
}

.auth-footer p {
  font-size: 12px;
  color: #bbb;
  margin: 0;
}
</style>
