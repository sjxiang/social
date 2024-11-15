package store

import (
	"context"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)


type Storage struct {
	Users    UserStore
	Posts    PostStore
	Comments CommentStore
	Plans    PlanStore

}

func NewStorage(users UserStore, posts PostStore, comments CommentStore, plans PlanStore) Storage {
	return Storage{
		Users: users,
		Posts: posts,
		Comments: comments,
		Plans: plans,
	}
}



type UserStore interface {
	GetOne(ctx context.Context, id int64) (*User, error)
	Exists(ctx context.Context, id int64) (bool, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, id int64) error
	CreateAndInvite(ctx context.Context, arg User, token string, invitationExpiry time.Duration) error
	Activate(ctx context.Context, token string) error
	ModPassword(ctx context.Context, arg User) error
}

type PostStore interface {
	GetOne(ctx context.Context, id int64) (*Post, error)
}

type CommentStore interface {
	GetOne(ctx context.Context, id int64) (*Comment, error)
}

type PlanStore interface {
	GetAll(ctx context.Context) ([]*Plan, error)
	GetOne(ctx context.Context, id int64) (*Plan, error)
	SubscribeUserToPlan(ctx context.Context, arg User) error
}

// 关注
type FollowerStore interface {
	Follow(ctx context.Context, userID, followerID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
	GetAllFollowed(ctx context.Context, userID int64) ([]*Follower, error)
	GetAllFollowedCount(ctx context.Context, userID int64) (int64, error)
	GetAllFollower(ctx context.Context, userID int64) ([]*Follower, error) 
	GetAllFollowerCount(ctx context.Context, userID int64) (int64, error)
}