package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// 帖子
type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments"`
	User      User       `json:"user"`
}

type PostWithMetadata struct {
	Post          Post `json:"post"`
	CommentsCount int  `json:"comments_count"`
}

type PostStoreImpl struct {
	db *sql.DB
}


func (p *PostStoreImpl)	GetOne(ctx context.Context, postID int64) (*Post, error) {
	query := `
		select id, user_id, title, content, created_at,  updated_at, tags, version
		from posts
		where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var i Post

	// 特殊处理
	var tags string

	row := p.db.QueryRowContext(ctx, query, postID)

	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&tags,  // 临时工, 顶住
		&i.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// 特殊处理, 分割
	i.Tags = strings.Split(tags, ",")

	return &i, nil 
}


func (p *PostStoreImpl) Create(ctx context.Context, params Post) error {
	stmt := `
		insert into posts (content, title, user_id, tags, created_at, updated_at, version)
		values (?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// 特殊处理, 拼接
	tags := strings.Join(params.Tags, ",")

	_, err := p.db.ExecContext(ctx, stmt, 
		params.Content, 
		params.Title,
		params.UserID, 
		tags,
		params.CreatedAt, 
		params.UpdatedAt, 
		params.Version,
	)
	return err
}


func (p *PostStoreImpl) Delete(ctx context.Context, postID int64) error {
	stmt := `
		delete from posts where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := p.db.ExecContext(ctx, stmt, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil 
}

func (p *PostStoreImpl) Update(ctx context.Context, params Post) error {
	stmt := `
		update posts
		set title = ?, content = ?, version = version + 1, updated_at = ?
		where id = ? and version = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := p.db.ExecContext(ctx, stmt, 
		params.Title, 
		params.Content, 
		params.UpdatedAt, 
		params.ID, 
		params.Version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil 
}

func (p *PostStoreImpl) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count 
		FROM 
			posts p
			LEFT JOIN comments c ON c.post_id = p.id
			LEFT JOIN users u ON p.user_id = u.id
			JOIN followers f ON f.follower_id = p.user_id OR p.user_id = ?
		WHERE 
			(article.title LIKE %?% OR p.content LIKE %?%) 
			AND f.user_id = ? 
		GROUP BY 
			p.id, u.username
		ORDER BY 
			p.created_at ?
		LIMIT ?, ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query, 
		userID, 
		fq.Search,
		fq.Search,
		fq.Sort,
		fq.Limit, 
		fq.Offset, 
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()


	var feed []PostWithMetadata

	// 特殊处理
	var tags string

	for rows.Next() {
		var i PostWithMetadata

		err := rows.Scan(
			&i.Post.ID,
			&i.Post.UserID,
			&i.Post.Title,
			&i.Post.Content,
			&i.Post.CreatedAt,
			&i.Post.Version,
			&tags,  // 暂代
			&i.Post.User.Username,
			&i.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		// 特殊处理, 分割
		i.Post.Tags = strings.Split(tags, ",")


		feed = append(feed, i)
	}

	return feed, nil
}
