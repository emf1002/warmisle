<template>
  <div class="dashboard" data-testid="dashboard-page">
    <!-- 月份切换 -->
    <div class="month-switcher">
      <a-button type="text" @click="goPrevMonth" class="month-arrow" aria-label="上个月" data-testid="month-prev">
        ◀
      </a-button>
      <span class="month-text" data-testid="current-month">{{ selectedMonth.format('YYYY年M月') }}</span>
      <a-button type="text" @click="goNextMonth" class="month-arrow" aria-label="下个月" data-testid="month-next">
        ▶
      </a-button>
    </div>

    <!-- 月度统计卡片 - Skeleton -->
    <div v-if="loading" class="summary-grid" data-testid="summary-grid">
      <div v-for="i in 3" :key="i" class="stat-card">
        <div class="stat-card-header">
          <SkeletonCard width="36px" height="36px" borderRadius="10px" />
          <SkeletonCard width="60px" height="14px" />
        </div>
        <SkeletonCard width="120px" height="28px" />
        <SkeletonCard width="100%" height="32px" />
      </div>
    </div>

    <!-- 月度统计卡片 - Loaded -->
    <div v-else class="summary-grid" data-testid="summary-grid">
      <div class="stat-card stat-card--income" data-testid="summary-income">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--income">
            <Icon name="TrendingUp" :size="18" />
          </div>
          <span class="stat-label">本月收入</span>
        </div>
        <div class="stat-amount amount-income">¥{{ (displayIncome / 100).toFixed(2) }}</div>
        <MiniSparkline :data="incomeTrend" color="var(--color-income)" />
      </div>

      <div class="stat-card stat-card--expense" data-testid="summary-expense">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--expense">
            <Icon name="TrendingDown" :size="18" />
          </div>
          <span class="stat-label">本月支出</span>
        </div>
        <div class="stat-amount amount-expense">¥{{ (displayExpense / 100).toFixed(2) }}</div>
        <MiniSparkline :data="expenseTrend" color="var(--color-expense)" />
      </div>

      <div class="stat-card stat-card--balance" data-testid="summary-balance">
        <div class="stat-card-header">
          <div class="stat-icon stat-icon--balance">
            <Icon name="Wallet" :size="18" />
          </div>
          <span class="stat-label">本月结余</span>
        </div>
        <div class="stat-amount amount-balance">¥{{ (displayBalance / 100).toFixed(2) }}</div>
        <div class="stat-trend">
          <span v-if="balanceChange >= 0" class="trend-up">较上月 +{{ balanceChange }}%</span>
          <span v-else class="trend-down">较上月 {{ balanceChange }}%</span>
        </div>
      </div>
    </div>

    <div class="dashboard-grid">
      <!-- 支出分类图表 -->
      <a-card title="支出分类" class="section-card">
        <div v-if="expenseChart.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无数据" />
        </div>
        <div v-else data-testid="expense-chart">
          <BarChart :data="barChartData" color="var(--color-expense)" :height="180" />
        </div>
      </a-card>

      <!-- 近期待办 -->
      <a-card title="近期待办" class="section-card">
        <div v-if="upcomingTodos.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无待办" />
        </div>
        <a-list v-else :data-source="upcomingTodos" size="small" data-testid="upcoming-todos">
          <template #renderItem="{ item: todo }">
            <a-list-item class="todo-item" :data-testid="'todo-link-' + todo.id" @click="$router.push('/todo')">
              <div class="todo-item-content">
                <span class="todo-title">{{ todo.title }}</span>
                <span class="todo-meta">
                  <a-tag v-if="todo.priority === 'urgent'" color="red">紧急</a-tag>
                  <a-tag v-else-if="todo.priority === 'important'" color="orange">重要</a-tag>
                  <span v-if="todo.assignee" class="todo-assignee">{{ todo.assignee.name }}</span>
                  <span v-if="todo.due_date" class="todo-due" :class="{ overdue: isOverdue(todo.due_date) }">
                    {{ formatDate(todo.due_date) }}
                  </span>
                </span>
              </div>
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>

    <div class="bottom-grid" style="margin-top: 16px">
      <!-- 愿望动态 -->
      <a-card title="愿望动态" class="section-card">
        <div v-if="wishTrends.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无愿望" />
        </div>
        <a-list v-else :data-source="wishTrends" size="small" data-testid="wish-trends">
          <template #renderItem="{ item: trend }">
            <a-list-item class="clickable-item" :data-testid="'wish-link-' + trend.id" @click="$router.push('/wish')">
              <a-list-item-meta>
                <template #title>{{ trend.title }}</template>
                <template #description>
                  <span>{{ trend.creator.name }}</span>
                  <span style="margin-left: 8px">{{ trend.vote_count }} 票</span>
                </template>
              </a-list-item-meta>
              <template #extra>
                <a-tag v-if="trend.status === 'pending'" color="default">待定</a-tag>
                <a-tag v-else-if="trend.status === 'agreed'" color="blue">已同意</a-tag>
                <a-tag v-else-if="trend.status === 'achieved'" color="green">已实现</a-tag>
                <a-tag v-else-if="trend.status === 'abandoned'">已放弃</a-tag>
              </template>
            </a-list-item>
          </template>
        </a-list>
      </a-card>

      <!-- 论坛热点 -->
      <a-card title="论坛热点" class="section-card">
        <div v-if="forumHot.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无动态" />
        </div>
        <a-list v-else :data-source="forumHot" size="small" data-testid="forum-hot">
          <template #renderItem="{ item: feed }">
            <a-list-item class="clickable-item" :data-testid="'topic-link-' + feed.id" @click="$router.push('/forum')">
              <a-list-item-meta>
                <template #title>
                  <a-tag v-if="feed.type === 'topic'" color="blue" size="small">话题</a-tag>
                  <a-tag v-else size="small">动态</a-tag>
                  <span>{{ feed.title || feed.content }}</span>
                </template>
                <template #description>{{ feed.creator.name }} · {{ timeAgo(feed.created_at) }}</template>
              </a-list-item-meta>
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { getSummary, getExpenseChart, getUpcomingTodos, getWishTrends, getForumHot } from '@/api/dashboard'
import EmptyState from '@/components/EmptyState.vue'
import Icon from '@/components/Icon.vue'
import MiniSparkline from '@/components/MiniSparkline.vue'
import BarChart from '@/components/BarChart.vue'
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

// Sparkline trend data (placeholder 7-point arrays for visual effect)
const incomeTrend = computed(() => {
  if (expenseChart.value.length === 0) return [0, 0, 0, 0, 0, 0, 0]
  // Generate 7-point trend from summary data for visual placeholder
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

function goPrevMonth() {
  selectedMonth.value = selectedMonth.value.subtract(1, 'month')
  fetchData()
}

function goNextMonth() {
  selectedMonth.value = selectedMonth.value.add(1, 'month')
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
  padding: var(--space-md);
  max-width: var(--max-content-width);
  margin: 0 auto;
}

/* Month Switcher */
.month-switcher {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  gap: 16px;
}

.month-arrow {
  min-width: 44px;
  min-height: 44px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.month-arrow:disabled {
  color: var(--color-text-disabled);
  cursor: not-allowed;
}

.month-text {
  font-size: 16px;
  font-weight: 600;
  min-width: 120px;
  text-align: center;
}

/* Summary Grid */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
}

.stat-card {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-lg);
  border: 1px solid var(--color-border);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-level-2);
}

.stat-card-header {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  margin-bottom: var(--space-sm);
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-icon--income {
  background: var(--color-income-bg);
  color: var(--color-income);
}
.stat-icon--expense {
  background: var(--color-expense-bg);
  color: var(--color-expense);
}
.stat-icon--balance {
  background: var(--color-balance-bg);
  color: var(--color-balance);
}

.stat-label {
  font-size: 13px;
  color: var(--color-muted);
  font-family: var(--font-display);
}

.stat-amount {
  font-size: 28px;
  font-weight: 700;
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  margin-bottom: var(--space-xs);
}
.amount-income { color: var(--color-income); }
.amount-expense { color: var(--color-expense); }
.amount-balance { color: var(--color-balance); }

.stat-trend {
  font-size: 12px;
  color: var(--color-muted);
}
.trend-up { color: var(--color-income); }
.trend-down { color: var(--color-expense); }

/* Dashboard Grids */
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.section-card {
  height: 100%;
  box-shadow: var(--shadow-level-1);
  border-radius: var(--radius-md);
  transition: box-shadow var(--duration-normal) ease;
}

.section-card:hover {
  box-shadow: var(--shadow-level-2);
}

.empty-section {
  padding: 24px 0;
}

/* Todo items in dashboard */
.todo-item {
  cursor: pointer;
  min-height: 44px;
}

.todo-item-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  gap: 8px;
}

.todo-title {
  flex: 1;
  font-size: 14px;
}

.todo-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.todo-assignee {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.todo-due {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.todo-due.overdue {
  color: var(--color-danger);
}

.clickable-item {
  cursor: pointer;
  min-height: 44px;
}

/* Responsive */
@media (max-width: 1023px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 767px) {
  .dashboard {
    padding: 12px;
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .stat-amount {
    font-size: 24px;
  }
}
</style>
