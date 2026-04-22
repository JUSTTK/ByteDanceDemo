package mocks

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockRedisClient mocks the Redis client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	args := m.Called(ctx, key, field)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, values)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	args := m.Called(ctx, key, fields)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringStringMapCmd)
}

// MockStringCmd mocks the Redis StringCmd
type MockStringCmd struct {
	mock.Mock
}

func (m *MockStringCmd) Result() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockStringCmd) Val() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStringCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStringCmd) Bytes() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

// MockStatusCmd mocks the Redis StatusCmd
type MockStatusCmd struct {
	mock.Mock
}

func (m *MockStatusCmd) Result() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockStatusCmd) Val() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockStatusCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockIntCmd mocks the Redis IntCmd
type MockIntCmd struct {
	mock.Mock
}

func (m *MockIntCmd) Result() (int64, error) {
	args := m.Called()
	return args.Int64(0), args.Error(1)
}

func (m *MockIntCmd) Val() int64 {
	args := m.Called()
	return args.Int64(0)
}

func (m *MockIntCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockBoolCmd mocks the Redis BoolCmd
type MockBoolCmd struct {
	mock.Mock
}

func (m *MockBoolCmd) Result() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockBoolCmd) Val() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockBoolCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

// MockStringStringMapCmd mocks the Redis StringStringMapCmd
type MockStringStringMapCmd struct {
	mock.Mock
}

func (m *MockStringStringMapCmd) Result() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockStringStringMapCmd) Val() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockStringStringMapCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	mockRedis := &MockRedisClient{}

	// Setup default behaviors
	mockRedis.On("Get", mock.Anything, mock.Anything).Return(&MockStringCmd{})
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(&MockStatusCmd{})
	mockRedis.On("Del", mock.Anything, mock.Anything).Return(&MockIntCmd{})
	mockRedis.On("Exists", mock.Anything, mock.Anything).Return(&MockIntCmd{})
	mockRedis.On("Expire", mock.Anything, mock.Anything, mock.Anything).Return(&MockBoolCmd{})

	return mockRedis
}

// SetupMockBehaviors allows setting up specific mock behaviors
func SetupMockBehaviors(redisClient *MockRedisClient) {
	// Setup successful GET operation
	redisClient.On("Get", mock.Anything, "user:1").Return(&MockStringCmd{
		mock.Mock{},
	}).Return("user:1", nil)

	// Setup successful SET operation
	redisClient.On("Set", mock.Anything, "user:1", mock.Anything, mock.Anything).Return(&MockStatusCmd{
		mock.Mock{},
	}).Return("OK", nil)

	// Setup user exists check
	redisClient.On("Exists", mock.Anything, "user:1").Return(&MockIntCmd{
		mock.Mock{},
	}).Return(1, nil)
}

// SetupErrorBehaviors allows setting up error scenarios
func SetupErrorBehaviors(redisClient *MockRedisClient) {
	// Setup Redis connection error
	redisClient.On("Get", mock.Anything, "nonexistent").Return(&MockStringCmd{
		mock.Mock{},
	}).Return("", errors.New("redis connection error"))

	// Setup Redis operation error
	redisClient.On("Set", mock.Anything, "invalid", mock.Anything, mock.Anything).Return(&MockStatusCmd{
		mock.Mock{},
	}).Return("", errors.New("redis operation error"))
}