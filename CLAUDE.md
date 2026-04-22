# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ByteDanceDemo is a simplified TikTok/Douyin clone backend built with Go (Golang). It's a microservices-oriented application providing video sharing, social interactions, and real-time messaging capabilities.

### Key Architecture Components

- **API Layer**: Gin-based RESTful API with middleware pipeline
- **Service Layer**: Business logic with concurrent operations and caching
- **Data Access Layer**: GORM-generated DAOs for database operations
- **Database**: MySQL for persistent storage, Redis for caching
- **Message Queue**: RabbitMQ for async processing
- **Authentication**: JWT-based auth with Casbin authorization

### Application Entry Point

The application uses Cobra CLI with two main commands:
- `./bin/app server` - Start API server (main entry point)
- `./bin/app migrate` - Run database migrations

Development workflow typically: `make run-api` (uses `cmd/api/service.go`)

## Common Development Commands

### Building
```bash
make build              # Build binary to bin/app
go build -o bin/app ./main.go
```

### Testing
```bash
make test-parallel      # Run all tests with 4 parallel workers
./test/run_tests.sh   # Flexible test runner with options
./test/run_tests.sh --type unit --coverage  # Unit tests with coverage
./test/run_tests.sh --type integration --verbose  # Integration tests
./test/run_tests.sh --type benchmark  # Performance benchmarks
```

### Running Services
```bash
make run-api           # Start API server only
make run-migrate        # Run migrations only
make run-parallel       # Start migrate and API together
./bin/app server -c config/settings.yml -m debug  # With custom config and mode
```

### Environment Setup
```bash
# Set China proxy (required for dependency downloads)
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=off

# Install dependencies
make deps  # or go mod download && go mod tidy
```

### Development Workflow
```bash
make ci    # Full CI workflow: clean -> deps -> test -> build
make dev   # Development workflow: test -> build
```

## Code Architecture

### Layer Structure

1. **Controller Layer** (`controller/`)
   - HTTP request handlers
   - Input validation
   - Response formatting
   - JWT token extraction

2. **Service Layer** (`service/`)
   - Business logic implementation
   - Concurrency management (goroutines)
   - Caching integration
   - Cross-service orchestration

3. **Data Access Layer** (`dao/`, `database/`)
   - GORM-generated query methods
   - Database connection management
   - Transaction handling

4. **Middleware** (`middleware/`)
   - Security headers
   - CORS handling
   - Rate limiting
   - JWT authentication
   - CSRF protection
   - Input validation
   - Casbin authorization

### Concurrency Pattern

Services use a pattern of creating 5-6 goroutines per operation for parallel execution:
- Video feed operations fetch author, comments, likes concurrently
- User profile operations fetch follow counts, metrics in parallel
- Use `sync.WaitGroup` for coordination
- **Note**: Uncontrolled goroutine creation is a performance concern - consider worker pools

### Caching Strategy

- Redis caching with variable TTL based on data size
- Cache keys: `user_followings:{id}`, `user_followers:{id}`, etc.
- Expiration varies from 5-60 seconds based on result count
- Cache warming on first access, then lazy refresh

## Database Schema

### Core Tables
- `users` - User accounts and profiles
- `videos` - Video metadata and content
- `comments` - Comment threads
- `likes` - Like/favorite records
- `relations` - Follow/following relationships
- `messages` - Chat messages

### DAO Generation

GORM gen is used to generate type-safe database access code:
- Generated files in `dao/` directory (comments.gen.go, users.gen.go, etc.)
- Run `go run model/main/mian.go` to regenerate after schema changes

## Configuration Management

### Configuration Files
- Primary: `config/settings.yml` (created from template for development)
- Template: `config/settings.yml.template`
- Environment-specific: `config/settings.{env}.yml` for overrides

### Key Configuration Sections
- `settings.application` - Rate limiting, timeouts, upload limits
- `settings.mysql` - Database connection parameters
- `settings.redis` - Cache configuration
- `settings.jwt` - Authentication token settings
- `settings.rabbitmq` - Message queue configuration
- `settings.log` - Logging levels and rotation

## API Endpoints

### Basic APIs
- `GET /douyin/feed/` - Video feed with pagination
- `GET /douyin/user/` - User profile info
- `POST /douyin/user/register/` - User registration
- `POST /douyin/user/login/` - User authentication
- `POST /douyin/publish/action/` - Video upload
- `GET /douyin/publish/list/` - User's published videos

### Extra APIs
- `POST /douyin/favorite/action/` - Like/unlike video
- `GET /douyin/favorite/list/` - User's favorited videos
- `POST /douyin/comment/action/` - Comment on video
- `GET /douyin/comment/list/` - Video comments
- `POST /douyin/relation/action/` - Follow/unfollow
- `GET /douyin/relation/follow/list/` - User's following
- `GET /douyin/relation/follower/list/` - User's followers
- `GET /douyin/relation/friend/list/` - User's friends (mutual follow)
- `GET /douyin/message/chat/` - Get chat messages
- `POST /douyin/message/action/` - Send message

Static files served from: `/static/*` (located in `public/` directory)

## Testing Strategy

### Test Organization
- `test/services/` - Service layer unit tests
- `test/controllers/` - Controller layer tests
- `test/integration/` - End-to-end API tests
- `test/middleware/` - Middleware unit tests
- `test/benchmarks/` - Performance benchmarks
- `test/mocks/` - Generated mock dependencies

### Test Execution
Tests run in parallel with 4 workers by default. Integration tests require the API server to be running on port 8080.

### Known Test Issues
- Integration tests fail if API server not running
- Some service tests require Redis connection
- Test server address configurable in `test/common.go` via `_serverAddr_`

## Security Considerations

### Current State (from security review)
- **Critical**: SQL injection risks in some query constructions
- **Critical**: Weak JWT secret (`123456`) in development config
- **High**: Input validation missing on some endpoints
- **High**: MD5 password hashing (should use bcrypt)
- **Medium**: No CSRF protection (middleware exists but may need review)

### Security Middleware Stack
- SecurityHeadersMiddleware - Adds security headers
- CORSMiddleware - CORS handling
- JWTMiddleware - Token validation
- CSRFMiddleware - CSRF token generation/validation
- CasbinMiddleware - RBAC authorization

## Performance Considerations

### Known Bottlenecks (from performance analysis)
- **N+1 Queries**: Video feed makes individual queries per video for author, comments, likes
- **Goroutine Control**: Uncontrolled goroutine creation in service layer
- **Connection Pooling**: MySQL and Redis lack proper pool configuration
- **Caching**: Inconsistent TTL strategies across services

### Recommendations
- Implement batch queries for feed operations
- Use worker pools instead of individual goroutines
- Configure connection pools for database and Redis
- Implement consistent caching strategy

## Documentation

### Available Documentation
- `docs/ARCHITECTURE.md` - System architecture overview
- `docs/CONFIGURATION.md` - Configuration reference
- `docs/DEPLOYMENT.md` - Deployment procedures
- `docs/CONTRIBUTING.md` - Contribution guidelines
- `docs/SECURITY.md` - Security best practices
- `docs/PERFORMANCE.md` - Performance tuning guide
- `docs/TROUBLESHOOTING.md` - Common issues and solutions
- `docs/api/openapi-spec.yaml` - OpenAPI specification

### Quick Documentation Access
```bash
cat docs/README.md  # Documentation index with quick links
```

## Development Tips

### Adding New Features
1. Add DAO methods in `dao/gen.go` for database operations
2. Implement business logic in appropriate `service/*` file
3. Create controller handler in `controller/` directory
4. Register route in `router/router.go`
5. Add tests in `test/` subdirectories

### Database Changes
- Modify schema using SQL scripts in `config/init.sql`
- Regenerate DAOs: `go run model/main/mian.go`
- Update model files as needed

### Adding Tests
- Create mock in `test/mocks/` using `gomock`
- Write test files in appropriate `test/` subdirectory
- Use `test/middleware_mock.go` for middleware helpers
- Run `./test/run_tests.sh --type unit` to verify

### Debugging
- Application runs on port 8080 by default
- Logs configured via `settings.log` section
- Use `--mode debug` flag for verbose logging
- Static video files accessible at `http://localhost:8080/static/`

## Important Notes

### User Authentication
- JWT tokens stored in memory (`usersLoginInfo` map) for development
- Tokens expire based on `settings.jwt.expirationTime`
- Production should use persistent storage (Redis/database)

### Video Storage
- Uploaded videos saved to `public/` directory
- Access via: `http://localhost:8080/static/{filename}`
- File upload size limited by `settings.application.maxUploadSize`

### Message Queue
- RabbitMQ used for async processing of comments, follows, messages
- Multiple exchanges: comment, follow, message
- Connection configured via `settings.rabbitmq` section

### Rate Limiting
- Default: 50 requests per minute
- Configured via `settings.application.rateLimit`
- Redis-backed for distributed systems

## Module Dependencies

### Key External Libraries
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/golang-jwt/jwt` - JWT implementation
- `github.com/rabbitmq/amqp091-go` - RabbitMQ client
- `github.com/casbin/casbin/v2` - Authorization framework
- `github.com/spf13/viper` - Configuration management
- `github.com/spf13/cobra` - CLI framework
- `go.uber.org/zap` - Structured logging

### Go Version
- Target: Go 1.20
- Minimum: Go 1.18 (for compatibility)
