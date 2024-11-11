package utils

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"  // MySQL driver
)

func NewMySQL(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// dsn
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}