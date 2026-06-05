<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  data: number[]
  color?: string
  height?: number
}>(), {
  color: 'var(--color-brand)',
  height: 32,
})

const visible = ref(false)

onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

const max = computed(() => Math.max(...props.data, 1))
</script>

<template>
  <div class="sparkline" :style="{ height: height + 'px' }">
    <div
      v-for="(val, i) in data"
      :key="i"
      class="sparkline-bar"
      :style="{
        height: visible ? ((val / max) * 100) + '%' : '0%',
        background: color,
        transitionDelay: (i * 50) + 'ms',
      }"
    />
  </div>
</template>

<style scoped>
.sparkline {
  display: flex;
  align-items: flex-end;
  gap: 2px;
}
.sparkline-bar {
  flex: 1;
  border-radius: 2px;
  opacity: 0.6;
  transition: height 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
