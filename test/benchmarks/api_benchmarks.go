package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"bytedancedemo/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupBenchmarkAPI() *gin.Engine {
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

		// Publish routes
		api.POST("/publish/action/", controller.PublishAction)
		api.GET("/publish/list/", controller.PublishList)

		// Comment routes
		api.POST("/comment/action/", controller.CommentAction)
		api.GET("/comment/list/", controller.CommentList)

		// Favorite routes
		api.POST("/favorite/action/", controller.FavoriteAction)
		api.GET("/favorite/list/", controller.FavoriteList)

		// Relation routes
		api.POST("/relation/action/", controller.FollowAction)
		api.GET("/relation/follow/list/", controller.FollowList)
		api.GET("/relation/follower/list/", controller.FollowerList)

		// Message routes
		api.POST("/message/action/", controller.MessageAction)
		api.GET("/message/chat/", controller.MessageChat)
	}

	return router
}

func BenchmarkUserRegistration(b *testing.B) {
	router := setupBenchmarkAPI()

	userData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
		"email":    "benchuser@example.com",
	}
	jsonData, _ := json.Marshal(userData)

	b.Run("Serial user registration", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)
		}
	})

	b.Run("Concurrent user registration", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(b, http.StatusOK, w.Code)
			}(i)
		}

		wg.Wait()
	})
}

func BenchmarkUserLogin(b *testing.B) {
	router := setupBenchmarkAPI()

	// First register a user
	userData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Login data
	loginData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)

	b.Run("Serial user login", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)
		}
	})

	b.Run("Concurrent user login", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(b, http.StatusOK, w.Code)
			}()
		}

		wg.Wait()
	})
}

func BenchmarkPublishVideo(b *testing.B) {
	router := setupBenchmarkAPI()

	// First register and login a user
	userData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Publish data
	publishData := map[string]interface{}{
		"title":       "Benchmark Video",
		"description": "This is a benchmark video",
	}
	jsonData, _ = json.Marshal(publishData)

	b.Run("Serial video publishing", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)
		}
	})

	b.Run("Concurrent video publishing", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(b, http.StatusOK, w.Code)
			}()
		}

		wg.Wait()
	})
}

func BenchmarkCommentOperations(b *testing.B) {
	router := setupBenchmarkAPI()

	// First register and login a user
	userData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Comment data
	commentData := map[string]interface{}{
		"video_id":     1,
		"comment_text": "Benchmark comment",
	}
	jsonData, _ = json.Marshal(commentData)

	b.Run("Serial comment operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Add comment
			req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)

			// Get comment list
			req, _ = http.NewRequest("GET", "/douyin/comment/list/?video_id=1", nil)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)
		}
	})

	b.Run("Concurrent comment operations", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Add comment
				req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(b, http.StatusOK, w.Code)
			}()
		}

		wg.Wait()
	})
}

func BenchmarkFollowOperations(b *testing.B) {
	router := setupBenchmarkAPI()

	// First register users
	for i := 1; i <= 2; i++ {
		userData := map[string]interface{}{
			"username": fmt.Sprintf("benchuser%d", i),
			"password": "password123",
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}

	// Login as user 1
	user1Login := map[string]interface{}{
		"username": "benchuser1",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(user1Login)

	req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// Follow data
	followData := map[string]interface{}{
		"user_id":      1,
		"to_user_id":    2,
		"action_type":  1, // Follow action
	}
	jsonData, _ = json.Marshal(followData)

	b.Run("Serial follow operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "/douyin/relation/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)

			// Get following list
			req, _ = http.NewRequest("GET", "/douyin/relation/follow/list/?user_id=1", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusOK, w.Code)
		}
	})

	b.Run("Concurrent follow operations", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				req, _ := http.NewRequest("POST", "/douyin/relation/action/", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(b, http.StatusOK, w.Code)
			}()
		}

		wg.Wait()
	})
}

func BenchmarkAPIResponseTime(b *testing.B) {
	router := setupBenchmarkAPI()

	// First register and login a user
	userData := map[string]interface{}{
		"username": "benchuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	endpoints := []struct {
		method string
		path   string
		data   map[string]interface{}
	}{
		{"GET", "/douyin/user/?user_id=1", nil},
		{"GET", "/douyin/publish/list/?user_id=1", nil},
		{"GET", "/douyin/comment/list/?video_id=1", nil},
		{"GET", "/douyin/favorite/list/?user_id=1", nil},
		{"GET", "/douyin/relation/follow/list/?user_id=1", nil},
	}

	b.Run("API response time benchmark", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, endpoint := range endpoints {
				req, _ := http.NewRequest(endpoint.method, endpoint.path, nil)
				if endpoint.data != nil {
					jsonData, _ := json.Marshal(endpoint.data)
					req.Body = io.NopCloser(bytes.NewReader(jsonData))
					req.Header.Set("Content-Type", "application/json")
				}
				if endpoint.path != "/douyin/user/" {
					req.Header.Set("Authorization", "Bearer "+token)
				}

				w := httptest.NewRecorder()
				start := time.Now()

				router.ServeHTTP(w, req)

				duration := time.Since(start)
				assert.Equal(b, http.StatusOK, w.Code)
				assert.Less(b, duration.Milliseconds(), int64(1000), "API response should be fast")
			}
		}
	})
}

func BenchmarkAPIThroughput(b *testing.B) {
	router := setupBenchmarkAPI()

	// Test throughput with multiple concurrent requests
	concurrentRequests := 100
	requestsPerGoroutine := 10

	b.Run("API throughput benchmark", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < requestsPerGoroutine; j++ {
					// Simple GET request
					req, _ := http.NewRequest("GET", "/douyin/user/?user_id=1", nil)
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					assert.Equal(b, http.StatusOK, w.Code)
				}
			}()
		}

		wg.Wait()
	})
}

func BenchmarkAPIErrorHandling(b *testing.B) {
	router := setupBenchmarkAPI()

	b.Run("Error handling benchmark", func(b *testing.B) {
		errorRequests := []struct {
			method string
			path   string
			data   map[string]interface{}
		}{
			{"GET", "/douyin/user/?user_id=0", nil}, // Invalid user ID
			{"POST", "/douyin/user/login/", nil},    // Missing credentials
			{"GET", "/douyin/comment/list/", nil},     // Missing video ID
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, req := range errorRequests {
				jsonData, _ := json.Marshal(req.data)
				httpReq, _ := http.NewRequest(req.method, req.path, bytes.NewBuffer(jsonData))
				if req.data != nil {
					httpReq.Header.Set("Content-Type", "application/json")
				}

				w := httptest.NewRecorder()
				router.ServeHTTP(w, httpReq)

				// Error responses should return appropriate status codes
				assert.Contains(b, []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound}, w.Code)
			}
		}
	})
}

func BenchmarkAPIMemoryUsage(b *testing.B) {
	router := setupBenchmarkAPI()

	// This is a basic memory usage benchmark
	// For more detailed memory profiling, you would need to use runtime.ReadMemStats

	b.Run("Memory usage benchmark", func(b *testing.B) {
		// Create a large request body
		largeData := map[string]interface{}{
			"content": strings.Repeat("x", 1024*1024), // 1MB of data
		}
		jsonData, _ := json.Marshal(largeData)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(b, http.StatusBadRequest, w.Code) // Should fail due to large data
		}
	})
}