---
trigger: always_on
---

# 并发与 Panic 处理规范

## 统一 Panic 捕获（gox.Recover）

**核心原则**：严禁手动编写 `defer func() { recover() }`，必须使用 `gox.Recover()`

### 场景1：仅记录日志

```go
func (s *Service) ProcessData(data []byte) {
    defer gox.Recover()
    result := parseData(data)
    saveToDB(result)
}
```

### 场景2：Panic 转 Error 返回

```go
func (s *Service) SaveUser(ctx context.Context, user *entity.User) (err error) {
    defer gox.Recover(func(e error) {
        err = fmt.Errorf("保存用户失败: %w", e)
    })
    if err := validateUser(user); err != nil {
        return err
    }
    return s.repo.Insert(ctx, user)
}
```

### 场景3：Goroutine 安全启动

```go
// ✅ 推荐
gox.Go(func() {
    sendNotification(userId, message)
})

// 🚫 禁止
go func() {
    sendNotification(userId, message)
}()
```

## Context 传递

所有阻塞操作必须接受 `context.Context`：

```go
func (d *dbAppImpl) SaveDb(ctx context.Context, entity *entity.Db) error {
    return d.GetRepo().Insert(ctx, entity)
}
```

## 错误组使用

```go
eg, ctx := errgroup.WithContext(context.Background())
for _, task := range tasks {
    eg.Go(func() error {
        return process(ctx, task)
    })
}
err := eg.Wait()
```
