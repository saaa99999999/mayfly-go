---
trigger: always_on
---

# 代码质量与 Git 规范

## 函数长度

- 单个函数不超过 100 行
- 复杂逻辑拆分为私有方法

## Error 处理

```go
// ✅ 完整处理
result, err := doSomething()
if err != nil {
    logx.Errorf("操作失败: %v", err)
    return errorx.NewBiz("操作失败")
}

// 🚫 忽略错误
result, _ := doSomething()
```

## 资源释放

```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()
```

## 魔法数字

```go
const MaxRetryCount = 3
if retry > MaxRetryCount { ... } // ✅
if retry > 3 { ... }             // 🚫
```

## Git 提交格式

```
<type>(<scope>): <subject>

<body>
```

**Type 类型**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

**示例**:
```
feat(db): 添加数据库备份功能

- 实现定时备份任务
- 支持增量备份和全量备份

Closes #123
```
