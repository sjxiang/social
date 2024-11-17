package store

import (
	"context"
	"testing"
)


func TestFollowed(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}

	tests := []struct {
		userID       int64
		name         string
		followerID   int64
		followerName string
    }{
        {1000, "阿瓜", 1001, "炫狗"},
        {1000, "阿瓜", 1002, "孙笑川"},
		{1000, "阿瓜", 1003, "娟妹"},
        {1000, "阿瓜", 1004, "李老八"},
		
		{2000, "梓神", 1001, "炫狗"},
		{2000, "梓神", 1002, "萧瑟"},

		{1001, "炫狗", 2000, "梓神"},
    }

    for _, test := range tests {
        err := followStore.Follow(context.TODO(), test.userID, test.followerID)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(test.name, "关注", test.followerName)
    }
    
	t.Log("测试, 批量关注成功")
}	


func TestUnfollowed(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}


	tests := []struct {
		userID       int64
		name         string
		followerID   int64
		followerName string
    }{
        {1000, "阿瓜", 1001, "炫狗"},
        {1000, "阿瓜", 1002, "孙笑川"},
		{1000, "阿瓜", 1003, "娟妹"},
        {1000, "阿瓜", 1004, "李老八"},
		
		{2000, "梓神", 1001, "炫狗"},
		{2000, "梓神", 1002, "萧瑟"},
		{1001, "炫狗", 2000, "梓神"},
    }

    for _, test := range tests {
        err := followStore.Unfollow(context.TODO(), test.followerID, test.userID)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(test.name, "取消关注", test.followerName)
    }
    
	t.Log("测试, 批量取消关注成功")
}	


func TestGetAllFollowed(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}

	followers, err := followStore.GetAllFollowed(context.TODO(), 1000)
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range followers {
		t.Logf("%+v\n", e)
	}

	t.Logf("%d\n", len(followers))
	t.Log("测试, 获取 gua 的关注列表, 成功")
}	


func TestGetAllFollowedCount(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}

	count, err := followStore.GetAllFollowedCount(context.TODO(), 1000)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("关注数, %d\n", count)
	t.Log("测试, 获取 gua 的关注数, 成功")
}	

func TestGetAllFollower(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}
	followers, err := followStore.GetAllFollower(context.TODO(), 1000)
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range followers {
		t.Logf("%+v\n", e)
	}

	t.Logf("%d\n", len(followers))
	t.Log("测试, 获取 gua 的粉丝列表, 成功")
}	

func TestGetAllFollowerCount(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}

	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	followStore := &FollowertoreImpl{db: db}
	
	count, err := followStore.GetAllFollowerCount(context.TODO(), 1001)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("粉丝数, %d\n", count)
	t.Log("测试, 获取 炫狗 的粉丝数, 成功")
}	