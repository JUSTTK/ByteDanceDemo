package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserDAO mocks the User DAO
type MockUserDAO struct {
	mock.Mock
}

func (m *MockUserDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockUserDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockUserDAO) Where(field string, value interface{}) *MockUserDAO {
	args := m.Called(field, value)
	return args.Get(0).(*MockUserDAO)
}

func (m *MockUserDAO) Eq(value interface{}) *MockUserDAO {
	args := m.Called(value)
	return args.Get(0).(*MockUserDAO)
}

func (m *MockUserDAO) ID(field string) *MockUserDAO {
	args := m.Called(field)
	return args.Get(0).(*MockUserDAO)
}

func (m *MockUserDAO) Pluck(column string, dest interface{}) error {
	args := m.Called(column, dest)
	return args.Error(0)
}

// MockVideoDAO mocks the Video DAO
type MockVideoDAO struct {
	mock.Mock
}

func (m *MockVideoDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockVideoDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockCommentDAO mocks the Comment DAO
type MockCommentDAO struct {
	mock.Mock
}

func (m *MockCommentDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockCommentDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockLikeDAO mocks the Like DAO
type MockLikeDAO struct {
	mock.Mock
}

func (m *MockLikeDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockLikeDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockMessageDAO mocks the Message DAO
type MockMessageDAO struct {
	mock.Mock
}

func (m *MockMessageDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockMessageDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockRelationDAO mocks the Relation DAO
type MockRelationDAO struct {
	mock.Mock
}

func (m *MockRelationDAO) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockRelationDAO) Find(conds ...interface{}) ([]interface{}, error) {
	args := m.Called(conds...)
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockDB mocks the database connection
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	args := m.Called(query, args)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Pluck(column string, dest interface{}) *gorm.DB {
	args := m.Called(column, dest)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}