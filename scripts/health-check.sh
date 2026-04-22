#!/bin/bash

# ByteDanceDemo 健康检查脚本
# 用途：快速诊断项目和服务状态
# 使用方法：./scripts/health-check.sh

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== ByteDanceDemo 健康检查 ===${NC}"
echo ""

# 函数定义

check_service() {
    local service_name=$1
    local check_command=$2

    echo -n -e "${YELLOW}检查 $service_name: ${NC}"
    if eval "$check_command" &>/dev/null 2>&1; then
        echo -e "${GREEN}✓ 运行中${NC}"
        return 0
    else
        echo -e "${RED}✗ 未运行${NC}"
        return 1
    fi
}

check_port() {
    local port=$1
    local service_name=$2

    echo -n -e "${YELLOW}检查端口 $port ($service_name): ${NC}"
    if command -v netstat >/dev/null 2>&1; then
        if netstat -tuln 2>/dev/null | grep -q ":$port"; then
            echo -e "${GREEN}✓ 在使用中${NC}"
            return 0
        else
            echo -e "${RED}✗ 未在使用${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 netstat${NC}"
        if command -v ss >/dev/null 2>&1; then
            if ss -tuln 2>/dev/null | grep -q ":$port"; then
                echo -e "${GREEN}✓ 在使用中${NC}"
                return 0
            else
                echo -e "${RED}✗ 未在使用${NC}"
                return 1
            fi
        fi
        return 1
    fi
}

check_disk() {
    local threshold=80  # 80% 使用率警告

    echo -e "${YELLOW}磁盘使用情况：${NC}"
    df -h | grep -vE "^/dev|^Filesystem" | while read line; do
        usage=$(echo "$line" | awk '{print $5}' | sed 's/%//')
        usage_int=${usage%%.*}

        if [ $usage_int -gt $threshold ]; then
            echo -e "  $line ${RED}(使用率: $usage%)${NC}"
        else
            echo -e "  $line ${GREEN}(使用率: $usage%)${NC}"
        fi
    done
}

check_memory() {
    local threshold=85  # 85% 使用率警告

    echo -e "${YELLOW}内存使用情况：${NC}"
    free -h

    # 检查 Go 进程内存使用
    if command -v pgrep >/dev/null 2>&1; then
        if pgrep -q "simple-demo"; then
            echo -e "  ${GREEN}Go 应用进程运行中${NC}"
            pid=$(pgrep "simple-demo" | head -1 | awk '{print $1}')
            if [ -n "$pid" ]; then
                echo -e "  进程 PID: $pid"
                ps -p "$pid" -o %mem= 2>/dev/null | tail -1
            fi
        else
            echo -e "  ${RED}Go 应用进程未运行${NC}"
        fi
    fi
}

check_mysql() {
    echo -e "${YELLOW}=== MySQL 数据库检查 ===${NC}"

    # 检查 MySQL 服务
    if command -v systemctl >/dev/null 2>&1; then
        if systemctl is-active --quiet mysql || systemctl is-active --quiet mysql-server || systemctl is-active --quiet mariadb; then
            echo -e "${GREEN}✓ MySQL 服务运行中${NC}"
        elif systemctl is-enabled --quiet mysql || systemctl is-enabled --quiet mysql-server || systemctl is-enabled --quiet mariadb; then
            echo -e "${YELLOW}⚠ MySQL 服务已启用但未运行${NC}"
        else
            echo -e "${YELLOW}⚠ MySQL 服务未在 systemd 中配置${NC}"
        fi
    fi

    # 检查 MySQL 连接
    if command -v mysqladmin >/dev/null 2>&1; then
        echo -n -e "${YELLOW}检查 MySQL 连接: ${NC}"
        if mysqladmin ping 2>/dev/null; then
            echo -e "${GREEN}✓ 可以连接${NC}"
        else
            echo -e "${RED}✗ 连接失败${NC}"
            echo -e "  请检查配置文件中的 MySQL 凭据"
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 mysqladmin${NC}"
    fi
}

check_redis() {
    echo -e "${YELLOW}=== Redis 缓存检查 ===${NC}"

    # 检查 Redis 服务
    if command -v systemctl >/dev/null 2>&1; then
        if systemctl is-active --quiet redis-server || systemctl is-active --quiet redis; then
            echo -e "${GREEN}✓ Redis 服务运行中${NC}"
        elif systemctl is-enabled --quiet redis-server || systemctl is-enabled --quiet redis; then
            echo -e "${YELLOW}⚠ Redis 服务已启用但未运行${NC}"
        else
            echo -e "${YELLOW}⚠ Redis 服务未在 systemd 中配置${NC}"
        fi
    fi

    # 检查 Redis 连接
    if command -v redis-cli >/dev/null 2>&1; then
        echo -n -e "${YELLOW}检查 Redis 连接: ${NC}"
        if redis-cli ping 2>/dev/null; then
            echo -e "${GREEN}✓ 可以连接 (PONG)${NC}"

            # 检查 Redis 内存使用
            memory_info=$(redis-cli INFO memory 2>/dev/null | grep used_memory_human)
            if [ -n "$memory_info" ]; then
                echo -e "  $memory_info"
            fi
        else
            echo -e "${RED}✗ 连接失败${NC}"
            echo -e "  请检查配置文件中的 Redis 地址"
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 redis-cli${NC}"
    fi
}

check_rabbitmq() {
    echo -e "${YELLOW}=== RabbitMQ 消息队列检查 ===${NC}"

    # 检查 RabbitMQ 服务
    if command -v systemctl >/dev/null 2>&1; then
        if systemctl is-active --quiet rabbitmq-server; then
            echo -e "${GREEN}✓ RabbitMQ 服务运行中${NC}"
        elif systemctl is-enabled --quiet rabbitmq-server; then
            echo -e "${YELLOW}⚠ RabbitMQ 服务已启用但未运行${NC}"
        else
            echo -e "${YELLOW}⚠ RabbitMQ 服务未在 systemd 中配置${NC}"
        fi
    fi

    # 检查 RabbitMQ 连接
    if command -v rabbitmq-diagnostics >/dev/null 2>&1; then
        echo -n -e "${YELLOW}检查 RabbitMQ 连接: ${NC}"
        if rabbitmq-diagnostics ping 2>/dev/null; then
            echo -e "${GREEN}✓ 可以连接${NC}"
        else
            echo -e "${RED}✗ 连接失败${NC}"
            echo -e "  请检查配置文件中的 RabbitMQ 凭据"
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 rabbitmq-diagnostics${NC}"
    fi
}

check_application() {
    echo -e "${YELLOW}=== 应用程序检查 ===${NC}"

    # 检查应用进程
    if pgrep -q "simple-demo"; then
        echo -e "${GREEN}✓ 应用进程运行中${NC}"
        pgrep "simple-demo" -a
    else
        echo -e "${RED}✗ 应用进程未运行${NC}"
    fi

    # 检查二进制文件
    if [ -f "./bin/simple-demo" ]; then
        echo -e "${GREEN}✓ 二进制文件存在${NC}"
        ls -lh ./bin/simple-demo
    else
        echo -e "${RED}✗ 二进制文件不存在${NC}"
        echo -e "  请运行 make build 构建应用"
    fi
}

check_api_health() {
    echo -e "${YELLOW}=== API 健康检查 ===${NC}"

    # 检查应用健康端点
    echo -n -e "${YELLOW}检查 API 健康端点: ${NC}"
    if command -v curl >/dev/null 2>&1; then
        response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null)
        if [ "$response" = "200" ]; then
            echo -e "${GREEN}✓ 健康检查通过 (200 OK)${NC}"
        else
            echo -e "${RED}✗ 健康检查失败 ($response)${NC}"
        fi

        # 检查 Feed API
        echo -n -e "${YELLOW}检查 Feed API: ${NC}"
        feed_response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/douyin/feed/ 2>/dev/null)
        if [ "$feed_response" = "200" ]; then
            echo -e "${GREEN}✓ Feed API 正常 (200 OK)${NC}"
        else
            echo -e "${YELLOW}⚠ Feed API 返回 $feed_response${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ 未安装 curl${NC}"
    fi
}

check_logs() {
    echo -e "${YELLOW}=== 日志文件检查 ===${NC}"

    # 检查日志目录
    if [ -d "./logs" ]; then
        echo -e "${GREEN}✓ 日志目录存在${NC}"

        # 检查最近的错误
        if [ -f "./logs/error.log" ]; then
            error_count=$(tail -n 50 ./logs/error.log | grep -i "error\|panic\|fatal" | wc -l)
            echo -e "  最近 50 行日志中的错误数: $error_count"

            if [ $error_count -gt 0 ]; then
                echo -e "${RED}最近的错误：${NC}"
                tail -n 5 ./logs/error.log | grep -i "error\|panic\|fatal"
            fi
        else
            echo -e "${YELLOW}⚠ 错误日志文件不存在${NC}"
        fi

        # 检查日志文件大小
        total_size=$(du -sh ./logs | cut -f1)
        echo -e "  日志目录总大小: $total_size"
    else
        echo -e "${YELLOW}⚠ 日志目录不存在${NC}"
        echo -e "  日志将在首次运行时创建"
    fi
}

check_configuration() {
    echo -e "${YELLOW}=== 配置文件检查 ===${NC}"

    # 检查配置文件
    if [ -f "./config/settings.yml" ]; then
        echo -e "${GREEN}✓ 配置文件存在${NC}"

        # 检查敏感配置
        echo -e "${YELLOW}敏感配置检查：${NC}"

        # 检查 JWT 密钥
        if grep -q "secretKey: 123456" ./config/settings.yml 2>/dev/null; then
            echo -e "${RED}⚠ 使用了默认 JWT 密钥 (123456)${NC}"
            echo -e "  请运行 ./scripts/setup_jwt.sh 生成安全的密钥"
        fi

        # 检查数据库密码
        if grep -q "password: \"password\"" ./config/settings.yml 2>/dev/null; then
            echo -e "${YELLOW}⚠ 使用了默认数据库密码${NC}"
            echo -e "  请在生产环境中修改数据库密码"
        fi

        # 检查 RabbitMQ 默认凭据
        if grep -q "username: \"guest\"" ./config/settings.yml 2>/dev/null; then
            echo -e "${YELLOW}⚠ 使用了默认 RabbitMQ 凭据${NC}"
            echo -e "  请在生产环境中修改 RabbitMQ 凭据"
        fi
    else
        echo -e "${RED}✗ 配置文件不存在${NC}"
        echo -e "  请运行 cp config/settings.yml.template config/settings.yml"
    fi
}

check_dependencies() {
    echo -e "${YELLOW}=== 依赖项检查 ===${NC}"

    # 检查 Go 版本
    if command -v go >/dev/null 2>&1; then
        go_version=$(go version | awk '{print $3}')
        echo -e "${GREEN}✓ Go 版本: $go_version${NC}"

        # 检查 Go 版本是否满足要求
        go_major=$(echo $go_version | cut -d. -f1)
        if [ "$go_major" -lt 1 ] || [ "$go_major" -eq 1 ] && [ "$(echo $go_version | cut -d. -f2)" -lt 20 ]; then
            echo -e "${RED}⚠ Go 版本过低 (需要 1.20+ )${NC}"
        else
            echo -e "${GREEN}✓ Go 版本满足要求${NC}"
        fi
    else
        echo -e "${RED}✗ Go 未安装${NC}"
    fi

    # 检查 go.mod
    if [ -f "./go.mod" ]; then
        echo -e "${GREEN}✓ go.mod 文件存在${NC}"
    else
        echo -e "${RED}✗ go.mod 文件不存在${NC}"
    fi

    # 检查 go.sum
    if [ -f "./go.sum" ]; then
        echo -e "${GREEN}✓ go.sum 文件存在${NC}"
    else
        echo -e "${RED}✗ go.sum 文件不存在${NC}"
        echo -e "  请运行 go mod download"
    fi
}

# 主检查流程

echo -e "${GREEN}=== 端口检查 ===${NC}"
check_port 8080 "API 服务"
        check_port 3306 "MySQL 数据库"
        check_port 6379 "Redis 缓存"
        check_port 5672 "RabbitMQ"

echo ""
check_dependencies
echo ""
check_configuration
echo ""
check_application
echo ""
check_mysql
echo ""
check_redis
echo ""
check_rabbitmq
echo ""
check_api_health
echo ""
check_logs
echo ""
check_disk
echo ""
check_memory

echo ""
echo -e "${GREEN}=== 健康检查完成 ===${NC}"
echo ""
echo -e "${YELLOW}如需更多帮助，请查看：${NC}"
echo -e "  - 文档: docs/TROUBLESHOOTING.md"
echo -e "  - FAQ: docs/FAQ.md"
echo -e "  - 支持: https://github.com/yourusername/ByteDanceDemo/issues"
echo ""
