# Mayfly-Go 服务管理脚本使用说明

## 概述

`mayfly-go.sh` 是 Mayfly-Go 的统一服务管理脚本，支持启动、停止、重启和查看状态等操作。

## 快速开始

### 查看帮助

```bash
./mayfly-go.sh
```

### 启动服务

```bash
./mayfly-go.sh start
```

启动后会：
- 检查服务是否已在运行
- 验证二进制文件是否存在
- 自动设置执行权限
- 后台启动服务
- 保存 PID 到 `mayfly-go.pid`

### 停止服务

```bash
./mayfly-go.sh stop
```

停止时会：
- 优雅关闭服务（发送 SIGTERM）
- 最多等待 10 秒
- 如无响应则强制关闭（SIGKILL）
- 清理 PID 文件

### 重启服务

```bash
./mayfly-go.sh restart
```

等同于先 stop 再 start。

### 查看状态

```bash
./mayfly-go.sh status
```

显示服务是否运行及进程 PID。

## 输出示例

### 启动
```
Starting mayfly-go...
✓ mayfly-go started successfully (PID: 12345)
```

### 停止
```
Stopping mayfly-go (PID: 12345)...
✓ mayfly-go stopped
```

### 状态（运行中）
```
● mayfly-go is running (PID: 12345)
```

### 状态（未运行）
```
○ mayfly-go is not running
```

## 日志查看

### 启动日志

启动日志记录服务的启动过程和错误信息：

```bash
# 查看启动日志
cat startup.log

# 实时查看启动日志
tail -f startup.log
```

### 运行日志

服务运行日志由 `config.yml` 中的日志配置决定，请查看配置文件中的日志设置。

```bash
# 根据配置的日志路径查看日志
tail -f <your-log-path>

# 查看最近 100 行
tail -n 100 <your-log-path>

# 搜索错误
grep -i error <your-log-path>
```

## 文件说明

| 文件 | 说明 |
|------|------|
| `bin/mayfly-go` | 主程序二进制文件 |
| `mayfly-go.sh` | 服务管理脚本 |
| `mayfly-go.pid` | 进程 PID 文件（自动创建） |
| `startup.log` | 启动日志文件（自动创建） |
| `config.yml` | 配置文件（包含日志配置） |

## 故障排查

### 启动失败

1. 检查二进制文件是否存在：`ls -la bin/mayfly-go`
2. 检查端口是否被占用：`lsof -i :8001`（假设使用 8001 端口）
3. 查看启动日志：`cat startup.log`
4. 查看配置文件中指定的运行日志路径，检查错误信息

### 无法正常停止

如果服务无法正常停止，可以手动清理：

```bash
# 查找进程
ps aux | grep mayfly-go

# 强制杀死进程
kill -9 <PID>

# 清理 PID 文件
rm -f mayfly-go.pid
```

### PID 文件不同步

如果 PID 文件存在但进程不存在，脚本会自动检测并重新创建：

```bash
# 手动清理
rm -f mayfly-go.pid

# 重新启动
./mayfly-go.sh start
```

## 高级用法

### 开机自启动

**systemd 方式（推荐）：**

创建 `/etc/systemd/system/mayfly-go.service`：

```ini
[Unit]
Description=Mayfly-Go Service
After=network.target

[Service]
Type=forking
User=mayfly
Group=mayfly
WorkingDirectory=/opt/mayfly-go
ExecStart=/opt/mayfly-go/mayfly-go.sh start
ExecStop=/opt/mayfly-go/mayfly-go.sh stop
ExecReload=/opt/mayfly-go/mayfly-go.sh restart
PIDFile=/opt/mayfly-go/mayfly-go.pid

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl enable mayfly-go
sudo systemctl start mayfly-go
sudo systemctl status mayfly-go
```

**crontab 方式：**

```bash
# 编辑 crontab
crontab -e

# 添加重启后自动启动
@reboot /opt/mayfly-go/mayfly-go.sh start
```

## 注意事项

1. 脚本需要在包含 `mayfly-go` 二进制文件的目录下运行
2. 确保有足够的权限执行二进制文件
3. 建议不要同时运行多个实例
4. 日志配置请在 `config.yml` 中设置
