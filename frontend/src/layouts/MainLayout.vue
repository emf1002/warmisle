<template>
  <div class="main-layout">
    <!-- ==================== 桌面端侧边栏 ==================== -->
    <aside
      class="sidebar"
      aria-label="主导航"
      :class="{
        'sidebar-collapsed': sidebarCollapsed && !isMobile,
        'sidebar-hidden': isMobile,
      }"
      @mouseenter="onSidebarMouseEnter"
      @mouseleave="onSidebarMouseLeave"
    >
      <div class="sidebar-logo" @click="$router.push('/')">
        <LogoIcon :size="32" />
        <span class="sidebar-logo-text">暖屿</span>
      </div>

      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        :theme="isDarkTheme ? 'dark' : 'light'"
        class="sidebar-menu"
        @click="onMenuClick"
      >
        <a-menu-item key="Dashboard">
          <template #icon><Icon name="LayoutDashboard" :size="18" /></template>
          <span>仪表盘</span>
        </a-menu-item>
        <a-menu-item key="Ledger">
          <template #icon><Icon name="Wallet" :size="18" /></template>
          <span>记账本</span>
        </a-menu-item>
        <a-menu-item key="Todo">
          <template #icon><Icon name="ListTodo" :size="18" /></template>
          <span>待办管理</span>
        </a-menu-item>
        <a-menu-item key="Wish">
          <template #icon><Icon name="Star" :size="18" /></template>
          <span>愿望清单</span>
        </a-menu-item>
        <a-menu-item key="Forum">
          <template #icon><Icon name="MessageSquare" :size="18" /></template>
          <span>家庭论坛</span>
        </a-menu-item>

        <a-menu-divider />

        <a-menu-item key="Members" v-if="isAdmin">
          <template #icon><Icon name="Users" :size="18" /></template>
          <span>成员管理</span>
        </a-menu-item>
        <a-menu-item key="Categories" v-if="isAdmin">
          <template #icon><Icon name="FolderOpen" :size="18" /></template>
          <span>分类管理</span>
        </a-menu-item>
        <a-menu-item key="Profile">
          <template #icon><Icon name="UserCircle" :size="18" /></template>
          <span>个人中心</span>
        </a-menu-item>
      </a-menu>
    </aside>

    <!-- ==================== 主内容区 ==================== -->
    <div class="main-area" :class="{ 'main-full': isMobile, 'main-compact': sidebarCollapsed && !isMobile }">
      <!-- 顶部栏 -->
      <header
        class="topbar"
        :class="{
          'topbar-mobile': isMobile,
          'topbar-compact': sidebarCollapsed && !isMobile,
        }"
      >
        <div v-if="isMobile" class="topbar-left">
          <span class="page-title">{{ pageTitle }}</span>
        </div>
        <div v-else class="topbar-spacer" />

        <div class="topbar-right">
          <ThemeToggle />
          <a-dropdown>
            <div class="user-trigger">
              <a-avatar :size="isMobile ? 28 : 32" class="user-avatar">
                <template #icon>
                  <span>{{ currentMember?.avatar || '👤' }}</span>
                </template>
              </a-avatar>
              <span v-if="!isMobile" class="user-name">{{ currentMember?.name || '用户' }}</span>
            </div>
            <template #overlay>
              <a-menu @click="onUserMenuClick">
                <a-menu-item key="Profile">
                  <Icon name="UserCircle" :size="14" />
                  <span style="margin-left: 8px">个人中心</span>
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="Logout">
                  <Icon name="LogOut" :size="14" />
                  <span style="margin-left: 8px">退出登录</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="content" :class="{ 'content-mobile': isMobile }">
        <slot />
      </main>
    </div>

    <!-- ==================== 移动端底部TabBar ==================== -->
    <nav v-if="isMobile" class="tabbar" aria-label="底部导航" data-testid="mobile-tabbar">
      <div
        v-for="tab in bottomTabs"
        :key="tab.key"
        class="tabbar-item"
        :class="{ 'tabbar-item-active': activeBottomTab === tab.key }"
        @click="onBottomTabClick(tab.key)"
      >
        <span class="tabbar-icon"><Icon :name="tab.icon" :size="22" /></span>
        <span class="tabbar-label">{{ tab.label }}</span>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getProfile } from '@/api/member'
import { useThemeStore } from '@/stores/theme'
import ThemeToggle from '@/components/ThemeToggle.vue'
import LogoIcon from '@/components/LogoIcon.vue'
import Icon from '@/components/Icon.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const isDarkTheme = computed(() => themeStore.theme === 'dark')

// ---- 响应式检测 ----
const isMobile = ref(window.innerWidth < 768)
const isTablet = ref(window.innerWidth >= 768 && window.innerWidth < 1024)
const sidebarCollapsed = ref(false)

function onResize() {
  isMobile.value = window.innerWidth < 768
  isTablet.value = window.innerWidth >= 768 && window.innerWidth < 1024
  // 平板端默认折叠
  if (isTablet.value) {
    sidebarCollapsed.value = true
  } else if (window.innerWidth >= 1024) {
    sidebarCollapsed.value = false
  }
}

// 初始化时设置
if (isTablet.value) {
  sidebarCollapsed.value = true
}

function onSidebarMouseEnter() {
  if (isTablet.value) {
    sidebarCollapsed.value = false
  }
}

function onSidebarMouseLeave() {
  if (isTablet.value) {
    sidebarCollapsed.value = true
  }
}

onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

// ---- 当前成员信息 ----
const currentMember = ref<any>(authStore.memberInfo)

async function fetchProfile() {
  try {
    const res: any = await getProfile()
    currentMember.value = res.data
    authStore.memberInfo = res.data
  } catch {
    // ignore
  }
}

onMounted(() => {
  if (!currentMember.value) {
    fetchProfile()
  }
})

const isAdmin = computed(() => currentMember.value?.role === 'admin')

// ---- 菜单选中项 ----
const routeToMenuKey: Record<string, string> = {
  Dashboard: 'Dashboard',
  Ledger: 'Ledger',
  Todo: 'Todo',
  Wish: 'Wish',
  Forum: 'Forum',
  TopicDetail: 'Forum',
  Members: 'Members',
  Categories: 'Categories',
  Profile: 'Profile',
}

const selectedKeys = ref<string[]>([routeToMenuKey[route.name as string] || 'Dashboard'])

watch(
  () => route.name,
  (name) => {
    selectedKeys.value = [routeToMenuKey[name as string] || 'Dashboard']
  }
)

function onMenuClick({ key }: { key: string }) {
  router.push({ name: key })
}

// ---- 页面标题（移动端） ----
const routeTitleMap: Record<string, string> = {
  Dashboard: '仪表盘',
  Ledger: '记账本',
  Todo: '待办管理',
  Wish: '愿望清单',
  Forum: '家庭论坛',
  TopicDetail: '家庭论坛',
  Members: '成员管理',
  Categories: '分类管理',
  Profile: '个人中心',
}

const pageTitle = computed(() => routeTitleMap[route.name as string] || '暖屿')

// ---- 用户下拉菜单 ----
function onUserMenuClick({ key }: { key: string }) {
  if (key === 'Logout') {
    authStore.logout()
    router.push('/login')
  } else if (key === 'Profile') {
    router.push({ name: 'Profile' })
  }
}

// ---- 底部TabBar ----
const bottomTabs = [
  { key: 'Dashboard', icon: 'LayoutDashboard', label: '仪表盘', routeName: 'Dashboard' },
  { key: 'Ledger', icon: 'Wallet', label: '记账', routeName: 'Ledger' },
  { key: 'Todo', icon: 'ListTodo', label: '待办', routeName: 'Todo' },
  { key: 'Forum', icon: 'MessageSquare', label: '论坛', routeName: 'Forum' },
  { key: 'Profile', icon: 'UserCircle', label: '我的', routeName: 'Profile' },
]

const activeBottomTab = ref(routeToMenuKey[route.name as string] || 'Dashboard')

watch(
  () => route.name,
  (name) => {
    activeBottomTab.value = routeToMenuKey[name as string] || 'Dashboard'
  }
)

function onBottomTabClick(key: string) {
  const tab = bottomTabs.find((t) => t.key === key)
  if (tab) {
    router.push({ name: tab.routeName })
  }
}
</script>

<style scoped>
/* ==================== 全局 ==================== */
.main-layout {
  min-height: 100vh;
  background: var(--color-bg-layout);
}

/* 暗色主题添加背景渐变 */
[data-theme="dark"] .main-layout {
  background-image: var(--bg-gradient);
}

/* ==================== 侧边栏 ==================== */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 220px;
  background: var(--color-bg-sidebar);
  border-right: 1px solid var(--sidebar-border);
  display: flex;
  flex-direction: column;
  z-index: var(--z-sticky);
  overflow-y: auto;
  transition: width var(--duration-normal) ease;
}

[data-theme="dark"] .sidebar {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.sidebar-collapsed {
  width: 64px;
}

.sidebar-collapsed .sidebar-logo-text,
.sidebar-collapsed .ant-menu-item .ant-menu-title-content {
  display: none;
}

.sidebar-hidden {
  display: none;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 56px;
  padding: 0 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--sidebar-border);
  overflow: hidden;
  white-space: nowrap;
}

.sidebar-logo-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: 1px;
}

.sidebar-menu {
  flex: 1;
  border-inline-end: none !important;
}

/* 亮色主题菜单覆盖 */
[data-theme="light"] .sidebar-menu :deep(.ant-menu-item) {
  color: var(--sidebar-text);
}

[data-theme="light"] .sidebar-menu :deep(.ant-menu-item:hover) {
  background: var(--sidebar-hover-bg);
  color: var(--color-text-primary);
}

[data-theme="light"] .sidebar-menu :deep(.ant-menu-item-selected) {
  background: var(--sidebar-active-bg) !important;
  color: var(--sidebar-text-active) !important;
}

/* 暗色主题菜单覆盖 */
[data-theme="dark"] .sidebar-menu :deep(.ant-menu-item-selected) {
  background: var(--sidebar-active-bg) !important;
}

.menu-icon {
  font-size: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
}

/* ==================== 主内容区 ==================== */
.main-area {
  margin-left: 220px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  transition: margin-left var(--duration-normal) ease;
}

.main-compact {
  margin-left: 64px;
}

.main-full {
  margin-left: 0;
}

/* ==================== 顶部栏 ==================== */
.topbar {
  position: fixed;
  top: 0;
  right: 0;
  left: 220px;
  height: 56px;
  background: var(--topbar-bg);
  box-shadow: var(--topbar-shadow);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 24px;
  z-index: calc(var(--z-sticky) - 1);
  transition: left var(--duration-normal) ease;
}

[data-theme="dark"] .topbar {
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
}

.topbar-compact {
  left: 64px;
}

.topbar-mobile {
  left: 0;
  height: 48px;
  justify-content: space-between;
  padding: 0 16px;
}

.topbar-spacer {
  flex: 1;
}

.topbar-left {
  display: flex;
  align-items: center;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.user-trigger {
  display: flex;
  align-items: center;
  cursor: pointer;
  min-height: 44px;
  padding: 4px 8px;
  border-radius: var(--radius-md);
  transition: background var(--duration-fast) ease;
}

.user-trigger:hover {
  background: var(--sidebar-hover-bg);
}

.user-avatar {
  flex-shrink: 0;
}

.user-name {
  margin-left: 8px;
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: nowrap;
}

/* ==================== 内容区 ==================== */
.content {
  margin-top: 56px;
  flex: 1;
}

.content-mobile {
  margin-top: 48px;
  margin-bottom: 56px;
}

/* ==================== 底部TabBar（移动端） ==================== */
.tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: var(--tabbar-bg);
  border-top: 1px solid var(--color-border-secondary);
  display: flex;
  align-items: center;
  z-index: var(--z-sticky);
  padding-bottom: env(safe-area-inset-bottom, 0);
  box-shadow: var(--shadow-level-2);
}

[data-theme="dark"] .tabbar {
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 44px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: color var(--duration-fast) ease;
  -webkit-tap-highlight-color: transparent;
  position: relative;
}

.tabbar-item-active {
  color: var(--color-brand);
}

.tabbar-item-active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 3px;
  background: var(--color-brand);
  border-radius: 2px;
  transition: all var(--duration-normal) ease;
}

.tabbar-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 12px;
  transition: all var(--duration-fast) ease;
  color: var(--color-muted);
  margin-bottom: 2px;
}

.tabbar-item-active .tabbar-icon {
  background: var(--color-brand-bg);
  color: var(--color-brand);
}

.tabbar-label {
  font-size: 11px;
  margin-top: 2px;
  line-height: 1;
}

/* ==================== 响应式 ==================== */
@media (max-width: 767px) {
  .sidebar {
    display: none;
  }

  .main-area {
    margin-left: 0;
  }

  .topbar {
    left: 0;
    height: 48px;
    justify-content: space-between;
    padding: 0 16px;
  }
}
</style>
