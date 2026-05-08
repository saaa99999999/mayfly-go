---
trigger: always_on
---

# Application 层规范

## 接口与实现

```go
type Db interface {
    base.App[*entity.Db]
    GetPageList(condition *entity.DbQuery, orderBy ...string) (*model.PageResult[*entity.DbListPO], error)
    SaveDb(ctx context.Context, entity *entity.Db) error
}

type dbAppImpl struct {
    base.AppImpl[*entity.Db, repository.Db]
    dbInstanceApp Instance       `inject:"T"`
    tagApp        tagapp.TagTree `inject:"T"`
}
var _ Db = (*dbAppImpl)(nil)

func (d *dbAppImpl) SaveDb(ctx context.Context, dbEntity *entity.Db) error {
    // 1. 参数校验（返回error）
    if dbEntity.Name == "" {
        return errorx.NewBiz("名称不能为空")
    }
    // 2. 业务检查
    oldDb := &entity.Db{Name: dbEntity.Name, InstanceId: dbEntity.InstanceId}
    if dbEntity.Id == 0 && d.GetByCond(oldDb) == nil {
        return errorx.NewBizI(ctx, imsg.ErrDbNameExist)
    }
    // 3. 持久化
    return d.Save(ctx, dbEntity)
}
```

## 错误处理

```go
// 普通业务错误
return errorx.NewBiz("数据库名称已存在")

// 国际化错误
return errorx.NewBizI(ctx, imsg.ErrDbNameExist)
```

## 边界

- ✅ **Always**: 参数校验后返回 error，禁止 panic
- 🚫 **Never**: 在 application 层使用 `biz.ErrIsNil` 或 `biz.IsTrue`
