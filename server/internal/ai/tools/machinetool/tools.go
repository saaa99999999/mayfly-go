package machinetool

import (
	"mayfly-go/internal/ai/tools"
	"mayfly-go/pkg/logx"
)

func Init() {
	if commandExecTool, err := GetCommandExec(); err != nil {
		logx.Errorf("agent tool - 获取MachineCommandExec工具失败: %v", err)
	} else {
		tools.DefaultRegistry.Register(commandExecTool)
		tools.RegisterTool(tools.ToolTypeMachine, commandExecTool)
	}
}
