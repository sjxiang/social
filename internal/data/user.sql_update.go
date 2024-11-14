package data

import (
	"context"
	"database/sql"
)


// 修改密码
func (m *MySQLUserStore) ModPassword(ctx context.Context, arg User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		update 
			users 
		set 
			password = ? 
		where 
			id = ?`

	_, err := m.db.ExecContext(ctx, stmt,
		arg.Password.hash,
		arg.ID,
	)

	return err
}

// 激活用户
func (m *MySQLUserStore) Activate(ctx context.Context, token string) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {
		// 1. 通过邀请码, 查找用户
		user, err := m.queryUserDetailsBasedOnInvitationToken(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. 更改用户状态
		user.IsActive = true
		if err := m.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. 删除邀请码
		if err := m.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}


// 修改用户信息
func (m *MySQLUserStore) update(ctx context.Context, tx *sql.Tx, arg *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		update 
			users 
		set 
			username = ?, 
			email = ?, 
			is_active = ?,
			updated_at = UTC_TIMESTAMP() 
		where 
			id = ?`

	_, err := tx.ExecContext(ctx, stmt, 
		arg.Username, 
		arg.Email, 
		arg.IsActive, 
		arg.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

