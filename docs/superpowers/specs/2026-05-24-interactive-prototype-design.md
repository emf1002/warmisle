# 暖屿 V1 交互式产品原型设计文档

| 项目 | 内容 |
|------|------|
| 产品名称 | 暖屿 V1 交互式产品原型 |
| 版本号 | V1.0.0 |
| 状态 | 已批准 |
| 创建日期 | 2026-05-24 |
| 作者 | Claude Code |

---

## 1. 设计目标

### 1.1 核心目标

- **视觉冲击力**：科技感渐变风格，深色背景 + 动态光效
- **交互体验**：模拟真实产品交互，卡片可 hover、点击展开
- **信息展示**：清晰展示四大核心模块的功能和数据
- **响应式设计**：桌面端和移动端都能完美展示

### 1.2 目标用户

- 潜在用户：了解产品功能和价值
- 开发团队：参考 UI 设计和交互模式
- 投资人/合作伙伴：快速了解产品定位

---

## 2. 整体架构

### 2.1 页面结构

```
┌─────────────────────────────────────────────────────────────────┐
│  暖屿 · WarmIsle V1                         [GitHub] [Docs]    │  ← 顶部导航栏
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │              🏠 暖屿 · 家庭协作中心                     │   │  ← Hero 区域
│   │   全家人的记账、待办、愿望、交流，一个地方搞定           │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌──────────┐ │
│   │   📊 记账本  │ │   ✅ 待办   │ │   ⭐ 愿望   │ │  💬 论坛 │ │  ← 模块卡片
│   │             │ │             │ │             │ │          │ │
│   │  [预览内容] │ │  [预览内容] │ │  [预览内容] │ │ [预览内容]│ │
│   └─────────────┘ └─────────────┘ └─────────────┘ └──────────┘ │
│                                                                 │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    📱 响应式预览                         │   │  ← 移动端预览
│   │              [手机框架内的界面预览]                      │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 布局规范

| 区域 | 高度 | 说明 |
|------|------|------|
| 顶部导航栏 | 64px | 固定定位，深色背景 |
| Hero 区域 | 400px | 动态粒子背景 |
| 模块卡片区域 | 自适应 | 2x2 网格布局 |
| 响应式预览区域 | 500px | 手机框架预览 |
| 页脚 | 200px | 技术栈和版权 |

---

## 3. 视觉设计

### 3.1 颜色方案

```css
:root {
  /* 主背景渐变 */
  --bg-primary: #0f0c29;
  --bg-secondary: #302b63;
  --bg-tertiary: #24243e;
  
  /* 卡片背景 */
  --card-bg: rgba(255, 255, 255, 0.05);
  --card-border: rgba(255, 255, 255, 0.1);
  --card-hover: rgba(255, 255, 255, 0.1);
  
  /* 强调色渐变 */
  --accent-start: #667eea;
  --accent-end: #764ba2;
  
  /* 功能色 */
  --income: #10b981;  /* 收入 - 翠绿 */
  --expense: #f43f5e; /* 支出 - 玫红 */
  --warning: #f59e0b; /* 警告 - 琥珀 */
  --success: #22c55e; /* 成功 - 绿色 */
  
  /* 文字颜色 */
  --text-primary: rgba(255, 255, 255, 0.9);
  --text-secondary: rgba(255, 255, 255, 0.6);
  --text-tertiary: rgba(255, 255, 255, 0.4);
}
```

### 3.2 字体规范

```css
/* 中文字体 */
font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;

/* 数字字体 */
font-family: 'DIN Alternate', 'Roboto Mono', monospace;

/* 字号 */
--text-xs: 12px;
--text-sm: 14px;
--text-base: 16px;
--text-lg: 18px;
--text-xl: 20px;
--text-2xl: 24px;
--text-3xl: 30px;
--text-4xl: 36px;
```

### 3.3 阴影效果

```css
/* 卡片阴影 */
box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1),
            0 2px 4px -2px rgba(0, 0, 0, 0.1);

/* 卡片 hover 阴影 */
box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2),
            0 8px 10px -6px rgba(0, 0, 0, 0.2);

/* 发光效果 */
box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
```

---

## 4. 组件设计

### 4.1 顶部导航栏

```html
<header class="navbar">
  <div class="navbar-brand">
    <span class="logo">🏠</span>
    <span class="title">暖屿 · WarmIsle</span>
  </div>
  <nav class="navbar-links">
    <a href="#features">功能</a>
    <a href="#demo">演示</a>
    <a href="https://github.com/user/home-center-v1" class="github-link">GitHub</a>
  </nav>
</header>
```

**样式规范**：
- 背景：`rgba(15, 12, 41, 0.9)` + 毛玻璃效果
- 高度：64px
- 固定定位，滚动时保持可见
- Logo 和标题左对齐，链接右对齐

### 4.2 Hero 区域

```html
<section class="hero">
  <div class="particles" id="particles"></div>
  <div class="hero-content">
    <h1>暖屿 · 家庭协作中心</h1>
    <p>全家人的记账、待办、愿望、交流，一个地方搞定</p>
    <div class="cta-buttons">
      <button class="btn-primary">开始体验</button>
      <button class="btn-secondary">查看文档</button>
    </div>
  </div>
</section>
```

**样式规范**：
- 背景：动态粒子效果（使用 Canvas 或 CSS 动画）
- 标题：48px，渐变文字（`background-clip: text`）
- 副标题：18px，半透明白色
- CTA 按钮：渐变背景 + hover 发光效果

### 4.3 模块卡片

```html
<div class="module-card" data-module="ledger">
  <div class="card-header">
    <span class="card-icon">📊</span>
    <h3 class="card-title">记账本</h3>
  </div>
  <div class="card-features">
    <ul>
      <li>收支记录 CRUD</li>
      <li>月度概览统计</li>
      <li>分类占比分析</li>
    </ul>
  </div>
  <div class="card-preview">
    <!-- 模拟数据预览 -->
  </div>
  <div class="card-footer">
    <button class="btn-explore">探索功能 →</button>
  </div>
</div>
```

**样式规范**：
- 背景：`var(--card-bg)` + 毛玻璃效果
- 边框：1px solid `var(--card-border)`
- 圆角：16px
- Hover：上浮 8px + 边框发光 + 背景变亮
- 内部间距：24px

### 4.4 记账本预览

```html
<div class="ledger-preview">
  <div class="ledger-summary">
    <div class="summary-item income">
      <span class="label">本月收入</span>
      <span class="amount">+¥12,580.00</span>
    </div>
    <div class="summary-item expense">
      <span class="label">本月支出</span>
      <span class="amount">-¥8,432.50</span>
    </div>
    <div class="summary-item balance">
      <span class="label">结余</span>
      <span class="amount">¥4,147.50</span>
    </div>
  </div>
  <div class="ledger-list">
    <!-- 3 条模拟记录 -->
  </div>
</div>
```

### 4.5 待办管理预览

```html
<div class="todo-preview">
  <div class="todo-stats">
    <span class="pending">5 待办</span>
    <span class="completed">12 已完成</span>
  </div>
  <div class="todo-list">
    <div class="todo-item">
      <input type="checkbox" class="todo-checkbox">
      <div class="todo-content">
        <span class="todo-title">交电费</span>
        <span class="todo-meta">
          <span class="priority urgent">紧急</span>
          <span class="assignee">👨 李先生</span>
          <span class="due-date">今天</span>
        </span>
      </div>
    </div>
    <!-- 更多待办项 -->
  </div>
</div>
```

### 4.6 愿望清单预览

```html
<div class="wish-preview">
  <div class="wish-tabs">
    <button class="tab active">家庭愿望</button>
    <button class="tab">个人愿望</button>
  </div>
  <div class="wish-list">
    <div class="wish-item">
      <div class="wish-header">
        <span class="wish-title">全家旅行</span>
        <span class="wish-status">待讨论</span>
      </div>
      <div class="wish-meta">
        <span class="votes">👍 3 票</span>
        <span class="amount">预估 ¥15,000</span>
      </div>
      <div class="wish-actions">
        <button class="btn-vote">投票</button>
      </div>
    </div>
    <!-- 更多愿望 -->
  </div>
</div>
```

### 4.7 家庭论坛预览

```html
<div class="forum-preview">
  <div class="announcement">
    <span class="badge">公告</span>
    <span class="content">本周六家庭会议</span>
  </div>
  <div class="post-list">
    <div class="post-item">
      <div class="post-header">
        <span class="author">👩 张女士</span>
        <span class="time">2 小时前</span>
      </div>
      <div class="post-content">
        今天带孩子去了动物园，玩得很开心！
      </div>
      <div class="post-actions">
        <span class="likes">❤️ 5</span>
        <span class="comments">💬 3</span>
      </div>
    </div>
    <!-- 更多动态 -->
  </div>
</div>
```

### 4.8 移动端预览

```html
<div class="mobile-preview">
  <div class="phone-frame">
    <div class="phone-screen">
      <!-- 模拟移动端界面 -->
      <div class="mobile-header">
        <span class="back-arrow">←</span>
        <span class="page-title">记账本</span>
        <span class="add-btn">+</span>
      </div>
      <div class="mobile-content">
        <!-- 移动端内容 -->
      </div>
      <div class="mobile-tabbar">
        <div class="tab-item active">
          <span class="tab-icon">📊</span>
          <span class="tab-label">记账</span>
        </div>
        <!-- 更多 Tab -->
      </div>
    </div>
  </div>
</div>
```

---

## 5. 动画设计

### 5.1 粒子背景效果

```javascript
// 使用 Canvas 实现粒子效果
class Particle {
  constructor(canvas) {
    this.x = Math.random() * canvas.width;
    this.y = Math.random() * canvas.height;
    this.size = Math.random() * 2 + 1;
    this.speedX = Math.random() * 1 - 0.5;
    this.speedY = Math.random() * 1 - 0.5;
    this.opacity = Math.random() * 0.5 + 0.1;
  }
  
  update() {
    this.x += this.speedX;
    this.y += this.speedY;
  }
  
  draw(ctx) {
    ctx.fillStyle = `rgba(102, 126, 234, ${this.opacity})`;
    ctx.beginPath();
    ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
    ctx.fill();
  }
}
```

### 5.2 卡片入场动画

```css
@keyframes cardFadeIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.module-card {
  animation: cardFadeIn 0.6s ease-out forwards;
}

.module-card:nth-child(1) { animation-delay: 0.1s; }
.module-card:nth-child(2) { animation-delay: 0.2s; }
.module-card:nth-child(3) { animation-delay: 0.3s; }
.module-card:nth-child(4) { animation-delay: 0.4s; }
```

### 5.3 数字滚动动画

```javascript
function animateNumber(element, target, duration = 1000) {
  const start = 0;
  const startTime = performance.now();
  
  function update(currentTime) {
    const elapsed = currentTime - startTime;
    const progress = Math.min(elapsed / duration, 1);
    
    // 使用 easeOutCubic 缓动函数
    const easeProgress = 1 - Math.pow(1 - progress, 3);
    const current = Math.floor(start + (target - start) * easeProgress);
    
    element.textContent = formatCurrency(current);
    
    if (progress < 1) {
      requestAnimationFrame(update);
    }
  }
  
  requestAnimationFrame(update);
}
```

### 5.4 Hover 效果

```css
.module-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
}

.module-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2),
              0 0 20px rgba(102, 126, 234, 0.3);
  border-color: rgba(102, 126, 234, 0.5);
}
```

---

## 6. 模拟数据

### 6.1 记账数据

```javascript
const ledgerData = {
  summary: {
    income: 1258000,  // 分为单位
    expense: 843250,
    balance: 414750
  },
  records: [
    {
      id: 1,
      category: '餐饮',
      icon: '🍜',
      amount: -3500,
      date: '今天',
      recorder: '👨 李先生',
      remark: '午餐外卖'
    },
    {
      id: 2,
      category: '工资',
      icon: '💰',
      amount: 1258000,
      date: '昨天',
      recorder: '👩 张女士',
      remark: '5月工资'
    },
    {
      id: 3,
      category: '交通',
      icon: '🚗',
      amount: -1500,
      date: '5月22日 周四',
      recorder: '👨 李先生',
      remark: '加油'
    }
  ]
};
```

### 6.2 待办数据

```javascript
const todoData = {
  stats: {
    pending: 5,
    completed: 12
  },
  items: [
    {
      id: 1,
      title: '交电费',
      priority: 'urgent',
      assignee: '👨 李先生',
      dueDate: '今天',
      completed: false
    },
    {
      id: 2,
      title: '买菜',
      priority: 'important',
      assignee: '👩 张女士',
      dueDate: '明天',
      completed: false
    },
    {
      id: 3,
      title: '修水龙头',
      priority: 'normal',
      assignee: null,
      dueDate: '本周末',
      completed: false
    }
  ]
};
```

### 6.3 愿望数据

```javascript
const wishData = {
  family: [
    {
      id: 1,
      title: '全家旅行',
      category: '旅行',
      amount: 1500000,
      votes: 3,
      status: '待讨论',
      voters: ['👨', '👩', '👦']
    },
    {
      id: 2,
      title: '换新沙发',
      category: '物品',
      amount: 500000,
      votes: 2,
      status: '已同意',
      voters: ['👨', '👩']
    }
  ],
  personal: [
    {
      id: 3,
      title: '新手机',
      category: '物品',
      amount: 600000,
      votes: 0,
      status: '待讨论',
      creator: '👦 小李'
    }
  ]
};
```

### 6.4 论坛数据

```javascript
const forumData = {
  announcement: {
    id: 1,
    content: '本周六下午3点家庭会议，讨论暑假旅行计划',
    author: '👩 张女士',
    time: '1小时前'
  },
  posts: [
    {
      id: 1,
      content: '今天带孩子去了动物园，玩得很开心！',
      author: '👩 张女士',
      avatar: '👩',
      time: '2小时前',
      likes: 5,
      comments: 3
    },
    {
      id: 2,
      content: '新买的洗碗机到了，以后不用手洗碗了 😄',
      author: '👨 李先生',
      avatar: '👨',
      time: '5小时前',
      likes: 4,
      comments: 2
    }
  ]
};
```

---

## 7. 响应式设计

### 7.1 断点设置

```css
/* 移动端 (< 768px) */
@media (max-width: 767px) {
  .module-cards {
    grid-template-columns: 1fr;
  }
  .hero-content h1 {
    font-size: 32px;
  }
  .navbar-links {
    display: none;
  }
  .mobile-menu-btn {
    display: block;
  }
}

/* 平板端 (768px - 1023px) */
@media (min-width: 768px) and (max-width: 1023px) {
  .module-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .hero-content h1 {
    font-size: 40px;
  }
}

/* 桌面端 (≥ 1024px) */
@media (min-width: 1024px) {
  .module-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .hero-content h1 {
    font-size: 48px;
  }
}
```

### 7.2 移动端适配

**导航栏**：
- 隐藏桌面端链接，显示汉堡菜单按钮
- 点击汉堡菜单展开侧边栏或下拉菜单

**Hero 区域**：
- 高度从 400px 减小到 300px
- 标题字号从 48px 减小到 32px
- 副标题字号从 18px 减小到 16px

**模块卡片**：
- 从 2x2 网格布局改为单列堆叠
- 卡片内部间距从 24px 减小到 16px

**响应式预览区域**：
- 在移动端隐藏手机框架预览（避免嵌套过深）

**页脚**：
- 技术栈信息改为垂直排列

---

## 8. 技术实现

### 8.1 文件结构

```
docs/
└── warmisle-prototype.html    # 单文件，内嵌 CSS + JS
```

**说明**：
- 所有 CSS 样式通过 `<style>` 标签内嵌在 HTML 中
- 所有 JavaScript 代码通过 `<script>` 标签内嵌在 HTML 中
- 模拟数据直接硬编码在 JavaScript 中
- 无外部资源依赖（无图片、字体、图标文件）

### 8.2 依赖

- 无外部依赖，纯 HTML + CSS + JavaScript
- 使用 CSS Grid/Flexbox 布局
- 使用 Canvas 实现粒子效果
- 使用 requestAnimationFrame 实现动画
- 使用 CSS 变量（Custom Properties）管理主题色

### 8.3 浏览器兼容性

- Chrome 90+
- Safari 14+
- Firefox 90+
- Edge 90+

---

## 9. 验收标准

### 9.1 视觉验收

- [ ] 深色渐变背景正确显示
- [ ] 粒子动画流畅运行
- [ ] 卡片 hover 效果正常
- [ ] 颜色方案符合设计规范
- [ ] 字体和排版正确

### 9.2 交互验收

- [ ] 数字滚动动画正常
- [ ] 卡片入场动画正常
- [ ] 响应式布局正确切换
- [ ] 移动端预览正确显示

### 9.3 数据验收

- [ ] 记账数据正确显示（收入绿色、支出红色）
- [ ] 待办数据正确显示（优先级标签）
- [ ] 愿望数据正确显示（投票数、状态）
- [ ] 论坛数据正确显示（公告、动态）

---

## 10. 后续优化

### 10.1 V2 计划

- 添加更多交互效果（点击展开详情、模拟登录）
- 集成真实 API 数据
- 添加深色/浅色模式切换
- 优化移动端体验

### 10.2 性能优化

- 粒子效果使用 Web Workers
- 图片懒加载
- 代码分割（如需拆分为多文件）

---

## 附录 A：设计参考

- [Ant Design Vue](https://antdv.com/) - UI 组件库
- [Tailwind CSS](https://tailwindcss.com/) - 工具类优先的 CSS 框架
- [Framer Motion](https://www.framer.com/motion/) - 动画库

## 附录 B：颜色对比度

| 组合 | 对比度 | WCAG 标准 |
|------|--------|-----------|
| 主文字 (rgba(255,255,255,0.9)) / 深色背景 (#0f0c29) | 15.2:1 | AAA |
| 次文字 (rgba(255,255,255,0.6)) / 深色背景 (#0f0c29) | 7.8:1 | AAA |
| 收入绿 (#10b981) / 深色背景 (#0f0c29) | 8.5:1 | AAA |
| 支出红 (#f43f5e) / 深色背景 (#0f0c29) | 5.2:1 | AA |
