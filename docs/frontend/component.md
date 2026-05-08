---
trigger: always_on
---

# 组件开发规范

## 代码组织顺序

```
Imports
Props/Emits
常量定义 (as const)
类型定义
响应式数据
计算属性
监听器
工具函数
事件处理方法 (on 开头)
```

## 命名规范

- **事件方法**: 必须以 `on` 开头（`onSubmit`, `onDelete`, `onEdit`）
- **变量/函数**: camelCase
- **常量**: UPPER_SNAKE_CASE + `as const`
- **组件**: PascalCase
- **文件**: 组件用 PascalCase，其他用小写

## Props & Emits

```vue
<script lang="ts" setup>
interface Props {
    visible?: boolean;
    data?: any;
}

const props = withDefaults(defineProps<Props>(), {
    visible: false,
    data: null,
});

const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void;
    (e: 'success'): void;
}>();
</script>
```

## 双向绑定规范

### 使用 defineModel (Vue 3.4+)

**必须使用 `defineModel` 实现双向绑定**，替代旧的 `computed` + `emit('update:xxx')` 模式。

#### 基本用法

```vue
<script lang="ts" setup>
// 单个 v-model
const modelValue = defineModel<string>('modelValue', {
    default: '',
});

// 命名 v-model
const authCertName = defineModel<string>('authCertName');
const machineId = defineModel<number>('machineId');
</script>
```

#### 内部字段联动更新

当组件内部有多个字段，需要联动更新外部的 `modelValue` 时，使用 `watch` 监听：

```vue
<script lang="ts" setup>
import { watch } from 'vue';

const authCertName = defineModel<string>('authCertName');
const machineName = defineModel<string>('machineName');
const selectNode = defineModel<string>('modelValue', { default: '' });

// 监听内部字段变化，自动更新 selectNode
watch(
    [authCertName, machineName],
    () => {
        selectNode.value = authCertName.value 
            ? `${machineName.value} > ${authCertName.value}` 
            : '';
    },
    { immediate: true }
);
</script>
```

#### 规范要点

- ✅ **Always**: 使用 `defineModel` 替代 `computed` + `emit('update:xxx')`
- ✅ **Always**: 为 `defineModel` 提供合适的 `default` 值
- 🚫 **Never**: 使用旧的 `computed` getter/setter 模式实现双向绑定

## 图标使用规范

### 统一使用 SvgIcon 组件

**所有图标必须使用 `SvgIcon` 组件**，禁止使用 `<el-icon>` 配合导入图标组件。

```vue
<!-- ✅ 正确：使用 SvgIcon -->
<SvgIcon name="Monitor" :size="20" />
<SvgIcon name="check" class="text-success" />

<!-- ❌ 错误：使用 el-icon + 导入 -->
<el-icon><Check /></el-icon>
```

**规范要点**：
- ✅ 使用 `name` 属性指定图标，`size` 属性控制大小
- ✅ 图标名称使用 PascalCase 或 kebab-case
- 🚫 禁止使用 `<el-icon>` 和导入 `@element-plus/icons-vue`
- 🚫 禁止通过 class 设置图标大小

### 自定义 SVG 图标

项目支持在 `assets/icon` 目录下添加自定义 SVG 图标。

#### 目录结构

```
frontend/src/assets/icon/
├── db/              # 数据库图标（mysql.svg, postgres.svg...）
├── machine/         # 机器图标
└── ...
```

#### 使用方法

**格式**: `name="icon {目录}/{文件名}"`（不含 .svg）

```vue
<SvgIcon name="icon db/mysql" :size="20" />
```

#### 添加步骤

1. 将 SVG 文件放到 `frontend/src/assets/icon/` 对应子目录
2. 文件名使用小写 + 连字符（如 `mysql.svg`）
3. 使用 `name="icon db/mysql"` 引用

#### 注意事项

- ✅ SVG 必须有 `viewBox` 属性
- ✅ 使用 `size` 属性控制大小
- ✅ 图标颜色继承当前元素的 `color`
- 🚫 不要在 SVG 中硬编码颜色值
- 🚫 文件名不要使用大写或下划线

## 边界

- ✅ **Always**: 使用 Composition API + `<script setup>`
- ✅ **Always**: 事件方法以 `on` 开头
- ✅ **Always**: 移除无用的导入（import）和无用的字段、变量、函数
- 🚫 **Never**: 保留未使用的代码或注释掉的代码
- 🚫 **Never**: 使用固定高度计算，优先用 Flexbox
