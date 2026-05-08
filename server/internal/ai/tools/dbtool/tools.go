package dbtool

import (
	"mayfly-go/internal/ai/tools"
	"mayfly-go/pkg/logx"
)

func Init() {
	if queryTableDDLTool, err := GetQueryTableDDL(); err != nil {
		logx.Errorf("agent tool - 获取QueryTableDDL工具失败: %v", err)
	} else {
		tools.DefaultRegistry.Register(queryTableDDLTool)
	}

	if queryTablesTool, err := GetQueryTables(); err != nil {
		logx.Errorf("agent tool - 获取QueryTables工具失败: %v", err)
	} else {
		tools.DefaultRegistry.Register(queryTablesTool)
	}

	if queryDataTool, err := GetQueryData(); err != nil {
		logx.Errorf("agent tool - 获取QueryData工具失败: %v", err)
	} else {
		tools.DefaultRegistry.Register(queryDataTool)
	}

	if sqlExecTool, err := GetSqlExec(); err != nil {
		logx.Errorf("agent tool - 获取ExecSql工具失败: %v", err)
	} else {
		tools.DefaultRegistry.Register(sqlExecTool)
	}
}
