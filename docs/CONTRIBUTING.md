# ByteDanceDemo - Contributing Guidelines

Thank you for your interest in contributing to ByteDanceDemo! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Code Style](#code-style)
5. [Testing Guidelines](#testing-guidelines)
6. [Documentation Guidelines](#documentation-guidelines)
7. [Pull Request Process](#pull-request-process)
8. [Issue Reporting](#issue-reporting)
9. [Release Process](#release-process)
10. [Community Guidelines](#community-guidelines)

## Code of Conduct

This project follows a standard Code of Conduct. Please be respectful and inclusive when participating in this project's development.

### Our Pledge

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone, regardless of age, body size, visible or invisible disability, ethnicity, sex characteristics, gender identity and expression, level of experience, education, socio-economic status, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

Examples of behavior that contributes to a positive environment for our community:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

Examples of unacceptable behavior by participants:

- The use of sexualized language or imagery and unwelcome sexual attention or advances
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information, such as a physical or electronic address, without explicit permission
- Other conduct which could reasonably be considered inappropriate in a professional setting

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git
- MySQL 8.0 or higher
- Redis 6.0 or higher
- RabbitMQ 3.9 or higher

### Setup Development Environment

1. **Fork the Repository**
   ```bash
   # Fork the repository on GitHub
   git clone https://github.com/your-username/ByteDanceDemo.git
   cd ByteDanceDemo
   git remote add upstream https://github.com/original-owner/ByteDanceDemo.git
   ```

2. **Install Dependencies**
   ```bash
   # Download Go modules
   go mod download
   go mod tidy

   # Install development tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **Setup Database**
   ```bash
   # Create development database
   mysql -u root -p -e "CREATE DATABASE sample_douyin_dev;"
   mysql -u root -p sample_douyin_dev < config/init.sql
   ```

4. **Create Configuration File**
   ```bash
   cp config/settings.yml.template config/settings.yml
   # Edit settings.yml for development
   ```

5. **Build and Test**
   ```bash
   # Build the application
   make build

   # Run tests
   make test

   # Run the application
   make run
   ```

### Development Tools

We recommend using the following tools for development:

- **IDE**: Visual Studio Code with Go extension
- **Go Plugins**:
  - Go (by Go team)
  - Go Doc (by Go team)
  - Go Test Explorer (by utku)
  - GitLens (by Microsoft)
- **Docker**: For containerized development
- **Postman**: For API testing
- **GitKraken**: For Git history visualization

## Development Workflow

### Branch Strategy

We use the GitFlow branching model:

```
main
├── develop
├── feature/feature-name
├── bugfix/bug-description
├── release/v1.0.0
└── hotfix/critical-fix
```

### Branch Naming Convention

- **Feature branches**: `feature/description` (e.g., `feature/user-profile`)
- **Bugfix branches**: `bugfix/description` (e.g., `bugfix/login-error`)
- **Release branches**: `release/vX.Y.Z` (e.g., `release/v1.0.0`)
- **Hotfix branches**: `hotfix/description` (e.g., `hotfix/security-fix`)

### Workflow Steps

1. **Create a new branch from develop**
   ```bash
   git checkout develop
   git pull upstream develop
   git checkout -b feature/your-feature-name
   ```

2. **Make changes**
   - Make your changes
   - Follow the [Code Style](#code-style) guidelines
   - Write tests for new functionality
   - Update documentation if needed

3. **Commit changes**
   ```bash
   git add .
   git commit -m "feat: add user profile feature"
   git commit -m "fix: resolve login error"
   git commit -m "docs: update API documentation"
   ```

4. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create pull request on GitHub
   ```

### Commit Message Convention

We follow the Conventional Commits specification:

```
<type>(<scope>): <description>

[body]

<footer>
```

#### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or modifying tests
- `chore`: Build process or auxiliary tool changes

#### Examples

```
feat(auth): add JWT token refresh mechanism

- Add refresh token endpoint
- Implement token rotation
- Update authentication middleware

Closes #123
```

```
fix(video): resolve video upload size validation

- Properly validate file size before upload
- Return appropriate error message
- Add unit test for validation

Fixes #456
```

## Code Style

### Go Code Style

We use the standard Go formatting conventions with some additions:

```go
// Package declaration
package main

// Import statements grouped by category
import (
    "context"
    "fmt"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    "bytedancedemo/model"
    "bytedancedemo/service"
)

// Type definitions
type UserService struct {
    db *gorm.DB
    cache *redis.Client
}

// Struct tags
type User struct {
    ID        int64  `json:"id" gorm:"primaryKey"`
    Name      string `json:"name" gorm:"not null;unique"`
    Email     string `json:"email" gorm:"unique"`
    Password  string `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Function documentation
// GetUserByID retrieves a user by ID from the database
// Returns the user and any error encountered
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    result := s.db.WithContext(ctx).First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

### Naming Conventions

- **Variables**: `camelCase` (e.g., `userName`, `videoId`)
- **Constants**: `SCREAMING_SNAKE_CASE` (e.g., `MAX_UPLOAD_SIZE`)
- **Functions**: `camelCase` (e.g., `GetUser`, `CreateVideo`)
- **Interfaces**: `camelCase` with `-er` suffix if applicable (e.g., `UserService`)
- **Packages**: lowercase, short, and descriptive (e.g., `user`, `video`)

### Code Organization

```
service/
├── user/
│   ├── user.go          # Interface definitions
│   ├── user_impl.go     # Implementation
│   └── user_test.go     # Tests
├── video/
│   ├── video.go
│   ├── video_impl.go
│   └── video_test.go
└── common.go            # Shared utilities
```

### Best Practices

1. **Error Handling**
   ```go
   // Good
   if err != nil {
       return nil, fmt.Errorf("failed to get user: %w", err)
   }
   
   // Bad
   if err != nil {
       return nil, err
   }
   ```

2. **Context Usage**
   ```go
   // Good
   func (s *Service) ProcessRequest(ctx context.Context, req Request) error {
       // Pass context to database calls
       user, err := s.repo.GetUser(ctx, req.UserID)
       if err != nil {
           return err
       }
       return nil
   }
   ```

3. **Interface Segregation**
   ```go
   // Good - separated interfaces
   type Repository interface {
       GetByID(id int64) (*model.User, error)
       Create(user *model.User) error
       Update(user *model.User) error
       Delete(id int64) error
   }
   
   type Cache interface {
       Get(key string) (string, error)
       Set(key, value string, ttl time.Duration) error
       Delete(key string) error
   }
   ```

4. **Dependency Injection**
   ```go
   // Good
   func NewUserService(db *gorm.DB, cache Cache) *UserService {
       return &UserService{
           db:    db,
           cache: cache,
       }
   }
   ```

### Code Formatting

Use `gofmt` and `goimports` for formatting:

```bash
# Format code
gofmt -w .

# Import formatting
goimports -w .

# Run all formatters
make format
```

### Linting

We use golangci-lint for static analysis:

```bash
# Run linter
golangci-lint run

# Run with specific issues
golangci-lint run --enable=errcheck --enable=goconst
```

## Testing Guidelines

### Testing Structure

```
service/
├── user/
│   ├── user.go
│   ├── user_impl.go
│   ├── user_test.go       # Unit tests
│   └── user_integration.go # Integration tests
```

### Unit Tests

```go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock implementation
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int64) (*model.User, error) {
    args := m.Called(id)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestUserService_GetUser_Success(t *testing.T) {
    // Setup
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    expectedUser := &model.User{
        ID:   1,
        Name: "testuser",
    }
    
    mockRepo.On("GetByID", int64(1)).Return(expectedUser, nil)
    
    // Execute
    user, err := service.GetUser(1)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
    // Setup
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    
    mockRepo.On("GetByID", int64(1)).Return(nil, gorm.ErrRecordNotFound)
    
    // Execute
    user, err := service.GetUser(1)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, user)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests

```go
package user

import (
    "testing"
    "bytedancedemo/database"
    "github.com/stretchr/testify/suite"
)

type UserServiceIntegrationTestSuite struct {
    suite.Suite
    userService *Service
    db         *gorm.DB
}

func (suite *UserServiceIntegrationTestSuite) SetupSuite() {
    // Setup test database
    db, err := database.InitTestDB()
    suite.NoError(err)
    suite.db = db
    
    // Setup service
    suite.userService = NewUserService(db)
}

func (suite *UserServiceIntegrationTestSuite) TearDownSuite() {
    // Cleanup test database
    suite.db.Exec("DROP TABLE users")
    suite.db.Close()
}

func (suite *UserServiceIntegrationTestSuite) SetupTest() {
    // Clean up before each test
    suite.db.Exec("DELETE FROM users")
}

func (suite *UserServiceIntegrationTestSuite) TestCreateUser_Success() {
    // Arrange
    user := &model.User{
        Name:  "testuser",
        Email: "test@example.com",
        Password: "hashedpassword",
    }
    
    // Act
    err := suite.userService.Create(user)
    
    // Assert
    suite.NoError(err)
    suite.True(user.ID > 0)
}

func TestUserServiceIntegrationSuite(t *testing.T) {
    suite.Run(t, new(UserServiceIntegrationTestSuite))
}
```

### Benchmark Tests

```go
package user

import (
    "testing"
    "bytedancedemo/database"
)

func BenchmarkGetUser(b *testing.B) {
    // Setup
    db, _ := database.InitTestDB()
    defer db.Close()
    
    userService := NewUserService(db)
    
    // Create test user
    user := &model.User{
        Name:  "benchmarkuser",
        Email: "benchmark@example.com",
        Password: "hashedpassword",
    }
    userService.Create(user)
    
    // Benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = userService.GetUser(user.ID)
    }
}
```

### Test Coverage

We maintain high test coverage standards:

```bash
# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
go test -covermode=count -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "total:" | awk '{print substr($3,1, length($3)-1)}'
```

### Test Data Management

Use test data factories:

```go
package testdata

import "time"

func User(id int64) *model.User {
    return &model.User{
        ID:        id,
        Name:      fmt.Sprintf("user%d", id),
        Email:     fmt.Sprintf("user%d@example.com", id),
        Password:  "hashedpassword",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func Video(id int64, authorID int64) *model.Video {
    return &model.Video{
        ID:        id,
        AuthorID:  authorID,
        PlayURL:   fmt.Sprintf("/videos/video%d.mp4", id),
        CoverURL:  fmt.Sprintf("/videos/cover%d.jpg", id),
        CreatedAt: time.Now(),
    }
}
```

## Documentation Guidelines

### Code Documentation

We use Go's standard documentation format:

```go
// Package user provides user management services for the ByteDanceDemo application.
// It includes user authentication, profile management, and related operations.
package user

import (
    "context"
    "gorm.io/gorm"
)

// UserService provides operations for managing users.
type UserService struct {
    db    *gorm.DB
    cache Cache
}

// NewUserService creates a new user service with the given dependencies.
// Parameters:
//   - db: Database connection
//   - cache: Cache client
// Returns:
//   - *UserService: New user service instance
func NewUserService(db *gorm.DB, cache Cache) *UserService {
    return &UserService{
        db:    db,
        cache: cache,
    }
}

// GetByID retrieves a user by their ID.
// This function caches user data in Redis for performance.
// Parameters:
//   - ctx: Request context for cancellation
//   - id: User ID to retrieve
// Returns:
//   - *User: User data
//   - error: Error if user not found or database error
func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    // First try to get from cache
    cacheKey := fmt.Sprintf("user:%d", id)
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        return unmarshalUser(cached)
    }
    
    // If not in cache, get from database
    var user User
    result := s.db.WithContext(ctx).First(&user, id)
    if result.Error != nil {
        return nil, fmt.Errorf("user not found: %w", result.Error)
    }
    
    // Cache the user data
    go s.cache.Set(ctx, cacheKey, marshalUser(&user), 5*time.Minute)
    
    return &user, nil
}
```

### API Documentation

Generate OpenAPI/Swagger documentation:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g main.go -o docs/api

# Update after changes
swag init -g main.go -o docs/api
```

### README Updates

When adding new features or breaking changes:

1. Update the main README.md with feature descriptions
2. Add relevant examples
3. Update installation instructions if needed
4. Add troubleshooting section if necessary

## Pull Request Process

### PR Checklist

Before creating a PR, ensure:

- [ ] Code follows the [Code Style](#code-style) guidelines
- [ ] All tests pass (`make test`)
- [ ] Code is linted (`golangci-lint run`)
- [ ] Documentation is updated if needed
- [ ] Change is tested (unit/integration tests)
- [ ] PR title follows [convention](#commit-message-convention)
- [ ] PR description clearly explains changes
- [ ] Relevant issues are linked
- [ ] Breaking changes are clearly documented

### PR Template

Create PR with this template:

```markdown
## Description
Brief description of the changes made.

## Changes Made
- Feature 1 description
- Bug fix 2 description
- Documentation update

## Testing
- Test case 1
- Test case 2
- Integration test results

## Breaking Changes
- If any, list them here

## Related Issues
Closes #123, #456

## Screenshots (if applicable)
![Screenshot of the feature]

## Additional Notes
Any additional information reviewers should know
```

### Review Process

1. **PR Submission**
   - Create PR from feature branch to `develop`
   - Fill in PR template
   - Link to relevant issues

2. **Automated Checks**
   - CI pipeline runs tests and linting
   - Build verification
   - Security scan

3. **Code Review**
   - At least one maintainer must approve
   - Address all review comments
   - Keep PR focused on single feature

4. **Merging**
   - Once approved, PR is merged to `develop`
   - Fast-forward merge without history rewriting
   - No force pushes after review

### PR Guidelines

- Keep PRs small and focused
- Respond to review comments promptly
- Update PR title/description based on feedback
- Don't merge until all issues are resolved
- Maintain clean commit history

## Issue Reporting

### Bug Reports

Use GitHub issues to report bugs:

```markdown
## Bug Description
Brief description of the bug.

## Steps to Reproduce
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior
What should happen.

## Actual Behavior
What happened instead.

## Environment
- Go version: 1.20
- OS: Ubuntu 20.04
- MySQL version: 8.0.26
- Redis version: 6.2.6

## Screenshots
![Screenshot of the bug]

## Additional Context
Any other context about the problem.
```

### Feature Requests

For feature requests:

```markdown
## Feature Description
Clear and concise description of the feature.

## Problem Statement
What problem does this feature solve?

## Proposed Solution
Detailed explanation of the feature implementation.

## Alternatives Considered
What other approaches did you consider?

## Additional Context
Any additional context or screenshots.
```

### Issue Labels

Use standard labels:
- `bug`: Bug report
- `enhancement`: Feature request
- `documentation`: Documentation issue
- `question`: Question or discussion
- `good first issue`: Suitable for new contributors
- `help wanted`: Needs additional contribution

## Release Process

### Version Management

We follow Semantic Versioning (SemVer):

- `MAJOR.MINOR.PATCH` (e.g., `1.2.3`)
- `PATCH` version for bug fixes
- `MINOR` version for backward-compatible features
- `MAJOR` version for breaking changes

### Release Steps

1. **Prepare Release Branch**
   ```bash
   git checkout develop
   git pull upstream develop
   git checkout -b release/vX.Y.Z
   ```

2. **Update Version**
   ```bash
   # Update version in go.mod
   go mod edit -go=1.20
   go mod tidy
   
   # Update version in version.go
   echo 'const Version = "vX.Y.Z"' > internal/version.go
   ```

3. **Update Changelog**
   ```bash
   # Update CHANGELOG.md with new features and fixes
   git commit -m "docs: Update changelog for vX.Y.Z"
   ```

4. **Tag Release**
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   git push upstream vX.Y.Z
   ```

5. **Merge to Main**
   ```bash
   git checkout main
   git pull upstream main
   git merge --no-ff release/vX.Y.Z
   git push upstream main
   ```

6. **Create GitHub Release**
   - Go to GitHub releases
   - Create new release
   - Attach changelog
   - Mark as latest release

### Deployment Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Version bumped
- [ ] Database migrations reviewed
- [ ] Configuration reviewed
- [ ] Security scan completed
- [ ] Performance tests passed
- [ ] Deployment to staging successful
- [ ] Monitoring alerts verified

## Community Guidelines

### Communication Channels

- **Issues**: GitHub issues for bugs and feature requests
- **Discussions**: GitHub discussions for general questions
- **Pull Requests**: Code review and contribution discussion
- **Email**: Team mailing list for announcements

### Getting Help

1. **Search existing issues** before creating new ones
2. **Be specific** in your questions
3. **Provide context** (error messages, environment info)
4. **Be patient** - community members may have timezones

### Participation Guidelines

1. **Be respectful** and professional
2. **Welcome newcomers** and help them learn
3. **Focus on constructive criticism**
4. **Listen to feedback** and be open to different perspectives
5. **Follow the code of conduct**

### Recognition

Contributors will be recognized in:
- Release notes for significant contributions
- Contributors list in README
- Special recognition for bug fixing and documentation

### Maintainers

Project maintainers:
- Review pull requests
- Merge approved PRs
- Manage releases
- Respond to issues
- Set project direction

### Becoming a Maintainer

Requirements for maintainers:
- Active contributions to the project
- Understanding of project architecture
- Ability to review code
- Commit to regular participation
- Community support

## Acknowledgments

Special thanks to all contributors who have helped improve ByteDanceDemo. Your contributions make this project better for everyone!

---

*For questions about contributing, please open an issue labeled "question" or contact the maintainers directly.*
