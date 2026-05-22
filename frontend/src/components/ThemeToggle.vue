<!-- frontend/src/components/ThemeToggle.vue -->
<template>
  <button
    class="theme-toggle"
    :title="isDark ? '切换到亮色主题' : '切换到暗色主题'"
    @click="themeStore.toggleTheme()"
  >
    <span class="theme-toggle-icon">{{ isDark ? '🌙' : '☀️' }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.theme === 'dark')
</script>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--duration-fast) ease;
  -webkit-tap-highlight-color: transparent;
}

.theme-toggle:hover {
  background: var(--sidebar-hover-bg);
}

.theme-toggle-icon {
  font-size: 18px;
  line-height: 1;
}
</style>
