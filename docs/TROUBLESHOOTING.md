# ByteDanceDemo - Troubleshooting Guide

## Table of Contents

1. [Quick Diagnostics](#quick-diagnostics)
2. [Common Issues](#common-issues)
3. [Database Issues](#database-issues)
4. [API Issues](#api-issues)
5. [Performance Issues](#performance-issues)
6. [Deployment Issues](#deployment-issues)
7. [Docker Issues](#docker-issues)
8. [Security Issues](#security-issues)
9. [Logging and Debugging](#logging-and-debugging)
10. [Emergency Procedures](#emergency-procedures)

## Quick Diagnostics

### Health Check Script

```bash
#!/bin/bash
# quick-check.sh - Quick diagnostic script

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
if netstat -tuln | grep -q ":8080"; then
    echo "✓ Port 8080 is in use"
else
    echo "✗ Port 8080 is NOT in use"
fi

# Check database connection
echo ""
echo "3. Database Connection:"
if mysql -u sample_douyin -psample_douyin -e "SELECT 1;" &> /dev/null; then
    echo "✓ MySQL connection successful"
else
    echo "✗ MySQL connection failed"
fi

# Check Redis connection
echo ""
echo "4. Redis Connection:"
if redis-cli ping &> /dev/null; then
    echo "✓ Redis connection successful"
else
    echo "✗ Redis connection failed"
fi

# Check RabbitMQ connection
echo ""
echo "5. RabbitMQ Connection:"
if rabbitmq-diagnostics ping &> /dev/null; then
    echo "✓ RabbitMQ connection successful"
else
    echo "✗ RabbitMQ connection failed"
fi

# Check disk space
echo ""
echo "6. Disk Space:"
df -h | grep -E "/$|/var|/home" | while read line; do
    echo "  $line"
done

# Check memory usage
echo ""
echo "7. Memory Usage:"
free -h

# Check recent errors
echo ""
echo "8. Recent Application Errors:"
if [ -f "/var/log/bytedancedemo/error.log" ]; then
    tail -20 /var/log/bytedancedemo/error.log | grep -i "error" || echo "No recent errors found"
else
    echo "Error log file not found"
fi

echo ""
echo "=== Diagnostics Complete ==="
```

### System Status Overview

```bash
# Check all system services
systemctl status bytedancedemo mysql redis-server rabbitmq-server

# Check application logs
journalctl -u bytedancedemo -n 50 --no-pager

# Check database processes
mysqladmin processlist -u sample_douyin -p

# Check Redis memory usage
redis-cli INFO memory

# Check RabbitMQ queues
rabbitmqctl list_queues name messages consumers
```

## Common Issues

### Application Won't Start

#### Symptoms
- Application fails to start
- Error: "address already in use"
- Error: "connection refused"

#### Diagnosis

```bash
# Check if port is already in use
sudo lsof -i :8080
sudo netstat -tuln | grep 8080

# Check application logs
tail -50 /var/log/bytedancedemo/error.log
journalctl -u bytedancedemo -n 100

# Check configuration file syntax
./bin/simple-demo server -c config/settings.yml -m debug
```

#### Solutions

```bash
# Solution 1: Kill existing process
sudo kill -9 $(lsof -t -i:8080)
sudo systemctl restart bytedancedemo

# Solution 2: Change port in configuration
# Edit config/settings.yml
settings:
  application:
    port: 8081  # Change to different port

# Solution 3: Check file permissions
sudo chown -R appuser:appuser /opt/bytedancedemo
sudo chmod 755 /opt/bytedancedemo/bin/simple-demo
```

### Configuration Errors

#### Symptoms
- "configuration file not found"
- "missing required field"
- "invalid configuration format"

#### Diagnosis

```bash
# Check if configuration file exists
ls -la config/settings.yml

# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('config/settings.yml'))"

# Check for environment variables
env | grep -E "MYSQL|REDIS|RABBITMQ|JWT"
```

#### Solutions

```bash
# Solution 1: Create missing configuration
cp config/settings.yml.template config/settings.yml

# Solution 2: Set missing environment variables
export MYSQL_PASSWORD="your_password"
export REDIS_PASSWORD="your_password"
export JWT_SECRET="your_secret_key"

# Solution 3: Fix YAML indentation
# Use proper YAML formatting (2 spaces indentation)
```

### Memory Issues

#### Symptoms
- Out of memory errors
- High memory usage
- Application crashes

#### Diagnosis

```bash
# Check memory usage
free -h
ps aux --sort=-%mem | head

# Check Go memory profiling
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Monitor memory over time
watch -n 1 'free -h'
```

#### Solutions

```bash
# Solution 1: Increase system memory or swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Solution 2: Tune MySQL memory settings
# Edit /etc/mysql/mysql.conf.d/mysqld.cnf
innodb_buffer_pool_size = 1G  # Reduce if needed

# Solution 3: Optimize Go garbage collection
export GOGC=100  # Adjust GC percentage
```

## Database Issues

### MySQL Connection Issues

#### Symptoms
- "connection refused"
- "access denied"
- "too many connections"

#### Diagnosis

```bash
# Test MySQL connection
mysql -u sample_douyin -p -h localhost

# Check MySQL status
sudo systemctl status mysql

# Check MySQL logs
sudo tail -50 /var/log/mysql/error.log

# Check connection count
mysql -u root -p -e "SHOW PROCESSLIST;"
```

#### Solutions

```bash
# Solution 1: Restart MySQL service
sudo systemctl restart mysql

# Solution 2: Check credentials
mysql -u sample_douyin -p -e "SELECT CURRENT_USER();"

# Solution 3: Increase connection limit
mysql -u root -p -e "SET GLOBAL max_connections = 200;"

# Solution 4: Grant necessary permissions
mysql -u root -p -e "
GRANT ALL PRIVILEGES ON sample_douyin.* TO 'sample_douyin'@'localhost';
FLUSH PRIVILEGES;"
```

### Slow Database Queries

#### Symptoms
- Slow API responses
- High database load
- Timeout errors

#### Diagnosis

```bash
# Enable slow query logging
mysql -u root -p -e "
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
SET GLOBAL log_queries_not_using_indexes = 'ON';
"

# Check slow queries
sudo tail -f /var/log/mysql/slow-query.log

# Analyze query performance
mysql -u sample_douyin -p -e "SHOW FULL PROCESSLIST;"

# Check table locks
mysql -u sample_douyin -p -e "SHOW OPEN TABLES WHERE In_use > 0;"
```

#### Solutions

```bash
# Solution 1: Add database indexes
mysql -u sample_douyin -p sample_douyin -e "
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_videos_author ON videos(author_id);
CREATE INDEX idx_comments_video ON comments(video_id);
"

# Solution 2: Optimize tables
mysql -u sample_douyin -p sample_douyin -e "OPTIMIZE TABLE users, videos, comments;"

# Solution 3: Update MySQL configuration
# Edit /etc/mysql/mysql.conf.d/mysqld.cnf
query_cache_type = 1
query_cache_size = 64M
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
```

### Database Migration Issues

#### Symptoms
- "migration failed"
- "table already exists"
- "foreign key constraint fails"

#### Diagnosis

```bash
# Check migration status
mysql -u sample_douyin -p sample_douyin -e "SELECT * FROM schema_migrations;"

# Check table structure
mysql -u sample_douyin -p sample_douyin -e "DESCRIBE users;"

# Check for orphaned records
mysql -u sample_douyin -p sample_douyin -e "
SELECT videos.id FROM videos 
LEFT JOIN users ON videos.author_id = users.id 
WHERE users.id IS NULL;"
```

#### Solutions

```bash
# Solution 1: Drop and recreate database
mysql -u root -p -e "
DROP DATABASE IF EXISTS sample_douyin;
CREATE DATABASE sample_douyin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
"

# Solution 2: Reset migration
mysql -u sample_douyin -p sample_douyin -e "DROP TABLE IF EXISTS schema_migrations;"

# Solution 3: Fix orphaned records
mysql -u sample_douyin -p sample_douyin -e "
DELETE FROM videos WHERE author_id NOT IN (SELECT id FROM users);"
```

## API Issues

### 500 Internal Server Error

#### Symptoms
- Generic 500 errors
- No detailed error messages
- Application crashes

#### Diagnosis

```bash
# Check application logs
tail -f /var/log/bytedancedemo/error.log
journalctl -u bytedancedemo -f

# Check HTTP errors
tail -f /var/log/nginx/bytedancedemo_error.log

# Check for Go panics
grep -r "panic:" /var/log/bytedancedemo/
```

#### Solutions

```bash
# Solution 1: Enable debug mode
./bin/simple-demo server -c config/settings.yml -m debug

# Solution 2: Check error handlers
# Review middleware/error.go for proper error handling

# Solution 3: Add detailed logging
# Add logging to controller methods
zap.L().Error("Error processing request",
    zap.Error(err),
    zap.String("path", c.Request.URL.Path),
)
```

### 401 Unauthorized Errors

#### Symptoms
- "invalid token"
- "token expired"
- "unauthorized access"

#### Diagnosis

```bash
# Test with valid token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/douyin/user/

# Check JWT configuration
grep -A 5 "jwt:" config/settings.yml

# Verify token expiration
echo "YOUR_TOKEN" | base64 -d | jq .
```

#### Solutions

```bash
# Solution 1: Generate new token
curl -X POST "http://localhost:8080/douyin/user/login/?username=user&password=pass"

# Solution 2: Check JWT secret
echo "Verify JWT_SECRET matches: $JWT_SECRET"

# Solution 3: Extend token expiration
# Edit config/settings.yml
settings:
  jwt:
    expirationTime: 48h  # Extend to 48 hours
```

### Rate Limiting Issues

#### Symptoms
- "rate limit exceeded"
- "too many requests"
- Requests blocked

#### Diagnosis

```bash
# Check rate limit configuration
grep -A 5 "rateLimit:" config/settings.yml

# Check Redis rate limit keys
redis-cli KEYS "ratelimit:*"
redis-cli TTL "ratelimit:USER_IP"

# Check rate limit status
curl -I http://localhost:8080/douyin/feed/
```

#### Solutions

```bash
# Solution 1: Clear rate limit cache
redis-cli --scan --pattern "ratelimit:*" | xargs redis-cli DEL

# Solution 2: Increase rate limit
# Edit config/settings.yml
settings:
  application:
    rateLimit: 200  # Increase from default

# Solution 3: Whitelist trusted IPs
# Add to rate limiter configuration
whitelist_ips:
  - "192.168.1.100"
  - "10.0.0.1"
```

## Performance Issues

### Slow API Responses

#### Symptoms
- High response times
- Timeouts
- Poor user experience

#### Diagnosis

```bash
# Measure response time
time curl http://localhost:8080/douyin/feed/

# Check CPU usage
top -p $(pgrep -f simple-demo)

# Profile application performance
curl http://localhost:8080/debug/pprof/profile?seconds=30 > profile.prof
go tool pprof profile.prof

# Check database query times
grep "Query Time" /var/log/bytedancedemo/error.log
```

#### Solutions

```bash
# Solution 1: Enable caching
# Add Redis caching for frequent queries
cachedUser, err := cache.Get("user:" + userID)
if err != nil {
    user, err := db.GetUser(userID)
    cache.Set("user:"+userID, user, 5*time.Minute)
}

# Solution 2: Optimize database queries
# Use SELECT with specific fields
db.Select("id", "name", "email").First(&user)

# Solution 3: Implement connection pooling
# Edit config/settings.yml
settings:
  mysql:
    maxOpenConns: 100
    maxIdleConns: 20
    connMaxLifetime: 5m
```

### High CPU Usage

#### Symptoms
- High CPU utilization
- Slow application performance
- System slowdown

#### Diagnosis

```bash
# Check CPU usage
top -p $(pgrep -f simple-demo)

# Profile CPU usage
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# Check goroutine count
curl http://localhost:8080/debug/pprof/goroutine?debug=1 | wc -l
```

#### Solutions

```bash
# Solution 1: Reduce worker count
# Edit config/settings.yml
settings:
  application:
    workerCount: 2  # Reduce from 4

# Solution 2: Optimize algorithms
# Use more efficient data structures
# Implement pagination
db.Limit(10).Offset(page * 10).Find(&videos)

# Solution 3: Add caching layers
// Implement Redis caching for expensive operations
func GetCachedData(key string) (interface{}, error) {
    cached, err := redis.Get(key)
    if err == nil {
        return cached, nil
    }
    data, err := computeExpensiveOperation()
    redis.Set(key, data, 1*time.Hour)
    return data, err
}
```

### Memory Leaks

#### Symptoms
- Increasing memory usage over time
- Out of memory errors
- Application crashes

#### Diagnosis

```bash
# Monitor memory over time
watch -n 10 'ps aux | grep simple-demo | awk "{print $6}"'

# Generate heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Check for unreleased connections
netstat -an | grep ESTABLISHED | grep 8080 | wc -l
```

#### Solutions

```bash
# Solution 1: Tune garbage collection
export GOGC=50  # More frequent GC

# Solution 2: Check for unclosed connections
// Ensure database connections are closed
defer db.Close()
defer resp.Body.Close()

# Solution 3: Implement connection pooling
// Use database connection pool properly
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(20)
```

## Deployment Issues

### Build Failures

#### Symptoms
- "build failed"
- "module not found"
- "compilation errors"

#### Diagnosis

```bash
# Check Go version
go version

# Check dependencies
go mod verify
go mod tidy

# Check build output
go build -v 2>&1 | tee build.log

# Check for errors
grep -i "error" build.log
```

#### Solutions

```bash
# Solution 1: Clean build cache
go clean -cache
go clean -modcache

# Solution 2: Update dependencies
go get -u ./...
go mod tidy

# Solution 3: Check for conflicting versions
go list -m all | grep "duplicates"

# Solution 4: Rebuild from scratch
rm -rf bin/
go build -o bin/simple-demo
```

### Service Startup Failures

#### Symptoms
- "service failed to start"
- "timeout waiting for service"
- Service crashes immediately

#### Diagnosis

```bash
# Check service status
sudo systemctl status bytedancedemo

# Check service logs
sudo journalctl -u bytedancedemo -n 100 --no-pager

# Check service configuration
sudo systemd-analyze verify /etc/systemd/system/bytedancedemo.service

# Test manual start
sudo -u appuser /opt/bytedancedemo/bin/simple-demo server \
  -c /opt/bytedancedemo/config/settings.yml -m debug
```

#### Solutions

```bash
# Solution 1: Check service dependencies
sudo systemctl status mysql redis-server rabbitmq-server
sudo systemctl start mysql redis-server rabbitmq-server

# Solution 2: Fix file permissions
sudo chown -R appuser:appuser /opt/bytedancedemo
sudo chmod 755 /opt/bytedancedemo/bin/simple-demo
sudo chmod 600 /opt/bytedancedemo/config/settings.yml

# Solution 3: Update service configuration
sudo nano /etc/systemd/system/bytedancedemo.service
sudo systemctl daemon-reload
sudo systemctl restart bytedancedemo
```

### Nginx Issues

#### Symptoms
- "502 Bad Gateway"
- "503 Service Unavailable"
- Nginx connection errors

#### Diagnosis

```bash
# Check Nginx status
sudo systemctl status nginx

# Check Nginx configuration
sudo nginx -t

# Check Nginx error logs
sudo tail -50 /var/log/nginx/error.log

# Check upstream connectivity
curl http://localhost:8080/health
```

#### Solutions

```bash
# Solution 1: Restart Nginx
sudo systemctl restart nginx

# Solution 2: Fix Nginx configuration
sudo nano /etc/nginx/sites-available/bytedancedemo
sudo nginx -t
sudo systemctl reload nginx

# Solution 3: Check upstream configuration
# Ensure application is running on correct port
netstat -tuln | grep 8080

# Solution 4: Increase Nginx timeout
# Add to upstream configuration
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;
```

## Docker Issues

### Container Won't Start

#### Symptoms
- "container exited"
- "error pulling image"
- "container not found"

#### Diagnosis

```bash
# Check container status
docker ps -a

# Check container logs
docker logs bytedancedemo_app

# Check Docker daemon
sudo systemctl status docker

# Check disk space
df -h
```

#### Solutions

```bash
# Solution 1: Restart Docker
sudo systemctl restart docker

# Solution 2: Remove and recreate container
docker-compose down
docker-compose up -d

# Solution 3: Check for port conflicts
docker ps
# Change ports in docker-compose.yml

# Solution 4: Rebuild image
docker-compose build --no-cache
docker-compose up -d
```

### Docker Network Issues

#### Symptoms
- "network not found"
- "connection refused"
- Services can't communicate

#### Diagnosis

```bash
# Check Docker networks
docker network ls

# Check container network settings
docker inspect bytedancedemo_app | grep -A 20 "Networks"

# Test connectivity between containers
docker exec bytedancedemo_app ping mysql
docker exec bytedancedemo_app ping redis
```

#### Solutions

```bash
# Solution 1: Recreate network
docker-compose down
docker network prune
docker-compose up -d

# Solution 2: Check service names
# Use service names as hostnames in configuration
mysql:
  host: mysql  # Not localhost
redis:
  addr: redis:6379  # Not localhost:6379

# Solution 3: Check firewall rules
sudo ufw allow from 172.16.0.0/12  # Allow Docker network
```

### Volume Mount Issues

#### Symptoms
- "permission denied"
- Data not persisting
- Volume not accessible

#### Diagnosis

```bash
# Check volume mounts
docker inspect bytedancedemo_app | grep -A 10 "Mounts"

# Check volume contents
docker exec bytedancedemo_app ls -la /app/logs

# Check host permissions
ls -la /opt/bytedancedemo/logs
```

#### Solutions

```bash
# Solution 1: Fix host permissions
sudo chown -R $USER:$USER /opt/bytedancedemo/logs
sudo chmod -R 755 /opt/bytedancedemo/logs

# Solution 2: Use named volumes
# In docker-compose.yml:
volumes:
  logs_data:
    driver: local

# Solution 3: Run as non-root user
# In Dockerfile:
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
```

## Security Issues

### Unauthorized Access

#### Symptoms
- Unknown users logged in
- Unusual API calls
- Data tampering

#### Diagnosis

```bash
# Check access logs
grep "401\|403" /var/log/bytedancedemo/access.log | tail -50

# Check authentication logs
grep "failed login\|invalid token" /var/log/bytedancedemo/error.log

# Check for suspicious IPs
awk '{print $1}' /var/log/bytedancedemo/access.log | sort | uniq -c | sort -rn
```

#### Solutions

```bash
# Solution 1: Review access logs
grep "suspicious_activity" /var/log/bytedancedemo/error.log

# Solution 2: Change secrets
# Regenerate JWT secret
export JWT_SECRET=$(openssl rand -base64 32)

# Solution 3: Block malicious IPs
sudo iptables -A INPUT -s SUSPICIOUS_IP -j DROP

# Solution 4: Rotate credentials
# Update all database and service passwords
```

### Security Headers Issues

#### Symptoms
- Security warnings in browser
- Content not loading
- Mixed content errors

#### Diagnosis

```bash
# Check HTTP headers
curl -I http://localhost:8080/douyin/feed/

# Check SSL configuration
curl -I https://api.example.com/douyin/feed/

# Check browser console for warnings
# Look for CSP violations, mixed content
```

#### Solutions

```bash
# Solution 1: Fix Content Security Policy
# Update in middleware/security.go
c.Header("Content-Security-Policy",
    "default-src 'self'; " +
    "script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
    "style-src 'self' 'unsafe-inline'; " +
    "img-src 'self' data: https:;")

# Solution 2: Fix mixed content
# Ensure all resources use HTTPS
<img src="https://cdn.example.com/image.jpg" />

# Solution 3: Update CORS settings
# Update in middleware/cors.go
allowOrigins: []string{"https://yourdomain.com"}
```

## Logging and Debugging

### Enable Debug Mode

```bash
# Start application with debug mode
./bin/simple-demo server -c config/settings.yml -m debug

# Or update configuration
# Edit config/settings.yml
settings:
  log:
    level: "debug"
    development: true
```

### Check Specific Logs

```bash
# Application logs
tail -f /var/log/bytedancedemo/bytedancedemo.log

# Error logs
tail -f /var/log/bytedancedemo/error.log

# Access logs
tail -f /var/log/nginx/bytedancedemo_access.log

# Database logs
tail -f /var/log/mysql/error.log

# Redis logs
tail -f /var/log/redis/redis-server.log
```

### Generate Diagnostic Report

```bash
#!/bin/bash
# generate-report.sh

REPORT_FILE="diagnostic-report-$(date +%Y%m%d-%H%M%S).txt"

{
    echo "=== ByteDanceDemo Diagnostic Report ==="
    echo "Generated: $(date)"
    echo ""

    echo "=== System Information ==="
    uname -a
    echo ""

    echo "=== Application Status ==="
    systemctl status bytedancedemo
    echo ""

    echo "=== Recent Application Logs ==="
    tail -50 /var/log/bytedancedemo/bytedancedemo.log
    echo ""

    echo "=== Error Logs ==="
    tail -50 /var/log/bytedancedemo/error.log
    echo ""

    echo "=== Database Status ==="
    systemctl status mysql
    mysql -u sample_douyin -p -e "SELECT NOW();"
    echo ""

    echo "=== Redis Status ==="
    redis-cli INFO server
    echo ""

    echo "=== System Resources ==="
    free -h
    df -h
    echo ""

    echo "=== Network Connections ==="
    netstat -tuln | grep -E "8080|3306|6379|5672"
    echo ""

    echo "=== Recent Processes ==="
    ps aux | grep -E "simple-demo|mysql|redis|rabbitmq"

} > "$REPORT_FILE"

echo "Diagnostic report generated: $REPORT_FILE"
```

## Emergency Procedures

### Application Not Responding

```bash
# 1. Check if process is running
ps aux | grep simple-demo

# 2. Check system resources
free -h
top

# 3. Restart application
sudo systemctl restart bytedancedemo

# 4. If restart fails, force kill and restart
sudo pkill -9 simple-demo
sudo systemctl start bytedancedemo

# 5. Check logs for cause
journalctl -u bytedancedemo -n 100 --no-pager
```

### Database Connection Lost

```bash
# 1. Check MySQL status
sudo systemctl status mysql

# 2. Test connection
mysql -u sample_douyin -p -e "SELECT 1;"

# 3. Check MySQL logs
tail -50 /var/log/mysql/error.log

# 4. Restart MySQL if needed
sudo systemctl restart mysql

# 5. Check connection pool
mysql -u root -p -e "SHOW PROCESSLIST;"
```

### Complete System Failure

```bash
# 1. Assess situation
# Check what services are down
sudo systemctl status bytedancedemo mysql redis-server rabbitmq-server nginx

# 2. Restart all services in order
sudo systemctl restart mysql
sudo systemctl restart redis-server
sudo systemctl restart rabbitmq-server
sudo systemctl restart bytedancedemo
sudo systemctl restart nginx

# 3. Verify services are running
sudo systemctl status bytedancedemo mysql redis-server rabbitmq-server nginx

# 4. Test application
curl http://localhost:8080/health

# 5. Check logs
journalctl -u bytedancedemo -f
```

### Data Recovery

```bash
# 1. Stop application to prevent further damage
sudo systemctl stop bytedancedemo

# 2. Restore from backup
gunzip < /backup/mysql/backup_latest.sql.gz | \
    mysql -u sample_douyin -p sample_douyin

# 3. Verify data integrity
mysql -u sample_douyin -p sample_douyin -e "
SELECT COUNT(*) as user_count FROM users;
SELECT COUNT(*) as video_count FROM videos;
"

# 4. Start application
sudo systemctl start bytedancedemo

# 5. Monitor for issues
journalctl -u bytedancedemo -f
```

### Getting Help

If you can't resolve an issue:

1. **Check documentation**: Review all documentation files in `/docs`
2. **Search logs**: Look for error messages and stack traces
3. **Check GitHub issues**: Search for similar problems
4. **Create diagnostic report**: Use the script above
5. **Contact support**: Include diagnostic report with your request

---

*This troubleshooting guide is regularly updated with new solutions based on community feedback.*
