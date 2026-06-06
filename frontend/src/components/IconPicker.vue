<template>
  <div class="icon-picker" data-testid="avatar-grid">
    <input
      v-model="search"
      class="icon-picker-search"
      type="text"
      placeholder="搜索图标..."
      data-testid="icon-search"
    />
    <div class="icon-picker-grid">
      <div
        v-for="icon in filteredIcons"
        :key="icon.name"
        class="icon-picker-item"
        :class="{ active: modelValue === icon.name }"
        :style="{ backgroundColor: icon.color }"
        :data-testid="`icon-item-${icon.name}`"
        @click="$emit('update:modelValue', icon.name)"
      >
        <Icon :name="icon.name" :size="20" color="#fff" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Icon from './Icon.vue'

defineProps<{
  modelValue: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const search = ref('')

interface IconEntry {
  name: string
  color: string
}

const icons: IconEntry[] = [
  // People / Identity
  { name: 'User', color: '#5B8FF9' },
  { name: 'Users', color: '#5B8FF9' },
  { name: 'UserCircle', color: '#4A90D9' },
  { name: 'Heart', color: '#E85D5D' },
  { name: 'Baby', color: '#FF8C42' },
  { name: 'Crown', color: '#F2B84B' },
  { name: 'Star', color: '#F2B84B' },
  { name: 'Smile', color: '#5AD8A6' },

  // Finance
  { name: 'Wallet', color: '#6BBAA7' },
  { name: 'Banknote', color: '#2ECC71' },
  { name: 'DollarSign', color: '#27AE60' },
  { name: 'TrendingUp', color: '#27AE60' },
  { name: 'TrendingDown', color: '#E85D5D' },
  { name: 'PieChart', color: '#9B59B6' },
  { name: 'BarChart3', color: '#4A90D9' },
  { name: 'PiggyBank', color: '#E88CCF' },
  { name: 'CreditCard', color: '#4A90D9' },
  { name: 'Receipt', color: '#95A5A6' },

  // Objects
  { name: 'Home', color: '#6BBAA7' },
  { name: 'Car', color: '#4ECDC4' },
  { name: 'ShoppingBag', color: '#E88CCF' },
  { name: 'Smartphone', color: '#4A90D9' },
  { name: 'BookOpen', color: '#F2B84B' },
  { name: 'Gamepad2', color: '#9B59B6' },
  { name: 'Gift', color: '#E74C3C' },
  { name: 'Sparkles', color: '#E88CCF' },
  { name: 'Package', color: '#95A5A6' },
  { name: 'Briefcase', color: '#F39C12' },
  { name: 'UtensilsCrossed', color: '#E87461' },
  { name: 'Dumbbell', color: '#2ECC71' },
  { name: 'PawPrint', color: '#8B6F47' },
  { name: 'ShieldCheck', color: '#3498DB' },
  { name: 'Mail', color: '#E74C3C' },
  { name: 'GraduationCap', color: '#F2B84B' },
  { name: 'HeartPulse', color: '#E85D5D' },
  { name: 'Stethoscope', color: '#E85D5D' },
  { name: 'Shirt', color: '#E88CCF' },
  { name: 'Plane', color: '#4ECDC4' },

  // Nature
  { name: 'Sun', color: '#F2B84B' },
  { name: 'Moon', color: '#9B59B6' },
  { name: 'Cloud', color: '#6DC8EC' },
  { name: 'Flame', color: '#E87461' },
  { name: 'Leaf', color: '#5AD8A6' },
  { name: 'TreePine', color: '#27AE60' },
  { name: 'Fish', color: '#4ECDC4' },
  { name: 'Bird', color: '#6DC8EC' },
  { name: 'Bug', color: '#5AD8A6' },
  { name: 'Flower2', color: '#E88CCF' },

  // Symbols
  { name: 'Zap', color: '#F2B84B' },
  { name: 'Target', color: '#E87461' },
  { name: 'Anchor', color: '#4A90D9' },
  { name: 'Compass', color: '#4ECDC4' },
  { name: 'Shield', color: '#3498DB' },
  { name: 'Key', color: '#F39C12' },
  { name: 'Bell', color: '#E87461' },
  { name: 'Flag', color: '#E74C3C' },
  { name: 'Globe', color: '#4A90D9' },
  { name: 'Coffee', color: '#8B6F47' },
  { name: 'Cake', color: '#E88CCF' },
  { name: 'Music', color: '#9B59B6' },
  { name: 'Camera', color: '#95A5A6' },
  { name: 'Phone', color: '#5AD8A6' },
  { name: 'MapPin', color: '#E85D5D' },
]

const filteredIcons = computed(() => {
  if (!search.value) return icons
  const q = search.value.toLowerCase()
  return icons.filter(i => i.name.toLowerCase().includes(q))
})
</script>

<style scoped>
.icon-picker {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.icon-picker-search {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--color-border, #e8e8e8);
  border-radius: var(--radius-sm, 8px);
  font-size: 14px;
  outline: none;
  background: var(--color-bg-container, #fff);
  color: var(--color-text-primary, #1f1f1f);
  transition: border-color 0.2s;
}

.icon-picker-search:focus {
  border-color: var(--color-brand, #1677ff);
}

.icon-picker-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.icon-picker-item {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm, 8px);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.15s, transform 0.15s;
}

.icon-picker-item:hover {
  transform: scale(1.08);
}

.icon-picker-item.active {
  border-color: var(--color-brand, #1677ff);
  transform: scale(1.08);
}
</style>
