package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bytedancedemo/controller"
	"bytedancedemo/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationAPI contains integration tests for the entire API
func TestIntegrationAPI(t *testing.T) {
	// Setup test environment
	router := setupIntegrationTestRouter()

	// Test users for integration tests
	testUsers := []struct {
		username string
		password string
		email    string
	}{
		{"integrationuser1", "password123", "user1@example.com"},
		{"integrationuser2", "password123", "user2@example.com"},
		{"integrationuser3", "password123", "user3@example.com"},
	}

	t.Run("User Registration and Login Flow", func(t *testing.T) {
		// Register user
		userData := map[string]interface{}{
			"username": testUsers[0].username,
			"password": testUsers[0].password,
			"email":    testUsers[0].email,
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "user_id")
		assert.Contains(t, response, "token")

		// Login with same user
		loginData := map[string]interface{}{
			"username": testUsers[0].username,
			"password": testUsers[0].password,
		}
		jsonData, _ = json.Marshal(loginData)

		req, _ = http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "token")
		assert.Contains(t, response, "user_id")
	})

	t.Run("Video Publishing Flow", func(t *testing.T) {
		// First register a user
		userData := map[string]interface{}{
			"username": testUsers[1].username,
			"password": testUsers[1].password,
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		token := response["token"].(string)
		userID := int(response["user_id"].(float64))

		// Publish a video
		publishData := map[string]interface{}{
			"title":       "Integration Test Video",
			"description": "This is a video created during integration testing",
		}
		jsonData, _ = json.Marshal(publishData)

		req, _ = http.NewRequest("POST", "/douyin/publish/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		videoID := int(response["video_id"].(float64))
		assert.Greater(t, videoID, 0)

		// Get user's video list
		req, _ = http.NewRequest("GET", fmt.Sprintf("/douyin/publish/list/?user_id=%d", userID), nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_list")

		// Verify the video is in the list
		videoList := response["video_list"].([]interface{})
		assert.Len(t, videoList, 1)

		video := videoList[0].(map[string]interface{})
		assert.Equal(t, videoID, int(video["id"].(float64)))
	})

	t.Run("Follow/Unfollow Flow", func(t *testing.T) {
		// Register users
		for _, user := range testUsers[:2] {
			userData := map[string]interface{}{
				"username": user.username,
				"password": user.password,
			}
			jsonData, _ := json.Marshal(userData)

			req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
		}

		// Get user tokens
		tokens := make([]string, 2)
		for i, user := range testUsers[:2] {
			loginData := map[string]interface{}{
				"username": user.username,
				"password": user.password,
			}
			jsonData, _ := json.Marshal(loginData)

			req, _ := http.NewRequest("POST", "/douyin/user/login/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			tokens[i] = response["token"].(string)
		}

		// Follow user 2 from user 1
		followData := map[string]interface{}{
			"user_id":      2,
			"to_user_id":    3,
			"action_type":  1, // Follow action
		}
		jsonData, _ := json.Marshal(followData)

		req, _ := http.NewRequest("POST", "/douyin/relation/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokens[0])
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))

		// Check follow status (this would depend on the actual implementation)
		req, _ = http.NewRequest("GET", "/douyin/relation/follow/list/?user_id=2", nil)
		req.Header.Set("Authorization", "Bearer "+tokens[0])
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "user_list")

		// Unfollow user
		unfollowData := map[string]interface{}{
			"user_id":      2,
			"to_user_id":    3,
			"action_type":  2, // Unfollow action
		}
		jsonData, _ = json.Marshal(unfollowData)

		req, _ = http.NewRequest("POST", "/douyin/relation/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokens[0])
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
	})

	t.Run("Comment Flow", func(t *testing.T) {
		// Register and login a user
		userData := map[string]interface{}{
			"username": testUsers[2].username,
			"password": testUsers[2].password,
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		token := response["token"].(string)

		// Add a comment to a video
		commentData := map[string]interface{}{
			"video_id":     1,
			"comment_text": "This is a test comment",
		}
		jsonData, _ = json.Marshal(commentData)

		req, _ = http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "comment")

		comment := response["comment"].(map[string]interface{})
		commentID := int(comment["id"].(float64))

		// Get video comments
		req, _ = http.NewRequest("GET", "/douyin/comment/list/?video_id=1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "comment_list")

		// Verify the comment is in the list
		commentList := response["comment_list"].([]interface{})
		commentFound := false
		for _, c := range commentList {
			c := c.(map[string]interface{})
			if int(c["id"].(float64)) == commentID {
				commentFound = true
				break
			}
		}
		assert.True(t, commentFound, "Comment should be in the list")

		// Delete the comment
		deleteData := map[string]interface{}{
			"action_type": 1, // Delete action
			"comment_id":   commentID,
		}
		jsonData, _ = json.Marshal(deleteData)

		req, _ = http.NewRequest("POST", "/douyin/comment/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
	})

	t.Run("Favorite Flow", func(t *testing.T) {
		// Register and login a user
		userData := map[string]interface{}{
			"username": testUsers[0].username,
			"password": testUsers[0].password,
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		token := response["token"].(string)

		// Favorite a video
		favoriteData := map[string]interface{}{
			"video_id":    1,
			"action_type": 1, // Favorite action
		}
		jsonData, _ := json.Marshal(favoriteData)

		req, _ = http.NewRequest("POST", "/douyin/favorite/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))

		// Get favorite list
		req, _ = http.NewRequest("GET", "/douyin/favorite/list/?user_id=1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "video_list")

		// Unfavorite the video
		unfavoriteData := map[string]interface{}{
			"video_id":    1,
			"action_type": 2, // Unfavorite action
		}
		jsonData, _ = json.Marshal(unfavoriteData)

		req, _ = http.NewRequest("POST", "/douyin/favorite/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
	})

	t.Run("Message Flow", func(t *testing.T) {
		// Register and login users
		userData := map[string]interface{}{
			"username": testUsers[1].username,
			"password": testUsers[1].password,
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/douyin/user/register/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		token := response["token"].(string)

		// Send a message
		messageData := map[string]interface{}{
			"to_user_id":  2,
			"content":     "Hello! This is a test message.",
		}
		jsonData, _ = json.Marshal(messageData)

		req, _ = http.NewRequest("POST", "/douyin/message/action/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "message_id")

		// Get message list
		req, _ = http.NewRequest("GET", "/douyin/message/chat/?to_user_id=2", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 0, int(response["status_code"].(float64)))
		assert.Contains(t, response, "message_list")
	})
}

func setupIntegrationTestRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Setup all routes
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