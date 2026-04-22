# ByteDanceDemo - Security Best Practices

## Table of Contents

1. [Security Overview](#security-overview)
2. [Authentication Security](#authentication-security)
3. [Authorization Security](#authorization-security)
4. [Data Protection](#data-protection)
5. [API Security](#api-security)
6. [Infrastructure Security](#infrastructure-security)
7. [Common Vulnerabilities](#common-vulnerabilities)
8. [Security Testing](#security-testing)
9. [Incident Response](#incident-response)
10. [Security Checklist](#security-checklist)

## Security Overview

### Security Principles

ByteDanceDemo follows these core security principles:

- **Defense in Depth**: Multiple layers of security controls
- **Least Privilege**: Minimum necessary access rights
- **Secure by Default**: Security features enabled automatically
- **Fail Securely**: Fail to secure state on errors
- **Transparency**: Open security practices and disclosures

### Security Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                         │
├─────────────────────────────────────────────────────────────┤
│  1. Network Security    (Firewall, DDoS protection)      │
│  2. Web Server         (Nginx security headers)          │
│  3. Application        (Input validation, sanitization)   │
│  4. Authentication     (JWT, bcrypt, rate limiting)       │
│  5. Authorization      (Casbin RBAC, role-based access)   │
│  6. Data Security      (Encryption, secure storage)       │
│  7. Logging           (Audit trails, monitoring)          │
└─────────────────────────────────────────────────────────────┘
```

### Threat Model

We consider the following threat vectors:

- **Authentication attacks**: Brute force, credential stuffing
- **Injection attacks**: SQL, NoSQL, command injection
- **XSS attacks**: Cross-site scripting
- **CSRF attacks**: Cross-site request forgery
- **Data breaches**: Unauthorized data access
- **DoS attacks**: Denial of service
- **MITM attacks**: Man-in-the-middle attacks

## Authentication Security

### Password Security

#### Password Hashing

We use bcrypt for password hashing:

```go
import "golang.org/x/crypto/bcrypt"

// HashPassword securely hashes a password using bcrypt
func HashPassword(password string) (string, error) {
    // Use bcrypt with cost factor 12 (recommended)
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(bytes), nil
}

// CheckPassword verifies a password against a hash
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

#### Password Requirements

```go
import "regexp"

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
    // Minimum 8 characters
    if len(password) < 8 {
        return false
    }
    
    // Contains uppercase letter
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    // Contains lowercase letter
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    // Contains digit
    hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
    // Contains special character
    hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
    
    return hasUpper && hasLower && hasDigit && hasSpecial
}

// Common password dictionary check
func IsCommonPassword(password string) bool {
    commonPasswords := []string{
        "password", "123456", "qwerty", "admin",
        // Add more common passwords
    }
    for _, common := range commonPasswords {
        if strings.EqualFold(password, common) {
            return true
        }
    }
    return false
}
```

### JWT Token Security

#### Token Generation

```go
import (
    "github.com/golang-jwt/jwt"
    "time"
)

// GenerateToken creates a secure JWT token
func GenerateToken(userID int64, secretKey string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
        "iat":     time.Now().Unix(),
        "nbf":     time.Now().Add(-1 * time.Minute).Unix(),
        "jti":     generateUniqueTokenID(), // JWT ID for revocation
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// generateUniqueTokenID generates a unique JWT ID
func generateUniqueTokenID() string {
    return uuid.New().String()
}
```

#### Token Validation

```go
// ValidateToken validates a JWT token
func ValidateToken(tokenString, secretKey string) (*jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return &claims, nil
    }
    
    return nil, fmt.Errorf("invalid token")
}

// IsTokenRevoked checks if token is in revocation list
func IsTokenRevoked(tokenID string) bool {
    // Check Redis or database for revoked tokens
    return false
}
```

#### Token Refresh Mechanism

```go
// RefreshTokenPair generates new access and refresh tokens
func RefreshTokenPair(refreshToken, secretKey string) (string, string, error) {
    // Validate refresh token
    claims, err := ValidateToken(refreshToken, secretKey)
    if err != nil {
        return "", "", err
    }
    
    // Check if refresh token is revoked
    if tokenID, ok := (*claims)["jti"].(string); ok {
        if IsTokenRevoked(tokenID) {
            return "", "", fmt.Errorf("token is revoked")
        }
    }
    
    userID := int64((*claims)["user_id"].(float64))
    
    // Generate new tokens
    accessToken, err := GenerateToken(userID, secretKey)
    if err != nil {
        return "", "", err
    }
    
    newRefreshToken, err := GenerateRefreshToken(userID, secretKey)
    if err != nil {
        return "", "", err
    }
    
    // Revoke old refresh token
    if tokenID, ok := (*claims)["jti"].(string); ok {
        RevokeToken(tokenID)
    }
    
    return accessToken, newRefreshToken, nil
}
```

### Rate Limiting

#### IP-based Rate Limiting

```go
import "github.com/redis/go-redis/v9"

// RateLimiter manages rate limiting
type RateLimiter struct {
    redis  *redis.Client
    config RateLimitConfig
}

type RateLimitConfig struct {
    RequestsPerMinute int
    BurstSize        int
}

// CheckRateLimit checks if request should be rate limited
func (rl *RateLimiter) CheckRateLimit(ip string) (bool, error) {
    key := fmt.Sprintf("ratelimit:%s", ip)
    now := time.Now().Unix()
    
    // Use Redis atomic operations for rate limiting
    result, err := rl.redis.Eval(context.Background(), `
        local key = KEYS[1]
        local now = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])
        local limit = tonumber(ARGV[3])
        
        redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)
        local count = redis.call('ZCARD', key)
        
        if count < limit then
            redis.call('ZADD', key, now, now)
            redis.call('EXPIRE', key, window)
            return 1
        else
            return 0
        end
    `, []string{key}, now, 60, rl.config.RequestsPerMinute).Int()
    
    if err != nil {
        return false, err
    }
    
    return result == 1, nil
}
```

#### User-based Rate Limiting

```go
// UserRateLimiter manages user-specific rate limits
type UserRateLimiter struct {
    redis  *redis.Client
    config UserRateLimitConfig
}

type UserRateLimitConfig struct {
    DefaultLimit int
    UserLimits   map[string]int // user_id -> custom limit
}

// CheckUserRateLimit checks rate limit for specific user
func (rl *UserRateLimiter) CheckUserRateLimit(userID string) (bool, error) {
    limit := rl.config.DefaultLimit
    if customLimit, exists := rl.config.UserLimits[userID]; exists {
        limit = customLimit
    }
    
    key := fmt.Sprintf("ratelimit:user:%s", userID)
    now := time.Now().Unix()
    
    // Similar to IP-based rate limiting
    result, err := rl.redis.Eval(context.Background(), `
        -- Same logic as IP-based but with user-specific limits
    `, []string{key}, now, 60, limit).Int()
    
    return result == 1, err
}
```

## Authorization Security

### Casbin RBAC Configuration

#### Model Configuration

`config/rbac_model.conf`:
```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

#### Policy Enforcement

```go
import "github.com/casbin/casbin/v2"

// AuthService handles authorization
type AuthService struct {
    enforcer *casbin.Enforcer
}

// CheckPermission checks if user has permission
func (as *AuthService) CheckPermission(userID string, resource, action string) (bool, error) {
    return as.enforcer.Enforce(userID, resource, action)
}

// AddRoleForUser assigns role to user
func (as *AuthService) AddRoleForUser(userID, role string) error {
    return as.enforcer.AddRoleForUser(userID, role)
}

// DefinePermission defines permission
func (as *AuthService) DefinePermission(role, resource, action string) error {
    return as.enforcer.AddPolicy(role, resource, action)
}
```

### Resource-level Authorization

```go
// ResourceOwnershipCheck verifies user owns the resource
func ResourceOwnershipCheck(userID, resourceType, resourceID string) (bool, error) {
    switch resourceType {
    case "video":
        return checkVideoOwnership(userID, resourceID)
    case "comment":
        return checkCommentOwnership(userID, resourceID)
    default:
        return false, fmt.Errorf("unknown resource type")
    }
}

func checkVideoOwnership(userID, videoID string) (bool, error) {
    var video model.Video
    result := db.First(&video, videoID)
    if result.Error != nil {
        return false, result.Error
    }
    return video.AuthorID == userID, nil
}
```

## Data Protection

### Encryption

#### Sensitive Data Encryption

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

// EncryptData encrypts sensitive data using AES-256-GCM
func EncryptData(plaintext, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts encrypted data
func DecryptData(ciphertext string, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    data, err := base64.URLEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
```

#### Database Encryption

```go
// EncryptedField represents an encrypted field in database
type EncryptedField string

// BeforeSave encrypts field before saving to database
func (ef *EncryptedField) BeforeSave(tx *gorm.DB) (err error) {
    if len(*ef) == 0 {
        return nil
    }
    
    encrypted, err := EncryptData([]byte(*ef), getEncryptionKey())
    if err != nil {
        return err
    }
    
    *ef = EncryptedField(encrypted)
    return nil
}

// AfterFind decrypts field after retrieving from database
func (ef *EncryptedField) AfterFind(tx *gorm.DB) (err error) {
    if len(*ef) == 0 {
        return nil
    }
    
    decrypted, err := DecryptData(string(*ef), getEncryptionKey())
    if err != nil {
        return err
    }
    
    *ef = EncryptedField(decrypted)
    return nil
}
```

### Input Validation and Sanitization

#### Input Sanitization

```go
import (
    "html"
    "regexp"
    "strings"
)

// SanitizeInput sanitizes user input
func SanitizeInput(input string) string {
    // Remove HTML tags
    sanitized := html.EscapeString(input)
    
    // Remove potential SQL injection patterns
    sanitized = regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|alter)`).ReplaceAllString(sanitized, "")
    
    // Remove control characters
    sanitized = strings.Map(func(r rune) rune {
        if r < 32 || r > 126 {
            return -1
        }
        return r
    }, sanitized)
    
    // Trim whitespace
    sanitized = strings.TrimSpace(sanitized)
    
    return sanitized
}

// ValidateUsername validates username format
func ValidateUsername(username string) error {
    if len(username) < 3 || len(username) > 20 {
        return fmt.Errorf("username must be between 3 and 20 characters")
    }
    
    // Only alphanumeric and underscore
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    if !matched {
        return fmt.Errorf("username can only contain letters, numbers, and underscores")
    }
    
    return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
    if !matched {
        return fmt.Errorf("invalid email format")
    }
    return nil
}
```

#### File Upload Security

```go
import (
    "mime/multipart"
    "path/filepath"
    "os"
    "strings"
)

// UploadConfig configures secure file upload
type UploadConfig struct {
    AllowedTypes      []string
    MaxFileSize       int64
    UploadPath       string
    AllowedExtensions []string
}

// SecureUploadHandler handles secure file uploads
func SecureUploadHandler(fileHeader *multipart.FileHeader, config UploadConfig) (string, error) {
    // Validate file size
    if fileHeader.Size > config.MaxFileSize {
        return "", fmt.Errorf("file size exceeds maximum allowed size")
    }
    
    // Validate file extension
    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    if !isAllowedExtension(ext, config.AllowedExtensions) {
        return "", fmt.Errorf("file type not allowed")
    }
    
    // Validate MIME type
    file, err := fileHeader.Open()
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    // Detect actual MIME type
    mimeType := detectMimeType(file)
    if !isAllowedMimeType(mimeType, config.AllowedTypes) {
        return "", fmt.Errorf("invalid MIME type")
    }
    
    // Generate secure filename
    filename := generateSecureFilename(fileHeader.Filename)
    
    // Save file
    filepath := filepath.Join(config.UploadPath, filename)
    if err := saveFile(file, filepath); err != nil {
        return "", err
    }
    
    return filename, nil
}

func generateSecureFilename(originalName string) string {
    ext := filepath.Ext(originalName)
    // Use UUID instead of original filename
    return uuid.New().String() + ext
}

func isAllowedExtension(ext string, allowed []string) bool {
    for _, allowedExt := range allowed {
        if strings.EqualFold(ext, allowedExt) {
            return true
        }
    }
    return false
}
```

## API Security

### Security Headers

```go
// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Content Security Policy
        c.Header("Content-Security-Policy",
            "default-src 'self'; " +
            "script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
            "style-src 'self' 'unsafe-inline'; " +
            "img-src 'self' data: https:; " +
            "font-src 'self'; " +
            "connect-src 'self'; " +
            "frame-ancestors 'none'")
        
        // X-Frame-Options
        c.Header("X-Frame-Options", "DENY")
        
        // X-Content-Type-Options
        c.Header("X-Content-Type-Options", "nosniff")
        
        // X-XSS-Protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Referrer-Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Strict-Transport-Security (only on HTTPS)
        if c.Request.TLS != nil {
            c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        }
        
        // Permissions-Policy
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        c.Next()
    }
}
```

### CORS Configuration

```go
import "github.com/gin-contrib/cors"

// CORSMiddleware configures CORS
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"https://yourdomain.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
        ExposeHeaders:    []string{"Content-Length", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

### CSRF Protection

```go
import (
    "crypto/rand"
    "encoding/hex"
)

// CSRFMiddleware provides CSRF protection
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip CSRF check for safe methods
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }
        
        // Get CSRF token from header
        token := c.GetHeader("X-CSRF-Token")
        
        // Get CSRF token from cookie
        cookie, err := c.Cookie("csrf_token")
        if err != nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
            c.Abort()
            return
        }
        
        // Validate CSRF token
        if token != cookie {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
```

## Infrastructure Security

### Network Security

#### Firewall Configuration

```bash
# UFW (Uncomplicated Firewall) configuration
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow database access from specific IPs only
sudo ufw allow from 192.168.1.100 to any port 3306

# Enable firewall
sudo ufw enable
```

#### SSL/TLS Configuration

```nginx
# Nginx SSL configuration
server {
    listen 443 ssl http2;
    server_name api.example.com;

    # SSL certificates
    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;

    # SSL protocols
    ssl_protocols TLSv1.2 TLSv1.3;

    # SSL ciphers
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers on;

    # SSL session cache
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # SSL stapling
    ssl_stapling on;
    ssl_stapling_verify on;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
}
```

### Database Security

#### MySQL Security

```sql
-- Create dedicated database user with limited privileges
CREATE USER 'bytedancedemo_app'@'localhost' IDENTIFIED BY 'secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON sample_douyin.* TO 'bytedancedemo_app'@'localhost';
GRANT EXECUTE ON PROCEDURE sample_douyin.* TO 'bytedancedemo_app'@'localhost';
FLUSH PRIVILEGES;

-- Remove test database and anonymous users
DROP DATABASE IF EXISTS test;
DELETE FROM mysql.user WHERE User='';
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');
FLUSH PRIVILEGES;

-- Enable query logging (for monitoring)
SET GLOBAL general_log = 'ON';
SET GLOBAL general_log_file = '/var/log/mysql/query.log';
```

#### Redis Security

```conf
# redis.conf security settings

# Require authentication
requirepass your_redis_password

# Disable dangerous commands
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command CONFIG ""
rename-command SHUTDOWN ""
rename-command DEBUG ""

# Bind to specific interface
bind 127.0.0.1

# Enable protected mode
protected-mode yes

# Disable access to sensitive data
rename-command KEYS "dangerous_CMD_NOT_ALLOWED"
```

### Application Security

#### Secure Configuration

```go
// SecureConfig provides secure configuration defaults
func SecureConfig() config.Settings {
    return config.Settings{
        JWT: config.JWTConfig{
            SecretKey:          os.Getenv("JWT_SECRET"),
            ExpirationTime:     24 * time.Hour,
            RefreshExpirationTime: 7 * 24 * time.Hour,
        },
        Database: config.DatabaseConfig{
            MaxOpenConns: 100,
            MaxIdleConns: 20,
            MaxLifetime:   5 * time.Minute,
        },
        Redis: config.RedisConfig{
            Password: os.Getenv("REDIS_PASSWORD"),
            DB:       0,
        },
        Security: config.SecurityConfig{
            EnableHTTPS:       true,
            EnableHSTS:       true,
            EnableCSRF:       true,
            EnableRateLimit:   true,
            MaxUploadSize:     50 * 1024 * 1024, // 50MB
            AllowedOrigins:    []string{"https://yourdomain.com"},
        },
    }
}
```

## Common Vulnerabilities

### SQL Injection Prevention

```go
// ❌ VULNERABLE - SQL Injection
func VulnerableGetUser(username string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    // This allows SQL injection
    result := db.Raw(query)
    // ...
}

// ✅ SECURE - Parameterized query
func SecureGetUser(username string) (*User, error) {
    var user User
    result := db.Where("username = ?", username).First(&user)
    return &user, result.Error
}
```

### XSS Prevention

```go
// ❌ VULNERABLE - XSS
func VulnerableRenderContent(content string) string {
    return content // Direct output allows XSS
}

// ✅ SECURE - HTML escaping
func SecureRenderContent(content string) string {
    return html.EscapeString(content)
}

// ✅ SECURE - Using safe HTML library
func SecureRenderHTML(content string) string {
    policy := bluemonday.UGCPolicy()
    return policy.Sanitize(content)
}
```

### CSRF Prevention

```go
// ✅ SECURE - CSRF protection implementation
func CSRFProtectedHandler(c *gin.Context) {
    // Verify CSRF token
    token := c.GetHeader("X-CSRF-Token")
    if !validateCSRFToken(token, c) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
        return
    }
    
    // Process request
    // ...
}

func validateCSRFToken(token string, c *gin.Context) bool {
    cookie, err := c.Cookie("csrf_token")
    if err != nil {
        return false
    }
    
    // Use timing-safe comparison
    return subtle.ConstantTimeCompare([]byte(token), []byte(cookie)) == 1
}
```

## Security Testing

### Security Testing Tools

```bash
# OWASP ZAP - Web application security scanner
docker run -t owasp/zap2docker-stable zap-cli quick-scan \
  --self-contained \
  --start-options '-config api.disablekey=true' \
  http://localhost:8080

# Gosec - Go security checker
gosec -fmt=json -out=security-report.json ./...

# Nuclei - Vulnerability scanner
nuclei -u http://localhost:8080 -t cves/ -o results.txt

# SQLMap - SQL injection tester
sqlmap -u "http://localhost:8080/douyin/user/?username=test" --batch

# Nmap - Network scanner
nmap -sV --script vuln localhost
```

### Penetration Testing Checklist

```markdown
## Authentication & Authorization
- [ ] Test weak password policies
- [ ] Test brute force protection
- [ ] Test session management
- [ ] Test privilege escalation
- [ ] Test JWT token manipulation

## Input Validation
- [ ] Test SQL injection
- [ ] Test XSS attacks
- [ ] Test CSRF attacks
- [ ] Test file upload vulnerabilities
- [ ] Test command injection

## API Security
- [ ] Test rate limiting
- [ ] Test API key management
- [ ] Test sensitive data exposure
- [ ] Test API versioning
- [ ] Test error handling

## Infrastructure
- [ ] Test firewall rules
- [ ] Test SSL/TLS configuration
- [ ] Test database security
- [ ] Test server hardening
- [ ] Test logging and monitoring
```

## Incident Response

### Security Incident Response Plan

```markdown
## Incident Response Steps

### 1. Detection and Analysis
- Identify the security incident
- Determine scope and impact
- Classify incident severity
- Preserve evidence

### 2. Containment
- Isolate affected systems
- Block malicious IPs
- Disable compromised accounts
- Temporarily shut down services if needed

### 3. Eradication
- Identify root cause
- Remove malicious content
- Patch vulnerabilities
- Secure systems

### 4. Recovery
- Restore from clean backups
- Update security measures
- Monitor for continued threats
- Document lessons learned

### 5. Post-Incident Activity
- Conduct post-mortem analysis
- Update security policies
- Implement new controls
- Train staff on lessons learned
```

### Security Monitoring

```go
// SecurityLogger logs security events
type SecurityLogger struct {
    logger *zap.Logger
}

func (sl *SecurityLogger) LogSecurityEvent(event SecurityEvent) {
    sl.logger.Warn("Security Event",
        zap.String("type", event.Type),
        zap.String("severity", event.Severity),
        zap.String("user_id", event.UserID),
        zap.String("ip", event.IP),
        zap.String("user_agent", event.UserAgent),
        zap.Time("timestamp", event.Timestamp),
        zap.Any("details", event.Details),
    )
    
    // Send to SIEM or security monitoring system
    if event.Severity == "critical" {
        sl.sendAlert(event)
    }
}

type SecurityEvent struct {
    Type       string
    Severity   string
    UserID     string
    IP         string
    UserAgent  string
    Timestamp  time.Time
    Details    map[string]interface{}
}
```

## Security Checklist

### Pre-Deployment Checklist

```markdown
## Authentication
- [ ] Passwords hashed with bcrypt (cost >= 12)
- [ ] JWT tokens properly signed and validated
- [ ] Token expiration configured
- [ ] Session management implemented
- [ ] Rate limiting enabled
- [ ] Account lockout policy configured

## Authorization
- [ ] Role-based access control implemented
- [ ] Resource ownership verification
- [ ] Admin interfaces secured
- [ ] API endpoints properly protected

## Data Protection
- [ ] Sensitive data encrypted at rest
- [ ] Data encrypted in transit
- [ ] Input validation implemented
- [ ] Output sanitization performed
- [ ] Secure file upload handling

## API Security
- [ ] CORS properly configured
- [ ] CSRF protection enabled
- [ ] Security headers implemented
- [ ] API versioning maintained
- [ ] Error messages don't leak info

## Infrastructure
- [ ] Firewall rules configured
- [ ] SSL/TLS certificates valid
- [ ] Database access restricted
- [ ] Redis password set
- [ ] Log rotation configured

## Testing
- [ ] Security testing performed
- [ ] Penetration testing completed
- [ ] Vulnerability scan passed
- [ ] Code review completed
- [ ] Security documentation updated
```

### Ongoing Security Practices

```markdown
## Regular Tasks
- [ ] Weekly security audit
- [ ] Monthly penetration testing
- [ ] Quarterly security review
- [ ] Annual security assessment
- [ ] Regular security training

## Continuous Monitoring
- [ ] Security logs reviewed
- [ ] Vulnerability scans run
- [ ] Access logs monitored
- [ ] Performance metrics tracked
- [ ] User activity audited

## Incident Preparedness
- [ ] Incident response plan documented
- [ ] Contact information updated
- [ ] Backup procedures tested
- [ ] Recovery plans validated
- [ ] Communication plans established
```

---

*This security guide should be reviewed and updated regularly to address new threats and best practices.*
