<script setup lang="ts">
withDefaults(defineProps<{
  width?: string
  height?: string
  borderRadius?: string
}>(), {
  width: '100%',
  height: '20px',
  borderRadius: 'var(--radius-md)',
})
</script>

<template>
  <div
    class="skeleton-block"
    :style="{ width, height, borderRadius }"
  />
</template>
