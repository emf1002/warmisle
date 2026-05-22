<template>
  <div class="wish-page" data-testid="wish-page">
    <div class="page-header">
      <h2>愿望清单</h2>
      <a-button type="primary" @click="openCreate" data-testid="add-btn">新建愿望</a-button>
    </div>

    <div class="filter-row">
      <a-segmented v-model:value="activeType" :options="typeOptions" @change="onTypeChange" data-testid="type-switch" />
      <a-select
        v-model:value="filters.status"
        placeholder="全部状态"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="status-filter"
      >
        <a-select-option value="pending">待定</a-select-option>
        <a-select-option value="agreed">已同意</a-select-option>
        <a-select-option value="achieved">已实现</a-select-option>
        <a-select-option value="abandoned">已放弃</a-select-option>
      </a-select>
    </div>

    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <div v-else-if="wishes.length === 0" class="empty-state">
      <EmptyState v-if="!hasFilters" type="no-data" description="暂无愿望">
        <template #action>
          <a-button type="primary" @click="openCreate">新建愿望</a-button>
        </template>
      </EmptyState>
      <EmptyState v-else type="no-result" description="没有找到匹配的愿望" @clear="clearFilters" />
    </div>

    <div v-else class="wish-grid" data-testid="wish-grid">
      <div v-for="wish in wishes" :key="wish.id" class="wish-card" :data-testid="'wish-card-' + wish.id">
        <div class="wish-card-header">
          <span class="wish-category-tag">
            {{ categoryLabel(wish.category) }}
          </span>
          <a-tag :color="priorityColor(wish.priority)" size="small">
            {{ priorityLabel(wish.priority) }}
          </a-tag>
          <a-tag :color="statusColor(wish.status)" size="small">
            {{ statusLabel(wish.status) }}
          </a-tag>
        </div>
        <h3 class="wish-title">{{ wish.title }}</h3>
        <p v-if="wish.description" class="wish-desc">{{ wish.description }}</p>
        <div v-if="wish.amount" class="wish-amount">
          ¥{{ (wish.amount / 100).toFixed(2) }}
        </div>
        <div class="wish-card-footer">
          <span class="wish-creator">{{ wish.creator.avatar }} {{ wish.creator.name }}</span>
          <div class="wish-card-actions">
            <a-button
              v-if="activeType === 'family'"
              type="text"
              size="small"
              @click="handleVote(wish)"
              data-testid="vote-btn"
            >
              👍 {{ wish.vote_count }}
            </a-button>
            <a-dropdown v-if="canEdit(wish)" :trigger="['click']">
              <a-button type="text" size="small">···</a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="openEdit(wish)">编辑</a-menu-item>
                  <a-menu-item
                    v-if="wish.type === 'personal' && wish.creator_id === currentUserId"
                    @click="handlePromote(wish)"
                  >
                    提升为家庭愿望
                  </a-menu-item>
                  <a-menu-item
                    v-if="isAdmin && wish.status !== 'agreed'"
                    @click="handleStatusChange(wish, 'agreed')"
                  >
                    标记为同意
                  </a-menu-item>
                  <a-menu-item
                    v-if="isAdmin && wish.status !== 'achieved'"
                    @click="handleStatusChange(wish, 'achieved')"
                  >
                    标记为已实现
                  </a-menu-item>
                  <a-menu-item
                    v-if="canAbandon(wish) && wish.status !== 'abandoned'"
                    @click="handleStatusChange(wish, 'abandoned')"
                  >
                    标记为放弃
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item danger @click="confirmDelete(wish)">删除</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>
    </div>

    <div v-if="total > pageSize" class="pagination-row">
      <a-pagination
        v-model:current="page"
        :total="total"
        :page-size="pageSize"
        size="small"
        @change="fetchWishes"
      />
    </div>

    <a-modal
      v-model:open="dialogOpen"
      :title="editingWish ? '编辑愿望' : '新建愿望'"
      :confirm-loading="submitting"
      width="480px"
      @ok="handleSubmit"
      @cancel="dialogOpen = false"
      data-testid="wish-modal"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model:value="form.title" :maxlength="100" placeholder="请输入愿望标题" data-testid="title-input" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" :maxlength="500" :rows="3" placeholder="可选，最多500字" data-testid="desc-input" />
        </a-form-item>
        <a-form-item label="分类">
          <a-select v-model:value="form.category" data-testid="category-select">
            <a-select-option value="item">物品</a-select-option>
            <a-select-option value="travel">旅行</a-select-option>
            <a-select-option value="experience">体验</a-select-option>
            <a-select-option value="other">其他</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model:value="form.priority" data-testid="priority-select">
            <a-select-option value="normal">普通</a-select-option>
            <a-select-option value="important">重要</a-select-option>
            <a-select-option value="urgent">紧急</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="预估金额（元）">
          <a-input-number v-model:value="form.amountYuan" :min="0" :precision="2" style="width: 100%" placeholder="可选" data-testid="amount-input" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  getWishList, createWish, updateWish, deleteWish, promoteWish,
  updateWishStatus, voteWish, unvoteWish,
} from '@/api/wish'
import EmptyState from '@/components/EmptyState.vue'

interface WishItem {
  id: number
  title: string
  description: string
  category: string
  amount: number | null
  priority: string
  type: string
  status: string
  creator_id: number
  vote_count: number
  creator: { id: number; name: string; avatar: string }
}

const loading = ref(false)
const submitting = ref(false)
const dialogOpen = ref(false)
const editingWish = ref<WishItem | null>(null)
const wishes = ref<WishItem[]>([])
const total = ref(0)
const pageSize = 12
const page = ref(1)
const activeType = ref('personal')

const filters = reactive({
  status: undefined as string | undefined,
})

const typeOptions = [
  { label: '个人愿望', value: 'personal' },
  { label: '家庭愿望', value: 'family' },
]

const hasFilters = computed(() => !!filters.status)

const currentUserId = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return 0
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return (payload.member_id as number) || 0
  } catch { return 0 }
})

const isAdmin = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return false
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role === 'admin'
  } catch { return false }
})

function canEdit(wish: WishItem): boolean {
  if (isAdmin.value) return true
  return wish.creator_id === currentUserId.value
}

function canAbandon(wish: WishItem): boolean {
  if (isAdmin.value) return true
  return wish.creator_id === currentUserId.value
}

function priorityColor(p: string): string {
  if (p === 'urgent') return 'red'
  if (p === 'important') return 'orange'
  return 'default'
}

function priorityLabel(p: string): string {
  if (p === 'urgent') return '紧急'
  if (p === 'important') return '重要'
  return '普通'
}

function statusColor(s: string): string {
  if (s === 'agreed') return 'blue'
  if (s === 'achieved') return 'green'
  if (s === 'abandoned') return 'default'
  return 'default'
}

function statusLabel(s: string): string {
  if (s === 'pending') return '待定'
  if (s === 'agreed') return '已同意'
  if (s === 'achieved') return '已实现'
  if (s === 'abandoned') return '已放弃'
  return s
}

function categoryLabel(c: string): string {
  if (c === 'item') return '🛍️ 物品'
  if (c === 'travel') return '✈️ 旅行'
  if (c === 'experience') return '🎯 体验'
  return '📦 其他'
}

async function fetchWishes() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: page.value,
      page_size: pageSize,
      type: activeType.value,
    }
    if (filters.status) params.status = filters.status
    const res: any = await getWishList(params as any)
    const data = res.data
    wishes.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function onTypeChange() {
  page.value = 1
  fetchWishes()
}

function onFilterChange() {
  page.value = 1
  fetchWishes()
}

function clearFilters() {
  filters.status = undefined
  page.value = 1
  fetchWishes()
}

function openCreate() {
  editingWish.value = null
  form.title = ''
  form.description = ''
  form.category = 'other'
  form.priority = 'normal'
  form.amountYuan = null
  dialogOpen.value = true
}

function openEdit(wish: WishItem) {
  editingWish.value = wish
  form.title = wish.title
  form.description = wish.description || ''
  form.category = wish.category
  form.priority = wish.priority
  form.amountYuan = wish.amount ? wish.amount / 100 : null
  dialogOpen.value = true
}

const form = reactive({
  title: '',
  description: '',
  category: 'other' as string,
  priority: 'normal' as string,
  amountYuan: null as number | null,
})

async function handleSubmit() {
  if (!form.title.trim()) {
    message.error('❌ 请输入标题')
    return
  }
  submitting.value = true
  try {
    const payload: any = {
      title: form.title.trim(),
      description: form.description.trim(),
      category: form.category,
      priority: form.priority,
    }
    if (form.amountYuan != null) {
      payload.amount = Math.round(form.amountYuan * 100)
    }
    if (editingWish.value) {
      await updateWish(editingWish.value.id, payload)
      message.success('✅ 更新成功')
    } else {
      await createWish(payload)
      message.success('✅ 创建成功')
    }
    dialogOpen.value = false
    fetchWishes()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  } finally {
    submitting.value = false
  }
}

async function handlePromote(wish: WishItem) {
  Modal.confirm({
    title: '❓ 提升为家庭愿望',
    content: '提升后将展示给所有成员且不可撤回，确定继续吗？',
    okText: '确定',
    cancelText: '取消',
    async onOk() {
      try {
        await promoteWish(wish.id)
        message.success('✅ 已提升为家庭愿望')
        fetchWishes()
      } catch (e: any) {
        if (e?.response?.data?.message) message.error(e.response.data.message)
        throw e
      }
    },
  })
}

async function handleStatusChange(wish: WishItem, status: string) {
  try {
    await updateWishStatus(wish.id, status)
    message.success('✅ 状态已更新')
    fetchWishes()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }
}

async function handleVote(wish: WishItem) {
  try {
    await voteWish(wish.id)
    message.success('✅ 投票成功')
    fetchWishes()
  } catch (e: any) {
    const msg = e?.response?.data?.message || ''
    if (msg.includes('已投票')) {
      Modal.confirm({
        title: '❓ 取消投票',
        content: '你已经投过票了，要取消投票吗？',
        okText: '取消投票',
        cancelText: '算了',
        async onOk() {
          try {
            await unvoteWish(wish.id)
            message.success('✅ 已取消投票')
            fetchWishes()
          } catch { /* ignore */ }
        },
      })
    } else {
      message.error(msg || '操作失败')
    }
  }
}

function confirmDelete(wish: WishItem) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: `确定要删除愿望「${wish.title}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteWish(wish.id)
        message.success('✅ 删除成功')
        dialogOpen.value = false
        fetchWishes()
      } catch (e: any) {
        if (e?.response?.data?.message) message.error(e.response.data.message)
        throw e
      }
    },
  })
}

onMounted(() => {
  fetchWishes()
})
</script>

<style scoped>
.wish-page {
  padding: 24px;
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 { margin: 0; }

.filter-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
  color: var(--color-text-secondary);
}

.empty-state {
  text-align: center;
  padding: 48px 0;
}

.wish-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.wish-card {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow 0.2s;
  display: flex;
  flex-direction: column;
}

.wish-card:hover { box-shadow: var(--shadow-level-2); }

.wish-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.wish-category-tag {
  font-size: 13px;
  margin-right: auto;
}

.wish-title {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 4px 0;
  word-break: break-word;
}

.wish-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0 0 8px 0;
  word-break: break-word;
}

.wish-amount {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-brand);
  margin-bottom: 8px;
}

.wish-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  min-height: 44px;
}

.wish-creator {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.wish-card-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.vote-ended {
	opacity: 0.7;
	pointer-events: none;
}

@media (max-width: 767px) {
  .wish-page { padding: 16px; }
  .wish-grid { grid-template-columns: 1fr; }
  .wish-title { font-size: 15px; }
}
</style>
