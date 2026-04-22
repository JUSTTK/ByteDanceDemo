# ByteDanceDemo - Configuration Reference

## Table of Contents

1. [Configuration Overview](#configuration-overview)
2. [Application Settings](#application-settings)
3. [Database Configuration](#database-configuration)
4. [Redis Configuration](#redis-configuration)
5. [Message Queue Configuration](#message-queue-configuration)
6. [Authentication Configuration](#authentication-configuration)
7. [Cloud Storage Configuration](#cloud-storage-configuration)
8. [Logging Configuration](#logging-configuration)
9. [Security Configuration](#security-configuration)
10. [Environment Variables](#environment-variables)
11. [Configuration Examples](#configuration-examples)

## Configuration Overview

ByteDanceDemo uses a YAML-based configuration system. The main configuration file is located at `config/settings.yml`. The configuration system supports:

- **Environment variable substitution** for sensitive values
- **Multiple configuration files** with environment-specific overrides
- **Hot-reload** of configuration in development mode
- **Validation** of required parameters

### Configuration Hierarchy

1. Default values (built into application)
2. Base configuration (`config/settings.yml`)
3. Environment-specific overrides (`config/settings.{env}.yml`)
4. Command-line parameters
5. Environment variables

## Application Settings

### Basic Application Configuration

```yaml
settings:
  application:
    rateLimit: 50             # Rate limit requests per minute
    port: 8080                # Server port
    host: 0.0.0.0           # Host to bind to
    timeout: 30              # Request timeout in seconds
    maxUploadSize: 50MB      # Maximum file upload size
    workerCount: 4           # Number of worker threads
    corsOrigins:            # CORS allowed origins
      - "http://localhost:3000"
      - "http://localhost:8080"
```

### Request Handling

```yaml
settings:
  application:
    enableGzip: true         # Enable gzip compression
    enablePProf: false       # Enable profiling endpoints
    readTimeout: 30s         # HTTP read timeout
    writeTimeout: 30s        # HTTP write timeout
    idleTimeout: 60s         # HTTP idle timeout
    maxHeaderBytes: 1048576 # Maximum header size (1MB)
    keepAliveTimeout: 30s    # Connection keep-alive timeout
```

### File Upload Settings

```yaml
settings:
  upload:
    path: "./public"         # Upload directory
    allowedTypes:           # Allowed file types
      - "mp4"
      - "mov"
      - "avi"
    maxFileSize: 104857600  # Maximum file size (100MB)
    tempDir: "/tmp"         # Temporary upload directory
    cleanupInterval: 1h     # Cleanup interval for temp files
```

## Database Configuration

### MySQL Configuration

```yaml
settings:
  mysql:
    host: "localhost"        # Database host
    port: 3306             # Database port
    schema: "sample_douyin" # Database name
    username: "sample_douyin" # Database username
    password: "password"     # Database password (use env var in production)
    charset: "utf8mb4"     # Character set
    parseTime: true        # Parse time values to Go types
    loc: "Local"           # Location for time parsing
    maxOpenConns: 100      # Maximum open connections
    maxIdleConns: 10       # Maximum idle connections
    connMaxLifetime: 3600s  # Connection maximum lifetime
    connMaxIdleTime: 600s  # Connection maximum idle time
    logLevel: 1            # Log level (0=off, 1=error, 2=info, 3=debug)
    slowThreshold: 100ms   # Query execution threshold for slow queries
```

### Connection Pool Settings

```yaml
settings:
  mysql:
    pool:
      retryCount: 3         # Connection retry count
      retryDelay: 5s       # Connection retry delay
      healthCheckInterval: 1m # Connection health check interval
      maxLifetime: 30m     # Maximum connection lifetime
      maxIdleTime: 5m      # Maximum idle time
      cleanWait: 1m       # Wait time for connection cleanup
```

### Database Migration Settings

```yaml
settings:
  database:
    migration:
      autoMigrate: true     # Enable automatic migration
      seed: true           # Run seed data after migration
      path: "./migrations" # Migration files directory
      tableName: "schema_migrations" # Migration tracking table
```

## Redis Configuration

### Basic Redis Settings

```yaml
settings:
  redis:
    addr: "localhost:6379"  # Redis address
    password: ""           # Redis password (use env var in production)
    db: 0                 # Redis database number
    poolSize: 10          # Connection pool size
    minIdleConns: 5       # Minimum idle connections
    maxRetries: 3         # Maximum connection retries
    dialTimeout: 5s      # Connection timeout
    readTimeout: 3s       # Read timeout
    writeTimeout: 3s      # Write timeout
    poolTimeout: 4s       # Pool timeout
```

### Redis Cluster Configuration (for production)

```yaml
settings:
  redis:
    cluster:
      enabled: true
      addrs:
        - "redis-node1:6379"
        - "redis-node2:6379"
        - "redis-node3:6379"
      readOnly: false
      maxRedirects: 16
      routeByLatency: true
      routeRandomly: true
```

### Redis Caching Configuration

```yaml
settings:
  redis:
    cache:
      defaultTTL: 5m       # Default cache TTL
      cleanupInterval: 5m   # Cache cleanup interval
      maxKeys: 10000       # Maximum number of cached keys
      maxMemory: "512MB"    # Maximum memory usage
      memoryPolicy: "allkeys-lru" # Eviction policy
    session:
      prefix: "session:"   # Session key prefix
      ttl: 24h             # Session TTL
      enable: true         # Enable session store
```

## Message Queue Configuration

### RabbitMQ Configuration

```yaml
settings:
  rabbitMQ:
    host: "localhost"      # RabbitMQ host
    port: 5672            # RabbitMQ port
    username: "guest"      # RabbitMQ username (use env var in production)
    password: "guest"      # RabbitMQ password (use env var in production)
    vhost: "/"            # RabbitMQ virtual host
    reconnectDelay: 5s    # Reconnection delay
    reconnectAttempts: 5 # Reconnection attempts
    exchange: "bytedancedemo" # Default exchange
    routingKey: "default"  # Default routing key
    queue: "bytedancedemo" # Default queue
    durable: true         # Queue durability
```

### Queue Configuration

```yaml
settings:
  rabbitMQ:
    queues:
      message:
        name: "messages"
        durable: true
        exclusive: false
        autoDelete: false
      comment:
        name: "comments"
        durable: true
        exclusive: false
        autoDelete: false
      notification:
        name: "notifications"
        durable: true
        exclusive: false
        autoDelete: false
```

### Exchange Configuration

```yaml
settings:
  rabbitMQ:
    exchanges:
      default:
        name: "bytedancedemo"
        type: "direct"
        durable: true
        autoDelete: false
      fanout:
        name: "events"
        type: "fanout"
        durable: true
        autoDelete: false
      topic:
        name: "topics"
        type: "topic"
        durable: true
        autoDelete: false
```

## Authentication Configuration

### JWT Configuration

```yaml
settings:
  jwt:
    secretKey: "your-secret-key" # JWT secret key (use env var in production)
    expirationTime: 24h          # Token expiration time
    refreshExpirationTime: 168h # Refresh token expiration (1 week)
    issuer: "bytedancedemo"     # Token issuer
    audience: "users"           # Token audience
    algorithm: "HS256"          # Signing algorithm
    leeway: 60s                 # Time leeway for clock skew
```

### Token Configuration

```yaml
settings:
  token:
    header: "Authorization"      # Authorization header name
    prefix: "Bearer"           # Token prefix
    key: "token"               # Token parameter name
    cookieName: "auth_token"   # Cookie name (for web clients)
    secureCookie: false       # Secure cookie flag
    sameSite: "lax"           # SameSite cookie policy
```

### Authorization Configuration

```yaml
settings:
  authorization:
    enabled: true              # Enable authorization
    model: "rbac"             # Authorization model (rbac, abac)
    log: true                  # Log authorization decisions
    enabledDomains:            # Enabled authorization domains
      - "api"
      - "resource"
```

### RBAC Configuration

```yaml
settings:
  rbac:
    modelPath: "./config/rbac_model.conf"  # Casbin model file
    policyPath: "./config/policy.csv"     # Casbin policy file
    autoReload: true                     # Auto reload policy changes
    enforcePriority: true                # Enforce policy priority
```

## Cloud Storage Configuration

### OSS/Cloud Storage

```yaml
settings:
  oss:
    provider: "local"          # Storage provider (local, s3, oss, minio)
    endpoint: ""              # Storage endpoint (for cloud providers)
    accessKey: ""             # Access key (for cloud providers)
    secretKey: ""             # Secret key (for cloud providers)
    bucket: ""                # Storage bucket name
    region: ""                # Storage region
    domain: ""                # Custom domain (CDN)
    acl: "public-read"        # Default ACL
    enableSSL: true           # Use HTTPS
    cacheControl: "public,max-age=86400" # Cache control header
```

### CDN Configuration

```yaml
settings:
  cdn:
    enabled: false            # Enable CDN
    domain: "cdn.example.com" # CDN domain
    ssl: true                # Use HTTPS for CDN
    cacheSettings:           # Cache settings
      video:
        ttl: 31536000         # 1 year
        query: "public,max-age=31536000"
      image:
        ttl: 2592000         # 30 days
        query: "public,max-age=2592000"
      static:
        ttl: 604800          # 7 days
        query: "public,max-age=604800"
```

## Logging Configuration

### Basic Logging Settings

```yaml
settings:
  log:
    path: "./logs"           # Log directory
    level: "info"            # Log level (debug, info, warn, error, fatal, panic)
    format: "json"          # Log format (json, text)
    output: "file"          # Output (file, stdout, both)
    timestamps: true        # Include timestamps
    caller: false           # Include caller information
    development: false      # Development mode (more verbose)
```

### Log Rotation Settings

```yaml
settings:
  log:
    rotation:
      maxSize: 100           # Maximum log size (MB)
      maxAge: 30            # Maximum log age (days)
      maxBackups: 10        # Maximum number of backup files
      compress: true        # Compress old log files
      localTime: true       # Use local time for timestamps
      removeCompress: true  # Remove compressed old logs
```

### Structured Logging

```yaml
settings:
  log:
    structured:
      enable: true          # Enable structured logging
      fields:               # Default fields
        service: "bytedancedemo"
        version: "1.0.0"
        environment: "development"
      context: true        # Include request context
      stacktraceLevel: "warn" # Level for stack traces
```

### Access Logging

```yaml
settings:
  accessLog:
    enabled: true           # Enable access logging
    path: "./logs/access.log"
    format: "combined"      # Log format (combined, common, json)
    ignorePaths:           # Paths to ignore
      - "/health"
      - "/metrics"
    maxRequests: 10000     # Rotate after N requests
    flushInterval: 1s      # Flush interval
```

## Security Configuration

### Security Headers

```yaml
settings:
  security:
    headers:
      xFrameOptions: "SAMEORIGIN"
      xContentTypeOptions: "nosniff"
      xSSProtection: "1; mode=block"
      xPoweredBy: "false"
      referrerPolicy: "strict-origin-when-cross-origin"
      contentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src * data:"
      strictTransportSecurity:
        enable: true
        maxAge: 31536000
        includeSubDomains: true
        preload: true
```

### Rate Limiting

```yaml
settings:
  ratelimit:
    enabled: true
    requests: 100           # Requests per minute
    burst: 50               # Burst requests
    clients:               # Client-specific limits
      default:
        requests: 100
        burst: 50
      "127.0.0.1":
        requests: 1000
        burst: 500
    store: "redis"         # Storage for rate limits (memory, redis)
    cleanupInterval: 1h    # Cleanup interval
```

### CSRF Protection

```yaml
settings:
  csrf:
    enabled: true          # Enable CSRF protection
    header: "X-CSRF-Token"  # CSRF header name
    cookie: "csrf_token"   # CSRF cookie name
    cookieDomain: ""        # Cookie domain
    cookieSecure: false     # Secure cookie flag
    cookieHTTPOnly: true    # HTTP-only cookie flag
    cookiePath: "/"        # Cookie path
    maxAge: 3600           # Cookie lifetime (seconds)
```

### Input Validation

```yaml
settings:
  validation:
    enabled: true          # Enable input validation
    maxStringLength: 65536 # Maximum string length
    maxArrayLength: 1000   # Maximum array length
    maxFileSize: 104857600 # Maximum file size (bytes)
    sanitizeInput: true     # Sanitize input data
    allowedOrigins:        # Allowed origins
      - "http://localhost:3000"
    allowedMethods:        # Allowed methods
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
    allowedHeaders:        # Allowed headers
      - "Content-Type"
      - "Authorization"
```

## Environment Variables

### Core Environment Variables

```bash
# Application Configuration
APP_ENV=production         # Environment (development, staging, production)
APP_DEBUG=false            # Debug mode
APP_PORT=8080              # Server port

# Database Configuration
MYSQL_HOST=localhost       # MySQL host
MYSQL_PORT=3306           # MySQL port
MYSQL_DATABASE=sample_douyin # Database name
MYSQL_USERNAME=sample_douyin # Username
MYSQL_PASSWORD=password    # Password

# Redis Configuration
REDIS_HOST=localhost      # Redis host
REDIS_PORT=6379           # Redis port
REDIS_PASSWORD=           # Redis password
REDIS_DB=0                # Redis database

# RabbitMQ Configuration
RABBITMQ_HOST=localhost   # RabbitMQ host
RABBITMQ_PORT=5672        # RabbitMQ port
RABBITMQ_USERNAME=guest   # Username
RABBITMQ_PASSWORD=guest  # Password

# JWT Configuration
JWT_SECRET=your-secret-key # JWT secret
JWT_EXPIRE=24h            # Token expiration

# Cloud Storage Configuration
OSS_ACCESS_KEY=           # OSS access key
OSS_SECRET_KEY=           # OSS secret key
OSS_ENDPOINT=             # OSS endpoint
OSS_BUCKET=               # OSS bucket
OSS_REGION=               # OSS region

# Logging Configuration
LOG_LEVEL=info           # Log level
LOG_PATH=./logs          # Log path
LOG_FORMAT=json          # Log format
```

### Optional Environment Variables

```bash
# Application Performance
WORKER_COUNT=4            # Number of worker threads
MAX_UPLOAD_SIZE=52428800 # Maximum upload size (50MB)

# Security
CORS_ENABLED=true         # Enable CORS
CORS_ORIGINS="http://localhost:3000" # CORS origins

# Monitoring
ENABLE_METRICS=true       # Enable metrics endpoint
HEALTH_CHECK_ENABLED=true # Enable health check
PPROF_ENABLED=false       # Enable profiling

# Development
SWAGGER_ENABLED=true      # Enable Swagger UI
ADMIN_PANEL_ENABLED=true  # Enable admin panel
```

## Configuration Examples

### Development Configuration

`config/settings.yml`:
```yaml
settings:
  application:
    rateLimit: 50
    debug: true
    enablePProf: true
    corsOrigins:
      - "http://localhost:3000"
      - "http://localhost:8080"
  mysql:
    host: "localhost"
    port: 3306
    schema: "sample_douyin"
    username: "sample_douyin"
    password: "dev_password"
    logLevel: 3  # Debug
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
  rabbitMQ:
    host: "localhost"
    port: 5672
    username: guest
    password: guest
  log:
    level: "debug"
    format: "text"
    output: "stdout"
    development: true
  security:
    headers:
      xFrameOptions: "ALLOWALL"
```

### Production Configuration

`config/settings.production.yml`:
```yaml
settings:
  application:
    rateLimit: 100
    debug: false
    enablePProf: false
    corsOrigins:
      - "https://yourdomain.com"
      - "https://api.yourdomain.com"
    workerCount: 8
  mysql:
    host: "${MYSQL_HOST}"
    port: 3306
    schema: "${MYSQL_DATABASE}"
    username: "${MYSQL_USERNAME}"
    password: "${MYSQL_PASSWORD}"
    logLevel: 1  # Error only
    maxOpenConns: 100
    maxIdleConns: 20
    maxLifetime: 5m
  redis:
    addr: "${REDIS_HOST}:${REDIS_PORT}"
    password: "${REDIS_PASSWORD}"
    db: ${REDIS_DB}
    poolSize: 20
    maxIdleConns: 10
  rabbitMQ:
    host: "${RABBITMQ_HOST}"
    port: 5672
    username: "${RABBITMQ_USERNAME}"
    password: "${RABBITMQ_PASSWORD}"
  jwt:
    secretKey: "${JWT_SECRET}"
    expirationTime: 24h
  log:
    level: "info"
    format: "json"
    output: "file"
    path: "${LOG_PATH:-/var/log/bytedancedemo}"
    maxSize: 100
    maxAge: 30
    maxBackups: 10
    compress: true
  security:
    headers:
      xFrameOptions: "SAMEORIGIN"
      xContentTypeOptions: "nosniff"
      xSSProtection: "1; mode=block"
    ratelimit:
      enabled: true
      requests: 100
      burst: 50
      store: "redis"
```

### Docker Configuration

`config/settings.docker.yml`:
```yaml
settings:
  application:
    port: 8080
    host: "0.0.0.0"
    enableGzip: true
  mysql:
    host: "mysql"
    port: 3306
    schema: "sample_douyin"
    username: "sample_douyin"
    password: "sample_douyin"
  redis:
    addr: "redis:6379"
    password: ""
    db: 0
  rabbitMQ:
    host: "rabbitmq"
    port: 5672
    username: guest
    password: guest
  log:
    level: "info"
    format: "text"
    output: "stdout"
```

### Kubernetes Configuration

Use ConfigMap and Secrets:

```bash
# Create ConfigMap
kubectl create configmap bytedancedemo-config \
  --from-file=config/settings.yml \
  --from-literal=configOverride='{"application": {"debug": false, "enableMetrics": true}}'

# Create Secret
kubectl create secret generic bytedancedemo-secrets \
  --from-literal=mysql-password=$(openssl rand -base64 32) \
  --from-literal=redis-password=$(openssl rand -base64 32) \
  --from-literal=rabbitmq-password=$(openssl rand -base64 32) \
  --from-literal=jwt-secret=$(openssl rand -base64 32)
```

## Configuration Validation

### Required Configuration Parameters

The following configuration parameters are required for proper operation:

```yaml
settings:
  mysql:
    host: required
    port: required
    schema: required
    username: required
    password: required
  redis:
    addr: required
  rabbitMQ:
    host: required
    port: required
    username: required
    password: required
  jwt:
    secretKey: required
```

### Configuration Validation Commands

```bash
# Test configuration file validity
go run cmd/api/service.go -c config/settings.yml -m debug

# Check required parameters
grep -E 'required|${' config/settings.yml

# Validate environment variables
echo $MYSQL_PASSWORD
echo $REDIS_PASSWORD
echo $JWT_SECRET
```

## Troubleshooting Configuration Issues

### Common Issues

1. **Database Connection Failed**
   - Verify MySQL credentials in settings.yml
   - Check MySQL server is running
   - Verify network connectivity
   ```bash
   mysql -h localhost -P 3306 -u sample_douyin -p
   ```

2. **Redis Connection Failed**
   - Verify Redis server is running
   - Check Redis configuration settings
   ```bash
   redis-cli ping
   redis-cli -h localhost -p 6379 ping
   ```

3. **JWT Token Issues**
   - Verify JWT secret is set
   - Check token expiration time
   ```bash
   echo $JWT_SECRET | wc -c  # Should be at least 32 characters
   ```

4. **File Upload Issues**
   - Verify upload directory exists and is writable
   - Check file size limits
   ```bash
   mkdir -p public
   chmod 755 public
   ```

### Configuration Hot Reload

In development mode, you can hot-reload configuration:

```bash
# Start with hot reload
./bin/simple-demo server -c config/settings.yml -m debug --hot-reload

# Wait for file changes
# Configuration will be automatically reloaded
```

### Environment Check Script

Create a script to verify all configurations:

`scripts/check-env.sh`:
```bash
#!/bin/bash

# Check required environment variables
required_vars=(
    "MYSQL_HOST"
    "MYSQL_PORT"
    "MYSQL_DATABASE"
    "MYSQL_USERNAME"
    "REDIS_HOST"
    "REDIS_PORT"
    "RABBITMQ_HOST"
    "RABBITMQ_PORT"
    "JWT_SECRET"
)

# Check each variable
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "ERROR: Environment variable $var is not set"
        exit 1
    else
        echo "✓ $var is set"
    fi
done

# Check services are running
echo "Checking services..."
mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USERNAME -p -e "SELECT 1;" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ MySQL is accessible"
else
    echo "✗ MySQL is not accessible"
    exit 1
fi

redis-cli -h $REDIS_HOST -p $REDIS_PORT ping > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ Redis is accessible"
else
    echo "✗ Redis is not accessible"
    exit 1
fi

echo "All checks passed!"
```

## Best Practices

1. **Never commit sensitive credentials** to version control
2. **Use environment variables** for sensitive data in production
3. **Rotate secrets** regularly
4. **Enable logging** in development, minimize in production
5. **Validate configuration** after changes
6. **Backup configuration files** before major changes
7. **Use separate configurations** for different environments
8. **Monitor configuration changes** with version control
9. **Document custom configurations** for team members
10. **Test configuration changes** in staging before production
