#!/bin/bash
# quick-check.sh - Quick diagnostic script for ByteDanceDemo
# Usage: bash docs/scripts/quick-check.sh

echo "=== ByteDanceDemo Quick Diagnostics ==="
echo ""

# Check if application is running
echo "1. Application Status:"
if pgrep -f "simple-demo" > /dev/null; then
    echo "✓ Application is running (PID: $(pgrep -f 'simple-demo'))"
else
    echo "✗ Application is NOT running"
fi

# Check port availability
echo ""
echo "2. Port Availability:"
if netstat -tuln 2>/dev/null | grep -q ":8080"; then
    echo "✓ Port 8080 is in use"
else
    echo "✗ Port 8080 is NOT in use"
fi

# Check database connection
echo ""
echo "3. Database Connection:"
if command -v mysql &> /dev/null; then
    if mysql -u sample_douyin -psample_douyin -e "SELECT 1;" &> /dev/null; then
        echo "✓ MySQL connection successful"
    else
        echo "✗ MySQL connection failed (check credentials and service)"
    fi
else
    echo "⚠ MySQL client not installed - cannot check database connection"
fi

# Check Redis connection
echo ""
echo "4. Redis Connection:"
if command -v redis-cli &> /dev/null; then
    if redis-cli ping &> /dev/null; then
        echo "✓ Redis connection successful"
    else
        echo "✗ Redis connection failed (check service is running)"
    fi
else
    echo "⚠ Redis CLI not installed - cannot check Redis connection"
fi

# Check RabbitMQ connection
echo ""
echo "5. RabbitMQ Connection:"
if command -v rabbitmq-diagnostics &> /dev/null; then
    if rabbitmq-diagnostics ping &> /dev/null; then
        echo "✓ RabbitMQ connection successful"
    else
        echo "✗ RabbitMQ connection failed (check service is running)"
    fi
else
    echo "⚠ RabbitMQ diagnostics not installed - cannot check connection"
fi

# Check disk space
echo ""
echo "6. Disk Space:"
df -h | grep -E "Filesystem|/$|/var|/home|/opt" | while read line; do
    echo "  $line"
done

# Check memory usage
echo ""
echo "7. Memory Usage:"
if command -v free &> /dev/null; then
    free -h
else
    echo "⚠ 'free' command not available"
fi

# Check CPU usage
echo ""
echo "8. CPU Load:"
if command -v uptime &> /dev/null; then
    uptime
else
    echo "⚠ 'uptime' command not available"
fi

# Check recent errors
echo ""
echo "9. Recent Application Errors:"
if [ -f "/var/log/bytedancedemo/error.log" ]; then
    recent_errors=$(tail -20 /var/log/bytedancedemo/error.log | grep -i "error")
    if [ -n "$recent_errors" ]; then
        echo "Recent errors found:"
        echo "$recent_errors"
    else
        echo "No recent errors found in error log"
    fi
else
    echo "Error log file not found at /var/log/bytedancedemo/error.log"
fi

# Check Go version
echo ""
echo "10. Go Environment:"
if command -v go &> /dev/null; then
    echo "✓ Go version: $(go version)"
else
    echo "✗ Go is not installed or not in PATH"
fi

echo ""
echo "=== Diagnostics Complete ==="
echo ""
echo "For detailed troubleshooting, see: docs/TROUBLESHOOTING.md"
