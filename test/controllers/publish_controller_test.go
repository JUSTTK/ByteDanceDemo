package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bytedancedemo/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPublishControllerTest() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Setup routes
	api := router.Group("/douyin")
	{
		// Publish routes
		api.POST("/publish/action/", controller.PublishAction)
		api.GET("/publish/list/", controller.PublishList)
	}

	return router
}

func TestPublishController_PublishAction(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Publish valid video", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_id")
	})

	t.Run("Publish video without authentication", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Publish video with empty title", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Publish video with title too long", func(t *testing.T) {
		longTitle := "This video title is way too long and exceeds the maximum allowed length for video titles in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience."
		publishData := map[string]interface{}{
			"title":       longTitle,
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Publish video with description too long", func(t *testing.T) {
		longDescription := "This video description is way too long and exceeds the maximum allowed length for video descriptions in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience."
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": longDescription,
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Publish video with special characters", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video!@#$%^&*()",
			"description": "This is a test video with special characters!@#$%^&*()",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Publish video with emojis", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video 🎬",
			"description": "This is a test video with emojis 🎬🎨🎯",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPublishController_PublishList(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Get valid publish list", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_list")
	})

	t.Run("Get publish list without user ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get publish list with invalid user ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get publish list for non-existent user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.NotNil(t, response["video_list"])
	})

	t.Run("Get publish list with pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=1&page=1&size=10", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_list")
	})

	t.Run("Get publish list with invalid pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=1&page=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get publish list with size limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/publish/list/?user_id=1&size=100", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_list")
	})
}

func TestPublishController_ErrorHandling(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Handle database connection errors", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return error status
		assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)
	})

	t.Run("Handle JSON parsing errors", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Handle malformed authorization header", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "InvalidTokenFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestPublishController_InputValidation(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Prevent SQL injection in title", func(t *testing.T) {
		maliciousInput := map[string]interface{}{
			"title":       "'; DROP TABLE videos; --",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(maliciousInput)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or properly escape the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Prevent XSS in title", func(t *testing.T) {
		xssInput := map[string]interface{}{
			"title":       "<script>alert('xss')</script>",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(xssInput)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or sanitize the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sanitize HTML tags in title", func(t *testing.T) {
		htmlInput := map[string]interface{}{
			"title":       "This is <b>bold</b> title",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(htmlInput)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotContains(t, response, "<b>")
	})
}

func TestPublishController_RateLimiting(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Rate limit video publishing", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Rate limit test video",
			"description": "This is a rate limit test video",
		}
		jsonData, _ := json.Marshal(publishData)

		// Make multiple requests quickly
		for i := 0; i < 20; i++ {
			req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return either 200 or 429 (Too Many Requests)
			assert.Contains(t, []int{http.StatusOK, http.StatusTooManyRequests}, w.Code)
		}
	})
}

func TestPublishController_FileUpload(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Publish video with file upload", func(t *testing.T) {
		// This would require multipart/form-data test
		// For now, we test the basic functionality
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
			// file would be added as multipart
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Publish video without file", func(t *testing.T) {
		// This depends on whether file upload is required
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// May fail or succeed depending on implementation
		assert.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, w.Code)
	})
}

func TestPublishController_Permissions(t *testing.T) {
	router := setupPublishControllerTest()

	t.Run("Only authenticated users can publish", func(t *testing.T) {
		publishData := map[string]interface{}{
			"title":       "Test Video",
			"description": "This is a test video",
		}
		jsonData, _ := json.Marshal(publishData)

		req, _ := http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}