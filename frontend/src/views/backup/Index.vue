<template>
  <div class="backup-page" data-testid="backup-page">
    <div class="page-header">
      <h2>网盘备份</h2>
    </div>

    <!-- 云盘配置卡片 -->
    <a-card title="阿里云盘配置" class="config-card" data-testid="config-card">
      <a-form :model="configForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="App ID">
              <a-input v-model:value="configForm.app_id" placeholder="输入阿里云盘应用的 App ID" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="App Secret">
              <a-input-password v-model:value="configForm.app_secret" placeholder="输入阿里云盘应用的 App Secret" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="回调地址 (Redirect URI)">
              <a-input v-model:value="configForm.redirect_uri" placeholder="例如 http://192.168.1.100:8080/api/backup/callback" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="云端备份目录">
              <a-input v-model:value="configForm.backup_dir" placeholder="默认 /warmisle-backups/" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="saveConfig" :loading="configSaving" data-testid="save-config-btn">保存配置</a-button>
            <a-button v-if="authStatus === 'pending_auth'" type="primary" @click="getAuthUrl" :loading="authLoading" data-testid="authorize-btn">
              授权阿里云盘
            </a-button>
            <a-tag v-if="authStatus === 'authorized'" color="success" data-testid="status-tag">已授权</a-tag>
            <a-tag v-else-if="authStatus === 'token_expired'" color="error" data-testid="status-tag">授权已过期</a-tag>
            <a-tag v-else-if="authStatus === 'pending_auth'" color="warning" data-testid="status-tag">待授权</a-tag>
            <a-tag v-else-if="authStatus === 'unconfigured'" data-testid="status-tag">未配置</a-tag>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 定时备份配置 -->
    <a-card title="定时备份" class="schedule-card" data-testid="schedule-card">
      <a-row :gutter="24">
        <a-col :span="6">
          <a-form-item label="启用自动备份">
            <a-switch v-model:checked="scheduleForm.schedule_enabled" data-testid="schedule-switch" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="备份时间">
            <a-time-picker v-model:value="scheduleTime" format="HH:mm" value-format="HH:mm" placeholder="选择时间" />
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item label="保留天数">
            <a-input-number v-model:value="scheduleForm.retention_days" :min="1" :max="365" data-testid="retention-input" />
          </a-form-item>
        </a-col>
        <a-col :span="4">
          <a-form-item label=" ">
            <a-button type="primary" @click="saveSchedule" :loading="scheduleSaving" data-testid="save-schedule-btn">保存</a-button>
          </a-form-item>
        </a-col>
      </a-row>
    </a-card>

    <!-- 手动备份 -->
    <a-card title="手动备份" class="backup-action-card" data-testid="backup-action-card">
      <a-space>
        <a-button type="primary" @click="triggerBackup" :loading="backing" :disabled="authStatus !== 'authorized'" data-testid="trigger-backup-btn">
          立即备份
        </a-button>
        <span v-if="lastBackupMsg" class="backup-msg">{{ lastBackupMsg }}</span>
      </a-space>
    </a-card>

    <!-- 备份历史 -->
    <a-card title="备份历史" class="history-card" data-testid="history-card">
      <a-table
        :columns="historyColumns"
        :data-source="historyList"
        :loading="historyLoading"
        :pagination="historyPagination"
        row-key="id"
        @change="onHistoryPageChange"
        data-testid="history-table"
      >
        <template #emptyText>
          <div data-testid="empty-state">暂无备份记录</div>
        </template>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'file_name'">
            {{ record.file_name || '-' }}
          </template>
          <template v-if="column.key === 'file_size'">
            {{ formatSize(record.file_size) }}
          </template>
          <template v-if="column.key === 'backup_type'">
            <a-tag v-if="record.backup_type === 'manual'" color="blue">手动</a-tag>
            <a-tag v-else-if="record.backup_type === 'scheduled'" color="green">定时</a-tag>
            <a-tag v-else-if="record.backup_type === 'pre_restore'" color="orange">恢复前备份</a-tag>
            <a-tag v-else>{{ record.backup_type }}</a-tag>
          </template>
          <template v-if="column.key === 'upload_status'">
            <a-tag v-if="record.upload_status === 'completed'" color="success">成功</a-tag>
            <a-tag v-else-if="record.upload_status === 'failed'" color="error">失败</a-tag>
            <a-tag v-else-if="record.upload_status === 'pending'" color="processing">进行中</a-tag>
            <a-tag v-else>{{ record.upload_status }}</a-tag>
          </template>
          <template v-if="column.key === 'created_at'">
            {{ formatTime(record.created_at) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-button v-if="record.upload_status === 'failed' && record.error_message" type="link" size="small" @click="showError(record.error_message)">
              查看错误
            </a-button>
            <a-button type="link" size="small" danger @click="deleteHistoryRecord(record)" data-testid="delete-history-btn">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 云端文件列表 -->
    <a-card title="云端备份文件" class="cloud-files-card" data-testid="cloud-files-card">
      <a-button type="default" @click="loadCloudFiles" :loading="cloudFilesLoading" style="margin-bottom: 16px" data-testid="refresh-cloud-btn">
        刷新云端文件列表
      </a-button>
      <a-table
        :columns="cloudFileColumns"
        :data-source="cloudFiles"
        :loading="cloudFilesLoading"
        row-key="file_id"
        :pagination="false"
        data-testid="cloud-files-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'file_name'">
            {{ record.file_name }}
          </template>
          <template v-if="column.key === 'size'">
            {{ formatSize(record.size) }}
          </template>
          <template v-if="column.key === 'create_time'">
            {{ formatTime(record.create_time) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="primary" size="small" danger @click="showRestoreModal(record)" data-testid="restore-btn">恢复到此备份</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 恢复确认弹窗 -->
    <a-modal
      v-model:open="restoreModalVisible"
      title="⚠️ 恢复备份"
      :confirm-loading="restoring"
      @ok="confirmRestore"
      :ok-button-props="{ danger: true }"
      ok-text="确认恢复"
      cancel-text="取消"
      data-testid="restore-modal"
    >
      <p><strong>警告：恢复将覆盖当前数据库，此操作不可逆！</strong></p>
      <p>建议先手动备份当前数据。</p>
      <p>将恢复文件：<strong>{{ restoreTargetName }}</strong></p>
      <a-form-item label='请输入确认文字："我已了解风险，确认恢复"'>
        <a-input v-model:value="restoreConfirmText" placeholder="请输入确认文字" data-testid="restore-confirm-input" />
      </a-form-item>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getConfig, saveConfig as apiSaveConfig, getAuthUrl as apiGetAuthUrl, triggerBackup as apiTriggerBackup, listCloudFiles, restoreBackup, listHistory, deleteHistory, getSchedule, saveSchedule as apiSaveSchedule } from '@/api/backup'

// ---- 配置 ----
const configForm = reactive({ app_id: '', app_secret: '', redirect_uri: '', backup_dir: '' })
const configSaving = ref(false)
const authStatus = ref('unconfigured')
const authLoading = ref(false)

async function loadConfig() {
  try {
    const res = await getConfig()
    const cfg = res.data
    if (cfg) {
      configForm.app_id = cfg.app_id || ''
      configForm.redirect_uri = cfg.redirect_uri || ''
      configForm.backup_dir = cfg.backup_dir || '/warmisle-backups/'
      authStatus.value = cfg.status || 'unconfigured'
    }
  } catch { /* ignore */ }
}

async function saveConfig() {
  configSaving.value = true
  try {
    const res = await apiSaveConfig({
      app_id: configForm.app_id,
      app_secret: configForm.app_secret,
      redirect_uri: configForm.redirect_uri,
      backup_dir: configForm.backup_dir
    })
    const cfg = res.data
    if (cfg) authStatus.value = cfg.status || 'pending_auth'
    message.success('配置已保存')
    configForm.app_secret = ''
  } catch (e: any) {
    message.error(e?.response?.data?.message || '保存配置失败')
  } finally {
    configSaving.value = false
  }
}

async function getAuthUrl() {
  authLoading.value = true
  try {
    const res = await apiGetAuthUrl()
    const { auth_url } = res.data
    if (auth_url) window.open(auth_url, '_blank')
  } catch (e: any) {
    message.error(e?.response?.data?.message || '获取授权链接失败')
  } finally {
    authLoading.value = false
  }
}

// ---- 定时备份 ----
const scheduleForm = reactive({ schedule_enabled: false, retention_days: 30 })
const scheduleTime = ref<any>(null)
const scheduleSaving = ref(false)

async function loadSchedule() {
  try {
    const res = await getSchedule()
    const d = res.data
    if (d) {
      scheduleForm.schedule_enabled = !!d.schedule_enabled
      scheduleForm.retention_days = d.retention_days || 30
      if (d.schedule_time) {
        scheduleTime.value = d.schedule_time
      }
    }
  } catch { /* ignore */ }
}

async function saveSchedule() {
  scheduleSaving.value = true
  try {
    const timeStr = scheduleTime.value || '03:00'
    await apiSaveSchedule({
      schedule_enabled: scheduleForm.schedule_enabled,
      schedule_time: timeStr,
      retention_days: scheduleForm.retention_days
    })
    message.success('定时备份配置已保存')
  } catch (e: any) {
    message.error(e?.response?.data?.message || '保存失败')
  } finally {
    scheduleSaving.value = false
  }
}

// ---- 手动备份 ----
const backing = ref(false)
const lastBackupMsg = ref('')

async function triggerBackup() {
  backing.value = true
  lastBackupMsg.value = ''
  try {
    await apiTriggerBackup()
    message.success('备份已上传到阿里云盘')
    lastBackupMsg.value = '备份成功'
    await loadHistoryData()
  } catch (e: any) {
    const errMsg = e?.response?.data?.message || '备份失败'
    message.error(errMsg)
    lastBackupMsg.value = '备份失败: ' + errMsg
    await loadHistoryData()
  } finally {
    backing.value = false
  }
}

// ---- 备份历史 ----
const historyColumns = [
  { title: '文件名', key: 'file_name' },
  { title: '大小', key: 'file_size' },
  { title: '类型', key: 'backup_type' },
  { title: '状态', key: 'upload_status' },
  { title: '时间', key: 'created_at' },
  { title: '操作', key: 'action' }
]
const historyList = ref<any[]>([])
const historyLoading = ref(false)
const historyPagination = reactive({ current: 1, pageSize: 20, total: 0 })

async function loadHistoryData(page = 1) {
  historyLoading.value = true
  try {
    const res = await listHistory({ page, page_size: historyPagination.pageSize })
    const d = res.data
    if (d) {
      historyList.value = d.list || []
      historyPagination.total = d.total || 0
      historyPagination.current = d.page || page
    }
  } catch { /* ignore */ }
  finally { historyLoading.value = false }
}

function onHistoryPageChange(pag: any) {
  loadHistoryData(pag.current)
}

async function deleteHistoryRecord(record: any) {
  try {
    await deleteHistory(record.id)
    message.success('已删除')
    await loadHistoryData()
  } catch (e: any) {
    message.error(e?.response?.data?.message || '删除失败')
  }
}

function showError(errMsg: string) {
  message.error(errMsg)
}

// ---- 云端文件 ----
const cloudFileColumns = [
  { title: '文件名', key: 'file_name' },
  { title: '大小', key: 'size' },
  { title: '时间', key: 'create_time' },
  { title: '操作', key: 'action' }
]
const cloudFiles = ref<any[]>([])
const cloudFilesLoading = ref(false)

async function loadCloudFiles() {
  cloudFilesLoading.value = true
  try {
    const res = await listCloudFiles()
    cloudFiles.value = res.data || []
  } catch (e: any) {
    message.error(e?.response?.data?.message || '获取云端文件列表失败')
  } finally {
    cloudFilesLoading.value = false
  }
}

// ---- 恢复 ----
const restoreModalVisible = ref(false)
const restoreTargetId = ref('')
const restoreTargetName = ref('')
const restoreConfirmText = ref('')
const restoring = ref(false)

function showRestoreModal(file: any) {
  restoreTargetId.value = file.file_id
  restoreTargetName.value = file.file_name
  restoreConfirmText.value = ''
  restoreModalVisible.value = true
}

async function confirmRestore() {
  if (restoreConfirmText.value !== '我已了解风险，确认恢复') {
    message.warning('请输入正确的确认文字')
    return
  }
  restoring.value = true
  try {
    await restoreBackup(restoreTargetId.value, restoreConfirmText.value)
    message.success('恢复已启动，系统即将重启...')
    restoreModalVisible.value = false
  } catch (e: any) {
    message.error(e?.response?.data?.message || '恢复失败')
  } finally {
    restoring.value = false
  }
}

// ---- 工具函数 ----
function formatSize(bytes: number): string {
  if (!bytes || bytes === 0) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function formatTime(t: string): string {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

// ---- 初始化 ----
onMounted(async () => {
  await loadConfig()
  await loadSchedule()
  await loadHistoryData()
})
</script>

<style scoped>
.backup-page {
  max-width: 900px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
}

.config-card,
.schedule-card,
.backup-action-card,
.history-card,
.cloud-files-card {
  margin-bottom: 16px;
}

.backup-msg {
  color: #666;
  font-size: 13px;
}
</style>
