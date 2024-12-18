package main

import (
	"fmt"
	"net/http"

	"github.com/sjxiang/social/internal/store"
)

type ctxKey string


const (
	userKey ctxKey = "user"
	roleKey ctxKey = "role"
	postKey ctxKey = "post"
)

// 从请求中获取用户(经过拦截器, 肯定有用户)

func getUserFromContext(r *http.Request) (*store.User, error) {
	v, ok := r.Context().Value(userKey).(*store.User)

	if !ok {
		return nil, fmt.Errorf("user not found")  // 类型断言失败
	}

	return v, nil 
}


func getPostFromContext(r *http.Request) (*store.Post, error) {
	v, exists := r.Context().Value(postKey).(*store.Post)

	if !exists {
		return nil, fmt.Errorf("post not found")
	}

	return v, nil
}


// var roleMap = map[int64][]string{
// 	"4: []string{"admin", "管理员", "可以删贴, 也可以补充说明"},
// 	3: []string{"moderator", "版主", "可以对帖子补充说明"},
// 	2: []string{"user", "用户", "可以发帖和评论"},
// 	1: []string{"guest", "访客", "只能浏览"},
// }

