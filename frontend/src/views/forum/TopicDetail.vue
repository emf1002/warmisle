<template>
  <div class="topic-detail-page">
    <!-- Back button -->
    <div class="back-row">
      <a-button type="text" @click="$router.push('/forum')">
        ← 返回论坛
      </a-button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <!-- Topic not found -->
    <div v-else-if="!topic" class="empty-state">
      <a-result status="404" title="话题不存在" sub-title="该话题可能已被删除">
        <template #extra>
          <a-button type="primary" @click="$router.push('/forum')">返回论坛</a-button>
        </template>
      </a-result>
    </div>

    <template v-else>
      <!-- Topic detail card -->
      <div class="topic-card">
        <div class="topic-header">
          <a-tag color="blue" size="small">话题</a-tag>
          <a-tag v-if="topic.tag" size="small">{{ topic.tag.name }}</a-tag>
          <a-tag v-if="topic.is_pinned" color="orange" size="small">📌 置顶</a-tag>
        </div>

        <h1 class="topic-title">{{ topic.title }}</h1>

        <div v-if="topic.content" class="topic-content">{{ topic.content }}</div>

        <div class="topic-meta">
          <span class="topic-creator">{{ topic.creator.avatar }} {{ topic.creator.name }}</span>
          <span class="topic-time">{{ formatTime(topic.created_at) }}</span>
          <span v-if="topic.updated_at !== topic.created_at" class="topic-edited">（已编辑）</span>
        </div>

        <div class="topic-actions">
          <a-button
            type="text"
            :class="{ liked: topic.liked }"
            @click="handleLike"
          >
            {{ topic.liked ? '❤️' : '🤍' }} {{ topic.like_count }}
          </a-button>
          <span class="action-stat">💬 {{ topic.comment_count }}</span>

          <a-dropdown v-if="canManage" :trigger="['click']">
            <a-button type="text" size="small" class="more-btn">···</a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item @click="openEditDialog">编辑</a-menu-item>
                <a-menu-item v-if="authStore.isAdmin" @click="handleTogglePin">
                  {{ topic.is_pinned ? '取消置顶' : '置顶' }}
                </a-menu-item>
                <a-menu-item danger @click="confirmDelete">删除话题</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>

      <!-- Comments section -->
      <div class="comments-section">
        <h3 class="comments-title">评论 ({{ comments.length }})</h3>

        <!-- Comment input -->
        <div class="comment-input-box">
          <a-textarea
            v-model:value="commentText"
            :maxlength="500"
            :rows="3"
            :placeholder="replyTarget ? `回复 ${replyTarget.creator.name}：` : '写下你的评论...'"
            show-count
          />
          <div class="comment-input-actions">
            <a-button v-if="replyTarget" size="small" @click="cancelReply">取消回复</a-button>
            <a-button
              type="primary"
              size="small"
              :loading="commentSubmitting"
              :disabled="!commentText.trim()"
              @click="submitComment"
            >
              发布
            </a-button>
          </div>
        </div>

        <!-- Comments list -->
        <div v-if="comments.length === 0" class="no-comments">
          <a-empty description="暂无评论，来说点什么吧" />
        </div>

        <div v-else class="comments-list">
          <div
            v-for="comment in comments"
            :key="comment.id"
            class="comment-item"
          >
            <!-- Top-level comment -->
            <div class="comment-body">
              <span class="comment-avatar">{{ comment.creator.avatar }}</span>
              <div class="comment-main">
                <div class="comment-header">
                  <span class="comment-author">{{ comment.creator.name }}</span>
                  <span class="comment-time">{{ timeAgo(comment.created_at) }}</span>
                </div>
                <div class="comment-content">{{ comment.content }}</div>
                <div class="comment-actions">
                  <a-button type="link" size="small" @click="startReply(comment)">
                    回复
                  </a-button>
                  <a-button
                    v-if="canDeleteComment(comment)"
                    type="link"
                    size="small"
                    danger
                    @click="confirmDeleteComment(comment)"
                  >
                    删除
                  </a-button>
                </div>

                <!-- Reply input (inline) -->
                <div v-if="replyTarget?.id === comment.id" class="reply-input-box">
                  <a-textarea
                    v-model:value="replyText"
                    :maxlength="500"
                    :rows="2"
                    :placeholder="`回复 ${comment.creator.name}：`"
                  />
                  <div class="reply-input-actions">
                    <a-button size="small" @click="cancelReply">取消</a-button>
                    <a-button
                      type="primary"
                      size="small"
                      :disabled="!replyText.trim()"
                      @click="submitReply(comment)"
                    >
                      回复
                    </a-button>
                  </div>
                </div>

                <!-- Children (replies) -->
                <div v-if="comment.children && comment.children.length > 0" class="replies-list">
                  <div
                    v-for="child in comment.children"
                    :key="child.id"
                    class="reply-item"
                  >
                    <span class="reply-avatar">{{ child.creator.avatar }}</span>
                    <div class="reply-main">
                      <div class="reply-header">
                        <span class="reply-author">{{ child.creator.name }}</span>
                        <span class="reply-time">{{ timeAgo(child.created_at) }}</span>
                      </div>
                      <div class="reply-content">{{ child.content }}</div>
                      <div class="reply-actions">
                        <a-button
                          v-if="canDeleteComment(child)"
                          type="link"
                          size="small"
                          danger
                          @click="confirmDeleteComment(child)"
                        >
                          删除
                        </a-button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Edit Topic Dialog -->
    <a-modal
      v-model:open="editDialogOpen"
      title="编辑话题"
      ok-text="保存"
      cancel-text="取消"
      :confirm-loading="editSubmitting"
      width="520px"
      @ok="handleEditSubmit"
      @cancel="editDialogOpen = false"
    >
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input
            v-model:value="editForm.title"
            :maxlength="100"
            placeholder="请输入话题标题"
          />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea
            v-model:value="editForm.content"
            :maxlength="2000"
            :rows="4"
            placeholder="可选，补充话题详情"
            show-count
          />
        </a-form-item>
        <a-form-item label="标签">
          <a-select
            v-model:value="editForm.tag_id"
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
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { timeAgo } from '@/utils/format'
import dayjs from 'dayjs'
import {
  getTopic,
  updateTopic,
  deleteTopic,
  togglePin as togglePinApi,
  getComments,
  createComment,
  deleteComment,
  toggleLike,
  getTags,
} from '@/api/forum'
import { useAuthStore } from '@/stores/auth'

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

interface TopicDetail {
  id: number
  title: string
  content: string
  tag_id: number | null
  creator_id: number
  is_pinned: boolean
  created_at: string
  updated_at: string
  creator: MemberInfo
  tag: TagInfo | null
  like_count: number
  comment_count: number
  liked: boolean
}

interface CommentItem {
  id: number
  target_type: string
  target_id: number
  parent_id: number | null
  content: string
  creator_id: number
  created_at: string
  creator: MemberInfo
  children: CommentItem[]
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const topicId = computed(() => Number(route.params.id))

const loading = ref(false)
const topic = ref<TopicDetail | null>(null)
const comments = ref<CommentItem[]>([])
const tags = ref<TagInfo[]>([])

// Comment input
const commentText = ref('')
const commentSubmitting = ref(false)
const replyTarget = ref<CommentItem | null>(null)
const replyText = ref('')

// Edit dialog
const editDialogOpen = ref(false)
const editSubmitting = ref(false)
const editForm = reactive({
  title: '',
  content: '',
  tag_id: undefined as number | undefined,
})

const canManage = computed(() => {
  if (!topic.value) return false
  if (authStore.isAdmin) return true
  return topic.value.creator_id === authStore.currentUserId
})

function canDeleteComment(comment: CommentItem): boolean {
  if (authStore.isAdmin) return true
  return comment.creator_id === authStore.currentUserId
}

function formatTime(date: string): string {
  const d = dayjs(date)
  const now = dayjs()
  if (d.isSame(now, 'day')) return `今天 ${d.format('HH:mm')}`
  if (d.isSame(now.subtract(1, 'day'), 'day')) return `昨天 ${d.format('HH:mm')}`
  return d.format('M月D日 HH:mm')
}

// --- Fetch ---

async function fetchTopic() {
  loading.value = true
  try {
    const res: any = await getTopic(topicId.value)
    topic.value = res.data
  } catch {
    topic.value = null
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  try {
    const res: any = await getComments({
      target_type: 'topic',
      target_id: topicId.value,
    })
    comments.value = res.data || []
  } catch {
    comments.value = []
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

// --- Like ---

async function handleLike() {
  if (!topic.value) return
  try {
    const res: any = await toggleLike({
      target_type: 'topic',
      target_id: topic.value.id,
    })
    const liked = res.data.liked
    topic.value.liked = liked
    topic.value.like_count += liked ? 1 : -1
  } catch {
    // error handled by interceptor
  }
}

// --- Comments ---

function startReply(comment: CommentItem) {
  replyTarget.value = comment
  replyText.value = ''
}

function cancelReply() {
  replyTarget.value = null
  replyText.value = ''
  commentText.value = ''
}

async function submitComment() {
  if (!commentText.value.trim()) return
  commentSubmitting.value = true
  try {
    await createComment({
      target_type: 'topic',
      target_id: topicId.value,
      content: commentText.value.trim(),
    })
    message.success('✅ 评论已发布')
    commentText.value = ''
    fetchComments()
    // Refresh topic to update comment count
    fetchTopic()
  } catch {
    // error handled by interceptor
  } finally {
    commentSubmitting.value = false
  }
}

async function submitReply(parentComment: CommentItem) {
  if (!replyText.value.trim()) return
  try {
    await createComment({
      target_type: 'topic',
      target_id: topicId.value,
      parent_id: parentComment.id,
      content: replyText.value.trim(),
    })
    message.success('✅ 回复已发布')
    cancelReply()
    fetchComments()
    fetchTopic()
  } catch {
    // error handled by interceptor
  }
}

function confirmDeleteComment(comment: CommentItem) {
  Modal.confirm({
    title: '❓ 确认删除',
    content: comment.children && comment.children.length > 0
      ? '删除该评论后，其下的回复也会一并删除，确定继续吗？'
      : '确定要删除这条评论吗？',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteComment(comment.id)
        message.success('✅ 删除成功')
        fetchComments()
        fetchTopic()
      } catch (e: any) {
        throw e // prevent modal from closing on error
      }
    },
  })
}

// --- Edit topic ---

function openEditDialog() {
  if (!topic.value) return
  editForm.title = topic.value.title
  editForm.content = topic.value.content || ''
  editForm.tag_id = topic.value.tag ? topic.value.tag.id : undefined
  editDialogOpen.value = true
}

async function handleEditSubmit() {
  if (!editForm.title.trim()) {
    message.error('❌ 标题不能为空')
    return
  }
  editSubmitting.value = true
  try {
    const payload: any = {
      title: editForm.title.trim(),
      content: editForm.content.trim(),
    }
    if (editForm.tag_id != null) {
      payload.tag_id = editForm.tag_id
    }
    await updateTopic(topicId.value, payload)
    message.success('✅ 更新成功')
    editDialogOpen.value = false
    fetchTopic()
  } catch {
    // error handled by interceptor
  } finally {
    editSubmitting.value = false
  }
}

// --- Delete topic ---

function confirmDelete() {
  Modal.confirm({
    title: '❓ 确认删除',
    content: '确定要删除这个话题吗？删除后所有评论也会一并删除。',
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteTopic(topicId.value)
        message.success('✅ 删除成功')
        router.push('/forum')
      } catch (e: any) {
        throw e // prevent modal from closing on error
      }
    },
  })
}

// --- Toggle pin ---

async function handleTogglePin() {
  if (!topic.value) return
  try {
    await togglePinApi(topic.value.id)
    message.success(topic.value.is_pinned ? '✅ 已取消置顶' : '✅ 已置顶')
    fetchTopic()
  } catch {
    // error handled by interceptor
  }
}

onMounted(() => {
  fetchTopic()
  fetchComments()
  fetchTags()
})
</script>

<style scoped>
.topic-detail-page {
  padding: 24px;
  max-width: 700px;
  margin: 0 auto;
}

.back-row {
  margin-bottom: 16px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
  color: var(--color-text-secondary);
}

.empty-state {
  padding: 48px 0;
}

/* Topic card */
.topic-card {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 20px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  margin-bottom: 24px;
}

.topic-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.topic-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 12px 0;
  word-break: break-word;
  line-height: 1.4;
}

.topic-content {
  font-size: 15px;
  color: var(--color-text-primary);
  line-height: 1.8;
  word-break: break-word;
  white-space: pre-wrap;
  margin-bottom: 16px;
}

.topic-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.topic-creator {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.topic-time {
  font-size: 12px;
  color: var(--color-text-disabled);
}

.topic-edited {
  font-size: 12px;
  color: var(--color-text-disabled);
}

.topic-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--color-bg-layout);
}

.topic-actions .ant-btn-text {
  min-height: 44px;
}

.topic-actions .ant-btn-text.liked {
  color: var(--color-danger);
}

.action-stat {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.more-btn {
  color: var(--color-text-secondary);
  min-width: 32px;
}

/* Comments section */
.comments-section {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 20px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
}

.comments-title {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 16px 0;
}

.comment-input-box {
  margin-bottom: 20px;
}

.comment-input-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.no-comments {
  padding: 24px 0;
}

/* Comment items */
.comments-list {
  display: flex;
  flex-direction: column;
}

.comment-item {
  border-bottom: 1px solid var(--color-bg-layout);
  padding-bottom: 12px;
  margin-bottom: 12px;
}

.comment-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.comment-body {
  display: flex;
  gap: 10px;
}

.comment-avatar,
.reply-avatar {
  font-size: 20px;
  flex-shrink: 0;
  width: 36px;
  text-align: center;
}

.comment-main {
  flex: 1;
  min-width: 0;
}

.comment-header,
.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-author,
.reply-author {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.comment-time,
.reply-time {
  font-size: 12px;
  color: var(--color-text-disabled);
}

.comment-content,
.reply-content {
  font-size: 14px;
  color: var(--color-text-primary);
  line-height: 1.6;
  word-break: break-word;
  white-space: pre-wrap;
  margin-bottom: 4px;
}

.comment-actions,
.reply-actions {
  display: flex;
  gap: 4px;
  min-height: 32px;
}

.comment-actions .ant-btn-link,
.reply-actions .ant-btn-link {
  padding: 0 8px;
  font-size: 12px;
  min-height: 28px;
}

/* Reply input */
.reply-input-box {
  margin-top: 8px;
  margin-bottom: 8px;
}

.reply-input-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 6px;
}

/* Replies list */
.replies-list {
  margin-top: 10px;
  padding: 10px 12px;
  background: var(--color-border-secondary);
  border-radius: var(--radius-md);
}

.reply-item {
  display: flex;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--color-border-secondary);
}

.reply-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.reply-main {
  flex: 1;
  min-width: 0;
}

@media (max-width: 767px) {
  .topic-detail-page {
    padding: 16px;
  }

  .topic-card {
    padding: 16px;
  }

  .topic-title {
    font-size: 18px;
  }

  .topic-content {
    font-size: 14px;
  }

  .comments-section {
    padding: 16px;
  }

  .comment-body {
    gap: 8px;
  }
}
</style>
