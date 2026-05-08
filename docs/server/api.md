---
trigger: always_on
---

# API 层规范

## Handler 标准结构

```go
type Db struct {
    dbApp  application.Db `inject:"T"`
    tagApp tagapp.TagTree `inject:"T"`
}

// @router /api/dbs [get]
func (d *Db) Dbs(rc *req.Ctx) {
    queryCond := req.BindQuery[entity.DbQuery](rc)  // 1. 绑定参数
    loginAccount := rc.GetLoginAccount()            // 2. 获取上下文
    result, err := d.dbApp.GetPageList(queryCond)   // 3. 调用应用层
    biz.ErrIsNil(err)                               // 4. 断言错误（仅API层）
    rc.ResData = result                             // 5. 返回结果
}
```

## 路由配置

```go
func (d *Db) ReqConfs() *req.Confs {
    return req.NewConfs("/dbs",
        req.NewGet("", d.Dbs),
        req.NewPost("", d.Save).Log(req.NewLogSaveI(imsg.LogDbSave)),
        req.NewDelete(":dbId", d.DeleteDb).Log(req.NewLogSaveI(imsg.LogDbDelete)),
    )
}
```

## 断言边界

**✅ API 层可用断言**：

```go
func (d *Db) Save(rc *req.Ctx) {
    form := req.BindFormAndValid[form.DbForm](rc)
    biz.IsTrue(form.InstanceId > 0, "实例ID不能为空")
    biz.ErrIsNil(d.dbApp.SaveDb(rc, &entity.Db{Name: form.Name}))
    rc.ResData = "保存成功"
}
```

**🚫 Application 层禁止断言，必须返回 error**：

```go
func (d *dbAppImpl) SaveDb(ctx context.Context, db *entity.Db) error {
    if db.Name == "" {
        return errorx.NewBiz("名称不能为空")
    }
    return d.Save(ctx, db)
}
```
