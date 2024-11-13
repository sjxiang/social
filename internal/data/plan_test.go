package data

import (
	"context"
	"testing"
	"time"
)


func TestGetAllPlans(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)
	items, err := store.Plans.GetAll(context.Background())
	if err!= nil {
		t.Fatal(err)
	}

	for _, item := range items {
		t.Log(item)
	}

	t.Log("获取全部的订阅计划, ok")
}

func TestGetOnePlan(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)
	plan, err := store.Plans.GetOne(context.Background(), 4)
	if err!= nil {
		t.Fatal(err)
	}

	t.Log(plan)

	t.Log("查询单个订阅计划, ok")
}


func TestSubscribeUserToPlan(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	
	db := initDB()
	defer db.Close()

	store := NewMySQLStorage(db)

	// 1. 创建用户
	params := User{
		ID: 7,
		Plan: Plan{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

	}
	if err := store.Plans.SubscribeUserToPlan(context.Background(), params); err!= nil {
		t.Fatal(err)
	}

	t.Log("订阅计划成功, ok")
}