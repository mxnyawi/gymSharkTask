package mocks

import (
	"github.com/mxnyawi/gymSharkTask/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockDBManager struct {
	mock.Mock
}

func (m *MockDBManager) GetDocument(bucketName, scopeName, collectionName, documentID string) (*db.DocumentHistory, error) {
	args := m.Called(bucketName, scopeName, collectionName, documentID)
	return args.Get(0).(*db.DocumentHistory), args.Error(1)
}

func (m *MockDBManager) GetUser(bucketName, scopeName, collectionName, documentID string) (*db.User, error) {
	args := m.Called(bucketName, scopeName, collectionName, documentID)
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockDBManager) WriteDocument(bucket, scope, collection, id string, data interface{}) error {
	args := m.Called(bucket, scope, collection, id, data)
	return args.Error(0)
}

func (m *MockDBManager) GetDBCreds() (string, string, string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.String(2), args.String(3), args.Error(4)
}

func (m *MockDBManager) CreateAdminUser(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockDBManager) CreateBucket(bucketName string) error {
	args := m.Called(bucketName)
	return args.Error(0)
}

func (m *MockDBManager) CreateScope(bucketName, scopeName string) error {
	args := m.Called(bucketName, scopeName)
	return args.Error(0)
}

func (m *MockDBManager) CreateCollection(bucketName, scopeName, collectionName string) error {
	args := m.Called(bucketName, scopeName, collectionName)
	return args.Error(0)
}

func (m *MockDBManager) GetClusterCredentials() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}
