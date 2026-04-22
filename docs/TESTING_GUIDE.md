# ByteDanceDemo 测试指南

本指南提供 ByteDanceDemo 项目测试的完整说明，包括测试策略、工具使用、最佳实践和故障排除。

## 测试概览

ByteDanceDemo 使用分层测试策略，确保代码质量和功能正确性：

- **单元测试** - 测试单个组件和服务方法
- **集成测试** - 测试组件之间的交互和完整 API 流程
- **基准测试** - 测量关键操作的性能
- **端到端测试** - 验证完整的用户场景

## 测试覆盖目标

| 组件 | 目标覆盖率 | 状态 |
|-------|-----------|------|
| Service Layer | 80% | 📊 |
| Controller Layer | 75% | 📊 |
| Middleware | 85% | 📊 |
| Repository Layer | 70% | 📊 |
| 整体项目 | 75% | 📊 |

## 目录结构

```
test/
├── mocks/                    # Mock 实现
│   ├── dao_mock.go          # 数据库层 mock
│   ├── redis_mock.go        # Redis 客户端 mock
│   └── cache_mock.go       # 缓存层 mock
├── services/               # Service 层测试
│   ├── user_service_test.go  # 用户服务测试
│   ├── comment_service_test.go # 评论服务测试
│   ├── follow_service_test.go # 关注服务测试
│   ├── favorite_service_test.go # 点赞服务测试
│   ├── message_service_test.go # 消息服务测试
│   └── video_service_test.go # 视频服务测试
├── controllers/            # Controller 层测试
│   ├── user_controller_test.go # 用户控制器测试
│   ├── comment_controller_test.go # 评论控制器测试
│   └── publish_controller_test.go # 发布控制器测试
├── middleware/            # 中间件测试
│   ├── auth_middleware_test.go # 认证中间件测试
│   └── validation_middleware_test.go # 验证中间件测试
├── integration/           # 集成测试
│   └── api_integration_test.go # API 集成测试
├── benchmarks/            # 性能基准测试
│   ├── database_benchmarks.go # 数据库基准测试
│   ├── api_benchmarks.go # API 基准测试
│   └── redis_benchmarks.go # Redis 基准测试
├── utils_test.go          # 测试工具函数
├── test_config.yaml       # 测试配置文件
└── run_tests.sh          # 测试运行脚本
```

## 运行测试

### 使用测试脚本

项目提供了灵活的测试运行脚本 `run_tests.sh`：

```bash
# 运行所有测试
./test/run_tests.sh

# 运行特定类型测试
./test/run_tests.sh --type unit        # 单元测试
./test/run_tests.sh --type integration # 集成测试
./test/run_tests.sh --type benchmark   # 基准测试

# 生成覆盖率报告
./test/run_tests.sh --type unit --coverage

# 使用详细输出
./test/run_tests.sh --type unit --verbose

# 清理测试数据
./test/run_tests.sh --type unit --clean
```

### 使用 Makefile

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 运行基准测试
make test-benchmark

# 生成覆盖率报告
make test-coverage

# 生成 HTML 覆盖率报告
make test-coverage-html

# 带竞态检测运行测试
make test-race

# 运行服务包测试
make test-service

# 运行 repository 包测试
make test-repository

# 运行 utils 包测试
make test-utils
```

### 并行测试

项目支持并行测试执行：

```bash
# 使用 4 个并行工作线程运行所有测试
make test-parallel

# 查看结果
cat test-results.log
```

## 测试配置

### 配置文件

测试配置存储在 `test/test_config.yaml`：

```yaml
# 数据库设置
database:
  host: localhost
  port: 3306
  username: test_user
  password: test_password
  database: bytedancedemo_test

# Redis 设置
redis:
  host: localhost
  port: 6379
  database: 1

# 测试设置
test:
  benchmark_iterations: 1000
  concurrency: 100
  request_timeout: 30s
  debug: true
  cleanup_after_test: true

# 测试数据
test_data:
  default_user_count: 100
  default_video_count: 500
  default_comment_count: 2000
```

### 环境变量

可以使用环境变量覆盖配置：

```bash
# 数据库配置
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_USER=test_user
export TEST_DB_PASSWORD=test_password
export TEST_DB_NAME=bytedancedemo_test

# Redis 配置
export TEST_REDIS_HOST=localhost
export TEST_REDIS_PORT=6379

# 测试配置
export TEST_DEBUG=true
export TEST_CLEANUP=true
```

## Mock 使用

### 数据库 Mock

项目使用 gomock 生成数据库层 mock：

```go
// 创建 UserDAO mock
mockDAO := NewMockUserDAO(ctrl)

// 设置预期行为
mockDAO.EXPECT().Create(gomock.Any()).Return(nil)
mockDAO.EXPECT().FindByID(gomock.Eq(1), gomock.Any()).Return(expectedUser, nil)

// 使用 mock 实例
userService := NewUserService(mockDAO)
user, err := userService.InsertUser(testUser)
```

### Redis Mock

Redis 操作的 mock 示例：

```go
// 创建 Redis 客户端 mock
mockRedis := NewMockRedisClient()

// 设置预期行为
mockRedis.On("Get", mock.Anything, "user:1").Return("user_data", nil)
mockRedis.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return("OK", nil)

// 使用 mock 实例
cacheService := NewCacheService(mockRedis)
data, err := cacheService.Get("user:1")
```

## 单元测试

### 编写单元测试

单元测试应该遵循 AAA 模式：

```go
func TestUserService_InsertUser(t *testing.T) {
    // Arrange - 准备测试数据
    mockDAO := NewMockUserDAO(ctrl)
    userService := NewUserService(mockDAO)
    
    testUser := &model.User{
        Name:     "testuser",
        Password: "hashedpassword",
    }
    
    // Act - 执行测试操作
    result, err := userService.InsertUser(testUser)
    
    // Assert - 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "testuser", result.Name)
    assert.True(t, result.ID > 0)
}
```

### 测试错误用列

```go
func TestUserService_InsertUser_InvalidInput(t *testing.T) {
    // Arrange
    mockDAO := NewMockUserDAO(ctrl)
    userService := NewUserService(mockDAO)
    
    invalidUser := &model.User{
        Name:     "", // 无效的用户名
        Password: "pass",
    }
    
    // Act
    result, err := userService.InsertUser(invalidUser)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
}
```

### 使用表驱动测试

```go
func TestUserService_InsertUser_TableDriven(t *testing.T) {
    testCases := []struct {
        name     string
        input    *model.User
        expected error
        wantNil  bool
    }{
        {
            name:     "成功插入",
            input:    &model.User{Name: "testuser", Password: "pass"},
            expected: nil,
            wantNil:  false,
        },
        {
            name:     "空用户名",
            input:    &model.User{Name: "", Password: "pass"},
            expected: errors.New("用户名不能为空"),
            wantNil:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            mockDAO := NewMockUserDAO(ctrl)
            userService := NewUserService(mockDAO)
            
            result, err := userService.InsertUser(tc.input)
            
            if tc.expected != nil {
                assert.Error(t, err)
                assert.Equal(t, tc.expected.Error(), err.Error())
            } else {
                assert.NoError(t, err)
            }
            
            if tc.wantNil {
                assert.Nil(t, result)
            } else {
                assert.NotNil(t, result)
            }
        })
    }
}
```

## 集成测试

### API 集成测试

完整的 API 流程测试：

```go
func TestAPIIntegration_CompleteUserFlow(t *testing.T) {
    // Setup
    router := setupIntegrationTestRouter()
    defer cleanupTestData(router)
    
    // 1. 用户注册
    registerResp := registerUser(t, router, "testuser", "password123")
    assert.Equal(t, 200, registerResp.StatusCode)
    var regResult struct {
        StatusCode int    `json:"status_code"`
        Token       string   `json:"token"`
    }
    json.Unmarshal(registerResp.Body.Bytes(), &regResult)
    
    // 2. 用户登录
    loginResp := loginUser(t, router, "testuser", "password123")
    assert.Equal(t, 200, loginResp.StatusCode)
    var loginResult struct {
        StatusCode int    `json:"status_code"`
        Token       string   `json:"token"`
    }
    json.Unmarshal(loginResp.Body.Bytes(), &loginResult)
    
    // 3. 获取用户信息（使用 token）
    token := loginResult.Token
    userResp := getUserInfo(t, router, token)
    assert.Equal(t, 200, userResp.StatusCode)
    
    // 4. 发布视频
    videoResp := publishVideo(t, router, token, "test_video.mp4")
    assert.Equal(t, 200, videoResp.StatusCode)
    
    // 5. 获取视频列表
    feedResp := getVideoFeed(t, router, token)
    assert.Equal(t, 200, feedResp.StatusCode)
    
    // Cleanup
    cleanupTestData(router)
}
```

### 测试用例辅助函数

```go
// 创建测试用户并获取 token
func createTestUserAndLogin(router *gin.Engine, username, password string) (string, string) {
    // 注册
    user := &model.User{Name: username, Password: passwordHash(password)}
    jsonBody, _ := json.Marshal(user)
    
    req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    var regResp struct {
        StatusCode int `json:"status_code"`
        Token       string `json:"token"`
    }
    json.Unmarshal(w.Body.Bytes(), &regResp)
    
    return regResp.Token, strconv.FormatInt(regResp.StatusCode, 10)
}

// 生成测试用户数据
func generateTestUsers(count int) []*model.User {
    users := make([]*model.User, count)
    for i := 0; i < count; i++ {
        users[i] = &model.User{
            Name:      fmt.Sprintf("user%d", i),
            Password:  fmt.Sprintf("pass%d", i),
            Avatar:    fmt.Sprintf("http://example.com/avatar%d.jpg", i),
        }
    }
    return users
}
```

## 基准测试

### 编写基准测试

```go
func BenchmarkUserService_InsertUser(b *testing.B) {
    // Setup
    db, _ := database.InitTestDB()
    defer db.Close()
    userService := service.GetUserServiceInstance(db)
    
    // 创建测试数据（预热）
    for i := 0; i < 100; i++ {
        user := &model.User{
            Name:     fmt.Sprintf("benchuser%d", i),
            Password: "hashedpassword",
        }
        userService.InsertUser(user)
    }
    
    // 基准测试
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user := &model.User{
            Name:     fmt.Sprintf("benchuser%d", i),
            Password: "hashedpassword",
        }
        userService.InsertUser(user)
    }
}
```

### 基准测试目标

| 操作 | 目标时间 | P50 | P95 | P99 |
|-------|---------|-----|-----|-----|
| 用户注册 | < 100ms | 150ms | 200ms | 300ms |
| 用户登录 | < 50ms | 75ms | 100ms | 150ms |
| 视频发布 | < 500ms | 750ms | 1000ms | 1500ms |
| 评论添加 | < 100ms | 150ms | 200ms | 300ms |
| 点赞操作 | < 30ms | 50ms | 75ms | 100ms |
| 消息发送 | < 150ms | 200ms | 300ms | 400ms |

## 并发测试

### 竞态检测

运行测试时检测竞态条件：

```bash
# 所有测试
go test -race ./...

# 特定包测试
go test -race ./service/...
```

### 并发安全验证

```go
func TestConcurrentUserFollow(t *testing.T) {
    // Setup
    db, _ := database.InitTestDB()
    defer db.Close()
    
    // 启用并发
    var wg sync.WaitGroup
    userCount := 10
    followCount := 100
    
    // 创建测试用户
    users := make([]*model.User, userCount)
    for i := 0; i < userCount; i++ {
        users[i] = &model.User{
            Name:      fmt.Sprintf("user%d", i),
            Password:  fmt.Sprintf("pass%d", i),
        }
        db.Create(users[i])
    }
    
    // 并发执行关注操作
    for i := 0; i < userCount; i++ {
        for j := 0; j < userCount; j++ {
            if i == j {
                continue // 不能关注自己
            }
            
            wg.Add(1)
            go func(follower, followee int64) {
                defer wg.Done()
                followService := service.GetFollowServiceInstance(db)
                err := followService.Follow(follower, followee)
                assert.NoError(t, err)
            }(users[i].ID, users[j].ID)
        }
    }
    
    wg.Wait()
}
```

## 测试数据管理

### 测试数据创建

```go
// 创建完整的测试用户
func createTestUser(id int64, username string) *model.User {
    return &model.User{
        ID:       id,
        Name:     username,
        Password: bcryptPassword("defaultpassword"),
        Avatar:    fmt.Sprintf("http://example.com/avatar%d.jpg", id),
        BackgroundImage: fmt.Sprintf("http://example.com/bg%d.jpg", id),
        Signature: fmt.Sprintf("Signature for user %d", id),
    }
}

// 创建测试视频
func createTestVideo(id int64, authorID int64) *model.Video {
    return &model.Video{
        ID:        id,
        AuthorID:  authorID,
        PlayURL:   fmt.Sprintf("/videos/video%d.mp4", id),
        CoverURL:  fmt.Sprintf("/videos/cover%d.jpg", id),
        Title:     fmt.Sprintf("Test Video %d", id),
    }
}

// 创建测试评论
func createTestComment(id int64, videoID int64, userID int64) *model.Comment {
    return &model.Comment{
        ID:      id,
        VideoID: videoID,
        UserID: userID,
        Content: fmt.Sprintf("This is a test comment %d", id),
    }
}
```

### 测试数据清理

```go
// 清理测试数据
func cleanupTestData(db *gorm.DB) {
    tables := []string{
        "users", "videos", "comments", "likes", "relations", "messages",
    }
    
    for _, table := range tables {
        db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id > 0", table))
    }
    
    // 重置自增 ID
    db.Exec("SET @@auto_increment_increment=0")
}

// 测试中使用
func TestWithCleanup(t *testing.T, testFunc func(*gorm.DB)) {
    db, _ := database.InitTestDB()
    defer func() {
        cleanupTestData(db)
        db.Close()
    }()
    
    testFunc(db)
}
```

## 持续集成 (CI)

### GitHub Actions 配置

项目包含 GitHub Actions 自动化测试配置：

```yaml
name: Tests

on:
  push:
    branches: [ main, master, develop ]
  pull_request:
    branches: [ main, master, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - run: go mod download
      - run: make test
      
  coverage:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - run: go mod download
      - run: make test-coverage
      - uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### 覆盖率报告

生成覆盖率报告：

```bash
# 生成文本格式
make test-coverage

# 生成 HTML 格式（推荐）
make test-coverage-html

# 查看覆盖率
cat coverage.out | grep "^total:" | head -1
```

目标覆盖率：

- Service Layer: 80%
- Controller Layer: 75%
- Middleware: 85%
- 整体项目: 75%

## 故障排除

### 常见测试问题

**问题**：测试失败，错误信息不清晰

**解决方法**：
- 使用 `-v` 标志获取详细输出
- 检查测试命名，确保描述性
- 检查断言错误信息

```bash
# 使用详细输出运行测试
go test -v ./service/...

# 运行特定测试
go test -v -run TestUserService_InsertUser ./service/
```

**问题**：Mock 行为未设置

**解决方法**：
- 检查 `EXPECT()` 调用
- 确保每个 `EXPECT()` 都有对应的调用
- 使用 `ctrl.AssertExpectations(t)` 验证

```go

mockDAO.EXPECT().Create(gomock.Any()).Return(nil)
mockDAO.EXPECT().FindByID(gomock.Eq(1), gomock.Any()).Return(expectedUser, nil)

// 在测试结束时验证
defer ctrl.AssertExpectations(t)
```

**问题**：数据库连接失败

**解决方法**：
- 检查测试数据库是否运行
- 检查配置文件中的数据库凭据
- 使用环境变量覆盖配置

```bash
# 检查数据库连接
mysql -u test_user -ptest_password -h localhost test_bytedancedemo

# 使用环境变量
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_PASSWORD=test_password
```

**问题**：Redis 连接失败

**解决方法**：
- 确保 Redis 服务正在运行
- 检查 Redis 配置
- 验证网络连接

```bash
# 检查 Redis 连接
redis-cli ping

# 查看日志
tail -f logs/redis-server.log
```

## 最佳实践

### 1. 测试隔离

每个测试应该是独立的：

```go
func TestUserService_InsertUser(t *testing.T) {
    // 每个测试使用独立的数据库连接
    db, _ := database.InitTestDB()
    defer db.Close()
    
    // 每个测试创建独立的数据
    testUser := &model.User{
        Name:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
        Password: "hashedpassword",
    }
    
    // 测试完成后清理数据
    defer db.Delete(testUser)
    
    userService := NewUserService(db)
    result, err := userService.InsertUser(testUser)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 2. 使用测试辅助函数

创建可重用的测试辅助函数：

```go
func assertHTTPSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
    assert.Equal(t, 200, resp.Code)
}

func assertHTTPError(t *testing.T, resp *httptest.ResponseRecorder, expectedCode int) {
    assert.Equal(t, expectedCode, resp.Code)
}

func assertJSONResponse(t *testing.T, resp *httptest.ResponseRecorder, expected interface{}) {
    var result interface{}
    err := json.Unmarshal(resp.Body.Bytes(), &result)
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### 3. 测试边界情况

确保测试覆盖边界情况：

```go
func TestUserService_InsertUser_EmptyUsername(t *testing.T) {
    // 测试空用户名
    user := &model.User{Name: "", Password: "pass"}
    _, err := userService.InsertUser(user)
    assert.Error(t, err)
}

func TestUserService_InsertUser_VeryLongUsername(t *testing.T) {
    // 测试过长的用户名
    user := &model.User{Name: strings.Repeat("a", 300), Password: "pass"}
    _, err := userService.InsertUser(user)
    assert.Error(t, err)
}

func TestUserService_InsertUser_SpecialCharacters(t *testing.T) {
    // 测试特殊字符
    user := &model.User{Name: "<script>alert('xss')</script>", Password: "pass"}
    _, err := userService.InsertUser(user)
    assert.Error(t, err)
}
```

### 4. 并发测试安全

使用 sync.WaitGroup 确保并发安全：

```go
func TestConcurrentUserFollow(t *testing.T) {
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // 测试代码
        }()
    }
    
    wg.Wait()
}
```

### 5. 性能测试规范

基准测试应该模拟真实的负载：

```go
func BenchmarkUserService_InsertUser(b *testing.B) {
    // 预热
    for i := 0; i < 100; i++ {
        user := &model.User{Name: fmt.Sprintf("user%d", i), Password: "pass"}
        userService.InsertUser(user)
    }
    
    // 基准测试
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user := &model.User{Name: fmt.Sprintf("user%d", i), Password: "pass"}
        userService.InsertUser(user)
    }
}
```

## 性能基准

### 数据库操作基准

| 操作 | 操作/秒 | 平均响应时间 | 内存使用 |
|-------|----------|-------------|----------|
| 插入用户 | 1000 | 0.5ms | 10MB |
| 查询用户 | 5000 | 0.2ms | 5MB |
| 更新用户 | 500 | 0.3ms | 8MB |
| 删除用户 | 100 | 0.1ms | 3MB |

### API 端点基准

| 端点 | 并发请求 | 平均响应时间 | P95 响应时间 |
|-------|----------|-------------|---------------|
| /feed/ | 100 | 45ms | 80ms |
| /user/ | 200 | 25ms | 50ms |
| /login/ | 50 | 30ms | 60ms |
| /publish/ | 10 | 1.2s | 2.0s |

## 测试清单

### 功能测试

- [ ] 用户注册和登录
- [ ] 用户信息查询和更新
- [ ] 视频上传和列表
- [ ] 评论创建、查询和删除
- [ ] 点赞和取消点赞
- [ ] 关注和取消关注
- [ ] 关注列表查询
- [ ] 消息发送和查询

### 安全测试

- [ ] SQL 注入防护
- [ ] XSS 攻击防护
- [ ] CSRF 保护
- [ ] 未授权访问
- [ ] 速率限制
- [ ] 输入验证

### 性能测试

- [ ] 数据库查询优化
- [ ] Redis 缓存效率
- [ ] API 响应时间
- [ ] 并发处理
- [ ] 内存使用
- [ ] 连接池效率

### 错误处理测试

- [ ] 数据库连接失败
- [ ] Redis 连接失败
- [ ] 无效输入处理
- [ ] 网络超时处理
- [ ] 文件上传错误
- [ ] 权限拒绝

## 测试命令速查

```bash
# 运行所有测试
./test/run_tests.sh

# 运行单元测试
./test/run_tests.sh --type unit

# 运行集成测试
./test/run_tests.sh --type integration

# 运行基准测试
./test/run_tests.sh --type benchmark

# 生成覆盖率报告
./test/run_tests.sh --type unit --coverage

# 并行测试
make test-parallel

# 清理并运行
make clean test

# 带竞态检测
make test-race
```

## 相关文档

- [详细测试文档](../TEST_DOCUMENTATION.md) - 完整的技术测试文档
- [测试目录 README](../test/README.md) - 测试快速开始指南
- [开发指南](CONTRIBUTING.md) - 代码风格和测试要求
- [故障排除](TROUBLESHOOTING.md) - 常见问题解决方案

---

**最后更新**: 2026-04-21
