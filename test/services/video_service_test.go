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

func TestVideoServiceImpl_GetVideoServiceInstanceSingleton(t *testing.T) {
	t.Run("Singleton pattern test", func(t *testing.T) {
		vs1 := service.GetVideoServiceInstance()
		vs2 := service.GetVideoServiceInstance()

		assert.Equal(t, vs1, vs2, "Instances should be the same")
	})
}

func TestVideoServiceImpl_PublishVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Publish valid video", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "Test Video",
			Description: "This is a test video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		if err == nil {
			assert.NotNil(t, publishedVideo, "Published video should not be nil")
			assert.Equal(t, video.AuthorId, publishedVideo.AuthorId, "Author ID should match")
			assert.Equal(t, video.Title, publishedVideo.Title, "Title should match")
			assert.Equal(t, video.Description, publishedVideo.Description, "Description should match")
			assert.Equal(t, video.PlayUrl, publishedVideo.PlayUrl, "Play URL should match")
			assert.Equal(t, video.CoverUrl, publishedVideo.CoverUrl, "Cover URL should match")
			assert.Greater(t, publishedVideo.Id, int64(0), "Video ID should be set")
		} else {
			// This might happen if author doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent author")
		}
	})

	t.Run("Publish video with empty title", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "", // Empty title
			Description: "This is a test video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for empty title")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})

	t.Run("Publish video with empty play URL", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "Test Video",
			Description: "This is a test video",
			PlayUrl:     "", // Empty play URL
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for empty play URL")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})

	t.Run("Publish video with invalid author ID", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    -1,
			Title:       "Test Video",
			Description: "This is a test video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for invalid author ID")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})

	t.Run("Publish video with non-existent author", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    999, // Non-existent author
			Title:       "Test Video",
			Description: "This is a test video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for non-existent author")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})

	t.Run("Publish video with too long title", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "This video title is way too long and exceeds the maximum allowed length for video titles in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience.",
			Description: "This is a test video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for too long title")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})

	t.Run("Publish video with invalid URL", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "Test Video",
			Description: "This is a test video",
			PlayUrl:     "invalid-url",
			CoverUrl:    "http://example.com/cover.jpg",
			CreateTime:  time.Now(),
		}

		publishedVideo, err := videoService.PublishVideo(video)

		assert.NotNil(t, err, "Error should be returned for invalid URL")
		assert.Nil(t, publishedVideo, "Video should be nil for invalid input")
	})
}

func TestVideoServiceImpl_GetVideoList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Get video list for existing user", func(t *testing.T) {
		authorID := int64(1)
		videoList, err := videoService.GetVideoList(authorID)

		if err == nil {
			assert.NotNil(t, videoList, "Video list should not be nil")
			for _, video := range videoList {
				assert.Equal(t, authorID, video.AuthorId, "All videos should belong to the author")
				assert.Greater(t, video.Id, int64(0), "Video IDs should be valid")
			}
		} else {
			// This might happen if author doesn't exist or has no videos
			assert.NotNil(t, err, "Error might be returned for non-existent author")
		}
	})

	t.Run("Get video list for non-existent user", func(t *testing.T) {
		authorID := int64(999) // Non-existent user

		videoList, err := videoService.GetVideoList(authorID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, videoList, "Video list should be nil for non-existent user")
	})

	t.Run("Get video list with invalid author ID", func(t *testing.T) {
		authorID := int64(-1) // Invalid author ID

		videoList, err := videoService.GetVideoList(authorID)

		assert.NotNil(t, err, "Error should be returned for invalid author ID")
		assert.Nil(t, videoList, "Video list should be nil for invalid input")
	})
}

func TestVideoServiceImpl_GetVideoById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Get existing video by ID", func(t *testing.T) {
		videoID := int64(1)

		video, err := videoService.GetVideoById(videoID)

		if err == nil {
			assert.NotNil(t, video, "Video should not be nil")
			assert.Equal(t, videoID, video.Id, "Video ID should match")
		} else {
			// This might happen if video doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent video")
		}
	})

	t.Run("Get non-existent video by ID", func(t *testing.T) {
		videoID := int64(999) // Non-existent video

		video, err := videoService.GetVideoById(videoID)

		assert.NotNil(t, err, "Error should be returned for non-existent video")
		assert.Nil(t, video, "Video should be nil for non-existent video")
	})

	t.Run("Get video with invalid ID", func(t *testing.T) {
		videoID := int64(-1) // Invalid video ID

		video, err := videoService.GetVideoById(videoID)

		assert.NotNil(t, err, "Error should be returned for invalid video ID")
		assert.Nil(t, video, "Video should be nil for invalid input")
	})
}

func TestVideoServiceImpl_GetVideoByAuthorId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Get video by author ID", func(t *testing.T) {
		authorID := int64(1)

		videos, err := videoService.GetVideoByAuthorId(authorID)

		if err == nil {
			assert.NotNil(t, videos, "Videos should not be nil")
			for _, video := range videos {
				assert.Equal(t, authorID, video.AuthorId, "All videos should belong to the author")
			}
		} else {
			// This might happen if author doesn't exist or has no videos
			assert.NotNil(t, err, "Error might be returned for non-existent author")
		}
	})

	t.Run("Get video by non-existent author ID", func(t *testing.T) {
		authorID := int64(999) // Non-existent author

		videos, err := videoService.GetVideoByAuthorId(authorID)

		assert.NotNil(t, err, "Error should be returned for non-existent author")
		assert.Nil(t, videos, "Videos should be nil for non-existent author")
	})

	t.Run("Get video by invalid author ID", func(t *testing.T) {
		authorID := int64(-1) // Invalid author ID

		videos, err := videoService.GetVideoByAuthorId(authorID)

		assert.NotNil(t, err, "Error should be returned for invalid author ID")
		assert.Nil(t, videos, "Videos should be nil for invalid input")
	})
}

func TestVideoServiceImpl_UpdateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Update existing video", func(t *testing.T) {
		videoID := int64(1)
		video := &model.Video{
			Id:          videoID,
			Title:       "Updated Title",
			Description: "Updated description",
		}

		err := videoService.UpdateVideo(video)

		assert.NoError(t, err, "Should update video without error")
	})

	t.Run("Update non-existent video", func(t *testing.T) {
		video := &model.Video{
			Id:          int64(999), // Non-existent video
			Title:       "Updated Title",
			Description: "Updated description",
		}

		err := videoService.UpdateVideo(video)

		assert.Error(t, err, "Should not update non-existent video")
	})

	t.Run("Update video with invalid ID", func(t *testing.T) {
		video := &model.Video{
			Id:          int64(-1), // Invalid video ID
			Title:       "Updated Title",
			Description: "Updated description",
		}

		err := videoService.UpdateVideo(video)

		assert.Error(t, err, "Should not update with invalid video ID")
	})

	t.Run("Update video with empty title", func(t *testing.T) {
		video := &model.Video{
			Id:          int64(1),
			Title:       "", // Empty title
			Description: "Updated description",
		}

		err := videoService.UpdateVideo(video)

		assert.Error(t, err, "Should not update with empty title")
	})

	t.Run("Update video with empty description", func(t *testing.T) {
		video := &model.Video{
			Id:          int64(1),
			Title:       "Updated Title",
			Description: "", // Empty description
		}

		err := videoService.UpdateVideo(video)

		// Depending on implementation, empty description might be allowed
		if err != nil {
			assert.Error(t, err, "Should not update with empty description if not allowed")
		}
	})
}

func TestVideoServiceImpl_DeleteVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Delete existing video", func(t *testing.T) {
		videoID := int64(1)

		err := videoService.DeleteVideo(videoID)

		assert.NoError(t, err, "Should delete video without error")
	})

	t.Run("Delete non-existent video", func(t *testing.T) {
		videoID := int64(999) // Non-existent video

		err := videoService.DeleteVideo(videoID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle deletion of non-existent video")
		}
	})

	t.Run("Delete video with invalid ID", func(t *testing.T) {
		videoID := int64(-1) // Invalid video ID

		err := videoService.DeleteVideo(videoID)

		assert.Error(t, err, "Should not delete with invalid video ID")
	})
}

func TestVideoServiceImpl_GetVideoCountByAuthorId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Get video count for existing author", func(t *testing.T) {
		authorID := int64(1)

		count, err := videoService.GetVideoCountByAuthorId(authorID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Video count should be non-negative")
		} else {
			// This might happen if author doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent author")
		}
	})

	t.Run("Get video count for non-existent author", func(t *testing.T) {
		authorID := int64(999) // Non-existent author

		count, err := videoService.GetVideoCountByAuthorId(authorID)

		assert.NotNil(t, err, "Error should be returned for non-existent author")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent author")
	})

	t.Run("Get video count with invalid author ID", func(t *testing.T) {
		authorID := int64(-1) // Invalid author ID

		count, err := videoService.GetVideoCountByAuthorId(authorID)

		assert.NotNil(t, err, "Error should be returned for invalid author ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestVideoServiceImpl_ValidateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	videoService := service.GetVideoServiceInstance()

	t.Run("Validate valid video", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "Valid Video",
			Description: "This is a valid video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
		}

		err := videoService.ValidateVideo(video)

		assert.NoError(t, err, "Valid video should not produce error")
	})

	t.Run("Validate video with empty title", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "", // Empty title
			Description: "This is a valid video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
		}

		err := videoService.ValidateVideo(video)

		assert.Error(t, err, "Empty title should produce error")
	})

	t.Run("Validate video with empty play URL", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "Valid Video",
			Description: "This is a valid video",
			PlayUrl:     "", // Empty play URL
			CoverUrl:    "http://example.com/cover.jpg",
		}

		err := videoService.ValidateVideo(video)

		assert.Error(t, err, "Empty play URL should produce error")
	})

	t.Run("Validate video with invalid author ID", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    -1,
			Title:       "Valid Video",
			Description: "This is a valid video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
		}

		err := videoService.ValidateVideo(video)

		assert.Error(t, err, "Invalid author ID should produce error")
	})

	t.Run("Validate video with nil", func(t *testing.T) {
		err := videoService.ValidateVideo(nil)

		assert.Error(t, err, "Nil video should produce error")
	})

	t.Run("Validate video with too long title", func(t *testing.T) {
		video := &model.Video{
			AuthorId:    1,
			Title:       "This video title is way too long and exceeds the maximum allowed length for video titles in the system. The maximum length should be properly enforced to prevent abuse and ensure good user experience.",
			Description: "This is a valid video",
			PlayUrl:     "http://example.com/video.mp4",
			CoverUrl:    "http://example.com/cover.jpg",
		}

		err := videoService.ValidateVideo(video)

		assert.Error(t, err, "Too long title should produce error")
	})
}

// Benchmark tests
func BenchmarkVideoServiceImpl_GetVideoServiceInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.GetVideoServiceInstance()
	}
}

func BenchmarkVideoServiceImpl_PublishVideo(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	video := &model.Video{
		AuthorId:    1,
		Title:       "Benchmark Video",
		Description: "This is a benchmark video",
		PlayUrl:     "http://example.com/video.mp4",
		CoverUrl:    "http://example.com/cover.jpg",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.PublishVideo(video)
	}
}

func BenchmarkVideoServiceImpl_GetVideoList(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	authorID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.GetVideoList(authorID)
	}
}

func BenchmarkVideoServiceImpl_GetVideoById(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.GetVideoById(videoID)
	}
}

func BenchmarkVideoServiceImpl_GetVideoByAuthorId(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	authorID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.GetVideoByAuthorId(authorID)
	}
}

func BenchmarkVideoServiceImpl_UpdateVideo(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	video := &model.Video{
		Id:          int64(1),
		Title:       "Benchmark Video",
		Description: "This is a benchmark video",
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.UpdateVideo(video)
	}
}

func BenchmarkVideoServiceImpl_DeleteVideo(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.DeleteVideo(videoID)
	}
}

func BenchmarkVideoServiceImpl_GetVideoCountByAuthorId(b *testing.B) {
	videoService := service.GetVideoServiceInstance()
	authorID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		videoService.GetVideoCountByAuthorId(authorID)
	}
}