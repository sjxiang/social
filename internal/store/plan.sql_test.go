package store

import (
	"context"
	"testing"
	"time"
)


func TestGetAllPlans(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	planStore := newPlanStore(db)

	items, err := planStore.GetAll(context.TODO())
	if err!= nil {
		t.Fatal(err)
	}

	if len(items) > 0 {
		for _, item := range items {
			t.Log(item)
		}
	
		t.Log("获取全部的订阅计划, ok")
	}

	t.Log("订阅计划为空")
}


func TestGetOnePlan(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	planStore := newPlanStore(db)

	plan, err := planStore.GetOne(context.TODO(), 4)
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
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	planStore := newPlanStore(db)

	params := User{
		ID: 7,
		Plan: Plan{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	if err := planStore.SubscribeUserToPlan(context.TODO(), params); err!= nil {
		t.Fatal(err)
	}

	t.Log("订阅成功, ok")
}