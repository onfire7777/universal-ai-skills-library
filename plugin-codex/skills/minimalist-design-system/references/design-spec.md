# 设计系统完整规范

## 目录

1. [设计哲学](#设计哲学)
2. [风格基因](#风格基因)
3. [响应式战略](#响应式战略)
4. [动效规范](#动效规范)
5. [无障碍指南](#无障碍指南)
6. [实施检查清单](#实施检查清单)

---

## 设计哲学

### 核心原则

**"以结构求清晰，以细节求性格。"**

这套设计系统在向极简主义基石致敬的同时，热烈拥抱现代Web布局和动态交互。它运作于一种深刻的张力之中：**数量上的极端克制与执行上的绝对自信。**

每一个出现在屏幕上的元素都必须通过"生存审计"证明其存在的必要性。然而，这些被保留下来的元素必须以极具表现力和精密度的笔触来执行。

- **留白即工具**：留白（Whitespace）并非空无一物，它是引导注意力的精密仪器
- **动效即沟通**：所有的动效（Motion）都不是为了装饰，而是为了向用户传递界面的层级、状态和响应逻辑
- **色彩即焦点**：色彩不应散漫分布，而应浓缩为单一、令人振奋的强调色

### 视觉意境 (Visual Vibe)

**专业且前卫，自信且艺术，精炼而富有生命力。**

这种设计风格代表了高科技SaaS产品的严谨精密度与创意设计机构大胆表现力的交汇。

**情感关键词：**
- **自信 (Confident)**：尺寸大胆，色彩鲜艳，动画意图明确
- **精妙 (Sophisticated)**：双字体排版系统、精挑细选的色彩比例
- **呼吸感 (Alive)**：微动画、脉动指示器和浮动元素
- **高端 (Premium)**：慷慨的留白与带色调偏移的阴影
- **当代 (Contemporary)**：渐变文字、玻璃拟态暗示和非对称布局

### 拒绝清单

- ❌ 拒绝冰冷与临床：通过暖色调标题和动效注入人文关怀
- ❌ 拒绝模板化：大胆的细节处理防止"另一个Bootstrap网站"
- ❌ 拒绝扁平与死板：纹理、阴影和深度的运用确保交互的物理感

---

## 风格基因

### 1. 签名级电光渐变

电光蓝（Electric Blue）渐变从 `#0052FF` 向 `#4D7CFF` 的流转是系统的灵魂。

**应用逻辑：**
- 主行动点（CTA）
- 关键词强调
- 特色卡片描边
- 暗示能量的流动，比纯色更具维度感

### 2. 反转对比章节

策略性地翻转配色方案：使用深石板色（Deep Slate）作为背景，配以亮色文字。

**战略意图：**
- 在用户滚动页面时创造戏剧性的节奏变化
- 打破白色背景的沉闷感
- 使数据、指标或Final CTA更有冲击力

### 3. 动态深度与"活体"元素

- **脉动指示器**：Badge内部柔和脉动的小圆点，传达"实时"或"在线"状态
- **正弦波浮动**：英雄区核心元素在缓慢正弦波上轻微摆动，营造失重感
- **极慢旋转环**：装饰性虚线圈以极慢速度（60s+）旋转，暗示系统的精密运行

### 4. 冲突美学：双字体系统

| 字体 | 角色 | 特点 | 用途 |
|------|------|------|------|
| **Calistoga** | Display | 温暖、充满性格、略带俏皮的衬线体 | H1/H2标题，"感性与故事"的载体 |
| **Inter** | UI/Body | 极其冷峻、高度易读、专业感十足 | 正文、标签、UI，"理性与效率"的工具 |

**这种对撞产生的火花**：让品牌在抓住用户注意力的同时，能以清晰、毫无阻碍的方式传递信息。

### 5. 纹理优于扁平化

对抗极简主义常见的"廉价感"：

- **微米级点阵**：深色区域背景加入2-3%透明度的`radial-gradient`点阵网格
- **径向发光**：章节角落放置巨大的、模糊处理的强调色晕光（Blur Circles）
- **分层阴影**：卡片拥有主阴影和环境遮蔽阴影，创造真实的物理悬浮感

---

## 响应式战略

### 缩放原则

拒绝简单的"堆叠"：

| 元素 | 桌面端 | 移动端 |
|------|--------|--------|
| Calistoga标题 | `5.25rem` | `2.75rem` |
| 章节间距 | `py-44` | `py-20` |
| 容器宽度 | `max-w-6xl` | `px-4` |

### 断点策略

```css
/* 移动优先 */
.base { /* 默认移动端样式 */ }

@media (min-width: 768px) {
  /* 平板 */
}

@media (min-width: 1024px) {
  /* 桌面 */
}

@media (min-width: 1280px) {
  /* 大屏 */
}
```

---

## 动效规范

### Framer Motion 配置

```javascript
const transition = {
  duration: 0.7,
  ease: [0.16, 1, 0.3, 1] // 平滑退出曲线
}
```

### 微交互

| 交互 | 效果 |
|------|------|
| 按钮悬停 | 向上平移`2px` + 阴影加深 |
| 按钮点击 | `scale(0.98)` |
| 箭头图标悬停 | 向右平移`4px` |
| 卡片悬停 | 阴影`md`→`xl`，背景渐变发光 |

### 减弱动效

**必须**监听 `prefers-reduced-motion` 媒体查询：

```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

开启时禁用：
- 正弦波浮动
- 背景旋转动画
- 脉动指示器

---

## 无障碍指南

### 对比度控制

- 辅助文本：至少 **4.5:1** 对比度
- 渐变文字：确保在背景上的可读性
- 遵循 **WCAG 2.1 AA** 标准

### 键盘导航

- 完美的DOM结构
- 所有交互元素可通过Tab访问
- 焦点状态清晰可见（`ring-2 ring-offset-2`）

### 语义化

- 使用正确的HTML5语义标签
- ARIA标签用于复杂组件
- 图片必须有alt文本

---

## 实施检查清单

### Tailwind 配置

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      fontFamily: {
        display: ['Calistoga', 'serif'],
        sans: ['Inter', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      colors: {
        background: '#FAFAFA',
        foreground: '#0F172A',
        muted: '#F1F5F9',
        accent: {
          DEFAULT: '#0052FF',
          secondary: '#4D7CFF',
        },
        border: '#E2E8F0',
        card: '#FFFFFF',
      },
      boxShadow: {
        'accent': '0 4px 14px rgba(0,82,255,0.25)',
      },
    },
  },
}
```

### CSS 变量导出

```css
:root {
  --color-background: #FAFAFA;
  --color-foreground: #0F172A;
  --color-muted: #F1F5F9;
  --color-accent: #0052FF;
  --color-accent-secondary: #4D7CFF;
  --color-border: #E2E8F0;
  --color-card: #FFFFFF;
}
```

### 字体加载

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Calistoga&family=Inter:wght@400;500;600&family=JetBrains+Mono&display=swap" rel="stylesheet">
```

---

## 结语

作为专家，你的代码不仅要能运行，更要能作为**活的设计规范**。

在每一步集成中，都要自问：

> "这个实现是否体现了'极简现代主义'的自信与性格？"

开始吧，让代码库通过这套系统焕发新生。
