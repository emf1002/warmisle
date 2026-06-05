# 前端品质升级设计文档

> 日期：2026-06-06
> 状态：待审批
> 范围：暖屿 V1 全站前端 UI 品质提升

## 目标

将暖屿前端从"功能可用"提升到"视觉精致"，重点解决以下问题：
- 导航和数据展示过度依赖 emoji，跨平台渲染不一致
- Dashboard 统计卡片视觉层次单薄，图表是手写 flex 布局
- 大部分页面缺少骨架屏，只有通用 loading spinner
- 页面 padding、内联样式等存在不统一
- 登录表单与精致的 AuthLayout 容器存在品质落差
- 字体缺乏品牌辨识度

## 设计决策总览

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 图标库 | Lucide Icons | 线性风格温暖友好，2000+ 图标，tree-shakeable |
| 显示字体 | Plus Jakarta Sans | 现代几何，专业干净，2024-25 流行趋势 |
| 正文中文字体 | Noto Sans SC | Google Fonts 免费，覆盖全字重 |
| Dashboard 风格 | 简洁数据风（方案 A） | 数据卡片 + 迷你趋势图 + 渐变柱状图 |
| 加载状态 | 页面定制骨架屏 | 替代通用 spinner，匹配各页面内容形态 |
| 执行策略 | 系统先行 | 先建基础设施，再统一应用，一致性最好 |

---

## 一、图标系统

### 1.1 Icon 包装组件

创建 `src/components/Icon.vue`，统一管理 Lucide 图标引用：

```vue
<template>
  <component :is="iconComponent" :size="size" :color="color" :stroke-width="strokeWidth" />
</template>
```

**Props：**
- `name`: string — Lucide 图标名称（PascalCase）
- `size`: number — 默认 20
- `color`: string — 默认 `currentColor`
- `strokeWidth`: number — 默认 1.5

### 1.2 图标映射

| 场景 | 当前 emoji | Lucide 图标 |
|------|-----------|-------------|
| 仪表盘 | 📊 | LayoutDashboard |
| 记账本 | 💰 | Wallet |
| 待办 | ✅ | ListTodo |
| 愿望 | ⭐ | Star |
| 论坛 | 💬 | MessageSquare |
| 成员 | 👥 | Users |
| 分类 | 📂 | FolderOpen |
| 个人 | 👤 | UserCircle |
| 设置 | ⚙️ | Settings |
| 主题切换 | 🌙/☀️ | Moon / Sun |
| 收入 | 📈 | TrendingUp |
| 支出 | 📉 | TrendingDown |
| 添加 | ➕ | Plus |
| 删除 | 🗑️ | Trash2 |
| 编辑 | ✏️ | Pencil |
| 搜索 | 🔍 | Search |
| 筛选 | 🔽 | Filter |
| 日历 | 📅 | Calendar |
| 金额 | 💲 | DollarSign |
| 标签 | 🏷️ | Tag |
| 评论 | 💬 | MessageCircle |
| 点赞 | ❤️ | Heart |
| 返回 | ← | ArrowLeft |
| 更多 | ⋯ | MoreHorizontal |
| 关闭 | ✕ | X |
| 成功 | ✓ | Check |
| 警告 | ⚠ | AlertTriangle |
| 错误 | ✕ | XCircle |
| 信息 | ℹ | Info |

### 1.3 侧边栏菜单图标样式

```css
.menu-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  color: var(--text-secondary);
  transition: color var(--duration-fast) ease;
}
.menu-item.active .menu-icon,
.menu-item:hover .menu-icon {
  color: var(--color-brand);
}
```

---

## 二、字体系统

### 2.1 字体引入

通过 Google Fonts 引入，在 `index.html` 的 `<head>` 中添加：

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&family=Noto+Sans+SC:wght@400;500;600;700&display=swap" rel="stylesheet">
```

### 2.2 字体变量

在 `themes.css` 中添加：

```css
:root {
  --font-display: 'Plus Jakarta Sans', 'Noto Sans SC', sans-serif;
  --font-body: 'Noto Sans SC', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  --font-mono: 'JetBrains Mono', 'Fira Code', monospace;
}
```

### 2.3 应用规则

| 元素 | 字体 | 字重 |
|------|------|------|
| 页面标题 h1 | var(--font-display) | 700 |
| 区域标题 h2/h3 | var(--font-display) | 600 |
| 金额数字 | var(--font-display) | 700, tabular-nums |
| 统计数值 | var(--font-display) | 700 |
| 正文/段落 | var(--font-body) | 400 |
| 按钮文字 | var(--font-display) | 600 |
| 导航菜单 | var(--font-display) | 500 |
| 标签/徽章 | var(--font-display) | 500 |

---

## 三、Dashboard 视觉升级

### 3.1 统计卡片

**结构：** 图标圆角方块 + 标签 + 金额 + 迷你趋势图

```html
<div class="stat-card">
  <div class="stat-card-header">
    <div class="stat-icon stat-icon--income">
      <Icon name="TrendingUp" :size="18" />
    </div>
    <span class="stat-label">本月收入</span>
  </div>
  <div class="stat-amount amount-income">¥12,580.00</div>
  <MiniSparkline :data="incomeTrend" color="var(--color-income)" />
</div>
```

**视觉规范：**
- 卡片圆角：`var(--radius-lg)` (16px)
- 图标容器：36x36px，圆角 10px，背景色为语义色 10% 透明度
- 金额字体：`var(--font-display)`，28px，700 weight
- 收入色：`#4CAF50`，支出色：`#F44336`，结余色：`#2196F3`（实施时新增 `--color-income` / `--color-expense` / `--color-balance` Token 到 themes.css）
- 迷你图高度：32px，7 个数据点

### 3.2 迷你趋势图组件

创建 `src/components/MiniSparkline.vue`：

- Props：`data: number[]`、`color: string`、`height: number` (默认 32)
- 渲染：CSS flex 布局的迷你柱状图，每根柱子带渐变
- 动画：入场时柱子从底部生长（stagger 50ms）

### 3.3 图表区

**柱状图组件** `src/components/BarChart.vue`：
- Props：`data: {label: string, value: number}[]`、`color: string`、`height: number`
- 特性：渐变柱子（品牌色 → 30% 透明度）、标签在底部、hover 显示具体数值
- 动画：入场时柱子从底部依次生长

**图表卡片：**
- 左侧品牌色装饰条（3px 宽，圆角）
- 标题 + 时间范围切换（周/月）
- 背景：`var(--bg-elevated)`

---

## 四、骨架屏系统

### 4.1 通用骨架组件

创建 `src/components/SkeletonCard.vue`：

```vue
<template>
  <div class="skeleton-card" :style="{ width, height, borderRadius }">
    <div class="skeleton-shimmer" />
  </div>
</template>
```

- shimmer 动画：从左到右的渐变扫光
- 颜色：`var(--bg-container)` → `var(--bg-elevated)` 色差

### 4.2 各页面骨架形态

| 页面 | 骨架内容 |
|------|---------|
| 仪表盘 | 3 个统计卡片骨架（图标方块 + 金额条 + 迷你图条） + 图表区骨架 |
| 待办 | 5 行列表骨架（圆形头像 + 标题条 + 标签条 + 日期条） |
| 愿望 | 2x2 卡片网格骨架（矩形图片区 + 标题条 + 进度条） |
| 论坛 | 3 个帖子卡片骨架（头像 + 标题 + 2 行内容 + 标签） |
| 成员 | 5 行表格骨架（头像 + 名称 + 角色标签 + 操作按钮） |
| 记账本 | 已有，保持现状 |
| 个人中心 | 头像区骨架 + 表单字段骨架 |

### 4.3 过渡动画

骨架屏 → 内容的过渡：
```css
.skeleton-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.skeleton-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
```

---

## 五、页面入场动效

### 5.1 数字 Count-up

Dashboard 统计金额首次加载时，从 0 滚动到目标值：
- 时长：800ms
- 缓动：ease-out
- 使用 `requestAnimationFrame` 实现，不引入第三方库

### 5.2 图表生长

柱状图入场动画：
- 柱子从底部生长到目标高度
- 每根柱子 stagger 80ms
- 时长：400ms per bar
- 缓动：`cubic-bezier(0.34, 1.56, 0.64, 1)`（带微弹效果）

### 5.3 卡片列表 Stagger

已有 `cardSlideUp` 动画，保持现有实现。

---

## 六、移动端底栏升级

### 6.1 结构改造

```html
<div class="tabbar">
  <router-link v-for="tab in tabs" :to="tab.path" class="tabbar-item">
    <div class="tabbar-icon" :class="{ active: isActive(tab) }">
      <Icon :name="tab.icon" :size="22" />
    </div>
    <span class="tabbar-label">{{ tab.label }}</span>
  </router-link>
</div>
```

### 6.2 视觉规范

- 图标尺寸：22px
- 选中态：品牌色图标 + 品牌色背景胶囊（圆角 12px，padding 4px 12px）
- 未选中态：`var(--text-muted)` 色
- 标签字体：11px，`var(--font-display)`，500 weight
- 切换动画：图标弹跳（scale 1 → 1.15 → 1，200ms）

### 6.3 Tab 映射

| Tab | 图标 | 路径 |
|-----|------|------|
| 仪表盘 | LayoutDashboard | / |
| 记账 | Wallet | /ledger |
| 待办 | ListTodo | /todo |
| 论坛 | MessageSquare | /forum |
| 我的 | UserCircle | /profile |

---

## 七、卡片视觉层次

### 7.1 统计卡片

- 背景：`var(--bg-elevated)` + 微妙的品牌色渐变底色（5% 透明度）
- 边框：`1px solid var(--border-light)`
- hover：`translateY(-2px)` + 阴影加深到 `var(--shadow-level-2)`

### 7.2 图表卡片

- 左侧装饰条：3px 宽，品牌色，圆角
- 标题区域与图表区有 `1px` 分隔线

### 7.3 列表卡片

- hover 上浮：`translateY(-1px)`
- 阴影变化：`var(--shadow-level-1)` → `var(--shadow-level-2)`
- 过渡：200ms ease

---

## 八、登录表单升级

### 8.1 输入框图标

```html
<a-input v-model:value="form.username" placeholder="用户名">
  <template #prefix><Icon name="User" :size="16" /></template>
</a-input>
<a-input-password v-model:value="form.password" placeholder="密码">
  <template #prefix><Icon name="Lock" :size="16" /></template>
</a-input-password>
```

### 8.2 毛玻璃效果

登录卡片增加与 AuthLayout 一致的毛玻璃背景：
```css
.login-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}
[data-theme="dark"] .login-card {
  background: rgba(45, 42, 38, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
```

### 8.3 记住我（UI 占位）

在登录按钮上方增加复选框。本次只实现 UI 展示，实际持久化（localStorage 或后端 JWT 延期）需后续单独实现。

```html
<div class="login-options">
  <a-checkbox v-model:checked="rememberMe">记住我</a-checkbox>
</div>
```

---

## 九、样式一致性

### 9.1 页面 padding 统一

| 断点 | 值 | Token |
|------|-----|-------|
| 桌面端 ≥768px | 24px | var(--space-lg) |
| 移动端 <768px | 16px | var(--space-md) |

所有页面统一使用 `var(--space-lg)` / `var(--space-md)`，替换硬编码的 `24px`、`16px`。

### 9.2 内联样式清理

将散落的 `style="margin-left: 8px"` 等内联样式迁移到 scoped CSS，使用 Token 变量。

### 9.3 空状态统一

所有列表/网格页面统一使用 `EmptyState` 组件，替换手动实现的空状态。

### 9.4 filter-row 统一

所有使用筛选栏的页面统一使用 `components.css` 中的 `.filter-row` 样式。

---

## 十、新增组件清单

| 组件 | 文件 | 用途 |
|------|------|------|
| Icon | `src/components/Icon.vue` | Lucide 图标包装 |
| MiniSparkline | `src/components/MiniSparkline.vue` | 迷你趋势柱状图 |
| BarChart | `src/components/BarChart.vue` | 柱状图（Dashboard 图表区） |
| SkeletonCard | `src/components/SkeletonCard.vue` | 通用骨架屏卡片 |

## 十一、修改文件清单

| 文件 | 改动 |
|------|------|
| `index.html` | 添加 Google Fonts 引入 |
| `styles/themes.css` | 添加字体变量 |
| `styles/global.css` | 清理内联样式工具类 |
| `styles/components.css` | 统一 filter-row 样式 |
| `layouts/MainLayout.vue` | 侧边栏/底栏 emoji → Lucide 图标 |
| `layouts/AuthLayout.vue` | 无改动（已足够精致） |
| `views/auth/Login.vue` | 输入框图标 + 毛玻璃 + 记住我 |
| `views/dashboard/Index.vue` | 统计卡片 + 图表全面升级 |
| `views/ledger/Index.vue` | padding 统一 |
| `views/todo/Index.vue` | 骨架屏 + padding 统一 |
| `views/wish/Index.vue` | 骨架屏 + padding 统一 |
| `views/forum/Index.vue` | 骨架屏 + padding 统一 |
| `views/member/Index.vue` | 骨架屏 + padding 统一 |
| `views/profile/Index.vue` | padding 统一 |
| `views/category/Index.vue` | padding 统一 |

## 十二、依赖变更

```bash
# 新增
npm install lucide-vue-next

# 移除（如果确认不再使用 emoji 相关逻辑）
# @ant-design/icons-vue — 保留，Ant Design 组件内部可能依赖
```

## 十三、不包含

以下内容不在本次升级范围内：
- 数据可视化增强（饼图、折线图）— 需要后端聚合 API
- 论坛文件拆分 — 属于代码重构，非 UI 品质
- 论坛点赞功能实现 — 属于功能开发
- 新增页面或路由

---

## 验收标准

1. **图标一致性**：全站无 emoji 图标，所有导航和功能图标使用 Lucide
2. **字体一致性**：标题和数字使用 Plus Jakarta Sans，正文使用 Noto Sans SC
3. **Dashboard 品质**：统计卡片有图标 + 迷你趋势图，图表区有渐变柱状图
4. **骨架屏覆盖**：除记账本外，仪表盘/待办/愿望/论坛/成员/个人中心均有定制骨架屏
5. **页面动效**：Dashboard 金额 count-up，图表柱子生长动画
6. **移动端底栏**：Lucide 图标 + 品牌色胶囊选中态
7. **样式统一**：全站页面 padding 一致，无内联样式，空状态统一使用 EmptyState
8. **登录表单**：输入框有图标前缀，卡片有毛玻璃效果
9. **双主题兼容**：所有改动在 light/dark 主题下均正常显示
10. **响应式**：桌面端/平板端/移动端三档断点下均正常显示
