package store

import (
	"context"
	"testing"
	"time"
)

func TestModPassword(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userStore := newUserStore(db)
	
	params := User{
		ID: 7,
	}
	if err := params.Password.Set("000111"); err != nil {
		t.Fatal(err)
	}

	if err := userStore.ModPassword(context.TODO(), params); err!= nil {
		t.Fatal(err)
	}

	t.Log("修改用户密码成功")
}

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userStore := newUserStore(db)

	params := User{
		Username: "xuan",
		Email:    "xuan@163.com",
		IsActive: true,
		Role:     2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := params.Password.Set("123456"); err != nil {
		t.Fatal(err)
	}

	id, err := userStore.Create(context.TODO(), params)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("用户id", id)
	t.Log("管理员 admin 创建用户, 成功")
}


func TestGetOneUser(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userStore := newUserStore(db)

	user, err := userStore.GetOne(context.TODO(), 7)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", user)
	t.Log("通过id查询用户, 成功")
}


func TestGetUserByEmail(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userStore := newUserStore(db)

	user, err := userStore.GetByEmail(context.TODO(), "xuan@163.com")  // xuan@163.com gua@vip.cn
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", user)
	t.Log("通过 email 查询用户, 成功")
}