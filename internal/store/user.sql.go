package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	IsActive  bool      `json:"is_active"`  
	Role      int64     `json:"role"`        // 角色, 0 游客, 1 用户, 2 版主, 3 管理员
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Plan      Plan      `json:"plan"`
}


type password struct {
	text *string  
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}


type UserStoreImpl struct {
	db *sql.DB
} 

func newUserStore(db *sql.DB) UserStore {
	return &UserStoreImpl{db: db}
}


func (u *UserStoreImpl) GetOne(ctx context.Context, userID int64) (*User, error) {
	query := `
		select id, username, email, password, is_active, role, created_at, updated_at
		from users
		where id = ? and is_active = 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var i User
	row := u.db.QueryRowContext(ctx, query, userID)

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

	i.Plan.ForDisplay()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	stmt := `
		select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at 
		from plans p
		left join user_plans up on (p.id = up.plan_id)
		where up.user_id = ?`

	err = u.db.QueryRowContext(ctx, stmt, userID).Scan(
		&i.Plan.ID,
		&i.Plan.PlanName,
		&i.Plan.PlanAmount,
		&i.Plan.CreatedAt,
		&i.Plan.UpdatedAt,
	)

	i.Plan.PlanName, i.Plan.PlanAmountFormatted = i.Plan.ForDisplay()

	switch err {
	case nil:
		return &i, nil
	case sql.ErrNoRows:
		return &i, nil
	default:
		return nil, err
	}
}

func (u *UserStoreImpl) Exists(ctx context.Context, userID int64) (bool, error) {
	query := `
		select exists
			(select true from users where id = ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var exists bool

	err := u.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	return exists, err
}

func (u *UserStoreImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		select id, username, email, password, is_active, role, created_at, updated_at
		from users
		where email = ? and is_active = 1
		limit 1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	
	var i User
	row := u.db.QueryRowContext(ctx, query, email)

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err 
		}
	}


	stmt := `
		select p.id, p.plan_name, p.plan_amount, p.created_at, p.updated_at 
		from plans p
		left join user_plans up on (p.id = up.plan_id)
		where up.user_id = ?`

	err = u.db.QueryRowContext(ctx, stmt, i.ID).Scan(
		&i.Plan.ID,
		&i.Plan.PlanName,
		&i.Plan.PlanAmount,
		&i.Plan.CreatedAt,
		&i.Plan.UpdatedAt,
	)
	if err != nil {
		// 订阅记录没查到, 但还有用户记录
		if errors.Is(err, sql.ErrNoRows) {
			return &i, nil
		} else {
			return nil, err
		}
	}
	i.Plan.PlanName, i.Plan.PlanAmountFormatted = i.Plan.ForDisplay()

	return &i, nil
}

func (u *UserStoreImpl) delete(ctx context.Context, tx *sql.Tx, userID int64) error {	
	stmt := `
		delete from users 
		where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, stmt, userID)
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

// admin
func (u *UserStoreImpl) Delete(ctx context.Context, userID int64) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {

		// 删除用户信息
		if err := u.delete(ctx, tx, userID); err != nil {
			return err
		}

		// bug, 激活码用完就删了, 除非你没有激活
		
		// 删除用户激活码
		if err := u.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}


// admin
func (u *UserStoreImpl) Create(ctx context.Context, params User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	stmt := `
		insert into users (username, email, password, is_active, role, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := u.db.ExecContext(ctx, stmt, 
		params.Username, params.Email, params.Password.hash, 
		params.IsActive, params.Role, 
		params.CreatedAt, params.UpdatedAt,
	)
	if err != nil {
	
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			switch {
			case strings.Contains(mysqlErr.Message, "users.idx_username"):
				return 0, ErrDuplicateUsername
			case strings.Contains(mysqlErr.Message, "users.idx_email"):
				return 0, ErrDuplicateEmail	
			}
		}
		
		return 0, err 
	}		
	
	return result.LastInsertId()
}


func (u *UserStoreImpl) ModPassword(ctx context.Context, params User) error {
	stmt := `
		update users 
		set password = ? 
		where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	
	result, err := u.db.ExecContext(ctx, stmt,
		params.Password.hash,
		params.ID,
	)
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

func (u *UserStoreImpl) create(ctx context.Context, tx *sql.Tx, params User) (int64, error) {
	stmt := `
		insert into users (username, email, password, is_active, role, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := tx.ExecContext(ctx, stmt, 
		params.Username, params.Email, params.Password.hash, 
		params.IsActive, params.Role,
		params.CreatedAt, params.UpdatedAt,
	)

	if err != nil {
		var mysqlError *mysql.MySQLError

		if errors.As(err, &mysqlError) {
			switch {
			case mysqlError.Number == 1062 && 
					strings.Contains(mysqlError.Message, "users.idx_username"):
				return 0, ErrDuplicateUsername
			case mysqlError.Number == 1062 && 
					strings.Contains(mysqlError.Message, "users.idx_email"):
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}

	return result.LastInsertId()
}


func (u *UserStoreImpl) update(ctx context.Context, tx *sql.Tx, params *User) error {
	stmt := `
		update users 
		set is_active = ?, updated_at = updated_at 
		where id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, stmt, 
		params.IsActive, 
		params.UpdatedAt, 
		params.ID,
	)
	return err
}