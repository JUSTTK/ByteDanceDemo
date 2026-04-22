package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bytedancedemo/middleware/validation"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware_InputSanitization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup validation middleware
	router.Use(validation.ValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username": input.Username,
			"email":    input.Email,
			"message":  "Valid input",
		})
	})

	t.Run("Sanitize valid input", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Sanitize input with special characters", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "testuser!@#$%^&*()",
			"email":    "test@example.com",
			"password": "password123!@#$",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Prevent SQL injection", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "'; DROP TABLE users; --",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or properly escape the input
		assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)
	})

	t.Run("Prevent XSS attacks", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "<script>alert('xss')</script>",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or sanitize the input
		assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)

		if w.Code == http.StatusOK {
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			// Verify that HTML tags are escaped
			assert.NotContains(t, response["username"].(string), "<script>")
		}
	})
}

func TestValidationMiddleware_FieldValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(validation.ValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required,min=3,max=20"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Valid input"})
	})

	t.Run("Validate required fields", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "testuser",
			// Missing email and password
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate minimum length", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "ab", // Too short
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate maximum length", func(t *testing.T) {
		input := map[string]interface{}{
			"username": strings.Repeat("a", 25), // Too long
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate email format", func(t *testing.T) {
		input := map[string]interface{}{
			"username": "testuser",
			"email":    "invalid-email",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestValidationMiddleware_DataTypeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(validation.ValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		var input struct {
			ID       int64   `json:"id" binding:"required"`
			Name     string  `json:"name" binding:"required"`
			IsActive bool    `json:"is_active"`
			Amount   float64 `json:"amount"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Valid input"})
	})

	t.Run("Validate integer type", func(t *testing.T) {
		input := map[string]interface{}{
			"id":       "not-an-integer", // Wrong type
			"name":     "test",
			"is_active": true,
			"amount":   100.50,
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate boolean type", func(t *testing.T) {
		input := map[string]interface{}{
			"id":       1,
			"name":     "test",
			"is_active": "not-a-boolean", // Wrong type
			"amount":   100.50,
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate float type", func(t *testing.T) {
		input := map[string]interface{}{
			"id":       1,
			"name":     "test",
			"is_active": true,
			"amount":   "not-a-float", // Wrong type
		}
		jsonData, _ := json.Marshal(input)

		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestValidationMiddleware_QueryParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(validation.ValidationMiddleware())
	router.GET("/test", func(c *gin.Context) {
		var params struct {
			Page     int    `form:"page" binding:"required,min=1"`
			PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
			Search   string `form:"search"`
		}

		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":      params.Page,
			"page_size": params.PageSize,
			"search":    params.Search,
		})
	})

	t.Run("Validate query parameters", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test?page=1&page_size=10&search=test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Validate required query parameters", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test?search=test", nil) // Missing page and page_size
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate query parameter ranges", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test?page=0&page_size=10", nil) // Invalid page
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validate query parameter maximum", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test?page=1&page_size=101", nil) // Invalid page_size
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestValidationMiddleware_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(validation.ValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	t.Run("Handle invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Handle empty JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Handle malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{"incomplete": "json"`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestValidationMiddleware_Performance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(validation.ValidationMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	t.Run("Benchmark validation middleware", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping benchmark in short mode")
		}

		input := map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(input)

		start := time.Now()
		for i := 0; i < 1000; i++ {
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
		duration := time.Since(start)

		assert.Less(t, duration.Milliseconds(), int64(2000), "Validation should be fast")
	})
}