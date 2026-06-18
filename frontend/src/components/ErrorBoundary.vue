<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'

const hasError = ref(false)
const errorMessage = ref('')

onErrorCaptured((err, _instance, info) => {
  hasError.value = true
  errorMessage.value = err.message || '未知错误'
  console.error('[ErrorBoundary]', err, info)
  // 返回 false 阻止错误继续向上传播
  return false
})

function reload() {
  hasError.value = false
  errorMessage.value = ''
  window.location.reload()
}
</script>

<template>
  <div v-if="hasError" class="error-boundary">
    <a-result status="error" title="页面出现了错误">
      <template #subTitle>
        <span>{{ errorMessage }}</span>
      </template>
      <template #extra>
        <a-button key="reload" type="primary" @click="reload">
          重新加载
        </a-button>
      </template>
    </a-result>
  </div>
  <slot v-else />
</template>

<style scoped>
.error-boundary {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 24px;
}
</style>
