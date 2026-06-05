<template>
  <div class="forum-page" data-testid="forum-page">
    <div class="page-header">
      <h2>家庭论坛</h2>
      <div class="header-actions">
        <a-button @click="openCreatePost" data-testid="create-post-btn">发动态</a-button>
        <a-button type="primary" @click="openCreateTopic" data-testid="create-topic-btn">发话题</a-button>
        <a-button v-if="authStore.isAdmin" @click="openCreateAnnouncement" data-testid="create-announcement-btn">发布公告</a-button>
        <a-button @click="openCreatePoll" data-testid="create-poll-btn">发起投票</a-button>
        <a-button v-if="authStore.isAdmin" @click="openManageTags" data-testid="manage-tags-btn">管理标签</a-button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <a-spin />
      <span style="margin-left: 8px">加载中...</span>
    </div>

    <template v-else-if="displayItems.length === 0 && pinnedItems.length === 0">
      <EmptyState type="no-data" description="暂无动态">
        <template #action>
          <a-button type="primary" @click="openCreatePost">发第一条动态</a-button>
        </template>
      </EmptyState>
    </template>

    <template v-else>
      <!-- Pinned topics -->
      <div v-if="pinnedItems.length > 0" class="pinned-section">
        <div v-for="item in pinnedItems" :key="'pinned-' + item.id" class="feed-card topic-card" :data-testid="'feed-card-' + item.id">
          <div class="feed-header">
            <span class="topic-pin-badge" data-testid="pinned-tag">📌 公告</span>
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
                  <a-menu-item v-if="authStore.isAdmin" key="pin" data-testid="unpin-btn">
                    {{ item.is_pinned ? '取消置顶' : '置顶' }}
                  </a-menu-item>
                  <a-menu-item key="edit">编辑</a-menu-item>
                  <a-menu-item key="delete" danger data-testid="delete-feed-btn">删除</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </div>

      <!-- Feed items -->
      <div class="feed-list" data-testid="feed-list">
        <div v-for="(item, index) in displayItems" :key="item.type + '-' + item.id">
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
                    <a-menu-item key="delete" danger data-testid="delete-feed-btn">删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>

          <!-- Topic card -->
          <div v-else-if="item.type === 'topic'" class="feed-card topic-card card-stagger" :style="{ animationDelay: `${index * 50}ms` }" :data-testid="'feed-card-' + item.id">
            <div v-if="item.is_pinned" class="feed-header">
              <span class="topic-pin-badge" data-testid="pinned-tag">📌 公告</span>
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
                    <a-menu-item v-if="authStore.isAdmin" key="pin" data-testid="unpin-btn">
                      {{ item.is_pinned ? '取消置顶' : '置顶' }}
                    </a-menu-item>
                    <a-menu-item key="edit">编辑</a-menu-item>
                    <a-menu-item key="delete" danger data-testid="delete-feed-btn">删除</a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>

          <!-- Poll card -->
          <div v-else-if="item.type === 'vote'" class="feed-card poll-card card-stagger" :style="{ animationDelay: `${index * 50}ms` }" :data-testid="'feed-card-' + item.id">
            <div class="feed-header">
              <span class="feed-author">
                <span class="feed-avatar" :aria-label="`${item.creator?.name || '用户'}的头像`">{{ item.creator?.avatar || '👤' }}</span>
                <span class="feed-name">{{ item.creator?.name }}</span>
              </span>
              <span class="feed-time">{{ timeAgo(item.created_at) }}</span>
            </div>
            <h3 class="poll-title">{{ item.title }}</h3>
            <div v-if="item.is_multi !== undefined" class="poll-type-hint">
              {{ item.is_multi ? '多选' : '单选' }}
            </div>
            <template v-if="!item.voted">
              <div class="poll-options">
                <div
                  v-for="(opt, optIndex) in item.options"
                  :key="opt.id"
                  class="poll-option-item"
                  :class="{ selected: selectedPollOptions[item.id]?.includes(opt.id) }"
                  :data-testid="'poll-option-' + optIndex"
                  @click="selectPollOption(item, opt.id)"
                >
                  <span v-if="item.is_multi" class="poll-option-check">
                    <a-checkbox :checked="selectedPollOptions[item.id]?.includes(opt.id)" @click.stop="selectPollOption(item, opt.id)" />
                  </span>
                  <span v-else class="poll-option-radio">
                    <a-radio :checked="selectedPollOptions[item.id]?.includes(opt.id)" @click.stop="selectPollOption(item, opt.id)" />
                  </span>
                  <span class="poll-option-text">{{ opt.text }}</span>
                </div>
              </div>
              <a-button
                type="primary"
                size="small"
                :disabled="!selectedPollOptions[item.id]?.length"
                @click="submitPollVote(item)"
                data-testid="poll-submit"
                class="poll-submit-btn"
              >
                提交投票
              </a-button>
            </template>
            <template v-else>
              <div class="poll-results" data-testid="poll-result">
                <div v-for="opt in item.options" :key="opt.id" class="poll-result-item">
                  <div class="poll-result-header">
                    <span class="poll-result-text">{{ opt.text }}</span>
                    <span class="poll-result-count">{{ opt.vote_count || 0 }}票</span>
                  </div>
                  <a-progress
                    :percent="getPollPercent(item, opt)"
                    :show-info="false"
                    :stroke-color="item.user_voted_options?.includes(opt.id) ? 'var(--color-brand)' : 'var(--color-border)'"
                    size="small"
                  />
                </div>
              </div>
            </template>
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
      ok-text="发布"
      cancel-text="取消"
      :confirm-loading="postSubmitting"
      @ok="handlePostSubmit"
      @cancel="postDialogOpen = false"
    >
      <div data-testid="post-modal">
      <a-textarea
        v-model:value="postForm.content"
        :maxlength="1000"
        :rows="4"
        placeholder="分享你的想法..."
        show-count
        data-testid="post-content"
      />
      </div>
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
        <div v-if="authStore.isAdmin" class="sheet-option" @click="openCreateAnnouncement(); showCreateSheet = false">
          <span class="sheet-option-icon">📢</span>
          <span class="sheet-option-label">发布公告</span>
        </div>
        <div class="sheet-option" @click="openCreatePoll(); showCreateSheet = false">
          <span class="sheet-option-icon">📊</span>
          <span class="sheet-option-label">发起投票</span>
        </div>
      </div>
    </a-drawer>

    <!-- Create/Edit Topic Dialog -->
    <a-modal
      v-model:open="topicDialogOpen"
      :title="editingTopicItem ? '编辑话题' : (isAnnouncement ? '发布公告' : '发话题')"
      ok-text="发布"
      cancel-text="取消"
      :confirm-loading="topicSubmitting"
      width="520px"
      @ok="handleTopicSubmit"
      @cancel="topicDialogOpen = false"
    >
      <div data-testid="topic-modal">
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
      </div>
    </a-modal>

    <!-- Create Poll Dialog -->
    <a-modal
      v-model:open="pollDialogOpen"
      title="发起投票"
      ok-text="创建"
      cancel-text="取消"
      :confirm-loading="pollSubmitting"
      width="520px"
      @ok="handlePollSubmit"
      @cancel="pollDialogOpen = false"
    >
      <div data-testid="poll-modal">
      <a-form layout="vertical">
        <a-form-item label="投票标题" required>
          <a-input
            v-model:value="pollForm.title"
            :maxlength="100"
            placeholder="请输入投票标题"
            data-testid="poll-title"
          />
        </a-form-item>
        <a-form-item label="选项">
          <div class="poll-form-options">
            <div v-for="(_opt, idx) in pollForm.options" :key="idx" class="poll-form-option-row">
              <a-input
                v-model:value="pollForm.options[idx]"
                :placeholder="'选项 ' + (idx + 1)"
                :data-testid="'option-input-' + (idx + 1)"
              />
              <a-button
                v-if="pollForm.options.length > 2"
                type="text"
                danger
                size="small"
                @click="pollForm.options.splice(idx, 1)"
              >
                删除
              </a-button>
            </div>
          </div>
          <a-button
            v-if="pollForm.options.length < 10"
            type="dashed"
            block
            style="margin-top: 8px"
            @click="pollForm.options.push('')"
            data-testid="add-option-btn"
          >
            + 添加选项
          </a-button>
        </a-form-item>
        <a-form-item>
          <a-checkbox v-model:checked="pollForm.is_multi" data-testid="poll-multi-select">
            允许多选
          </a-checkbox>
        </a-form-item>
      </a-form>
      </div>
    </a-modal>

    <!-- Manage Tags Dialog -->
    <a-modal
      v-model:open="tagDialogOpen"
      title="管理标签"
      :footer="null"
      width="480px"
    >
      <div class="tag-manage-list">
        <div v-for="tag in tags" :key="tag.id" data-testid="tag-item" class="tag-manage-item">
          <span class="tag-manage-name">{{ tag.name }}</span>
          <a-button
            size="small"
            danger
            :disabled="usedTagIds.has(tag.id)"
            data-testid="delete-tag-btn"
            @click="handleDeleteTag(tag)"
          >
            删除
          </a-button>
        </div>
        <div v-if="tags.length === 0" class="tag-empty-hint">暂无标签</div>
      </div>
      <div class="tag-add-section">
        <a-button v-if="!showAddTagForm" type="dashed" block @click="showAddTagForm = true" data-testid="add-tag-btn">
          + 添加标签
        </a-button>
        <div v-else class="tag-add-form">
          <a-input
            v-model:value="newTagName"
            placeholder="输入标签名称"
            data-testid="tag-name-input"
            @keyup.enter="handleCreateTag"
          />
          <a-button type="primary" :disabled="!newTagName.trim()" data-testid="tag-submit-btn" @click="handleCreateTag">
            添加
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { truncate, timeAgo } from '@/utils/format'
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
  createTag,
  deleteTag,
  createVote,
  vote as submitVoteApi,
  getVote,
} from '@/api/forum'
import EmptyState from '@/components/EmptyState.vue'
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

interface PollOption {
  id: number
  text: string
  vote_count?: number
}

interface PollItem {
  type: 'vote'
  id: number
  title: string
  options: PollOption[]
  is_multi?: boolean
  creator: MemberInfo
  created_at: string
  voted?: boolean
  user_voted_options?: number[]
  total_votes?: number
}

type DisplayItem = FeedItem | PollItem

const router = useRouter()
const authStore = useAuthStore()

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
const isAnnouncement = ref(false)

// Tag management dialog
const tagDialogOpen = ref(false)
const showAddTagForm = ref(false)
const newTagName = ref('')

// Poll dialog
const pollDialogOpen = ref(false)
const pollSubmitting = ref(false)
const pollItems = ref<PollItem[]>([])
const selectedPollOptions = ref<Record<number, number[]>>({})
const pollForm = reactive({
  title: '',
  options: ['', ''] as string[],
  is_multi: false,
})

const displayItems = computed<DisplayItem[]>(() => {
  const items: DisplayItem[] = [...feedItems.value, ...pollItems.value]
  items.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  return items
})

const usedTagIds = computed(() => {
  const ids = new Set<number>()
  const allItems = [...pinnedItems.value, ...feedItems.value]
  for (const item of allItems) {
    if (item.tag) {
      ids.add(item.tag.id)
    }
  }
  return ids
})

function goToDetail(item: FeedItem) {
  if (item.type === 'topic') {
    router.push(`/forum/topic/${item.id}`)
  }
}

function canManage(item: FeedItem): boolean {
  if (authStore.isAdmin) return true
  return item.creator.id === authStore.currentUserId
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
  } catch {
    // error handled by interceptor
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
        throw e // prevent modal from closing on error
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
  isAnnouncement.value = false
  topicDialogOpen.value = true
}

function openCreateAnnouncement() {
  editingTopicItem.value = null
  topicForm.title = ''
  topicForm.content = ''
  isAnnouncement.value = true
  // Auto-select the "公告" tag if it exists
  const announcementTag = tags.value.find(t => t.name === '公告')
  topicForm.tag_id = announcementTag ? announcementTag.id : undefined
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
      const res: any = await createTopic(payload)
      // Auto-pin if creating an announcement
      if (isAnnouncement.value) {
        const newTopicId = res?.data?.id
        if (newTopicId) {
          await togglePinApi(newTopicId)
        }
      }
      message.success('✅ 发布成功')
    }
    isAnnouncement.value = false
    topicDialogOpen.value = false
    fetchFeed()
  } catch {
    // error handled by interceptor
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
        throw e // prevent modal from closing on error
      }
    },
  })
}

async function handleTogglePin(item: FeedItem) {
  try {
    await togglePinApi(item.id)
    message.success(item.is_pinned ? '✅ 已取消置顶' : '✅ 已置顶')
    fetchFeed()
  } catch {
    // error handled by interceptor
  }
}

// --- Tag management ---

function openManageTags() {
  showAddTagForm.value = false
  newTagName.value = ''
  tagDialogOpen.value = true
}

async function handleCreateTag() {
  const name = newTagName.value.trim()
  if (!name) return
  try {
    await createTag({ name })
    message.success('✅ 标签已创建')
    newTagName.value = ''
    showAddTagForm.value = false
    fetchTags()
  } catch {
    // error handled by interceptor
  }
}

function handleDeleteTag(tag: TagInfo) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除标签「${tag.name}」吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteTag(tag.id)
        message.success('✅ 标签已删除')
        fetchTags()
      } catch (e: any) {
        throw e // prevent modal from closing on error
      }
    },
  })
}

// --- Poll actions ---

function openCreatePoll() {
  pollForm.title = ''
  pollForm.options = ['', '']
  pollForm.is_multi = false
  pollDialogOpen.value = true
}

async function handlePollSubmit() {
  if (!pollForm.title.trim()) {
    message.error('投票标题不能为空')
    return
  }
  const validOptions = pollForm.options.map(o => o.trim()).filter(Boolean)
  if (validOptions.length < 2) {
    message.error('至少需要2个选项')
    return
  }
  pollSubmitting.value = true
  try {
    const res: any = await createVote({
      title: pollForm.title.trim(),
      options: validOptions,
      is_multi: pollForm.is_multi,
    })
    const voteData = res.data
    pollItems.value.unshift({
      type: 'vote',
      id: voteData.id,
      title: voteData.title,
      options: voteData.options || validOptions.map((text, i) => ({ id: i + 1, text, vote_count: 0 })),
      is_multi: voteData.is_multi ?? pollForm.is_multi,
      creator: voteData.creator || { id: authStore.currentUserId, name: authStore.memberInfo?.name || '我', avatar: authStore.memberInfo?.avatar || '👤' },
      created_at: voteData.created_at || new Date().toISOString(),
      voted: false,
    })
    pollDialogOpen.value = false
    message.success('投票已创建')
  } catch {
    // error handled by interceptor
  } finally {
    pollSubmitting.value = false
  }
}

function selectPollOption(poll: PollItem, optionId: number) {
  if (poll.voted) return
  if (!selectedPollOptions.value[poll.id]) {
    selectedPollOptions.value[poll.id] = []
  }
  const opts = selectedPollOptions.value[poll.id]
  const idx = opts.indexOf(optionId)
  if (poll.is_multi) {
    if (idx >= 0) {
      opts.splice(idx, 1)
    } else {
      opts.push(optionId)
    }
  } else {
    selectedPollOptions.value[poll.id] = [optionId]
  }
}

async function submitPollVote(poll: PollItem) {
  const selected = selectedPollOptions.value[poll.id]
  if (!selected || selected.length === 0) return
  try {
    for (const optionId of selected) {
      await submitVoteApi(poll.id, { option_id: optionId })
    }
    // Fetch updated vote data
    const res: any = await getVote(poll.id)
    const voteData = res.data
    const idx = pollItems.value.findIndex(p => p.id === poll.id)
    if (idx >= 0) {
      pollItems.value[idx] = {
        ...pollItems.value[idx],
        voted: true,
        options: voteData.options || pollItems.value[idx].options,
        user_voted_options: selected,
        total_votes: voteData.total_votes,
      }
    }
    message.success('投票成功')
  } catch {
    // If backend rejects duplicate vote, still show results
    try {
      const res: any = await getVote(poll.id)
      const voteData = res.data
      const idx = pollItems.value.findIndex(p => p.id === poll.id)
      if (idx >= 0) {
        pollItems.value[idx] = {
          ...pollItems.value[idx],
          voted: true,
          options: voteData.options || pollItems.value[idx].options,
          user_voted_options: selected,
          total_votes: voteData.total_votes,
        }
      }
    } catch {
      // ignore
    }
  }
}

function getPollPercent(poll: PollItem, opt: PollOption): number {
  const total = poll.options.reduce((sum, o) => sum + (o.vote_count || 0), 0)
  if (total === 0) return 0
  return Math.round(((opt.vote_count || 0) / total) * 100)
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

.header-actions {
  display: flex;
  gap: 8px;
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

/* ==================== Poll Card ==================== */
.poll-card {
  border-left: 3px solid var(--color-brand);
}

.poll-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 8px 0;
}

.poll-type-hint {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 12px;
}

.poll-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.poll-option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border: 1px solid var(--color-border-secondary);
  border-radius: var(--radius-md);
  cursor: pointer;
  min-height: 44px;
  transition: border-color var(--duration-fast) ease, background var(--duration-fast) ease;
}

.poll-option-item:hover {
  border-color: var(--color-brand);
  background: var(--color-brand-light);
}

.poll-option-item.selected {
  border-color: var(--color-brand);
  background: var(--color-brand-light);
}

.poll-option-text {
  font-size: 14px;
  color: var(--color-text-primary);
}

.poll-submit-btn {
  margin-top: 4px;
}

.poll-results {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.poll-result-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.poll-result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.poll-result-text {
  font-size: 14px;
  color: var(--color-text-primary);
}

.poll-result-count {
  font-size: 12px;
  color: var(--color-text-secondary);
}

/* ==================== Poll Form ==================== */
.poll-form-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.poll-form-option-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.poll-form-option-row .ant-input {
  flex: 1;
}

/* ==================== Tag Management ==================== */
.tag-manage-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
  max-height: 300px;
  overflow-y: auto;
}

.tag-manage-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--color-bg-container);
  border: 1px solid var(--color-border-secondary);
  border-radius: var(--radius-md);
}

.tag-manage-name {
  font-size: 14px;
  color: var(--color-text-primary);
}

.tag-empty-hint {
  text-align: center;
  color: var(--color-text-secondary);
  padding: 16px 0;
}

.tag-add-section {
  border-top: 1px solid var(--color-border-secondary);
  padding-top: 12px;
}

.tag-add-form {
  display: flex;
  gap: 8px;
}

.tag-add-form .ant-input {
  flex: 1;
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
