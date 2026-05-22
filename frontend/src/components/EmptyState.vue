<template>
  <div class="empty-state">
    <div class="empty-icon">
      <svg aria-hidden="true" width="64" height="64" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="8" y="12" width="48" height="40" rx="4" stroke="var(--color-border)" stroke-width="2" fill="var(--color-border-secondary)" />
        <line x1="20" y1="24" x2="44" y2="24" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="20" y1="32" x2="38" y2="32" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="20" y1="40" x2="34" y2="40" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
      </svg>
    </div>
    <p class="empty-description">{{ displayDescription }}</p>
    <div v-if="type === 'no-data'" class="empty-action">
      <slot name="action" />
    </div>
    <a v-if="type === 'no-result'" class="clear-link" @click="$emit('clear')">清除筛选条件</a>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  type?: 'no-data' | 'no-result'
  description?: string
}>(), {
  type: 'no-data'
})

defineEmits<{
  clear: []
}>()

const displayDescription = computed(() => {
  if (props.description) return props.description
  return props.type === 'no-data' ? '暂无数据，快来创建第一个吧' : '未找到匹配的结果'
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 16px;
  text-align: center;
}

.empty-icon {
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 16px 0;
}

.empty-action {
  min-height: 44px;
}

.clear-link {
  color: var(--color-brand);
  cursor: pointer;
  font-size: 14px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  text-decoration: none;
}

.clear-link:hover {
  color: var(--color-brand-hover);
}
</style>
