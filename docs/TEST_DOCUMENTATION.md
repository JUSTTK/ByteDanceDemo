# ByteDanceDemo - Test Documentation

## Overview

This document provides comprehensive information about the testing setup, execution, and maintenance for the ByteDanceDemo project.

## Test Structure

```
test/
├── mocks/                    # Mock implementations
│   ├── dao_mock.go          # Database mocks
│   ├── redis_mock.go        # Redis mocks
│   └── cache_mock.go       # Cache mocks
├── services/               # Service layer tests
│   ├── user_service_test.go
│   ├── comment_service_test.go
│   ├── follow_service_test.go
│   ├── favorite_service_test.go
│   ├── message_service_test.go
│   └── video_service_test.go
├── controllers/            # Controller layer tests
│   ├── user_controller_test.go
│   ├── comment_controller_test.go
│   ├── publish_controller_test.go
│   └── ...
├── middleware/            # Middleware tests
│   ├── auth_middleware_test.go
│   └── validation_middleware_test.go
├── integration/           # Integration tests
│   └── api_integration_test.go
├── benchmarks/            # Performance benchmarks
│   ├── database_benchmarks.go
│   ├── api_benchmarks.go
│   └── redis_benchmarks.go
├── utils_test.go          # Test utilities
└── test_config.yaml       # Test configuration
```

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run benchmark tests
make test-benchmark

# Run tests with coverage
make test-coverage

# Run tests with HTML coverage report
make test-coverage-html

# Run tests with race detector
make test-race

# Run tests in short mode (skip long tests)
make test-short

# Run tests with verbose output
make test-verbose
```

### Running Specific Tests

```bash
# Run tests for a specific package
make test-package PACKAGE=test/services

# Run a specific test function
make test-function FUNCTION=TestUserServiceInsertUser

# Run tests with specific tags
make test-tags TAGS="integration"
```

## Test Coverage

### Coverage Requirements

- **Service Layer**: Minimum 80% coverage
- **Controller Layer**: Minimum 75% coverage
- **Middleware**: Minimum 85% coverage
- **Overall Project**: Minimum 75% coverage

### Generating Coverage Reports

```bash
# Generate coverage report
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

Coverage reports will be generated as:
- `coverage.out` - Plain text format
- `coverage.html` - Interactive HTML report

## Mock Implementations

### Database Mocks

The project uses mock implementations for database operations to enable isolated unit testing:

```go
// Example: Using database mocks
mockDAO := &MockUserDAO{}
mockDAO.EXPECT().Create(user).Return(nil)
mockDAO.EXPECT().Where("name", "test").Return(mockDAO)
mockDAO.EXPECT().Find().Return([]interface{}{user}, nil)
```

### Redis Mocks

Redis operations are mocked to prevent dependency on external Redis instances during unit tests:

```go
// Example: Using Redis mocks
mockRedis := NewMockRedisClient()
mockRedis.On("Get", mock.Anything, "user:1").Return("user_data", nil)
mockRedis.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return("OK", nil)
```

## Test Categories

### Unit Tests

Unit tests focus on individual components in isolation:

- **Service Layer**: Tests business logic without external dependencies
- **Controller Layer**: Tests HTTP handlers with mocked services
- **Middleware**: Tests request/response processing

Example:
```go
func TestUserServiceInsertUser(t *testing.T) {
    userService := service.GetUserServiceInstance()
    user := &model.User{
        Name:     "testuser",
        Password: "hashed_password",
    }

    result, success := userService.InsertUser(user)

    assert.True(t, success)
    assert.NotNil(t, result)
    assert.Equal(t, "testuser", result.Name)
}
```

### Integration Tests

Integration tests verify that components work together correctly:

```go
func TestIntegrationAPI(t *testing.T) {
    router := setupIntegrationTestRouter()

    // Test complete user registration flow
    t.Run("User Registration and Login Flow", func(t *testing.T) {
        // Register user
        // Login user
        // Verify token generation
        // Assert successful flow
    })
}
```

### Benchmark Tests

Benchmark tests measure performance of critical operations:

```go
func BenchmarkUserServiceInsertUser(b *testing.B) {
    userService := service.GetUserServiceInstance()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        user := &model.User{
            Name:     "benchuser",
            Password: "password",
        }
        userService.InsertUser(user)
    }
}
```

## Test Configuration

### Test Configuration File

The `test/test_config.yaml` file contains test-specific settings:

```yaml
# Database settings for tests
database:
  driver: mysql
  host: localhost
  port: 3306
  username: test_user
  password: test_password
  database: bytedancedemo_test

# Test settings
test:
  benchmark_iterations: 1000
  concurrency: 100
  request_timeout: 30s
  debug: true
  cleanup_after_test: true
```

### Environment Variables

Tests can be configured using environment variables:

```bash
# Test database
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_USER=test_user
export TEST_DB_PASSWORD=test_password
export TEST_DB_NAME=bytedancedemo_test

# Test Redis
export TEST_REDIS_HOST=localhost
export TEST_REDIS_PORT=6379

# Test settings
export TEST_DEBUG=true
export TEST_CLEANUP=true
```

## Test Data Management

### Test Data Creation

Test data is created using helper functions:

```go
// Create test users
users := generateTestUsers(100)

// Create test videos
videos := generateTestVideos(1000)

// Create test comments
comments := generateTestComments(5000)
```

### Test Data Cleanup

Test data cleanup is handled automatically:

```go
func cleanupTestData(router *gin.Engine) {
    // Clean up any test data created during tests
    // Example: Delete test users, videos, comments, etc.
}
```

### Data Isolation

Each test should:
1. Create its own test data
2. Use unique identifiers to avoid conflicts
3. Clean up after completion
4. Not depend on data from other tests

## Error Handling Tests

### Test Categories

1. **Database Errors**: Test handling of database connection failures, query errors
2. **Validation Errors**: Test input validation failures
3. **Authentication Errors**: Test invalid token handling
4. **Authorization Errors**: Test permission failures
5. **Network Errors**: Test timeout and connection issues

### Example

```go
func TestUserController_RegisterUser_InvalidInput(t *testing.T) {
    router := setupUserControllerTest()

    // Test with missing fields
    userData := map[string]interface{}{
        "username": "", // Empty username
        "password": "password123",
    }
    jsonData, _ := json.Marshal(userData)

    req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

## Authentication & Authorization Tests

### JWT Token Tests

```go
func TestJWTMiddleware_TokenValidation(t *testing.T) {
    // Test valid token
    token, _ := jwt.GenerateToken(1, "testuser")
    claims, err := jwt.ParseToken(token)
    assert.NoError(t, err)
    assert.Equal(t, int64(1), claims.UserID)

    // Test invalid token
    claims, err = jwt.ParseToken("invalid.token.here")
    assert.Error(t, err)
    assert.Nil(t, claims)
}
```

### Authorization Tests

```go
func TestCommentController_DeleteComment_Unauthorized(t *testing.T) {
    // Test that only comment owners can delete comments
    // Test that users cannot delete other users' comments
    // Test permission checks
}
```

## Performance Testing

### Benchmark Categories

1. **Database Operations**: Query performance, connection pooling
2. **Redis Operations**: Cache performance, pipeline operations
3. **API Performance**: Response times, throughput
4. **Concurrent Operations**: Multi-user scenarios

### Running Benchmarks

```bash
# Run all benchmarks
make test-benchmark

# Run specific benchmark
go test -bench=BenchmarkUserServiceInsertUser -benchtime=10s ./...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./...
```

### Performance Benchmarks

The project should meet these performance benchmarks:

- **User Registration**: < 100ms
- **User Login**: < 50ms
- **Video Publishing**: < 500ms
- **Comment Addition**: < 100ms
- **Follow/Unfollow**: < 100ms
- **Message Sending**: < 150ms
- **Database Queries**: < 10ms (simple), < 100ms (complex)
- **Redis Operations**: < 5ms (get/set), < 50ms (complex)

## Continuous Integration

### CI Pipeline

The project includes CI pipeline for automated testing:

```yaml
# Example CI pipeline
stages:
  - test
  - coverage
  - benchmark

test:
  script:
    - make test
    - make test-race

coverage:
  script:
    - make test-coverage
    - make test-coverage-html

benchmark:
  script:
    - make test-benchmark
```

### Test Reports

CI generates test reports:

- **Test Results**: XML/JSON format for CI tools
- **Coverage Reports**: HTML and Cobertura formats
- **Benchmark Results**: JSON format for trend analysis
- **Lint Results**: Code quality reports

## Test Best Practices

### Writing Good Tests

1. **Independent Tests**: Each test should be self-contained
2. **Clear Naming**: Use descriptive test names
3. **Proper Setup**: Initialize test data properly
4. **Complete Cleanup**: Clean up after tests
5. **Error Messages**: Provide clear failure messages

### Example of Good Test

```go
func TestUserServiceInsertUser_ValidInput(t *testing.T) {
    // Arrange
    userService := service.GetUserServiceInstance()
    user := &model.User{
        Name:     "testuser",
        Password: "hashed_password",
    }

    // Act
    result, success := userService.InsertUser(user)

    // Assert
    assert.True(t, success, "User insertion should succeed")
    assert.NotNil(t, result, "Result should not be nil")
    assert.Equal(t, "testuser", result.Name, "Username should match")
}
```

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   - Check database is running
   - Verify connection settings in test config
   - Ensure test database exists

2. **Redis Connection Failures**
   - Check Redis is running
   - Verify connection settings
   - Ensure test Redis database exists

3. **Test Flakiness**
   - Check for race conditions
   - Verify proper cleanup between tests
   - Use proper synchronization for concurrent tests

4. **Performance Degradation**
   - Check database connection pool settings
   - Verify Redis connection pool settings
   - Analyze benchmark results

## Maintenance

### Adding New Tests

1. Create test file in appropriate directory
2. Follow naming convention: `*_test.go`
3. Use proper test structure (Arrange, Act, Assert)
4. Add test documentation
5. Update this documentation

### Updating Tests

1. Review test changes with code changes
2. Update test data if needed
3. Update test configuration if needed
4. Update documentation
5. Run all tests to ensure no regressions

## Resources

### Documentation

- [Go Testing](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [Ginkgo](https://onsi.github.io/ginkgo/)

### Tools

- [golang/mock](https://github.com/golang/mock) - Mock generation
- [golangci-lint](https://golangci-lint.run/) - Linting
- [go test](https://golang.org/cmd/go/) - Testing framework

## Contact

For questions or issues related to testing, please refer to the project documentation or contact the development team.

---

Last Updated: 2025-04-15