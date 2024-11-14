package data

import (
	"context"
	"database/sql"
)

/*
	
	Delete(ctx context.Context, id int64) error

 */


// 注销用户
func (m *MySQLUserStore) Delete(ctx context.Context, id int64) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {

		// 删除用户信息
		if err := m.delete(ctx, tx, id); err != nil {
			return err
		}

		// 删除用户邀请码
		if err := m.deleteUserInvitations(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})
}

// 删除用户信息
func (m *MySQLUserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `delete from users where id = ?`

	_, err := tx.ExecContext(ctx, stmt, id)
	return err
}

// 删除用户邀请码
func (m *MySQLUserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `delete from user_invitations where user_id =?`
	
	_, err := tx.ExecContext(ctx, stmt, id)
	return err
}
