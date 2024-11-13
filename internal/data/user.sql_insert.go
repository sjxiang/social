package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"strings"

	"github.com/go-sql-driver/mysql"
)


func (m *MySQLUserStore) CreateAndInvite(ctx context.Context, arg User, token string, invitationExp time.Duration) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {
		
		// 创建用户
		if err := m.create(ctx, tx, arg); err != nil {
			return err
		}

		// 创建邀请码
		if err := m.createUserInvitation(ctx, tx, token, invitationExp, arg.ID); err != nil {
			return err
		}

		return nil
	})
}

// 创建用户
func (m *MySQLUserStore) create(ctx context.Context, tx *sql.Tx, arg User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `insert into users (username, email, password, is_active, role, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.ExecContext(ctx, stmt, 
		arg.Username, 
		arg.Email,
		arg.Password.hash, 
		arg.IsActive,
		arg.Role,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	if err != nil {
		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) {
			switch {
			case mysqlError.Number == 1062 && 
					strings.Contains(mysqlError.Message, "users.idx_username"):
				return ErrDuplicateUsername
			case mysqlError.Number == 1062 && 
					strings.Contains(mysqlError.Message, "users.idx_email"):
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}


// 创建邀请码
func (m *MySQLUserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `insert into user_invitations (token, user_id, expiry) 
		values (?, ?, ?)`

	_, err := tx.ExecContext(ctx, stmt, token, userID, time.Now().Add(exp))
	return err 
}

