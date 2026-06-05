# Frontend Quality Upgrade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade WarmIsle frontend from "functional" to "polished" — Lucide icons, Plus Jakarta Sans typography, Dashboard visual overhaul, per-page skeleton screens, mobile tabbar upgrade, login form polish, and style consistency.

**Architecture:** System-first approach — build reusable components (Icon, SkeletonCard, MiniSparkline, BarChart) first, then apply them across all views. All styling uses existing CSS custom property system. No new libraries except `lucide-vue-next`.

**Tech Stack:** Vue 3.5 (Composition API), Ant Design Vue 4.2, Lucide Vue Next, CSS Custom Properties, Google Fonts (Plus Jakarta Sans + Noto Sans SC)

---

## File Structure

### New Files
| File | Responsibility |
|------|---------------|
| `frontend/src/components/Icon.vue` | Lucide icon wrapper with name/size/color/strokeWidth props |
| `frontend/src/components/SkeletonCard.vue` | Reusable skeleton placeholder with shimmer animation |
| `frontend/src/components/MiniSparkline.vue` | Mini bar chart for Dashboard stat cards |
| `frontend/src/components/BarChart.vue` | Bar chart component for Dashboard chart section |

### Modified Files
| File | Changes |
|------|---------|
| `frontend/index.html:5` | Add Google Fonts preconnect + stylesheet links |
| `frontend/src/styles/themes.css:1-74,76-153` | Add font variables + income/expense/balance color tokens |
| `frontend/src/styles/global.css:11-12` | Update font-family to use `--font-display` / `--font-body` |
| `frontend/src/styles/components.css` | Add skeleton shimmer animation + chart card styles |
| `frontend/src/App.vue:51` | Update font-family in Ant Design theme config |
| `frontend/src/layouts/MainLayout.vue:26-61,113-126,279-565` | Replace emoji with Lucide icons in sidebar + tabbar |
| `frontend/src/views/dashboard/Index.vue` | Full overhaul: stat cards, charts, skeleton, count-up |
| `frontend/src/views/auth/Login.vue:3-46,81-109` | Add input icons, glassmorphism card, remember-me checkbox |
| `frontend/src/views/todo/Index.vue:35-38` | Replace `<a-spin>` with custom skeleton |
| `frontend/src/views/wish/Index.vue:25-28` | Replace `<a-spin>` with custom skeleton |
| `frontend/src/views/forum/Index.vue:14-17` | Replace `<a-spin>` with custom skeleton |
| `frontend/src/views/member/Index.vue` | Add loading skeleton |
| `frontend/src/views/profile/Index.vue` | Add loading skeleton + standardize padding |
| `frontend/src/views/ledger/Index.vue` | Standardize padding |
| `frontend/src/views/category/Index.vue` | Replace emoji icons with Lucide + standardize padding |

---

## Task Dependency Graph

```
Task 1 (Setup + CSS)
  ├── Task 2 (Icon.vue)
  │     ├── Task 5 (MainLayout icons)
  │     ├── Task 7 (Login form)
  │     └── Task 9 (Category icons)
  ├── Task 3 (SkeletonCard.vue)
  │     ├── Task 6 (Dashboard)
  │     └── Task 8 (Other pages skeletons)
  └── Task 4 (MiniSparkline + BarChart)
        └── Task 6 (Dashboard)
```

Tasks 2, 3, 4 can run in parallel after Task 1.
Tasks 5, 7, 8, 9 can run in parallel after their dependencies.

---

### Task 1: Setup — Dependencies, Fonts, CSS Variables

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/index.html:5`
- Modify: `frontend/src/styles/themes.css:1-6,76-81`
- Modify: `frontend/src/styles/global.css:1-12`
- Modify: `frontend/src/styles/components.css`
- Modify: `frontend/src/App.vue:51`

- [ ] **Step 1: Install lucide-vue-next**

```bash
cd frontend && npm install lucide-vue-next
```

- [ ] **Step 2: Add Google Fonts to index.html**

In `frontend/index.html`, insert after line 5 (after the favicon link):

```html
    <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&family=Noto+Sans+SC:wght@400;500;600;700&display=swap" rel="stylesheet">
```

- [ ] **Step 3: Add font variables to themes.css**

In `frontend/src/styles/themes.css`, add font variables at the top of `[data-theme="light"]` (after line 2, before `--color-brand`):

```css
  /* Typography */
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', sans-serif;
  --font-body: 'Noto Sans SC', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
```

Do the same in `[data-theme="dark"]` (after line 78, before `--color-brand`):

```css
  /* Typography */
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', sans-serif;
  --font-body: 'Noto Sans SC', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
```

- [ ] **Step 4: Add semantic color tokens to themes.css**

In `[data-theme="light"]`, add after the chart colors section (after line 67):

```css
  /* Semantic data colors */
  --color-income: #4CAF50;
  --color-income-bg: rgba(76, 175, 80, 0.1);
  --color-expense: #F44336;
  --color-expense-bg: rgba(244, 67, 54, 0.1);
  --color-balance: #2196F3;
  --color-balance-bg: rgba(33, 150, 243, 0.1);
```

In `[data-theme="dark"]`, add after the chart colors section (after line 142):

```css
  /* Semantic data colors */
  --color-income: #66BB6A;
  --color-income-bg: rgba(102, 187, 106, 0.15);
  --color-expense: #EF5350;
  --color-expense-bg: rgba(239, 83, 80, 0.15);
  --color-balance: #42A5F5;
  --color-balance-bg: rgba(66, 165, 245, 0.15);
```

- [ ] **Step 5: Update global.css font-family**

In `frontend/src/styles/global.css`, replace lines 11-12:

```css
  /* Before: */
  --font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif;

  /* After: */
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', sans-serif;
  --font-body: 'Noto Sans SC', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
```

Update the body style (line 36) to use `--font-body`:

```css
body {
  font-family: var(--font-body);
  /* ... rest unchanged */
}
```

- [ ] **Step 6: Add skeleton shimmer to components.css**

Append to `frontend/src/styles/components.css`:

```css
/* Skeleton shimmer animation */
.skeleton-block {
  background: var(--bg-container);
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
}
.skeleton-block::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    var(--bg-elevated) 50%,
    transparent 100%
  );
  animation: shimmer 1.5s ease-in-out infinite;
}
@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

/* Chart card with left accent bar */
.chart-card {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-light);
  padding: var(--space-lg);
  position: relative;
}
.chart-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 12px;
  bottom: 12px;
  width: 3px;
  background: var(--color-brand);
  border-radius: 0 2px 2px 0;
}

/* Count-up number utility */
.count-up {
  font-variant-numeric: tabular-nums;
  font-family: var(--font-display);
}
```

- [ ] **Step 7: Update App.vue font-family**

In `frontend/src/App.vue`, replace line 51:

```ts
// Before:
fontFamily: '"PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", -apple-system, sans-serif',

// After:
fontFamily: "'Plus Jakarta Sans', 'Noto Sans SC', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', -apple-system, sans-serif",
```

- [ ] **Step 8: Verify dev server starts**

```bash
cd frontend && npm run dev
```

Open http://localhost:3000 — verify no console errors, fonts load (check Network tab for Google Fonts requests).

- [ ] **Step 9: Commit**

```bash
git add frontend/package.json frontend/package-lock.json frontend/index.html \
  frontend/src/styles/themes.css frontend/src/styles/global.css \
  frontend/src/styles/components.css frontend/src/App.vue
git commit -m "feat(ui): add Lucide dependency, Google Fonts, and CSS design tokens"
```

---

### Task 2: Icon Component

**Files:**
- Create: `frontend/src/components/Icon.vue`

- [ ] **Step 1: Create Icon.vue**

Create `frontend/src/components/Icon.vue`:

```vue
<script setup lang="ts">
import { computed } from 'vue'
import {
  LayoutDashboard,
  Wallet,
  ListTodo,
  Star,
  MessageSquare,
  Users,
  FolderOpen,
  UserCircle,
  Settings,
  Moon,
  Sun,
  TrendingUp,
  TrendingDown,
  Plus,
  Trash2,
  Pencil,
  Search,
  Filter,
  Calendar,
  DollarSign,
  Tag,
  MessageCircle,
  Heart,
  ArrowLeft,
  MoreHorizontal,
  X,
  Check,
  AlertTriangle,
  XCircle,
  Info,
  ChevronLeft,
  ChevronRight,
  Eye,
  EyeInvisible,
  Lock,
  User,
  Bell,
  Home,
  Send,
  ThumbsUp,
  Clock,
  FileText,
  PieChart,
  BarChart3,
  Minus,
  type LucideProps,
} from 'lucide-vue-next'

const iconMap: Record<string, any> = {
  LayoutDashboard,
  Wallet,
  ListTodo,
  Star,
  MessageSquare,
  Users,
  FolderOpen,
  UserCircle,
  Settings,
  Moon,
  Sun,
  TrendingUp,
  TrendingDown,
  Plus,
  Trash2,
  Pencil,
  Search,
  Filter,
  Calendar,
  DollarSign,
  Tag,
  MessageCircle,
  Heart,
  ArrowLeft,
  MoreHorizontal,
  X,
  Check,
  AlertTriangle,
  XCircle,
  Info,
  ChevronLeft,
  ChevronRight,
  Eye,
  EyeInvisible,
  Lock,
  User,
  Bell,
  Home,
  Send,
  ThumbsUp,
  Clock,
  FileText,
  PieChart,
  BarChart3,
  Minus,
}

const props = withDefaults(defineProps<{
  name: string
  size?: number
  color?: string
  strokeWidth?: number
}>(), {
  size: 20,
  color: 'currentColor',
  strokeWidth: 1.5,
})

const iconComponent = computed(() => iconMap[props.name])
</script>

<template>
  <component
    v-if="iconComponent"
    :is="iconComponent"
    :size="size"
    :color="color"
    :stroke-width="strokeWidth"
  />
</template>
```

- [ ] **Step 2: Verify Icon renders in dev server**

In `frontend/src/views/dashboard/Index.vue`, temporarily add to the template:

```html
<Icon name="LayoutDashboard" :size="24" />
```

Check dev server — icon should render. Then remove the test line.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/Icon.vue
git commit -m "feat(ui): add Icon wrapper component for Lucide icons"
```

---

### Task 3: SkeletonCard Component

**Files:**
- Create: `frontend/src/components/SkeletonCard.vue`

- [ ] **Step 1: Create SkeletonCard.vue**

Create `frontend/src/components/SkeletonCard.vue`:

```vue
<script setup lang="ts">
withDefaults(defineProps<{
  width?: string
  height?: string
  borderRadius?: string
}>(), {
  width: '100%',
  height: '20px',
  borderRadius: 'var(--radius-md)',
})
</script>

<template>
  <div
    class="skeleton-block"
    :style="{ width, height, borderRadius }"
  />
</template>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/SkeletonCard.vue
git commit -m "feat(ui): add SkeletonCard component with shimmer animation"
```

---

### Task 4: MiniSparkline and BarChart Components

**Files:**
- Create: `frontend/src/components/MiniSparkline.vue`
- Create: `frontend/src/components/BarChart.vue`

- [ ] **Step 1: Create MiniSparkline.vue**

Create `frontend/src/components/MiniSparkline.vue`:

```vue
<script setup lang="ts">
import { onMounted, ref } from 'vue'

const props = withDefaults(defineProps<{
  data: number[]
  color?: string
  height?: number
}>(), {
  color: 'var(--color-brand)',
  height: 32,
})

const visible = ref(false)

onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

const max = computed(() => Math.max(...props.data, 1))

import { computed } from 'vue'
</script>

<template>
  <div class="sparkline" :style="{ height: height + 'px' }">
    <div
      v-for="(val, i) in data"
      :key="i"
      class="sparkline-bar"
      :style="{
        height: visible ? ((val / max) * 100) + '%' : '0%',
        background: color,
        transitionDelay: (i * 50) + 'ms',
      }"
    />
  </div>
</template>

<style scoped>
.sparkline {
  display: flex;
  align-items: flex-end;
  gap: 2px;
}
.sparkline-bar {
  flex: 1;
  border-radius: 2px;
  opacity: 0.6;
  transition: height 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
```

- [ ] **Step 2: Create BarChart.vue**

Create `frontend/src/components/BarChart.vue`:

```vue
<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  data: { label: string; value: number }[]
  color?: string
  height?: number
}>(), {
  color: 'var(--color-brand)',
  height: 160,
})

const visible = ref(false)

onMounted(() => {
  requestAnimationFrame(() => { visible.value = true })
})

const max = computed(() => Math.max(...props.data.map(d => d.value), 1))
</script>

<template>
  <div class="bar-chart" :style="{ height: height + 'px' }">
    <div
      v-for="(item, i) in data"
      :key="i"
      class="bar-col"
    >
      <div class="bar-tooltip">{{ item.value }}</div>
      <div
        class="bar-fill"
        :style="{
          height: visible ? ((item.value / max) * 100) + '%' : '0%',
          transitionDelay: (i * 80) + 'ms',
        }"
      >
        <div class="bar-gradient" :style="{ background: `linear-gradient(180deg, ${color}, ${color}33)` }" />
      </div>
      <span class="bar-label">{{ item.label }}</span>
    </div>
  </div>
</template>

<style scoped>
.bar-chart {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 0 8px;
}
.bar-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}
.bar-fill {
  width: 100%;
  transition: height 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
  border-radius: 6px 6px 0 0;
  overflow: hidden;
}
.bar-gradient {
  width: 100%;
  height: 100%;
}
.bar-label {
  margin-top: 6px;
  font-size: 11px;
  color: var(--text-muted);
  font-family: var(--font-display);
}
.bar-tooltip {
  display: none;
  position: absolute;
  top: -28px;
  background: var(--text-primary);
  color: var(--bg-container);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.bar-col:hover .bar-tooltip {
  display: block;
}
</style>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/MiniSparkline.vue frontend/src/components/BarChart.vue
git commit -m "feat(ui): add MiniSparkline and BarChart components for Dashboard"
```

---

### Task 5: MainLayout — Sidebar & Tabbar Icons

**Files:**
- Modify: `frontend/src/layouts/MainLayout.vue:1-127,279-565`

- [ ] **Step 1: Add Icon import to script**

In `frontend/src/layouts/MainLayout.vue`, add after line 130 (inside `<script setup>`):

```ts
import Icon from '@/components/Icon.vue'
```

- [ ] **Step 2: Replace sidebar menu item icons**

Replace each menu item's `<template #icon>` block. For each item, change from:

```html
<template #icon><span class="menu-icon">EMOJI</span></template>
```

To:

```html
<template #icon><Icon name="IconName" :size="18" /></template>
```

Mapping (lines 26-61):
- Line 27 `🏠` → `<Icon name="LayoutDashboard" :size="18" />`
- Line 31 `💰` → `<Icon name="Wallet" :size="18" />`
- Line 35 `✅` → `<Icon name="ListTodo" :size="18" />`
- Line 39 `💫` → `<Icon name="Star" :size="18" />`
- Line 43 `💬` → `<Icon name="MessageSquare" :size="18" />`
- Line 49 `👥` → `<Icon name="Users" :size="18" />`
- Line 53 `📂` → `<Icon name="FolderOpen" :size="18" />`
- Line 58 `👤` → `<Icon name="UserCircle" :size="18" />`

- [ ] **Step 3: Replace mobile tabbar icons**

In the mobile tabbar section (lines 113-126), replace each emoji span with Icon component:

```html
<!-- Before (example): -->
<span class="tabbar-icon">🏠</span>

<!-- After: -->
<span class="tabbar-icon"><Icon name="LayoutDashboard" :size="22" /></span>
```

Mapping:
- `🏠` → `<Icon name="LayoutDashboard" :size="22" />`
- `💰` → `<Icon name="Wallet" :size="22" />`
- `✅` → `<Icon name="ListTodo" :size="22" />`
- `💬` → `<Icon name="MessageSquare" :size="22" />`
- `👤` → `<Icon name="UserCircle" :size="22" />`

- [ ] **Step 4: Update tabbar styles for icon-based design**

In the `<style scoped>` section, update `.tabbar-icon` to support Lucide icons:

```css
.tabbar-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 12px;
  transition: all var(--duration-fast) ease;
  color: var(--text-muted);
  margin-bottom: 2px;
}
.tabbar-item.router-link-active .tabbar-icon {
  background: var(--color-brand-bg, rgba(232, 116, 97, 0.12));
  color: var(--color-brand);
}
```

Add `--color-brand-bg` token to themes.css light theme:

```css
  --color-brand-bg: rgba(232, 116, 97, 0.12);
```

And dark theme:

```css
  --color-brand-bg: rgba(229, 168, 75, 0.15);
```

- [ ] **Step 5: Replace ThemeToggle emoji**

In MainLayout.vue, find the ThemeToggle component usage. If it uses emoji sun/moon, update `ThemeToggle.vue` to use Icon:

In `frontend/src/components/ThemeToggle.vue`, replace the emoji with:

```html
<Icon :name="isDark ? 'Sun' : 'Moon'" :size="18" />
```

Add the import:

```ts
import Icon from '@/components/Icon.vue'
```

- [ ] **Step 6: Verify in dev server**

Start dev server, check:
- Desktop sidebar shows Lucide icons (not emoji)
- Mobile tabbar shows Lucide icons with brand color capsule on active item
- Theme toggle shows Sun/Moon icons
- Tablet collapsed sidebar shows icons correctly

- [ ] **Step 7: Commit**

```bash
git add frontend/src/layouts/MainLayout.vue frontend/src/components/ThemeToggle.vue \
  frontend/src/styles/themes.css
git commit -m "feat(ui): replace emoji with Lucide icons in sidebar and tabbar"
```

---

### Task 6: Dashboard Visual Overhaul

**Files:**
- Modify: `frontend/src/views/dashboard/Index.vue`

- [ ] **Step 1: Read current dashboard structure**

Read `frontend/src/views/dashboard/Index.vue` fully to understand the current template, script, and style sections before making changes.

- [ ] **Step 2: Add imports to script section**

Add at the top of `<script setup>`:

```ts
import Icon from '@/components/Icon.vue'
import MiniSparkline from '@/components/MiniSparkline.vue'
import BarChart from '@/components/BarChart.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import { ref, onMounted, computed } from 'vue'
```

Add loading state and count-up logic:

```ts
const loading = ref(true)

// Count-up animation for stat amounts
const displayIncome = ref(0)
const displayExpense = ref(0)
const displayBalance = ref(0)

function countUp(target: number, setter: (v: number) => void, duration = 800) {
  const start = performance.now()
  const step = (now: number) => {
    const progress = Math.min((now - start) / duration, 1)
    const eased = 1 - Math.pow(1 - progress, 3) // ease-out cubic
    setter(Math.round(target * eased))
    if (progress < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}
```

In the existing `fetchData` or `onMounted`, after data is loaded:

```ts
loading.value = false
// Trigger count-up after a tick
nextTick(() => {
  countUp(summary.income, v => displayIncome.value = v)
  countUp(summary.expense, v => displayExpense.value = v)
  countUp(summary.balance, v => displayBalance.value = v)
})
```

- [ ] **Step 3: Replace stat cards template**

Replace the existing stat cards section (lines 14-35) with:

```html
<div v-if="loading" class="summary-grid">
  <div v-for="i in 3" :key="i" class="stat-card">
    <div class="stat-card-header">
      <SkeletonCard width="36px" height="36px" borderRadius="10px" />
      <SkeletonCard width="60px" height="14px" />
    </div>
    <SkeletonCard width="120px" height="28px" />
    <SkeletonCard width="100%" height="32px" />
  </div>
</div>

<div v-else class="summary-grid">
  <div class="stat-card stat-card--income">
    <div class="stat-card-header">
      <div class="stat-icon stat-icon--income">
        <Icon name="TrendingUp" :size="18" />
      </div>
      <span class="stat-label">本月收入</span>
    </div>
    <div class="stat-amount amount-income">
      ¥{{ (displayIncome / 100).toFixed(2) }}
    </div>
    <MiniSparkline :data="incomeTrend" color="var(--color-income)" />
  </div>

  <div class="stat-card stat-card--expense">
    <div class="stat-card-header">
      <div class="stat-icon stat-icon--expense">
        <Icon name="TrendingDown" :size="18" />
      </div>
      <span class="stat-label">本月支出</span>
    </div>
    <div class="stat-amount amount-expense">
      ¥{{ (displayExpense / 100).toFixed(2) }}
    </div>
    <MiniSparkline :data="expenseTrend" color="var(--color-expense)" />
  </div>

  <div class="stat-card stat-card--balance">
    <div class="stat-card-header">
      <div class="stat-icon stat-icon--balance">
        <Icon name="Wallet" :size="18" />
      </div>
      <span class="stat-label">本月结余</span>
    </div>
    <div class="stat-amount amount-balance">
      ¥{{ (displayBalance / 100).toFixed(2) }}
    </div>
    <div class="stat-trend">
      <span v-if="balanceChange >= 0" class="trend-up">较上月 +{{ balanceChange }}%</span>
      <span v-else class="trend-down">较上月 {{ balanceChange }}%</span>
    </div>
  </div>
</div>
```

- [ ] **Step 4: Replace chart section template**

Replace the expense chart section with:

```html
<div v-if="loading" class="chart-card">
  <SkeletonCard width="120px" height="18px" />
  <SkeletonCard width="100%" height="160px" style="margin-top: 16px;" />
</div>

<div v-else class="chart-card">
  <div class="chart-header">
    <span class="chart-title">近 7 天支出趋势</span>
    <div class="chart-tabs">
      <button
        :class="['chart-tab', { active: chartRange === 'week' }]"
        @click="chartRange = 'week'"
      >周</button>
      <button
        :class="['chart-tab', { active: chartRange === 'month' }]"
        @click="chartRange = 'month'"
      >月</button>
    </div>
  </div>
  <BarChart :data="chartData" color="var(--color-brand)" />
</div>
```

- [ ] **Step 5: Add stat card styles**

Add to the `<style scoped>` section:

```css
.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-md);
  margin-bottom: var(--space-lg);
}

.stat-card {
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  padding: var(--space-lg);
  border: 1px solid var(--border-light);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-level-2);
}

.stat-card-header {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  margin-bottom: var(--space-sm);
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-icon--income {
  background: var(--color-income-bg);
  color: var(--color-income);
}
.stat-icon--expense {
  background: var(--color-expense-bg);
  color: var(--color-expense);
}
.stat-icon--balance {
  background: var(--color-balance-bg);
  color: var(--color-balance);
}

.stat-label {
  font-size: 13px;
  color: var(--text-secondary);
  font-family: var(--font-display);
}

.stat-amount {
  font-size: 28px;
  font-weight: 700;
  font-family: var(--font-display);
  font-variant-numeric: tabular-nums;
  margin-bottom: var(--space-xs);
}
.amount-income { color: var(--color-income); }
.amount-expense { color: var(--color-expense); }
.amount-balance { color: var(--color-balance); }

.stat-trend {
  font-size: 12px;
  color: var(--text-muted);
}
.trend-up { color: var(--color-income); }
.trend-down { color: var(--color-expense); }

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}
.chart-title {
  font-weight: 600;
  font-family: var(--font-display);
}
.chart-tabs {
  display: flex;
  gap: var(--space-xs);
}
.chart-tab {
  padding: 4px 12px;
  border-radius: 6px;
  border: none;
  background: var(--bg-container);
  color: var(--text-secondary);
  font-size: 12px;
  font-family: var(--font-display);
  cursor: pointer;
  transition: all 0.15s ease;
}
.chart-tab.active {
  background: var(--color-brand);
  color: white;
}

@media (max-width: 767px) {
  .summary-grid {
    grid-template-columns: 1fr;
  }
  .stat-amount {
    font-size: 24px;
  }
}
```

- [ ] **Step 6: Add missing script variables**

Add to script section (these may need to be computed from existing data):

```ts
const chartRange = ref<'week' | 'month'>('week')

// These should be computed from the actual API data
// Adjust based on existing data structure
const incomeTrend = computed(() => {
  // Return 7 data points for the sparkline
  // Use existing data or default to zeros
  return weeklyData.value?.map(d => d.income) || [0, 0, 0, 0, 0, 0, 0]
})

const expenseTrend = computed(() => {
  return weeklyData.value?.map(d => d.expense) || [0, 0, 0, 0, 0, 0, 0]
})

const chartData = computed(() => {
  return weeklyData.value?.map(d => ({
    label: d.dayLabel,
    value: d.expense,
  })) || []
})

const balanceChange = computed(() => {
  // Calculate month-over-month change
  return 0 // Placeholder — implement based on available data
})
```

- [ ] **Step 7: Verify in dev server**

Check Dashboard:
- Loading shows skeleton cards (not blank)
- Stat cards render with icons, amounts count up from 0
- MiniSparkline bars animate in
- BarChart bars animate in with stagger
- Hover on stat cards shows lift effect
- Hover on chart bars shows tooltip
- Mobile: stat cards stack vertically

- [ ] **Step 8: Commit**

```bash
git add frontend/src/views/dashboard/Index.vue
git commit -m "feat(ui): overhaul Dashboard with stat cards, charts, skeleton, and count-up"
```

---

### Task 7: Login Form Upgrade

**Files:**
- Modify: `frontend/src/views/auth/Login.vue`

- [ ] **Step 1: Add Icon import**

In `<script setup>` section, add:

```ts
import Icon from '@/components/Icon.vue'
```

- [ ] **Step 2: Add input icon prefixes**

Replace the username input (around line 15):

```html
<!-- Before: -->
<a-input v-model:value="form.username" placeholder="请输入用户名" />

<!-- After: -->
<a-input v-model:value="form.username" placeholder="请输入用户名">
  <template #prefix><Icon name="User" :size="16" color="var(--text-muted)" /></template>
</a-input>
```

Replace the password input (around line 24):

```html
<!-- Before: -->
<a-input-password v-model:value="form.password" placeholder="请输入密码" />

<!-- After: -->
<a-input-password v-model:value="form.password" placeholder="请输入密码">
  <template #prefix><Icon name="Lock" :size="16" color="var(--text-muted)" /></template>
</a-input-password>
```

- [ ] **Step 3: Add remember-me checkbox**

Before the submit button (around line 32), add:

```html
<div class="login-options">
  <a-checkbox v-model:checked="rememberMe">记住我</a-checkbox>
</div>
```

Add to script:

```ts
const rememberMe = ref(false)
```

- [ ] **Step 4: Add glassmorphism card styles**

Update the `<style scoped>` section. Add to `.auth-card`:

```css
.auth-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: var(--radius-lg);
  padding: var(--space-xl);
  box-shadow: var(--shadow-level-2);
}

:root[data-theme="dark"] .auth-card {
  background: rgba(45, 42, 38, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}
```

- [ ] **Step 5: Verify in dev server**

Check login page:
- Input fields have User/Lock icons as prefixes
- "记住我" checkbox appears above submit button
- Card has glassmorphism effect (translucent + blur)
- Dark mode card looks correct

- [ ] **Step 6: Commit**

```bash
git add frontend/src/views/auth/Login.vue
git commit -m "feat(ui): upgrade login form with input icons, glassmorphism, remember-me"
```

---

### Task 8: Skeleton Screens for Todo, Wish, Forum, Member, Profile

**Files:**
- Modify: `frontend/src/views/todo/Index.vue:35-38`
- Modify: `frontend/src/views/wish/Index.vue:25-28`
- Modify: `frontend/src/views/forum/Index.vue:14-17`
- Modify: `frontend/src/views/member/Index.vue`
- Modify: `frontend/src/views/profile/Index.vue`

- [ ] **Step 1: Add SkeletonCard import to each file**

In each of the 5 files, add to `<script setup>`:

```ts
import SkeletonCard from '@/components/SkeletonCard.vue'
```

- [ ] **Step 2: Replace todo loading state**

In `frontend/src/views/todo/Index.vue`, replace lines 35-38:

```html
<!-- Before: -->
<div v-if="loading" class="loading-state">
  <a-spin />
  <span style="margin-left: 8px">加载中...</span>
</div>

<!-- After: -->
<div v-if="loading" class="skeleton-list">
  <div v-for="i in 5" :key="i" class="skeleton-todo-item">
    <SkeletonCard width="20px" height="20px" borderRadius="4px" />
    <div class="skeleton-todo-content">
      <SkeletonCard width="60%" height="16px" />
      <SkeletonCard width="40%" height="12px" />
    </div>
    <SkeletonCard width="60px" height="24px" borderRadius="6px" />
  </div>
</div>
```

Add styles:

```css
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}
.skeleton-todo-item {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
}
.skeleton-todo-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-xxs);
}
```

- [ ] **Step 3: Replace wish loading state**

In `frontend/src/views/wish/Index.vue`, replace lines 25-28:

```html
<!-- After: -->
<div v-if="loading" class="skeleton-wish-grid">
  <div v-for="i in 4" :key="i" class="skeleton-wish-card">
    <SkeletonCard width="100%" height="20px" />
    <SkeletonCard width="70%" height="16px" />
    <SkeletonCard width="100%" height="12px" />
    <SkeletonCard width="50%" height="12px" />
  </div>
</div>
```

Add styles:

```css
.skeleton-wish-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--space-md);
}
.skeleton-wish-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  padding: var(--space-lg);
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
}
```

- [ ] **Step 4: Replace forum loading state**

In `frontend/src/views/forum/Index.vue`, replace lines 14-17:

```html
<!-- After: -->
<div v-if="loading" class="skeleton-forum">
  <div v-for="i in 3" :key="i" class="skeleton-post-card">
    <div class="skeleton-post-header">
      <SkeletonCard width="36px" height="36px" borderRadius="50%" />
      <SkeletonCard width="80px" height="14px" />
      <SkeletonCard width="50px" height="12px" />
    </div>
    <SkeletonCard width="100%" height="16px" />
    <SkeletonCard width="80%" height="16px" />
    <div class="skeleton-post-footer">
      <SkeletonCard width="40px" height="20px" borderRadius="6px" />
      <SkeletonCard width="40px" height="20px" borderRadius="6px" />
    </div>
  </div>
</div>
```

Add styles:

```css
.skeleton-forum {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}
.skeleton-post-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
  padding: var(--space-lg);
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
}
.skeleton-post-header {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
}
.skeleton-post-footer {
  display: flex;
  gap: var(--space-sm);
}
```

- [ ] **Step 5: Add member loading state**

In `frontend/src/views/member/Index.vue`, add before the `<a-table>` (after the page header):

```html
<div v-if="loading" class="skeleton-member">
  <div v-for="i in 5" :key="i" class="skeleton-member-row">
    <SkeletonCard width="32px" height="32px" borderRadius="50%" />
    <SkeletonCard width="100px" height="16px" />
    <SkeletonCard width="60px" height="24px" borderRadius="6px" />
    <SkeletonCard width="50px" height="24px" borderRadius="6px" />
  </div>
</div>
```

Add a `loading` ref to the script if not present:

```ts
const loading = ref(true)
// Set loading.value = false after data fetch
```

Add styles:

```css
.skeleton-member {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}
.skeleton-member-row {
  display: flex;
  align-items: center;
  gap: var(--space-md);
  padding: var(--space-md);
  background: var(--bg-elevated);
  border-radius: var(--radius-md);
}
```

- [ ] **Step 6: Add profile loading state**

In `frontend/src/views/profile/Index.vue`, wrap the existing card content with a loading check:

```html
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

<a-card v-else ...>
  <!-- existing content -->
</a-card>
```

Add a `loading` ref if not present. Add styles:

```css
.skeleton-profile {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  padding: var(--space-lg);
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
}
.skeleton-profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-sm);
  margin-bottom: var(--space-md);
}
```

- [ ] **Step 7: Standardize padding**

In each of the 5 files, check the page container's padding. Replace any hardcoded `24px` or `16px` with `var(--space-lg)` / `var(--space-md)`. The pattern should be:

```css
.page-name {
  padding: var(--space-lg);
}
@media (max-width: 767px) {
  .page-name {
    padding: var(--space-md);
  }
}
```

- [ ] **Step 8: Verify in dev server**

For each page, verify:
- Loading state shows skeleton shapes (not spinner)
- Skeleton shimmer animation plays
- After data loads, skeleton transitions to content
- Padding is consistent at 24px desktop / 16px mobile

- [ ] **Step 9: Commit**

```bash
git add frontend/src/views/todo/Index.vue frontend/src/views/wish/Index.vue \
  frontend/src/views/forum/Index.vue frontend/src/views/member/Index.vue \
  frontend/src/views/profile/Index.vue
git commit -m "feat(ui): add custom skeleton screens and standardize padding"
```

---

### Task 9: Category Page Icons + Ledger Padding

**Files:**
- Modify: `frontend/src/views/category/Index.vue`
- Modify: `frontend/src/views/ledger/Index.vue`

- [ ] **Step 1: Add Icon import to category**

In `frontend/src/views/category/Index.vue`, add to `<script setup>`:

```ts
import Icon from '@/components/Icon.vue'
```

- [ ] **Step 2: Replace category emoji icons**

Find each `<span class="category-icon">EMOJI</span>` and replace with the appropriate Lucide icon. Categories use emoji for their icon field — since categories are user-configurable with emoji, we need to keep emoji for category icons (they're data, not UI chrome). However, the page header and action buttons should use Lucide:

- "添加分类" button: add `<Icon name="Plus" :size="16" />` before the text
- Edit buttons: replace with `<Icon name="Pencil" :size="14" />`
- Delete buttons: replace with `<Icon name="Trash2" :size="14" />`

- [ ] **Step 3: Standardize category padding**

Update the category page container styles to use `var(--space-lg)` / `var(--space-md)`.

- [ ] **Step 4: Standardize ledger padding**

In `frontend/src/views/ledger/Index.vue`, check and update any hardcoded padding values to use CSS variables.

- [ ] **Step 5: Verify in dev server**

Check category page: edit/delete buttons show Lucide icons, padding consistent.
Check ledger page: padding consistent.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/views/category/Index.vue frontend/src/views/ledger/Index.vue
git commit -m "feat(ui): category page Lucide icons + ledger/category padding standardization"
```

---

### Task 10: Final Integration Verification

**Files:** None (verification only)

- [ ] **Step 1: Full visual audit**

Start dev server and visit every page:

1. **Login** (`/login`) — glassmorphism card, input icons, remember-me checkbox
2. **Dashboard** (`/`) — skeleton → stat cards with icons + sparklines + count-up, bar chart with animation
3. **Ledger** (`/ledger`) — consistent padding, existing skeleton works
4. **Todo** (`/todo`) — custom skeleton, consistent padding
5. **Wish** (`/wish`) — custom skeleton, consistent padding
6. **Forum** (`/forum`) — custom skeleton, consistent padding
7. **Members** (`/members`) — skeleton, consistent padding
8. **Categories** (`/categories`) — Lucide icons on buttons, consistent padding
9. **Profile** (`/profile`) — skeleton, consistent padding

For each page, verify:
- Light theme looks correct
- Dark theme looks correct
- Desktop layout (≥1024px) works
- Tablet layout (768-1023px) works
- Mobile layout (<768px) works
- No emoji icons in navigation/chrome (data emoji in categories/wishes is OK)
- Plus Jakarta Sans loads for headings/numbers
- Noto Sans SC loads for body text
- No console errors

- [ ] **Step 2: Check for remaining inline styles**

```bash
cd frontend/src
grep -rn 'style="' views/ --include="*.vue" | grep -v 'v-bind' | grep -v 'data-testid'
```

Review results — any remaining hardcoded inline styles that bypass the token system should be migrated to scoped CSS.

- [ ] **Step 3: Check for remaining emoji in navigation**

```bash
cd frontend/src
grep -rn '[🏠💰✅💫💬👥📂👤🌙☀️]' layouts/ components/ --include="*.vue"
```

Should return no results (emoji in data/views is acceptable).

- [ ] **Step 4: Final commit if any fixes needed**

```bash
git add -A
git commit -m "fix(ui): final polish and consistency fixes from integration audit"
```
