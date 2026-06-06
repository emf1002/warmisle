<template>
  <div class="empty-state" data-testid="empty-state">
    <!-- 空状态插图 -->
    <div class="empty-illustration">
      <!-- 无数据插图 -->
      <svg
        v-if="type === 'no-data'"
        aria-hidden="true"
        width="200"
        height="160"
        viewBox="0 0 200 160"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        class="empty-svg"
      >
        <!-- 背景装饰 -->
        <circle cx="40" cy="120" r="8" fill="var(--color-brand-light)" opacity="0.5" />
        <circle cx="160" cy="40" r="6" fill="var(--color-chart-color-3, #F2B84B)" opacity="0.4" />
        <rect x="20" y="60" width="12" height="12" rx="3" fill="var(--color-border-secondary)" opacity="0.6" />

        <!-- 主要插画：文件夹/列表 -->
        <rect x="60" y="50" width="80" height="60" rx="8" fill="var(--color-bg-container)" stroke="var(--color-border)" stroke-width="2" />
        <line x1="75" y1="70" x2="125" y2="70" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="75" y1="82" x2="115" y2="82" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="75" y1="94" x2="120" y2="94" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />

        <!-- 装饰元素 -->
        <circle cx="135" cy="65" r="4" fill="var(--color-brand)" opacity="0.6" />
        <path d="M 50 45 L 55 35 L 60 45" stroke="var(--color-chart-color-2, #6BBAA7)" stroke-width="2" fill="none" opacity="0.5" />
      </svg>

      <!-- 无结果插图 -->
      <svg
        v-else-if="type === 'no-result'"
        aria-hidden="true"
        width="200"
        height="160"
        viewBox="0 0 200 160"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        class="empty-svg"
      >
        <!-- 放大镜 -->
        <circle cx="90" cy="70" r="35" fill="none" stroke="var(--color-border)" stroke-width="3" />
        <line x1="118" y1="98" x2="140" y2="120" stroke="var(--color-border)" stroke-width="3" stroke-linecap="round" />

        <!-- 问号 -->
        <text x="82" y="80" font-size="32" fill="var(--color-text-disabled)" font-weight="600" font-family="var(--font-display)">?</text>

        <!-- 装饰 dots -->
        <circle cx="50" cy="40" r="3" fill="var(--color-brand)" opacity="0.4" />
        <circle cx="150" cy="50" r="4" fill="var(--color-chart-color-3, #F2B84B)" opacity="0.4" />
        <circle cx="160" cy="110" r="3" fill="var(--color-chart-color-2, #6BBAA7)" opacity="0.4" />
      </svg>

      <!-- 网络错误插图 -->
      <svg
        v-else-if="type === 'error'"
        aria-hidden="true"
        width="200"
        height="160"
        viewBox="0 0 200 160"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        class="empty-svg"
      >
        <!-- 云服务器 -->
        <rect x="70" y="45" width="60" height="50" rx="8" fill="var(--color-bg-container)" stroke="var(--color-border)" stroke-width="2" />
        <line x1="85" y1="60" x2="115" y2="60" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="85" y1="70" x2="105" y2="70" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <circle cx="110" cy="70" r="3" fill="var(--color-danger)" opacity="0.6" />

        <!-- 断开的连接 -->
        <path d="M 50 90 Q 100 60 150 90" stroke="var(--color-border)" stroke-width="2" fill="none" stroke-dasharray="6 4" />
        <circle cx="50" cy="90" r="5" fill="var(--color-chart-color-2, #6BBAA7)" opacity="0.6" />
        <circle cx="150" cy="90" r="5" fill="var(--color-danger)" opacity="0.6" />
      </svg>

      <!-- 默认插图 -->
      <svg
        v-else
        aria-hidden="true"
        width="200"
        height="160"
        viewBox="0 0 200 160"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        class="empty-svg"
      >
        <rect x="50" y="40" width="100" height="80" rx="12" fill="var(--color-bg-container)" stroke="var(--color-border)" stroke-width="2" />
        <line x1="70" y1="65" x2="130" y2="65" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="70" y1="80" x2="115" y2="80" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
        <line x1="70" y1="95" x2="125" y2="95" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
      </svg>
    </div>

    <!-- 描述文本 -->
    <p class="empty-description">{{ displayDescription }}</p>

    <!-- 操作区域 -->
    <div v-if="type === 'no-data'" class="empty-action">
      <slot name="action" />
    </div>

    <!-- 清除筛选链接 -->
    <a
      v-if="type === 'no-result'"
      class="clear-link"
      @click="$emit('clear')"
      data-testid="clear-link"
    >
      <Icon name="RotateCcw" :size="14" />
      <span>清除筛选条件</span>
    </a>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Icon from './Icon.vue'

const props = withDefaults(defineProps<{
  type?: 'no-data' | 'no-result' | 'error'
  description?: string
}>(), {
  type: 'no-data'
})

defineEmits<{
  clear: []
}>()

const displayDescription = computed(() => {
  if (props.description) return props.description
  const map = {
    'no-data': '暂无数据，快来创建第一个吧',
    'no-result': '未找到匹配的结果',
    'error': '加载失败，请稍后重试'
  }
  return map[props.type] || '暂无数据'
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-2xl, 48px) var(--space-lg, 24px);
  text-align: center;
  animation: fadeIn var(--duration-slow, 300ms) ease-out;
}

/* ==================== 插图 ==================== */
.empty-illustration {
  margin-bottom: var(--space-lg, 24px);
  animation: slideUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) 100ms both;
}

.empty-svg {
  width: 200px;
  height: 160px;
  transition: transform var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1);
}

.empty-state:hover .empty-svg {
  transform: translateY(-4px);
}

/* ==================== 描述文本 ==================== */
.empty-description {
  font-size: var(--text-base, 1rem);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-md, 16px) 0;
  line-height: var(--leading-relaxed, 1.625);
  animation: fadeIn var(--duration-slow, 300ms) ease-out 200ms both;
}

/* ==================== 操作区域 ==================== */
.empty-action {
  min-height: 44px;
  animation: fadeIn var(--duration-slow, 300ms) ease-out 300ms both;
}

/* ==================== 清除链接 ==================== */
.clear-link {
  color: var(--color-brand);
  cursor: pointer;
  font-size: var(--text-sm, 0.875rem);
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs, 8px);
  text-decoration: none;
  padding: var(--space-xs, 8px) var(--space-sm, 12px);
  border-radius: var(--radius-md, 12px);
  transition: all var(--duration-fast, 150ms) ease-out;
  animation: fadeIn var(--duration-slow, 300ms) ease-out 200ms both;
}

.clear-link:hover {
  color: var(--color-brand-hover);
  background: var(--color-brand-light);
}

/* ==================== 动画 ==================== */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ==================== 响应式 ==================== */
@media (max-width: 480px) {
  .empty-illustration {
    margin-bottom: var(--space-md, 16px);
  }

  .empty-svg {
    width: 160px;
    height: 128px;
  }

  .empty-description {
    font-size: var(--text-sm, 0.875rem);
  }
}
</style>
