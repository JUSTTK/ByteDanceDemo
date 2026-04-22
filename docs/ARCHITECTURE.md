# ByteDanceDemo - Project Architecture

## Overview

ByteDanceDemo is a simplified TikTok/Douyin clone backend built with Go (Golang), featuring a microservices-oriented architecture with clean separation of concerns. The system provides video sharing, social interactions, and real-time messaging capabilities.

## Technology Stack

### Core Framework
- **Gin Framework**: HTTP web framework for building RESTful APIs
- **Go 1.20**: Core programming language
- **GORM**: ORM library for database operations
- **MySQL**: Primary relational database
- **Redis**: Caching and session management

### Authentication & Security
- **JWT (golang-jwt/jwt)**: Token-based authentication
- **bcrypt**: Password hashing
- **Casbin**: Authorization and access control
- **Sensitive**: Content filtering and moderation

### Message Queue
- **RabbitMQ**: Asynchronous message processing for real-time features

### Logging & Monitoring
- **Zap**: High-performance logging library
- **Lumberjack**: Log rotation and archival

### Testing
- **Testify**: Testing framework
- **Gomock**: Mock generation for testing
- **Httpexpect**: HTTP testing library

### Configuration
- **Viper**: Configuration management
- **Cobra**: Command-line interface

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│                   (Mobile/Web Applications)                   │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ HTTP/HTTPS
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                      API Gateway                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Middleware Layer                        │   │
│  │  • Security Headers                                  │   │
│  │  • CORS Configuration                                │   │
│  │  • Rate Limiting                                     │   │
│  │  • Request Logging                                   │   │
│  │  • Error Handling                                    │   │
│  │  • Input Validation                                  │   │
│  │  • CSRF Protection                                   │   │
│  │  • JWT Authentication                                │   │
│  │  • Casbin Authorization                              │   │
│  └──────────────────────────────────────────────────────┘   │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   Controller  │ │   Controller  │ │   Controller  │
│     Layer     │ │     Layer     │ │     Layer     │
└───────┬───────┘ └───────┬───────┘ └───────┬───────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                      Service Layer                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  User    │  │  Video   │  │  Social  │  │ Message  │   │
│  │ Service  │  │ Service  │  │ Service  │  │ Service  │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   Repository  │ │   Repository  │ │   Repository  │
│     Layer     │ │     Layer     │ │     Layer     │
└───────┬───────┘ └───────┬───────┘ └───────┬───────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Data Access Layer                        │
│  ┌──────────────┐      ┌──────────────┐                   │
│  │   MySQL DB   │      │    Redis     │                   │
│  │              │      │    Cache     │                   │
│  │  • Users     │      │  • Sessions  │                   │
│  │  • Videos    │      │  • Rate Limit│                   │
│  │  • Likes     │      │  • Cache     │                   │
│  │  • Comments  │      │              │                   │
│  │  • Relations │      │              │                   │
│  │  • Messages  │      │              │                   │
│  └──────────────┘      └──────────────┘                   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Message Queue (RabbitMQ)                │   │
│  │  • Real-time message processing                      │   │
│  │  • Asynchronous task handling                       │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
ByteDanceDemo/
├── cmd/                    # Command-line interface
│   ├── api/               # API server command
│   │   └── service.go     # Server initialization
│   ├── migrate/           # Database migration commands
│   └── cobra.go           # Root command configuration
├── config/                # Configuration files
│   ├── config.go          # Configuration loading logic
│   ├── settings.yml       # Application settings
│   ├── init.sql           # Database schema
│   └── rbac_model.conf    # Casbin RBAC model
├── controller/            # HTTP request handlers
│   ├── user.go            # User authentication & profile
│   ├── video.go           # Video upload & listing
│   ├── publish.go         # Video publishing
│   ├── feed.go            # Video feed
│   ├── favorite.go        # Like/unlike videos
│   ├── comment.go         # Comments management
│   ├── relation.go        # Follow/unfollow users
│   ├── message.go         # Messaging system
│   └── common.go          # Shared data structures
├── service/               # Business logic layer
│   ├── userService.go     # User service interface
│   ├── userServiceImpl.go # User service implementation
│   ├── VideoService.go    # Video service interface
│   ├── VideoServiceImpl.go # Video service implementation
│   ├── favorite_service.go # Favorite service interface
│   ├── favorite_action.go # Favorite action implementation
│   ├── favorite_list.go  # Favorite list implementation
│   ├── commentService.go  # Comment service interface
│   ├── commentServiceImpl.go # Comment service implementation
│   ├── followService.go  # Follow service interface
│   ├── followServiceImpl.go # Follow service implementation
│   └── messageService.go # Message service implementation
├── repository/            # Data access layer
│   ├── userRepository.go  # User data access
│   ├── videoRepository.go # Video data access
│   └── ...               # Other repositories
├── model/                 # Database models (GORM generated)
│   ├── users.gen.go       # User model
│   ├── videos.gen.go      # Video model
│   ├── comments.gen.go    # Comment model
│   ├── likes.gen.go       # Like model
│   ├── relations.gen.go  # Relation model
│   ├── messages.gen.go    # Message model
│   └── casbin_rule.gen.go # Casbin rule model
├── dao/                   # Data access objects
│   └── user.go            # User DAO
├── middleware/            # HTTP middleware
│   ├── auth.go            # Authentication middleware
│   ├── cors.go            # CORS middleware
│   ├── rate_limit.go      # Rate limiting middleware
│   ├── logger.go          # Logging middleware
│   ├── error.go           # Error handling middleware
│   ├── validation.go      # Input validation middleware
│   ├── csrf.go            # CSRF protection
│   ├── casbin.go          # Authorization middleware
│   └── rabbitmq/          # RabbitMQ middleware
├── router/                # Route definitions
│   └── router.go          # API routes setup
├── database/              # Database initialization
│   ├── mysql/             # MySQL connection
│   └── redis/             # Redis connection
├── utils/                 # Utility functions
│   ├── log/               # Logging utilities
│   ├── token/             # JWT token utilities
│   ├── encryption/        # Password encryption
│   └── ...               # Other utilities
├── test/                  # Test files
│   ├── integration/       # Integration tests
│   ├── unit/              # Unit tests
│   └── benchmark/         # Benchmark tests
├── public/                # Static assets
│   └── video/             # Uploaded videos
├── docs/                  # Documentation
│   ├── api/               # API documentation
│   ├── architecture.md   # Architecture documentation
│   ├── deployment.md     # Deployment guide
│   ├── configuration.md  # Configuration reference
│   ├── contributing.md   # Contributing guidelines
│   ├── security.md       # Security best practices
│   ├── troubleshooting.md # Troubleshooting guide
│   └── performance.md    # Performance tuning
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go dependencies checksum
└── Makefile              # Build automation
```

## Core Components

### 1. Controller Layer
**Responsibility**: Handle HTTP requests and responses
- Parse request parameters
- Validate input data
- Call appropriate service methods
- Format and return responses
- Handle HTTP-specific concerns

**Key Controllers**:
- `user.go`: Authentication, registration, user profile
- `publish.go`: Video upload and publishing
- `feed.go`: Video feed generation
- `favorite.go`: Like/unlike functionality
- `comment.go`: Comment management
- `relation.go`: Follow/unfollow functionality
- `message.go`: Real-time messaging

### 2. Service Layer
**Responsibility**: Implement business logic
- Orchestrate business operations
- Coordinate between multiple repositories
- Implement transaction management
- Apply business rules and validation
- Handle complex business scenarios

**Key Services**:
- `UserService`: User management, authentication
- `VideoService`: Video CRUD operations
- `FavoriteService`: Like/unlike operations
- `CommentService`: Comment management
- `FollowService`: Follow relationships
- `MessageService`: Message handling

### 3. Repository Layer
**Responsibility**: Data access and persistence
- Abstract database operations
- Provide CRUD operations
- Handle complex queries
- Manage caching logic
- Optimize database interactions

### 4. Middleware Layer
**Responsibility**: Cross-cutting concerns
- Security headers injection
- CORS handling
- Rate limiting
- Request/response logging
- Error handling and recovery
- Input validation
- CSRF protection
- JWT authentication
- Casbin authorization

## Data Flow

### Authentication Flow
```
1. Client sends credentials to /douyin/user/login/
2. Controller receives request
3. Service validates credentials
4. Repository queries database
5. If valid, Service generates JWT token
6. Controller returns token to client
7. Client includes token in subsequent requests
8. Middleware validates token
9. Request proceeds to protected endpoints
```

### Video Upload Flow
```
1. Client sends video with token to /douyin/publish/action/
2. Middleware validates JWT and CSRF
3. Controller receives multipart form data
4. Service processes video (validation, storage)
5. Repository saves metadata to database
6. Video file stored in /public/video/
7. Service triggers async processing via RabbitMQ
8. Controller returns success response
```

### Message Flow
```
1. Client sends message via /douyin/message/action/
2. Middleware validates authentication
3. Service validates message content
4. Message published to RabbitMQ queue
5. Consumer processes message asynchronously
6. Repository persists message to database
7. Real-time push to recipient (if online)
8. Controller confirms message sent
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(191) NOT NULL,
    password VARCHAR(191) NOT NULL,
    role VARCHAR(191) NOT NULL,
    avatar VARCHAR(191) DEFAULT 'http://yourserver.com/default_avatar.jpg',
    background_image VARCHAR(191) DEFAULT 'http://yourserver.com/default_background.jpg',
    signature TEXT,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3)
);
```

### Videos Table
```sql
CREATE TABLE videos (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    author_id BIGINT UNSIGNED NOT NULL,
    play_url VARCHAR(191) NOT NULL,
    cover_url VARCHAR(191),
    title VARCHAR(191),
    favorite_count BIGINT DEFAULT 0,
    comment_count BIGINT DEFAULT 0,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    FOREIGN KEY (author_id) REFERENCES users(id)
);
```

### Comments Table
```sql
CREATE TABLE comments (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    FOREIGN KEY (video_id) REFERENCES videos(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Likes Table
```sql
CREATE TABLE likes (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    video_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3),
    UNIQUE KEY unique_like (video_id, user_id),
    FOREIGN KEY (video_id) REFERENCES videos(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Relations Table
```sql
CREATE TABLE relations (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    follower_id BIGINT UNSIGNED NOT NULL,
    followee_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3),
    UNIQUE KEY unique_relation (follower_id, followee_id),
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followee_id) REFERENCES users(id)
);
```

### Messages Table
```sql
CREATE TABLE messages (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    from_user_id BIGINT UNSIGNED NOT NULL,
    to_user_id BIGINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME(3),
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
);
```

## Caching Strategy

### Redis Cache Usage
1. **User Sessions**: Store active user sessions with TTL
2. **Rate Limiting**: Track request counts per user/IP
3. **Feed Cache**: Cache video feed to reduce database load
4. **User Info**: Cache frequently accessed user profiles
5. **Like Status**: Cache like status for current user

### Cache Invalidation
- Time-based expiration (TTL)
- Manual invalidation on data updates
- Write-through pattern for critical data
- Cache-aside pattern for read-heavy operations

## Security Architecture

### Authentication
- JWT-based token authentication
- bcrypt password hashing
- Token expiration management
- Refresh token mechanism (future enhancement)

### Authorization
- Casbin RBAC model
- Role-based access control
- Resource-level permissions
- API-level authorization

### Security Measures
- SQL injection prevention (parameterized queries)
- XSS protection (input sanitization)
- CSRF token validation
- Rate limiting (prevent abuse)
- Security headers (CSP, X-Frame-Options, etc.)
- Content filtering (sensitive word detection)

## Scalability Considerations

### Horizontal Scaling
- Stateless application design
- Database connection pooling
- Redis cluster support
- Load balancer ready

### Performance Optimization
- Database indexing on frequently queried fields
- Query optimization and N+1 prevention
- Asynchronous processing (RabbitMQ)
- Caching strategy (Redis)
- Connection pooling

### Monitoring & Observability
- Structured logging with Zap
- Request/response logging
- Error tracking
- Performance metrics (future enhancement)
- Distributed tracing (future enhancement)

## Deployment Architecture

### Development Environment
- Single server deployment
- Local MySQL and Redis
- In-memory RabbitMQ (for testing)

### Production Environment
- Load balancer (Nginx/HAProxy)
- Multiple application instances
- Master-slave MySQL replication
- Redis cluster with persistence
- RabbitMQ cluster
- CDN for static assets
- Object storage for videos (OSS)

## Future Enhancements

### Planned Features
- WebSocket support for real-time updates
- Video processing pipeline (transcoding, thumbnail generation)
- Recommendation engine
- Analytics and reporting
- Admin dashboard
- Push notifications
- Live streaming support

### Technical Improvements
- GraphQL API alternative
- Microservices decomposition
- Event-driven architecture
- Circuit breakers and resilience patterns
- Advanced monitoring and alerting
- CI/CD pipeline enhancement
