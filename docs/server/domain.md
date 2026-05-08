---
trigger: always_on
---

# Domain 层规范

## 实体定义

```go
package entity

import "mayfly-go/pkg/model"

type Db struct {
    model.Model      // 必须嵌入基础模型
    model.ExtraData  // 辅助字段（展示用、非查询条件）

    Code       string `json:"code" gorm:"size:32;not null;index:idx_db_code"`
    Name       string `json:"name" gorm:"size:255;not null;"`
    InstanceId uint64 `json:"instanceId" gorm:"not null;"`
}

type Status int8
const (
    StatusActive   Status = 1
    StatusInactive Status = 0
)
```

## ExtraData 使用原则

- ✅ **使用 ExtraData**: 前端展示字段、关联名称、状态文本、可选扩展信息
- 🚫 **必须独立字段**: 查询条件、排序字段、分组统计、索引字段、核心业务字段

## Repository 接口

```go
package repository

type Db interface {
    base.Repo[*entity.Db]
    GetDbList(condition *entity.DbQuery, orderBy ...string) (*model.PageResult[*entity.DbListPO], error)
}
```
