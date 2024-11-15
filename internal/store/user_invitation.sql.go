package store

import (
	"context"
	"database/sql"
	"time"
)

// 用户激活码


// user
func (u *UserStoreImpl) CreateAndInvite(ctx context.Context, params User, token string, invitationExpiry time.Duration) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		
		// 创建用户
		userID, err := u.create(ctx, tx, params)
		if err != nil {
			return err
		}

		// 创建激活码
		if err := u.createUserInvitation(ctx, tx, token, invitationExpiry, userID); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStoreImpl) Activate(ctx context.Context, token string) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		
		// 1. 通过激活码, 查找用户
		user, err := u.findUserByToken(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. 更改用户状态
		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. 删除邀请码
		if err := u.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}


// 删除用户激活码
func (u *UserStoreImpl) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	stmt := `
		delete from user_invitations 
		where user_id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	
	_, err := tx.ExecContext(ctx, stmt, userID)
	return err
}

// 创建激活码
func (u *UserStoreImpl) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExpiry time.Duration, userID int64) error {
	stmt := `
		insert into user_invitations (token, user_id, expiry) 
		values (?, ?, ?)
	`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, stmt, token, userID, time.Now().Add(invitationExpiry))
	return err 
}

// 根据激活码, 查询用户
func (u *UserStoreImpl) findUserByToken(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	// 筛选, 邀请码有效
	query := `
		select u.id, u.username, u.email, u.is_active, u.role, u.created_at, u.updated_at
		from users u
		left join user_invitations ui on (u.id = ui.user_id)
		where ui.token = ? and ui.expiry > ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var i User
	row := u.db.QueryRowContext(ctx, query, token, time.Now())

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
