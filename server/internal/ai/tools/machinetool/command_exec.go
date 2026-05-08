package machinetool

import (
	"context"
	"fmt"
	"mayfly-go/internal/ai/imsg"
	"mayfly-go/internal/ai/tools"
	"mayfly-go/internal/machine/application"
	"mayfly-go/pkg/i18n"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// CommandExecParam 命令执行参数
type CommandExecParam struct {
	AuthCertName string `json:"authCertName" jsonschema_description:"授权凭证名称"`
	Command      string `json:"command" jsonschema_description:"要执行的命令"`
	Remark       string `json:"remark" jsonschema_description:"命令作用说明，简要描述该命令的用途或目的"`
}

// CommandExecOutput 命令执行输出
type CommandExecOutput struct {
	AuthCertName string `json:"authCertName" jsonschema_description:"授权凭证名称"`
	MachineId    uint64 `json:"machineId" jsonschema_description:"机器ID"`
	MachineName  string `json:"machineName" jsonschema_description:"机器名称"`
	MachineIp    string `json:"machineIp" jsonschema_description:"机器IP地址"`
	MachinePort  int    `json:"machinePort" jsonschema_description:"机器端口"`
	Username     string `json:"username" jsonschema_description:"连接用户名"`
	Output       string `json:"output" jsonschema_description:"命令执行输出"`
	Success      bool   `json:"success" jsonschema_description:"是否执行成功"`
}

// GetCommandExec 获取命令执行工具
func GetCommandExec() (tool.InvokableTool, error) {
	return utils.InferTool("MachineCommandExec",
		i18n.T(imsg.MachineCommandExecToolInfo),
		func(ctx context.Context, param *CommandExecParam) (*CommandExecOutput, error) {
			toolDesc := i18n.TC(ctx, imsg.MachineCommandExecToolDesc)

			// 检查必要参数，触发参数完善
			if param.AuthCertName == "" {
				if err := tools.InterruptOrResumeParamCompletion(ctx, toolDesc, param, i18n.TC(ctx, imsg.MachineInfoIncomplete), "machine", []tools.CompletionParamInfo{
					{Param: "authCertName", Name: "授权凭证名称"},
				}); err != nil {
					return nil, err
				}
			}

			// 检查命令是否为空
			if param.Command == "" {
				return nil, tools.NewToolError(fmt.Errorf("%s", i18n.TC(ctx, imsg.MissingRequiredParams)), tools.RecoverRetry)
			}

			// 危险命令检测
			if isDangerousCommand(param.Command) {
				// 触发审批中断
				if err := tools.InterruptOrResumeApproval(ctx, toolDesc, param, i18n.TC(ctx, imsg.CommandExecApprovalReason)); err != nil {
					return nil, err
				}
			}

			// 获取机器客户端
			cli, err := application.GetMachineApp().GetCliByAc(ctx, param.AuthCertName)
			if err != nil {
				return nil, tools.NewToolError(err, tools.RecoverRetry)
			}

			// 执行命令
			output, err := cli.Run(param.Command)
			success := err == nil

			// 从 CLI 中获取机器信息
			machineInfo := cli.Info

			// 即使执行失败也返回输出，让AI能看到错误信息
			return &CommandExecOutput{
				AuthCertName: param.AuthCertName,
				MachineId:    machineInfo.Id,
				MachineName:  machineInfo.Name,
				MachineIp:    machineInfo.Ip,
				MachinePort:  machineInfo.Port,
				Username:     machineInfo.Username,
				Output:       output,
				Success:      success,
			}, nil
		},
	)
}

// isDangerousCommand 判断命令是否包含危险操作
// 使用词法分析遍历所有命令，检测危险命令模式
func isDangerousCommand(cmd string) bool {
	if cmd == "" {
		return false
	}

	// 定义危险命令规则
	type dangerousRule struct {
		command string // 命令名
		pattern string // 危险参数模式（空表示命令本身就危险）
	}

	dangerousRules := []dangerousRule{
		{command: "rm", pattern: ""},   // rm 删除
		{command: "mkfs", pattern: ""}, // 所有格式化命令
		{command: "mkfs.ext2", pattern: ""},
		{command: "mkfs.ext3", pattern: ""},
		{command: "mkfs.ext4", pattern: ""},
		{command: "mkfs.xfs", pattern: ""},
		{command: "mkfs.vfat", pattern: ""},
		{command: "dd", pattern: "if="},      // dd if= 磁盘写入
		{command: "dd", pattern: "of=/dev/"}, // dd of=/dev/ 写入设备
		{command: "shutdown", pattern: ""},   // 关机命令
		{command: "reboot", pattern: ""},     // 重启命令
		{command: "halt", pattern: ""},       // 停机命令
		{command: "poweroff", pattern: ""},   // 断电命令
		{command: "fdisk", pattern: ""},      // 分区操作
		{command: "parted", pattern: ""},     // 分区工具
		{command: "debugfs", pattern: ""},    // 文件系统调试
		{command: "ddrescue", pattern: ""},   // 磁盘救援
	}

	// 词法分析命令
	tokens := tokenize(cmd)

	// 遍历 tokens，查找命令及其参数
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// 跳过操作符和引号内容
		if isOperator(token) || isQuotedString(token) {
			continue
		}

		// 提取命令名（去除路径）
		cmdName := token
		if idx := strings.LastIndex(token, "/"); idx >= 0 {
			cmdName = token[idx+1:]
		}

		// 检查是否匹配危险命令
		for _, rule := range dangerousRules {
			if cmdName != rule.command {
				continue
			}

			// 如果规则没有 pattern，命令本身就危险
			if rule.pattern == "" {
				return true
			}

			// 检查后续参数是否包含危险模式
			if checkArgsContainPattern(tokens[i+1:], rule.pattern) {
				return true
			}
		}

		// 特殊处理：chmod 777 /
		if cmdName == "chmod" && checkChmod777Root(tokens[i+1:]) {
			return true
		}

		// 特殊处理：重定向到危险设备
		if checkRedirectToDevice(tokens[i:]) {
			return true
		}
	}

	return false
}

// tokenize 将命令字符串分割成 tokens
func tokenize(cmd string) []string {
	var tokens []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	for i := 0; i < len(cmd); i++ {
		ch := cmd[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' && !inSingleQuote {
			current.WriteByte(ch)
			escaped = true
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			current.WriteByte(ch)
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			current.WriteByte(ch)
			continue
		}

		if inSingleQuote || inDoubleQuote {
			current.WriteByte(ch)
			continue
		}

		if isOperatorChar(ch) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			operator := string(ch)
			if i+1 < len(cmd) && isOperatorChar(cmd[i+1]) {
				nextCh := cmd[i+1]
				if ch == nextCh {
					operator += string(nextCh)
					i++
				}
			}

			tokens = append(tokens, operator)
			continue
		}

		if ch == ' ' || ch == '\t' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// isOperatorChar 判断字符是否为操作符字符
func isOperatorChar(ch byte) bool {
	return ch == '|' || ch == '&' || ch == ';' || ch == '>' || ch == '<'
}

// isOperator 判断 token 是否为操作符
func isOperator(token string) bool {
	if len(token) == 0 {
		return false
	}
	return isOperatorChar(token[0])
}

// checkArgsContainPattern 检查参数列表是否包含指定模式
func checkArgsContainPattern(args []string, pattern string) bool {
	for _, arg := range args {
		// 遇到操作符就停止
		if isOperator(arg) {
			break
		}

		// 检查参数
		if strings.Contains(arg, pattern) {
			return true
		}

		// 处理组合选项，如 -rf 拆分为 -r -f
		if strings.HasPrefix(arg, "-") && len(arg) > 2 {
			opts := arg[1:]
			// 检查是否包含 pattern 中的字符（如 -rf 包含 r 和 f）
			if pattern[0] == '-' {
				patternChars := pattern[1:]
				hasAll := true
				for _, pc := range patternChars {
					if !strings.ContainsRune(opts, pc) {
						hasAll = false
						break
					}
				}
				if hasAll {
					return true
				}
			}
		}
	}
	return false
}

// checkChmod777Root 检查是否为 chmod 777 /
func checkChmod777Root(args []string) bool {
	for i, arg := range args {
		if isOperator(arg) {
			break
		}

		// 跳过选项
		if strings.HasPrefix(arg, "-") {
			continue
		}

		// 检查权限模式
		if arg == "777" || arg == "0777" {
			// 检查目标是否为根目录
			if i+1 < len(args) {
				target := args[i+1]
				if !isOperator(target) && (target == "/" || target == "/*") {
					return true
				}
			}
		}
	}
	return false
}

// checkRedirectToDevice 检查是否重定向到危险设备
func checkRedirectToDevice(tokens []string) bool {
	for i, token := range tokens {
		// 检查重定向操作符
		if token == ">" || token == ">>" || token == "1>" || token == "2>" || token == "&>" {
			// 检查下一个 token 是否为危险设备
			if i+1 < len(tokens) {
				target := tokens[i+1]
				if strings.HasPrefix(target, "/dev/sd") ||
					strings.HasPrefix(target, "/dev/hd") ||
					strings.HasPrefix(target, "/dev/nvme") ||
					strings.HasPrefix(target, "/dev/vd") {
					return true
				}
			}
		}
	}
	return false
}

// isQuotedString 判断是否为引号字符串
func isQuotedString(token string) bool {
	if len(token) < 2 {
		return false
	}
	// 单引号或双引号包裹
	return (token[0] == '\'' && token[len(token)-1] == '\'') ||
		(token[0] == '"' && token[len(token)-1] == '"')
}
