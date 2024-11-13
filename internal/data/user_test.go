package data

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"

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

func TestCreateAndInvite(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)
	
	params := User{
		Username:  "gua",
		Email:     "gua@vip.cn",
		Role:      1,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := params.Password.Set("123456"); err != nil {
		t.Fatal(err)
	}

	plainToken := uuid.New().String()
	hashToken := utils.HashifyStr(plainToken)

	if err := store.Users.CreateAndInvite(context.TODO(), params, hashToken, time.Hour); err != nil {
		t.Fatal(err)
	}
	t.Log("创建用户和邀请码成功")
}


