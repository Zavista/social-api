package store

import (
	"context"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
		Posts: &MockPostStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	return &User{ID: userID}, nil
}

func (m *MockUserStore) GetByIDWithRole(ctx context.Context, userID int64) (*UserWithRole, error) {
	return &UserWithRole{User: User{ID: userID}}, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{Email: email}, nil
}

func (m *MockUserStore) Create(ctx context.Context, u *User) error {
	return nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}

type MockPostStore struct{}

func (m *MockPostStore) Create(ctx context.Context, post *Post) error {
	post.ID = 1
	return nil
}

func (m *MockPostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	return &Post{ID: id}, nil
}

func (m *MockPostStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *MockPostStore) Update(ctx context.Context, post *Post) error {
	return nil
}

func (m *MockPostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	return nil, nil
}
