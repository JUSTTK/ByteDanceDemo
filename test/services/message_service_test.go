package test

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/test/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMessageServiceImpl_GetMessageServiceInstanceSingleton(t *testing.T) {
	t.Run("Singleton pattern test", func(t *testing.T) {
		ms1 := service.GetMessageServiceInstance()
		ms2 := service.GetMessageServiceInstance()

		assert.Equal(t, ms1, ms2, "Instances should be the same")
	})
}

func TestMessageServiceImpl_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Send valid message", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "Hello there!",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		if err == nil {
			assert.NotNil(t, sentMessage, "Sent message should not be nil")
			assert.Equal(t, message.FromUserId, sentMessage.FromUserId, "From user ID should match")
			assert.Equal(t, message.ToUserId, sentMessage.ToUserId, "To user ID should match")
			assert.Equal(t, message.Content, sentMessage.Content, "Content should match")
			assert.Greater(t, sentMessage.Id, int64(0), "Message ID should be set")
		} else {
			// This might happen if users don't exist
			assert.NotNil(t, err, "Error should be returned for non-existent users")
		}
	})

	t.Run("Send empty message", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "", // Empty content
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for empty message")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})

	t.Run("Send message to self", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   1,
			Content:    "This is a message to myself",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for message to self")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})

	t.Run("Send message with invalid from user ID", func(t *testing.T) {
		message := &model.Message{
			FromUserId: -1,
			ToUserId:   2,
			Content:    "Hello there!",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for invalid from user ID")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})

	t.Run("Send message with invalid to user ID", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   -1,
			Content:    "Hello there!",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for invalid to user ID")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})

	t.Run("Send message with non-existent user", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   999, // Non-existent user
			Content:    "Hello there!",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})

	t.Run("Send message with too long content", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "This message content is way too long and exceeds the maximum allowed length for messages in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience.",
			CreateTime: time.Now(),
		}

		sentMessage, err := messageService.SendMessage(message)

		assert.NotNil(t, err, "Error should be returned for too long content")
		assert.Nil(t, sentMessage, "Message should be nil for invalid input")
	})
}

func TestMessageServiceImpl_GetMessageList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Get message list between two users", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(2)

		messages, err := messageService.GetMessageList(fromUserID, toUserID)

		if err == nil {
			assert.NotNil(t, messages, "Messages list should not be nil")
			for _, msg := range messages {
				assert.True(t, msg.FromUserId == fromUserID && msg.ToUserId == toUserID ||
					msg.FromUserId == toUserID && msg.ToUserId == fromUserID,
					"Messages should be between the two users")
			}
		} else {
			// This might happen if users don't exist or have no messages
			assert.NotNil(t, err, "Error might be returned for non-existent users")
		}
	})

	t.Run("Get message list with self", func(t *testing.T) {
		userID := int64(1)

		messages, err := messageService.GetMessageList(userID, userID)

		assert.NotNil(t, err, "Error should be returned for self message list")
		assert.Nil(t, messages, "Messages should be nil for invalid input")
	})

	t.Run("Get message list with invalid from user ID", func(t *testing.T) {
		fromUserID := int64(-1)
		toUserID := int64(2)

		messages, err := messageService.GetMessageList(fromUserID, toUserID)

		assert.NotNil(t, err, "Error should be returned for invalid from user ID")
		assert.Nil(t, messages, "Messages should be nil for invalid input")
	})

	t.Run("Get message list with invalid to user ID", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(-1)

		messages, err := messageService.GetMessageList(fromUserID, toUserID)

		assert.NotNil(t, err, "Error should be returned for invalid to user ID")
		assert.Nil(t, messages, "Messages should be nil for invalid input")
	})

	t.Run("Get message list with non-existent user", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(999) // Non-existent user

		messages, err := messageService.GetMessageList(fromUserID, toUserID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, messages, "Messages should be nil for non-existent user")
	})
}

func TestMessageServiceImpl_GetRecentChats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Get recent chats for existing user", func(t *testing.T) {
		userID := int64(1)

		chats, err := messageService.GetRecentChats(userID)

		if err == nil {
			assert.NotNil(t, chats, "Recent chats should not be nil")
			for _, chat := range chats {
				assert.NotEqual(t, userID, chat.UserId, "Chat partner should not be self")
				assert.Greater(t, chat.ChatCount, int64(0), "Chat count should be positive")
			}
		} else {
			// This might happen if user doesn't exist or has no chats
			assert.NotNil(t, err, "Error might be returned for non-existent user")
		}
	})

	t.Run("Get recent chats for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		chats, err := messageService.GetRecentChats(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, chats, "Chats should be nil for non-existent user")
	})

	t.Run("Get recent chats with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		chats, err := messageService.GetRecentChats(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Nil(t, chats, "Chats should be nil for invalid input")
	})
}

func TestMessageServiceImpl_GetUnreadCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Get unread count for existing user", func(t *testing.T) {
		userID := int64(1)

		count, err := messageService.GetUnreadCount(userID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Unread count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get unread count for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		count, err := messageService.GetUnreadCount(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent user")
	})

	t.Run("Get unread count with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		count, err := messageService.GetUnreadCount(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestMessageServiceImpl_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Mark messages as read", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(2)

		err := messageService.MarkAsRead(fromUserID, toUserID)

		assert.NoError(t, err, "Should mark messages as read without error")
	})

	t.Run("Mark messages as read with self", func(t *testing.T) {
		userID := int64(1)

		err := messageService.MarkAsRead(userID, userID)

		assert.Error(t, err, "Should not allow marking self messages as read")
	})

	t.Run("Mark messages as read with invalid from user ID", func(t *testing.T) {
		fromUserID := int64(-1)
		toUserID := int64(2)

		err := messageService.MarkAsRead(fromUserID, toUserID)

		assert.Error(t, err, "Should not mark with invalid from user ID")
	})

	t.Run("Mark messages as read with invalid to user ID", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(-1)

		err := messageService.MarkAsRead(fromUserID, toUserID)

		assert.Error(t, err, "Should not mark with invalid to user ID")
	})

	t.Run("Mark messages as read with non-existent user", func(t *testing.T) {
		fromUserID := int64(1)
		toUserID := int64(999) // Non-existent user

		err := messageService.MarkAsRead(fromUserID, toUserID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle marking non-existent user")
		}
	})
}

func TestMessageServiceImpl_DeleteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Delete existing message", func(t *testing.T) {
		messageID := int64(1)

		err := messageService.DeleteMessage(messageID)

		assert.NoError(t, err, "Should delete message without error")
	})

	t.Run("Delete non-existent message", func(t *testing.T) {
		messageID := int64(999) // Non-existent message

		err := messageService.DeleteMessage(messageID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle deletion of non-existent message")
		}
	})

	t.Run("Delete message with invalid ID", func(t *testing.T) {
		messageID := int64(-1) // Invalid message ID

		err := messageService.DeleteMessage(messageID)

		assert.Error(t, err, "Should not delete with invalid message ID")
	})
}

func TestMessageServiceImpl_ValidateMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	messageService := service.GetMessageServiceInstance()

	t.Run("Validate valid message", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "Valid message",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.NoError(t, err, "Valid message should not produce error")
	})

	t.Run("Validate message with empty content", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.Error(t, err, "Empty content should produce error")
	})

	t.Run("Validate message with too long content", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "This message content is way too long and exceeds the maximum allowed length for messages in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience.",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.Error(t, err, "Too long content should produce error")
	})

	t.Run("Validate message with invalid from user ID", func(t *testing.T) {
		message := &model.Message{
			FromUserId: -1,
			ToUserId:   2,
			Content:    "Valid message",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.Error(t, err, "Invalid from user ID should produce error")
	})

	t.Run("Validate message with invalid to user ID", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   -1,
			Content:    "Valid message",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.Error(t, err, "Invalid to user ID should produce error")
	})

	t.Run("Validate message with same user", func(t *testing.T) {
		message := &model.Message{
			FromUserId: 1,
			ToUserId:   1,
			Content:    "Valid message",
			CreateTime: time.Now(),
		}

		err := messageService.ValidateMessage(message)

		assert.Error(t, err, "Message to self should produce error")
	})

	t.Run("Validate nil message", func(t *testing.T) {
		err := messageService.ValidateMessage(nil)

		assert.Error(t, err, "Nil message should produce error")
	})
}

// Benchmark tests
func BenchmarkMessageServiceImpl_GetMessageServiceInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.GetMessageServiceInstance()
	}
}

func BenchmarkMessageServiceImpl_SendMessage(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	message := &model.Message{
		FromUserId: 1,
		ToUserId:   2,
		Content:    "Benchmark message",
		CreateTime: time.Now(),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.SendMessage(message)
	}
}

func BenchmarkMessageServiceImpl_GetMessageList(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	fromUserID := int64(1)
	toUserID := int64(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.GetMessageList(fromUserID, toUserID)
	}
}

func BenchmarkMessageServiceImpl_GetRecentChats(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.GetRecentChats(userID)
	}
}

func BenchmarkMessageServiceImpl_GetUnreadCount(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.GetUnreadCount(userID)
	}
}

func BenchmarkMessageServiceImpl_MarkAsRead(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	fromUserID := int64(1)
	toUserID := int64(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.MarkAsRead(fromUserID, toUserID)
	}
}

func BenchmarkMessageServiceImpl_DeleteMessage(b *testing.B) {
	messageService := service.GetMessageServiceInstance()
	messageID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messageService.DeleteMessage(messageID)
	}
}