package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockAuthMiddleware mocks the authentication middleware
type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) MockTokenValidation(token string, isValid bool, userId int64) {
	m.On("ValidateToken", token).Return(userId, isValid)
}

func (m *MockAuthMiddleware) ValidateToken(token string) (int64, bool) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Bool(1)
}

// SetupTestRoute creates a test route with mock middleware
func SetupTestRoute(handler gin.HandlerFunc, method, path string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add mock middleware
	router.Use(func(c *gin.Context) {
		// Mock token validation
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			mockAuth := &MockAuthMiddleware{}
			mockAuth.MockTokenValidation(token, true, 1) // Default to valid token with user ID 1
			userId, _ := mockAuth.ValidateToken(token)
			c.Set("user_id", userId)
		}
		c.Next()
	})

	// Add the handler
	router.Handle(method, path, handler)

	// Create test request
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	return w
}

// CreateAuthRequest creates a request with authentication header
func CreateAuthRequest(method, path, token string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

// MockRateLimiter mocks the rate limiting middleware
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) CheckLimit(key string) bool {
	args := m.Called(key)
	return args.Bool(1)
}

func (m *MockRateLimiter) SetLimit(key string, limit int) {
	m.Called(key, limit)
}

// MockCSRFMiddleware mocks the CSRF middleware
type MockCSRFMiddleware struct {
	mock.Mock
}

func (m *MockCSRFMiddleware) ValidateToken(token string) bool {
	args := m.Called(token)
	return args.Bool(1)
}

func (m *MockCSRFMiddleware) GenerateToken() string {
	return "mock-csrf-token"
}

// MockValidator mocks the validation middleware
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateRequest(req interface{}) error {
	args := m.Called(req)
	return args.Error(0)
}

// SetupValidValidation sets up valid validation
func (m *MockValidator) SetupValidValidation() {
	m.On("ValidateRequest", mock.Anything).Return(nil)
}

// SetupInvalidValidation sets up invalid validation
func (m *MockValidator) SetupInvalidValidation() {
	m.On("ValidateRequest", mock.Anything).Return(mock.Error)
}