# UI 风格指南 V2.0.0 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将前端从 Ant Design 默认蓝主题全面升级为暖橙品牌色的「家庭感」视觉语言，包含品牌色、布局、卡片设计、动效、无障碍等全方位改进。

**Architecture:** 分层改造 —— 先改全局主题和 CSS 变量基础，再改布局组件（AuthLayout / MainLayout），最后逐页面改造组件样式。动效和无障碍作为独立层叠加。所有改动集中在 CSS 和模板层，不涉及业务逻辑。

**Tech Stack:** Vue 3 + Ant Design Vue 4 + CSS custom properties + CSS animations

**说明:** 本计划假设工程师已阅读 `docs/superpowers/specs/2026-05-16-home-center-ui-style.md`（V2.0.0 风格指南）全文。

---

### Task 1: 全局主题与 CSS 变量基础

**Files:**
- Modify: `frontend/src/App.vue`
- Create: `frontend/src/styles/global.css`

- [ ] **Step 1: 更新 App.vue 中的 Ant Design 主题配置**

将 `frontend/src/App.vue` 中 `themeConfig` 替换为 V2.0.0 指南第 9 节的配置：

```typescript
const themeConfig = {
  token: {
    // 品牌色（暖橙）
    colorPrimary: '#E8734A',
    colorSuccess: '#389E0D',
    colorWarning: '#FAAD14',
    colorError: '#FF4D4F',
    colorInfo: '#E8734A',

    // 圆角
    borderRadius: 8,
    fontSize: 14,

    // 字体（中文优先）
    fontFamily: '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif',

    // 背景
    colorBgLayout: '#FAF9F7',
    colorBgContainer: '#FFFFFF',
  },
  components: {
    Button: {
      borderRadius: 8,
      controlHeight: 44,
    },
    Input: {
      borderRadius: 8,
      controlHeight: 44,
    },
    Select: {
      borderRadius: 8,
      controlHeight: 44,
    },
  },
}
```

- [ ] **Step 2: 创建全局 CSS 变量文件**

创建 `frontend/src/styles/global.css`，定义所有不在 Ant Design 体系内的自定义 Token：

```css
/* ==================== 全局 CSS 变量 ==================== */
:root {
  /* 品牌色 */
  --color-brand: #E8734A;
  --color-brand-light: #FFF2EC;

  /* 语义色 */
  --color-success: #389E0D;
  --color-success-light: #F6FFED;
  --color-warning: #FAAD14;
  --color-warning-light: #FFFBE6;
  --color-danger: #FF4D4F;
  --color-danger-light: #FFF2F0;

  /* 中性色 */
  --color-text-primary: #262626;
  --color-text-secondary: #8C8C8C;
  --color-text-disabled: #BFBFBF;
  --color-muted: #8C8C8C;

  /* 背景色 */
  --color-bg-layout: #FAF9F7;
  --color-bg-container: #FFFFFF;
  --color-bg-sidebar: #2C3E50;

  /* 边框 */
  --color-border: #D9D9D9;
  --color-border-secondary: #F0F0F0;

  /* 阴影层级 */
  --shadow-level-0: none;
  --shadow-level-1: 0 1px 2px rgba(0, 0, 0, 0.06);
  --shadow-level-2: 0 2px 8px rgba(0, 0, 0, 0.08);
  --shadow-level-3: 0 4px 12px rgba(0, 0, 0, 0.12);
  --shadow-level-4: 0 8px 24px rgba(0, 0, 0, 0.15);

  /* 间距 */
  --space-xxs: 4px;
  --space-xs: 8px;
  --space-sm: 12px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;

  /* 字体 */
  --font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif;
  --font-title: 600 24px/1.3 var(--font-family);
  --font-card-title: 600 16px/1.4 var(--font-family);
  --font-body: 400 14px/1.6 var(--font-family);
  --font-small: 400 12px/1.5 var(--font-family);
  --font-amount-large: 700 28px/1.2 var(--font-family);
  --font-amount-normal: 600 16px/1.2 var(--font-family);
  --font-amount-small: 500 12px/1.2 var(--font-family);

  /* 圆角 */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;

  /* z-index */
  --z-base: 0;
  --z-sticky: 10;
  --z-dropdown: 20;
  --z-sticky-header: 30;
  --z-overlay: 40;
  --z-modal: 50;
  --z-message: 60;
  --z-notification: 70;

  /* 动效时长 */
  --duration-fast: 150ms;
  --duration-normal: 200ms;
  --duration-slow: 300ms;
}

/* ==================== 全局重置 ==================== */
*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: var(--font-family);
  font-size: 14px;
  line-height: 1.6;
  color: var(--color-text-primary);
  background: var(--color-bg-layout);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 横屏溢出防护 */
#app {
  max-width: 100vw;
  overflow-x: hidden;
}

/* ==================== 金额数字排版 ==================== */
.amount {
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

.amount-income {
  composes: amount;
  color: var(--color-success);
  font-weight: 600;
}

.amount-expense {
  composes: amount;
  color: var(--color-danger);
  font-weight: 600;
}

/* ==================== 焦点可见（无障碍） ==================== */
:focus-visible {
  outline: 2px solid var(--color-brand);
  outline-offset: 2px;
}

/* ==================== 微交互 ==================== */
/* 按钮 hover 上移效果（仅桌面端） */
@media (hover: hover) {
  .ant-btn:not(.ant-btn-link):not(.ant-btn-text):hover {
    transform: translateY(-1px);
  }
}

/* 卡片 hover 阴影提升（仅桌面端） */
@media (hover: hover) {
  .card-hoverable:hover {
    box-shadow: var(--shadow-level-2);
  }
}

/* 列表项按下反馈 */
.list-item-active:active {
  background-color: var(--color-bg-layout);
}

/* ==================== 金额颜色 ==================== */
.income-color {
  color: var(--color-success);
}

.expense-color {
  color: var(--color-danger);
}

/* ==================== 页面过渡动画 ==================== */
.page-fade-enter-active {
  transition: opacity var(--duration-fast) ease-out;
}

.page-fade-leave-active {
  transition: opacity var(--duration-fast) ease-out;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

/* ==================== 卡片列表 stagger 动画 ==================== */
@keyframes cardSlideUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.card-stagger {
  opacity: 0;
  animation: cardSlideUp var(--duration-slow) ease-out forwards;
}

/* ==================== 点赞心跳动画 ==================== */
@keyframes likeBounce {
  0% { transform: scale(1); }
  50% { transform: scale(1.3); }
  100% { transform: scale(1); }
}

.like-bounce {
  animation: likeBounce var(--duration-fast) ease-out;
}

/* ==================== 待办完成淡出 ==================== */
.todo-complete-enter-active {
  transition: all var(--duration-slow) ease-out;
}

.todo-complete-leave-active {
  transition: all var(--duration-slow) ease-out;
}

.todo-complete-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.todo-complete-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

.todo-complete-move {
  transition: transform var(--duration-slow) ease;
}

/* ==================== 金额变化淡入淡出 ==================== */
.amount-fade-enter-active,
.amount-fade-leave-active {
  transition: opacity var(--duration-normal) ease;
}

.amount-fade-enter-from,
.amount-fade-leave-to {
  opacity: 0;
}

/* ==================== 打印样式预留 ==================== */
@media print {
  .sidebar,
  .tabbar,
  .topbar button,
  .fab {
    display: none !important;
  }
  body {
    background: white;
    color: black;
  }
  .card,
  .ledger-item,
  .todo-item {
    box-shadow: none;
    border-radius: 0;
    border: 1px solid #ccc;
  }
}

/* ==================== 暗色模式预留（V2） ==================== */
/* @media (prefers-color-scheme: dark) {
  :root {
    --color-bg-layout: #141414;
    --color-bg-container: #1F1F1F;
    --color-text-primary: #E8E8E8;
    --color-text-secondary: #A0A0A0;
    --color-text-disabled: #666;
    --color-border: #434343;
    --color-border-secondary: #333;
  }
} */
```

- [ ] **Step 3: 在 main.ts 中引入全局 CSS**

在 `frontend/src/main.ts` 顶部添加：
```typescript
import '@/styles/global.css'
```

- [ ] **Step 4: 验证全局主题生效**

运行 `make dev`，检查：
- 按钮颜色是否变为暖橙 `#E8734A`
- 页面背景是否变为暖白 `#FAF9F7`
- 按钮/输入框圆角是否变为 8px
- 按钮/输入框高度是否变为 44px

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.vue frontend/src/styles/global.css frontend/src/main.ts
git commit -m "feat: update global theme to warm orange brand and add CSS foundation"
```

---

### Task 2: AuthLayout 与登录/初始化页重设计

**Files:**
- Modify: `frontend/src/layouts/AuthLayout.vue`
- Modify: `frontend/src/views/auth/Login.vue`
- Modify: `frontend/src/views/auth/Init.vue`

- [ ] **Step 1: 重写 AuthLayout 样式**

替换 `frontend/src/layouts/AuthLayout.vue` 的 `<style scoped>` 部分：

```css
<style scoped>
.auth-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  min-height: 100dvh;
  background: linear-gradient(135deg, #E8734A 0%, #F0884A 50%, #F5A67A 100%);
  padding: 16px;
}

.auth-container {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-logo {
  font-size: 64px;
  display: block;
  margin-bottom: 12px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.auth-title {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
  letter-spacing: 2px;
}

.auth-subtitle {
  font-size: 14px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.75);
  margin: 0;
}

.auth-content {
  width: 100%;
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-level-4);
  padding: 32px;
}

.auth-footer {
  margin-top: 24px;
}

.auth-footer p {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}
</style>
```

- [ ] **Step 2: 更新 AuthLayout 模板添加副标题**

将 `AuthLayout.vue` 模板中 `<h1>` 后添加副标题：

```html
<div class="auth-header">
  <span class="auth-logo">🏠</span>
  <h1 class="auth-title">家庭数字中心</h1>
  <p class="auth-subtitle">全家人的记账、待办、愿望、交流</p>
</div>
```

- [ ] **Step 3: 更新登录页按钮文案为纯动词**

在 `frontend/src/views/auth/Login.vue` 中将登录按钮文案从"登录"改为全宽品牌色按钮"登 录"（带空格增加视觉宽度）：

找到登录按钮的 `a-button`，确认其已经是 `type="primary"` 且 `html-type="submit"` 且 `block`（全宽）。

如果按钮不是全宽，修改为：
```html
<a-button type="primary" html-type="submit" :loading="loading" block size="large">
  登 录
</a-button>
```

同时确保按钮上方的密码输入框支持回车提交（已有 `@keypress.enter` 则保留）。

- [ ] **Step 4: 更新初始化页标题**

在 `frontend/src/views/auth/Init.vue` 中，将页面标题（如果有的话）从"初始化"改为"创建你的家庭数字中心"。确认表单包含「确认密码」和「姓名」字段（按 PRD 已有）。

- [ ] **Step 5: Commit**

```bash
git add frontend/src/layouts/AuthLayout.vue frontend/src/views/auth/Login.vue frontend/src/views/auth/Init.vue
git commit -m "feat: redesign auth pages with brand gradient and warm styling"
```

---

### Task 3: MainLayout 侧边栏重设计 + 平板折叠

**Files:**
- Modify: `frontend/src/layouts/MainLayout.vue`

- [ ] **Step 1: 更新侧边栏背景色和 Logo 样式**

在 `<style scoped>` 中修改侧边栏相关样式：

```css
/* ==================== 侧边栏 ==================== */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 220px;
  background: var(--color-bg-sidebar);
  display: flex;
  flex-direction: column;
  z-index: var(--z-sticky);
  overflow-y: auto;
  transition: width var(--duration-normal) ease;
}

.sidebar-collapsed {
  width: 64px;
}

.sidebar-collapsed .sidebar-logo-text,
.sidebar-collapsed .ant-menu-item span:not(.menu-icon) {
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
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
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
  color: #fff;
  letter-spacing: 1px;
}

.sidebar-menu {
  flex: 1;
  border-inline-end: none !important;
}

/* 侧边栏菜单选中态：品牌色左侧指示条 */
.sidebar-menu :deep(.ant-menu-item-selected) {
  background-color: rgba(232, 115, 74, 0.2) !important;
}

.menu-icon {
  font-size: 16px;
  display: inline-block;
  width: 20px;
  text-align: center;
}
```

- [ ] **Step 2: 添加平板断点响应式逻辑**

在 `<script setup>` 中修改响应式检测，增加 `isTablet` 状态：

```typescript
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

onMounted(() => {
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})
```

- [ ] **Step 3: 更新侧边栏模板绑定 class**

修改 `<aside>` 标签的 class 绑定：

```html
<aside
  class="sidebar"
  :class="{
    'sidebar-collapsed': sidebarCollapsed && !isMobile,
    'sidebar-hidden': isMobile,
  }"
  @mouseenter="isTablet ? sidebarCollapsed = false : null"
  @mouseleave="isTablet ? sidebarCollapsed = true : null"
>
```

- [ ] **Step 4: 更新主内容区 margin-left 响应侧边栏宽度**

修改 `.main-area` 样式和模板绑定：

```html
<div class="main-area" :class="{ 'main-full': isMobile, 'main-compact': sidebarCollapsed && !isMobile }">
```

```css
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
```

顶部栏同理：
```css
.topbar {
  position: fixed;
  top: 0;
  right: 0;
  left: 220px;
  /* ... */
  transition: left var(--duration-normal) ease;
}
```

模板绑定：
```html
<header
  class="topbar"
  :class="{
    'topbar-mobile': isMobile,
    'topbar-compact': sidebarCollapsed && !isMobile,
  }"
>
```

```css
.topbar-compact {
  left: 64px;
}
```

- [ ] **Step 5: 更新移动端 TabBar 选中色为品牌色**

```css
.tabbar-item-active {
  color: var(--color-brand);
}
```

同时添加 Tab 切换指示条滑动过渡（通过 `a-menu` 或自定义实现，由于当前 TabBar 是纯 div 不是 a-menu，需要改为使用 CSS 伪元素模拟）：

```css
.tabbar-item {
  position: relative;
  /* ... existing styles ... */
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
```

- [ ] **Step 6: 验证**

运行 `make dev`，检查：
- 桌面端（>= 1024px）侧边栏 220px 展开，背景 `#2C3E50`
- 平板端（768px-1023px）侧边栏默认 64px 折叠，hover 展开
- 移动端（< 768px）隐藏侧边栏，显示底部 TabBar
- TabBar 选中项为品牌色

- [ ] **Step 7: Commit**

```bash
git add frontend/src/layouts/MainLayout.vue
git commit -m "feat: redesign sidebar with dark background and collapsible tablet mode"
```

---

### Task 4: 仪表盘 CSS Grid 布局

**Files:**
- Modify: `frontend/src/views/dashboard/Index.vue`

- [ ] **Step 1: 替换统计卡片样式**

将 `.stat-card` 阴影从硬编码改为 CSS 变量，添加金额样式：

```css
.stat-card {
  text-align: center;
  box-shadow: var(--shadow-level-1);
  border-radius: var(--radius-md);
  transition: box-shadow var(--duration-normal) ease;
}

.stat-card:hover {
  box-shadow: var(--shadow-level-2);
}

.stat-card :deep(.ant-statistic-content-value) {
  font-size: 28px !important;
  font-weight: 700 !important;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

.income-prefix {
  color: #389E0D;
  font-size: 14px;
}

.expense-prefix {
  color: #FF4D4F;
  font-size: 14px;
}
```

- [ ] **Step 2: 将概览卡片从 a-row/a-col 改为 CSS Grid**

将顶部三张概览卡片区域改为 CSS Grid：

```html
<div class="summary-grid">
  <a-card :bordered="false" class="stat-card">
    <!-- 收入卡片内容不变 -->
  </a-card>
  <a-card :bordered="false" class="stat-card">
    <!-- 支出卡片内容不变 -->
  </a-card>
  <a-card :bordered="false" class="stat-card">
    <!-- 结余卡片内容不变 -->
  </a-card>
</div>
```

```css
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}
```

- [ ] **Step 3: 将中下部区域改为 CSS Grid 布局**

将饼图+待办区域改为 Grid `2fr 1fr`：

```html
<div class="dashboard-grid">
  <a-card title="支出分类" class="section-card">
    <!-- 饼图内容不变 -->
  </a-card>
  <a-card title="近期待办" class="section-card">
    <!-- 待办列表内容不变 -->
  </a-card>
</div>

<div class="bottom-grid" style="margin-top: 16px">
  <a-card title="愿望动态" class="section-card">
    <!-- 愿望列表内容不变 -->
  </a-card>
  <a-card title="论坛热点" class="section-card">
    <!-- 论坛列表内容不变 -->
  </a-card>
</div>
```

```css
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
```

- [ ] **Step 4: 添加平板/移动端响应式 Grid**

```css
/* 平板端：饼图和待办堆叠 */
@media (max-width: 1023px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

/* 移动端：全部单列 */
@media (max-width: 767px) {
  .dashboard {
    padding: 12px;
  }

  .summary-grid {
    gap: 8px;
  }

  .dashboard-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .stat-card :deep(.ant-statistic-content-value) {
    font-size: 24px !important;
  }

  .stat-card :deep(.ant-statistic-content) {
    font-size: 24px;
  }
}
```

移除旧的 `<a-row :gutter="[16, 16]">` 和 `<a-col>` 包装（概览卡片、饼图区、底部区）。

- [ ] **Step 5: 更新 section-card 阴影**

```css
.section-card {
  height: 100%;
  box-shadow: var(--shadow-level-1);
  border-radius: var(--radius-md);
  transition: box-shadow var(--duration-normal) ease;
}

.section-card:hover {
  box-shadow: var(--shadow-level-2);
}
```

- [ ] **Step 6: 更新图表条颜色为暖橙**

```css
.chart-bar {
  height: 100%;
  background: linear-gradient(90deg, #FF4D4F, #E8734A);
  border-radius: 6px;
  min-width: 2px;
  transition: width var(--duration-slow) ease;
}
```

- [ ] **Step 7: 更新月份选择器为月份切换器样式**

将现有的 `<a-month-picker>` + "月份：" 标签替换为左右箭头月份切换器（符合指南 4.8e）：

```html
<div class="month-switcher">
  <a-button type="text" :disabled="!hasPrevMonth" @click="goPrevMonth" class="month-arrow">
    ◀
  </a-button>
  <span class="month-text">{{ selectedMonth.format('YYYY年M月') }}</span>
  <a-button type="text" @click="goNextMonth" class="month-arrow">
    ▶
  </a-button>
</div>
```

```css
.month-switcher {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  gap: 16px;
}

.month-arrow {
  min-width: 44px;
  min-height: 44px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.month-arrow:disabled {
  color: var(--color-text-disabled);
  cursor: not-allowed;
}

.month-text {
  font-size: 16px;
  font-weight: 600;
  min-width: 120px;
  text-align: center;
}
```

在 script 中添加：
```typescript
const hasPrevMonth = computed(() => {
  // V1 不限制范围，始终可切换
  return true
})

function goPrevMonth() {
  selectedMonth.value = selectedMonth.value.subtract(1, 'month')
  fetchData()
}

function goNextMonth() {
  selectedMonth.value = selectedMonth.value.add(1, 'month')
  fetchData()
}
```

移除旧的 `month-selector` div 和 `onMonthChange`（改用 `goPrevMonth`/`goNextMonth`）。

- [ ] **Step 8: Commit**

```bash
git add frontend/src/views/dashboard/Index.vue
git commit -m "feat: redesign dashboard with CSS Grid layout and month switcher"
```

---

### Task 5: 记账记录卡片 + 月份切换器 + 筛选栏重设计

**Files:**
- Modify: `frontend/src/views/ledger/Index.vue`

- [ ] **Step 1: 重写记录卡片 HTML 结构为指南 4.8b 规范**

将 `.ledger-item` 的 HTML 结构改为：

```html
<div
  v-for="item in group.items"
  :key="item.id"
  class="ledger-item"
  :class="{ clickable: canEdit(item) }"
  @click="onItemClick(item)"
>
  <div class="item-top">
    <span class="item-category">
      <span class="item-cat-icon">{{ item.category.icon }}</span>
      <span class="item-cat-name">{{ item.category.name }}</span>
    </span>
    <span class="item-creator-line">
      <span class="creator-avatar">{{ item.creator.avatar }}</span>
      <span class="creator-name">{{ item.creator.name }}</span>
      <span class="creator-label">记录</span>
    </span>
  </div>
  <div class="item-body">
    <span v-if="item.note" class="item-note">{{ truncate(item.note, 30) }}</span>
  </div>
  <div class="item-bottom">
    <span v-if="item.members && item.members.length > 0" class="item-members">
      关联：<span v-for="m in item.members.slice(0, 4)" :key="m.id" class="member-tag">{{ m.avatar }}</span>
      <span v-if="item.members.length > 4" class="member-more">等 {{ item.members.length }} 人</span>
    </span>
    <span class="item-amount" :class="item.category.type === 'income' ? 'amount-income' : 'amount-expense'">
      {{ item.category.type === 'income' ? '+' : '-' }}¥{{ (item.amount / 100).toFixed(2) }}
    </span>
  </div>
</div>
```

- [ ] **Step 2: 更新记录卡片 CSS**

```css
/* ==================== Ledger Item Card ==================== */
.ledger-item {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  margin-bottom: 8px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow var(--duration-normal) ease;
  min-height: 44px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ledger-item.clickable {
  cursor: pointer;
}

.ledger-item.clickable:hover {
  box-shadow: var(--shadow-level-2);
}

/* Top row: category + creator */
.item-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-category {
  display: flex;
  align-items: center;
  gap: 6px;
}

.item-cat-icon {
  font-size: 18px;
}

.item-cat-name {
  font-size: 14px;
  font-weight: 600;
}

.item-creator-line {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.creator-avatar {
  font-size: 14px;
}

/* Body: note */
.item-body {
  /* note line */
}

.item-note {
  font-size: 14px;
  color: var(--color-text-primary);
  word-break: break-word;
}

/* Bottom: members + amount */
.item-bottom {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.item-members {
  font-size: 12px;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 2px;
}

.member-tag {
  font-size: 14px;
}

.member-more {
  font-size: 11px;
  color: var(--color-muted);
}

.item-amount {
  font-size: 16px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
  flex-shrink: 0;
}
```

- [ ] **Step 3: 更新日期分组标头样式**

```css
.date-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: var(--color-bg-layout);
  border-radius: var(--radius-md);
  margin-bottom: 8px;
}

.date-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.date-total {
  font-size: 14px;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}
```

- [ ] **Step 4: 替换月份选择器为左右箭头切换器**

与 Task 4 中的月份切换器一致：

```html
<div class="month-row">
  <div class="month-switcher">
    <a-button type="text" @click="goPrevMonth" class="month-arrow" :disabled="false">
      ◀
    </a-button>
    <span class="month-text">{{ selectedMonth.format('YYYY年M月') }}</span>
    <a-button type="text" @click="goNextMonth" class="month-arrow">
      ▶
    </a-button>
  </div>
  <a-button type="primary" @click="openCreate()">记一笔</a-button>
</div>
```

```css
.month-switcher {
  display: flex;
  align-items: center;
  gap: 8px;
}

.month-arrow {
  min-width: 44px;
  min-height: 44px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.month-text {
  font-size: 16px;
  font-weight: 600;
  min-width: 110px;
  text-align: center;
}
```

在 script 中添加：
```typescript
function goPrevMonth() {
  selectedMonth.value = selectedMonth.value.subtract(1, 'month')
  fetchLedgers()
}

function goNextMonth() {
  selectedMonth.value = selectedMonth.value.add(1, 'month')
  fetchLedgers()
}
```

移除旧的 `onMonthChange` 方法和 `<a-date-picker>`。

- [ ] **Step 5: 更新筛选栏样式（桌面端）**

```css
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
```

- [ ] **Step 6: 更新金额颜色使用 CSS 变量**

将硬编码的 `#52c41a` / `#ff4d4f` 替换为 CSS 变量：

```css
.income-amount {
  color: var(--color-success);
  font-weight: 600;
  font-size: 16px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}

.expense-amount {
  color: var(--color-danger);
  font-weight: 600;
  font-size: 16px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
}
```

- [ ] **Step 7: 更新移动端 FAB 按钮样式为品牌色**

```css
@media (max-width: 767px) {
  .month-row :deep(.ant-btn-primary) {
    position: fixed;
    bottom: 72px;
    right: 20px;
    z-index: var(--z-overlay);
    width: 56px;
    height: 56px;
    border-radius: 50%;
    font-size: 24px;
    box-shadow: var(--shadow-level-3);
  }
}
```

- [ ] **Step 8: Commit**

```bash
git add frontend/src/views/ledger/Index.vue
git commit -m "feat: redesign ledger record cards with new card layout"
```

---

### Task 6: 待办卡片重设计

**Files:**
- Modify: `frontend/src/views/todo/Index.vue`

- [ ] **Step 1: 重写待办卡片 HTML 结构为指南 4.8c 规范**

将 `.todo-item` 的 HTML 结构改为：

```html
<div
  v-for="todo in todos"
  :key="todo.id"
  class="todo-item"
  :class="{ completed: todo.status === 'completed' }"
>
  <div class="todo-main">
    <a-checkbox
      :checked="todo.status === 'completed'"
      class="todo-checkbox"
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
          <span v-else class="todo-unassigned">未指派</span>
        </span>
        <span v-if="todo.due_date" class="todo-due" :class="{ overdue: isOverdue(todo) }">
          📅 {{ formatDate(todo.due_date) }}
        </span>
      </div>
    </div>
    <div v-if="canEdit(todo)" class="todo-actions">
      <a-button type="link" size="small" @click.stop="openEdit(todo)">编辑</a-button>
      <a-button type="link" size="small" danger @click.stop="confirmDelete(todo)">删除</a-button>
    </div>
  </div>
</div>
```

- [ ] **Step 2: 更新待办卡片 CSS**

```css
.todo-item {
  background: var(--color-bg-container);
  border-radius: var(--radius-md);
  padding: 12px 16px;
  border: 1px solid var(--color-border-secondary);
  box-shadow: var(--shadow-level-1);
  transition: box-shadow var(--duration-normal) ease;
}

.todo-item:hover {
  box-shadow: var(--shadow-level-2);
}

.todo-item.completed {
  background: #f9fafb;
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
```

- [ ] **Step 3: 更新已完成待办勾选动效**

添加待办完成后的淡出动画（通过 CSS transition）：卡片已完成时自动应用 `opacity: 0.75` + 删除线。不需要 Vue `<TransitionGroup>` 因为待办不会移出列表（保持原位置，仅样式变化）。

- [ ] **Step 4: 更新移动端样式**

```css
@media (max-width: 767px) {
  .todo-page {
    padding: 16px;
  }
}
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/views/todo/Index.vue
git commit -m "feat: redesign todo cards with new card layout and completion animation"
```

---

### Task 7: 论坛信息流卡片重设计

**Files:**
- Modify: `frontend/src/views/forum/Index.vue`

- [ ] **Step 1: 重写动态卡片 HTML 结构为指南 4.8d 规范**

替换论坛信息流中的 topic/post 卡片。将 template 中 `v-for="feed in feeds"` 内的列表项改为：

**动态卡片（post 类型）：**
```html
<div v-if="feed.type === 'post'" class="feed-card">
  <div class="feed-header">
    <span class="feed-author">
      <span class="feed-avatar">{{ feed.creator?.avatar || '👤' }}</span>
      <span class="feed-name">{{ feed.creator?.name }}</span>
    </span>
    <span class="feed-time">{{ timeAgo(feed.created_at) }}</span>
  </div>
  <div class="feed-body">
    <p class="feed-content">{{ truncate(feed.content, 200) }}</p>
  </div>
  <div class="feed-actions">
    <span class="feed-action" @click.stop="handleLike(feed)">
      <span :class="{ 'like-bounce': feed._liking }">{{ feed.is_liked ? '❤️' : '🤍' }}</span>
      <span>{{ feed.like_count || 0 }}</span>
    </span>
    <span class="feed-action" @click.stop="goToDetail(feed)">
      💬 {{ feed.comment_count || 0 }}
    </span>
    <a-dropdown v-if="canManage(feed)" trigger="click">
      <span class="feed-action" @click.stop>⋯</span>
      <template #overlay>
        <a-menu @click="({ key }) => handleFeedAction(key, feed)">
          <a-menu-item v-if="feed.type === 'topic' && isAdmin" key="pin">
            {{ feed.is_pinned ? '取消置顶' : '置顶' }}
          </a-menu-item>
          <a-menu-item v-if="canManage(feed)" key="edit">编辑</a-menu-item>
          <a-menu-item v-if="canManage(feed)" key="delete" danger>删除</a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</div>
```

**话题卡片（topic 类型）：**
```html
<div v-else-if="feed.type === 'topic'" class="feed-card topic-card">
  <div class="feed-header" v-if="feed.is_pinned">
    <span class="topic-pin-badge">📌 公告</span>
  </div>
  <h3 class="topic-title" @click="goToDetail(feed)">{{ feed.title }}</h3>
  <p class="feed-content topic-excerpt">{{ truncate(feed.content, 150) }}</p>
  <div class="topic-footer">
    <span v-if="feed.tag" class="topic-tag">#{{ feed.tag }}</span>
    <span class="feed-author topic-author">
      <span class="feed-avatar">{{ feed.creator?.avatar || '👤' }}</span>
      <span class="feed-name">{{ feed.creator?.name }}</span>
      <span class="feed-time">{{ timeAgo(feed.created_at) }}</span>
    </span>
  </div>
  <div class="feed-actions">
    <!-- 点赞评论操作栏同动态卡片 -->
  </div>
</div>
```

- [ ] **Step 2: 新增论坛卡片 CSS**

```css
/* ==================== Feed Card ==================== */
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
  background: #f5f5f5;
  border-radius: var(--radius-sm);
}

.topic-author {
  margin-left: auto;
}

/* Feed Actions */
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
```

- [ ] **Step 3: 添加移动端 FAB 悬浮按钮**

在论坛页面底部添加 FAB 按钮（仅移动端显示）：

```html
<div v-if="isMobile" class="forum-fab" @click="showCreateSheet = true">
  <span class="fab-icon">+</span>
</div>

<!-- 底部 Sheet 选择发帖类型 -->
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
    <div class="sheet-option" @click="openCreatePoll(); showCreateSheet = false">
      <span class="sheet-option-icon">🗳️</span>
      <span class="sheet-option-label">发起投票</span>
    </div>
  </div>
</a-drawer>
```

```css
.forum-fab {
  position: fixed;
  bottom: calc(56px + 16px + env(safe-area-inset-bottom, 0));
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
```

在 script 中添加：
```typescript
const isMobile = ref(window.innerWidth < 768)
const showCreateSheet = ref(false)
```

- [ ] **Step 4: 验证现有置顶公告和分页功能不受影响**

确认置顶话题的黄色背景区（`#fffbe6`）改为使用 CSS 变量 `var(--color-warning-light)`：

```css
.pinned-section {
  background: var(--color-warning-light);
  border-radius: var(--radius-md);
  padding: 12px;
  margin-bottom: 16px;
}
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/views/forum/Index.vue
git commit -m "feat: redesign forum feed cards and add mobile FAB"
```

---

### Task 8: 点赞动画 + 页面过渡 + 列表 stagger 加载

**Files:**
- Modify: `frontend/src/views/forum/Index.vue` (点赞动画)
- Modify: `frontend/src/views/forum/TopicDetail.vue` (点赞动画)
- Modify: `frontend/src/App.vue` (页面过渡)
- Modify: `frontend/src/views/ledger/Index.vue` (stagger 加载)
- Modify: `frontend/src/views/todo/Index.vue` (stagger 加载)
- Modify: `frontend/src/views/forum/Index.vue` (stagger 加载)

- [ ] **Step 1: 添加点赞心跳动画**

在 `forum/Index.vue` 和 `forum/TopicDetail.vue` 中，为点赞点击添加 `like-bounce` class 触发动画：

在 script 中添加：
```typescript
const likingItems = ref<Record<string, boolean>>({})

async function triggerLike(item: any) {
  const key = `${item.type}_${item.id}`
  likingItems.value[key] = true
  setTimeout(() => {
    likingItems.value[key] = false
  }, 150)
  // ... existing like logic
}
```

在模板中：
```html
<span
  class="feed-action"
  :class="{ 'like-bounce': likingItems[`${feed.type}_${feed.id}`] }"
  @click.stop="triggerLike(feed)"
>
```

- [ ] **Step 2: 添加页面切换淡入过渡**

在 `App.vue` 中，为 `<router-view />` 包裹 `<Transition>`：

```html
<template>
  <a-config-provider :theme="themeConfig">
    <component :is="layout">
      <router-view v-slot="{ Component: RouteComponent }">
        <transition name="page-fade" mode="out-in">
          <component :is="RouteComponent" />
        </transition>
      </router-view>
    </component>
  </a-config-provider>
</template>
```

添加样式：
```css
<style>
/* 页面过渡 */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 150ms ease-out;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
```

- [ ] **Step 3: 添加列表卡片 stagger 加载动画**

在 `ledger/Index.vue` 中，为记录卡片添加 stagger 动画：

```html
<div
  v-for="(item, index) in group.items"
  :key="item.id"
  class="ledger-item"
  :style="{ animationDelay: `${index * 50}ms` }"
  :class="{ 'card-stagger': true, ... }"
>
```

在 `todo/Index.vue` 中：
```html
<div
  v-for="(todo, index) in todos"
  :key="todo.id"
  class="todo-item"
  :style="{ animationDelay: `${index * 50}ms` }"
  :class="{ 'card-stagger': true, ... }"
>
```

在 `forum/Index.vue` 中：
```html
<div
  v-for="(feed, index) in feeds"
  :key="`${feed.type}_${feed.id}`"
  class="feed-card"
  :style="{ animationDelay: `${index * 50}ms` }"
  :class="{ 'card-stagger': true }"
>
```

stagger 动画的 CSS 已在 `global.css` 中定义。

- [ ] **Step 4: 验证动效**

运行 `make dev`：
- 点赞心形图标应有 scale 弹跳动画（1→1.3→1）
- 页面切换时应有 150ms 淡入效果
- 列表加载时卡片应从下方逐个滑入

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.vue frontend/src/views/forum/Index.vue frontend/src/views/forum/TopicDetail.vue frontend/src/views/ledger/Index.vue frontend/src/views/todo/Index.vue
git commit -m "feat: add like animation, page transitions, and card stagger loading"
```

---

### Task 9: 无障碍改进

**Files:**
- Modify: `frontend/src/layouts/MainLayout.vue`
- Modify: `frontend/src/views/ledger/Index.vue`
- Modify: `frontend/src/views/todo/Index.vue`
- Modify: `frontend/src/views/forum/Index.vue`
- Modify: `frontend/src/components/EmptyState.vue`

- [ ] **Step 1: 为侧边栏和 TabBar 添加语义标签**

在 `MainLayout.vue` 中将 `<aside>` 改为添加 `aria-label`：
```html
<aside class="sidebar" aria-label="主导航">
```

将 `<nav class="tabbar">` 改为添加 `aria-label`：
```html
<nav v-if="isMobile" class="tabbar" aria-label="底部导航">
```

为菜单项添加 `role` 和 `aria-label`。反正是 emoji 图标，ant-design-vue 的 a-menu-item 已有基本无障碍支持。

- [ ] **Step 2: 为 emoji 头像添加 aria-label**

在 `ledger/Index.vue` 中的头像添加 `aria-label`：

```html
<span class="creator-avatar" :aria-label="`${item.creator.name}的头像`">{{ item.creator.avatar }}</span>
```

在所有页面的成员头像处批量添加。

- [ ] **Step 3: 为图标按钮添加 aria-label**

在 `ledger/Index.vue` 的月份切换器中：
```html
<a-button type="text" aria-label="上个月" @click="goPrevMonth" class="month-arrow">◀</a-button>
<a-button type="text" aria-label="下个月" @click="goNextMonth" class="month-arrow">▶</a-button>
```

- [ ] **Step 4: 更新 EmptyState 组件装饰 SVG 添加 aria-hidden**

```html
<svg aria-hidden="true" ...>
```

- [ ] **Step 5: 确认焦点环样式**

`global.css` 中已定义 `:focus-visible` 样式（2px 品牌色 outline）。

- [ ] **Step 6: Commit**

```bash
git add frontend/src/layouts/MainLayout.vue frontend/src/views/ledger/Index.vue frontend/src/views/todo/Index.vue frontend/src/views/forum/Index.vue frontend/src/components/EmptyState.vue
git commit -m "feat: add accessibility improvements (aria labels, focus outlines, semantic HTML)"
```

---

### Task 10: 其他页面样式同步与最终验证

**Files:**
- Modify: `frontend/src/views/wish/Index.vue`
- Modify: `frontend/src/views/profile/Index.vue`
- Modify: `frontend/src/views/member/Index.vue`
- Modify: `frontend/src/views/category/Index.vue`
- Modify: `frontend/src/views/forum/TopicDetail.vue`

- [ ] **Step 1: 更新其他页面中硬编码的颜色值为 CSS 变量**

在所有页面中进行以下替换：
- `#52c41a` → `var(--color-success)` 或类名 `income-color`
- `#ff4d4f` → `var(--color-danger)` 或类名 `expense-color`
- `#1677ff` / `#1890ff` → `var(--color-brand)`
- `#f5f5f5` → `var(--color-bg-layout)`
- `#f0f0f0` → `var(--color-border-secondary)`
- `#999` / `#666` → `var(--color-text-secondary)`
- `#333` / `#1a1a1a` / `#262626` → `var(--color-text-primary)`
- `box-shadow: 0 1px 2px rgba(0,0,0,0.06)` → `var(--shadow-level-1)`
- `box-shadow: 0 2px 8px rgba(0,0,0,0.08)` → `var(--shadow-level-2)`

同时在 wish 页面添加投票截止状态的视觉规范（`opacity: 0.7` 遮罩）。

- [ ] **Step 2: 更新愿望卡片样式**

在 `wish/Index.vue` 中确保卡片使用 CSS 变量：
- 卡片阴影：`var(--shadow-level-1)`，hover 升级为 `var(--shadow-level-2)`
- 投票截止卡片添加 `.vote-ended` class：`opacity: 0.7; pointer-events: none;`
- 进度条颜色使用品牌色

- [ ] **Step 3: 更新个人中心页面**

在 `profile/Index.vue` 中：
- 按钮使用品牌色
- 角色标签使用品牌色
- 列表项最小高度 44px

- [ ] **Step 4: 构建验证**

```bash
make build
```

确认：
- 构建成功，无 CSS 编译错误
- 单二进制正常启动
- 所有页面布局正常

- [ ] **Step 5: Commit**

```bash
git add frontend/src/views/wish/Index.vue frontend/src/views/profile/Index.vue frontend/src/views/member/Index.vue frontend/src/views/category/Index.vue frontend/src/views/forum/TopicDetail.vue
git commit -m "feat: sync remaining pages with new color tokens and card styles"
```

---

## 实施顺序

```
Task 1 (全局主题)
  └─> Task 2 (AuthLayout)
  └─> Task 3 (MainLayout)
        └─> Task 4 (Dashboard)
        └─> Task 5 (Ledger)
        └─> Task 6 (Todo)
        └─> Task 7 (Forum)
              └─> Task 8 (Animations)
              └─> Task 9 (Accessibility)
                    └─> Task 10 (Other pages sync)
```

Task 1 必须先完成（所有后续任务依赖 CSS 变量）。Task 2-3 可以与 1 部分并行。Task 4-7 依赖 Task 1 和 3。Task 8-9 可以与 4-7 并行。Task 10 在最后统一处理。
