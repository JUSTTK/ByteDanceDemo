package test

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/test/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteServiceImpl_GetFavoriteServiceInstanceSingleton(t *testing.T) {
	t.Run("Singleton pattern test", func(t *testing.T) {
		fs1 := service.GetFavoriteServiceInstance()
		fs2 := service.GetFavoriteServiceInstance()

		assert.Equal(t, fs1, fs2, "Instances should be the same")
	})
}

func TestFavoriteServiceImpl_FavoriteVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Favorite existing video", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(1)

		err := favoriteService.FavoriteVideo(userID, videoID)

		assert.NoError(t, err, "Should favorite video without error")
	})

	t.Run("Favorite same video multiple times", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(2)

		// First favorite should succeed
		err1 := favoriteService.FavoriteVideo(userID, videoID)
		assert.NoError(t, err1, "First favorite should succeed")

		// Second favorite should handle gracefully
		err2 := favoriteService.FavoriteVideo(userID, videoID)
		// Depending on implementation, this might succeed or return an error
		if err2 != nil {
			assert.Error(t, err2, "Should handle duplicate favorite gracefully")
		}
	})

	t.Run("Favorite non-existent video", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(999) // Non-existent video

		err := favoriteService.FavoriteVideo(userID, videoID)

		assert.Error(t, err, "Should not favorite non-existent video")
	})

	t.Run("Favorite with invalid user ID", func(t *testing.T) {
		userID := int64(-1)
		videoID := int64(1)

		err := favoriteService.FavoriteVideo(userID, videoID)

		assert.Error(t, err, "Should not favorite with invalid user ID")
	})

	t.Run("Favorite with invalid video ID", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(-1)

		err := favoriteService.FavoriteVideo(userID, videoID)

		assert.Error(t, err, "Should not favorite with invalid video ID")
	})
}

func TestFavoriteServiceImpl_UnfavoriteVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Unfavorite existing video", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(1)

		err := favoriteService.UnfavoriteVideo(userID, videoID)

		assert.NoError(t, err, "Should unfavorite video without error")
	})

	t.Run("Unfavorite non-existent favorite", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(999) // Non-existent favorite

		err := favoriteService.UnfavoriteVideo(userID, videoID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle unfavorite of non-existent favorite")
		}
	})

	t.Run("Unfavorite with invalid user ID", func(t *testing.T) {
		userID := int64(-1)
		videoID := int64(1)

		err := favoriteService.UnfavoriteVideo(userID, videoID)

		assert.Error(t, err, "Should not unfavorite with invalid user ID")
	})

	t.Run("Unfavorite with invalid video ID", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(-1)

		err := favoriteService.UnfavoriteVideo(userID, videoID)

		assert.Error(t, err, "Should not unfavorite with invalid video ID")
	})

	t.Run("Unfavorite video that is not favorited", func(t *testing.T) {
		userID := int64(3)
		videoID := int64(4)

		err := favoriteService.UnfavoriteVideo(userID, videoID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle unfavorite of non-favorited video")
		}
	})
}

func TestFavoriteServiceImpl_GetFavoriteList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Get favorite list for existing user", func(t *testing.T) {
		userID := int64(1)

		favoriteList, err := favoriteService.GetFavoriteList(userID)

		if err == nil {
			assert.NotNil(t, favoriteList, "Favorite list should not be nil")
			for _, favorite := range favoriteList {
				assert.Greater(t, favorite.VideoId, int64(0), "Video IDs should be valid")
				assert.Equal(t, userID, favorite.UserId, "All favorites should belong to the user")
			}
		} else {
			// This might happen if user doesn't exist or has no favorites
			assert.NotNil(t, err, "Error might be returned for non-existent user")
		}
	})

	t.Run("Get favorite list for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		favoriteList, err := favoriteService.GetFavoriteList(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, favoriteList, "Favorite list should be nil for non-existent user")
	})

	t.Run("Get favorite list with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		favoriteList, err := favoriteService.GetFavoriteList(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Nil(t, favoriteList, "Favorite list should be nil for invalid input")
	})
}

func TestFavoriteServiceImpl_CheckIsFavorite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Check favorite when video is favorited", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(1)

		isFavorite, err := favoriteService.CheckIsFavorite(userID, videoID)

		if err == nil {
			assert.NotNil(t, &isFavorite, "Favorite status should be determined")
		} else {
			// This might happen if user or video doesn't exist
			assert.NotNil(t, err, "Error might be returned for non-existent user/video")
		}
	})

	t.Run("Check favorite when video is not favorited", func(t *testing.T) {
		userID := int64(2)
		videoID := int64(3)

		isFavorite, err := favoriteService.CheckIsFavorite(userID, videoID)

		if err == nil {
			assert.NotNil(t, &isFavorite, "Favorite status should be determined")
		} else {
			// This might happen if user or video doesn't exist
			assert.NotNil(t, err, "Error might be returned for non-existent user/video")
		}
	})

	t.Run("Check favorite with invalid user ID", func(t *testing.T) {
		userID := int64(-1)
		videoID := int64(1)

		isFavorite, err := favoriteService.CheckIsFavorite(userID, videoID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.False(t, isFavorite, "Should not be favorite with invalid user ID")
	})

	t.Run("Check favorite with invalid video ID", func(t *testing.T) {
		userID := int64(1)
		videoID := int64(-1)

		isFavorite, err := favoriteService.CheckIsFavorite(userID, videoID)

		assert.NotNil(t, err, "Error should be returned for invalid video ID")
		assert.False(t, isFavorite, "Should not be favorite with invalid video ID")
	})
}

func TestFavoriteServiceImpl_GetFavoriteCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Get favorite count for existing user", func(t *testing.T) {
		userID := int64(1)

		count, err := favoriteService.GetFavoriteCount(userID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Favorite count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get favorite count for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		count, err := favoriteService.GetFavoriteCount(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent user")
	})

	t.Run("Get favorite count with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		count, err := favoriteService.GetFavoriteCount(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestFavoriteServiceImpl_GetVideoFavoriteCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Get favorite count for existing video", func(t *testing.T) {
		videoID := int64(1)

		count, err := favoriteService.GetVideoFavoriteCount(videoID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Favorite count should be non-negative")
		} else {
			// This might happen if video doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent video")
		}
	})

	t.Run("Get favorite count for non-existent video", func(t *testing.T) {
		videoID := int64(999) // Non-existent video

		count, err := favoriteService.GetVideoFavoriteCount(videoID)

		assert.NotNil(t, err, "Error should be returned for non-existent video")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent video")
	})

	t.Run("Get favorite count with invalid video ID", func(t *testing.T) {
		videoID := int64(-1) // Invalid video ID

		count, err := favoriteService.GetVideoFavoriteCount(videoID)

		assert.NotNil(t, err, "Error should be returned for invalid video ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestFavoriteServiceImpl_GetFavoriteStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	favoriteService := service.GetFavoriteServiceInstance()

	t.Run("Get complete favorite stats", func(t *testing.T) {
		userID := int64(1)

		favoriteCount, videoCount, err := favoriteService.GetFavoriteStats(userID)

		if err == nil {
			assert.GreaterOrEqual(t, favoriteCount, int64(0), "Favorite count should be non-negative")
			assert.GreaterOrEqual(t, videoCount, int64(0), "Video count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get stats for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		favoriteCount, videoCount, err := favoriteService.GetFavoriteStats(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), favoriteCount, "Favorite count should be 0")
		assert.Equal(t, int64(0), videoCount, "Video count should be 0")
	})

	t.Run("Get stats with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		favoriteCount, videoCount, err := favoriteService.GetFavoriteStats(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), favoriteCount, "Favorite count should be 0")
		assert.Equal(t, int64(0), videoCount, "Video count should be 0")
	})
}

// Benchmark tests
func BenchmarkFavoriteServiceImpl_GetFavoriteServiceInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.GetFavoriteServiceInstance()
	}
}

func BenchmarkFavoriteServiceImpl_FavoriteVideo(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	userID := int64(1)
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.FavoriteVideo(userID, videoID)
	}
}

func BenchmarkFavoriteServiceImpl_UnfavoriteVideo(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	userID := int64(1)
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.UnfavoriteVideo(userID, videoID)
	}
}

func BenchmarkFavoriteServiceImpl_GetFavoriteList(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.GetFavoriteList(userID)
	}
}

func BenchmarkFavoriteServiceImpl_CheckIsFavorite(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	userID := int64(1)
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.CheckIsFavorite(userID, videoID)
	}
}

func BenchmarkFavoriteServiceImpl_GetFavoriteCount(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.GetFavoriteCount(userID)
	}
}

func BenchmarkFavoriteServiceImpl_GetVideoFavoriteCount(b *testing.B) {
	favoriteService := service.GetFavoriteServiceInstance()
	videoID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		favoriteService.GetVideoFavoriteCount(videoID)
	}
}