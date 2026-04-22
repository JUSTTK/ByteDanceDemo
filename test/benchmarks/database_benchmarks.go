package test

import (
	"bytedancedemo/config"
	"bytedancedemo/database/mysql"
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"bytedancedemo/service"
	"math/rand"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
)

func setupBenchmarkDB() {
	config.Init("../config/settings.yml")
	mysql.Init()
	dao.SetDefault(mysql.DB)
}

func TestBenchmarkUserCreation(b *testing.B) {
	setupBenchmarkDB()

	b.Run("Sequential user creation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			user := &model.User{
				Name:     fmt.Sprintf("user%d", i),
				Password: "hashed_password",
				Email:    fmt.Sprintf("user%d@example.com", i),
			}

			userService := service.GetUserServiceInstance()
			_, success := userService.InsertUser(user)

			if !success {
				b.Errorf("Failed to create user %d", i)
			}
		}
	})

	b.Run("Concurrent user creation", func(b *testing.B) {
		var wg sync.WaitGroup
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				user := &model.User{
					Name:     fmt.Sprintf("user%d", i),
					Password: "hashed_password",
					Email:    fmt.Sprintf("user%d@example.com", i),
				}

				userService := service.GetUserServiceInstance()
				_, success := userService.InsertUser(user)

				if !success {
					b.Errorf("Failed to create user %d", i)
				}
			}(i)
		}

		wg.Wait()
	})

	b.Run("Batch user creation", func(b *testing.B) {
		users := make([]*model.User, 100)
		for i := 0; i < 100; i++ {
			users[i] = &model.User{
				Name:     fmt.Sprintf("user%d", i),
				Password: "hashed_password",
				Email:    fmt.Sprintf("user%d@example.com", i),
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N/100; i++ {
			var wg sync.WaitGroup

			for _, user := range users {
				wg.Add(1)
				go func(u *model.User) {
					defer wg.Done()
					userService := service.GetUserServiceInstance()
					_, _ = userService.InsertUser(u)
				}(user)
			}

			wg.Wait()
		}
	})
}

func TestBenchmarkUserQueries(b *testing.B) {
	setupBenchmarkDB()

	// Create test users first
	createTestUsers(1000)

	b.Run("Get user by ID", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		userID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = userService.GetUserDetailsById(userID, nil)
		}
	})

	b.Run("Get user by username", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		username := "user1"

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = userService.GetUserBasicByPassword(username, "hashed_password")
		}
	})

	b.Run("Get user name", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		userID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = userService.GetUserName(userID)
		}
	})

	b.Run("Concurrent user queries", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(id int64) {
				defer wg.Done()
				_, _ = userService.GetUserDetailsById(id, nil)
			}(int64(rand.Intn(1000) + 1))
		}

		wg.Wait()
	})
}

func TestBenchmarkUserUpdates(b *testing.B) {
	setupBenchmarkDB()

	// Create test users first
	createTestUsers(100)

	b.Run("Update user profile", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		userID := int64(1)

		newSignature := "Updated signature"
		newAvatar := "new_avatar.jpg"

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// This would require implementing update functionality in the service
			// For now, we just benchmark the existing methods
			user, err := userService.GetUserDetailsById(userID, nil)
			if err == nil {
				user.Signature = newSignature + fmt.Sprintf("_%d", i)
				user.Avatar = newAvatar
			}
		}
	})

	b.Run("Concurrent user updates", func(b *testing.B) {
		userService := service.GetUserServiceInstance()
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(id int64) {
				defer wg.Done()
				user, _ := userService.GetUserDetailsById(id, nil)
				if user != nil {
					user.Signature = "Updated signature"
				}
			}(int64(rand.Intn(100) + 1))
		}

		wg.Wait()
	})
}

func TestBenchmarkUserDeletions(b *testing.B) {
	setupBenchmarkDB()

	b.Run("Delete user", func(b *testing.B) {
		// Note: This test would require implementing user deletion in the service
		// For now, it's a placeholder
		userService := service.GetUserServiceInstance()
		userID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// This would be userService.DeleteUser(userID)
			_, _ = userService.GetUserDetailsById(userID, nil)
		}
	})
}

func TestBenchmarkVideoOperations(b *testing.B) {
	setupBenchmarkDB()

	// Create test users first
	createTestUsers(10)

	b.Run("Publish video", func(b *testing.B) {
		videoService := service.GetVideoServiceInstance()
		authorID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			video := &model.Video{
				AuthorId:    authorID,
				Title:       fmt.Sprintf("Video %d", i),
				Description: fmt.Sprintf("Description for video %d", i),
				PlayUrl:     fmt.Sprintf("http://example.com/video%d.mp4", i),
				CoverUrl:    fmt.Sprintf("http://example.com/cover%d.jpg", i),
			}

			_, _ = videoService.PublishVideo(video)
		}
	})

	b.Run("Get video list", func(b *testing.B) {
		videoService := service.GetVideoServiceInstance()
		authorID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = videoService.GetVideoList(authorID)
		}
	})

	b.Run("Get video by ID", func(b *testing.B) {
		videoService := service.GetVideoServiceInstance()
		videoID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = videoService.GetVideoById(videoID)
		}
	})

	b.Run("Get video count", func(b *testing.B) {
		videoService := service.GetVideoServiceInstance()
		authorID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = videoService.GetVideoCountByAuthorId(authorID)
		}
	})
}

func TestBenchmarkFollowOperations(b *testing.B) {
	setupBenchmarkDB()

	// Create test users first
	createTestUsers(100)

	b.Run("Follow user", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		followerID := int64(1)
		followingID := int64(2)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = followService.FollowUser(followerID, followingID)
		}
	})

	b.Run("Unfollow user", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		followerID := int64(1)
		followingID := int64(2)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = followService.UnfollowUser(followerID, followingID)
		}
	})

	b.Run("Check if following", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		followerID := int64(1)
		followingID := int64(2)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = followService.CheckIsFollowing(followerID, followingID)
		}
	})

	b.Run("Get following list", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		userID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = followService.GetFollowingList(userID)
		}
	})

	b.Run("Get follower list", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		userID := int64(1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = followService.GetFollowerList(userID)
		}
	})

	b.Run("Concurrent follow operations", func(b *testing.B) {
		followService := service.GetFollowServiceInstance()
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				followerID := int64(rand.Intn(100) + 1)
				followingID := int64(rand.Intn(100) + 1)
				_, _ = followService.FollowUser(followerID, followingID)
			}()
		}

		wg.Wait()
	})
}

func TestBenchmarkDatabaseConnections(b *testing.B) {
	setupBenchmarkDB()

	b.Run("Database connection pool performance", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Simulate database query
				userService := service.GetUserServiceInstance()
				_, _ = userService.GetUserDetailsById(int64(1), nil)
			}()
		}

		wg.Wait()
	})

	b.Run("Connection stress test", func(b *testing.B) {
		var wg sync.WaitGroup
		concurrentGoroutines := 100

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < concurrentGoroutines; j++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()

					// Simulate different operations
					switch j % 5 {
					case 0:
						userService := service.GetUserServiceInstance()
						_, _ = userService.GetUserDetailsById(int64(1), nil)
					case 1:
						followService := service.GetFollowServiceInstance()
						_, _ = followService.GetFollowingList(int64(1))
					case 2:
						videoService := service.GetVideoServiceInstance()
						_, _ = videoService.GetVideoList(int64(1))
					case 3:
						commentService := service.GetCommentServiceInstance()
						_, _ = commentService.GetCommentList(int64(1))
					case 4:
						messageService := service.GetMessageServiceInstance()
						_, _ = messageService.GetMessageList(int64(1), int64(2))
					}
				}(j)
			}
		}

		wg.Wait()
	})
}

func createTestUsers(count int) {
	userService := service.GetUserServiceInstance()

	for i := 0; i < count; i++ {
		user := &model.User{
			Name:     fmt.Sprintf("user%d", i),
			Password: "hashed_password",
			Email:    fmt.Sprintf("user%d@example.com", i),
		}

		_, _ = userService.InsertUser(user)
	}
}

func cleanupDatabase() {
	setupBenchmarkDB()

	// Clean up test data
	db := mysql.DB
	db.Exec("DELETE FROM users WHERE name LIKE 'user%'")
	db.Exec("DELETE FROM videos WHERE title LIKE 'Video %'")
	db.Exec("DELETE FROM comments WHERE content LIKE 'Test comment %'")
	db.Exec("DELETE FROM favorites WHERE comment_id IN (SELECT id FROM comments WHERE content LIKE 'Test comment %')")
	db.Exec("DELETE FROM relations WHERE follower_id IN (SELECT id FROM users WHERE name LIKE 'user%')")
}

// Benchmark cleanup
func BenchmarkDatabaseCleanup(b *testing.B) {
	setupBenchmarkDB()

	for i := 0; i < b.N; i++ {
		cleanupDatabase()
	}
}

func init() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
}