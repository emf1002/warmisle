<template>
  <div
    class="category-icon"
    :style="{
      width: size + 'px',
      height: size + 'px',
      backgroundColor: getCategoryColor(icon),
      fontSize: (size * 0.5) + 'px',
    }"
  >
    <Icon
      v-if="!isEmoji(icon) && icon"
      :name="icon"
      :size="size * 0.55"
      color="#fff"
    />
    <span v-else class="category-icon-emoji" :style="{ fontSize: (size * 0.55) + 'px' }">
      {{ icon || '...' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import Icon from './Icon.vue'
import { getCategoryColor, isEmoji } from '@/utils/icon-utils'

withDefaults(defineProps<{
  icon: string
  categoryId?: number | null
  size?: number
}>(), {
  size: 32,
})
</script>

<style scoped>
.category-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm, 8px);
  color: #fff;
  flex-shrink: 0;
  overflow: hidden;
  line-height: 1;
}
.category-icon-emoji {
  line-height: 1;
}
</style>
