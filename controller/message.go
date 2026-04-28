package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var tempChat = map[string][]Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	// Validate user ID
	userIdB, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid user ID format"})
		return
	}

	// Sanitize content
	if len(content) > 500 {
		content = content[:500]
	}

	chatKey := genChatKey(userId, userIdB)

	atomic.AddInt64(&messageIdSequence, 1)
	curMessage := Message{
		Id:         messageIdSequence,
		Content:    content,
		CreateTime: time.Now().Format(time.Kitchen),
	}

	if messages, exist := tempChat[chatKey]; exist {
		tempChat[chatKey] = append(messages, curMessage)
	} else {
		tempChat[chatKey] = []Message{curMessage}
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserId := c.Query("to_user_id")

	userIdB, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid user ID format"})
		return
	}

	chatKey := genChatKey(userId, userIdB)

	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}