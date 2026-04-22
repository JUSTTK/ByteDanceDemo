# ByteDanceDemo 常见问题

## 安装和配置

### Q: 如何设置项目？
A: 按照以下步骤设置：
1. 克隆项目：`git clone https://github.com/yourusername/ByteDanceDemo.git`
2. 复制配置文件：`cp config/settings.yml.template config/settings.yml`
3. 安装依赖：`go mod download`
4. 创建数据库：`mysql -u root -p < config/init.sql`
5. 运行迁移：`make run-migrate`
6. 启动服务：`make run-api`

详情请查看 [README.md](../README.md)。

### Q: 如何配置数据库连接？
A: 编辑 `config/settings.yml` 文件：
```yaml
settings:
  mysql:
    host: "localhost"
    port: 3306
    schema: "sample_douyin"
    username: "sample_douyin"
    password: "your_password"
```
或使用环境变量：
```bash
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_DATABASE=sample_douyin
export MYSQL_USERNAME=sample_douyin
export MYSQL_PASSWORD=your_password
```

### Q: 如何解决 "address already in use" 错误？
A: 端口 8080 已被占用。解决方法：
1. 停止占用端口的进程：
   ```bash
   sudo lsof -ti:8080
   sudo kill -9 $(lsof -ti:8080)
   ```
2. 或修改配置文件中的端口号：
   ```yaml
   settings:
     application:
       port: 8081  # 使用其他端口
   ```

### Q: 如何配置国内 Go 代理？
A: 设置环境变量：
```bash
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off
```
或添加到 `~/.bashrc` 或 `~/.zshrc` 使其永久生效。

### Q: Redis 连接失败怎么办？
A: 检查步骤：
1. 确认 Redis 正在运行：`redis-cli ping`
2. 检查配置文件中的 Redis 地址：`grep -A 5 redis: config/settings.yml`
3. 检查防火墙设置：`sudo ufw allow 6379/tcp`
4. 确认 Redis 密码是否正确

## API 使用

### Q: 如何获取 JWT token？
A: 首先注册用户，然后登录：
```bash
# 注册
curl -X POST "http://localhost:8080/douyin/user/register/" \
  -H "Content-Type: multipart/form-data" \
  -F "username=testuser" \
  -F "password=secure123"

# 登录
curl -X POST "http://localhost:8080/douyin/user/login/" \
  -H "Content-Type: multipart/form-data" \
  -F "username=testuser" \
  -F "password=secure123"
```
登录成功后会返回 JWT token，包含在响应的 `token` 字段中。

### Q: 如何在请求中使用 JWT token？
A: 在请求头中添加 Authorization 字段：
```bash
curl -X GET "http://localhost:8080/douyin/user/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
或在查询参数中：
```bash
curl -X GET "http://localhost:8080/douyin/user/?token=YOUR_JWT_TOKEN"
```

### Q: 如何上传视频？
A: 使用 multipart/form-data 格式：
```bash
curl -X POST "http://localhost:8080/douyin/publish/action/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "data=@video.mp4" \
  -F "token=YOUR_JWT_TOKEN"
```
视频会保存到 `public/` 目录。

### Q: API 返回 401 Unauthorized 怎么办？
A: 可能的原因：
1. JWT token 缺失或无效
2. Token 已过期
3. Token 格式错误

解决方法：
1. 确认 token 正确包含在请求中
2. 检查 token 是否过期：`echo "YOUR_JWT_TOKEN" | base64 -d | jq .`
3. 重新登录获取新 token
4. 检查配置文件中的 JWT 密钥是否正确

## 性能问题

### Q: API 响应很慢怎么办？
A: 可能的原因和解决方法：
1. **数据库查询慢**
   - 检查 MySQL 慢查询日志：`tail -f /var/log/mysql/slow-query.log`
   - 添加数据库索引
   - 使用 `EXPLAIN` 分析查询

2. **缓存未命中**
   - 检查 Redis 连接状态
   - 验证缓存配置
   - 增加缓存 TTL

3. **并发请求过多**
   - 检查 goroutine 数量：`curl http://localhost:8080/debug/pprof/goroutine`
   - 实现连接池
   - 添加速率限制

4. **网络延迟**
   - 检查数据库服务器连接
   - 优化网络配置

### Q: 如何启用 pprof 进行性能分析？
A: 应用已内置 pprof 支持：
```bash
# CPU 分析
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof

# 内存分析
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# goroutine 分析
curl http://localhost:8080/debug/pprof/goroutine?debug=1 > goroutine.prof
go tool pprof goroutine.prof
```

### Q: 如何优化数据库连接池？
A: 在配置文件中设置连接池参数：
```yaml
settings:
  mysql:
    maxOpenConns: 100      # 最大打开连接数
    maxIdleConns: 20       # 最大空闲连接数
    connMaxLifetime: 30m    # 连接最大生命周期
    connMaxIdleTime: 5m     # 连接最大空闲时间
```
或使用环境变量：
```bash
export MYSQL_MAX_OPEN_CONNS=100
export MYSQL_MAX_IDLE_CONNS=20
```

## 安全问题

### Q: 如何生成安全的 JWT 密钥？
A: 使用提供的脚本：
```bash
./scripts/setup_jwt.sh
```
或手动生成：
```bash
# 生成 32 字节的随机密钥
SECRET_KEY=$(openssl rand -base64 32 | tr -d '=' | tr -d '/' | cut -c1-32)
echo "Generated JWT secret: $SECRET_KEY"

# 更新配置文件
sed -i "s/secretKey:.*/secretKey: $SECRET_KEY/" config/settings.yml
```

### Q: 如何检查 JWT 密钥安全性？
A: 使用安全检查工具：
```bash
go run scripts/check_jwt_security.go
```
检查项包括：
- 密钥长度（建议至少 32 字符/256 位）
- 弱密钥模式检测
- Base64 编码检查
- 哈希模式检测
- 非法配置验证

### Q: 如何防止 SQL 注入？
A: 项目已实施多层防护：
1. 使用 GORM ORM（自动参数化）
2. 输入验证中间件
3. 使用 `c.PostForm()` 而非 `c.Query()`
4. 参数化查询而非原始 SQL

安全示例：
```go
// ✅ 安全 - 使用 GORM 参数化
var user User
result := db.Where("username = ?", username).First(&user)

// ❌ 不安全 - 直接拼接 SQL
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
result := db.Raw(query)
```

### Q: 如何启用 CSRF 保护？
A: 项目已包含 CSRF 中间件，确保：
1. 前端在请求中包含 CSRF token
2. 在表单中添加隐藏字段：
   ```html
   <input type="hidden" name="csrf_token" value="{{ .csrf_token }}">
   ```
3. 在请求头中包含：
   ```bash
   curl -X POST "http://localhost:8080/api" \
     -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
   ```

## 开发问题

### Q: 如何运行单元测试？
A: 使用 Makefile 或直接运行：
```bash
# 使用 Makefile
make test-unit

# 直接运行
go test -v ./service/...

# 运行特定测试
go test -v ./service/userService_test.go -run TestUserServiceInsertUser
```

### Q: 如何生成测试覆盖率报告？
A: 使用 Makefile 目标：
```bash
make test-coverage
```
覆盖率报告会生成在 `coverage.html` 文件中。

### Q: 如何调试 Go 代码？
A: 启动应用时使用 debug 模式：
```bash
./bin/simple-demo server -c config/settings.yml -m debug
```
或在 IDE 中设置断点和调试器。

### 使用 dlv：
```bash
# 安装 dlv
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试运行中的应用
dlv attach $(pgrep simple-demo)

# 或调试测试
dlv test ./service/userService_test.go
```

### Q: 如何生成数据库 DAO？
A: 运行生成器：
```bash
go run model/main/main.go
```
修改数据模型后需要重新生成 DAO。

## 测试问题

### Q: 集成测试失败怎么办？
A: 常见原因和解决方法：
1. **API 服务未运行**
   - 启动 API 服务：`make run-api` 或使用 docker-compose
   - 检查服务健康状态：`curl http://localhost:8080/health`

2. **数据库未初始化**
   - 运行迁移：`make run-migrate`
   - 检查数据库是否存在：`mysql -u sample_douyin -p -e "SHOW DATABASES;"`

3. **Redis 未运行**
   - 启动 Redis：`sudo systemctl start redis-server`
   - 或使用 Docker：`docker run -d -p 6379:6379 redis:7-alpine`

4. **配置不匹配**
   - 检查测试配置文件：`test/test_config.yaml`
   - 使用环境变量覆盖配置

### Q: 如何运行性能基准测试？
A: 使用 benchmark 标志：
```bash
make test-benchmark
```
或运行特定基准：
```bash
go test -bench=BenchmarkUserServiceInsertUser -benchtime=10s ./service/
```

## Docker 问题

### Q: 如何使用 Docker 运行项目？
A: 使用 docker-compose：
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止所有服务
docker-compose down
```

### Q: Docker 容器无法连接数据库怎么办？
A: 检查网络配置：
1. 确认服务在同一网络：`docker network ls`
2. 检查 docker-compose.yml 中的服务名称
3. 使用服务名称而非 localhost 连接：
   ```yaml
   mysql:
     host: mysql  # 而非 localhost
   ```

### Q: 如何进入 Docker 容器调试？
A: 使用 docker-compose：
```bash
# 进入应用容器
docker-compose exec app sh

# 进入数据库容器
docker-compose exec mysql bash

# 查看容器日志
docker-compose logs app
```

## 更多帮助

如果以上 FAQ 没有解决你的问题，请：

1. 查看 [故障排除文档](TROUBLESHOOTING.md)
2. 查看 [架构文档](ARCHITECTURE.md)
3. 查看 [配置文档](CONFIGURATION.md)
4. 在 [GitHub Issues](https://github.com/yourusername/ByteDanceDemo/issues) 中提问

---

**最后更新**: 2026-04-21
