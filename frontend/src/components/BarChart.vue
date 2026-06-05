<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  data: { label: string; value: number }[]
  color?: string
  height?: number
}>(), {
  color: 'var(--color-brand)',
  height: 160,
})

const visible = ref(false)

onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

const max = computed(() => Math.max(...props.data.map(d => d.value), 1))
</script>

<template>
  <div class="bar-chart" :style="{ height: height + 'px' }">
    <div
      v-for="(item, i) in data"
      :key="i"
      class="bar-col"
    >
      <div class="bar-tooltip">{{ item.value }}</div>
      <div
        class="bar-fill"
        :style="{
          height: visible ? ((item.value / max) * 100) + '%' : '0%',
          transitionDelay: (i * 80) + 'ms',
        }"
      >
        <div class="bar-gradient" :style="{ background: `linear-gradient(180deg, ${color}, ${color}33)` }" />
      </div>
      <span class="bar-label">{{ item.label }}</span>
    </div>
  </div>
</template>

<style scoped>
.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 0 8px;
}
.bar-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}
.bar-fill {
  width: 100%;
  transition: height 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
  border-radius: 6px 6px 0 0;
  overflow: hidden;
}
.bar-gradient {
  width: 100%;
  height: 100%;
}
.bar-label {
  margin-top: 6px;
  font-size: 11px;
  color: var(--color-muted);
  font-family: var(--font-display);
}
.bar-tooltip {
  display: none;
  position: absolute;
  top: -28px;
  background: var(--color-text-primary);
  color: var(--color-bg-container);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.bar-col:hover .bar-tooltip {
  display: block;
}
</style>
