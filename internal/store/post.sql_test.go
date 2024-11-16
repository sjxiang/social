package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Tests(t *testing.T) {
	tests := []struct {
		name     string
		input    string
        expected string
    }{
        {"正常", "2024-11-15 17:18:00", "2024-11-15 17:18:00"},
        {"异常", "invalid", ""},
    }

    for _, test := range tests {
        result := parseTime(test.input)
        if result != test.expected {
            t.Errorf("Expected %s, got %s", test.expected, result)
        }
    }
    
	t.Logf("\nParseTime, 测试完成")
}


func TestCreatePost(t *testing.T) {

    if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	postStore := newPostStore(db)

    params := Post{
        Title:    "T1 win",
        Content:  "shit ...",
        Tags:     []string{"lol", "live", "game"},
        UserID:    7,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Version:   1,
    }
    if err := postStore.Create(context.TODO(), params); err != nil {
        t.Log("数据库繁忙")
        t.Fatal(err)
    }

    t.Log("插入帖子成功")
}


func TestGetOnePost(t *testing.T) {

    if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	postStore := newPostStore(db)

    post, err := postStore.GetOne(context.TODO(), 2)
    require.NoError(t, err)

    t.Logf("%+v\n", post)
}

func TestUpdatePost(t *testing.T) {

    if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
    
    postStore := newPostStore(db)
    
    if err := postStore.Update(context.TODO(), Post{
        Title:     "DK Fail",
        Content:   "shit ...",
        UpdatedAt: time.Now(),
        ID:        2,
        Version:   1,
    }); err != nil {
        t.Fatal(err)
    }

    t.Log("更新帖子成功")
}

func TestDeletePost(t *testing.T) {

    if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
    
    postStore := newPostStore(db)
    
    if err := postStore.Delete(context.Background(), 2); err != nil {
        t.Fatal(err)
    }

    t.Log("删除成功")
}


func TestGetUserFeed(t *testing.T) {

}