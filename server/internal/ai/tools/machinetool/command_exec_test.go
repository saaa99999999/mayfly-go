package machinetool

import (
	"testing"
)

func TestIsDangerousCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		// 危险命令测试
		{
			name:     "rm -rf 根目录",
			command:  "rm -rf /",
			expected: true,
		},
		{
			name:     "rm -fr 目录",
			command:  "rm -fr /tmp/test",
			expected: true,
		},
		{
			name:     "rm -Rf 目录",
			command:  "rm -Rf /tmp/test",
			expected: true,
		},
		{
			name:     "组合命令中的rm -rf",
			command:  "ls -la && rm -rf /tmp/test",
			expected: true,
		},
		{
			name:     "管道符后的危险命令",
			command:  "echo test | rm -rf /",
			expected: true,
		},
		{
			name:     "分号分隔的危险命令",
			command:  "cd /tmp; rm -rf *",
			expected: true,
		},
		{
			name:     "mkfs格式化命令",
			command:  "mkfs.ext4 /dev/sdb1",
			expected: true,
		},
		{
			name:     "mkfs命令",
			command:  "mkfs -t ext4 /dev/sdb1",
			expected: true,
		},
		{
			name:     "dd磁盘写入",
			command:  "dd if=/dev/zero of=/dev/sda",
			expected: true,
		},
		{
			name:     "dd of参数",
			command:  "dd of=/dev/sda bs=1M",
			expected: true,
		},
		{
			name:     "shutdown关机",
			command:  "shutdown -h now",
			expected: true,
		},
		{
			name:     "reboot重启",
			command:  "reboot",
			expected: true,
		},
		{
			name:     "halt停机",
			command:  "halt",
			expected: true,
		},
		{
			name:     "poweroff关机",
			command:  "poweroff",
			expected: true,
		},
		{
			name:     "chmod 777 根目录",
			command:  "chmod 777 /",
			expected: true,
		},
		{
			name:     "chmod 0777 根目录",
			command:  "chmod 0777 /",
			expected: true,
		},
		{
			name:     "写入磁盘设备",
			command:  "echo test > /dev/sda",
			expected: true,
		},
		{
			name:     "追加重定向到磁盘",
			command:  "echo test >> /dev/sdb1",
			expected: true,
		},
		{
			name:     "重定向到nvme设备",
			command:  "cat file > /dev/nvme0n1",
			expected: true,
		},
		{
			name:     "fdisk分区操作",
			command:  "fdisk /dev/sda",
			expected: true,
		},

		// 安全命令测试（不应被拦截）
		{
			name:     "普通rm删除单个文件",
			command:  "rm /tmp/test.txt",
			expected: false,
		},
		{
			name:     "rm -i交互式删除",
			command:  "rm -i /tmp/test.txt",
			expected: false,
		},
		{
			name:     "rm -r不带f参数",
			command:  "rm -r /tmp/test",
			expected: false,
		},
		{
			name:     "ls命令",
			command:  "ls -la",
			expected: false,
		},
		{
			name:     "包含firmware文本",
			command:  "echo 'firmware update'",
			expected: false,
		},
		{
			name:     "grep查找rm相关日志",
			command:  "grep 'rm process' /var/log/syslog",
			expected: false,
		},
		{
			name:     "chmod设置合理权限",
			command:  "chmod 755 /usr/local/bin/app",
			expected: false,
		},
		{
			name:     "chmod 777非根目录",
			command:  "chmod 777 /tmp/test",
			expected: false,
		},
		{
			name:     "组合安全命令",
			command:  "ls -la && cat file.txt | grep test",
			expected: false,
		},
		{
			name:     "查看系统状态",
			command:  "ps aux | grep nginx",
			expected: false,
		},
		{
			name:     "查看磁盘使用",
			command:  "df -h && du -sh /var/log",
			expected: false,
		},
		{
			name:     "重定向到普通文件",
			command:  "echo test > /tmp/output.txt",
			expected: false,
		},
		{
			name:     "dd不带if参数",
			command:  "dd status=progress",
			expected: false,
		},
		{
			name:     "带引号的命令",
			command:  "echo \"rm -rf / is dangerous\"",
			expected: false,
		},
		{
			name:     "复杂管道命令",
			command:  "cat /var/log/syslog | grep error | wc -l",
			expected: false,
		},
		{
			name:     "带转义字符的命令",
			command:  "echo \"hello world\"",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDangerousCommand(tt.command)
			if result != tt.expected {
				t.Errorf("isDangerousCommand(%q) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []string
	}{
		{
			name:     "简单命令",
			command:  "ls -la",
			expected: []string{"ls", "-la"},
		},
		{
			name:     "&&连接符",
			command:  "ls -la && cat file.txt",
			expected: []string{"ls", "-la", "&&", "cat", "file.txt"},
		},
		{
			name:     "管道符",
			command:  "ps aux | grep nginx",
			expected: []string{"ps", "aux", "|", "grep", "nginx"},
		},
		{
			name:     "分号分隔",
			command:  "cd /tmp; ls -la",
			expected: []string{"cd", "/tmp", ";", "ls", "-la"},
		},
		{
			name:     "混合分隔符",
			command:  "ls && cat file | grep test; echo done",
			expected: []string{"ls", "&&", "cat", "file", "|", "grep", "test", ";", "echo", "done"},
		},
		{
			name:     "双引号字符串",
			command:  "echo \"hello world\"",
			expected: []string{"echo", "\"hello world\""},
		},
		{
			name:     "单引号字符串",
			command:  "echo 'hello world'",
			expected: []string{"echo", "'hello world'"},
		},
		{
			name:     "重定向",
			command:  "echo test > output.txt",
			expected: []string{"echo", "test", ">", "output.txt"},
		},
		{
			name:     "追加重定向",
			command:  "echo test >> output.txt",
			expected: []string{"echo", "test", ">>", "output.txt"},
		},
		{
			name:     "带转义字符",
			command:  "echo \"hello\\\"world\"",
			expected: []string{"echo", "\"hello\\\"world\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenize(tt.command)
			if len(result) != len(tt.expected) {
				t.Errorf("tokenize(%q) returned %d tokens, expected %d\ngot: %v\nexpected: %v",
					tt.command, len(result), len(tt.expected), result, tt.expected)
				return
			}
			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("tokenize(%q)[%d] = %q, expected %q\nfull result: %v",
						tt.command, i, token, tt.expected[i], result)
				}
			}
		})
	}
}
