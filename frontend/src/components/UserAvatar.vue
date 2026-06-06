<template>
  <div
    class="user-avatar"
    :style="{
      width: size + 'px',
      height: size + 'px',
      backgroundColor: isEmoji(icon) ? getAvatarColor(memberId) : getAvatarColor(memberId),
      fontSize: (size * 0.5) + 'px',
    }"
  >
    <Icon
      v-if="!isEmoji(icon) && icon"
      :name="icon"
      :size="size * 0.55"
      color="#fff"
    />
    <span v-else class="user-avatar-emoji" :style="{ fontSize: (size * 0.55) + 'px' }">
      {{ icon || '...' }}
    </span>
  </div>
</template>

<script setup lang="ts">
import Icon from './Icon.vue'
import { getAvatarColor, isEmoji } from '@/utils/icon-utils'

withDefaults(defineProps<{
  icon: string
  memberId?: number | null
  size?: number
}>(), {
  size: 32,
})
</script>

<style scoped>
.user-avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  flex-shrink: 0;
  overflow: hidden;
  line-height: 1;
}
.user-avatar-emoji {
  line-height: 1;
}
</style>
