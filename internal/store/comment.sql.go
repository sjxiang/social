package store

import (
	"context"
	"database/sql"
	"time"
)

// 评论
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type CommentStoreImpl struct {
	db *sql.DB
}


func (c *CommentStoreImpl) GetByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	query := `
		select c.id, c.post_id, c.user_id, c.content, c.created_at, u.username  
		from comments c
		left join users u on (u.id = c.user_id)
		where c.post_id = ?
		order by c.created_at desc;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	
	var comments []*Comment

	for rows.Next() {
		var i Comment

		err := rows.Scan(
			&i.ID, 
			&i.PostID, 
			&i.UserID, 
			&i.Content, 
			&i.CreatedAt, 
			&i.User.Username,
		)
		if err != nil {
			return nil, err
		}
		
		comments = append(comments, &i)
	}

	return comments, nil
}

func (c *CommentStoreImpl) Create(ctx context.Context, params Comment) error {
	query := `
		insert into comments (post_id, user_id, content, created_at)
		values (?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := c.db.ExecContext(
		ctx,
		query,
		params.PostID,
		params.UserID,
		params.Content,
		params.CreatedAt,
	)

	return err
}
