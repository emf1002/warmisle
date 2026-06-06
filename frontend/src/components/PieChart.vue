<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'

const props = withDefaults(defineProps<{
  data: { label: string; value: number }[]
  size?: number
}>(), {
  size: 200,
})

const visible = ref(false)
onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

const total = computed(() => props.data.reduce((sum, d) => sum + d.value, 0))

const palette = [
  '#F44336', '#FF9800', '#FFC107', '#4CAF50', '#2196F3',
  '#9C27B0', '#E91E63', '#00BCD4', '#8BC34A', '#FF5722',
]

const slices = computed(() => {
  if (total.value === 0) return []
  let cumulative = 0
  return props.data.map((item, i) => {
    const start = cumulative
    const pct = item.value / total.value
    cumulative += pct
    return {
      label: item.label,
      value: item.value,
      pct,
      color: palette[i % palette.length],
      // SVG arc endpoints
      startAngle: start * 2 * Math.PI - Math.PI / 2,
      endAngle: cumulative * 2 * Math.PI - Math.PI / 2,
    }
  })
})

function arcPath(cx: number, cy: number, r: number, startAngle: number, endAngle: number): string {
  // If nearly full circle, draw two arcs (SVG can't do >180° with one arc)
  const diff = endAngle - startAngle
  if (diff >= 2 * Math.PI - 0.001) {
    // Full circle: draw two half arcs
    const mid = startAngle + Math.PI
    const x1 = cx + r * Math.cos(startAngle)
    const y1 = cy + r * Math.sin(startAngle)
    const x2 = cx + r * Math.cos(mid)
    const y2 = cy + r * Math.sin(mid)
    return `M ${x1} ${y1} A ${r} ${r} 0 1 1 ${x2} ${y2} A ${r} ${r} 0 1 1 ${x1} ${y1} Z`
  }
  const largeArc = diff > Math.PI ? 1 : 0
  const x1 = cx + r * Math.cos(startAngle)
  const y1 = cy + r * Math.sin(startAngle)
  const x2 = cx + r * Math.cos(endAngle)
  const y2 = cy + r * Math.sin(endAngle)
  return `M ${cx} ${cy} L ${x1} ${y1} A ${r} ${r} 0 ${largeArc} 1 ${x2} ${y2} Z`
}
</script>

<template>
  <div class="pie-chart-wrapper">
    <div class="pie-chart" :style="{ width: size + 'px', height: size + 'px' }">
      <svg :width="size" :height="size" :viewBox="`0 0 ${size} ${size}`">
        <path
          v-for="(s, i) in slices"
          :key="i"
          :d="arcPath(size / 2, size / 2, size / 2 - 2, s.startAngle, s.endAngle)"
          :fill="s.color"
          class="pie-slice"
          :style="{
            opacity: visible ? 1 : 0,
            transition: `opacity 0.4s ease ${i * 100}ms`,
          }"
        />
      </svg>
    </div>
    <div class="pie-legend">
      <div v-for="(s, i) in slices" :key="i" class="legend-item">
        <span class="legend-dot" :style="{ background: s.color }" />
        <span class="legend-label">{{ s.label }}</span>
        <span class="legend-amount">¥{{ (s.value).toFixed(2) }}</span>
        <span class="legend-pct">{{ (s.pct * 100).toFixed(0) }}%</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pie-chart-wrapper {
  display: flex;
  align-items: center;
  gap: 24px;
}

.pie-chart {
  flex-shrink: 0;
}

.pie-slice {
  transition: opacity 0.4s ease;
}

.pie-legend {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-label {
  color: var(--color-text-primary);
  flex: 1;
}

.legend-amount {
  color: var(--color-text-secondary);
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

.legend-pct {
  color: var(--color-text-secondary);
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  min-width: 32px;
  text-align: right;
}

@media (max-width: 767px) {
  .pie-chart-wrapper {
    flex-direction: column;
    align-items: center;
  }
}
</style>
