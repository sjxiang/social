package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"strings"

	"github.com/go-sql-driver/mysql"
)


type MySQLUserStore struct {
	db *sql.DB
}

func NewMySQLUserStore(db *sql.DB) *MySQLUserStore {
	return &MySQLUserStore{db: db}
}


// 检查用户是否存在
func (s *MySQLUserStore) Exists(ctx context.Context, userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	

	return false, nil
}

// 修改密码
func (s *MySQLUserStore) ModPassword(ctx context.Context, arg *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return nil
}


// 创建用户
func (m *MySQLUserStore) Create(ctx context.Context, tx *sql.Tx, arg *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		INSERT INTO users 
			(username, password, email, created_at, is_active, role_id) 
		VALUES
			(?, ?, ?, UTC_TIMESTAMP(), 0, ?)`
	
	// 设置用户角色
	arg.SetRole()

	// 插入用户信息
	_, err := tx.ExecContext(ctx, stmt, 
		arg.Username, 
		arg.Password.hash, 
		arg.Email, 
		arg.RoleID,
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


// 创建用户的邀请码
func (m *MySQLUserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		INSERT INTO user_invitations 
			(token, user_id, expiry) 
		VALUES 
			(?, ?, ?)`

	_, err := tx.ExecContext(ctx, stmt, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}


// 根据 id 获取用户信息
func (m *MySQLUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT 
			u.id, u.username, u.email, u.password, u.created_at, r.id, r.name, r.level, r.description
		FROM 
			users u
		LEFT JOIN 
			roles r ON u.role_id = r.id
		WHERE 
			u.id = ? AND u.is_active = true`

	var i User
	row := m.db.QueryRowContext(ctx, query, userID)

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password.hash,
		&i.CreatedAt,
		&i.Role.ID,
		&i.Role.Name,
		&i.Role.Level,
		&i.Role.Description,
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

// 创建用户并创建邀请码
func (m *MySQLUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {
		
		if err := m.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := m.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// 激活用户
func (m *MySQLUserStore) Activate(ctx context.Context, token string) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user that this token belongs to
		user, err := m.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = true
		if err := m.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. clean the invitations
		if err := m.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// 根据邀请码查找用户信息
func (m *MySQLUserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, hashToken string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// 筛选条件 邀请码没过期
	query := `
		SELECT 
			u.id, u.username, u.email, u.created_at, u.is_active
		FROM 
			users u
		LEFT JOIN 
			user_invitations ui ON u.id = ui.user_id
		WHERE 
			ui.token = ? AND ui.expiry > ?`

	var i User
	row := tx.QueryRowContext(ctx, query, hashToken, time.Now())

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.IsActive,
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


func (m *MySQLUserStore) update(ctx context.Context, tx *sql.Tx, arg *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		UPDATE 
			users 
		SET 
			username = ?, 
			email = ?, 
			is_active = ? 
		WHERE 
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

// 删除用户的邀请码
func (m *MySQLUserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		DELETE FROM 
			user_invitations 
		WHERE 
			user_id = ?`
	
	_, err := tx.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}

	return nil
}

// 注销用户
func (m *MySQLUserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(m.db, ctx, func(tx *sql.Tx) error {
		if err := m.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := m.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

// 删除用户
func (m *MySQLUserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		DELETE FROM 
			users 
		WHERE 
			id = ?`

	_, err := tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}


func (m *MySQLRoleStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT 
			id, username, email, password, created_at
		FROM
			users
		WHERE
			email = ? AND is_active = 1
		LIMIT
			1`

	var i User

	row := m.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password.hash,
		&i.CreatedAt,
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
