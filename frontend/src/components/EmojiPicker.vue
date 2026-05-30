<template>
  <div class="emoji-picker">
    <div
      v-for="emoji in emojiList"
      :key="emoji"
      class="emoji-item"
      :class="{ active: modelValue === emoji }"
      @click="$emit('update:modelValue', emoji)"
    >
      {{ emoji }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const emojiList = [
  '😀', '😂', '🥰', '😎', '🤔', '😴', '🥳', '😇',
  '🤗', '😋', '🤩', '😏', '😤', '🤯', '🫡', '🤭',
  '😌', '🤓', '🐶', '🐱', '🐼', '🦊', '🐸', '🐵',
  '🏠', '🌟', '🔥', '❤️', '💪', '🎉', '🌈', '☕'
]
</script>
