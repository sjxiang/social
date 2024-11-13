package data

import (
	"database/sql"

	"github.com/sjxiang/social/internal/utils"
)

func initDB() *sql.DB {
	dsn := "root:my-secret-pw@tcp(127.0.0.1:13306)/social?charset=utf8&parseTime=True&loc=Local"

	db, err := utils.NewMySQL(dsn, 30, 30, "15m")
	if err != nil {
		panic(err)
	}

	return db
}	
