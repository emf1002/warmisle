<template>
  <a-config-provider :theme="themeConfig">
    <component :is="layout">
      <router-view />
    </component>
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'

const route = useRoute()

const layoutMap: Record<string, any> = {
  auth: AuthLayout,
  main: MainLayout,
}

const layout = computed(() => {
  const layoutName = (route.meta.layout as string) || 'main'
  return layoutMap[layoutName] || MainLayout
})

const themeConfig = {
  token: {
    colorPrimary: '#1677FF',
    colorSuccess: '#52C41A',
    colorWarning: '#FAAD14',
    colorError: '#FF4D4F',
    borderRadius: 6,
    fontSize: 14,
  },
}
</script>
