---
trigger: always_on
---

# 安全与权限规范

## 权限控制

```go
// 路由级别
req.NewPost(":dbId/exec-sql", d.ExecSql).RequiredPermissionCode("db:sqlscript:run")

// 代码级别
biz.IsTrue(account.HasPermission("db:sqlscript:run"), "无权限执行SQL")
```

## 敏感信息

- 资源密码使用 AES 加密存储
- `aes.key` 和 `jwt.key` 必须使用随机字符串

## OWASP 安全准则

- 防范 SQL 注入：使用参数化查询
- 防范 XSS：输出转义
- 防范 CSRF：配合前端同源策略
