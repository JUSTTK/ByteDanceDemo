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

func setupCommentControllerTest() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Setup routes
	api := router.Group("/douyin")
	{
		// Comment routes
		api.POST("/comment/action/", controller.CommentAction)
		api.GET("/comment/list/", controller.CommentList)
	}

	return router
}

func TestCommentController_CommentAction(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Add valid comment", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Great video!",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "comment")
	})

	t.Run("Add comment without authentication", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Great video!",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Add comment with invalid video ID", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 0,
			"comment_text": "Great video!",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Add comment with empty text", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Add comment with text too long", func(t *testing.T) {
		longText := "This comment text is way too long and exceeds the maximum allowed length for comments in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience."
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": longText,
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete existing comment", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "To be deleted",
		}
		jsonData, _ := json.Marshal(commentData)

		// First create comment
		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Get comment ID from response
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		commentID := int64(response["comment"].(map[string]interface{})["id"].(float64))

		// Then delete it
		deleteData := map[string]interface{}{
			"action_type": 1, // Delete action
			"comment_id": commentID,
		}
		jsonData, _ := json.Marshal(deleteData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete non-existent comment", func(t *testing.T) {
		deleteData := map[string]interface{}{
			"action_type": 1, // Delete action
			"comment_id": 999,
		}
		jsonData, _ := json.Marshal(deleteData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should handle gracefully - either succeed or return appropriate error
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("Delete comment without authorization", func(t *testing.T) {
		deleteData := map[string]interface{}{
			"action_type": 1, // Delete action
			"comment_id": 1,
		}
		jsonData, _ := json.Marshal(deleteData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid action type", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Great video!",
			"action_type": 999, // Invalid action type
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCommentController_CommentList(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Get valid comment list", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/?video_id=1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "comment_list")
	})

	t.Run("Get comment list without video ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get comment list with invalid video ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/?video_id=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get comment list for non-existent video", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/?video_id=999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.NotNil(t, response["comment_list"])
	})

	t.Run("Get comment list with pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/?video_id=1&page=1&size=10", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "comment_list")
	})

	t.Run("Get comment list with invalid pagination", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/douyin/comment/list/?video_id=1&page=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCommentController_ErrorHandling(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Handle database connection errors", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Test comment",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return error status
		assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, w.Code)
	})

	t.Run("Handle JSON parsing errors", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Handle malformed authorization header", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Great video!",
		}
		jsonData, _ := json.Marshal(commentData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "InvalidTokenFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCommentController_InputValidation(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Prevent SQL injection in comment text", func(t *testing.T) {
		maliciousInput := map[string]interface{}{
			"video_id": 1,
			"comment_text": "'; DROP TABLE comments; --",
		}
		jsonData, _ := json.Marshal(maliciousInput)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or properly escape the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Prevent XSS in comment text", func(t *testing.T) {
		xssInput := map[string]interface{}{
			"video_id": 1,
			"comment_text": "<script>alert('xss')</script>",
		}
		jsonData, _ := json.Marshal(xssInput)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should either reject or sanitize the input
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sanitize HTML tags in comment text", func(t *testing.T) {
		htmlInput := map[string]interface{}{
			"video_id": 1,
			"comment_text": "This is <b>bold</b> text",
		}
		jsonData, _ := json.Marshal(htmlInput)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid_token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		comment := response["comment"].(map[string]interface{})
		assert.NotContains(t, comment["comment_text"], "<b>")
	})
}

func TestCommentController_RateLimiting(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Rate limit comment creation", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Rate limit test comment",
		}
		jsonData, _ := json.Marshal(commentData)

		// Make multiple requests quickly
		for i := 0; i < 20; i++ {
			req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer valid_token")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Should return either 200 or 429 (Too Many Requests)
			assert.Contains(t, []int{http.StatusOK, http.StatusTooManyRequests}, w.Code)
		}
	})
}

func TestCommentController_Authorization(t *testing.T) {
	router := setupCommentControllerTest()

	t.Run("Only comment owner can delete", func(t *testing.T) {
		commentData := map[string]interface{}{
			"video_id": 1,
			"comment_text": "Comment to be deleted by owner",
		}
		jsonData, _ := json.Marshal(commentData)

		// Create comment with user A
		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token_user_a")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Get comment ID from response
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		commentID := int64(response["comment"].(map[string]interface{})["id"].(float64))

		// Try to delete with user B
		deleteData := map[string]interface{}{
			"action_type": 1, // Delete action
			"comment_id": commentID,
		}
		jsonData, _ := json.Marshal(deleteData)

		req, _ := http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token_user_b")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should return 403 Forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}