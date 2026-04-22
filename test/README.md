# ByteDanceDemo Test Suite

This directory contains comprehensive tests for the ByteDanceDemo project.

## Test Organization

```
test/
├── mocks/                    # Mock implementations for external dependencies
│   ├── dao_mock.go          # Database layer mocks
│   ├── redis_mock.go        # Redis client mocks
│   └── cache_mock.go       # Cache layer mocks
├── services/               # Service layer unit tests
│   ├── user_service_test.go         # User service tests
│   ├── comment_service_test.go      # Comment service tests
│   ├── follow_service_test.go       # Follow service tests
│   ├── favorite_service_test.go     # Favorite service tests
│   ├── message_service_test.go      # Message service tests
│   └── video_service_test.go       # Video service tests
├── controllers/            # Controller layer integration tests
│   ├── user_controller_test.go       # User controller tests
│   ├── comment_controller_test.go    # Comment controller tests
│   └── publish_controller_test.go   # Publish controller tests
├── middleware/            # Middleware tests
│   ├── auth_middleware_test.go      # JWT authentication tests
│   └── validation_middleware_test.go # Input validation tests
├── integration/           # End-to-end integration tests
│   └── api_integration_test.go       # API integration tests
├── benchmarks/            # Performance benchmarks
│   ├── database_benchmarks.go       # Database operation benchmarks
│   ├── api_benchmarks.go           # API endpoint benchmarks
│   └── redis_benchmarks.go        # Redis operation benchmarks
├── utils_test.go          # Test utilities and helpers
├── test_config.yaml       # Test configuration file
└── run_tests.sh          # Test runner script
```

## Quick Start

### Running All Tests

```bash
# Using the test runner script
./test/run_tests.sh

# Using Makefile
make test

# Direct Go command
go test -v ./test/...
```

### Running Specific Test Types

```bash
# Unit tests only
./test/run_tests.sh --type unit

# Integration tests only
./test/run_tests.sh --type integration

# Benchmark tests only
./test/run_tests.sh --type benchmark

# With coverage report
./test/run_tests.sh --type unit --coverage
```

## Test Categories

### 1. Unit Tests
- **Location**: `test/services/`
- **Purpose**: Test individual service methods in isolation
- **Dependencies**: Mocked database and Redis
- **Coverage**: ~80% target

**Examples**:
- User service CRUD operations
- Comment service validation and creation
- Follow/unfollow functionality
- Favorite management
- Message operations
- Video publishing and retrieval

### 2. Controller Tests
- **Location**: `test/controllers/`
- **Purpose**: Test HTTP endpoints and request handling
- **Dependencies**: Mocked services
- **Coverage**: ~75% target

**Examples**:
- User registration and login
- Comment creation and deletion
- Video publishing
- Follow/follow/unfollow actions
- Favorite actions

### 3. Middleware Tests
- **Location**: `test/middleware/`
- **Purpose**: Test authentication, authorization, and validation
- **Dependencies**: JWT tokens, validation rules
- **Coverage**: ~85% target

**Examples**:
- JWT token generation and validation
- Request authentication
- Input validation and sanitization
- SQL injection prevention
- XSS protection

### 4. Integration Tests
- **Location**: `test/integration/`
- **Purpose**: Test end-to-end API flows
- **Dependencies**: Real database and Redis (or well-configured mocks)
- **Coverage**: ~70% target

**Examples**:
- Complete user registration and login flow
- Video publishing and commenting flow
- Follow/favorite message flows
- Multi-user interaction scenarios

### 5. Benchmark Tests
- **Location**: `test/benchmarks/`
- **Purpose**: Measure performance of critical operations
- **Dependencies**: Real database and Redis
- **Coverage**: Performance monitoring

**Examples**:
- Database query performance
- Redis cache performance
- API response times
- Concurrent operation handling
- Memory usage analysis

## Test Configuration

### Configuration File

Edit `test/test_config.yaml` to configure test settings:

```yaml
# Database settings
database:
  host: localhost
  port: 3306
  username: test_user
  password: test_password
  database: bytedancedemo_test

# Redis settings
redis:
  host: localhost
  port: 6379
  database: 1

# Test settings
test:
  benchmark_iterations: 1000
  concurrency: 100
  request_timeout: 30s
  debug: true
  cleanup_after_test: true
```

### Environment Variables

```bash
# Database
export TEST_DB_HOST=localhost
export TEST_DB_PORT=3306
export TEST_DB_USER=test_user
export TEST_DB_PASSWORD=test_password
export TEST_DB_NAME=bytedancedemo_test

# Redis
export TEST_REDIS_HOST=localhost
export TEST_REDIS_PORT=6379

# Test settings
export TEST_DEBUG=true
export TEST_CLEANUP=true
```

## Mock Implementations

The project uses comprehensive mocks to enable isolated unit testing:

### Database Mocks

```go
// Example: Mocking UserDAO
mockDAO := &MockUserDAO{}
mockDAO.EXPECT().Create(user).Return(nil)
mockDAO.EXPECT().Where("name", "test").Return(mockDAO)
mockDAO.EXPECT().Find().Return([]interface{}{user}, nil)
```

### Redis Mocks

```go
// Example: Mocking Redis client
mockRedis := NewMockRedisClient()
mockRedis.On("Get", mock.Anything, "user:1").Return("user_data", nil)
mockRedis.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return("OK", nil)
```

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage report
./test/run_tests.sh --type unit --coverage

# Generate HTML coverage report
make test-coverage-html
```

### Coverage Targets

- **Service Layer**: 80% minimum
- **Controller Layer**: 75% minimum
- **Middleware**: 85% minimum
- **Overall Project**: 75% minimum

## Performance Benchmarks

### Running Benchmarks

```bash
# Run all benchmarks
./test/run_tests.sh --type benchmark

# Run specific benchmark
go test -bench=BenchmarkUserServiceInsertUser -benchtime=10s ./test/services/

# Run with memory profiling
go test -bench=. -benchmem ./test/benchmarks/
```

### Performance Targets

- **User Registration**: < 100ms
- **User Login**: < 50ms
- **Video Publishing**: < 500ms
- **Comment Addition**: < 100ms
- **Follow/Unfollow**: < 100ms
- **Message Sending**: < 150ms

## Test Utilities

The `test/utils_test.go` file provides helper functions:

- `createTestUserAndLogin()` - Create test user and get authentication token
- `generateTestUsers()` - Generate test user data
- `generateTestVideos()` - Generate test video data
- `generateTestComments()` - Generate test comment data
- `simulateLoad()` - Simulate load testing scenarios
- `simulateStress()` - Simulate stress testing scenarios

## Continuous Integration

### GitHub Actions Integration

The project includes CI configuration for automated testing:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: make test
      - run: make test-coverage
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   ```bash
   # Check if test database exists
   mysql -u test_user -p -e "CREATE DATABASE IF NOT EXISTS bytedancedemo_test;"
   ```

2. **Redis Connection Errors**
   ```bash
   # Check if Redis is running
   redis-cli ping

   # Use test Redis database
   redis-cli SELECT 1
   ```

3. **Test Flakiness**
   ```bash
   # Run tests with race detector
   ./test/run_tests.sh --race

   # Increase timeout
   go test -timeout=10m ./test/...
   ```

## Best Practices

1. **Write Isolated Tests**: Each test should be self-contained
2. **Use Descriptive Names**: Test names should clearly indicate what they test
3. **Proper Cleanup**: Always clean up test data after tests
4. **Use Table-Driven Tests**: Test multiple scenarios efficiently
5. **Mock External Dependencies**: Don't rely on external services
6. **Test Error Cases**: Don't just test happy paths
7. **Keep Tests Fast**: Avoid slow operations in unit tests

## Contributing

When adding new features:

1. Write tests for new functionality
2. Update existing tests if needed
3. Ensure coverage targets are met
4. Run all tests before committing
5. Update this README if needed

## Additional Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Mock Generation](https://github.com/golang/mock)

## Support

For questions or issues related to testing, please refer to the main project documentation or contact the development team.

---

**Last Updated**: 2025-04-15