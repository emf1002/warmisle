<template>
  <div class="ledger-page" data-testid="ledger-page">
    <!-- Page Header -->
    <div class="page-header">
      <h2>记账本</h2>
    </div>

    <!-- Month Navigation -->
    <div class="month-switcher">
      <a-button type="text" @click="goPrevMonth" class="month-arrow" data-testid="month-prev">◀</a-button>
      <span class="month-text" data-testid="current-month">{{ selectedMonth.format('YYYY年M月') }}</span>
      <a-button type="text" @click="goNextMonth" class="month-arrow" data-testid="month-next">▶</a-button>
    </div>

    <!-- Date Range Picker -->
    <div class="month-row">
      <a-range-picker
        v-model:value="dateRange"
        format="YYYY-MM-DD"
        :presets="rangePresets"
        @change="onDateRangeChange"
        data-testid="date-range-picker"
      />
      <a-button type="primary" @click="openCreate()" data-testid="add-btn"><span class="btn-text-full">记一笔</span><span class="btn-text-short">记</span></a-button>
    </div>

    <!-- Summary Bar -->
    <div class="summary-bar">
      <div class="summary-item">
        <span class="summary-label">收入</span>
        <span class="income-amount" data-testid="summary-income">+{{ formatYuan(summary.income) }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">支出</span>
        <span class="expense-amount" data-testid="summary-expense">-{{ formatYuan(summary.expense) }}</span>
      </div>
      <div class="summary-item">
        <span class="summary-label">结余</span>
        <span :class="summary.balance >= 0 ? 'income-amount' : 'expense-amount'" data-testid="summary-balance">
          {{ summary.balance >= 0 ? '+' : '-' }}{{ formatYuan(Math.abs(summary.balance)) }}
        </span>
      </div>
    </div>

    <!-- Filters -->
    <div class="filter-row">
      <a-select
        v-model:value="filters.category_id"
        placeholder="全部分类"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="filter-category"
      >
        <template v-if="expenseCategories.length > 0">
          <a-select-option disabled :value="-1" class="category-group-label">支出</a-select-option>
          <a-select-option v-for="c in expenseCategories" :key="c.id" :value="c.id">
            {{ c.icon }} {{ c.name }}
          </a-select-option>
        </template>
        <template v-if="incomeCategories.length > 0">
          <a-select-option disabled :value="-2" class="category-group-label">收入</a-select-option>
          <a-select-option v-for="c in incomeCategories" :key="c.id" :value="c.id">
            {{ c.icon }} {{ c.name }}
          </a-select-option>
        </template>
      </a-select>
      <a-select
        v-model:value="filters.creator_id"
        placeholder="全部创建者"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="filter-creator"
      >
        <a-select-option v-for="m in members" :key="m.id" :value="m.id">
          {{ m.avatar }} {{ m.name }}
        </a-select-option>
      </a-select>
      <a-button size="small" @click="clearFilters" data-testid="clear-filters">清除筛选</a-button>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading" class="skeleton-state">
      <a-skeleton active :paragraph="{ rows: 6 }" :title="false" />
    </div>

    <!-- Empty State -->
    <div v-else-if="groups.length === 0" class="empty-state" data-testid="empty-state">
      <p>当前时间段还没有记录，记一笔吧</p>
      <a-button type="primary" @click="openCreate()">记一笔</a-button>
    </div>

    <!-- Record List -->
    <div v-else class="record-list" data-testid="record-list">
      <div v-for="group in groups" :key="group.date" class="date-group" data-testid="date-group">
        <!-- Date Header -->
        <div class="date-header">
          <span class="date-text">{{ formatDate(group.date) }}</span>
          <span :class="group.daily_total >= 0 ? 'income-amount' : 'expense-amount'" class="date-total" data-testid="daily-total">
            {{ group.daily_total >= 0 ? '+' : '-' }}{{ formatYuan(Math.abs(group.daily_total)) }}
          </span>
        </div>

        <!-- Items -->
        <div
          v-for="(item, index) in group.items"
          :key="item.id"
          class="ledger-item card-stagger"
          :class="{ clickable: canEdit(item) }"
          :style="{ animationDelay: `${index * 50}ms` }"
          :data-testid="'ledger-item-' + item.id"
          @click="onItemClick(item)"
        >
          <div class="item-top">
            <span class="item-category">
              <span class="item-cat-icon">{{ item.category.icon }}</span>
              <span class="item-cat-name">{{ item.category.name }}</span>
            </span>
            <span class="item-creator-line">
              <span class="creator-avatar" :aria-label="`${item.creator.name}的头像`">{{ item.creator.avatar }}</span>
              <span class="creator-name" data-testid="creator-name">{{ item.creator.name }}</span>
              <span class="creator-label">记录</span>
            </span>
          </div>
          <div v-if="item.note" class="item-body">
            <span class="item-note">{{ truncate(item.note, 30) }}</span>
          </div>
          <div class="item-bottom">
            <span class="item-amount" :class="item.category.type === 'income' ? 'income-amount' : 'expense-amount'">
              {{ item.category.type === 'income' ? '+' : '-' }}¥{{ (item.amount / 100).toFixed(2) }}
            </span>
            <a-button
              v-if="canEdit(item)"
              type="link"
              size="small"
              @click.stop="onItemClick(item)"
              data-testid="edit-btn"
            >
              编辑
            </a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- Infinite scroll sentinel -->
    <div ref="sentinelRef" v-if="hasMore" class="load-sentinel" data-testid="load-sentinel">
      <a-spin v-if="loadingMore" size="small" />
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
            data-testid="delete-btn"
          >
            删除
          </a-button>
          <div style="display: flex; gap: 8px">
            <a-button @click="dialogOpen = false">取消</a-button>
            <a-button type="primary" :loading="submitting" @click="handleSubmit" data-testid="submit-btn">
              {{ editingRecord ? '保存' : '记账' }}
            </a-button>
          </div>
        </div>
      </template>
      <div data-testid="ledger-modal">
      <a-form :model="form" layout="vertical">
        <a-form-item label="分类" required>
          <a-tabs v-model:activeKey="categoryTab" size="small">
            <a-tab-pane key="expense" tab="支出">
              <div class="category-grid-picker">
                <div
                  v-for="cat in expenseCategories"
                  :key="cat.id"
                  :class="['category-pick-item', { active: form.category_id === cat.id }]"
                  @click="form.category_id = cat.id"
                  :data-testid="'cat-expense-' + cat.id"
                >
                  <span class="category-pick-icon">{{ cat.icon }}</span>
                  <span class="category-pick-name">{{ cat.name }}</span>
                </div>
              </div>
            </a-tab-pane>
            <a-tab-pane key="income" tab="收入">
              <div class="category-grid-picker">
                <div
                  v-for="cat in incomeCategories"
                  :key="cat.id"
                  :class="['category-pick-item', { active: form.category_id === cat.id }]"
                  @click="form.category_id = cat.id"
                  :data-testid="'cat-income-' + cat.id"
                >
                  <span class="category-pick-icon">{{ cat.icon }}</span>
                  <span class="category-pick-name">{{ cat.name }}</span>
                </div>
              </div>
            </a-tab-pane>
          </a-tabs>
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
            data-testid="amount-input"
          />
        </a-form-item>
        <a-form-item label="日期">
          <a-date-picker
            v-model:value="form.occurred_at"
            show-time
            format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
            data-testid="date-picker"
          />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea
            v-model:value="form.note"
            :maxlength="200"
            :rows="2"
            placeholder="可选，最多200字"
            data-testid="note-input"
          />
        </a-form-item>
      </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import { formatDate, truncate } from '@/utils/format'
import {
  getLedgers,
  createLedger,
  updateLedger,
  deleteLedger,
} from '@/api/ledger'
import { useCategoriesStore } from '@/stores/categories'
import { useMembersStore } from '@/stores/members'
import { useAuthStore } from '@/stores/auth'

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
  category_id: number | undefined
  creator_id: number | undefined
}

// State
const authStore = useAuthStore()
const categoriesStore = useCategoriesStore()
const membersStore = useMembersStore()
const loading = ref(false)
const loadingMore = ref(false)
const submitting = ref(false)
const dialogOpen = ref(false)
const editingRecord = ref<LedgerItem | null>(null)

const selectedMonth = ref(dayjs())
const dateRange = ref<[Dayjs, Dayjs]>([dayjs().startOf('month'), dayjs().endOf('month')])
const groups = ref<LedgerGroup[]>([])
const summary = ref<LedgerSummary>({ income: 0, expense: 0, balance: 0 })
const nextCursor = ref<string | null>(null)
const hasMore = ref(false)
const sentinelRef = ref<HTMLElement | null>(null)
let abortController: AbortController | null = null
let observer: IntersectionObserver | null = null

const rangePresets = [
  { label: '本月', value: [dayjs().startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '上月', value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] as [Dayjs, Dayjs] },
  { label: '近三个月', value: [dayjs().subtract(2, 'month').startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '近半年', value: [dayjs().subtract(5, 'month').startOf('month'), dayjs().endOf('month')] as [Dayjs, Dayjs] },
  { label: '今年', value: [dayjs().startOf('year'), dayjs().endOf('year')] as [Dayjs, Dayjs] },
]

const filters = reactive<Filters>({
  category_id: undefined,
  creator_id: undefined,
})

const form = reactive({
  category_id: undefined as number | undefined,
  amount: null as number | null,
  occurred_at: null as Dayjs | null,
  note: '',
})

const categoryTab = ref<'expense' | 'income'>('expense')

// Computed
const members = computed(() => membersStore.members)
const categories = computed(() => categoriesStore.categories)

const expenseCategories = computed(() =>
  categories.value.filter((c) => c.type === 'expense')
)

const incomeCategories = computed(() =>
  categories.value.filter((c) => c.type === 'income')
)

// Helpers
function debounce<F extends (...args: any[]) => any>(fn: F, delay: number): F {
  let timer: ReturnType<typeof setTimeout>
  return ((...args: any[]) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }) as any
}

function formatYuan(cents: number): string {
  return `\u00A5${(cents / 100).toFixed(2)}`
}

function canEdit(item: LedgerItem): boolean {
  return item.creator_id === authStore.currentUserId || authStore.isAdmin
}

// Methods
function onDateRangeChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    dateRange.value = dates
    fetchLedgers()
  }
}

function goPrevMonth() {
  selectedMonth.value = selectedMonth.value.subtract(1, 'month')
  dateRange.value = [selectedMonth.value.startOf('month'), selectedMonth.value.endOf('month')]
  fetchLedgers()
}

function goNextMonth() {
  selectedMonth.value = selectedMonth.value.add(1, 'month')
  dateRange.value = [selectedMonth.value.startOf('month'), selectedMonth.value.endOf('month')]
  fetchLedgers()
}

async function fetchLedgers(append = false) {
  // Cancel any in-flight request
  if (abortController) abortController.abort()
  abortController = new AbortController()

  if (!append) {
    loading.value = true
    nextCursor.value = null
    hasMore.value = false
  } else {
    loadingMore.value = true
  }

  try {
    const params: Record<string, unknown> = {
      start_date: dateRange.value[0].format('YYYY-MM-DD'),
      end_date: dateRange.value[1].add(1, 'day').format('YYYY-MM-DD'),
      limit: 20,
    }
    if (filters.category_id) params.category_id = filters.category_id
    if (filters.creator_id) params.creator_id = filters.creator_id
    if (append && nextCursor.value) params.cursor = nextCursor.value

    const res: any = await getLedgers(params as any, abortController.signal)
    const data = res.data
    if (append) {
      groups.value = [...groups.value, ...(data.groups || [])]
    } else {
      groups.value = data.groups || []
    }
    summary.value = data.summary || { income: 0, expense: 0, balance: 0 }
    nextCursor.value = data.next_cursor ?? null
    hasMore.value = data.has_more ?? false
  } catch (e: any) {
    if (e?.code !== 'ERR_CANCELED' && !append) {
      groups.value = []
      summary.value = { income: 0, expense: 0, balance: 0 }
      nextCursor.value = null
      hasMore.value = false
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const debouncedFetchLedgers = debounce(() => fetchLedgers(false), 300)

function onFilterChange() {
  debouncedFetchLedgers()
}

function clearFilters() {
  filters.category_id = undefined
  filters.creator_id = undefined
  nextCursor.value = null
  hasMore.value = false
  fetchLedgers(false)
}

function openCreate() {
  editingRecord.value = null
  form.category_id = undefined
  form.amount = null
  form.occurred_at = dayjs()
  form.note = ''
  categoryTab.value = 'expense'
  dialogOpen.value = true
}

function openEdit(record: LedgerItem) {
  editingRecord.value = record
  form.category_id = record.category_id
  form.amount = record.amount / 100
  form.occurred_at = dayjs(record.occurred_at)
  form.note = record.note || ''
  const cat = categories.value.find(c => c.id === record.category_id)
  categoryTab.value = cat?.type === 'income' ? 'income' : 'expense'
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

  submitting.value = true
  try {
    const payload: any = {
      amount: Math.round(form.amount! * 100),
      note: form.note,
      category_id: form.category_id,
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
  } catch {
    // error handled by interceptor
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
    await Promise.all([
      categoriesStore.fetchCategories(),
      membersStore.fetchMembers(),
    ])
  } catch {
    // error handled by interceptor
  }
  await fetchLedgers(false)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
  if (abortController) abortController.abort()
})

// Watch sentinelRef to (re-)attach IntersectionObserver whenever it appears/disappears
watch(sentinelRef, (el, oldEl) => {
  if (oldEl && observer) observer.unobserve(oldEl)
  if (!el) return
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting && hasMore.value && !loadingMore.value) {
        fetchLedgers(true)
      }
    },
    { rootMargin: '200px' }
  )
  observer.observe(el)
})
</script>

<style scoped>
.ledger-page {
  padding: var(--space-lg);
  max-width: 800px;
  margin: 0 auto;
}

/* Month Row */
/* Month Switcher */
.month-switcher {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  gap: 16px;
}

.month-arrow {
  min-width: 44px;
  min-height: 44px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.month-text {
  font-size: 16px;
  font-weight: 600;
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
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 16px 24px;
  margin-bottom: 16px;
  box-shadow: var(--shadow-level-1);
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
  color: var(--color-text-secondary);
}

.category-group-label {
  font-weight: 600;
  color: var(--color-text-secondary);
}

/* Skeleton loading */
.skeleton-state {
  padding: 16px;
}

/* Amount Colors */
.income-amount {
  color: var(--color-success);
  font-weight: 600;
  font-size: 16px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

.expense-amount {
  color: var(--color-danger);
  font-weight: 600;
  font-size: 16px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

/* Date Group */
.date-group {
  margin-bottom: 0;
}

.date-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: var(--color-bg-layout);
  border-radius: var(--radius-md);
  margin-bottom: 8px;
}

.date-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.date-total {
  font-size: 14px;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

/* ==================== Ledger Item Card ==================== */
.ledger-item {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  margin-bottom: 8px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow var(--duration-normal) ease;
  min-height: 44px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ledger-item.clickable {
  cursor: pointer;
}

.ledger-item.clickable:hover {
  box-shadow: var(--shadow-level-2);
}

/* Top row: category + creator */
.item-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-category {
  display: flex;
  align-items: center;
  gap: 6px;
}

.item-cat-icon {
  font-size: 18px;
}

.item-cat-name {
  font-size: 14px;
  font-weight: 600;
}

.item-creator-line {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.creator-avatar {
  font-size: 14px;
}

/* Body: note */
.item-note {
  font-size: 14px;
  color: var(--color-text-primary);
  word-break: break-word;
}

/* Bottom: amount */
.item-bottom {
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
}

.item-amount {
  font-weight: 600;
  font-size: 16px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
  flex-shrink: 0;
}

/* Category Grid Picker (in modal) */
.category-grid-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 0;
}

.category-pick-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 8px 12px;
  border-radius: var(--radius-md);
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 64px;
  background: var(--color-bg-container);
}

.category-pick-item:hover {
  border-color: var(--color-border);
  background: var(--color-bg-layout);
}

.category-pick-item.active {
  border-color: var(--color-brand);
  background: var(--color-brand-light);
}

.category-pick-icon {
  font-size: 24px;
  line-height: 1;
}

.category-pick-name {
  font-size: 12px;
  color: var(--color-text-primary);
  white-space: nowrap;
}

/* Infinite scroll sentinel */
.load-sentinel {
  display: flex;
  justify-content: center;
  padding: 16px 0;
}

/* Mobile */
.btn-text-short { display: none; }

@media (max-width: 767px) {
  .btn-text-full { display: none; }
  .btn-text-short { display: inline; }

  .ledger-page {
    padding: var(--space-md);
    padding-bottom: 80px;
  }

  .month-row :deep(.ant-btn-primary) {
    position: fixed;
    bottom: 72px;
    right: 20px;
    z-index: var(--z-overlay);
    width: 56px;
    height: 56px;
    border-radius: 50%;
    font-size: 24px;
    box-shadow: var(--shadow-level-3);
  }
}
</style>
