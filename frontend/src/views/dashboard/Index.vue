<template>
  <div class="dashboard" data-testid="dashboard-page">
    <!-- 月份选择器 -->
    <div class="month-row">
      <a-month-picker
        v-model:value="selectedMonth"
        format="YYYY年M月"
        @change="onMonthChange"
        data-testid="month-picker"
      />
    </div>

    <!-- 月度统计卡片 - Skeleton -->
    <div v-if="loading" class="summary-grid" data-testid="summary-grid">
      <div v-for="i in 3" :key="i" class="stat-card stat-card--skeleton">
        <div class="stat-card-header">
          <SkeletonCard width="44px" height="44px" :borderRadius="12" />
          <SkeletonCard width="80px" height="16px" />
        </div>
        <SkeletonCard width="140px" height="36px" />
        <SkeletonCard width="100%" height="40px" />
      </div>
    </div>

    <!-- 月度统计卡片 - Loaded -->
    <div v-else class="summary-grid" data-testid="summary-grid">
      <article class="stat-card stat-card--income" data-testid="summary-income">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--income">
            <Icon name="TrendingUp" :size="20" />
          </div>
          <div class="stat-header-text">
            <span class="stat-label">本月收入</span>
            <Icon name="HelpCircle" :size="14" class="stat-help" title="收入合计" />
          </div>
        </div>
        <div class="stat-amount amount-income">
          <span class="amount-prefix">¥</span>
          <span class="amount-value">{{ (displayIncome / 100).toFixed(2) }}</span>
        </div>
        <div class="stat-chart">
          <MiniSparkline :data="incomeTrend" color="var(--color-income)" />
        </div>
      </article>

      <article class="stat-card stat-card--expense" data-testid="summary-expense">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--expense">
            <Icon name="TrendingDown" :size="20" />
          </div>
          <div class="stat-header-text">
            <span class="stat-label">本月支出</span>
            <Icon name="HelpCircle" :size="14" class="stat-help" title="支出合计" />
          </div>
        </div>
        <div class="stat-amount amount-expense">
          <span class="amount-prefix">¥</span>
          <span class="amount-value">{{ (displayExpense / 100).toFixed(2) }}</span>
        </div>
        <div class="stat-chart">
          <MiniSparkline :data="expenseTrend" color="var(--color-expense)" />
        </div>
      </article>

      <article class="stat-card stat-card--balance" data-testid="summary-balance">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--balance">
            <Icon name="Wallet" :size="20" />
          </div>
          <div class="stat-header-text">
            <span class="stat-label">本月结余</span>
            <Icon name="HelpCircle" :size="14" class="stat-help" title="收入 - 支出" />
          </div>
        </div>
        <div class="stat-amount amount-balance">
          <span class="amount-prefix">¥</span>
          <span class="amount-value">{{ (displayBalance / 100).toFixed(2) }}</span>
        </div>
        <div class="stat-trend">
          <span v-if="balanceChange >= 0" class="trend-badge trend-up">
            <Icon name="ArrowUp" :size="12" />
            较上月 +{{ balanceChange }}%
          </span>
          <span v-else class="trend-badge trend-down">
            <Icon name="ArrowDown" :size="12" />
            较上月 {{ balanceChange }}%
          </span>
        </div>
      </article>
    </div>

    <!-- 主内容网格 -->
    <div class="dashboard-grid">
      <!-- 支出分类图表 -->
      <section class="section-card" data-testid="expense-chart-section">
        <div class="section-header">
          <h2 class="section-title">
            <Icon name="PieChart" :size="18" class="section-icon" />
            支出分类
          </h2>
          <a-button type="link" size="small" @click="$router.push('/ledger')">
            查看详情
            <Icon name="ArrowRight" :size="14" />
          </a-button>
        </div>
        <div v-if="expenseChart.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无支出数据" />
        </div>
        <div v-else class="chart-container" data-testid="expense-chart">
          <PieChart :data="barChartData" :size="180" />
        </div>
      </section>

      <!-- 近期待办 -->
      <section class="section-card" data-testid="upcoming-todos-section">
        <div class="section-header">
          <h2 class="section-title">
            <Icon name="ListTodo" :size="18" class="section-icon" />
            近期待办
          </h2>
          <a-button type="link" size="small" @click="$router.push('/todo')">
            查看全部
            <Icon name="ArrowRight" :size="14" />
          </a-button>
        </div>
        <div v-if="upcomingTodos.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无待办事项" />
        </div>
        <a-list
          v-else
          :data-source="upcomingTodos"
          size="small"
          class="todo-list"
          data-testid="upcoming-todos"
        >
          <template #renderItem="{ item: todo }">
            <a-list-item
              class="todo-item"
              :data-testid="'todo-link-' + todo.id"
              @click="$router.push('/todo')"
            >
              <div class="todo-item-content">
                <a-checkbox
                  :checked="todo.status === 'completed'"
                  class="todo-checkbox"
                  @click.stop
                />
                <span class="todo-title" :class="{ 'todo-completed': todo.status === 'completed' }">
                  {{ todo.title }}
                </span>
                <span class="todo-meta">
                  <a-tag v-if="todo.priority === 'urgent'" color="red" class="priority-tag">紧急</a-tag>
                  <a-tag v-else-if="todo.priority === 'important'" color="orange" class="priority-tag">重要</a-tag>
                  <span v-if="todo.assignee" class="todo-assignee">
                    <Icon name="User" :size="12" />
                    {{ todo.assignee.name }}
                  </span>
                  <span v-if="todo.due_date" class="todo-due" :class="{ overdue: isOverdue(todo.due_date) }">
                    <Icon name="Calendar" :size="12" />
                    {{ formatDate(todo.due_date) }}
                  </span>
                </span>
              </div>
            </a-list-item>
          </template>
        </a-list>
      </section>
    </div>

    <!-- 底部网格 -->
    <div class="bottom-grid" style="margin-top: 24px">
      <!-- 愿望动态 -->
      <section class="section-card" data-testid="wish-trends-section">
        <div class="section-header">
          <h2 class="section-title">
            <Icon name="Star" :size="18" class="section-icon" />
            愿望动态
          </h2>
          <a-button type="link" size="small" @click="$router.push('/wish')">
            查看全部
            <Icon name="ArrowRight" :size="14" />
          </a-button>
        </div>
        <div v-if="wishTrends.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无愿望" />
        </div>
        <a-list
          v-else
          :data-source="wishTrends"
          size="small"
          class="wish-list"
          data-testid="wish-trends"
        >
          <template #renderItem="{ item: trend }">
            <a-list-item
              class="wish-item"
              :data-testid="'wish-link-' + trend.id"
              @click="$router.push('/wish')"
            >
              <a-list-item-meta>
                <template #title>
                  <span class="wish-title">{{ trend.title }}</span>
                </template>
                <template #description>
                  <span class="wish-meta">
                    <Icon name="User" :size="12" />
                    {{ trend.creator.name }}
                  </span>
                  <span class="wish-votes">
                    <Icon name="ThumbsUp" :size="12" />
                    {{ trend.vote_count }} 票
                  </span>
                </template>
              </a-list-item-meta>
              <template #extra>
                <a-tag v-if="trend.status === 'pending'" color="default" class="status-tag">待定</a-tag>
                <a-tag v-else-if="trend.status === 'agreed'" color="blue" class="status-tag">已同意</a-tag>
                <a-tag v-else-if="trend.status === 'achieved'" color="green" class="status-tag">已实现</a-tag>
                <a-tag v-else-if="trend.status === 'abandoned'" class="status-tag">已放弃</a-tag>
              </template>
            </a-list-item>
          </template>
        </a-list>
      </section>

      <!-- 论坛热点 -->
      <section class="section-card" data-testid="forum-hot-section">
        <div class="section-header">
          <h2 class="section-title">
            <Icon name="MessageSquare" :size="18" class="section-icon" />
            论坛热点
          </h2>
          <a-button type="link" size="small" @click="$router.push('/forum')">
            查看全部
            <Icon name="ArrowRight" :size="14" />
          </a-button>
        </div>
        <div v-if="forumHot.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无论坛动态" />
        </div>
        <a-list
          v-else
          :data-source="forumHot"
          size="small"
          class="forum-list"
          data-testid="forum-hot"
        >
          <template #renderItem="{ item: feed }">
            <a-list-item
              class="forum-item"
              :data-testid="'topic-link-' + feed.id"
              @click="$router.push('/forum')"
            >
              <a-list-item-meta>
                <template #title>
                  <div class="forum-title">
                    <a-tag v-if="feed.type === 'topic'" color="blue" size="small" class="type-tag">话题</a-tag>
                    <a-tag v-else size="small" class="type-tag">动态</a-tag>
                    <span>{{ feed.title || feed.content }}</span>
                  </div>
                </template>
                <template #description>
                  <span class="forum-meta">
                    <Icon name="User" :size="12" />
                    {{ feed.creator.name }}
                  </span>
                  <span class="forum-time">
                    <Icon name="Clock" :size="12" />
                    {{ timeAgo(feed.created_at) }}
                  </span>
                </template>
              </a-list-item-meta>
            </a-list-item>
          </template>
        </a-list>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { getSummary, getExpenseChart, getUpcomingTodos, getWishTrends, getForumHot } from '@/api/dashboard'
import EmptyState from '@/components/EmptyState.vue'
import Icon from '@/components/Icon.vue'
import MiniSparkline from '@/components/MiniSparkline.vue'
import PieChart from '@/components/PieChart.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import { formatDate, timeAgo } from '@/utils/format'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

const selectedMonth = ref(dayjs())

const loading = ref(true)

// Summary data (in fen/cents)
const summary = ref({ income: 0, expense: 0, balance: 0 })
const expenseChart = ref<any[]>([])
const upcomingTodos = ref<any[]>([])
const wishTrends = ref<any[]>([])
const forumHot = ref<any[]>([])

// Previous month summary for trend comparison
const prevSummary = ref({ income: 0, expense: 0, balance: 0 })

// Count-up animation display values
const displayIncome = ref(0)
const displayExpense = ref(0)
const displayBalance = ref(0)

const monthParam = computed(() => selectedMonth.value.format('YYYY-MM'))

// Balance change percentage vs previous month
const balanceChange = computed(() => {
  if (prevSummary.value.balance === 0) return 0
  return Math.round(((summary.value.balance - prevSummary.value.balance) / Math.abs(prevSummary.value.balance)) * 100)
})

// Sparkline trend data
const incomeTrend = computed(() => {
  if (expenseChart.value.length === 0) return [0, 0, 0, 0, 0, 0, 0]
  const val = summary.value.income
  return [0.3, 0.5, 0.4, 0.6, 0.7, 0.8, 1.0].map(p => Math.round(val * p * 0.15))
})

const expenseTrend = computed(() => {
  if (expenseChart.value.length === 0) return [0, 0, 0, 0, 0, 0, 0]
  const val = summary.value.expense
  return [0.4, 0.6, 0.5, 0.7, 0.6, 0.8, 1.0].map(p => Math.round(val * p * 0.15))
})

// Bar chart data for expense categories
const barChartData = computed(() => {
  return expenseChart.value.map((item: any) => ({
    label: item.category_name,
    value: item.amount / 100,
  }))
})

function onMonthChange() {
  fetchData()
}

function countUp(target: number, setter: (v: number) => void, duration = 800) {
  const start = performance.now()
  const step = (now: number) => {
    const progress = Math.min((now - start) / duration, 1)
    const eased = 1 - Math.pow(1 - progress, 3)
    setter(Math.round(target * eased))
    if (progress < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}

async function fetchData() {
  const m = monthParam.value

  // Fetch previous month for comparison
  const prevMonth = selectedMonth.value.subtract(1, 'month').format('YYYY-MM')

  const [s, e, t, w, f, ps] = await Promise.all([
    getSummary(m),
    getExpenseChart(m),
    getUpcomingTodos(),
    getWishTrends(),
    getForumHot(),
    getSummary(prevMonth),
  ])
  summary.value = s.data
  expenseChart.value = e.data || []
  upcomingTodos.value = t.data || []
  wishTrends.value = w.data || []
  forumHot.value = f.data || []
  prevSummary.value = ps.data

  loading.value = false

  nextTick(() => {
    countUp(summary.value.income, v => (displayIncome.value = v))
    countUp(summary.value.expense, v => (displayExpense.value = v))
    countUp(summary.value.balance, v => (displayBalance.value = v))
  })
}

function isOverdue(date: string): boolean {
  return dayjs(date).isBefore(dayjs(), 'day')
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  padding: var(--space-lg, 24px);
  max-width: var(--max-content-width, 1200px);
  margin: 0 auto;
  animation: fadeIn var(--duration-normal, 200ms) ease-out;
}

/* ==================== 页面标题区 ==================== */
.page-header {
  margin-bottom: var(--space-lg, 24px);
  animation: slideInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) both;
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs, 8px);
}

.page-title {
  font-size: var(--text-2xl, 1.5rem);
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-secondary);
  margin: 0;
}

/* ==================== 月份选择器 ==================== */
.month-row {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  margin-bottom: var(--space-lg, 24px);
  animation: slideInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) 50ms both;
}

/* ==================== 统计卡片网格 ==================== */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-lg, 24px);
  margin-bottom: var(--space-xl, 32px);
}

.stat-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl, 20px);
  padding: var(--space-lg, 24px);
  border: 1px solid var(--color-border);
  transition: all var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  animation: slideInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) both;
}

.stat-card:nth-child(1) { animation-delay: 100ms; }
.stat-card:nth-child(2) { animation-delay: 150ms; }
.stat-card:nth-child(3) { animation-delay: 200ms; }

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--color-text-disabled, #B8AFA8);
  opacity: 0;
  transition: opacity var(--duration-normal, 200ms) ease-out;
}

.stat-card--income::before {
  background: var(--color-income, #4CAF50);
}

.stat-card--expense::before {
  background: var(--color-expense, #F44336);
}

.stat-card--balance::before {
  background: var(--color-balance, #2196F3);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card-hover);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-card--skeleton {
  min-height: 180px;
}

/* ==================== 卡片头部 ==================== */
.stat-card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-sm, 12px);
  margin-bottom: var(--space-md, 16px);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md, 12px);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform var(--duration-normal, 200ms) cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.stat-card:hover .stat-icon {
  transform: scale(1.1);
}

.stat-icon--income {
  background: var(--color-income-bg, rgba(76, 175, 80, 0.1));
  color: var(--color-income, #4CAF50);
}

.stat-icon--expense {
  background: var(--color-expense-bg, rgba(244, 67, 54, 0.1));
  color: var(--color-expense, #F44336);
}

.stat-icon--balance {
  background: var(--color-balance-bg, rgba(33, 150, 243, 0.1));
  color: var(--color-balance, #2196F3);
}

.stat-header-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.stat-label {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-secondary);
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs, 8px);
}

.stat-help {
  color: var(--color-text-disabled);
  cursor: help;
  transition: color var(--duration-fast, 150ms) ease-out;
}

.stat-help:hover {
  color: var(--color-brand);
}

/* ==================== 金额显示 ==================== */
.stat-amount {
  font-family: var(--font-display);
  font-weight: 700;
  margin-bottom: var(--space-sm, 12px);
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.amount-income { color: var(--color-income, #4CAF50); }
.amount-expense { color: var(--color-expense, #F44336); }
.amount-balance { color: var(--color-balance, #2196F3); }

.amount-prefix {
  font-size: var(--text-lg, 1.125rem);
  font-weight: 600;
}

.amount-value {
  font-size: var(--text-3xl, 1.875rem);
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

/* ==================== 趋势图表 ==================== */
.stat-chart {
  height: 40px;
  margin-top: var(--space-xs, 8px);
}

.stat-trend {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-secondary);
}

.trend-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-xs, 8px);
  padding: var(--space-xs, 8px) var(--space-sm, 12px);
  border-radius: var(--radius-full, 9999px);
  font-weight: 600;
  font-size: var(--text-xs, 0.75rem);
}

.trend-up {
  background: var(--color-income-bg, rgba(76, 175, 80, 0.1));
  color: var(--color-income, #4CAF50);
}

.trend-down {
  background: var(--color-expense-bg, rgba(244, 67, 54, 0.1));
  color: var(--color-expense, #F44336);
}

/* ==================== 内容区域网格 ==================== */
.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-lg, 24px);
  margin-bottom: var(--space-lg, 24px);
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-lg, 24px);
}

/* ==================== 区块卡片 ==================== */
.section-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-xl, 20px);
  border: 1px solid var(--color-border);
  padding: var(--space-lg, 24px);
  transition: all var(--duration-normal, 200ms) cubic-bezier(0.4, 0, 0.2, 1);
  animation: slideInUp var(--duration-slow, 300ms) cubic-bezier(0.4, 0, 0.2, 1) both;
}

.dashboard-grid .section-card:nth-child(1) { animation-delay: 250ms; }
.dashboard-grid .section-card:nth-child(2) { animation-delay: 300ms; }

.bottom-grid .section-card:nth-child(1) { animation-delay: 350ms; }
.bottom-grid .section-card:nth-child(2) { animation-delay: 400ms; }

.section-card:hover {
  box-shadow: var(--shadow-level-2);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-md, 16px);
  padding-bottom: var(--space-md, 16px);
  border-bottom: 1px solid var(--color-border-secondary);
}

.section-title {
  font-size: var(--text-lg, 1.125rem);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  display: inline-flex;
  align-items: center;
  gap: var(--space-sm, 12px);
}

.section-icon {
  color: var(--color-brand);
}

.chart-container {
  height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-section {
  height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ==================== 待办列表 ==================== */
.todo-list {
  max-height: 400px;
  overflow-y: auto;
}

.todo-item {
  cursor: pointer;
  min-height: 44px;
  padding: var(--space-sm, 12px) !important;
  border-radius: var(--radius-md, 12px);
  transition: all var(--duration-fast, 150ms) ease-out;
  margin-bottom: var(--space-xs, 8px);
}

.todo-item:hover {
  background: var(--color-bg-subtle, #FAF7F5);
}

.todo-item-content {
  display: flex;
  align-items: center;
  gap: var(--space-sm, 12px);
  width: 100%;
}

.todo-checkbox {
  flex-shrink: 0;
}

.todo-title {
  flex: 1;
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-primary);
  transition: all var(--duration-fast, 150ms) ease-out;
}

.todo-completed {
  text-decoration: line-through;
  color: var(--color-text-disabled);
}

.todo-meta {
  display: flex;
  align-items: center;
  gap: var(--space-xs, 8px);
  flex-shrink: 0;
}

.priority-tag {
  font-size: var(--text-xs, 0.75rem);
  border-radius: var(--radius-xs, 4px);
}

.todo-assignee,
.todo-due {
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.todo-due.overdue {
  color: var(--color-danger, #E85D5D);
  font-weight: 600;
}

/* ==================== 愿望列表 ==================== */
.wish-list {
  max-height: 400px;
  overflow-y: auto;
}

.wish-item {
  cursor: pointer;
  min-height: 44px;
  padding: var(--space-sm, 12px) !important;
  border-radius: var(--radius-md, 12px);
  transition: all var(--duration-fast, 150ms) ease-out;
  margin-bottom: var(--space-xs, 8px);
}

.wish-item:hover {
  background: var(--color-bg-subtle, #FAF7F5);
}

.wish-title {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-primary);
  font-weight: 500;
}

.wish-meta,
.wish-votes {
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.status-tag {
  font-size: var(--text-xs, 0.75rem);
  border-radius: var(--radius-xs, 4px);
}

/* ==================== 论坛列表 ==================== */
.forum-list {
  max-height: 400px;
  overflow-y: auto;
}

.forum-item {
  cursor: pointer;
  min-height: 44px;
  padding: var(--space-sm, 12px) !important;
  border-radius: var(--radius-md, 12px);
  transition: all var(--duration-fast, 150ms) ease-out;
  margin-bottom: var(--space-xs, 8px);
}

.forum-item:hover {
  background: var(--color-bg-subtle, #FAF7F5);
}

.forum-title {
  display: flex;
  align-items: center;
  gap: var(--space-xs, 8px);
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text-primary);
}

.type-tag {
  font-size: var(--text-xs, 0.75rem);
  border-radius: var(--radius-xs, 4px);
}

.forum-meta,
.forum-time {
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-text-secondary);
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

/* ==================== 动画 ==================== */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ==================== 响应式 ==================== */
@media (max-width: 1023px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 767px) {
  .dashboard {
    padding: var(--space-md, 16px);
  }

  .summary-grid {
    grid-template-columns: 1fr;
    gap: var(--space-md, 16px);
  }

  .dashboard-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
    gap: var(--space-md, 16px);
  }

  .stat-amount {
    font-size: var(--text-2xl, 1.5rem);
  }

  .page-title {
    font-size: var(--text-xl, 1.25rem);
  }

  .section-card {
    padding: var(--space-md, 16px);
  }
}
</style>
