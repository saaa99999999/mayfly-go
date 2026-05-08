---
trigger: always_on
---

# Go 分层架构与目录规范

## 分层目录

```
internal/{module}/
├── api/              # HTTP请求处理、参数绑定、响应返回
│   ├── form/         # 请求表单结构体
│   └── vo/           # 响应视图对象
├── application/      # 业务逻辑编排、事务控制
│   └── dto/          # 数据传输对象
├── domain/           # 核心业务逻辑、实体定义
│   ├── entity/       # 领域实体
│   └── repository/   # 仓储接口定义
├── infra/            # 数据持久化、外部服务调用
│   └── persistence/  # 仓储实现
├── imsg/             # 国际化消息定义
└── init/             # 模块初始化（依赖注册、路由注册）
```

## 命名规范

- **模块/包名**: 小写无分隔符（`machine`, `dbinstance`）
- **文件名**: 小写+下划线（`db.go`, `db_sql_exec.go`）
- **结构体/常量**: PascalCase
- **接口**: 以 `er` 结尾或名词（`Reader`, `Repository`）
- **变量/函数**: camelCase

## IOC 依赖注入

```go
// 1. 定义接口
type Db interface {
    base.App[*entity.Db]
    GetPageList(condition *entity.DbQuery, orderBy ...string) (*model.PageResult[*entity.DbListPO], error)
}

// 2. 实现接口并注入依赖
type dbAppImpl struct {
    base.AppImpl[*entity.Db, repository.Db]
    dbInstanceApp Instance       `inject:"T"`  // T=按类型注入
    tagApp        tagapp.TagTree `inject:"T"`
}
var _ Db = (*dbAppImpl)(nil)

// 3. 模块初始化时注册
func init() {
    ioc.Register(&dbAppImpl{})
}
```

## 边界

- ✅ **Always**: 依赖单向流动，上层依赖下层接口，禁止反向依赖
- 🚫 **Never**: 跨层直接调用具体实现，必须通过接口
