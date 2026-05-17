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
    colorPrimary: '#E8734A',
    colorSuccess: '#389E0D',
    colorWarning: '#FAAD14',
    colorError: '#FF4D4F',
    colorInfo: '#E8734A',
    borderRadius: 8,
    fontSize: 14,
    fontFamily: '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif',
    colorBgLayout: '#FAF9F7',
    colorBgContainer: '#FFFFFF',
  },
  components: {
    Button: {
      borderRadius: 8,
      controlHeight: 44,
    },
    Input: {
      borderRadius: 8,
      controlHeight: 44,
    },
    Select: {
      borderRadius: 8,
      controlHeight: 44,
    },
  },
}
</script>
