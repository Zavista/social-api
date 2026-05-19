package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

// DBTX is implemented by both *sql.DB and *sql.Tx.
//
// This allows repository helper methods to work with either:
//   - the base database connection pool (*sql.DB)
//   - an active transaction (*sql.Tx)
//
// We use this pattern so shared query logic can be reused inside
// transactions without duplicating methods like:
//
//	Create()
//	CreateWithTx()
//	CreateInsideTransaction()
//
// Example:
//
//	s.create(ctx, s.db, user) // normal query
//	s.create(ctx, tx, user)   // transactional query
type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type PostRepository interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
	GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
}

type UserRepository interface {
	GetByID(context.Context, int64) (*User, error)
	Create(context.Context, *User) error
	CreateAndInvite(context.Context, *User, string, time.Duration) error
	Activate(context.Context, string) error
	Delete(context.Context, int64) error
}

type CommentRepository interface {
	Create(context.Context, *Comment) error
	GetByPostID(context.Context, int64) ([]Comment, error)
}

type FollowerRepository interface {
	Follow(ctx context.Context, followedID, followerID int64) error
	Unfollow(ctx context.Context, followedID, followerID int64) error
}

type Storage struct {
	Posts     PostRepository
	Users     UserRepository
	Comments  CommentRepository
	Followers FollowerRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
