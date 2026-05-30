<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 32 32"
    xmlns="http://www.w3.org/2000/svg"
    class="logo-icon"
  >
    <defs>
      <linearGradient id="logo-badge-grad" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" :style="{ stopColor: isDark ? '#FF9F5A' : '#FF8C42' }" />
        <stop offset="100%" :style="{ stopColor: isDark ? '#FF8C42' : '#E85D2C' }" />
      </linearGradient>
      <linearGradient id="logo-bg-grad" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" :style="{ stopColor: isDark ? '#2D2D44' : '#FFF8F0' }" />
        <stop offset="100%" :style="{ stopColor: isDark ? '#1A1A2E' : '#FFF0E0' }" />
      </linearGradient>
    </defs>

    <!-- Badge background -->
    <rect x="1" y="1" width="30" height="30" rx="6" fill="url(#logo-bg-grad)" stroke="url(#logo-badge-grad)" stroke-width="1.5" />

    <!-- Sun rays -->
    <line x1="16" y1="4" x2="16" y2="7" stroke="#FFD93D" stroke-width="1.5" stroke-linecap="round" />
    <line x1="16" y1="25" x2="16" y2="28" stroke="#FFD93D" stroke-width="1.5" stroke-linecap="round" />
    <line x1="4" y1="16" x2="7" y2="16" stroke="#FFD93D" stroke-width="1.5" stroke-linecap="round" />
    <line x1="25" y1="16" x2="28" y2="16" stroke="#FFD93D" stroke-width="1.5" stroke-linecap="round" />

    <!-- House roof -->
    <polygon points="16,9 9,16 23,16" fill="url(#logo-badge-grad)" />

    <!-- House body -->
    <rect x="11" y="16" width="10" height="7" rx="0.5" fill="url(#logo-badge-grad)" />

    <!-- Door -->
    <rect x="14" y="18" width="4" height="5" rx="0.5" :fill="isDark ? '#2D2D44' : '#FFF8F0'" />

    <!-- Island wave -->
    <path d="M8,25 Q12,23 16,25 Q20,27 24,25" :stroke="isDark ? '#5DE8D4' : '#4ECDC4'" stroke-width="1.5" fill="none" stroke-linecap="round" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const props = withDefaults(defineProps<{
  size?: number
  theme?: 'light' | 'dark' | 'auto'
}>(), {
  size: 32,
  theme: 'auto'
})

const themeStore = useThemeStore()

const isDark = computed(() => {
  if (props.theme === 'dark') return true
  if (props.theme === 'light') return false
  return themeStore.theme === 'dark'
})
</script>