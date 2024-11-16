package store

import (
	"context"
	"database/sql"
)

// 帖子
type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post          Post `json:"post"`
	CommentsCount int  `json:"comments_count"`
}

type PostStoreImpl struct {
	db *sql.DB
}

func newPostStore(db *sql.DB) PostStore {
	return &PostStoreImpl{db: db}
}

// func (s *PostStoreImpl) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {

// 	return nil, nil
// }

// 	query := `
// 		SELECT 
// 			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
// 			u.username,
// 			COUNT(c.id) AS comments_count
// 		FROM posts p
// 		LEFT JOIN comments c ON c.post_id = p.id
// 		LEFT JOIN users u ON p.user_id = u.id
// 		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
// 		WHERE 
// 			f.user_id = $1 AND
// 			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
// 			(p.tags @> $5 OR $5 = '{}')
// 		GROUP BY p.id, u.username
// 		ORDER BY p.created_at ` + fq.Sort + `
// 		LIMIT $2 OFFSET $3
// 	`

// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var feed []PostWithMetadata
// 	for rows.Next() {
// 		var p PostWithMetadata
// 		err := rows.Scan(
// 			&p.ID,
// 			&p.UserID,
// 			&p.Title,
// 			&p.Content,
// 			&p.CreatedAt,
// 			&p.Version,
// 			pq.Array(&p.Tags),
// 			&p.User.Username,
// 			&p.CommentsCount,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		feed = append(feed, p)
// 	}

// 	return feed, nil
// }

// func (s *PostStore) Create(ctx context.Context, post *Post) error {
// 	query := `
// 		INSERT INTO posts (content, title, user_id, tags)
// 		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
// 	`

// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	err := s.db.QueryRowContext(
// 		ctx,
// 		query,
// 		post.Content,
// 		post.Title,
// 		post.UserID,
// 		pq.Array(post.Tags),
// 	).Scan(
// 		&post.ID,
// 		&post.CreatedAt,
// 		&post.UpdatedAt,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
// 	query := `
// 		SELECT id, user_id, title, content, created_at,  updated_at, tags, version
// 		FROM posts
// 		WHERE id = $1
// 	`

// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	var post Post
// 	err := s.db.QueryRowContext(ctx, query, id).Scan(
// 		&post.ID,
// 		&post.UserID,
// 		&post.Title,
// 		&post.Content,
// 		&post.CreatedAt,
// 		&post.UpdatedAt,
// 		pq.Array(&post.Tags),
// 		&post.Version,
// 	)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrNotFound
// 		default:
// 			return nil, err
// 		}
// 	}

// 	return &post, nil
// }

// func (s *PostStore) Delete(ctx context.Context, postID int64) error {
// 	query := `DELETE FROM posts WHERE id = $1`

// 	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
// 	defer cancel()

// 	res, err := s.db.ExecContext(ctx, query, postID)
// 	if err != nil {
// 		return err
// 	}

// 	rows, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rows == 0 {
// 		return ErrNotFound
// 	}

// 	return nil
// }


func (p *PostStoreImpl)	GetOne(ctx context.Context, postID int64) (*Post, error) {
	return nil, nil 
}

func (p *PostStoreImpl) Create(ctx context.Context, params Post) error {
	return nil 
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
	return nil, nil 
}