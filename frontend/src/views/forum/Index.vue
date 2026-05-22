<template>
  <div class="forum-page" data-testid="forum-page">
    <div class="page-header">
      <h2>家庭论坛</h2>
      <div class="header-actions">
        <a-button @click="openCreatePost" data-testid="create-post-btn">发动态</a-button>
        <a-button type="primary" @click="openCreateTopic" data-testid="create-topic-btn">发话题</a-button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <template v-else-if="feedItems.length === 0 && pinnedItems.length === 0">
      <EmptyState type="no-data" description="暂无动态">
        <template #action>
          <a-button type="primary" @click="openCreatePost">发第一条动态</a-button>
        </template>
      </EmptyState>
    </template>

    <template v-else>
      <!-- Pinned topics -->
      <div v-if="pinnedItems.length > 0" class="pinned-section">
        <div v-for="item in pinnedItems" :key="'pinned-' + item.id" class="feed-card topic-card">
          <div class="feed-header">
            <span class="topic-pin-badge">📌 公告</span>
          </div>
          <h3 class="topic-title" @click="goToDetail(item)">{{ item.title }}</h3>
          <p class="feed-content topic-excerpt">{{ truncate(item.content, 150) }}</p>
          <div class="topic-footer">
            <span v-if="item.tag" class="topic-tag">#{{ item.tag.name }}</span>
            <span class="feed-author topic-author">
              <span class="feed-avatar" :aria-label="`${item.creator?.name || '用户'}的头像`">{{ item.creator?.avatar || '👤' }}</span>
              <span class="feed-name">{{ item.creator?.name }}</span>
              <span class="feed-time">{{ timeAgo(item.created_at) }}</span>
            </span>
          </div>
          <div class="feed-actions">
            <span class="feed-action" @click.stop="handleLike(item)" data-testid="like-btn">
              <span>{{ item.is_liked ? '❤️' : '🤍' }}</span>
              <span>{{ item.like_count || 0 }}</span>
            </span>
            <span class="feed-action" @click.stop="goToDetail(item)" data-testid="comment-btn">
              💬 {{ item.comment_count || 0 }}
            </span>
            <a-dropdown v-if="canManage(item)" trigger="click">
              <span class="feed-action" @click.stop>⋯</span>
              <template #overlay>
                <a-menu @click="(e: any) => onMenuClick(e, item)">
                  <a-menu-item v-if="isAdmin" key="pin">
                    {{ item.is_pinned ? '取消置顶' : '置顶' }}
                  </a-menu-item>
                  <a-menu-item key="edit">编辑</a-menu-item>
                  <a-menu-item key="delete" danger>删除</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>

      <!-- Feed items -->
      <div class="feed-list" data-testid="feed-list">
        <div v-for="(item, index) in feedItems" :key="item.type + '-' + item.id">
          <!-- Post card -->
          <div v-if="item.type === 'post'" class="feed-card card-stagger" :style="{ animationDelay: `${index * 50}ms` }" :data-testid="'feed-card-' + item.id">
            <div class="feed-header">
              <span class="feed-author">
                <span class="feed-avatar" :aria-label="`${item.creator?.name || '用户'}的头像`">{{ item.creator?.avatar || '👤' }}</span>
                <span class="feed-name">{{ item.creator?.name }}</span>
              </span>
              <span class="feed-time">{{ timeAgo(item.created_at) }}</span>
            </div>
            <div class="feed-body">
              <p class="feed-content">{{ truncate(item.content, 200) }}</p>
            </div>
            <div class="feed-actions">
              <span class="feed-action" :class="{ 'like-bounce': likingItems[`${item.type}_${item.id}`] }" @click.stop="handleLike(item)" data-testid="like-btn">
                <span>{{ item.is_liked ? '❤️' : '🤍' }}</span>
                <span>{{ item.like_count || 0 }}</span>
              </span>
              <span class="feed-action" @click.stop="goToDetail(item)" data-testid="comment-btn">
                💬 {{ item.comment_count || 0 }}
              </span>
              <a-dropdown v-if="canManage(item)" trigger="click">
                <span class="feed-action" @click.stop>⋯</span>
                <template #overlay>
                  <a-menu @click="(e: any) => onMenuClick(e, item)">
                    <a-menu-item key="edit">编辑</a-menu-item>
                    <a-menu-item key="delete" danger>删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>

          <!-- Topic card -->
          <div v-else-if="item.type === 'topic'" class="feed-card topic-card card-stagger" :style="{ animationDelay: `${index * 50}ms` }" :data-testid="'feed-card-' + item.id">
            <div v-if="item.is_pinned" class="feed-header">
              <span class="topic-pin-badge">📌 公告</span>
            </div>
            <h3 class="topic-title" @click="goToDetail(item)">{{ item.title }}</h3>
            <p class="feed-content topic-excerpt">{{ truncate(item.content, 150) }}</p>
            <div class="topic-footer">
              <span v-if="item.tag" class="topic-tag">#{{ item.tag.name }}</span>
              <span class="feed-author topic-author">
                <span class="feed-avatar" :aria-label="`${item.creator?.name || '用户'}的头像`">{{ item.creator?.avatar || '👤' }}</span>
                <span class="feed-name">{{ item.creator?.name }}</span>
                <span class="feed-time">{{ timeAgo(item.created_at) }}</span>
              </span>
            </div>
            <div class="feed-actions">
              <span class="feed-action" :class="{ 'like-bounce': likingItems[`${item.type}_${item.id}`] }" @click.stop="handleLike(item)" data-testid="like-btn">
                <span>{{ item.is_liked ? '❤️' : '🤍' }}</span>
                <span>{{ item.like_count || 0 }}</span>
              </span>
              <span class="feed-action" @click.stop="goToDetail(item)" data-testid="comment-btn">
                💬 {{ item.comment_count || 0 }}
              </span>
              <a-dropdown v-if="canManage(item)" trigger="click">
                <span class="feed-action" @click.stop>⋯</span>
                <template #overlay>
                  <a-menu @click="(e: any) => onMenuClick(e, item)">
                    <a-menu-item v-if="isAdmin" key="pin">
                      {{ item.is_pinned ? '取消置顶' : '置顶' }}
                    </a-menu-item>
                    <a-menu-item key="edit">编辑</a-menu-item>
                    <a-menu-item key="delete" danger>删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="pagination-row">
        <a-pagination
          v-model:current="page"
          :total="total"
          :page-size="pageSize"
          size="small"
          @change="fetchFeed"
        />
      </div>
    </template>

    <!-- Create/Edit Post Dialog -->
    <a-modal
      v-model:open="postDialogOpen"
      :title="editingPostItem ? '编辑动态' : '发动态'"
      :confirm-loading="postSubmitting"
      @ok="handlePostSubmit"
      @cancel="postDialogOpen = false"
      data-testid="post-modal"
    >
      <a-textarea
        v-model:value="postForm.content"
        :maxlength="1000"
        :rows="4"
        placeholder="分享你的想法..."
        show-count
        data-testid="post-content"
      />
    </a-modal>

    <!-- Mobile FAB -->
    <div v-if="isMobile" class="forum-fab" @click="showCreateSheet = true">
      <span class="fab-icon">+</span>
    </div>

    <a-drawer
      v-model:open="showCreateSheet"
      placement="bottom"
      height="auto"
      title="发布内容"
    >
      <div class="create-sheet-options">
        <div class="sheet-option" @click="openCreatePost(); showCreateSheet = false">
          <span class="sheet-option-icon">💬</span>
          <span class="sheet-option-label">发动态</span>
        </div>
        <div class="sheet-option" @click="openCreateTopic(); showCreateSheet = false">
          <span class="sheet-option-icon">📝</span>
          <span class="sheet-option-label">发话题</span>
        </div>
      </div>
    </a-drawer>

    <!-- Create/Edit Topic Dialog -->
    <a-modal
      v-model:open="topicDialogOpen"
      :title="editingTopicItem ? '编辑话题' : '发话题'"
      :confirm-loading="topicSubmitting"
      width="520px"
      @ok="handleTopicSubmit"
      @cancel="topicDialogOpen = false"
      data-testid="topic-modal"
    >
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input
            v-model:value="topicForm.title"
            :maxlength="100"
            placeholder="请输入话题标题"
            data-testid="topic-title"
          />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea
            v-model:value="topicForm.content"
            :maxlength="2000"
            :rows="4"
            placeholder="可选，补充话题详情"
            show-count
            data-testid="topic-content"
          />
        </a-form-item>
        <a-form-item label="标签">
          <a-select
            v-model:value="topicForm.tag_id"
            placeholder="选择标签（可选）"
            allow-clear
            style="width: 100%"
            data-testid="topic-tag"
          >
            <a-select-option
              v-for="tag in tags"
              :key="tag.id"
              :value="tag.id"
            >
              {{ tag.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'
import {
  getFeed,
  createPost,
  updatePost,
  deletePost,
  createTopic,
  updateTopic,
  deleteTopic,
  togglePin as togglePinApi,
  getTags,
} from '@/api/forum'
import EmptyState from '@/components/EmptyState.vue'

interface MemberInfo {
  id: number
  name: string
  avatar: string
}

interface TagInfo {
  id: number
  name: string
  preset: boolean
}

interface FeedItem {
  type: 'post' | 'topic'
  id: number
  title: string
  content: string
  creator: MemberInfo
  tag: TagInfo | null
  is_pinned: boolean
  is_liked?: boolean
  like_count?: number
  comment_count?: number
  created_at: string
}

const router = useRouter()

const isMobile = ref(window.innerWidth < 768)
const showCreateSheet = ref(false)

const likingItems = ref<Record<string, boolean>>({})
const loading = ref(false)
const feedItems = ref<FeedItem[]>([])
const pinnedItems = ref<FeedItem[]>([])
const total = ref(0)
const pageSize = 20
const page = ref(1)

const tags = ref<TagInfo[]>([])

// Post dialog
const postDialogOpen = ref(false)
const postSubmitting = ref(false)
const editingPostItem = ref<FeedItem | null>(null)
const postForm = reactive({
  content: '',
})

// Topic dialog
const topicDialogOpen = ref(false)
const topicSubmitting = ref(false)
const editingTopicItem = ref<FeedItem | null>(null)
const topicForm = reactive({
  title: '',
  content: '',
  tag_id: undefined as number | undefined,
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

const isAdmin = computed(() => {
  const token = localStorage.getItem('token')
  if (!token) return false
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.role === 'admin'
  } catch {
    return false
  }
})

function truncate(text: string, maxLen: number): string {
  if (!text) return ''
  if (text.length <= maxLen) return text
  return text.slice(0, maxLen) + '...'
}

function timeAgo(date: string): string {
  const d = dayjs(date)
  const now = dayjs()
  const diffMins = now.diff(d, 'minute')
  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  const diffHours = now.diff(d, 'hour')
  if (diffHours < 24) return `${diffHours}小时前`
  const diffDays = now.diff(d, 'day')
  if (diffDays < 7) return `${diffDays}天前`
  return d.format('M月D日')
}

function goToDetail(item: FeedItem) {
  if (item.type === 'topic') {
    router.push(`/forum/topic/${item.id}`)
  }
}

function canManage(item: FeedItem): boolean {
  if (isAdmin.value) return true
  return item.creator.id === currentUserId.value
}

function handleFeedAction(key: string, item: FeedItem) {
  switch (key) {
    case 'edit':
      if (item.type === 'topic') openEditTopic(item)
      else openEditPost(item)
      break
    case 'delete':
      if (item.type === 'topic') confirmDeleteTopic(item)
      else confirmDeletePost(item)
      break
    case 'pin':
      handleTogglePin(item)
      break
  }
}

function onMenuClick(e: { key: string }, item: FeedItem) {
  handleFeedAction(e.key, item)
}

function handleLike(item: FeedItem) {
  const key = `${item.type}_${item.id}`
  likingItems.value[key] = true
  setTimeout(() => {
    likingItems.value[key] = false
  }, 150)
  // V1: like functionality is read-only display, no toggle API yet
  message.info('点赞功能将在后续版本上线')
}

// --- Fetch ---

async function fetchFeed() {
  loading.value = true
  try {
    const res: any = await getFeed({ page: page.value, page_size: pageSize })
    const data = res.data
    pinnedItems.value = data.pinned || []
    feedItems.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

async function fetchTags() {
  try {
    const res: any = await getTags()
    tags.value = res.data || []
  } catch {
    // tags are optional
  }
}

// --- Post actions ---

function openCreatePost() {
  editingPostItem.value = null
  postForm.content = ''
  postDialogOpen.value = true
}

function openEditPost(item: FeedItem) {
  editingPostItem.value = item
  postForm.content = item.content
  postDialogOpen.value = true
}

async function handlePostSubmit() {
  if (!postForm.content.trim()) {
    message.error('❌ 内容不能为空')
    return
  }
  postSubmitting.value = true
  try {
    if (editingPostItem.value) {
      await updatePost(editingPostItem.value.id, { content: postForm.content.trim() })
      message.success('✅ 更新成功')
    } else {
      await createPost({ content: postForm.content.trim() })
      message.success('✅ 发布成功')
    }
    postDialogOpen.value = false
    fetchFeed()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  } finally {
    postSubmitting.value = false
  }
}

function confirmDeletePost(item: FeedItem) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: '确定要删除这条动态吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deletePost(item.id)
        message.success('✅ 删除成功')
        fetchFeed()
      } catch (e: any) {
        if (e?.response?.data?.message) message.error(e.response.data.message)
        throw e
      }
    },
  })
}

// --- Topic actions ---

function openCreateTopic() {
  editingTopicItem.value = null
  topicForm.title = ''
  topicForm.content = ''
  topicForm.tag_id = undefined
  topicDialogOpen.value = true
}

function openEditTopic(item: FeedItem) {
  editingTopicItem.value = item
  topicForm.title = item.title
  topicForm.content = item.content || ''
  topicForm.tag_id = item.tag ? item.tag.id : undefined
  topicDialogOpen.value = true
}

async function handleTopicSubmit() {
  if (!topicForm.title.trim()) {
    message.error('❌ 标题不能为空')
    return
  }
  topicSubmitting.value = true
  try {
    const payload: any = {
      title: topicForm.title.trim(),
      content: topicForm.content.trim(),
    }
    if (topicForm.tag_id != null) {
      payload.tag_id = topicForm.tag_id
    }
    if (editingTopicItem.value) {
      await updateTopic(editingTopicItem.value.id, payload)
      message.success('✅ 更新成功')
    } else {
      await createTopic(payload)
      message.success('✅ 发布成功')
    }
    topicDialogOpen.value = false
    fetchFeed()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  } finally {
    topicSubmitting.value = false
  }
}

function confirmDeleteTopic(item: FeedItem) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: `确定要删除话题「${item.title}」吗？删除后话题下的所有评论也会一并删除。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteTopic(item.id)
        message.success('✅ 删除成功')
        fetchFeed()
      } catch (e: any) {
        if (e?.response?.data?.message) message.error(e.response.data.message)
        throw e
      }
    },
  })
}

async function handleTogglePin(item: FeedItem) {
  try {
    await togglePinApi(item.id)
    message.success(item.is_pinned ? '✅ 已取消置顶' : '✅ 已置顶')
    fetchFeed()
  } catch (e: any) {
    if (e?.response?.data?.message) message.error(e.response.data.message)
  }
}

onMounted(() => {
  fetchFeed()
  fetchTags()
})
</script>

<style scoped>
.forum-page {
  padding: 24px;
  max-width: 700px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
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

/* ==================== Pinned Section ==================== */
.pinned-section {
  background: var(--color-warning-light);
  border-radius: var(--radius-md);
  padding: 12px;
  margin-bottom: 16px;
}

/* ==================== Feed Card ==================== */
.feed-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feed-card {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 16px;
  margin-bottom: 16px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow var(--duration-normal) ease;
}

.feed-card:hover {
  box-shadow: var(--shadow-level-2);
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.feed-author {
  display: flex;
  align-items: center;
  gap: 6px;
}

.feed-avatar {
  font-size: 18px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.feed-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.feed-time {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.feed-body {
  margin-bottom: 12px;
}

.feed-content {
  font-size: 14px;
  line-height: 1.6;
  color: var(--color-text-primary);
  margin: 0;
  word-break: break-word;
}

.topic-pin-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  background: var(--color-brand-light);
  color: var(--color-brand);
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
}

.topic-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.topic-title:hover {
  color: var(--color-brand);
}

.topic-excerpt {
  color: var(--color-text-secondary);
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.topic-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.topic-tag {
  font-size: 12px;
  color: var(--color-text-secondary);
  padding: 2px 8px;
  background: var(--color-border-secondary);
  border-radius: var(--radius-sm);
}

.topic-author {
  margin-left: auto;
}

/* ==================== Feed Actions ==================== */
.feed-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 8px;
  border-top: 1px solid var(--color-border-secondary);
}

.feed-action {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: var(--color-text-secondary);
  cursor: pointer;
  min-height: 44px;
  min-width: 44px;
  transition: color var(--duration-fast) ease;
  user-select: none;
}

.feed-action:hover {
  color: var(--color-brand);
}

/* ==================== Forum FAB ==================== */
.forum-fab {
  position: fixed;
  bottom: calc(56px + 16px + env(safe-area-inset-bottom, 0px));
  right: 20px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--color-brand);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-level-3);
  cursor: pointer;
  z-index: var(--z-overlay);
  transition: transform var(--duration-fast) ease, box-shadow var(--duration-fast) ease;
}

.forum-fab:hover {
  transform: scale(1.05);
  box-shadow: var(--shadow-level-4);
}

.forum-fab:active {
  transform: scale(0.95);
}

.fab-icon {
  font-size: 28px;
  line-height: 1;
}

.create-sheet-options {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-bottom: 16px;
}

.sheet-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-md);
  cursor: pointer;
  min-height: 44px;
  transition: background var(--duration-fast) ease;
}

.sheet-option:hover {
  background: var(--color-brand-light);
}

.sheet-option-icon {
  font-size: 24px;
}

.sheet-option-label {
  font-size: 16px;
  font-weight: 500;
}

/* ==================== Pagination ==================== */
.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* ==================== Mobile ==================== */
@media (max-width: 767px) {
  .forum-page {
    padding: 16px;
  }

  .header-actions {
    gap: 6px;
  }

  .header-actions .ant-btn {
    padding: 0 12px;
    font-size: 13px;
  }

  .feed-card {
    padding: 14px;
  }

  .topic-title {
    font-size: 15px;
  }
}
</style>
