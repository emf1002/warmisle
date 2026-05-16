<template>
  <div class="ledger-page">
    <!-- Page Header -->
    <div class="page-header">
      <h2>记账本</h2>
    </div>

    <!-- Month Picker -->
    <div class="month-row">
      <a-date-picker
        v-model:value="selectedMonth"
        picker="month"
        format="YYYY年M月"
        @change="onMonthChange"
        style="width: 160px"
      />
      <a-button type="primary" @click="openCreate()">记一笔</a-button>
    </div>

    <!-- Summary Bar -->
    <div class="summary-bar">
      <div class="summary-item">
        <span class="summary-label">收入</span>
        <span class="income-amount">+{{ formatYuan(summary.income) }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">支出</span>
        <span class="expense-amount">-{{ formatYuan(summary.expense) }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">结余</span>
        <span :class="summary.balance >= 0 ? 'income-amount' : 'expense-amount'">
          {{ summary.balance >= 0 ? '+' : '-' }}{{ formatYuan(Math.abs(summary.balance)) }}
        </span>
      </div>
    </div>

    <!-- Filters -->
    <div class="filter-row">
      <a-select
        v-model:value="filters.member_id"
        placeholder="全部成员"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
      >
        <a-select-option v-for="m in members" :key="m.id" :value="m.id">
          {{ m.avatar }} {{ m.name }}
        </a-select-option>
      </a-select>
      <a-select
        v-model:value="filters.category_id"
        placeholder="全部分类"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
      >
        <a-select-opt-group label="支出">
          <a-select-option v-for="c in expenseCategories" :key="c.id" :value="c.id">
            {{ c.icon }} {{ c.name }}
          </a-select-option>
        </a-select-opt-group>
        <a-select-opt-group label="收入">
          <a-select-option v-for="c in incomeCategories" :key="c.id" :value="c.id">
            {{ c.icon }} {{ c.name }}
          </a-select-option>
        </a-select-opt-group>
      </a-select>
      <a-select
        v-model:value="filters.creator_id"
        placeholder="全部创建者"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
      >
        <a-select-option v-for="m in members" :key="m.id" :value="m.id">
          {{ m.avatar }} {{ m.name }}
        </a-select-option>
      </a-select>
      <a-button size="small" @click="clearFilters">清除筛选</a-button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <!-- Empty State -->
    <div v-else-if="groups.length === 0" class="empty-state">
      <p>本月还没有记录，记一笔吧</p>
      <a-button type="primary" @click="openCreate()">记一笔</a-button>
    </div>

    <!-- Record List -->
    <div v-else class="record-list">
      <div v-for="group in groups" :key="group.date" class="date-group">
        <!-- Date Header -->
        <div class="date-header">
          <span class="date-text">{{ formatDate(group.date) }}</span>
          <span :class="group.daily_total >= 0 ? 'income-amount' : 'expense-amount'" class="date-total">
            {{ group.daily_total >= 0 ? '+' : '-' }}{{ formatYuan(Math.abs(group.daily_total)) }}
          </span>
        </div>

        <!-- Items -->
        <div
          v-for="item in group.items"
          :key="item.id"
          class="ledger-item"
          :class="{ clickable: canEdit(item) }"
          @click="onItemClick(item)"
        >
          <div class="item-main">
            <span class="item-icon">{{ item.category.icon }}</span>
            <div class="item-info">
              <span class="item-name">{{ item.category.name }}</span>
              <span v-if="item.note" class="item-note">{{ truncate(item.note, 30) }}</span>
            </div>
            <span :class="item.category.type === 'income' ? 'income-amount' : 'expense-amount'" class="item-amount">
              {{ formatAmount(item.amount, item.category.type) }}
            </span>
          </div>
          <div class="item-meta">
            <span class="item-creator">{{ item.creator.avatar }} {{ item.creator.name }}</span>
            <span v-if="item.members && item.members.length > 0" class="item-members">
              <span v-for="m in item.members" :key="m.id" class="member-avatar">{{ m.avatar }}</span>
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Load More -->
    <div v-if="hasMore" class="load-more">
      <a-button @click="loadMore" :loading="loadingMore">加载更多</a-button>
    </div>

    <!-- Create/Edit Dialog -->
    <a-modal
      v-model:open="dialogOpen"
      :title="editingRecord ? '编辑记录' : '记一笔'"
      :confirm-loading="submitting"
      width="480px"
    >
      <template #footer>
        <div style="display: flex; justify-content: space-between">
          <a-button
            v-if="editingRecord"
            danger
            @click="confirmDelete"
          >
            删除
          </a-button>
          <div style="display: flex; gap: 8px">
            <a-button @click="dialogOpen = false">取消</a-button>
            <a-button type="primary" :loading="submitting" @click="handleSubmit">
              {{ editingRecord ? '保存' : '记账' }}
            </a-button>
          </div>
        </div>
      </template>
      <a-form :model="form" layout="vertical">
        <a-form-item label="分类" required>
          <a-select v-model:value="form.category_id" placeholder="选择分类">
            <a-select-opt-group label="支出分类">
              <a-select-option
                v-for="cat in expenseCategories"
                :key="cat.id"
                :value="cat.id"
              >
                {{ cat.icon }} {{ cat.name }}
              </a-select-option>
            </a-select-opt-group>
            <a-select-opt-group label="收入分类">
              <a-select-option
                v-for="cat in incomeCategories"
                :key="cat.id"
                :value="cat.id"
              >
                {{ cat.icon }} {{ cat.name }}
              </a-select-option>
            </a-select-opt-group>
          </a-select>
        </a-form-item>
        <a-form-item label="金额（元）" required>
          <a-input-number
            v-model:value="form.amount"
            :min="0.01"
            :step="0.01"
            :precision="2"
            placeholder="0.00"
            style="width: 100%"
            inputmode="decimal"
          />
        </a-form-item>
        <a-form-item label="关联成员" required>
          <a-select
            v-model:value="form.member_ids"
            mode="multiple"
            placeholder="选择成员"
          >
            <a-select-option v-for="m in members" :key="m.id" :value="m.id">
              {{ m.avatar }} {{ m.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="日期">
          <a-date-picker
            v-model:value="form.occurred_at"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea
            v-model:value="form.note"
            :maxlength="200"
            :rows="2"
            placeholder="可选，最多200字"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import {
  getLedgers,
  createLedger,
  updateLedger,
  deleteLedger,
} from '@/api/ledger'
import { getCategories } from '@/api/category'
import { getMembers } from '@/api/member'

// Types
interface Member {
  id: number
  name: string
  avatar: string
}

interface Category {
  id: number
  name: string
  icon: string
  type: string
}

interface LedgerItem {
  id: number
  amount: number
  note: string
  category_id: number
  creator_id: number
  occurred_at: string
  category: Category
  creator: Member
  members: Member[]
}

interface LedgerGroup {
  date: string
  daily_total: number
  items: LedgerItem[]
}

interface LedgerSummary {
  income: number
  expense: number
  balance: number
}

interface Filters {
  member_id: number | undefined
  category_id: number | undefined
  creator_id: number | undefined
}

// State
const loading = ref(false)
const loadingMore = ref(false)
const submitting = ref(false)
const dialogOpen = ref(false)
const editingRecord = ref<LedgerItem | null>(null)

const selectedMonth = ref<Dayjs>(dayjs())
const members = ref<Member[]>([])
const categories = ref<Category[]>([])
const groups = ref<LedgerGroup[]>([])
const summary = ref<LedgerSummary>({ income: 0, expense: 0, balance: 0 })
const total = ref(0)
const page = ref(1)
const pageSize = 20

const filters = reactive<Filters>({
  member_id: undefined,
  category_id: undefined,
  creator_id: undefined,
})

const form = reactive({
  category_id: undefined as number | undefined,
  amount: null as number | null,
  member_ids: [] as number[],
  occurred_at: null as Dayjs | null,
  note: '',
})

// Computed
const hasMore = computed(() => groups.value.length > 0 && page.value * pageSize < total.value)

const expenseCategories = computed(() =>
  categories.value.filter((c) => c.type === 'expense')
)

const incomeCategories = computed(() =>
  categories.value.filter((c) => c.type === 'income')
)

const currentUserRole = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return ''
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role || ''
  } catch {
    return ''
  }
})

const currentUserId = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return 0
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return (payload.member_id as number) || 0
  } catch {
    return 0
  }
})

// Helpers
function formatYuan(cents: number): string {
  return `\u00A5${(cents / 100).toFixed(2)}`
}

function formatAmount(cents: number, categoryType: string): string {
  const yuan = (cents / 100).toFixed(2)
  return categoryType === 'income' ? `+\u00A5${yuan}` : `-\u00A5${yuan}`
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  const target = new Date(d.getFullYear(), d.getMonth(), d.getDate())

  if (target.getTime() === today.getTime()) return '今天'
  if (target.getTime() === yesterday.getTime()) return '昨天'

  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  return `${d.getMonth() + 1}月${d.getDate()}日 周${weekDays[d.getDay()]}`
}

function truncate(text: string, maxLen: number): string {
  if (text.length <= maxLen) return text
  return text.slice(0, maxLen) + '...'
}

function canEdit(item: LedgerItem): boolean {
  return item.creator_id === currentUserId.value || currentUserRole.value === 'admin'
}

// Methods
function getMonthParam(): string {
  return selectedMonth.value.format('YYYY-MM')
}

async function fetchLedgers(isLoadMore = false) {
  if (!isLoadMore) {
    loading.value = true
    page.value = 1
  } else {
    loadingMore.value = true
  }

  try {
    const params: Record<string, unknown> = {
      month: getMonthParam(),
      page: page.value,
      page_size: pageSize,
    }
    if (filters.member_id) params.member_id = filters.member_id
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.creator_id) params.creator_id = filters.creator_id

    const res: any = await getLedgers(params as any)
    const data = res.data
    if (isLoadMore) {
      groups.value = [...groups.value, ...(data.groups || [])]
    } else {
      groups.value = data.groups || []
    }
    summary.value = data.summary || { income: 0, expense: 0, balance: 0 }
    total.value = data.total || 0
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

async function loadMore() {
  page.value++
  await fetchLedgers(true)
}

function onMonthChange() {
  fetchLedgers()
}

function onFilterChange() {
  fetchLedgers()
}

function clearFilters() {
  filters.member_id = undefined
  filters.category_id = undefined
  filters.creator_id = undefined
  fetchLedgers()
}

function openCreate() {
  editingRecord.value = null
  form.category_id = undefined
  form.amount = null
  form.member_ids = []
  form.occurred_at = dayjs()
  form.note = ''
  dialogOpen.value = true
}

function openEdit(record: LedgerItem) {
  editingRecord.value = record
  form.category_id = record.category_id
  form.amount = record.amount / 100
  form.member_ids = record.members ? record.members.map((m) => m.id) : []
  form.occurred_at = dayjs(record.occurred_at)
  form.note = record.note || ''
  dialogOpen.value = true
}

function onItemClick(item: LedgerItem) {
  if (canEdit(item)) {
    openEdit(item)
  }
}

async function handleSubmit() {
  if (!form.category_id) {
    message.error('❌ 请选择分类')
    return
  }
  if (!form.amount || form.amount <= 0) {
    message.error('❌ 金额必须大于 0')
    return
  }
  if (form.member_ids.length === 0) {
    message.error('❌ 请至少选择一位关联成员')
    return
  }

  submitting.value = true
  try {
    const payload: any = {
      amount: form.amount!,
      note: form.note,
      category_id: form.category_id,
      member_ids: form.member_ids,
    }
    if (form.occurred_at) {
      payload.occurred_at = form.occurred_at.toISOString()
    }

    if (editingRecord.value) {
      if (payload.amount !== undefined) {
        payload.amount = form.amount!
      }
      await updateLedger(editingRecord.value.id, payload)
      message.success('✅ 更新成功')
    } else {
      await createLedger(payload)
      message.success('✅ 记账成功')
    }
    dialogOpen.value = false
    fetchLedgers()
  } catch (e: any) {
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  } finally {
    submitting.value = false
  }
}

function confirmDelete() {
  if (!editingRecord.value) return
  Modal.confirm({
    title: '❓ 确认删除',
    content: '确定要删除这条记账记录吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteLedger(editingRecord.value!.id)
        message.success('✅ 删除成功')
        dialogOpen.value = false
        fetchLedgers()
      } catch (e: any) {
        if (e?.response?.data?.message) {
          message.error(e.response.data.message)
        }
        throw e
      }
    },
  })
}

// Lifecycle
onMounted(async () => {
  try {
    const [catRes, memRes]: any[] = await Promise.all([
      getCategories(),
      getMembers(),
    ])
    categories.value = catRes.data || []
    members.value = memRes.data || []
  } catch {
    // error handled by interceptor
  }
  fetchLedgers()
})
</script>

<style scoped>
.ledger-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}

.month-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

/* Summary Bar */
.summary-bar {
  display: flex;
  background: #fff;
  border-radius: 8px;
  padding: 16px 24px;
  margin-bottom: 16px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

.summary-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.summary-label {
  font-size: 13px;
  color: #999;
}

/* Filter Row */
.filter-row {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}

/* Amount Colors */
.income-amount {
  color: #52c41a;
  font-weight: 600;
  font-size: 16px;
}

.expense-amount {
  color: #ff4d4f;
  font-weight: 600;
  font-size: 16px;
}

/* Loading & Empty */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 48px 0;
  color: #999;
}

/* Record List */
.record-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.date-group {
  margin-bottom: 0;
}

.date-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: #f5f5f5;
  border-radius: 6px;
  margin-bottom: 8px;
}

.date-text {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.date-total {
  font-size: 14px;
  font-weight: 500;
}

/* Ledger Item */
.ledger-item {
  background: #fff;
  border-radius: 8px;
  padding: 14px 16px;
  margin-bottom: 8px;
  border: 1px solid #f0f0f0;
  transition: box-shadow 0.2s;
  min-height: 44px;
}

.ledger-item.clickable {
  cursor: pointer;
}

.ledger-item.clickable:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.item-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.item-icon {
  font-size: 24px;
  flex-shrink: 0;
  width: 36px;
  text-align: center;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.item-name {
  font-size: 15px;
  font-weight: 500;
}

.item-note {
  font-size: 12px;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-amount {
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #fafafa;
}

.item-creator {
  font-size: 12px;
  color: #666;
}

.item-members {
  display: flex;
  gap: 2px;
}

.member-avatar {
  font-size: 14px;
}

/* Load More */
.load-more {
  text-align: center;
  padding: 16px 0;
}

/* Mobile */
@media (max-width: 767px) {
  .ledger-page {
    padding: 16px;
    padding-bottom: 80px;
  }

  .item-name { font-size: 14px; }

  .month-row :deep(.ant-btn-primary) {
    position: fixed;
    bottom: 72px;
    right: 20px;
    z-index: 100;
    width: 56px;
    height: 56px;
    border-radius: 50%;
    font-size: 24px;
    box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
  }
}
</style>
