---
name: design-to-code
description: Implements UI from design mockups (Figma, Sketch, or image) with pixel-accurate layout, responsive behavior, and design tokens. Use when 还原设计图, implementing designs, 切图, 设计稿转代码, or converting mockups to frontend code.
---

# 还原设计图（Design to Code）

将设计稿（Figma/Sketch/图片）高保真还原为前端代码，保证布局、间距、字体、颜色与交互一致。

## 触发场景

- 用户说「还原设计图」「按设计稿实现」「切图」「设计转代码」
- 提供 Figma/Sketch 链接、设计图截图或标注
- 需要实现某个页面/组件的 UI

## 执行流程

### 1. 解析设计稿

- **标注**：尺寸、间距（padding/margin/gap）、字号与行高、颜色（含透明度）、圆角、阴影、边框
- **层级**：组件划分、嵌套关系、可复用模块
- **状态**：默认 / hover / focus / disabled / 加载 / 空态 / 错误
- **响应式**：断点、不同宽度下的布局变化（栅格、换行、隐藏）

### 2. 设计 token 对齐

- 颜色、字号、间距尽量映射到项目已有 CSS 变量或 Tailwind 配置
- 若无现成 token，在实现时用变量命名（如 `--color-primary`），便于后续统一

### 3. 实现优先级

1. **布局**：Flex/Grid 先搭骨架，保证对齐与间距
2. **字体排版**：字体、字号、字重、行高、颜色
3. **视觉**：背景、边框、圆角、阴影
4. **交互状态**：hover/focus/disabled 等
5. **响应式**：断点与弹性布局
6. **动效**：若有设计说明的过渡/动画再补

### 4. 还原度自检

- [ ] 关键尺寸与设计稿一致（可容忍 1–2px 差异）
- [ ] 字体与颜色与设计一致
- [ ] 主要断点下布局合理、无错位
- [ ] 可交互元素有明确状态反馈

## 输出约定

- 使用项目现有技术栈（如 Next.js、Tailwind、SCSS、组件库）
- 组件化：可复用部分拆成组件并命名清晰
- 语义化 HTML + 合理 ARIA（按钮、链接、表单）
- 必要时注明：某处与设计稿差异及原因（如兼容性、可访问性）

## 常用对照

| 设计稿 | 实现方式 |
|--------|----------|
| 8px 栅格 | 间距用 8 的倍数（8/16/24/32） |
| 字体层级 | 对应 heading/body/caption 等语义类或 design token |
| 模糊/毛玻璃 | backdrop-filter + 半透明背景 |
| 多行截断 | line-clamp 或 -webkit-line-clamp |
| 安全区域 | padding 配合 env(safe-area-inset-*) |
