package store

import (
	"testing"
	"time"
	"context"
	
	"github.com/google/uuid"

	"github.com/sjxiang/social/internal/utils"
)

func TestCreateAndInvite(t *testing.T) {
	
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
		Username:  "zishen",
		Email:     "zishen@vip.cn",
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

	if err := userStore.CreateAndInvite(context.TODO(), params, hashToken, time.Hour); err != nil {
		t.Fatal(err)
	}

	t.Log("创建用户梓神和激活码, 成功")
}


func TestActivate(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userStore := newUserStore(db)

	if err := userStore.Activate(context.TODO(), "cc9a52360aca6902d7eb2ba346e6a775"); err!= nil {
		t.Fatal(err)
	}

	t.Log("激活用户成功")
}

