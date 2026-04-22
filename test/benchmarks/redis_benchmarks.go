package test

import (
	"context"
	"fmt"
	"time"

	"bytedancedemo/database/redis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupBenchmarkRedis() *redis.Client {
	// Initialize Redis connection
	config := redis.Config{
		Addr:     "localhost:6379",
		Password: "", // No password
		DB:       0,  // Use default DB
	}

	db, err := redis.NewRedis(&config)
	if err != nil {
		panic(err)
	}

	// Clear any existing test data
	db.FlushDB(context.Background())

	return db
}

func BenchmarkRedisSetGet(b *testing.B) {
	db := setupBenchmarkRedis()

	key := "test_key"
	value := "test_value"

	b.Run("Single Set/Get operation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Set value
			err := db.Set(context.Background(), key, value, 0).Err()
			assert.NoError(b, err)

			// Get value
			val, err := db.Get(context.Background(), key).Result()
			assert.NoError(b, err)
			assert.Equal(b, value, val)
		}
	})

	b.Run("Multiple concurrent Set/Get operations", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				key := fmt.Sprintf("key_%d", i)
				value := fmt.Sprintf("value_%d", i)

				// Set value
				err := db.Set(context.Background(), key, value, 0).Err()
				assert.NoError(b, err)

				// Get value
				val, err := db.Get(context.Background(), key).Result()
				assert.NoError(b, err)
				assert.Equal(b, value, val)
			}(i)
		}

		wg.Wait()
	})

	b.Run("Pipeline Set/Get operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pipe := db.Pipeline()

			// Multiple sets in pipeline
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("pipeline_key_%d_%d", i, j)
				value := fmt.Sprintf("pipeline_value_%d_%d", i, j)
				pipe.Set(context.Background(), key, value, 0)
			}

			// Execute pipeline
			_, err := pipe.Exec(context.Background())
			assert.NoError(b, err)

			// Get values
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("pipeline_key_%d_%d", i, j)
				val, err := db.Get(context.Background(), key).Result()
				assert.NoError(b, err)
				assert.Contains(b, val, fmt.Sprintf("pipeline_value_%d_%d", i, j))
			}
		}
	})
}

func BenchmarkRedisCacheOperations(b *testing.B) {
	db := setupBenchmarkRedis()

	// Simulate user data
	userID := int64(1)
	userData := map[string]interface{}{
		"id":       userID,
		"name":     "John Doe",
		"email":    "john@example.com",
		"age":      30,
		"created":  time.Now().Unix(),
	}

	cacheKey := fmt.Sprintf("user:%d", userID)

	b.Run("Cache user data", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Serialize user data to JSON
			data, err := json.Marshal(userData)
			assert.NoError(b, err)

			// Set to cache
			err = db.Set(context.Background(), cacheKey, data, 10*time.Minute).Err()
			assert.NoError(b, err)

			// Get from cache
			cachedData, err := db.Get(context.Background(), cacheKey).Result()
			assert.NoError(b, err)

			// Deserialize
			var cachedUser map[string]interface{}
			err = json.Unmarshal([]byte(cachedData), &cachedUser)
			assert.NoError(b, err)
			assert.Equal(b, userID, int64(cachedUser["id"].(float64)))
		}
	})

	b.Run("Cache with expiration", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			expiration := 30 * time.Second

			// Set with expiration
			data, err := json.Marshal(userData)
			assert.NoError(b, err)

			err = db.Set(context.Background(), cacheKey, data, expiration).Err()
			assert.NoError(b, err)

			// Get TTL
			ttl, err := db.TTL(context.Background(), cacheKey).Result()
			assert.NoError(b, err)
			assert.GreaterOrEqual(b, ttl, 29*time.Second) // Should be close to 30 seconds
		}
	})

	b.Run("Cache hit vs miss", func(b *testing.B) {
		b.Run("Cache hit", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Ensure data is in cache
				data, err := json.Marshal(userData)
				assert.NoError(b, err)

				err = db.Set(context.Background(), cacheKey, data, 10*time.Minute).Err()
				assert.NoError(b, err)

				// Get from cache
				_, err = db.Get(context.Background(), cacheKey).Result()
				assert.NoError(b, err)
			}
		})

		b.Run("Cache miss", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				nonExistentKey := fmt.Sprintf("nonexistent:%d", i)

				// Get from cache (should miss)
				_, err := db.Get(context.Background(), nonExistentKey).Result()
				assert.ErrorIs(b, redis.Nil, err)
			}
		})
	})
}

func BenchmarkRedisListOperations(b *testing.B) {
	db := setupBenchmarkRedis()

	listKey := "test_list"
	items := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		items[i] = fmt.Sprintf("item_%d", i)
	}

	b.Run("List push/pop operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear list
			db.Del(context.Background(), listKey)

			// Push items to list
			for _, item := range items {
				db.RPush(context.Background(), listKey, item)
			}

			// Pop items from list
			for j := 0; j < len(items); j++ {
				_, err := db.LPop(context.Background(), listKey).Result()
				assert.NoError(b, err)
			}
		}
	})

	b.Run("List range operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear list
			db.Del(context.Background(), listKey)

			// Push items to list
			for _, item := range items {
				db.RPush(context.Background(), listKey, item)
			}

			// Get range of items
			result, err := db.LRange(context.Background(), listKey, 0, 10).Result()
			assert.NoError(b, err)
			assert.Equal(b, 11, len(result))
		}
	})

	b.Run("List length operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear list
			db.Del(context.Background(), listKey)

			// Push items to list
			for _, item := range items {
				db.RPush(context.Background(), listKey, item)
			}

			// Get list length
			length, err := db.LLen(context.Background(), listKey).Result()
			assert.NoError(b, err)
			assert.Equal(b, 100, int(length))
		}
	})
}

func BenchmarkRedisHashOperations(b *testing.B) {
	db := setupBenchmarkRedis()

	hashKey := "user_hash"
	fields := make(map[string]interface{})
	for i := 0; i < 50; i++ {
		fields[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	b.Run("Hash set/get operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear hash
			db.Del(context.Background(), hashKey)

			// Set fields
			for field, value := range fields {
				err := db.HSet(context.Background(), hashKey, field, value).Err()
				assert.NoError(b, err)
			}

			// Get fields
			for field := range fields {
				_, err := db.HGet(context.Background(), hashKey, field).Result()
				assert.NoError(b, err)
			}
		}
	})

	b.Run("Hash get all operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear hash
			db.Del(context.Background(), hashKey)

			// Set fields
			for field, value := range fields {
				db.HSet(context.Background(), hashKey, field, value)
			}

			// Get all fields
			allFields, err := db.HGetAll(context.Background(), hashKey).Result()
			assert.NoError(b, err)
			assert.Equal(b, 50, len(allFields))
		}
	})

	b.Run("Hash increment operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear hash
			db.Del(context.Background(), hashKey)

			// Increment counter
			for j := 0; j < 10; j++ {
				_, err := db.HIncrBy(context.Background(), hashKey, "counter", 1).Result()
				assert.NoError(b, err)
			}

			// Get counter value
			value, err := db.HGet(context.Background(), hashKey, "counter").Result()
			assert.NoError(b, err)
			assert.Equal(b, "10", value)
		}
	})
}

func BenchmarkRedisSetOperations(b *testing.B) {
	db := setupBenchmarkRedis()

	setKey := "test_set"
	members := make([]interface{}, 50)
	for i := 0; i < 50; i++ {
		members[i] = fmt.Sprintf("member_%d", i)
	}

	b.Run("Set add/remove operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear set
			db.Del(context.Background(), setKey)

			// Add members to set
			for _, member := range members {
				db.SAdd(context.Background(), setKey, member)
			}

			// Remove members from set
			for j := 0; j < 10; j++ {
				db.SRem(context.Background(), setKey, fmt.Sprintf("member_%d", j))
			}
		}
	})

	b.Run("Set membership operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear set
			db.Del(context.Background(), setKey)

			// Add members to set
			for _, member := range members {
				db.SAdd(context.Background(), setKey, member)
			}

			// Check membership
			for j := 0; j < 50; j++ {
				exists, err := db.SIsMember(context.Background(), setKey, fmt.Sprintf("member_%d", j)).Result()
				assert.NoError(b, err)
				assert.True(b, exists)
			}
		}
	})

	b.Run("Set size operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear set
			db.Del(context.Background(), setKey)

			// Add members to set
			for _, member := range members {
				db.SAdd(context.Background(), setKey, member)
			}

			// Get set size
			size, err := db.SCard(context.Background(), setKey).Result()
			assert.NoError(b, err)
			assert.Equal(b, 50, int(size))
		}
	})
}

func BenchmarkRedisZSetOperations(b *testing.B) {
	db := setupBenchmarkRedis()

	zsetKey := "test_zset"
	members := make([]redis.Z, 50)
	for i := 0; i < 50; i++ {
		members[i] = redis.Z{
			Score:  float64(i),
			Member: fmt.Sprintf("member_%d", i),
		}
	}

	b.Run("ZSet add operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear zset
			db.Del(context.Background(), zsetKey)

			// Add members to zset
			for _, member := range members {
				db.ZAdd(context.Background(), zsetKey, member)
			}
		}
	})

	b.Run("ZSet range operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear zset
			db.Del(context.Background(), zsetKey)

			// Add members to zset
			for _, member := range members {
				db.ZAdd(context.Background(), zsetKey, member)
			}

			// Get range
			result, err := db.ZRange(context.Background(), zsetKey, 0, 10).Result()
			assert.NoError(b, err)
			assert.Equal(b, 11, len(result))
		}
	})

	b.Run("ZSet rank operations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Clear zset
			db.Del(context.Background(), zsetKey)

			// Add members to zset
			for _, member := range members {
				db.ZAdd(context.Background(), zsetKey, member)
			}

			// Get rank
			for j := 0; j < 50; j++ {
				rank, err := db.ZRank(context.Background(), zsetKey, fmt.Sprintf("member_%d", j)).Result()
				assert.NoError(b, err)
				assert.Equal(b, int64(j), rank)
			}
		}
	})
}

func BenchmarkRedisPubSub(b *testing.B) {
	db := setupBenchmarkRedis()

	channel := "test_channel"
	message := "test_message"

	b.Run("Publish/Subscribe operations", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Subscribe to channel
		sub := db.Subscribe(ctx, channel)
		ch := sub.Channel()

		// Goroutine to receive messages
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range ch {
				// Process message
				_ = msg
			}
		}()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Publish message
			err := db.Publish(ctx, channel, fmt.Sprintf("%s_%d", message, i)).Err()
			assert.NoError(b, err)
		}

		// Wait for all messages to be processed
		time.Sleep(100 * time.Millisecond)
		cancel()
		wg.Wait()
	})
}

func BenchmarkRedisConnectionPool(b *testing.B) {
	db := setupBenchmarkRedis()

	b.Run("Connection stress test", func(b *testing.B) {
		var wg sync.WaitGroup

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				// Perform operations with different keys
				key := fmt.Sprintf("stress_key_%d", i)
				value := fmt.Sprintf("stress_value_%d", i)

				err := db.Set(context.Background(), key, value, 0).Err()
				assert.NoError(b, err)

				val, err := db.Get(context.Background(), key).Result()
				assert.NoError(b, err)
				assert.Equal(b, value, val)
			}(i)
		}

		wg.Wait()
	})
}

func BenchmarkRedisMemoryUsage(b *testing.B) {
	db := setupBenchmarkRedis()

	b.Run("Memory usage test", func(b *testing.B) {
		// Clear database
		db.FlushDB(context.Background())

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Store 1000 keys
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("mem_key_%d_%d", i, j)
				value := fmt.Sprintf("mem_value_%d_%d", i, j)
				db.Set(context.Background(), key, value, 0)
			}

			// Get memory info
			info, err := db.Info(context.Background(), "memory").Result()
			assert.NoError(b, err)
			_ = info
		}
	})
}