<template>
  <div class="profile-page" data-testid="profile-page">
    <div v-if="loading" class="skeleton-profile">
      <div class="skeleton-profile-header">
        <SkeletonCard width="64px" height="64px" borderRadius="50%" />
        <SkeletonCard width="120px" height="20px" />
        <SkeletonCard width="80px" height="16px" />
      </div>
      <SkeletonCard width="100%" height="40px" />
      <SkeletonCard width="100%" height="40px" />
      <SkeletonCard width="100%" height="40px" />
    </div>

    <a-card v-else :bordered="false" class="profile-card">
      <!-- User Info Header -->
      <div class="profile-header">
        <UserAvatar :icon="profile.avatar || 'User'" :member-id="profile.id" :size="64" />
        <div class="profile-name" data-testid="profile-name">{{ profile.name }}</div>
        <a-tag :color="profile.role === 'admin' ? 'blue' : 'default'" class="profile-role-tag">
          {{ profile.role === 'admin' ? '管理员' : '成员' }}
        </a-tag>
        <div class="profile-username">@{{ profile.username }}</div>
      </div>

      <a-divider />

      <!-- Info Section -->
      <div class="profile-info">
        <div class="info-item">
          <span class="info-label">用户名</span>
          <span class="info-value">{{ profile.username }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">姓名</span>
          <span class="info-value">{{ profile.name }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">角色</span>
          <span class="info-value">
            <a-tag :color="profile.role === 'admin' ? 'blue' : 'default'">
              {{ profile.role === 'admin' ? '管理员' : '成员' }}
            </a-tag>
          </span>
        </div>
      </div>

      <!-- Navigation -->
      <div class="profile-nav">
        <div class="nav-item" @click="goCategories">
          <span class="nav-item-icon">📂</span>
          <span class="nav-item-label">账单分类</span>
          <span class="nav-item-arrow">›</span>
        </div>
        <div class="nav-item" @click="goMembers">
          <span class="nav-item-icon">👥</span>
          <span class="nav-item-label">家庭成员</span>
          <span class="nav-item-arrow">›</span>
        </div>
      </div>

      <a-divider />

      <!-- Actions -->
      <div class="profile-actions">
        <a-button type="primary" block size="large" @click="openEdit" :style="{ minHeight: '44px' }" data-testid="edit-profile-btn">
          ✏️ 修改个人信息
        </a-button>
        <a-button
          block
          size="large"
          @click="handleChangePassword"
          :style="{ minHeight: '44px', marginTop: '12px' }"
          data-testid="change-pwd-btn"
        >
          🔒 修改密码
        </a-button>
        <a-button
          danger
          block
          size="large"
          @click="confirmLogout"
          :style="{ minHeight: '44px', marginTop: '12px' }"
          data-testid="logout-btn"
        >
          🚪 退出登录
        </a-button>
      </div>
    </a-card>

    <!-- Edit Profile Dialog -->
    <a-modal
      v-model:open="dialogOpen"
      title="修改个人信息"
      @ok="handleSubmit"
      :confirm-loading="submitting"
      cancel-text="取消"
      ok-text="保存"
      :ok-button-props="({ 'data-testid': 'modal-submit-btn' } as any)"
    >
      <div data-testid="profile-modal">
      <a-form ref="profileFormRef" :model="form" :rules="profileRules" layout="vertical">
        <a-form-item label="姓名" name="name">
          <a-input
            v-model:value="form.name"
            placeholder="1-20位显示名称"
            data-testid="name-input"
          />
        </a-form-item>
        <a-form-item label="头像">
          <div data-testid="avatar-picker">
            <IconPicker v-model="form.avatar" />
          </div>
        </a-form-item>
      </a-form>
      </div>
    </a-modal>

    <!-- Change Password Dialog -->
    <a-modal
      v-model:open="pwdDialogOpen"
      title="🔒 修改密码"
      @ok="handlePwdSubmit"
      :confirm-loading="pwdSubmitting"
      cancel-text="取消"
      ok-text="保存"
      :ok-button-props="({ 'data-testid': 'modal-submit-btn' } as any)"
    >
      <div data-testid="pwd-modal">
      <a-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" layout="vertical">
        <a-form-item label="当前密码" name="old_password" required :help="oldPasswordError" :validate-status="oldPasswordError ? 'error' : ''">
          <a-input-password v-model:value="pwdForm.old_password" placeholder="请输入当前密码" data-testid="old-pwd-input" />
        </a-form-item>
        <a-form-item label="新密码" name="new_password" required>
          <a-input-password v-model:value="pwdForm.new_password" placeholder="6-32位新密码" data-testid="new-pwd-input" />
        </a-form-item>
        <a-form-item label="确认新密码" name="confirm_password" required>
          <a-input-password v-model:value="pwdForm.confirm_password" placeholder="请再次输入新密码" data-testid="confirm-pwd-input" />
        </a-form-item>
      </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import { getProfile, updateProfile, changePassword } from '@/api/member'
import { useAuthStore } from '@/stores/auth'
import IconPicker from '@/components/IconPicker.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const profile = ref<any>({
  id: 0,
  username: '',
  name: '',
  avatar: 'User',
  role: 'member',
})

const dialogOpen = ref(false)
const submitting = ref(false)
const profileFormRef = ref<FormInstance>()

const form = reactive({
  name: '',
  avatar: 'User',
})

const profileRules = {
  name: [
    { required: true, message: '姓名不能为空', trigger: 'blur' },
    { max: 20, message: '请输入1-20字符', trigger: 'blur' },
  ],
}

onMounted(() => {
  fetchProfile()
})

async function fetchProfile() {
  loading.value = true
  try {
    const res: any = await getProfile()
    if (res.data) {
      profile.value = res.data
    }
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

function openEdit() {
  form.name = profile.value.name
  form.avatar = profile.value.avatar
  dialogOpen.value = true
}

async function handleSubmit() {
  try {
    await profileFormRef.value?.validate()
  } catch {
    return
  }
  submitting.value = true
  try {
    const res: any = await updateProfile({
      name: form.name.trim(),
      avatar: form.avatar,
    })
    if (res.data) {
      profile.value = res.data
    }
    message.success('✅ 修改成功')
    dialogOpen.value = false
  } catch (e: any) {
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  } finally {
    submitting.value = false
  }
}

function confirmLogout() {
  Modal.confirm({
    title: '❓ 确认退出',
    content: '确定要退出登录吗？',
    okText: '退出',
    okType: 'danger',
    cancelText: '取消',
    okButtonProps: { 'data-testid': 'modal-confirm-btn' } as any,
    onOk() {
      authStore.logout()
      message.success('✅ 已退出登录')
      router.push('/login')
    },
  })
}

function goCategories() {
  router.push('/categories')
}

function goMembers() {
  router.push('/members')
}

const pwdDialogOpen = ref(false)
const pwdSubmitting = ref(false)
const pwdFormRef = ref<FormInstance>()
const oldPasswordError = ref('')
const pwdForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

const pwdRules = {
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 32, message: '新密码长度需在6-32位之间', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value && value === pwdForm.old_password) {
          return Promise.reject('新密码不能与旧密码相同')
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value && value !== pwdForm.new_password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      },
      trigger: 'blur',
    },
  ],
}

function handleChangePassword() {
  pwdForm.old_password = ''
  pwdForm.new_password = ''
  pwdForm.confirm_password = ''
  oldPasswordError.value = ''
  pwdDialogOpen.value = true
}

async function handlePwdSubmit() {
  try {
    await pwdFormRef.value?.validate()
  } catch {
    return
  }
  pwdSubmitting.value = true
  try {
    await changePassword({
      old_password: pwdForm.old_password,
      new_password: pwdForm.new_password,
    })
    message.success('✅ 密码修改成功')
    pwdDialogOpen.value = false
  } catch (e: any) {
    if (e?.response?.data?.message) {
      oldPasswordError.value = '原密码错误'
    }
  } finally {
    pwdSubmitting.value = false
  }
}
</script>

<style scoped>
.profile-page {
  padding: var(--space-lg);
  display: flex;
  justify-content: center;
}

.profile-card {
  width: 100%;
  max-width: 600px;
}

/* Header */
.profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 0 8px;
}

.profile-avatar-large {
  font-size: 64px;
  line-height: 1;
  margin-bottom: 12px;
}

.profile-name {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 6px;
}

.profile-role-tag {
  margin-bottom: 6px;
}

.profile-username {
  font-size: 13px;
  color: var(--color-text-secondary);
}

/* Info */
.profile-info {
  padding: 0 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.info-item + .info-item {
  border-top: 1px solid var(--color-border-secondary);
}

.info-label {
  color: var(--color-text-secondary);
  font-size: 14px;
}

.info-value {
  font-size: 14px;
  color: var(--color-text-primary);
}

/* Actions */
.profile-actions {
  padding: 0 8px;
}

/* Navigation */
.profile-nav {
  padding: 0 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 0;
  min-height: 44px;
  cursor: pointer;
  border-bottom: 1px solid var(--color-border-secondary);
  transition: background 0.2s;
}

.nav-item:last-child {
  border-bottom: none;
}

.nav-item:hover {
  background: var(--color-border-secondary);
  margin: 0 -8px;
  padding: 14px 8px;
  border-radius: 6px;
}

.nav-item-icon {
  font-size: 18px;
}

.nav-item-label {
  flex: 1;
  font-size: 15px;
  color: var(--color-text-primary);
}

.nav-item-arrow {
  font-size: 18px;
  color: var(--color-text-disabled);
}


@media (max-width: 767px) {
  .profile-page {
    padding: 0;
  }

  .profile-card {
    max-width: 100%;
    border-radius: 0;
  }

  :deep(.ant-card) {
    border-radius: 0;
  }

  .skeleton-profile {
    border-radius: 0;
  }
}

/* ==================== Skeleton ==================== */
.skeleton-profile {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  padding: var(--space-lg);
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  width: 100%;
  max-width: 600px;
  box-sizing: border-box;
}

.skeleton-profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-sm);
  margin-bottom: var(--space-md);
}
</style>
