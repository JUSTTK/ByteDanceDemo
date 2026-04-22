package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bytedancedemo/controller"
	"bytedancedemo/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserControllerTest() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Setup routes
	api := router.Group("/douyin")
	{
		// User routes
		api.POST("/user/register/", controller.RegisterUser)
		api.POST("/user/login/", controller.LoginUser)
		api.GET("/user/", controller.GetUser)
		api.PUT("/user/", controller.UpdateUser)
	}

	return router
}

func TestUserController_RegisterUser(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Register valid user", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
			"email":    "test@example.com",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response, "user_id")
		assert.Contains(t, response, "token")
	})

	t.Run("Register user with missing fields", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Register user with duplicate username", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "admin",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// May return 409 Conflict or 400 Bad Request depending on implementation
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Register user with weak password", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "testuser",
			"password": "123", // Too short password
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Register user with invalid email", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
			"email":    "invalid-email",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_LoginUser(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Login with valid credentials", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "admin",
			"password": "123456",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user_id")
	})

	t.Run("Login with invalid credentials", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "admin",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Login with non-existent user", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "nonexistent",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Login with missing fields", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_GetUser(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Get user with valid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/user/?user_id=1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Contains(t, response, "id")
		assert.Contains(t, response, "name")
		assert.Contains(t, response, "follow_count")
		assert.Contains(t, response, "follower_count")
	})

	t.Run("Get user with invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/user/?user_id=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/user/?user_id=999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Get user without ID parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/user/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Update valid user", func(t *testing.T) {
		updateData := map[string]interface{}{
			"signature": "New signature",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/douyin/user/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update user without authentication", func(t *testing.T) {
		updateData := map[string]interface{}{
			"signature": "New signature",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/douyin/user/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Update user with invalid token", func(t *testing.T) {
		updateData := map[string]interface{}{
			"signature": "New signature",
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/douyin/user/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Update user with invalid data", func(t *testing.T) {
		updateData := map[string]interface{}{
			"signature": strings.Repeat("a", 1001), // Too long
		}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/douyin/user/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_ErrorHandling(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Handle database connection errors", func(t *testing.T) {
		// This would require mocking the database layer
		// For now, we test the response format
		userData := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Even if there's a database error, the response should be valid JSON
		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	})

	t.Run("Handle JSON parsing errors", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_RateLimiting(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Test rate limiting", func(t *testing.T) {
		userData := map[string]interface{}{
			"username": "ratelimituser",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		// Make multiple requests quickly
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return either 200 or 429 (Too Many Requests)
			assert.Contains(t, []int{http.StatusOK, http.StatusTooManyRequests}, w.Code)
		}
	})
}

func TestUserController_CORSHandling(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Test CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/douyin/user/register/", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should include CORS headers if enabled
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	})
}

func TestUserController_InputSanitization(t *testing.T) {
	router := setupUserControllerTest()

	t.Run("Test XSS prevention", func(t *testing.T) {
		maliciousInput := map[string]interface{}{
			"username": "<script>alert('xss')</script>",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(maliciousInput)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or sanitize the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Test SQL injection prevention", func(t *testing.T) {
		maliciousInput := map[string]interface{}{
			"username": "admin'; DROP TABLE users; --",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(maliciousInput)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or properly escape the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}