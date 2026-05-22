<template>
  <a-config-provider :theme="themeConfig">
    <component :is="layout">
      <router-view v-slot="{ Component: RouteComponent }">
        <transition name="page-fade" mode="out-in">
          <component :is="RouteComponent" />
        </transition>
      </router-view>
    </component>
  </a-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { theme } from 'ant-design-vue'
import MainLayout from '@/layouts/MainLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const themeStore = useThemeStore()

const layoutMap: Record<string, any> = {
  auth: AuthLayout,
  main: MainLayout,
}

const layout = computed(() => {
  const layoutName = (route.meta.layout as string) || 'main'
  return layoutMap[layoutName] || MainLayout
})

const themeConfig = computed(() => {
  const isDark = themeStore.theme === 'dark'
  return {
    algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: isDark ? '#E5A84B' : '#E87461',
      colorSuccess: isDark ? '#52C41A' : '#389E0D',
      colorWarning: '#FAAD14',
      colorError: isDark ? '#FF6B6B' : '#E85D5D',
      colorInfo: isDark ? '#E5A84B' : '#E87461',
      borderRadius: isDark ? 10 : 12,
      fontSize: 14,
      fontFamily: '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif',
      colorBgLayout: isDark ? '#0F1729' : '#FDF8F4',
      colorBgContainer: isDark ? '#1A2340' : '#FFFFFF',
      colorBgElevated: isDark ? '#1A2340' : '#FFFFFF',
      colorText: isDark ? '#E8E4DE' : '#3D3530',
      colorTextSecondary: isDark ? '#8A94A8' : '#8A807A',
      colorBorder: isDark ? 'rgba(255,255,255,0.06)' : 'rgba(61,53,48,0.08)',
      colorBorderSecondary: isDark ? 'rgba(255,255,255,0.04)' : 'rgba(61,53,48,0.05)',
    },
    components: {
      Button: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Input: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Select: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Card: {
        borderRadiusLG: isDark ? 12 : 14,
      },
      Modal: {
        borderRadiusLG: isDark ? 12 : 14,
      },
    },
  }
})
</script>
