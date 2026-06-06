<template>
  <div
    class="skeleton-card"
    :style="styleObject"
    :class="{ 'skeleton-card--circle': circle }"
  ></div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  width?: string | number;
  height?: string | number;
  borderRadius?: string | number;
  circle?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '16px',
  borderRadius: 8,
  circle: false,
});

const styleObject = computed(() => {
  const style: Record<string, string> = {
    width: typeof props.width === 'number' ? `${props.width}px` : String(props.width),
    height: typeof props.height === 'number' ? `${props.height}px` : String(props.height),
  };

  if (props.circle) {
    style.borderRadius = '50%';
  } else {
    const radius = typeof props.borderRadius === 'number'
      ? `${props.borderRadius}px`
      : String(props.borderRadius);
    style.borderRadius = radius;
  }

  return style;
});
</script>

<style scoped>
.skeleton-card {
  background: linear-gradient(
    90deg,
    var(--color-bg-subtle, #f0f0f0) 0%,
    var(--color-bg-muted, #e0e0e0) 50%,
    var(--color-bg-subtle, #f0f0f0) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
  flex-shrink: 0;
}

.skeleton-card--circle {
  border-radius: 50% !important;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>
