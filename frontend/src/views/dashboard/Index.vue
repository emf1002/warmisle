<template>
  <div class="dashboard">
    <!-- 月份切换 -->
    <div class="month-switcher">
      <a-button type="text" @click="goPrevMonth" class="month-arrow" aria-label="上个月">
        ◀
      </a-button>
      <span class="month-text">{{ selectedMonth.format('YYYY年M月') }}</span>
      <a-button type="text" @click="goNextMonth" class="month-arrow" aria-label="下个月">
        ▶
      </a-button>
    </div>

    <!-- 月度统计卡片 -->
    <div class="summary-grid">
      <a-card :bordered="false" class="stat-card">
        <a-statistic title="收入" :value="summary.income / 100" :precision="2">
          <template #prefix><span class="income-prefix">+¥</span></template>
        </a-statistic>
      </a-card>
      <a-card :bordered="false" class="stat-card">
        <a-statistic title="支出" :value="summary.expense / 100" :precision="2">
          <template #prefix><span class="expense-prefix">-¥</span></template>
        </a-statistic>
      </a-card>
      <a-card :bordered="false" class="stat-card">
        <a-statistic title="结余" :value="Math.abs(summary.balance / 100)" :precision="2">
          <template #prefix>
            <span :class="summary.balance >= 0 ? 'income-prefix' : 'expense-prefix'">
              {{ summary.balance >= 0 ? '+' : '-' }}¥
            </span>
          </template>
        </a-statistic>
      </a-card>
    </div>

    <div class="dashboard-grid">
      <!-- 支出分类图表 -->
      <a-card title="支出分类" class="section-card">
        <div v-if="expenseChart.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无数据" />
        </div>
        <div v-else class="expense-chart">
          <div v-for="item in expenseChart" :key="item.category_id" class="chart-row">
            <span class="chart-icon">{{ item.icon }}</span>
            <span class="chart-name">{{ item.category_name }}</span>
            <div class="chart-bar-wrapper">
              <div class="chart-bar" :style="{ width: getPercent(item.amount) + '%' }"></div>
            </div>
            <span class="chart-amount">-¥{{ (item.amount / 100).toFixed(2) }}</span>
          </div>
        </div>
      </a-card>

      <!-- 近期待办 -->
      <a-card title="近期待办" class="section-card">
        <div v-if="upcomingTodos.length === 0" class="empty-section">
          <EmptyState type="no-data" description="暂无待办" />
        </div>
        <a-list v-else :data-source="upcomingTodos" size="small">
          <template #renderItem="{ item: todo }">
            <a-list-item class="todo-item" @click="$router.push('/todo')">
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
        <a-list v-else :data-source="wishTrends" size="small">
          <template #renderItem="{ item: trend }">
            <a-list-item class="clickable-item" @click="$router.push('/wish')">
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
        <a-list v-else :data-source="forumHot" size="small">
          <template #renderItem="{ item: feed }">
            <a-list-item class="clickable-item" @click="$router.push('/forum')">
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
import { ref, onMounted, computed } from 'vue'
import { getSummary, getExpenseChart, getUpcomingTodos, getWishTrends, getForumHot } from '@/api/dashboard'
import EmptyState from '@/components/EmptyState.vue'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

const selectedMonth = ref(dayjs())

const summary = ref({ income: 0, expense: 0, balance: 0 })
const expenseChart = ref<any[]>([])
const upcomingTodos = ref<any[]>([])
const wishTrends = ref<any[]>([])
const forumHot = ref<any[]>([])

const monthParam = computed(() => selectedMonth.value.format('YYYY-MM'))

function goPrevMonth() {
  selectedMonth.value = selectedMonth.value.subtract(1, 'month')
  fetchData()
}

function goNextMonth() {
  selectedMonth.value = selectedMonth.value.add(1, 'month')
  fetchData()
}

async function fetchData() {
  const m = monthParam.value
  const [s, e, t, w, f] = await Promise.all([
    getSummary(m),
    getExpenseChart(m),
    getUpcomingTodos(),
    getWishTrends(),
    getForumHot(),
  ])
  summary.value = s.data
  expenseChart.value = e.data || []
  upcomingTodos.value = t.data || []
  wishTrends.value = w.data || []
  forumHot.value = f.data || []
}

function getPercent(amount: number): number {
  const max = Math.max(...expenseChart.value.map((i: any) => i.amount), 1)
  return (amount / max) * 100
}

function isOverdue(date: string): boolean {
  return dayjs(date).isBefore(dayjs(), 'day')
}

function formatDate(date: string): string {
  const d = dayjs(date)
  const today = dayjs()
  if (d.isSame(today, 'day')) return '今天'
  if (d.isSame(today.subtract(1, 'day'), 'day')) return '昨天'
  return d.format('M月D日')
}

function timeAgo(date: string): string {
  const d = dayjs(date)
  const now = dayjs()
  const diffMins = now.diff(d, 'minute')
  if (diffMins < 60) return `${diffMins}分钟前`
  const diffHours = now.diff(d, 'hour')
  if (diffHours < 24) return `${diffHours}小时前`
  const diffDays = now.diff(d, 'day')
  if (diffDays < 7) return `${diffDays}天前`
  return d.format('M月D日')
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
  gap: 16px;
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
  box-shadow: var(--shadow-level-1);
  border-radius: var(--radius-md);
  transition: box-shadow var(--duration-normal) ease;
}

.stat-card:hover {
  box-shadow: var(--shadow-level-2);
}

.stat-card :deep(.ant-statistic-content-value) {
  font-size: 28px !important;
  font-weight: 700 !important;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

.income-prefix {
  color: var(--color-success);
  font-size: 14px;
}

.expense-prefix {
  color: var(--color-danger);
  font-size: 14px;
}

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

/* Expense Chart */
.expense-chart {
  padding: 8px 0;
}

.chart-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  gap: 8px;
}

.chart-icon {
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.chart-name {
  width: 80px;
  flex-shrink: 0;
  font-size: 13px;
  color: var(--color-text-primary);
}

.chart-bar-wrapper {
  flex: 1;
  height: 12px;
  background: var(--color-border-secondary);
  border-radius: 6px;
  overflow: hidden;
}

.chart-bar {
  height: 100%;
  background: linear-gradient(90deg, var(--color-danger), var(--color-brand));
  border-radius: 6px;
  min-width: 2px;
  transition: width var(--duration-slow) ease;
}

.chart-amount {
  width: 100px;
  flex-shrink: 0;
  text-align: right;
  font-size: 13px;
  color: var(--color-danger);
  font-weight: 500;
  font-variant-numeric: tabular-nums;
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
    gap: 8px;
  }

  .dashboard-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .stat-card :deep(.ant-statistic-content-value) {
    font-size: 24px !important;
  }
}
</style>
