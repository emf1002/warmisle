<template>
  <div class="member-page">
    <div class="page-header">
      <h2>成员管理</h2>
      <a-button v-if="isAdmin" type="primary" @click="openCreate">
        添加成员
      </a-button>
    </div>

    <a-table :columns="columns" :data-source="members" row-key="id" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <span class="member-name">
            <span class="member-avatar">{{ record.avatar }}</span>
            {{ record.name }}
          </span>
        </template>
        <template v-if="column.key === 'username'">
          <span class="member-username">{{ record.username }}</span>
        </template>
        <template v-if="column.key === 'role'">
          <a-tag :color="record.role === 'admin' ? 'blue' : 'green'">
            {{ record.role === 'admin' ? '管理员' : '普通成员' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 'disabled' ? 'red' : 'green'">
            {{ record.status === 'disabled' ? '已禁用' : '正常' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'action'">
          <a-space v-if="isAdmin">
            <a-button type="link" size="small" @click="openEdit(record)">编辑</a-button>
            <a-button
              v-if="record.status === 'active'"
              type="link"
              size="small"
              danger
              @click="confirmDisable(record)"
            >
              禁用
            </a-button>
            <a-button
              v-else
              type="link"
              size="small"
              @click="handleEnable(record.id)"
            >
              启用
            </a-button>
            <a-button type="link" size="small" @click="confirmResetPwd(record)">
              重置密码
            </a-button>
            <a-button type="link" size="small" danger @click="confirmDelete(record)">
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- Create/Edit Dialog -->
    <a-modal
      v-model:open="dialogOpen"
      :title="editingMember ? '编辑成员' : '添加成员'"
      @ok="handleSubmit"
      :confirm-loading="submitting"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item v-if="!editingMember" label="用户名" required>
          <a-input
            v-model:value="form.username"
            placeholder="3-20位字母、数字、下划线"
            :maxlength="20"
          />
        </a-form-item>
        <a-form-item v-if="!editingMember" label="密码" required>
          <a-input-password
            v-model:value="form.password"
            placeholder="6-32位密码"
            :maxlength="32"
          />
        </a-form-item>
        <a-form-item label="名称">
          <a-input
            v-model:value="form.name"
            placeholder="1-20位显示名称"
            :maxlength="20"
          />
        </a-form-item>
        <a-form-item label="头像">
          <div class="emoji-picker">
            <span
              v-for="e in emojiList"
              :key="e"
              :class="['emoji-item', { active: form.avatar === e }]"
              @click="form.avatar = e"
            >{{ e }}</span>
          </div>
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model:value="form.role">
            <a-select-option value="member">普通成员</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  getMembers,
  createMember,
  updateMember,
  deleteMember,
  disableMember,
  enableMember,
  resetPassword,
} from '@/api/member'

const emojiList = [
  '👨','👩','👦','👧','👶','👴','👵','🐶','🐱','🏠',
  '💪','🎯','📚','🎮','🎨','🏀','🏊','🚗','✈️','🎵',
  '📷','💰','🔑','🌟','🔥','❤️','🍀','⭐','🎪','🐕','🐈','🏸'
]

const members = ref<any[]>([])
const dialogOpen = ref(false)
const editingMember = ref<any>(null)
const submitting = ref(false)

const form = reactive({
  username: '',
  password: '',
  name: '',
  avatar: '👨',
  role: 'member',
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

const columns = computed(() => {
  const cols: any[] = [
    { title: '成员', key: 'name', dataIndex: 'name' },
    { title: '用户名', key: 'username', dataIndex: 'username' },
    { title: '角色', key: 'role', dataIndex: 'role' },
    { title: '状态', key: 'status', dataIndex: 'status' },
  ]
  if (isAdmin.value) {
    cols.push({ title: '操作', key: 'action', width: 280 })
  }
  return cols
})

onMounted(() => {
  fetchMembers()
})

async function fetchMembers() {
  try {
    const res: any = await getMembers()
    members.value = res.data || []
  } catch {
    // error handled by interceptor
  }
}

function openCreate() {
  editingMember.value = null
  form.username = ''
  form.password = ''
  form.name = ''
  form.avatar = '👨'
  form.role = 'member'
  dialogOpen.value = true
}

function openEdit(record: any) {
  editingMember.value = record
  form.username = record.username
  form.password = ''
  form.name = record.name
  form.avatar = record.avatar
  form.role = record.role
  dialogOpen.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingMember.value) {
      await updateMember(editingMember.value.id, {
        name: form.name,
        avatar: form.avatar,
        role: form.role,
      })
      message.success('更新成功')
    } else {
      if (!form.username || !form.password) {
        message.error('用户名和密码不能为空')
        submitting.value = false
        return
      }
      await createMember({
        username: form.username,
        password: form.password,
        name: form.name,
        avatar: form.avatar,
        role: form.role,
      })
      message.success('添加成功')
    }
    dialogOpen.value = false
    fetchMembers()
  } catch (e: any) {
    // error handled by interceptor; show fallback if data has message
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  } finally {
    submitting.value = false
  }
}

function confirmDelete(record: any) {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除成员「${record.name}」吗？如果有活动记录则无法删除。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteMember(record.id)
        message.success('删除成功')
        fetchMembers()
      } catch (e: any) {
        if (e?.response?.data?.message) {
          message.error(e.response.data.message)
        }
        throw e // prevent modal from closing on error
      }
    },
  })
}

function confirmDisable(record: any) {
  Modal.confirm({
    title: '确认禁用',
    content: `确定要禁用成员「${record.name}」吗？`,
    okText: '禁用',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await disableMember(record.id)
        message.success('已禁用')
        fetchMembers()
      } catch (e: any) {
        if (e?.response?.data?.message) {
          message.error(e.response.data.message)
        }
        throw e
      }
    },
  })
}

async function handleEnable(id: number) {
  try {
    await enableMember(id)
    message.success('已启用')
    fetchMembers()
  } catch (e: any) {
    if (e?.response?.data?.message) {
      message.error(e.response.data.message)
    }
  }
}

function confirmResetPwd(record: any) {
  Modal.confirm({
    title: '重置密码',
    content: `确定要将成员「${record.name}」的密码重置为默认密码吗？`,
    okText: '确认',
    cancelText: '取消',
    async onOk() {
      try {
        await resetPassword(record.id)
        message.success('密码已重置')
      } catch (e: any) {
        if (e?.response?.data?.message) {
          message.error(e.response.data.message)
        }
        throw e
      }
    },
  })
}
</script>

<style scoped>
.member-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}

.member-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.member-avatar {
  font-size: 20px;
  line-height: 1;
}

.emoji-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.emoji-item {
  font-size: 24px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  border: 2px solid transparent;
  transition: border-color 0.2s;
  min-width: 40px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-item:hover {
  border-color: #d9d9d9;
}

.emoji-item.active {
  border-color: #1890ff;
  background: #e6f7ff;
}
</style>
