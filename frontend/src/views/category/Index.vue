<template>
  <div class="category-page" data-testid="category-page">
    <div class="page-header">
      <h2>分类管理</h2>
      <a-button v-if="authStore.isAdmin" type="primary" @click="openCreate" data-testid="add-btn">
        <Icon name="Plus" :size="16" /> 添加分类
      </a-button>
    </div>

    <!-- Expense categories -->
    <div class="category-section" data-testid="expense-categories">
      <h3 class="section-title">
        <Icon name="TrendingDown" :size="16" /> 支出分类
      </h3>
      <div v-if="expenseCategories.length === 0" class="empty-hint">暂无支出分类</div>
      <div class="category-grid">
        <div
          v-for="cat in expenseCategories"
          :key="cat.id"
          class="category-card"
          :data-testid="'category-card-' + cat.id"
        >
          <CategoryIcon :icon="cat.icon" :category-id="cat.id" :size="24" />
          <span class="category-name">{{ cat.name }}</span>
          <span v-if="authStore.isAdmin" class="category-actions">
            <a-button type="link" size="small" @click="openEdit(cat)" data-testid="edit-btn"><Icon name="Pencil" :size="14" /> 编辑</a-button>
            <a-button type="link" size="small" danger @click="confirmDelete(cat)" data-testid="delete-btn"><Icon name="Trash2" :size="14" /> 删除</a-button>
          </span>
        </div>
      </div>
    </div>

    <!-- Income categories -->
    <div class="category-section" data-testid="income-categories">
      <h3 class="section-title">
        <Icon name="TrendingUp" :size="16" /> 收入分类
      </h3>
      <div v-if="incomeCategories.length === 0" class="empty-hint">暂无收入分类</div>
      <div class="category-grid">
        <div
          v-for="cat in incomeCategories"
          :key="cat.id"
          class="category-card"
          :data-testid="'category-card-' + cat.id"
        >
          <CategoryIcon :icon="cat.icon" :category-id="cat.id" :size="24" />
          <span class="category-name">{{ cat.name }}</span>
          <span v-if="authStore.isAdmin" class="category-actions">
            <a-button type="link" size="small" @click="openEdit(cat)" data-testid="edit-btn"><Icon name="Pencil" :size="14" /> 编辑</a-button>
            <a-button type="link" size="small" danger @click="confirmDelete(cat)" data-testid="delete-btn"><Icon name="Trash2" :size="14" /> 删除</a-button>
          </span>
        </div>
      </div>
    </div>

    <!-- Create/Edit Dialog -->
    <a-modal
      v-model:open="dialogOpen"
      :title="editingCategory ? '编辑分类' : '添加分类'"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSubmit"
      :confirm-loading="submitting"
      :ok-button-props="({ 'data-testid': 'modal-submit-btn' } as any)"
    >
      <div data-testid="category-modal">
      <a-form :model="form" layout="vertical">
        <a-form-item label="类型" required>
          <a-select v-model:value="form.type" data-testid="type-select">
            <a-select-option value="expense">支出</a-select-option>
            <a-select-option value="income">收入</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="名称" required>
          <a-input
            v-model:value="form.name"
            placeholder="分类名称"
            :maxlength="20"
            data-testid="name-input"
          />
        </a-form-item>
        <a-form-item label="图标">
          <IconPicker v-model="form.icon" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number
            v-model:value="form.sort_order"
            :min="0"
            :max="999"
            placeholder="数字越小越靠前"
            style="width: 100%"
            data-testid="sort-input"
          />
        </a-form-item>
      </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,
} from '@/api/category'
import { useAuthStore } from '@/stores/auth'
import IconPicker from '@/components/IconPicker.vue'
import Icon from '@/components/Icon.vue'
import CategoryIcon from '@/components/CategoryIcon.vue'

const authStore = useAuthStore()

const categories = ref<any[]>([])
const dialogOpen = ref(false)
const editingCategory = ref<any>(null)
const submitting = ref(false)

const form = reactive({
  type: 'expense',
  name: '',
  icon: 'Package',
  sort_order: 0,
})

const expenseCategories = computed(() =>
  categories.value.filter((c: any) => c.type === 'expense')
)

const incomeCategories = computed(() =>
  categories.value.filter((c: any) => c.type === 'income')
)

onMounted(() => {
  fetchCategories()
})

async function fetchCategories() {
  try {
    const res: any = await getCategories()
    categories.value = res.data || []
  } catch {
    // error handled by interceptor
  }
}

function openCreate() {
  editingCategory.value = null
  form.type = 'expense'
  form.name = ''
  form.icon = 'Package'
  form.sort_order = 0
  dialogOpen.value = true
}

function openEdit(record: any) {
  editingCategory.value = record
  form.type = record.type
  form.name = record.name
  form.icon = record.icon
  form.sort_order = record.sort_order
  dialogOpen.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingCategory.value) {
      await updateCategory(editingCategory.value.id, {
        type: form.type,
        name: form.name,
        icon: form.icon,
        sort_order: form.sort_order,
      })
      message.success('✅ 更新成功')
    } else {
      if (!form.name) {
        message.error('❌ 分类名称不能为空')
        submitting.value = false
        return
      }
      await createCategory({
        type: form.type,
        name: form.name,
        icon: form.icon,
        sort_order: form.sort_order,
      })
      message.success('✅ 添加成功')
    }
    dialogOpen.value = false
    fetchCategories()
  } catch {
    // error handled by interceptor
  } finally {
    submitting.value = false
  }
}

function confirmDelete(record: any) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: `确定要删除分类「${record.name}」吗？如果该分类下有记账记录则无法删除。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    okButtonProps: { 'data-testid': 'modal-confirm-btn' } as any,
    async onOk() {
      try {
        await deleteCategory(record.id)
        message.success('✅ 删除成功')
        fetchCategories()
      } catch (e: any) {
        throw e // prevent modal from closing on error
      }
    },
  })
}
</script>

<style scoped>
.category-page {
  padding: var(--space-lg);
}

.category-section {
  margin-bottom: 32px;
}

.section-title {
  font-size: 16px;
  margin: 0 0 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-icon {
  font-size: 18px;
}

.empty-hint {
  color: var(--color-text-secondary);
  font-size: 14px;
  padding: 12px 0;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.category-card {
  display: flex;
  align-items: center;
  padding: var(--space-sm) var(--space-md);
  background: var(--color-bg-container);
  border: 1px solid var(--color-border-secondary);
  border-radius: var(--radius-md);
  gap: 12px;
  min-height: 44px;
  box-shadow: var(--shadow-level-1);
}

.category-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.category-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.category-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

</style>
