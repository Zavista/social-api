package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type PostRepository interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
}

type UserRepository interface {
	Create(context.Context, *User) error
}

type CommentRepository interface {
	Create(context.Context, *Comment) error
	GetByPostID(context.Context, int64) ([]Comment, error)
}

type Storage struct {
	Posts    PostRepository
	Users    UserRepository
	Comments CommentRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
