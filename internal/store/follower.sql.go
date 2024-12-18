package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// 粉丝, 我关注了峰哥亡命天涯
type Follower struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"` 
	CreatedAt  string `json:"created_at"`
}

type FollowertoreImpl struct {
	db *sql.DB
}


// 关注
func (f *FollowertoreImpl) Follow(ctx context.Context, userID, followerID int64) error {
	stmt := `
		insert into followers (user_id, follower_id, created_at) 
		values (?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, stmt, userID, followerID, time.Now())
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && 
					strings.Contains(mysqlErr.Message, "followers.unique_follow") {
				return ErrConflict
			}
		}

		return err
	}

	return nil 
}

// 取关(例, 我取关了峰哥亡命天涯)
func (f *FollowertoreImpl) Unfollow(ctx context.Context, followerID, userID int64) error {
	stmt := `
		delete from followers 
		where user_id = ? and follower_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := f.db.ExecContext(ctx, stmt, userID, followerID)
	if err != nil {
		return err
	}

	rowsAffected, err:= result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil 
}


// 关注列表
func (f *FollowertoreImpl) GetAllFollowed(ctx context.Context, userID int64) ([]*Follower, error) {
	query := `
		select id, user_id, follower_id, created_at from followers 
		where user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := f.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var followers []*Follower
	
	for rows.Next() {
		var i Follower
		
		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FollowerID,
			&i.CreatedAt,
		)
		
		if err != nil {
			return nil, err
		}

		followers = append(followers, &i)
	}

	return followers, nil 
}

// 关注数
func (f *FollowertoreImpl) GetAllFollowedCount(ctx context.Context, userID int64) (int64, error) {
	query := `
		select count(*) from followers 
		where user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var count int64
	
	if err := f.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err 
	}

	return count, nil 
}

// 粉丝列表
func (f *FollowertoreImpl) GetAllFollower(ctx context.Context, userID int64) ([]*Follower, error) {
	query := `
		select id, user_id, created_at from followers 
		where follower_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := f.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	var followers []*Follower
	
	for rows.Next() {
		var i Follower
		
		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FollowerID,
			&i.CreatedAt,
		)
		
		if err != nil {
			return nil, err
		}

		followers = append(followers, &i)
	}

	return followers, nil 
}

// 粉丝数
func (f *FollowertoreImpl) GetAllFollowerCount(ctx context.Context, userID int64) (int64, error) {
	query := `
		select count(*) from followers 
		where follower_id = ?
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var count int64
	
	if err := f.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err 
	}

	return count, nil 
}



