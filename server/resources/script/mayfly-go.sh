#!/bin/bash

#==============================================
# Mayfly-Go Service Manager
# 服务管理脚本：支持 start/stop/restart/status
#==============================================

set -e

# 配置
BINARY_NAME="mayfly-go"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="${SCRIPT_DIR}/bin/${BINARY_NAME}"
PID_FILE="${SCRIPT_DIR}/${BINARY_NAME}.pid"
CONFIG_FILE="${SCRIPT_DIR}/config.yml"
STARTUP_LOG="${SCRIPT_DIR}/startup.log"

# 颜色输出
print_info() {
    echo -e "\033[1;34m$1\033[0m"
}

print_success() {
    echo -e "\033[1;32m$1\033[0m"
}

print_error() {
    echo -e "\033[1;31m$1\033[0m" >&2
}

print_warning() {
    echo -e "\033[1;33m$1\033[0m"
}

# 获取进程 PID
get_pid() {
    if [ -f "${PID_FILE}" ]; then
        local pid=$(cat "${PID_FILE}")
        if ps -p "${pid}" > /dev/null 2>&1; then
            echo "${pid}"
            return 0
        fi
    fi
    
    # 回退：通过进程名查找
    local pid=$(ps aux | grep "${BINARY_PATH}" | grep -v grep | grep -v "$(basename "$0")" | awk '{print $2}' | head -n 1)
    if [ -n "${pid}" ]; then
        echo "${pid}"
        return 0
    fi
    
    return 1
}

# 检查服务状态
do_status() {
    local pid
    if pid=$(get_pid); then
        print_success "● ${BINARY_NAME} is running (PID: ${pid})"
        return 0
    else
        print_warning "○ ${BINARY_NAME} is not running"
        return 1
    fi
}

# 启动服务
do_start() {
    # 检查是否已在运行
    local pid
    if pid=$(get_pid); then
        print_warning "${BINARY_NAME} is already running (PID: ${pid})"
        return 0
    fi
    
    # 检查二进制文件
    if [ ! -f "${BINARY_PATH}" ]; then
        print_error "Error: ${BINARY_PATH} not found!"
        return 1
    fi
    
    # 确保可执行权限
    if [ ! -x "${BINARY_PATH}" ]; then
        print_info "Setting execute permission..."
        chmod +x "${BINARY_PATH}"
    fi
    
    print_info "Starting ${BINARY_NAME}..."
    nohup "${BINARY_PATH}" -e "${CONFIG_FILE}" >> "${STARTUP_LOG}" 2>&1 &
    local new_pid=$!
    echo "${new_pid}" > "${PID_FILE}"
    
    # 等待启动
    sleep 1
    if ps -p "${new_pid}" > /dev/null 2>&1; then
        print_success "✓ ${BINARY_NAME} started successfully (PID: ${new_pid})"
        print_info "  Startup log: ${STARTUP_LOG}"
        return 0
    else
        print_error "✗ Failed to start ${BINARY_NAME}"
        print_error "  Check startup log: ${STARTUP_LOG}"
        rm -f "${PID_FILE}"
        return 1
    fi
}

# 停止服务
do_stop() {
    local pid
    if ! pid=$(get_pid); then
        print_warning "${BINARY_NAME} is not running"
        return 0
    fi
    
    print_info "Stopping ${BINARY_NAME} (PID: ${pid})..."
    
    # 优雅关闭
    kill "${pid}" 2>/dev/null || true
    
    # 等待进程退出
    local count=0
    while ps -p "${pid}" > /dev/null 2>&1; do
        sleep 0.5
        count=$((count + 1))
        if [ ${count} -ge 20 ]; then  # 最多等待 10 秒
            print_warning "Process not responding, force killing..."
            kill -9 "${pid}" 2>/dev/null || true
            break
        fi
    done
    
    rm -f "${PID_FILE}"
    print_success "✓ ${BINARY_NAME} stopped"
    return 0
}

# 重启服务
do_restart() {
    print_info "Restarting ${BINARY_NAME}..."
    do_stop
    sleep 1
    do_start
}

# 显示帮助
do_help() {
    echo "Usage: $(basename "$0") {start|stop|restart|status}"
    echo ""
    echo "Commands:"
    echo "  start     Start the service"
    echo "  stop      Stop the service"
    echo "  restart   Restart the service"
    echo "  status    Check service status"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0") start"
    echo "  $(basename "$0") stop"
    echo "  $(basename "$0") restart"
    echo "  $(basename "$0") status"
}

# 主函数
main() {
    local command="${1:-}"
    
    case "${command}" in
        start)
            do_start
            ;;
        stop)
            do_stop
            ;;
        restart)
            do_restart
            ;;
        status)
            do_status
            ;;
        *)
            do_help
            exit 1
            ;;
    esac
}

main "$@"