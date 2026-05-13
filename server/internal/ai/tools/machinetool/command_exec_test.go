package machinetool

import (
	"testing"
)

func TestIsWhitelistCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool // true表示在白名单中，可以自动执行
	}{
		// 白名单命令测试（不应被拦截）
		{
			name:     "ls命令",
			command:  "ls -la",
			expected: true,
		},
		{
			name:     "free命令",
			command:  "free -m",
			expected: true,
		},
		{
			name:     "df命令",
			command:  "df -h",
			expected: true,
		},
		{
			name:     "cat查看文件",
			command:  "cat /etc/passwd",
			expected: true,
		},
		{
			name:     "ps查看进程",
			command:  "ps aux",
			expected: true,
		},
		{
			name:     "组合安全命令",
			command:  "ls -la && cat file.txt | grep test",
			expected: true,
		},
		{
			name:     "查看系统状态",
			command:  "ps aux | grep nginx",
			expected: true,
		},
		{
			name:     "查看磁盘使用",
			command:  "df -h && du -sh /var/log",
			expected: true,
		},
		{
			name:     "带引号的命令",
			command:  "echo \"hello world\"",
			expected: true,
		},
		{
			name:     "复杂管道命令",
			command:  "cat /var/log/syslog | grep error | wc -l",
			expected: true,
		},
		{
			name:     "uname系统信息",
			command:  "uname -a",
			expected: true,
		},
		{
			name:     "ping网络测试",
			command:  "ping -c 4 google.com",
			expected: true,
		},

		// 非白名单命令测试（需要审批）
		{
			name:     "rm删除命令",
			command:  "rm /tmp/test.txt",
			expected: false,
		},
		{
			name:     "rm -rf强制删除",
			command:  "rm -rf /tmp/test",
			expected: false,
		},
		{
			name:     "shutdown关机",
			command:  "shutdown -h now",
			expected: false,
		},
		{
			name:     "reboot重启",
			command:  "reboot",
			expected: false,
		},
		{
			name:     "dd磁盘写入",
			command:  "dd if=/dev/zero of=/dev/sda",
			expected: false,
		},
		{
			name:     "mkfs格式化",
			command:  "mkfs.ext4 /dev/sda1",
			expected: false,
		},
		{
			name:     "fdisk分区",
			command:  "fdisk /dev/sda",
			expected: false,
		},
		{
			name:     "chmod修改权限",
			command:  "chmod 755 /usr/local/bin/app",
			expected: false,
		},
		{
			name:     "echo重定向",
			command:  "echo test > /tmp/output.txt",
			expected: true, // echo在白名单中，重定向到普通文件是允许的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWhitelistCommand(tt.command)
			if result != tt.expected {
				t.Errorf("isWhitelistCommand(%q) = %v, expected %v", tt.command, result, tt.expected)
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
