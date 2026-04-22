# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation system improvements
- Enhanced README.md with project overview, features, and quick start guide
- Added API usage examples (curl, JavaScript, Go)
- Created FAQ documentation
- Added health check script for system monitoring
- Implemented GitHub Actions CI/CD pipeline
- Added Docker support with Dockerfile and docker-compose.yml
- Created developer quick start guide
- Unified terminology across project
- Added documentation update dates to all files

### Changed
- Restructured test documentation to reduce duplication
- Consolidated configuration documentation
- Improved code comments and documentation coverage
- Enhanced test utilities with better documentation

### Security
- Fixed SQL injection vulnerabilities in multiple controllers
- Upgraded password hashing from MD5 to bcrypt
- Strengthened JWT secret key management
- Added input validation middleware
- Implemented CSRF protection
- Enhanced security headers

### Fixed
- Fixed duplicate main function in model/main/
- Resolved go vet warnings in controller/feed.go
- Fixed database package imports
- Corrected field naming issues (user.ID vs user.Id)
- Fixed string conversion issues

### Performance
- Identified and documented N+1 query bottlenecks
- Added batch query optimization recommendations
- Documented goroutine pooling strategy
- Added Redis connection pooling configuration
- Documented caching strategy improvements

## [1.1.0] - 2026-04-20

### Added
- Complete testing suite with unit, integration, and benchmark tests
- Mock implementations for database, Redis, and cache
- Test utilities and helper functions
- Test configuration management
- Automated test runner script with multiple options

### Security
- SQL injection prevention implementation
- MD5 password hashing upgraded to bcrypt
- JWT secret key generation and security checks
- Input validation enhancement
- CSRF protection middleware
- Security headers middleware

### Documentation
- Architecture documentation with system design
- Configuration reference with all settings
- Deployment guide for development, Docker, and Kubernetes
- Contributing guidelines with code style and workflow
- Security best practices guide
- Performance tuning guide
- Troubleshooting guide with common issues
- OpenAPI specification for all endpoints

## [1.0.0] - 2026-04-15

### Added
- Initial project release
- User authentication and authorization (JWT + Casbin RBAC)
- Video upload and sharing functionality
- Comment system
- Like/favorite system
- Follow/unfollow system
- Real-time messaging
- Security middleware (JWT, CSRF, rate limiting)
- Database models (users, videos, comments, likes, relations, messages)
- Service layer with business logic
- Repository layer for data access
- Controller layer for HTTP endpoints

### Security
- JWT token-based authentication
- bcrypt password hashing
- Casbin authorization model
- SQL injection prevention
- XSS protection
- CSRF protection
- Rate limiting

### Technology Stack
- Go 1.20
- Gin web framework
- GORM ORM
- MySQL 8.0
- Redis 6.0
- RabbitMQ 3.9
- Zap logging
