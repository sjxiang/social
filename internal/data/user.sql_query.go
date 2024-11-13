package data

import (
	"context"
	"database/sql"
	"time"
	"errors"
)

/*

	GetOne(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Exists(ctx context.Context, id int64) (bool, error)

*/


 func (m *MySQLUserStore) GetOne(ctx context.Context, id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query1 := `
		select id, username, email, password, is_active, role, created_at, updated_at
		from users
		where id = ? and is_active = 1`

	var i User
	row := m.db.QueryRowContext(ctx, query1, id)

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password.hash,
		&i.IsActive,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	query2 := `
		select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at 
		from plans p
		left join user_plans up on (p.id = up.plan_id)
		where up.user_id = ?`

	rowAgain := m.db.QueryRowContext(ctx, query2, id)

	if err := rowAgain.Scan(
		&i.Plan.ID,
		&i.Plan.PlanName,
		&i.Plan.PlanAmount,
		&i.Plan.CreatedAt,
		&i.Plan.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &i, nil
		} else {
			return nil, err
		}
	}

	return &i, nil
}


func (m *MySQLUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query1 := `
		select id, username, email, password, is_active, role, created_at, updated_at
		from users
		where email = ? and is_active = 1
		limit 1`

	var i User
	row := m.db.QueryRowContext(ctx, query1, email)

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password.hash,
		&i.IsActive,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}


	query2 := `
		select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at 
		from plans p
		left join user_plans up on (p.id = up.plan_id)
		where up.user_id = ?`

	rowAgain := m.db.QueryRowContext(ctx, query2, i.ID)

	if err := rowAgain.Scan(
		&i.Plan.ID,
		&i.Plan.PlanName,
		&i.Plan.PlanAmount,
		&i.Plan.CreatedAt,
		&i.Plan.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &i, nil
		} else {
			return nil, err
		}
	}

	return &i, nil
}


// 检查用户是否存在
func (m *MySQLUserStore) Exists(ctx context.Context, id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var exists bool

	query := `select exists(select true from users where id = ?)`

	err := m.db.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}


// 根据邀请码查询用户的详细信息
func (m *MySQLUserStore) queryUserDetailsBasedOnInvitationToken(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// 筛选, 邀请码有效
	query := `
		select u.id, u.username, u.email, u.is_active, u.role, u.created_at, u.updated_at
		from users u
		left join user_invitations ui on (u.id = ui.user_id)
		where ui.token = ? and ui.expiry > ?`

	var i User
	row := m.db.QueryRowContext(ctx, query, token, time.Now())

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.IsActive,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &i, nil
}
