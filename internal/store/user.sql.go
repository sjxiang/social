package store

import (
	"time"
	 "database/sql"

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