package data

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/sjxiang/social/internal/utils"
)


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


func TestActivate(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)

	if err := store.Users.Activate(context.TODO(), "bc8fb4b8efcd46753c8a9499a1fb6f2d"); err!= nil {
		t.Fatal(err)
	}

	t.Log("激活用户成功")
}



func TestGetOneUser(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)

	user, err := store.Users.GetOne(context.TODO(), 7)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", user)
	t.Log("查询单个用户, ok")
}

func TestGetOneUserByEmail(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)

	user, err := store.Users.GetByEmail(context.TODO(), "gua@vip.cn")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", user)
	t.Log("查询单个用户, ok")
}


func TestModPassword(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)

	params := User{
		ID: 7,
	}

	if err := params.Password.Set("654321"); err != nil {
		t.Fatal(err)
	}

	if err := store.Users.ModPassword(context.TODO(), params); err!= nil {
		t.Fatal(err)
	}

	t.Log("修改用户密码成功")
}

