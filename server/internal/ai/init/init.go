package init

import (
	"mayfly-go/internal/ai/api"
	"mayfly-go/internal/ai/application"
	"mayfly-go/internal/ai/infra/persistence"
	"mayfly-go/internal/ai/tools/dbtool"
	"mayfly-go/internal/ai/tools/machinetool"
	"mayfly-go/pkg/starter"
)

func init() {
	// 注册AI模块的IoC组件
	starter.AddInitIocFunc(func() {
		persistence.InitIoc()
		application.Init()
		api.InitIoc()
	})

	// 数据库工具初始化
	dbtool.Init()
	// 机器工具初始化
	machinetool.Init()
}
