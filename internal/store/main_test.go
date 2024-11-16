package store

import (
	"time"
	"database/sql"
	"context"

	"github.com/sjxiang/social/internal/utils"
)

func initdb() (*sql.DB, error) {
	dsn := "root:my-secret-pw@tcp(127.0.0.1:13306)/social?charset=utf8&parseTime=True&loc=Local"

	db, err := utils.NewMySQL(dsn, 30, 30, time.Minute*15)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}	
