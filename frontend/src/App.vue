<template>
  <component :is="layout">
    <router-view />
  </component>
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
</script>
