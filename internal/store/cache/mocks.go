package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/zavista/social-api/internal/store"
)

func NewMockStore() CacheStorage {
	return CacheStorage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.UserWithRole, error) {
	args := m.Called(userID)

	user, _ := args.Get(0).(*store.UserWithRole)
	return user, args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context, user *store.UserWithRole) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) {
	m.Called(userID)
}
