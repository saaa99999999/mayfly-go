---
trigger: always_on
---

# 样式与 UI 规范

## Tailwind CSS

优先使用 Tailwind 工具类，支持 `dark:` 前缀：

```vue
<template>
    <div class="flex items-center justify-between p-4 bg-white dark:bg-gray-800">
        <span class="text-sm text-gray-600 dark:text-gray-300">Label</span>
    </div>
</template>
```

## 权限控制

按钮权限使用 `v-auth` 指令：

```vue
<el-button v-auth="'account:add'" type="primary" @click="onAdd">新增</el-button>
```

## 类型安全

- 避免 `any`，使用可选链 `?.`
- 使用 TypeScript 严格模式

## 边界

- ✅ **Always**: 优先使用 Tailwind CSS
- 🚫 **Never**: 使用固定高度计算，优先用 Flexbox
