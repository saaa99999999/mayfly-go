package dbtool

import (
	"context"

	"mayfly-go/internal/ai/imsg"
	"mayfly-go/internal/ai/tools"
	"mayfly-go/internal/db/application"
	"mayfly-go/internal/db/dbm/dbi"
	"mayfly-go/pkg/i18n"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type QueryTablesParam struct {
	DbId   int64  `json:"dbId" jsonschema_description:"数据库ID。取值逻辑：1. 用户本次明确指定；2. 从前序工具的输入输出中继承已选定的数据库ID；3. 若均无，传0以触发参数补全。禁止凭空猜测。"`
	DbName string `json:"dbName" jsonschema_description:"数据库名称。取值逻辑：1. 用户本次明确指定；2. 从前序工具的输入输出中继承已选定的数据库名称；3. 若均无，留空以触发参数补全。禁止凭空猜测。"`
}

type QueryTablesOutput struct {
	DbId   int64       `json:"dbId" jsonschema_description:"数据库ID"`
	DbName string      `json:"dbName" jsonschema_description:"数据库名称"`
	DbType string      `json:"dbType" jsonschema_description:"数据库类型，如mysql、postgresql等"`
	Tables []dbi.Table `json:"tables" jsonschema_description:"数据库表列表"`
}

func GetQueryTables() (tool.InvokableTool, error) {
	return utils.InferTool("DbQueryTables",
		i18n.T(imsg.DbQueryTablesToolInfo),
		func(ctx context.Context, param *QueryTablesParam) (*QueryTablesOutput, error) {
			toolDesc := i18n.TC(ctx, imsg.DbQueryTablesToolDesc)
			// 检查必要参数，触发参数完善
			if param.DbId == 0 || param.DbName == "" {
				if err := tools.InterruptOrResumeParamCompletion(ctx, toolDesc, param, i18n.TC(ctx, imsg.DbInfoIncomplete), "db", []tools.CompletionParamInfo{
					{Param: "dbId", Name: "数据库ID"},
					{Param: "dbName", Name: "数据库名称"},
				}); err != nil {
					return nil, err
				}
			}

			conn, err := application.GetDbApp().GetDbConn(ctx, uint64(param.DbId), param.DbName)
			if err != nil {
				return nil, tools.NewToolError(err, tools.RecoverRetry)
			}

			tables, err := conn.GetMetadata().GetTables()
			if err != nil {
				return nil, tools.NewToolError(err, tools.RecoverRetry)
			}

			output := &QueryTablesOutput{
				DbId:   param.DbId,
				DbName: param.DbName,
				DbType: string(conn.Info.Type),
				Tables: tables,
			}
			return output, nil
		},
	)
}
