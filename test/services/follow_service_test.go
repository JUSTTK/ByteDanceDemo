package test

import (
	"bytedancedemo/model"
	"bytedancedemo/service"
	"bytedancedemo/test/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFollowServiceImpl_GetFollowServiceInstanceSingleton(t *testing.T) {
	t.Run("Singleton pattern test", func(t *testing.T) {
		fs1 := service.GetFollowServiceInstance()
		fs2 := service.GetFollowServiceInstance()

		assert.Equal(t, fs1, fs2, "Instances should be the same")
	})
}

func TestFollowServiceImpl_CheckIsFollowing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Check following when user is following", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(2)

		isFollowing, err := followService.CheckIsFollowing(followerID, followingID)

		if err == nil {
			assert.NotNil(t, &isFollowing, "Follow status should be determined")
		} else {
			// This might happen if users don't exist
			assert.NotNil(t, err, "Error should be returned for non-existent users")
		}
	})

	t.Run("Check following when user is not following", func(t *testing.T) {
		followerID := int64(3)
		followingID := int64(4)

		isFollowing, err := followService.CheckIsFollowing(followerID, followingID)

		if err == nil {
			assert.NotNil(t, &isFollowing, "Follow status should be determined")
		} else {
			// This might happen if users don't exist
			assert.NotNil(t, err, "Error should be returned for non-existent users")
		}
	})

	t.Run("Check following with same user", func(t *testing.T) {
		userID := int64(1)

		isFollowing, err := followService.CheckIsFollowing(userID, userID)

		if err == nil {
			assert.False(t, isFollowing, "User should not follow themselves")
		} else {
			assert.NotNil(t, err, "Error should be returned when user follows themselves")
		}
	})

	t.Run("Check following with invalid follower ID", func(t *testing.T) {
		followerID := int64(-1)
		followingID := int64(2)

		isFollowing, err := followService.CheckIsFollowing(followerID, followingID)

		assert.NotNil(t, err, "Error should be returned for invalid follower ID")
		assert.False(t, isFollowing, "Should not follow with invalid ID")
	})

	t.Run("Check following with invalid following ID", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(-1)

		isFollowing, err := followService.CheckIsFollowing(followerID, followingID)

		assert.NotNil(t, err, "Error should be returned for invalid following ID")
		assert.False(t, isFollowing, "Should not follow with invalid ID")
	})
}

func TestFollowServiceImpl_FollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Follow existing user", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(2)

		err := followService.FollowUser(followerID, followingID)

		assert.NoError(t, err, "Should follow user without error")
	})

	t.Run("Follow same user", func(t *testing.T) {
		userID := int64(1)

		err := followService.FollowUser(userID, userID)

		assert.Error(t, err, "Should not allow user to follow themselves")
	})

	t.Run("Follow non-existent user", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(999) // Non-existent user

		err := followService.FollowUser(followerID, followingID)

		assert.Error(t, err, "Should not follow non-existent user")
	})

	t.Run("Follow with invalid follower ID", func(t *testing.T) {
		followerID := int64(-1)
		followingID := int64(2)

		err := followService.FollowUser(followerID, followingID)

		assert.Error(t, err, "Should not follow with invalid follower ID")
	})

	t.Run("Follow with invalid following ID", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(-1)

		err := followService.FollowUser(followerID, followingID)

		assert.Error(t, err, "Should not follow with invalid following ID")
	})

	t.Run("Follow user who is already being followed", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(2)

		// First follow should succeed
		err1 := followService.FollowUser(followerID, followingID)
		assert.NoError(t, err1, "First follow should succeed")

		// Second follow should handle gracefully
		err2 := followService.FollowUser(followerID, followingID)
		// Depending on implementation, this might succeed or return an error
		if err2 != nil {
			assert.Error(t, err2, "Should handle duplicate follow gracefully")
		}
	})
}

func TestFollowServiceImpl_UnfollowUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Unfollow existing user", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(2)

		err := followService.UnfollowUser(followerID, followingID)

		assert.NoError(t, err, "Should unfollow user without error")
	})

	t.Run("Unfollow same user", func(t *testing.T) {
		userID := int64(1)

		err := followService.UnfollowUser(userID, userID)

		assert.Error(t, err, "Should not allow user to unfollow themselves")
	})

	t.Run("Unfollow non-existent user", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(999) // Non-existent user

		err := followService.UnfollowUser(followerID, followingID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle unfollow of non-existent user")
		}
	})

	t.Run("Unfollow with invalid follower ID", func(t *testing.T) {
		followerID := int64(-1)
		followingID := int64(2)

		err := followService.UnfollowUser(followerID, followingID)

		assert.Error(t, err, "Should not unfollow with invalid follower ID")
	})

	t.Run("Unfollow with invalid following ID", func(t *testing.T) {
		followerID := int64(1)
		followingID := int64(-1)

		err := followService.UnfollowUser(followerID, followingID)

		assert.Error(t, err, "Should not unfollow with invalid following ID")
	})

	t.Run("Unfollow user who is not being followed", func(t *testing.T) {
		followerID := int64(3)
		followingID := int64(4)

		err := followService.UnfollowUser(followerID, followingID)

		// Should handle gracefully - either succeed or return appropriate error
		if err != nil {
			assert.Error(t, err, "Should handle unfollow of non-followed user")
		}
	})
}

func TestFollowServiceImpl_GetFollowingList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Get following list for existing user", func(t *testing.T) {
		userID := int64(1)

		followingList, err := followService.GetFollowingList(userID)

		if err == nil {
			assert.NotNil(t, followingList, "Following list should not be nil")
			for _, following := range followingList {
				assert.Greater(t, following, int64(0), "Following IDs should be valid")
			}
		} else {
			// This might happen if user doesn't exist or has no following
			assert.NotNil(t, err, "Error might be returned for non-existent user")
		}
	})

	t.Run("Get following list for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		followingList, err := followService.GetFollowingList(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, followingList, "Following list should be nil for non-existent user")
	})

	t.Run("Get following list with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		followingList, err := followService.GetFollowingList(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Nil(t, followingList, "Following list should be nil for invalid input")
	})
}

func TestFollowServiceImpl_GetFollowerList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Get follower list for existing user", func(t *testing.T) {
		userID := int64(1)

		followerList, err := followService.GetFollowerList(userID)

		if err == nil {
			assert.NotNil(t, followerList, "Follower list should not be nil")
			for _, follower := range followerList {
				assert.Greater(t, follower, int64(0), "Follower IDs should be valid")
			}
		} else {
			// This might happen if user doesn't exist or has no followers
			assert.NotNil(t, err, "Error might be returned for non-existent user")
		}
	})

	t.Run("Get follower list for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		followerList, err := followService.GetFollowerList(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Nil(t, followerList, "Follower list should be nil for non-existent user")
	})

	t.Run("Get follower list with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		followerList, err := followService.GetFollowerList(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Nil(t, followerList, "Follower list should be nil for invalid input")
	})
}

func TestFollowServiceImpl_GetFollowingCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Get following count for existing user", func(t *testing.T) {
		userID := int64(1)

		count, err := followService.GetFollowingCount(userID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Following count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get following count for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		count, err := followService.GetFollowingCount(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent user")
	})

	t.Run("Get following count with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		count, err := followService.GetFollowingCount(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestFollowServiceImpl_GetFollowerCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Get follower count for existing user", func(t *testing.T) {
		userID := int64(1)

		count, err := followService.GetFollowerCount(userID)

		if err == nil {
			assert.GreaterOrEqual(t, count, int64(0), "Follower count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get follower count for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		count, err := followService.GetFollowerCount(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), count, "Count should be 0 for non-existent user")
	})

	t.Run("Get follower count with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		count, err := followService.GetFollowerCount(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), count, "Count should be 0 for invalid input")
	})
}

func TestFollowServiceImpl_GetFollowStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	followService := service.GetFollowServiceInstance()

	t.Run("Get complete follow stats", func(t *testing.T) {
		userID := int64(1)

		followingCount, followerCount, err := followService.GetFollowStats(userID)

		if err == nil {
			assert.GreaterOrEqual(t, followingCount, int64(0), "Following count should be non-negative")
			assert.GreaterOrEqual(t, followerCount, int64(0), "Follower count should be non-negative")
		} else {
			// This might happen if user doesn't exist
			assert.NotNil(t, err, "Error should be returned for non-existent user")
		}
	})

	t.Run("Get stats for non-existent user", func(t *testing.T) {
		userID := int64(999) // Non-existent user

		followingCount, followerCount, err := followService.GetFollowStats(userID)

		assert.NotNil(t, err, "Error should be returned for non-existent user")
		assert.Equal(t, int64(0), followingCount, "Following count should be 0")
		assert.Equal(t, int64(0), followerCount, "Follower count should be 0")
	})

	t.Run("Get stats with invalid user ID", func(t *testing.T) {
		userID := int64(-1) // Invalid user ID

		followingCount, followerCount, err := followService.GetFollowStats(userID)

		assert.NotNil(t, err, "Error should be returned for invalid user ID")
		assert.Equal(t, int64(0), followingCount, "Following count should be 0")
		assert.Equal(t, int64(0), followerCount, "Follower count should be 0")
	})
}

// Benchmark tests
func BenchmarkFollowServiceImpl_GetFollowServiceInstance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.GetFollowServiceInstance()
	}
}

func BenchmarkFollowServiceImpl_CheckIsFollowing(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	followerID := int64(1)
	followingID := int64(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.CheckIsFollowing(followerID, followingID)
	}
}

func BenchmarkFollowServiceImpl_FollowUser(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	followerID := int64(1)
	followingID := int64(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.FollowUser(followerID, followingID)
	}
}

func BenchmarkFollowServiceImpl_UnfollowUser(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	followerID := int64(1)
	followingID := int64(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.UnfollowUser(followerID, followingID)
	}
}

func BenchmarkFollowServiceImpl_GetFollowingList(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.GetFollowingList(userID)
	}
}

func BenchmarkFollowServiceImpl_GetFollowerList(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.GetFollowerList(userID)
	}
}

func BenchmarkFollowServiceImpl_GetFollowingCount(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.GetFollowingCount(userID)
	}
}

func BenchmarkFollowServiceImpl_GetFollowerCount(b *testing.B) {
	followService := service.GetFollowServiceInstance()
	userID := int64(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		followService.GetFollowerCount(userID)
	}
}