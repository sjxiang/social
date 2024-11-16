package store

import (
	"context"
	"testing"
	"time"
)


func TestGetCommentByPostID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	commentStore := newCommentStore(db)
	
	comments, err := commentStore.GetByPostID(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}

	for i := range comments {
		t.Logf("%+v\n", comments[i].Content)		
	}

	t.Log("查看帖子下面所有评论")
}


func TestCreateComment(t *testing.T) {
	
	if testing.Short() {
		t.Skip()
	}
	
	db, err := initdb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	commentStore := newCommentStore(db)
	if err := commentStore.Create(context.TODO(), Comment{
		PostID: 3,
		UserID: 8,
		Content: "fuck",
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	t.Log("评论成功")
}