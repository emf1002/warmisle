# 明暗主题切换 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为暖屿添加明暗两套主题（晨屿 Light / 夜屿 Dark），支持手动切换，持久化用户偏好。暗色主题使用 Ant Design darkAlgorithm + 自定义 CSS 变量，亮主题沿用现有风格并升级为预览 V3 设计。登录页根据主题渲染不同背景场景。

**Architecture:** 在 `<html>` 元素上设置 `data-theme="light|dark"` 属性，通过 CSS 变量驱动所有主题相关颜色。Pinia store 管理主题状态并持久化到 localStorage。App.vue 的 `themeConfig` 改为 computed 属性，根据当前主题动态切换 Ant Design token 和 algorithm。MainLayout 侧边栏菜单根据主题切换 dark/light 模式。AuthLayout 和 Login 页面根据主题渲染不同背景（亮：有机色块 + 品牌介绍；暗：星空月亮 + 玻璃拟态卡片）。

**Tech Stack:** Vue 3 Composition API, Pinia, Ant Design Vue 4.x (ConfigProvider + darkAlgorithm), CSS Custom Properties

---

## File Structure

| Action | File | Responsibility |
|--------|------|----------------|
| Create | `frontend/src/stores/theme.ts` | 主题状态管理，localStorage 持久化 |
| Create | `frontend/src/styles/themes.css` | 两套主题的 CSS 变量定义 |
| Create | `frontend/src/components/ThemeToggle.vue` | 主题切换按钮组件 |
| Modify | `frontend/src/styles/global.css` | 删除旧变量定义，导入 themes.css，更新 body 背景 |
| Modify | `frontend/src/main.ts` | 导入 themes.css |
| Modify | `frontend/src/App.vue` | themeConfig 改为 computed，接入 theme store |
| Modify | `frontend/src/layouts/MainLayout.vue` | 侧边栏/顶栏/TabBar 主题适配，添加 ThemeToggle |
| Modify | `frontend/src/layouts/AuthLayout.vue` | 双主题背景：亮=有机色块，暗=星空场景 |
| Modify | `frontend/src/views/auth/Login.vue` | 双主题卡片样式 |
| Modify | `frontend/src/views/auth/Init.vue` | 同 Login.vue 风格适配 |
| Modify | `frontend/src/components/EmptyState.vue` | 硬编码颜色替换为 CSS 变量 |
| Modify | `frontend/src/views/dashboard/Index.vue` | 图表渐变色主题适配 |
| Modify | `frontend/src/views/ledger/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/todo/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/wish/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/forum/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/forum/TopicDetail.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/member/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/category/Index.vue` | 硬编码颜色检查 |
| Modify | `frontend/src/views/profile/Index.vue` | 硬编码颜色检查 |

---

### Task 1: 创建主题 Store

**Files:**
- Create: `frontend/src/stores/theme.ts`

- [ ] **Step 1: 创建 theme store 文件**

```typescript
// frontend/src/stores/theme.ts
import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const STORAGE_KEY = 'warmisle-theme'

  // 从 localStorage 读取，默认 light
  function getInitialTheme(): ThemeMode {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === 'light' || saved === 'dark') return saved
    return 'light'
  }

  const theme = ref<ThemeMode>(getInitialTheme())

  // 切换主题
  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  // 设置指定主题
  function setTheme(mode: ThemeMode) {
    theme.value = mode
  }

  // 同步到 <html data-theme="..."> 和 localStorage
  function applyTheme(mode: ThemeMode) {
    document.documentElement.setAttribute('data-theme', mode)
    localStorage.setItem(STORAGE_KEY, mode)
  }

  // 初始化时立即应用
  applyTheme(theme.value)

  // 监听变化自动应用
  watch(theme, (newTheme) => {
    applyTheme(newTheme)
  })

  return { theme, toggleTheme, setTheme }
})
```

- [ ] **Step 2: 验证 store 编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无 theme.ts 相关错误

- [ ] **Step 3: Commit**

```bash
git add frontend/src/stores/theme.ts
git commit -m "feat: add theme store with localStorage persistence"
```

---

### Task 2: 创建主题 CSS 变量系统

**Files:**
- Create: `frontend/src/styles/themes.css`
- Modify: `frontend/src/styles/global.css:1-47`
- Modify: `frontend/src/main.ts:1`

- [ ] **Step 1: 创建 themes.css — 两套完整变量定义**

```css
/* frontend/src/styles/themes.css */

/* ==================== 亮色主题（晨屿） ==================== */
[data-theme="light"] {
  /* 品牌色 */
  --color-brand: #E87461;
  --color-brand-light: rgba(232, 116, 97, 0.1);
  --color-brand-hover: #D05A48;

  /* 语义色 */
  --color-success: #389E0D;
  --color-success-light: #F6FFED;
  --color-warning: #E8A830;
  --color-warning-light: #FFFBE6;
  --color-danger: #E85D5D;
  --color-danger-light: #FFF2F0;

  /* 文字 */
  --color-text-primary: #3D3530;
  --color-text-secondary: #8A807A;
  --color-text-disabled: #B8AFA8;
  --color-muted: #8A807A;

  /* 背景 */
  --color-bg-layout: #FDF8F4;
  --color-bg-container: #FFFFFF;
  --color-bg-sidebar: #FDF8F4;
  --color-bg-elevated: #FFFFFF;

  /* 边框 */
  --color-border: rgba(61, 53, 48, 0.08);
  --color-border-secondary: rgba(61, 53, 48, 0.05);

  /* 阴影 */
  --shadow-level-0: none;
  --shadow-level-1: 0 1px 4px rgba(61, 53, 48, 0.05);
  --shadow-level-2: 0 2px 8px rgba(61, 53, 48, 0.08);
  --shadow-level-3: 0 4px 16px rgba(61, 53, 48, 0.08);
  --shadow-level-4: 0 8px 24px rgba(61, 53, 48, 0.1);

  /* 圆角 */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;

  /* 特有变量 */
  --sidebar-text: var(--color-text-secondary);
  --sidebar-text-active: var(--color-brand);
  --sidebar-border: var(--color-border);
  --sidebar-hover-bg: rgba(232, 116, 97, 0.06);
  --sidebar-active-bg: rgba(232, 116, 97, 0.1);
  --topbar-bg: var(--color-bg-container);
  --topbar-shadow: 0 1px 4px rgba(61, 53, 48, 0.05);
  --tabbar-bg: var(--color-bg-container);
  --card-bg: var(--color-bg-container);
  --card-border: var(--color-border);
  --input-bg: var(--color-bg-container);
  --input-border: var(--color-border);
  --dropdown-bg: var(--color-bg-container);
  --dropdown-shadow: var(--shadow-level-3);

  /* 图表色 */
  --chart-color-1: #E87461;
  --chart-color-2: #6BBAA7;
  --chart-color-3: #F2B84B;
  --chart-color-4: #B8D4E8;
  --chart-color-5: #C5B3D9;
  --chart-color-6: #F09888;
  --chart-color-7: #EDE9E4;

  /* 认证页 */
  --auth-bg: var(--color-bg-layout);
  --auth-card-bg: rgba(255, 255, 255, 0.72);
  --auth-card-border: rgba(255, 255, 255, 0.6);
  --auth-card-shadow: 0 4px 24px rgba(61, 53, 48, 0.06);
}

/* ==================== 暗色主题（夜屿） ==================== */
[data-theme="dark"] {
  /* 品牌色 */
  --color-brand: #E5A84B;
  --color-brand-light: rgba(229, 168, 75, 0.12);
  --color-brand-hover: #F0C273;

  /* 语义色 */
  --color-success: #52C41A;
  --color-success-light: rgba(82, 196, 26, 0.12);
  --color-warning: #FAAD14;
  --color-warning-light: rgba(250, 173, 20, 0.12);
  --color-danger: #FF6B6B;
  --color-danger-light: rgba(255, 107, 107, 0.12);

  /* 文字 */
  --color-text-primary: #E8E4DE;
  --color-text-secondary: #8A94A8;
  --color-text-disabled: #5A6478;
  --color-muted: #5A6478;

  /* 背景 */
  --color-bg-layout: #0F1729;
  --color-bg-container: rgba(26, 35, 64, 0.65);
  --color-bg-sidebar: rgba(15, 23, 41, 0.95);
  --color-bg-elevated: #1A2340;

  /* 边框 */
  --color-border: rgba(255, 255, 255, 0.06);
  --color-border-secondary: rgba(255, 255, 255, 0.04);

  /* 阴影 */
  --shadow-level-0: none;
  --shadow-level-1: 0 1px 4px rgba(0, 0, 0, 0.2);
  --shadow-level-2: 0 2px 8px rgba(0, 0, 0, 0.3);
  --shadow-level-3: 0 4px 16px rgba(0, 0, 0, 0.3);
  --shadow-level-4: 0 8px 24px rgba(0, 0, 0, 0.4);

  /* 圆角 */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 14px;

  /* 特有变量 */
  --sidebar-text: var(--color-text-secondary);
  --sidebar-text-active: var(--color-brand);
  --sidebar-border: rgba(255, 255, 255, 0.06);
  --sidebar-hover-bg: rgba(255, 255, 255, 0.05);
  --sidebar-active-bg: rgba(229, 168, 75, 0.12);
  --topbar-bg: rgba(15, 23, 41, 0.8);
  --topbar-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
  --tabbar-bg: rgba(15, 23, 41, 0.95);
  --card-bg: rgba(26, 35, 64, 0.65);
  --card-border: rgba(255, 255, 255, 0.06);
  --input-bg: rgba(255, 255, 255, 0.04);
  --input-border: rgba(255, 255, 255, 0.08);
  --dropdown-bg: #1A2340;
  --dropdown-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);

  /* 图表色 */
  --chart-color-1: #E5A84B;
  --chart-color-2: #F0C273;
  --chart-color-3: #6BBAA7;
  --chart-color-4: #B8D4E8;
  --chart-color-5: #C5B3D9;
  --chart-color-6: #FF6B6B;
  --chart-color-7: rgba(255, 255, 255, 0.08);

  /* 认证页 */
  --auth-bg: #0F1729;
  --auth-card-bg: rgba(26, 35, 64, 0.85);
  --auth-card-border: rgba(255, 255, 255, 0.08);
  --auth-card-shadow: 0 8px 40px rgba(0, 0, 0, 0.4);

  /* 背景渐变（暗色主题用） */
  --bg-gradient: radial-gradient(ellipse 80% 60% at 20% 10%, rgba(36, 48, 84, 0.5) 0%, transparent 60%),
                 radial-gradient(ellipse 60% 50% at 80% 90%, rgba(40, 30, 60, 0.3) 0%, transparent 50%);
}

/* ==================== 全局 backdrop-filter 支持 ==================== */
[data-theme="dark"] .ant-card,
[data-theme="dark"] .sidebar,
[data-theme="dark"] .topbar {
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}
```

- [ ] **Step 2: 更新 global.css — 替换旧变量为导入，保留非变量部分**

将 `global.css` 的第 1-47 行（`:root` 变量块）替换为对 `themes.css` 的引用说明，并在 `body` 中使用变量。实际操作：删除旧 `:root` 块，body 保持不变（已经使用变量）。

修改 `global.css`：删除第 2-47 行的 `:root { ... }` 块（保留第 1 行注释和第 48 行以后的所有内容）。在文件顶部添加：

```css
/* 主题变量定义见 themes.css（由 main.ts 导入） */
```

具体替换：将整个文件的第 1-47 行（从 `/* ==================== 全局 CSS 变量 ==================== */` 到 `}` 结束）替换为：

```css
/* ==================== 全局 CSS 变量 ==================== */
/* 主题变量定义见 themes.css（由 main.ts 导入） */
/* 以下为非主题相关的全局变量 */
:root {
  --space-xxs: 4px;
  --space-xs: 8px;
  --space-sm: 12px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif;
  --z-base: 0;
  --z-sticky: 10;
  --z-dropdown: 20;
  --z-sticky-header: 30;
  --z-overlay: 40;
  --z-modal: 50;
  --z-message: 60;
  --z-notification: 70;
  --duration-fast: 150ms;
  --duration-normal: 200ms;
  --duration-slow: 300ms;
  --max-content-width: 1200px;
}
```

同时删除 global.css 末尾的暗色模式注释块（第 167-176 行），因为已由 themes.css 实现。

- [ ] **Step 3: 更新 main.ts — 导入 themes.css**

在 `main.ts` 第 1 行之前添加 `themes.css` 导入：

```typescript
import '@/styles/themes.css'
import '@/styles/global.css'
```

- [ ] **Step 4: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无新错误

- [ ] **Step 5: Commit**

```bash
git add frontend/src/styles/themes.css frontend/src/styles/global.css frontend/src/main.ts
git commit -m "feat: add dual-theme CSS variable system (Morning Isle light / Night Isle dark)"
```

---

### Task 3: 更新 App.vue — 动态 Ant Design 主题

**Files:**
- Modify: `frontend/src/App.vue:13-58`

- [ ] **Step 1: 重写 App.vue script 部分**

将 `App.vue` 的 `<script setup lang="ts">` 部分替换为：

```vue
<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { theme } from 'ant-design-vue'
import MainLayout from '@/layouts/MainLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { useThemeStore } from '@/stores/theme'

const route = useRoute()
const themeStore = useThemeStore()

const layoutMap: Record<string, any> = {
  auth: AuthLayout,
  main: MainLayout,
}

const layout = computed(() => {
  const layoutName = (route.meta.layout as string) || 'main'
  return layoutMap[layoutName] || MainLayout
})

const themeConfig = computed(() => {
  const isDark = themeStore.theme === 'dark'
  return {
    algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: isDark ? '#E5A84B' : '#E87461',
      colorSuccess: isDark ? '#52C41A' : '#389E0D',
      colorWarning: '#FAAD14',
      colorError: isDark ? '#FF6B6B' : '#E85D5D',
      colorInfo: isDark ? '#E5A84B' : '#E87461',
      borderRadius: isDark ? 10 : 12,
      fontSize: 14,
      fontFamily: '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif',
      colorBgLayout: isDark ? '#0F1729' : '#FDF8F4',
      colorBgContainer: isDark ? '#1A2340' : '#FFFFFF',
      colorBgElevated: isDark ? '#1A2340' : '#FFFFFF',
      colorText: isDark ? '#E8E4DE' : '#3D3530',
      colorTextSecondary: isDark ? '#8A94A8' : '#8A807A',
      colorBorder: isDark ? 'rgba(255,255,255,0.06)' : 'rgba(61,53,48,0.08)',
      colorBorderSecondary: isDark ? 'rgba(255,255,255,0.04)' : 'rgba(61,53,48,0.05)',
    },
    components: {
      Button: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Input: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Select: {
        borderRadius: isDark ? 10 : 12,
        controlHeight: 44,
      },
      Card: {
        borderRadiusLG: isDark ? 12 : 14,
      },
      Modal: {
        borderRadiusLG: isDark ? 12 : 14,
      },
    },
  }
})
</script>
```

- [ ] **Step 2: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add frontend/src/App.vue
git commit -m "feat: dynamic Ant Design theme config with dark algorithm"
```

---

### Task 4: 创建主题切换组件

**Files:**
- Create: `frontend/src/components/ThemeToggle.vue`

- [ ] **Step 1: 创建 ThemeToggle 组件**

```vue
<!-- frontend/src/components/ThemeToggle.vue -->
<template>
  <button
    class="theme-toggle"
    :title="isDark ? '切换到亮色主题' : '切换到暗色主题'"
    @click="themeStore.toggleTheme()"
  >
    <span class="theme-toggle-icon">{{ isDark ? '🌙' : '☀️' }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.theme === 'dark')
</script>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--duration-fast) ease;
  -webkit-tap-highlight-color: transparent;
}

.theme-toggle:hover {
  background: var(--sidebar-hover-bg);
}

.theme-toggle-icon {
  font-size: 18px;
  line-height: 1;
}
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/ThemeToggle.vue
git commit -m "feat: add ThemeToggle component"
```

---

### Task 5: 改造 MainLayout — 侧边栏/顶栏/TabBar 主题适配

**Files:**
- Modify: `frontend/src/layouts/MainLayout.vue` (全文件改造)

- [ ] **Step 1: 更新 MainLayout template — 侧边栏菜单动态 theme**

将 template 中的 `<a-menu>` 部分（第 19-61 行）的 `theme="dark"` 改为动态：

```html
<a-menu
  v-model:selectedKeys="selectedKeys"
  mode="inline"
  :theme="isDarkTheme ? 'dark' : 'light'"
  class="sidebar-menu"
  @click="onMenuClick"
>
```

在 `<script setup>` 中添加 import 和 computed：

```typescript
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const isDarkTheme = computed(() => themeStore.theme === 'dark')
```

在 topbar 右侧区域的 dropdown 之前添加 ThemeToggle：

```html
<div class="topbar-right">
  <ThemeToggle />
  <a-dropdown>
    <!-- ... existing user trigger ... -->
  </a-dropdown>
</div>
```

添加 import：

```typescript
import ThemeToggle from '@/components/ThemeToggle.vue'
```

- [ ] **Step 2: 重写 MainLayout styles — 全面使用 CSS 变量**

将 `<style scoped>` 部分替换为以下（完整替换）：

```css
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
  height: 56px;
  padding: 0 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--sidebar-border);
  overflow: hidden;
  white-space: nowrap;
}

.sidebar-logo-icon {
  font-size: 20px;
  margin-right: 8px;
  flex-shrink: 0;
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
  display: inline-block;
  width: 20px;
  text-align: center;
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
  font-size: 20px;
  line-height: 1;
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
```

- [ ] **Step 3: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add frontend/src/layouts/MainLayout.vue
git commit -m "feat: adapt MainLayout sidebar/topbar/tabbar to dual theme"
```

---

### Task 6: 改造 AuthLayout — 双主题认证页背景

**Files:**
- Modify: `frontend/src/layouts/AuthLayout.vue` (全文件改造)

- [ ] **Step 1: 重写 AuthLayout — 支持双主题背景**

将整个文件替换为：

```vue
<!-- frontend/src/layouts/AuthLayout.vue -->
<template>
  <div class="auth-layout" :class="isDark ? 'auth-dark' : 'auth-light'">
    <!-- 暗色主题：星空背景 -->
    <template v-if="isDark">
      <div class="stars" ref="starsRef"></div>
      <div class="moon">
        <div class="moon-crater moon-crater-1"></div>
        <div class="moon-crater moon-crater-2"></div>
        <div class="moon-crater moon-crater-3"></div>
      </div>
      <div class="cloud cloud-1"></div>
      <div class="cloud cloud-2"></div>
      <div class="landscape">
        <svg class="landscape-svg" viewBox="0 0 1440 200" preserveAspectRatio="none">
          <path d="M0 200 L0 160 Q100 120 200 140 Q280 90 350 100 Q400 60 440 80 L460 75 L462 30 L468 30 L470 75 Q480 78 500 90 Q560 60 640 80 Q720 100 800 90 Q880 70 960 85 Q1040 100 1100 75 Q1150 50 1200 80 Q1300 110 1400 95 L1440 100 L1440 200Z" fill="#0B1120"/>
          <rect x="445" y="48" width="30" height="27" rx="2" fill="#162040"/>
          <path d="M440 50 L460 30 L480 50" fill="none" stroke="#1A2340" stroke-width="2"/>
          <rect x="452" y="55" width="8" height="8" rx="1" fill="rgba(229,168,75,0.4)"/>
          <rect x="464" y="55" width="8" height="8" rx="1" fill="rgba(229,168,75,0.3)"/>
          <circle cx="420" cy="72" r="12" fill="#0F1729"/>
          <circle cx="495" cy="68" r="10" fill="#0F1729"/>
          <circle cx="410" cy="78" r="8" fill="#0D1425"/>
        </svg>
      </div>
      <div class="lighthouse-glow"></div>
      <div class="water">
        <div class="water-line" style="bottom: 55px; --dur: 5s;"></div>
        <div class="water-line" style="bottom: 40px; --dur: 6s;"></div>
        <div class="water-line" style="bottom: 25px; --dur: 7s;"></div>
        <div class="water-line" style="bottom: 12px; --dur: 5.5s;"></div>
      </div>
    </template>

    <!-- 亮色主题：有机色块背景 -->
    <template v-else>
      <div class="bg-shapes">
        <div class="blob blob-1"></div>
        <div class="blob blob-2"></div>
        <div class="blob blob-3"></div>
        <div class="blob blob-4"></div>
      </div>
      <div class="bg-grid"></div>
      <div class="float-deco float-circle fc-1"></div>
      <div class="float-deco float-circle fc-2"></div>
      <div class="float-deco fc-3"></div>
      <div class="float-deco float-circle fc-4"></div>
      <div class="float-deco fc-5"></div>
    </template>

    <div class="auth-container" :class="{ 'auth-container-split': showBrand && !isMobile }">
      <!-- 亮色主题大屏：左侧品牌介绍 -->
      <div v-if="showBrand && !isMobile" class="brand-section">
        <div class="brand-badge">
          <span class="brand-badge-dot"></span>
          家庭数字中心
        </div>
        <div class="brand-heading">暖 屿</div>
        <div class="brand-heading-en">Warm Isle</div>
        <div class="brand-desc">
          全家人的记账、待办、愿望、交流。<br>
          一个温暖的地方，搞定所有家庭琐事。
        </div>
        <div class="brand-features">
          <div class="feature-item"><span class="feature-icon">💰</span><span>记账</span></div>
          <div class="feature-item"><span class="feature-icon">✅</span><span>待办</span></div>
          <div class="feature-item"><span class="feature-icon">💫</span><span>愿望</span></div>
          <div class="feature-item"><span class="feature-icon">💬</span><span>论坛</span></div>
        </div>
      </div>

      <!-- 右侧：卡片区 -->
      <div class="auth-card-wrapper">
        <div v-if="!showBrand || isMobile" class="auth-header-compact">
          <span class="auth-logo-icon">🏠</span>
          <h1 class="auth-title-compact">暖屿 · WarmIsle</h1>
        </div>
        <div class="auth-card">
          <slot />
        </div>
        <div class="auth-footer">
          <p>&copy; 2026 暖屿 · WarmIsle</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const isDark = computed(() => themeStore.theme === 'dark')

// 亮色主题在大屏显示品牌介绍
const showBrand = ref(window.innerWidth >= 860)
const isMobile = ref(window.innerWidth < 520)

function onResize() {
  showBrand.value = window.innerWidth >= 860
  isMobile.value = window.innerWidth < 520
}
onMounted(() => {
  window.addEventListener('resize', onResize)
  // 暗色主题生成星星
  if (isDark.value && starsRef.value) {
    generateStars()
  }
})

const starsRef = ref<HTMLElement | null>(null)

function generateStars() {
  if (!starsRef.value) return
  for (let i = 0; i < 80; i++) {
    const star = document.createElement('div')
    star.className = 'star'
    star.style.left = Math.random() * 100 + '%'
    star.style.top = Math.random() * 70 + '%'
    star.style.setProperty('--dur', (2 + Math.random() * 4) + 's')
    star.style.setProperty('--max-opacity', String(0.3 + Math.random() * 0.7))
    star.style.animationDelay = Math.random() * 5 + 's'
    const size = (1 + Math.random() * 2) + 'px'
    star.style.width = size
    star.style.height = size
    starsRef.value.appendChild(star)
  }
}
</script>

<style scoped>
/* ==================== 基础布局 ==================== */
.auth-layout {
  position: relative;
  width: 100vw;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  padding: 16px;
}

.auth-light {
  background: var(--color-bg-layout);
}

.auth-dark {
  background: linear-gradient(180deg, #0B1120 0%, #0F1729 40%, #162040 100%);
}

/* ==================== 亮色主题：有机色块背景 ==================== */
.bg-shapes {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0;
  animation: blobIn 2s ease-out forwards;
}

.blob-1 { width: 500px; height: 500px; top: -15%; left: -10%; background: radial-gradient(circle, #FDDAC6 0%, transparent 70%); animation-delay: 0.2s; }
.blob-2 { width: 400px; height: 400px; top: 10%; right: -8%; background: radial-gradient(circle, #B8D4E8 0%, transparent 70%); animation-delay: 0.4s; }
.blob-3 { width: 350px; height: 350px; bottom: -10%; left: 20%; background: radial-gradient(circle, #A8D8CA 0%, transparent 70%); animation-delay: 0.6s; }
.blob-4 { width: 300px; height: 300px; bottom: 5%; right: 10%; background: radial-gradient(circle, #C5B3D9 0%, transparent 70%); animation-delay: 0.8s; }

@keyframes blobIn {
  from { opacity: 0; transform: scale(0.7); }
  to { opacity: 0.5; transform: scale(1); }
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(61, 53, 48, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(61, 53, 48, 0.02) 1px, transparent 1px);
  background-size: 48px 48px;
  pointer-events: none;
}

.float-deco {
  position: absolute;
  pointer-events: none;
  opacity: 0;
  animation: floatIn 1s ease-out forwards;
}

.float-circle { border-radius: 50%; border: 2px solid; }

.fc-1 { width: 20px; height: 20px; top: 15%; left: 12%; border-color: #F09888; animation-delay: 1.2s; animation-name: floatIn, floatBob1; animation-duration: 1s, 6s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.2s, 2.2s; }
.fc-2 { width: 14px; height: 14px; top: 25%; right: 18%; border-color: #6BBAA7; animation-delay: 1.4s; animation-name: floatIn, floatBob2; animation-duration: 1s, 7s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.4s, 2.4s; }
.fc-3 { width: 10px; height: 10px; bottom: 20%; left: 8%; background: #F7D08A; border-radius: 50%; animation-delay: 1.6s; animation-name: floatIn, floatBob3; animation-duration: 1s, 8s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.6s, 2.6s; }
.fc-4 { width: 16px; height: 16px; bottom: 30%; right: 10%; border-color: #C5B3D9; border-radius: 4px; transform: rotate(45deg); animation-delay: 1.8s; animation-name: floatIn, floatBob1; animation-duration: 1s, 9s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 1.8s, 2.8s; }
.fc-5 { width: 8px; height: 8px; top: 60%; left: 25%; background: #F09888; border-radius: 50%; animation-delay: 2s; animation-name: floatIn, floatBob2; animation-duration: 1s, 5s; animation-timing-function: ease-out, ease-in-out; animation-fill-mode: forwards, none; animation-iteration-count: 1, infinite; animation-delay: 2s, 3s; }

@keyframes floatIn { from { opacity: 0; } to { opacity: 0.6; } }
@keyframes floatBob1 { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-12px); } }
@keyframes floatBob2 { 0%, 100% { transform: translateY(0) rotate(0deg); } 50% { transform: translateY(-8px) rotate(5deg); } }
@keyframes floatBob3 { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-15px); } }

/* ==================== 暗色主题：星空背景 ==================== */
.stars { position: absolute; inset: 0; pointer-events: none; }

.star {
  position: absolute;
  background: #FFF;
  border-radius: 50%;
  animation: twinkle var(--dur) ease-in-out infinite;
  opacity: 0;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.1; transform: scale(0.8); }
  50% { opacity: var(--max-opacity); transform: scale(1.2); }
}

.moon {
  position: absolute;
  top: 8%;
  right: 12%;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: radial-gradient(circle at 35% 35%, #FFF8E7 0%, #F0D890 40%, #D4A850 100%);
  box-shadow: 0 0 40px rgba(240, 216, 144, 0.3), 0 0 80px rgba(240, 216, 144, 0.15);
  opacity: 0;
  animation: moonRise 2s ease-out 0.5s forwards;
}

.moon-crater { position: absolute; border-radius: 50%; background: rgba(180, 150, 80, 0.15); }
.moon-crater-1 { width: 15px; height: 15px; top: 20px; left: 25px; }
.moon-crater-2 { width: 10px; height: 10px; top: 40px; left: 15px; }
.moon-crater-3 { width: 8px; height: 8px; top: 30px; left: 45px; }

@keyframes moonRise { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }

.cloud { position: absolute; background: rgba(255, 255, 255, 0.03); border-radius: 100px; filter: blur(20px); }
.cloud-1 { width: 300px; height: 60px; top: 18%; left: -50px; animation: cloudDrift 60s linear infinite; }
.cloud-2 { width: 200px; height: 40px; top: 30%; right: -30px; animation: cloudDrift 80s linear infinite reverse; }

@keyframes cloudDrift { from { transform: translateX(-100%); } to { transform: translateX(calc(100vw + 100%)); } }

.landscape { position: absolute; bottom: 0; left: 0; width: 100%; height: 200px; pointer-events: none; }
.landscape-svg { position: absolute; bottom: 0; width: 100%; height: 100%; }

.lighthouse-glow {
  position: absolute;
  bottom: 130px;
  left: 50%;
  transform: translateX(-50%);
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(229, 168, 75, 0.15) 0%, transparent 70%);
  animation: glowPulse 4s ease-in-out infinite;
  pointer-events: none;
}

@keyframes glowPulse {
  0%, 100% { opacity: 0.6; transform: translateX(-50%) scale(1); }
  50% { opacity: 1; transform: translateX(-50%) scale(1.1); }
}

.water { position: absolute; bottom: 0; left: 0; width: 100%; height: 80px; background: linear-gradient(180deg, rgba(15, 23, 41, 0) 0%, rgba(20, 30, 55, 0.8) 100%); pointer-events: none; }

.water-line {
  position: absolute;
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, rgba(229, 168, 75, 0.15) 30%, rgba(229, 168, 75, 0.25) 50%, rgba(229, 168, 75, 0.15) 70%, transparent 100%);
  animation: waterShimmer var(--dur) ease-in-out infinite;
}

@keyframes waterShimmer {
  0%, 100% { opacity: 0.3; transform: scaleX(0.8); }
  50% { opacity: 0.8; transform: scaleX(1); }
}

/* ==================== 容器 ==================== */
.auth-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-container-split {
  max-width: 860px;
  flex-direction: row;
  gap: 60px;
  align-items: center;
}

/* ==================== 亮色大屏品牌区 ==================== */
.brand-section {
  flex: 1;
  opacity: 0;
  animation: slideRight 1s cubic-bezier(0.22, 1, 0.36, 1) 0.4s forwards;
}

@keyframes slideRight {
  from { opacity: 0; transform: translateX(-30px); }
  to { opacity: 1; transform: translateX(0); }
}

.brand-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  background: rgba(232, 116, 97, 0.08);
  border: 1px solid rgba(232, 116, 97, 0.15);
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
  color: var(--color-brand);
  margin-bottom: 28px;
  letter-spacing: 0.03em;
}

.brand-badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-brand);
  animation: dotPulse 2s ease-in-out infinite;
}

@keyframes dotPulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.3); }
}

.brand-heading {
  font-size: 48px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
  margin-bottom: 8px;
  letter-spacing: 0.04em;
}

.brand-heading-en {
  font-size: 14px;
  font-weight: 300;
  color: var(--color-text-disabled);
  letter-spacing: 0.3em;
  text-transform: uppercase;
  margin-bottom: 24px;
}

.brand-desc {
  font-size: 16px;
  color: var(--color-text-secondary);
  line-height: 1.8;
  max-width: 340px;
  margin-bottom: 40px;
}

.brand-features {
  display: flex;
  gap: 24px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.feature-icon {
  font-size: 18px;
}

/* ==================== 卡片区 ==================== */
.auth-card-wrapper {
  flex: 0 0 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 400px;
}

.auth-header-compact {
  text-align: center;
  margin-bottom: 24px;
}

.auth-logo-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 8px;
}

.auth-title-compact {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: 2px;
}

[data-theme="dark"] .auth-title-compact {
  color: #E8E4DE;
}

.auth-card {
  width: 100%;
  background: var(--auth-card-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--auth-card-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--auth-card-shadow);
  padding: 36px 32px;
  opacity: 0;
  animation: cardAppear 0.9s cubic-bezier(0.22, 1, 0.36, 1) 0.6s forwards;
}

@keyframes cardAppear {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 暗色主题卡片顶部高光线 */
[data-theme="dark"] .auth-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 3px;
  background: linear-gradient(90deg, transparent, var(--color-brand), transparent);
  border-radius: 0 0 3px 3px;
}

.auth-footer {
  margin-top: 24px;
}

.auth-footer p {
  font-size: 12px;
  color: var(--color-text-disabled);
  margin: 0;
}

[data-theme="dark"] .auth-footer p {
  color: rgba(255, 255, 255, 0.3);
}

/* ==================== 响应式 ==================== */
@media (max-width: 860px) {
  .auth-container-split {
    flex-direction: column;
    gap: 24px;
  }
  .brand-section { text-align: center; }
  .brand-desc { margin-left: auto; margin-right: auto; }
  .brand-features { justify-content: center; }
  .auth-card-wrapper { flex: none; }
}

@media (max-width: 520px) {
  .auth-card { padding: 28px 24px; }
  .moon { width: 50px; height: 50px; top: 5%; right: 8%; }
}
</style>
```

- [ ] **Step 2: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无错误

- [ ] **Step 3: Commit**

```bash
git add frontend/src/layouts/AuthLayout.vue
git commit -m "feat: dual-theme AuthLayout with night sky (dark) and organic blobs (light)"
```

---

### Task 7: 改造 Login 和 Init 页面

**Files:**
- Modify: `frontend/src/views/auth/Login.vue`
- Modify: `frontend/src/views/auth/Init.vue`

- [ ] **Step 1: 更新 Login.vue — 移除硬编码颜色**

Login.vue 的 template 和 script 保持不变，只更新 `<style scoped>` 部分：

```css
<style scoped>
.auth-container {
  width: 100%;
}

.auth-card {
  width: 100%;
  background: transparent;
  border: none;
  box-shadow: none;
  padding: 0;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 24px;
  margin: 0 0 8px;
  color: var(--color-text-primary);
}

[data-theme="dark"] .auth-header h1 {
  color: #E8E4DE;
}

.auth-header p {
  color: var(--color-text-secondary);
  margin: 0;
}
</style>
```

- [ ] **Step 2: 读取 Init.vue 并更新样式**

读取 `frontend/src/views/auth/Init.vue`，将其中所有硬编码颜色（`#f5f5f5`, `#fff`, `#1a1a1a`, `#666` 等）替换为 CSS 变量。修改方式与 Login.vue 相同：移除 `.auth-container` 的 `background`，将 `.auth-card` 的 `background` 改为 `transparent`，文字颜色改为变量。

- [ ] **Step 3: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add frontend/src/views/auth/Login.vue frontend/src/views/auth/Init.vue
git commit -m "feat: adapt Login and Init pages to theme variables"
```

---

### Task 8: 更新 EmptyState 组件

**Files:**
- Modify: `frontend/src/components/EmptyState.vue:39-77`

- [ ] **Step 1: 替换硬编码颜色**

将 EmptyState.vue 的 `<style scoped>` 中硬编码颜色替换为变量：

```css
.empty-description {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0 0 16px 0;
}

.clear-link {
  color: var(--color-brand);
  cursor: pointer;
  font-size: 14px;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  text-decoration: none;
}

.clear-link:hover {
  color: var(--color-brand-hover);
}
```

同时将 SVG 中的硬编码颜色 `#d9d9d9`、`#fafafa`、`#e8e8e8` 替换为使用 `currentColor` 或 CSS 变量（通过 `style` 绑定）：

```html
<svg aria-hidden="true" width="64" height="64" viewBox="0 0 64 64" fill="none">
  <rect x="8" y="12" width="48" height="40" rx="4" stroke="var(--color-border)" stroke-width="2" fill="var(--color-border-secondary)" />
  <line x1="20" y1="24" x2="44" y2="24" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
  <line x1="20" y1="32" x2="38" y2="32" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
  <line x1="20" y1="40" x2="34" y2="40" stroke="var(--color-border)" stroke-width="2" stroke-linecap="round" />
</svg>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/EmptyState.vue
git commit -m "fix: replace hardcoded colors in EmptyState with CSS variables"
```

---

### Task 9: 更新 Dashboard 和所有 View 组件硬编码颜色

**Files:**
- Modify: `frontend/src/views/dashboard/Index.vue:338`
- Modify: `frontend/src/views/ledger/Index.vue` (如有硬编码)
- Modify: `frontend/src/views/todo/Index.vue` (如有硬编码)
- Modify: `frontend/src/views/wish/Index.vue` (如有硬编码)
- Modify: `frontend/src/views/forum/Index.vue:593`
- Modify: `frontend/src/views/forum/TopicDetail.vue:770,809`
- Modify: `frontend/src/views/member/Index.vue` (如有硬编码)
- Modify: `frontend/src/views/category/Index.vue` (如有硬编码)
- Modify: `frontend/src/views/profile/Index.vue` (如有硬编码)

- [ ] **Step 1: 搜索所有视图文件中的硬编码颜色**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend/src && grep -rn "#[0-9a-fA-F]\{3,6\}" views/ --include="*.vue" | grep -v "var(" | grep -v "^\s*//"`
Expected: 列出所有需要替换的硬编码颜色位置

- [ ] **Step 2: 逐个替换硬编码颜色**

对每个文件中的硬编码颜色进行替换：

| 硬编码值 | 替换为 |
|----------|--------|
| `#f5f5f5` (背景) | `var(--color-bg-layout)` |
| `#fff` / `#ffffff` (背景) | `var(--color-bg-container)` |
| `#1a1a1a` / `#262626` / `#333` (文字) | `var(--color-text-primary)` |
| `#666` / `#8c8c8c` / `#999` (次要文字) | `var(--color-text-secondary)` |
| `#d9d9d9` (边框) | `var(--color-border)` |
| `#f0f0f0` / `#fafafa` (次级背景) | `var(--color-border-secondary)` |
| `#1677ff` / `#4096ff` (链接) | `var(--color-brand)` |
| `#FF4D4F` (图表渐变中的危险色) | `var(--color-danger)` |

Dashboard 的图表渐变色（第 338 行）：
```css
/* 旧 */ background: linear-gradient(90deg, #FF4D4F, var(--color-brand));
/* 新 */ background: linear-gradient(90deg, var(--color-danger), var(--color-brand));
```

- [ ] **Step 3: 验证编译**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck 2>&1 | head -20`
Expected: 无错误

- [ ] **Step 4: Commit**

```bash
git add frontend/src/views/
git commit -m "fix: replace all hardcoded colors in view components with CSS variables"
```

---

### Task 10: 验证与测试

- [ ] **Step 1: 运行完整类型检查**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vue-tsc --noEmit --skipLibCheck`
Expected: 无错误

- [ ] **Step 2: 运行现有单元测试**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npx vitest run`
Expected: 所有测试通过

- [ ] **Step 3: 运行构建**

Run: `cd D:/Projects/my_projects/home-center-v1/frontend && npm run build`
Expected: 构建成功

- [ ] **Step 4: Commit（如有修复）**

```bash
git add -A
git commit -m "fix: address type errors and test failures from theme refactor"
```
