package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytedancedemo/middleware/jwt"
	"bytedancedemo/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware_TokenGeneration(t *testing.T) {
	t.Run("Generate valid token", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		token, err := jwt.GenerateToken(userID, username)

		assert.NoError(t, err, "Should generate token without error")
		assert.NotEmpty(t, token, "Token should not be empty")
	})

	t.Run("Generate token with empty username", func(t *testing.T) {
		userID := int64(1)
		username := ""

		token, err := jwt.GenerateToken(userID, username)

		assert.Error(t, err, "Should return error for empty username")
		assert.Empty(t, token, "Token should be empty for invalid input")
	})

	t.Run("Generate token with invalid user ID", func(t *testing.T) {
		userID := int64(-1)
		username := "testuser"

		token, err := jwt.GenerateToken(userID, username)

		assert.Error(t, err, "Should return error for invalid user ID")
		assert.Empty(t, token, "Token should be empty for invalid input")
	})
}

func TestJWTMiddleware_TokenValidation(t *testing.T) {
	t.Run("Validate valid token", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		token, _ := jwt.GenerateToken(userID, username)

		claims, err := jwt.ParseToken(token)

		assert.NoError(t, err, "Should validate valid token without error")
		assert.NotNil(t, claims, "Claims should not be nil")
		assert.Equal(t, userID, claims.UserID, "User ID should match")
		assert.Equal(t, username, claims.Username, "Username should match")
	})

	t.Run("Validate invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		claims, err := jwt.ParseToken(invalidToken)

		assert.Error(t, err, "Should return error for invalid token")
		assert.Nil(t, claims, "Claims should be nil for invalid token")
	})

	t.Run("Validate empty token", func(t *testing.T) {
		emptyToken := ""

		claims, err := jwt.ParseToken(emptyToken)

		assert.Error(t, err, "Should return error for empty token")
		assert.Nil(t, claims, "Claims should be nil for empty token")
	})

	t.Run("Validate malformed token", func(t *testing.T) {
		malformedToken := "Bearer invalid.token"

		claims, err := jwt.ParseToken(malformedToken)

		assert.Error(t, err, "Should return error for malformed token")
		assert.Nil(t, claims, "Claims should be nil for malformed token")
	})
}

func TestJWTMiddleware_TokenExpiry(t *testing.T) {
	t.Run("Token expiration", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		// This test would require mocking time or setting a very short expiration
		// For now, we test the basic functionality
		token, err := jwt.GenerateToken(userID, username)

		assert.NoError(t, err, "Should generate token without error")

		// Check if token is still valid immediately
		claims, err := jwt.ParseToken(token)

		assert.NoError(t, err, "Token should be valid immediately after generation")
		assert.NotNil(t, claims, "Claims should not be nil")
	})
}

func TestJWTMiddleware_GinIntegration(t *testing.T) {
	// Setup Gin router with JWT middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add JWT middleware to a protected route
	protected := router.Group("/protected")
	protected.Use(jwt.JWTAuthMiddleware())
	{
		protected.GET("/test", func(c *gin.Context) {
			userID := c.GetInt64("user_id")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"message": "Access granted",
			})
		})
	}

	t.Run("Access protected route with valid token", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		token, _ := jwt.GenerateToken(userID, username)

		req, _ := http.NewRequest("GET", "/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Access protected route without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Access protected route with invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		req, _ := http.NewRequest("GET", "/protected/test", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Access protected route with malformed authorization header", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		token, _ := jwt.GenerateToken(userID, username)

		req, _ := http.NewRequest("GET", "/protected/test", nil)
		req.Header.Set("Authorization", token) // Missing "Bearer " prefix
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestJWTMiddleware_ContextExtraction(t *testing.T) {
	// Setup Gin router with JWT middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(jwt.JWTAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		username := c.GetString("username")

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": username,
		})
	})

	t.Run("Extract user ID from token", func(t *testing.T) {
		userID := int64(1)
		username := "testuser"

		token, _ := jwt.GenerateToken(userID, username)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response contains correct user ID
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, userID, int64(response["user_id"].(float64)))
		assert.Equal(t, username, response["username"])
	})
}

func TestJWTMiddleware_ErrorHandling(t *testing.T) {
	// Setup Gin router with JWT middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(jwt.JWTAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted",
		})
	})

	t.Run("Handle missing authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Handle token with invalid signature", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Handle token with wrong secret", func(t *testing.T) {
		// This would require mocking the JWT secret
		// For now, we test with a clearly invalid token
		token := "invalid.token.with.wrong.secret"

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestJWTMiddleware_Performance(t *testing.T) {
	t.Run("Benchmark token generation", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping benchmark in short mode")
		}

		userID := int64(1)
		username := "testuser"

		start := time.Now()
		for i := 0; i < 1000; i++ {
			_, err := jwt.GenerateToken(userID, username)
			assert.NoError(t, err)
		}
		duration := time.Since(start)

		assert.Less(t, duration.Milliseconds(), int64(1000), "Token generation should be fast")
	})

	t.Run("Benchmark token validation", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping benchmark in short mode")
		}

		userID := int64(1)
		username := "testuser"
		token, _ := jwt.GenerateToken(userID, username)

		start := time.Now()
		for i := 0; i < 1000; i++ {
			_, err := jwt.ParseToken(token)
			assert.NoError(t, err)
		}
		duration := time.Since(start)

		assert.Less(t, duration.Milliseconds(), int64(1000), "Token validation should be fast")
	})
}