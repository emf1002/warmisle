<template>
  <div class="forum-page">
    <div class="page-header">
      <h2>家庭论坛</h2>
      <div class="header-actions">
        <a-button @click="openCreatePost">发动态</a-button>
        <a-button type="primary" @click="openCreateTopic">发话题</a-button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <template v-else-if="feedItems.length === 0 && pinnedItems.length === 0">
      <div class="empty-state">
        <a-empty description="暂无动态">
          <template #children>
            <a-button type="primary" @click="openCreatePost">发第一条动态</a-button>
          </template>
        </a-empty>
      </div>
    </template>

    <template v-else>
      <!-- Pinned topics -->
      <div v-if="pinnedItems.length > 0" class="pinned-section">
        <div v-for="item in pinnedItems" :key="'pinned-' + item.id" class="feed-card pinned-card">
          <div class="card-pin-badge">📌 置顶</div>
          <div class="card-header" @click="goTopic(item.id)">
            <span class="card-type-tag">
              <a-tag color="blue" size="small">话题</a-tag>
            </span>
            <a-tag v-if="item.tag" size="small" class="card-topic-tag">{{ item.tag.name }}</a-tag>
          </div>
          <h3 class="card-title" @click="goTopic(item.id)">{{ item.title }}</h3>
          <p v-if="item.content" class="card-content" @click="goTopic(item.id)">
            {{ truncate(item.content, 150) }}
          </p>
          <div class="card-footer">
            <span class="card-author">{{ item.creator.avatar }} {{ item.creator.name }}</span>
            <span class="card-time">{{ timeAgo(item.created_at) }}</span>
            <a-dropdown v-if="canEditTopic(item)" :trigger="['click']">
              <a-button type="text" size="small" class="card-more-btn">···</a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="openEditTopic(item)">编辑</a-menu-item>
                  <a-menu-item v-if="isAdmin" @click="handleTogglePin(item)">
                    {{ item.is_pinned ? '取消置顶' : '置顶' }}
                  </a-menu-item>
                  <a-menu-item danger @click="confirmDeleteTopic(item)">删除</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>

      <!-- Feed items -->
      <div class="feed-list">
        <div v-for="item in feedItems" :key="item.type + '-' + item.id" class="feed-card">
          <!-- Topic card -->
          <template v-if="item.type === 'topic'">
            <div class="card-header" @click="goTopic(item.id)">
              <span class="card-type-tag">
                <a-tag color="blue" size="small">话题</a-tag>
              </span>
              <a-tag v-if="item.tag" size="small" class="card-topic-tag">{{ item.tag.name }}</a-tag>
            </div>
            <h3 class="card-title" @click="goTopic(item.id)">{{ item.title }}</h3>
            <p v-if="item.content" class="card-content" @click="goTopic(item.id)">
              {{ truncate(item.content, 150) }}
            </p>
            <div class="card-footer">
              <span class="card-author">{{ item.creator.avatar }} {{ item.creator.name }}</span>
              <span class="card-time">{{ timeAgo(item.created_at) }}</span>
              <a-dropdown v-if="canEditTopic(item)" :trigger="['click']">
                <a-button type="text" size="small" class="card-more-btn">···</a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="openEditTopic(item)">编辑</a-menu-item>
                    <a-menu-item v-if="isAdmin" @click="handleTogglePin(item)">
                      {{ item.is_pinned ? '取消置顶' : '置顶' }}
                    </a-menu-item>
                    <a-menu-item danger @click="confirmDeleteTopic(item)">删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </template>

          <!-- Post card -->
          <template v-else>
            <div class="card-header">
              <span class="card-type-tag">
                <a-tag size="small">动态</a-tag>
              </span>
            </div>
            <p class="card-content card-content-post">
              {{ truncate(item.content, 200) }}
            </p>
            <div class="card-footer">
              <span class="card-author">{{ item.creator.avatar }} {{ item.creator.name }}</span>
              <span class="card-time">{{ timeAgo(item.created_at) }}</span>
              <a-dropdown v-if="canEditPost(item)" :trigger="['click']">
                <a-button type="text" size="small" class="card-more-btn">···</a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="openEditPost(item)">编辑</a-menu-item>
                    <a-menu-item danger @click="confirmDeletePost(item)">删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </template>
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
    >
      <a-textarea
        v-model:value="postForm.content"
        :maxlength="1000"
        :rows="4"
        placeholder="分享你的想法..."
        show-count
      />
    </a-modal>

    <!-- Create/Edit Topic Dialog -->
    <a-modal
      v-model:open="topicDialogOpen"
      :title="editingTopicItem ? '编辑话题' : '发话题'"
      :confirm-loading="topicSubmitting"
      width="520px"
      @ok="handleTopicSubmit"
      @cancel="topicDialogOpen = false"
    >
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input
            v-model:value="topicForm.title"
            :maxlength="100"
            placeholder="请输入话题标题"
          />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea
            v-model:value="topicForm.content"
            :maxlength="2000"
            :rows="4"
            placeholder="可选，补充话题详情"
            show-count
          />
        </a-form-item>
        <a-form-item label="标签">
          <a-select
            v-model:value="topicForm.tag_id"
            placeholder="选择标签（可选）"
            allow-clear
            style="width: 100%"
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
  created_at: string
}

const router = useRouter()

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

function canEditPost(item: FeedItem): boolean {
  if (isAdmin.value) return true
  return item.creator.id === currentUserId.value
}

function canEditTopic(item: FeedItem): boolean {
  if (isAdmin.value) return true
  return item.creator.id === currentUserId.value
}

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

function goTopic(id: number) {
  router.push(`/forum/topic/${id}`)
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
    message.error('内容不能为空')
    return
  }
  postSubmitting.value = true
  try {
    if (editingPostItem.value) {
      await updatePost(editingPostItem.value.id, { content: postForm.content.trim() })
      message.success('更新成功')
    } else {
      await createPost({ content: postForm.content.trim() })
      message.success('发布成功')
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
    title: '确认删除',
    content: '确定要删除这条动态吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deletePost(item.id)
        message.success('删除成功')
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
    message.error('标题不能为空')
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
      message.success('更新成功')
    } else {
      await createTopic(payload)
      message.success('发布成功')
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
    title: '确认删除',
    content: `确定要删除话题「${item.title}」吗？删除后话题下的所有评论也会一并删除。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteTopic(item.id)
        message.success('删除成功')
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
    message.success(item.is_pinned ? '已取消置顶' : '已置顶')
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
  color: #999;
}

.empty-state {
  text-align: center;
  padding: 48px 0;
}

/* Pinned section */
.pinned-section {
  margin-bottom: 16px;
}

.pinned-card {
  background: #fffbe6;
  border-color: #ffe58f;
}

.card-pin-badge {
  font-size: 12px;
  color: #ad8b00;
  margin-bottom: 8px;
}

/* Feed cards */
.feed-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feed-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  border: 1px solid #f0f0f0;
  transition: box-shadow 0.2s;
}

.feed-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.card-type-tag {
  flex-shrink: 0;
}

.card-topic-tag {
  flex-shrink: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 6px 0;
  cursor: pointer;
  word-break: break-word;
  line-height: 1.5;
}

.card-title:hover {
  color: #1677ff;
}

.card-content {
  font-size: 14px;
  color: #555;
  margin: 0 0 10px 0;
  line-height: 1.6;
  word-break: break-word;
  cursor: pointer;
}

.card-content-post {
  cursor: default;
}

.card-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 32px;
}

.card-author {
  font-size: 12px;
  color: #999;
}

.card-time {
  font-size: 12px;
  color: #bbb;
  flex: 1;
}

.card-more-btn {
  color: #999;
  min-width: 32px;
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Pagination */
.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

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

  .card-title {
    font-size: 15px;
  }
}
</style>
