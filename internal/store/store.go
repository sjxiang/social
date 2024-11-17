package store

import (
	"context"
	"database/sql"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)


type Storage struct {
	UserStore interface {
		GetOne(ctx context.Context, userID int64) (*User, error)
		Exists(ctx context.Context, userID int64) (bool, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		
		// admin
		Delete(ctx context.Context, userID int64) error
		Create(ctx context.Context, params User) (int64, error)
	
		// user
		CreateAndInvite(ctx context.Context, params User, token string, invitationExpiry time.Duration) error
		Activate(ctx context.Context, token string) error
		ModPassword(ctx context.Context, params User) error
	}

	PostStore interface {
		GetOne(ctx context.Context, postID int64) (*Post, error)
		Create(ctx context.Context, params Post) error
		Delete(ctx context.Context, postID int64) error
		Update(ctx context.Context, params Post) error
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	
	CommentStore interface {
		Create(ctx context.Context, params Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]*Comment, error)
	}
	
	PlanStore interface {
		GetAll(ctx context.Context) ([]*Plan, error)
		GetOne(ctx context.Context, planID int64) (*Plan, error)
		SubscribeUserToPlan(ctx context.Context, params User) error
	}

	FollowerStore interface {
		Follow(ctx context.Context, userID, followerID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
		GetAllFollowed(ctx context.Context, userID int64) ([]*Follower, error)  
		GetAllFollowedCount(ctx context.Context, userID int64) (int64, error)
		GetAllFollower(ctx context.Context, userID int64) ([]*Follower, error) 
		GetAllFollowerCount(ctx context.Context, userID int64) (int64, error)
	}

}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		UserStore: &UserStoreImpl{db: db},
		PostStore: &PostStoreImpl{db: db},
	}
}



// todo
type MySQLRepo interface {

}
