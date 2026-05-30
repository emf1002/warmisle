<template>
  <div class="todo-page" data-testid="todo-page">
    <div class="page-header">
      <h2>待办管理</h2>
      <a-button type="primary" @click="openCreate" data-testid="add-btn">新建待办</a-button>
    </div>

    <div class="filter-row">
      <a-select
        v-model:value="filters.status"
        placeholder="全部状态"
        allow-clear
        style="width: 140px"
        @change="onFilterChange"
        data-testid="status-filter"
      >
        <a-select-option value="pending">待完成</a-select-option>
        <a-select-option value="completed">已完成</a-select-option>
      </a-select>
      <a-select
        v-model:value="filters.assignee_id"
        placeholder="全部成员"
        allow-clear
        style="width: 160px"
        @change="onFilterChange"
      >
        <a-select-option v-for="m in members" :key="m.id" :value="m.id">
          {{ m.avatar }} {{ m.name }}
        </a-select-option>
      </a-select>
      <a-button size="small" @click="clearFilters" data-testid="clear-filters">清除筛选</a-button>
    </div>

    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <div v-else-if="todos.length === 0" class="empty-state">
      <EmptyState v-if="!hasFilters" type="no-data" description="暂无待办事项">
        <template #action>
          <a-button type="primary" @click="openCreate">新建待办</a-button>
        </template>
      </EmptyState>
      <EmptyState v-else type="no-result" description="没有找到匹配的待办事项" @clear="clearFilters" />
    </div>

    <div v-else class="todo-list" data-testid="todo-list">
      <div
        v-for="(todo, index) in todos"
        :key="todo.id"
        class="todo-item card-stagger"
        :class="{ completed: todo.status === 'completed' }"
        :style="{ animationDelay: `${index * 50}ms` }"
        :data-testid="'todo-item-' + todo.id"
      >
        <div class="todo-main">
          <a-checkbox
            :checked="todo.status === 'completed'"
            class="todo-checkbox"
            data-testid="todo-checkbox"
            @change="handleToggle(todo)"
          />
          <div class="todo-content" @click="canEdit(todo) ? openEdit(todo) : undefined">
            <div class="todo-title-row">
              <span class="todo-title" :class="{ 'title-done': todo.status === 'completed' }">
                {{ todo.title }}
              </span>
              <a-tag :color="priorityColor(todo.priority)" class="todo-priority">
                {{ priorityLabel(todo.priority) }}
              </a-tag>
            </div>
            <div v-if="todo.description" class="todo-desc">{{ todo.description }}</div>
            <div class="todo-meta">
              <span class="todo-assignee-line">
                <span class="meta-avatar">{{ todo.creator.avatar }}</span>
                <span v-if="todo.assignee">
                  <span class="meta-arrow">→</span>
                  <span class="meta-avatar">{{ todo.assignee.avatar }}</span>
                  <span class="meta-name">{{ todo.assignee.name }}</span>
                </span>
                <span v-else class="todo-unassigned">
                  未指派
                  <a-button type="link" size="small" @click.stop="handleClaim(todo)" data-testid="claim-btn">认领</a-button>
                </span>
              </span>
              <span v-if="todo.due_date" class="todo-due" :class="{ overdue: isOverdue(todo) }">
                📅 {{ formatDate(todo.due_date) }}
              </span>
            </div>
          </div>
          <div v-if="canEdit(todo)" class="todo-actions">
            <a-button type="link" size="small" @click.stop="openEdit(todo)" data-testid="edit-btn">编辑</a-button>
            <a-button type="link" size="small" danger @click.stop="confirmDelete(todo)" data-testid="delete-btn">删除</a-button>
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
        @change="fetchTodos"
      />
    </div>

    <a-modal
      v-model:open="dialogOpen"
      :title="editingTodo ? '编辑待办' : '新建待办'"
      :confirm-loading="submitting"
      width="480px"
      @cancel="dialogOpen = false"
      data-testid="todo-modal"
    >
      <template #footer>
        <div style="display: flex; justify-content: space-between">
          <a-button v-if="editingTodo" danger @click="confirmDelete(editingTodo)">删除</a-button>
          <div style="display: flex; gap: 8px">
            <a-button @click="dialogOpen = false">取消</a-button>
            <a-button type="primary" :loading="submitting" @click="handleSubmit" data-testid="submit-btn">
              {{ editingTodo ? '保存' : '创建' }}
            </a-button>
          </div>
        </div>
      </template>
      <a-form :model="form" layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model:value="form.title" :maxlength="100" placeholder="请输入待办标题" data-testid="title-input" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" :maxlength="500" :rows="3" placeholder="可选，最多500字" data-testid="desc-input" />
        </a-form-item>
        <a-form-item label="优先级">
          <a-select v-model:value="form.priority" data-testid="priority-select">
            <a-select-option value="normal">普通</a-select-option>
            <a-select-option value="important">重要</a-select-option>
            <a-select-option value="urgent">紧急</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="负责人">
          <a-select v-model:value="form.assignee_id" placeholder="选择负责人" allow-clear data-testid="assignee-select">
            <a-select-option v-for="m in members" :key="m.id" :value="m.id">
              {{ m.avatar }} {{ m.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="截止日期">
          <a-date-picker v-model:value="form.due_date" format="YYYY-MM-DD" style="width: 100%" data-testid="due-date-picker" />
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
import { getTodoList, createTodo, updateTodo, deleteTodo, toggleTodo, claimTodo } from '@/api/todo'
import { getMembers } from '@/api/member'
import EmptyState from '@/components/EmptyState.vue'
import { useAuthStore } from '@/stores/auth'

interface Member {
  id: number
  name: string
  avatar: string
}

interface TodoItem {
  id: number
  title: string
  description: string
  priority: string
  status: string
  assignee_id: number | null
  creator_id: number
  due_date: string | null
  completed_at: string | null
  created_at: string
  assignee: Member | null
  creator: Member
}

const authStore = useAuthStore()

const loading = ref(false)
const submitting = ref(false)
const dialogOpen = ref(false)
const editingTodo = ref<TodoItem | null>(null)
const members = ref<Member[]>([])
const todos = ref<TodoItem[]>([])
const total = ref(0)
const pageSize = 20
const page = ref(1)

const filters = reactive({
  status: undefined as string | undefined,
  assignee_id: undefined as number | undefined,
})

const hasFilters = computed(() => filters.status || filters.assignee_id)

function canEdit(todo: TodoItem): boolean {
  if (authStore.currentUserRole === 'admin') return true
  if (todo.creator_id === authStore.currentUserId) return true
  if (todo.assignee_id === authStore.currentUserId) return true
  return false
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

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  const target = new Date(d.getFullYear(), d.getMonth(), d.getDate())
  if (target.getTime() === today.getTime()) return '今天截止'
  if (target.getTime() === yesterday.getTime()) return '昨天截止'
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  return `${d.getMonth() + 1}月${d.getDate()}日 周${weekDays[d.getDay()]}截止`
}

function isOverdue(todo: TodoItem): boolean {
  if (!todo.due_date || todo.status === 'completed') return false
  const due = new Date(todo.due_date)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return due < today
}

async function fetchTodos() {
  loading.value = true
  try {
    const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
    if (filters.status) params.status = filters.status
    if (filters.assignee_id) params.assignee_id = filters.assignee_id
    const res: any = await getTodoList(params as any)
    const data = res.data
    todos.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

function onFilterChange() {
  page.value = 1
  fetchTodos()
}

function clearFilters() {
  filters.status = undefined
  filters.assignee_id = undefined
  page.value = 1
  fetchTodos()
}

function openCreate() {
  editingTodo.value = null
  form.title = ''
  form.description = ''
  form.priority = 'normal'
  form.assignee_id = undefined
  form.due_date = null
  dialogOpen.value = true
}

function openEdit(todo: TodoItem) {
  editingTodo.value = todo
  form.title = todo.title
  form.description = todo.description || ''
  form.priority = todo.priority
  form.assignee_id = todo.assignee_id || undefined
  form.due_date = todo.due_date ? dayjs(todo.due_date) : null
  dialogOpen.value = true
}

const form = reactive({
  title: '',
  description: '',
  priority: 'normal' as string,
  assignee_id: undefined as number | undefined,
  due_date: null as Dayjs | null,
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
      priority: form.priority,
    }
    if (form.assignee_id) payload.assignee_id = form.assignee_id
    if (form.due_date) payload.due_date = form.due_date.format('YYYY-MM-DD')
    if (editingTodo.value) {
      await updateTodo(editingTodo.value.id, payload)
      message.success('✅ 更新成功')
    } else {
      await createTodo(payload)
      message.success('✅ 创建成功')
    }
    dialogOpen.value = false
    fetchTodos()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  } finally {
    submitting.value = false
  }
}

async function handleToggle(todo: TodoItem) {
  try {
    await toggleTodo(todo.id)
    message.success(todo.status === 'pending' ? '✅ 已完成' : '✅ 已恢复')
    fetchTodos()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }
}

async function handleClaim(todo: TodoItem) {
  try {
    await claimTodo(todo.id)
    message.success('✅ 认领成功')
    fetchTodos()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }
}

function confirmDelete(todo: TodoItem) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: `确定要删除待办「${todo.title}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteTodo(todo.id)
        message.success('✅ 删除成功')
        dialogOpen.value = false
        fetchTodos()
      } catch (e: any) {
        if (e?.response?.data?.message) message.error(e.response.data.message)
        throw e
      }
    },
  })
}

onMounted(async () => {
  try {
    const res: any = await getMembers()
    members.value = res.data || []
  } catch { /* handled by interceptor */ }
  fetchTodos()
})
</script>

<style scoped>
.todo-page {
  padding: var(--space-lg);
  max-width: 800px;
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
  gap: var(--space-xs);
  margin-bottom: var(--space-md);
  flex-wrap: wrap;
  align-items: center;
  padding: var(--space-sm);
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-level-1);
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

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* ==================== Todo Item Card ==================== */
.todo-item {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow var(--duration-normal) ease, opacity var(--duration-normal) ease;
}

.todo-item:hover {
  box-shadow: var(--shadow-level-2);
}

.todo-item.completed {
  background: var(--color-bg-layout);
  opacity: 0.75;
}

.todo-main {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.todo-checkbox {
  min-height: 44px;
  display: flex;
  align-items: flex-start;
  padding-top: 2px;
}

.todo-content {
  flex: 1;
  min-width: 0;
  cursor: default;
}

.todo-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
}

.todo-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  word-break: break-word;
}

.todo-title.title-done {
  text-decoration: line-through;
  color: var(--color-text-disabled);
}

.todo-priority {
  flex-shrink: 0;
}

.todo-desc {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Meta row */
.todo-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
  flex-wrap: wrap;
  min-height: 24px;
}

.todo-assignee-line {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.meta-avatar {
  font-size: 14px;
}

.meta-arrow {
  color: var(--color-muted);
  margin: 0 2px;
}

.meta-name {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.todo-unassigned {
  font-size: 12px;
  color: var(--color-muted);
  font-style: italic;
}

.todo-due {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.todo-due.overdue {
  color: var(--color-danger);
  font-weight: 500;
}

.todo-actions {
  display: flex;
  align-items: center;
  gap: 0;
  flex-shrink: 0;
  min-height: 44px;
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 767px) {
  .todo-page {
    padding: var(--space-md);
  }
}
</style>
