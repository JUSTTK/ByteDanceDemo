package mocks

import (
	"bytedancedemo/database/redis"
	"github.com/stretchr/testify/mock"
)

// MockCache mocks the cache operations
type MockCache struct {
	mock.Mock
}

// NewMockCache creates a new mock cache
func NewMockCache() *MockCache {
	return &MockCache{}
}

// SetupDefaultBehaviors sets up default mock behaviors
func (m *MockCache) SetupDefaultBehaviors() {
	// Mock successful user cache
	m.On("GetUser", mock.Anything, int64(1)).Return(map[string]interface{}{
		"id":         int64(1),
		"name":       "testuser",
		"avatar":     "avatar.jpg",
		"follow_cnt": int64(100),
	}, nil)

	// Mock empty cache
	m.On("GetUser", mock.Anything, int64(999)).Return(nil, nil)
}

// SetupErrorBehaviors sets up error scenarios
func (m *MockCache) SetupErrorBehaviors() {
	// Mock cache error
	m.On("GetUser", mock.Anything, int64(0)).Return(nil, redis.ErrNil)

	// Mock DB error
	m.On("SetUser", mock.Anything, int64(0), mock.Anything).Return(redis.ErrNil)
}

// Cache methods
func (m *MockCache) GetUser(db *redis.Client, userId int64) (map[string]interface{}, error) {
	args := m.Called(db, userId)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockCache) SetUser(db *redis.Client, userId int64, data map[string]interface{}) error {
	args := m.Called(db, userId, data)
	return args.Error(0)
}

func (m *MockCache) GetUserFollowers(db *redis.Client, userId int64) ([]int64, error) {
	args := m.Called(db, userId)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockCache) GetUserFollowing(db *redis.Client, userId int64) ([]int64, error) {
	args := m.Called(db, userId)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockCache) GetUserVideoCount(db *redis.Client, userId int64) (int64, error) {
	args := m.Called(db, userId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) GetUserFavoriteCount(db *redis.Client, userId int64) (int64, error) {
	args := m.Called(db, userId)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) SetUserFollowers(db *redis.Client, userId int64, followers []int64) error {
	args := m.Called(db, userId, followers)
	return args.Error(0)
}

func (m *MockCache) SetUserFollowing(db *redis.Client, userId int64, following []int64) error {
	args := m.Called(db, userId, following)
	return args.Error(0)
}

func (m *MockCache) SetUserVideoCount(db *redis.Client, userId int64, count int64) error {
	args := m.Called(db, userId, count)
	return args.Error(0)
}

func (m *MockCache) SetUserFavoriteCount(db *redis.Client, userId int64, count int64) error {
	args := m.Called(db, userId, count)
	return args.Error(0)
}

// MockDBOperation mocks the database operations
type MockDBOperation struct {
	mock.Mock
}

// NewMockDBOperation creates a new mock DB operation
func NewMockDBOperation() *MockDBOperation {
	return &MockDBOperation{}
}

// SetupDefaultBehaviors sets up default mock behaviors
func (m *MockDBOperation) SetupDefaultBehaviors() {
	// Mock successful user creation
	m.On("CreateUser", mock.Anything).Return(int64(1), nil)

	// Mock successful user query
	m.On("GetUserById", mock.Anything, int64(1)).Return(map[string]interface{}{
		"id":    int64(1),
		"name":  "testuser",
		"email": "test@example.com",
	}, nil)

	// Mock user not found
	m.On("GetUserById", mock.Anything, int64(999)).Return(nil, redis.ErrNil)
}

// SetupErrorBehaviors sets up error scenarios
func (m *MockDBOperation) SetupErrorBehaviors() {
	// Mock DB error
	m.On("CreateUser", mock.Anything).Return(int64(0), redis.ErrNil)

	// Mock query error
	m.On("GetUserById", mock.Anything, int64(0)).Return(nil, redis.ErrNil)
}

// DBOperation methods
func (m *MockDBOperation) CreateUser(user map[string]interface{}) (int64, error) {
	args := m.Called(user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDBOperation) GetUserById(db *redis.Client, userId int64) (map[string]interface{}, error) {
	args := m.Called(db, userId)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockDBOperation) UpdateUser(db *redis.Client, userId int64, updates map[string]interface{}) error {
	args := m.Called(db, userId, updates)
	return args.Error(0)
}

func (m *MockDBOperation) DeleteUser(db *redis.Client, userId int64) error {
	args := m.Called(db, userId)
	return args.Error(0)
}