package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"bytedancedemo/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to create multipart form data for file uploads
func createMultipartFormData(fieldName, filePath string, params map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range params {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, "", err
		}
	}

	// Add file
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filePath)
		if err != nil {
			return nil, "", err
		}

		_, err = bytes.CopyBuffer(part, file, make([]byte, 1024))
		if err != nil {
			return nil, "", err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.Boundary(), nil
}

// Helper function to create test user and get token
func createTestUserAndLogin(router *gin.Engine, username, password string) (string, int64) {
	// Register user
	userData := map[string]interface{}{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	token := response["token"].(string)
	userID := int64(response["user_id"].(float64))

	return token, userID
}

// Helper function to benchmark API endpoints
func benchmarkEndpoint(router *gin.Engine, method, path string, headers map[string]string) func(b *testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest(method, path, nil)
			for key, value := range headers {
				req.Header.Set(key, value)
			}
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Request failed: %s %s - Status: %d", method, path, w.Code)
			}

			if duration > 1*time.Second {
				b.Logf("Slow request: %s %s took %v", method, path, duration)
			}
		}
	}
}

// Helper function to benchmark concurrent API calls
func benchmarkConcurrentEndpoint(router *gin.Engine, method, path string, headers map[string]string, concurrent int) func(b *testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup

			for j := 0; j < concurrent; j++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()

					req, _ := http.NewRequest(method, path, nil)
					for key, value := range headers {
						req.Header.Set(key, value)
					}
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					if w.Code != http.StatusOK {
						b.Errorf("Concurrent request %d failed: %s %s - Status: %d", j, method, path, w.Code)
					}
				}(j)
			}

			wg.Wait()
		}
	}
}

// Test data generation helpers
func generateTestUsers(count int) []map[string]interface{} {
	users := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		users[i] = map[string]interface{}{
			"username": fmt.Sprintf("testuser%d", i),
			"password": "password123",
			"email":    fmt.Sprintf("testuser%d@example.com", i),
		}
	}
	return users
}

func generateTestVideos(count int) []map[string]interface{} {
	videos := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		videos[i] = map[string]interface{}{
			"title":       fmt.Sprintf("Test Video %d", i),
			"description": fmt.Sprintf("Description for test video %d", i),
			"author_id":   1,
		}
	}
	return videos
}

func generateTestComments(count int) []map[string]interface{} {
	comments := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		comments[i] = map[string]interface{}{
			"video_id":     1,
			"user_id":      1,
			"comment_text": fmt.Sprintf("This is test comment %d", i),
		}
	}
	return comments
}

// Test validation helpers
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidUsername(username string) bool {
	return len(username) >= 3 && len(username) <= 20
}

func isValidPassword(password string) bool {
	return len(password) >= 8
}

// Test cleanup helper
func cleanupTestData(router *gin.Engine) {
	// This would clean up any test data created during tests
	// For example, you could delete test users, videos, etc.
}

// Benchmark helpers
func benchmarkMemoryUsage(b *testing.B, f func()) {
	var m1, m2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	b.ReportMetric(float64(m2.Alloc-m1.Alloc)/float64(b.N), "bytes/op")
}

func benchmarkCPUUsage(b *testing.B, f func()) {
	start := time.Now()
	for i := 0; i < b.N; i++ {
		f()
	}
	elapsed := time.Since(start)
	b.ReportMetric(float64(elapsed)/float64(b.N), "ns/op")
}

// Performance test helpers
func assertResponseTime(t *testing.T, w *httptest.ResponseRecorder, maxDuration time.Duration) {
	duration := time.Since(w.Result().Request.Context().Value("startTime").(time.Time))
	assert.Less(t, duration, maxDuration, fmt.Sprintf("Response time %v exceeds maximum %v", duration, maxDuration))
}

func assertMemoryUsage(t *testing.T, currentMB float64, maxMB float64) {
	assert.Less(t, currentMB, maxMB, fmt.Sprintf("Memory usage %.2fMB exceeds maximum %.2fMB", currentMB, maxMB))
}

// Load testing helpers
func simulateLoad(router *gin.Engine, endpoint string, duration time.Duration, concurrency int) error {
	var wg sync.WaitGroup
	start := time.Now()
	done := make(chan bool)

	// Start load test
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					req, _ := http.NewRequest("GET", endpoint, nil)
					w := httptest.NewRecorder()
					router.ServeHTTP(w, req)

					if w.Code != http.StatusOK {
						fmt.Printf("Request failed from goroutine %d: %d\n", id, w.Code)
					}

					time.Sleep(time.Millisecond * 10) // Small delay
				}
			}
		}(i)
	}

	// Wait for duration
	time.Sleep(duration)
	close(done)
	wg.Wait()

	fmt.Printf("Load test completed. Duration: %v, Concurrency: %d\n", duration, concurrency)
	return nil
}

// Stress testing helpers
func simulateStress(router *gin.Engine, endpoints []string, concurrency int, duration time.Duration) error {
	var wg sync.WaitGroup
	start := time.Now()
	done := make(chan bool)

	// Start stress test
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				default:
					// Pick random endpoint
					endpoint := endpoints[rand.Intn(len(endpoints))]
					req, _ := http.NewRequest("GET", endpoint, nil)
					w := httptest.NewRecorder()
					router.ServeHTTP(w, req)

					if w.Code != http.StatusOK {
						fmt.Printf("Request failed from goroutine %d: %s - %d\n", id, endpoint, w.Code)
					}

					// Random delay between 1ms and 100ms
					delay := time.Duration(rand.Intn(100)+1) * time.Millisecond
					time.Sleep(delay)
				}
			}
		}(i)
	}

	// Wait for duration
	time.Sleep(duration)
	close(done)
	wg.Wait()

	fmt.Printf("Stress test completed. Duration: %v, Concurrency: %d\n", duration, concurrency)
	return nil
}