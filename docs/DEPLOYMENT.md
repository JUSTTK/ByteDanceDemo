# ByteDanceDemo - Deployment Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Setup](#development-setup)
3. [Production Deployment](#production-deployment)
4. [Docker Deployment](#docker-deployment)
5. [Cloud Deployment](#cloud-deployment)
6. [Monitoring and Maintenance](#monitoring-and-maintenance)
7. [Rollback Procedures](#rollback-procedures)

## Prerequisites

### System Requirements

#### Minimum Requirements
- **CPU**: 2 cores
- **Memory**: 4GB RAM
- **Disk**: 20GB free space
- **OS**: Linux (Ubuntu 20.04+ recommended), macOS, or Windows with WSL2

#### Recommended Requirements (Production)
- **CPU**: 4+ cores
- **Memory**: 8GB+ RAM
- **Disk**: 100GB+ SSD
- **Network**: 1Gbps connection

### Software Dependencies

#### Required Software
- **Go**: 1.20 or higher
- **MySQL**: 8.0 or higher
- **Redis**: 6.0 or higher
- **RabbitMQ**: 3.9 or higher

#### Optional Tools
- **Git**: Version control
- **Docker**: Containerization
- **Docker Compose**: Multi-container management
- **Nginx**: Reverse proxy and load balancer
- **Make**: Build automation

## Development Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/yourusername/ByteDanceDemo.git
cd ByteDanceDemo
```

### Step 2: Install Go Dependencies

```bash
# Ensure Go is installed
go version  # Should be go1.20 or higher

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### Step 3: Setup MySQL Database

```bash
# Install MySQL (Ubuntu/Debian)
sudo apt update
sudo apt install mysql-server -y

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure MySQL installation (optional)
sudo mysql_secure_installation
```

Create the database and user:

```bash
# Login to MySQL
mysql -u root -p

# In MySQL prompt:
CREATE DATABASE sample_douyin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'sample_douyin'@'localhost' IDENTIFIED BY 'sample_douyin';
GRANT ALL PRIVILEGES ON sample_douyin.* TO 'sample_douyin'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

Import the database schema:

```bash
mysql -u sample_douyin -psample_douyin sample_douyin < config/init.sql
```

### Step 4: Setup Redis

```bash
# Install Redis (Ubuntu/Debian)
sudo apt install redis-server -y

# Start Redis service
sudo systemctl start redis-server
sudo systemctl enable redis-server

# Verify Redis is running
redis-cli ping  # Should return PONG
```

### Step 5: Setup RabbitMQ

```bash
# Install RabbitMQ (Ubuntu/Debian)
sudo apt install rabbitmq-server -y

# Start RabbitMQ service
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# Enable management plugin (optional)
sudo rabbitmq-plugins enable rabbitmq_management

# Create admin user (optional)
sudo rabbitmqctl add_user admin admin
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

### Step 6: Configure Application

Edit the configuration file:

```bash
cp config/settings.yml.template config/settings.yml
nano config/settings.yml
```

Update the following settings:

```yaml
settings:
  application:
    rateLimit: 50  # Adjust based on your needs
  mysql:
    host: localhost
    port: 3306
    schema: sample_douyin
    username: sample_douyin
    password: sample_douyin  # Change in production
    logLevel: 1
  jwt:
    secretKey: your-super-secret-key-change-this  # Change in production
    expirationTime: 24
  redis:
    addr: localhost:6379
    password: ""
    expirationTime: 5
  rabbitMQ:
    host: localhost
    port: 5672
    username: guest
    password: guest  # Change in production
  log:
    path: ./logs
    level: info
    maxSize: 100
    maxAge: 30
    maxBackups: 10
    compress: false
    mode: debug
```

### Step 7: Build and Run

```bash
# Build the application
go build -o bin/simple-demo

# Run the application
./bin/simple-demo server -c config/settings.yml -m debug

# Or use make
make build
make run
```

### Step 8: Verify Installation

```bash
# Test the API
curl http://localhost:8080/douyin/feed/

# Expected response:
{
  "status_code": 0,
  "status_msg": "",
  "video_list": []
}
```

## Production Deployment

### Step 1: Security Hardening

#### Environment Variables
Set sensitive configuration via environment variables:

```bash
export JWT_SECRET="your-production-secret-key"
export MYSQL_PASSWORD="your-mysql-password"
export REDIS_PASSWORD="your-redis-password"
export RABBITMQ_PASSWORD="your-rabbitmq-password"
```

#### File Permissions
```bash
# Secure configuration files
chmod 600 config/settings.yml
chmod 700 logs/

# Set proper ownership
sudo chown -R appuser:appuser /opt/bytedancedemo
```

#### Firewall Configuration
```bash
# Configure UFW firewall
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### Step 2: Database Optimization

#### MySQL Configuration
Edit `/etc/mysql/mysql.conf.d/mysqld.cnf`:

```ini
[mysqld]
# Connection Settings
max_connections = 500
max_connect_errors = 100000

# InnoDB Settings
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT

# Query Cache
query_cache_type = 1
query_cache_size = 64M

# Binary Logging
log_bin = /var/log/mysql/mysql-bin.log
expire_logs_days = 7
max_binlog_size = 100M

# Slow Query Log
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow-query.log
long_query_time = 2
```

Restart MySQL:
```bash
sudo systemctl restart mysql
```

#### Redis Configuration
Edit `/etc/redis/redis.conf`:

```conf
# Memory Management
maxmemory 2gb
maxmemory-policy allkeys-lru

# Persistence
save 900 1
save 300 10
save 60 10000

# Security
requirepass your-redis-password
bind 127.0.0.1

# Logging
loglevel notice
logfile /var/log/redis/redis-server.log
```

Restart Redis:
```bash
sudo systemctl restart redis-server
```

### Step 3: Application Configuration

Create production configuration `/opt/bytedancedemo/config/settings.yml`:

```yaml
settings:
  application:
    rateLimit: 100  # Adjust based on traffic
  mysql:
    host: mysql-db.example.com
    port: 3306
    schema: sample_douyin
    username: bytedancedemo_user
    password: ${MYSQL_PASSWORD}  # Use environment variable
    logLevel: 2  # Error only
  jwt:
    secretKey: ${JWT_SECRET}  # Use environment variable
    expirationTime: 24
  oss:
    avatar: https://cdn.example.com/avatars/
    backgroundImage: https://cdn.example.com/backgrounds/
    signature: https://cdn.example.com/
  redis:
    addr: redis.example.com:6379
    password: ${REDIS_PASSWORD}
    expirationTime: 5
  rabbitMQ:
    host: rabbitmq.example.com
    port: 5672
    username: bytedancedemo
    password: ${RABBITMQ_PASSWORD}
  log:
    path: /var/log/bytedancedemo
    level: info
    maxSize: 100
    maxAge: 30
    maxBackups: 10
    compress: true
    mode: release
```

### Step 4: Systemd Service

Create systemd service file `/etc/systemd/system/bytedancedemo.service`:

```ini
[Unit]
Description=ByteDanceDemo API Server
After=network.target mysql.service redis.service rabbitmq-server.service

[Service]
Type=simple
User=bytedancedemo
Group=bytedancedemo
WorkingDirectory=/opt/bytedancedemo
ExecStart=/opt/bytedancedemo/bin/simple-demo server -c /opt/bytedancedemo/config/settings.yml -m release
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security Settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/bytedancedemo /opt/bytedancedemo/public

# Environment
Environment="JWT_SECRET=${JWT_SECRET}"
Environment="MYSQL_PASSWORD=${MYSQL_PASSWORD}"
Environment="REDIS_PASSWORD=${REDIS_PASSWORD}"
Environment="RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD}"

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable bytedancedemo
sudo systemctl start bytedancedemo
sudo systemctl status bytedancedemo
```

### Step 5: Nginx Reverse Proxy

Install and configure Nginx:

```bash
# Install Nginx
sudo apt install nginx -y

# Create site configuration
sudo nano /etc/nginx/sites-available/bytedancedemo
```

Nginx configuration:

```nginx
upstream bytedancedemo_backend {
    server 127.0.0.1:8080;
    # Add more servers for load balancing
    # server 127.0.0.1:8081;
    # server 127.0.0.1:8082;
}

server {
    listen 80;
    server_name api.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Static Files
    location /static/ {
        alias /opt/bytedancedemo/public/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # API Proxy
    location / {
        proxy_pass http://bytedancedemo_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        # Buffer Settings
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
    }

    # Health Check
    location /health {
        proxy_pass http://bytedancedemo_backend/health;
        access_log off;
    }

    # Logging
    access_log /var/log/nginx/bytedancedemo_access.log;
    error_log /var/log/nginx/bytedancedemo_error.log;
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/bytedancedemo /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Step 6: SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d api.example.com

# Auto-renewal
sudo certbot renew --dry-run
```

## Docker Deployment

### Dockerfile

Create `Dockerfile` in the project root:

```dockerfile
# Build Stage
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simple-demo .

# Runtime Stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Shanghai

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/simple-demo .

# Copy configuration
COPY config/settings.yml /app/config/

# Create necessary directories
RUN mkdir -p /app/public /app/logs && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./simple-demo", "server", "-c", "/app/config/settings.yml", "-m", "release"]
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  # MySQL Database
  mysql:
    image: mysql:8.0
    container_name: bytedancedemo_mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: sample_douyin
      MYSQL_USER: sample_douyin
      MYSQL_PASSWORD: sample_douyin
    volumes:
      - mysql_data:/var/lib/mysql
      - ./config/init.sql:/docker-entrypoint-initdb.d/init.sql
      - ./mysql/conf.d:/etc/mysql/conf.d
    ports:
      - "3306:3306"
    networks:
      - bytedancedemo_network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-prootpassword"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: bytedancedemo_redis
    restart: always
    command: redis-server --appendonly yes --requirepass redispassword
    volumes:
      - redis_data:/data
      - ./redis/redis.conf:/usr/local/etc/redis/redis.conf
    ports:
      - "6379:6379"
    networks:
      - bytedancedemo_network
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  # RabbitMQ Message Queue
  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: bytedancedemo_rabbitmq
    restart: always
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: adminpassword
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    ports:
      - "5672:5672"   # AMQP port
      - "15672:15672" # Management UI
    networks:
      - bytedancedemo_network
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Application
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: bytedancedemo_app
    restart: always
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    environment:
      - MYSQL_PASSWORD=sample_douyin
      - REDIS_PASSWORD=redispassword
      - RABBITMQ_PASSWORD=adminpassword
      - JWT_SECRET=your-production-secret-key
    volumes:
      - ./config/settings.yml:/app/config/settings.yml:ro
      - ./logs:/app/logs
      - ./public:/app/public
    ports:
      - "8080:8080"
    networks:
      - bytedancedemo_network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    container_name: bytedancedemo_nginx
    restart: always
    depends_on:
      - app
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./public:/usr/share/nginx/html/static
      - ./nginx/ssl:/etc/nginx/ssl
    ports:
      - "80:80"
      - "443:443"
    networks:
      - bytedancedemo_network

volumes:
  mysql_data:
    driver: local
  redis_data:
    driver: local
  rabbitmq_data:
    driver: local

networks:
  bytedancedemo_network:
    driver: bridge
```

### Build and Run with Docker

```bash
# Build the Docker images
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Check service status
docker-compose ps

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Cloud Deployment

### AWS Deployment

#### 1. EC2 Setup
```bash
# Launch EC2 instance (Ubuntu 20.04)
# Choose appropriate instance type (t3.medium for development, t3.xlarge for production)
# Configure security groups:
#   - SSH (22) from your IP
#   - HTTP (80) from 0.0.0.0/0
#   - HTTPS (443) from 0.0.0.0/0
```

#### 2. RDS MySQL Setup
```bash
# Create RDS MySQL instance
# Choose MySQL 8.0
# Configure instance class (db.t3.micro for development, db.t3.medium for production)
# Set up security group to allow EC2 access
# Note the endpoint URL
```

#### 3. ElastiCache Redis Setup
```bash
# Create ElastiCache Redis cluster
# Choose Redis 6.x or 7.x
# Configure node type (cache.t3.micro for development, cache.t3.medium for production)
# Set up security group to allow EC2 access
# Note the endpoint URL
```

#### 4. Amazon MQ (RabbitMQ) Setup
```bash
# Create Amazon MQ broker
# Choose RabbitMQ engine
# Configure broker instance type (mq.t3.micro for development, mq.t3.medium for production)
# Set up security group to allow EC2 access
# Note the endpoint URL
```

#### 5. S3 for Static Files
```bash
# Create S3 bucket for video storage
# Configure bucket policy for public read access
# Set up CORS configuration
# Use AWS SDK in application for file uploads
```

#### 6. Update Configuration
Update `config/settings.yml` with AWS endpoints:

```yaml
settings:
  mysql:
    host: your-rds-endpoint.rds.amazonaws.com
    port: 3306
    # ...
  redis:
    addr: your-elasticache-endpoint.abc123.use1.cache.amazonaws.com:6379
    # ...
  rabbitMQ:
    host: your-mq-endpoint.amazonaws.com
    port: 5671
    # ...
  oss:
    # Use S3 URLs
    avatar: https://your-bucket.s3.amazonaws.com/avatars/
    backgroundImage: https://your-bucket.s3.amazonaws.com/backgrounds/
    signature: https://your-bucket.s3.amazonaws.com/
```

### Kubernetes Deployment

#### 1. Create Kubernetes Manifests

`k8s/deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bytedancedemo
  labels:
    app: bytedancedemo
spec:
  replicas: 3
  selector:
    matchLabels:
      app: bytedancedemo
  template:
    metadata:
      labels:
        app: bytedancedemo
    spec:
      containers:
      - name: bytedancedemo
        image: your-registry/bytedancedemo:latest
        ports:
        - containerPort: 8080
        env:
        - name: MYSQL_HOST
          valueFrom:
            secretKeyRef:
              name: bytedancedemo-secrets
              key: mysql-host
        - name: MYSQL_PASSWORD
          valueFrom:
            secretKeyRef:
              name: bytedancedemo-secrets
              key: mysql-password
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: bytedancedemo-secrets
              key: redis-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: bytedancedemo-secrets
              key: jwt-secret
        volumeMounts:
        - name: config
          mountPath: /app/config
        - name: logs
          mountPath: /app/logs
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
      volumes:
      - name: config
        configMap:
          name: bytedancedemo-config
      - name: logs
        emptyDir: {}
```

`k8s/service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: bytedancedemo-service
spec:
  selector:
    app: bytedancedemo
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

#### 2. Deploy to Kubernetes

```bash
# Create secrets
kubectl create secret generic bytedancedemo-secrets \
  --from-literal=mysql-host=mysql.example.com \
  --from-literal=mysql-password=your-password \
  --from-literal=redis-password=your-password \
  --from-literal=jwt-secret=your-secret

# Create config map
kubectl create configmap bytedancedemo-config \
  --from-file=config/settings.yml

# Deploy
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Check status
kubectl get pods
kubectl get services
kubectl logs -f deployment/bytedancedemo
```

## Monitoring and Maintenance

### Log Management

```bash
# View application logs
sudo journalctl -u bytedancedemo -f

# View Nginx logs
sudo tail -f /var/log/nginx/bytedancedemo_access.log
sudo tail -f /var/log/nginx/bytedancedemo_error.log

# View MySQL slow queries
sudo tail -f /var/log/mysql/slow-query.log

# Log rotation is handled by Lumberjack (configured in settings.yml)
```

### Backup Strategy

#### Database Backup
```bash
# Create backup script /opt/scripts/backup.sh
#!/bin/bash
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
MYSQL_USER="sample_douyin"
MYSQL_PASSWORD="sample_douyin"
MYSQL_DATABASE="sample_douyin"

mkdir -p $BACKUP_DIR
mysqldump -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Keep only last 7 days of backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
```

Add to crontab:
```bash
# Add daily backup at 2 AM
0 2 * * * /opt/scripts/backup.sh
```

#### Redis Backup
```bash
# Redis saves snapshots automatically (RDB)
# Also enable AOF persistence in redis.conf
appendonly yes
appendfilename "appendonly.aof"
```

### Performance Monitoring

```bash
# CPU and Memory
top
htop

# Disk Usage
df -h
du -sh /var/log/bytedancedemo

# Network Connections
netstat -tunlp
ss -tunlp

# Database Connections
mysql -u sample_douyin -p -e "SHOW PROCESSLIST;"

# Redis Statistics
redis-cli INFO
redis-cli INFO stats
```

### Health Checks

```bash
# Application health check
curl http://localhost:8080/health

# Database connection check
mysqladmin ping -u sample_douyin -p

# Redis connection check
redis-cli ping

# RabbitMQ connection check
rabbitmq-diagnostics ping
```

## Rollback Procedures

### Application Rollback

```bash
# Stop current version
sudo systemctl stop bytedancedemo

# Backup current version
sudo cp /opt/bytedancedemo/bin/simple-demo /opt/bytedancedemo/bin/simple-demo.backup

# Restore previous version
sudo cp /opt/bytedancedemo/backups/simple-demo-v1.0.0 /opt/bytedancedemo/bin/simple-demo

# Start application
sudo systemctl start bytedancedemo

# Verify
sudo systemctl status bytedancedemo
curl http://localhost:8080/health
```

### Database Rollback

```bash
# Stop application
sudo systemctl stop bytedancedemo

# Restore from backup
gunzip < /backup/mysql/backup_20230415_020000.sql.gz | mysql -u sample_douyin -p sample_douyin

# Start application
sudo systemctl start bytedancedemo
```

### Docker Rollback

```bash
# Stop current containers
docker-compose down

# Restore previous images
docker tag bytedancedemo:v1.0.1 bytedancedemo:latest

# Start containers
docker-compose up -d

# Verify
docker-compose ps
docker-compose logs
```

## Troubleshooting

### Common Issues

#### Application won't start
```bash
# Check logs
sudo journalctl -u bytedancedemo -n 50

# Check configuration
./bin/simple-demo server -c config/settings.yml -m debug

# Verify dependencies
systemctl status mysql
systemctl status redis-server
systemctl status rabbitmq-server
```

#### Database connection errors
```bash
# Test MySQL connection
mysql -u sample_douyin -p -h localhost sample_douyin

# Check MySQL logs
sudo tail -f /var/log/mysql/error.log

# Verify credentials
grep -A 10 "mysql:" config/settings.yml
```

#### High memory usage
```bash
# Check memory usage
free -h

# Monitor Go garbage collection
curl http://localhost:8080/debug/pprof/heap

# Adjust MySQL buffer pool
# Edit /etc/mysql/mysql.conf.d/mysqld.cnf
innodb_buffer_pool_size = 1G  # Reduce if needed
```

#### Slow API responses
```bash
# Enable slow query logging
# Edit MySQL config
slow_query_log = 1
long_query_time = 1

# Check slow queries
sudo tail -f /var/log/mysql/slow-query.log

# Check Redis hit rate
redis-cli INFO | grep hits
```

## Support and Resources

- **Documentation**: Check `/docs` directory for detailed guides
- **Logs**: `/var/log/bytedancedemo/` for application logs
- **Monitoring**: Set up Prometheus/Grafana for production monitoring
- **Alerts**: Configure alerts for critical failures
- **Backup**: Ensure regular backups are scheduled and tested
