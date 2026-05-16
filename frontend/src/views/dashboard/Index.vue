<template>
  <div class="dashboard">
    <!-- 月份选择 -->
    <div class="month-selector">
      <span class="month-label">月份：</span>
      <a-month-picker v-model:value="selectedMonth" format="YYYY-MM" placeholder="选择月份" @change="onMonthChange" />
    </div>

    <!-- 月度统计卡片 -->
    <a-row :gutter="[16, 16]" class="summary-cards">
      <a-col :xs="8" :sm="8">
        <a-card :bordered="false" class="stat-card">
          <a-statistic title="收入" :value="summary.income / 100" :precision="2" prefix="+¥">
            <template #prefix>
              <span class="income-prefix">+¥</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="8" :sm="8">
        <a-card :bordered="false" class="stat-card">
          <a-statistic title="支出" :value="summary.expense / 100" :precision="2" prefix="-¥">
            <template #prefix>
              <span class="expense-prefix">-¥</span>
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="8" :sm="8">
        <a-card :bordered="false" class="stat-card">
          <a-statistic title="结余" :value="Math.abs(summary.balance / 100)" :precision="2" :prefix="summary.balance >= 0 ? '+¥' : '-¥'">
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]">
      <!-- 支出分类图表 -->
      <a-col :xs="24" :lg="14">
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
      </a-col>

      <!-- 近期待办 -->
      <a-col :xs="24" :lg="10">
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
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <!-- 愿望动态 -->
      <a-col :xs="24" :lg="12">
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
      </a-col>

      <!-- 论坛热点 -->
      <a-col :xs="24" :lg="12">
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
      </a-col>
    </a-row>
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

function onMonthChange() {
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
  padding: 16px;
  max-width: 1200px;
  margin: 0 auto;
}

.month-selector {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.month-label {
  margin-right: 8px;
  font-size: 14px;
  color: #666;
}

.summary-cards {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

.stat-card :deep(.ant-statistic-content-value) {
  font-size: 28px !important;
  font-weight: 700 !important;
}

.stat-card :deep(.ant-statistic-content) {
  font-size: 28px;
}

.income-prefix {
  color: #52c41a;
  font-size: 14px;
}

.expense-prefix {
  color: #ff4d4f;
  font-size: 14px;
}

.section-card {
  height: 100%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
}

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
  color: #333;
}

.chart-bar-wrapper {
  flex: 1;
  height: 12px;
  background: #f0f0f0;
  border-radius: 6px;
  overflow: hidden;
}

.chart-bar {
  height: 100%;
  background: linear-gradient(90deg, #ff4d4f, #ff7875);
  border-radius: 6px;
  min-width: 2px;
  transition: width 0.3s ease;
}

.chart-amount {
  width: 100px;
  flex-shrink: 0;
  text-align: right;
  font-size: 13px;
  color: #ff4d4f;
  font-weight: 500;
}

.empty-section {
  padding: 24px 0;
}

.todo-item {
  cursor: pointer;
  min-height: 44px;
}

.todo-item:hover {
  background: #fafafa;
}

@media (max-width: 767px) {
  .dashboard {
    padding: 12px;
  }

  .stat-card :deep(.ant-statistic-content-value) {
    font-size: 24px !important;
  }

  .stat-card :deep(.ant-statistic-content) {
    font-size: 24px;
  }
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
  color: #999;
}

.todo-due {
  font-size: 12px;
  color: #999;
}

.todo-due.overdue {
  color: #ff4d4f;
}

.clickable-item {
  cursor: pointer;
  min-height: 44px;
}

.clickable-item:hover {
  background: #fafafa;
}
</style>
