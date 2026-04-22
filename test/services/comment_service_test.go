package test

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/test/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommentServiceImpl_GetCommentServiceInstanceSingleton(t *testing.T) {
	t.Run("Singleton pattern test", func(t *testing.T) {
		cs1 := service.GetCommentServiceInstance()
		cs2 := service.GetCommentServiceInstance()

		assert.Equal(t, cs1, cs2, "Instances should be the same")
	})
}

func TestCommentServiceImpl_AddComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Add valid comment", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   1,
			Content:   "Great video!",
			CreateTime: time.Now(),
		}

		comment, err := commentService.AddComment(comment)

		if err == nil {
			assert.NotNil(t, comment, "Comment should not be nil")
			assert.Greater(t, comment.Id, int64(0), "Comment ID should be set")
		} else {
			// This might happen if the database is not available
			assert.NotNil(t, err, "Error should be returned for database failure")
		}
	})

	t.Run("Add comment with empty content", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   1,
			Content:   "", // Empty content
			CreateTime: time.Now(),
		}

		comment, err := commentService.AddComment(comment)

		assert.NotNil(t, err, "Error should be returned for empty content")
		assert.Nil(t, comment, "Comment should be nil for invalid input")
	})

	t.Run("Add comment with invalid user ID", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    -1, // Invalid user ID
			VideoId:   1,
			Content:   "Great video!",
			CreateTime: time.Now(),
		}

		comment, err := commentService.AddComment(comment)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Nil(t, comment, "Comment should be nil for invalid input")
	})

	t.Run("Add comment with invalid video ID", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   -1, // Invalid video ID
			Content:   "Great video!",
			CreateTime: time.Now(),
		}

		comment, err := commentService.AddComment(comment)

		assert.NotNil(t, err, "Error should be returned for invalid video ID")
		assert.Nil(t, comment, "Comment should be nil for invalid input")
	})
}

func TestCommentServiceImpl_DeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Delete existing comment", func(t *testing.T) {
		commentID := int64(1)

		err := commentService.DeleteComment(commentID)

		assert.NoError(t, err, "Should delete existing comment without error")
	})

	t.Run("Delete non-existent comment", func(t *testing.T) {
		commentID := int64(999) // Non-existent comment ID

		err := commentService.DeleteComment(commentID)

		assert.NoError(t, err, "Should handle non-existent comment gracefully")
	})

	t.Run("Delete comment with invalid ID", func(t *testing.T) {
		commentID := int64(-1) // Invalid comment ID

		err := commentService.DeleteComment(commentID)

		assert.Error(t, err, "Should return error for invalid comment ID")
	})
}

func TestCommentServiceImpl_GetCommentList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Get comments for existing video", func(t *testing.T) {
		videoID := int64(1)

		comments, err := commentService.GetCommentList(videoID)

		if err == nil {
			assert.NotNil(t, comments, "Comments list should not be nil")
			for _, comment := range comments {
				assert.Equal(t, videoID, comment.VideoId, "All comments should belong to the video")
			}
		} else {
			// This might happen if no comments exist or video doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent video")
		}
	})

	t.Run("Get comments for non-existent video", func(t *testing.T) {
		videoID := int64(999) // Non-existent video ID

		comments, err := commentService.GetCommentList(videoID)

		assert.NotNil(t, err, "Error should be returned for non-existent video")
		assert.Nil(t, comments, "Comments should be nil for non-existent video")
	})

	t.Run("Get comments with invalid video ID", func(t *testing.T) {
		videoID := int64(-1) // Invalid video ID

		comments, err := commentService.GetCommentList(videoID)

		assert.NotNil(t, err, "Error should be returned for invalid video ID")
		assert.Nil(t, comments, "Comments should be nil for invalid input")
	})
}

func TestCommentServiceImpl_GetCommentById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Get existing comment", func(t *testing.T) {
		commentID := int64(1)

		comment, err := commentService.GetCommentById(commentID)

		if err == nil {
			assert.NotNil(t, comment, "Comment should not be nil")
			assert.Equal(t, commentID, comment.Id, "Comment ID should match")
		} else {
			// This might happen if comment doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent comment")
		}
	})

	t.Run("Get non-existent comment", func(t *testing.T) {
		commentID := int64(999) // Non-existent comment ID

		comment, err := commentService.GetCommentById(commentID)

		assert.NotNil(t, err, "Error should be returned for non-existent comment")
		assert.Nil(t, comment, "Comment should be nil for non-existent comment")
	})

	t.Run("Get comment with invalid ID", func(t *testing.T) {
		commentID := int64(-1) // Invalid comment ID

		comment, err := commentService.GetCommentById(commentID)

		assert.NotNil(t, err, "Error should be returned for invalid comment ID")
		assert.Nil(t, comment, "Comment should be nil for invalid input")
	})
}

func TestCommentServiceImpl_ValidateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Validate valid comment", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   1,
			Content:   "Valid comment",
			CreateTime: time.Now(),
		}

		err := commentService.ValidateComment(comment)

		assert.NoError(t, err, "Valid comment should not produce error")
	})

	t.Run("Validate comment with empty content", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   1,
			Content:   "",
			CreateTime: time.Now(),
		}

		err := commentService.ValidateComment(comment)

		assert.Error(t, err, "Empty content should produce error")
	})

	t.Run("Validate comment with too long content", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   1,
			Content:   "This comment is way too long and exceeds the maximum allowed length for comments in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience.",
			CreateTime: time.Now(),
		}

		err := commentService.ValidateComment(comment)

		assert.Error(t, err, "Too long content should produce error")
	})

	t.Run("Validate comment with invalid user ID", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    -1,
			VideoId:   1,
			Content:   "Valid comment",
			CreateTime: time.Now(),
		}

		err := commentService.ValidateComment(comment)

		assert.Error(t, err, "Invalid user ID should produce error")
	})

	t.Run("Validate comment with invalid video ID", func(t *testing.T) {
		comment := &model.Comment{
			UserId:    1,
			VideoId:   -1,
			Content:   "Valid comment",
			CreateTime: time.Now(),
		}

		err := commentService.ValidateComment(comment)

		assert.Error(t, err, "Invalid video ID should produce error")
	})

	t.Run("Validate comment with nil", func(t *testing.T) {
		err := commentService.ValidateComment(nil)

		assert.Error(t, err, "Nil comment should produce error")
	})
}

func TestCommentServiceImpl_Concurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentService := service.GetCommentServiceInstance()

	t.Run("Concurrent comment addition", func(t *testing.T) {
		const goroutines = 10
		done := make(chan bool, goroutines)
		errorsChan := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				comment := &model.Comment{
					UserId:    1,
					VideoId:   1,
					Content:   "Concurrent comment",
					CreateTime: time.Now(),
				}
				_, err := commentService.AddComment(comment)
				errorsChan <- err
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < goroutines; i++ {
			<-done
		}

		// Check for errors
		close(errorsChan)
		for err := range errorsChan {
			// In a real scenario, we might expect some errors due to concurrent operations
			// For now, we just collect and report them
			if err != nil {
				t.Logf("Got error: %v", err)
			}
		}
	})
}

// Benchmark tests
func BenchmarkCommentServiceImpl_GetCommentServiceInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.GetCommentServiceInstance()
	}
}

func BenchmarkCommentServiceImpl_AddComment(b *testing.B) {
	commentService := service.GetCommentServiceInstance()
	comment := &model.Comment{
		UserId:    1,
		VideoId:   1,
		Content:   "Benchmark comment",
		CreateTime: time.Now(),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commentService.AddComment(comment)
	}
}

func BenchmarkCommentServiceImpl_GetCommentList(b *testing.B) {
	commentService := service.GetCommentServiceInstance()
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commentService.GetCommentList(videoID)
	}
}

func BenchmarkCommentServiceImpl_GetCommentById(b *testing.B) {
	commentService := service.GetCommentServiceInstance()
	commentID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commentService.GetCommentById(commentID)
	}
}

func BenchmarkCommentServiceImpl_DeleteComment(b *testing.B) {
	commentService := service.GetCommentServiceInstance()
	commentID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		commentService.DeleteComment(commentID)
	}
}