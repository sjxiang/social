package main

import (
	"fmt"
	"net/http"

	"github.com/sjxiang/social/internal/data"
)

type ctxKey string


const (
	ctxKeyUser ctxKey = "user"
	ctxKeyRole ctxKey = "role"
	ctxKeyPost ctxKey = "post"
)

// 从请求中获取用户
func getUserFromContext(r *http.Request) (data.User, error) {
	v, ok := r.Context().Value(ctxKeyUser).(data.User)

	if !ok {
		return data.User{}, fmt.Errorf("user not found")  // 类型断言失败
	}

	return v, nil 
}


// var RoleSet = []Role{
// 	{	
// 		ID: 1,
// 		Name: "admin",
// 		Level: 4,
// 		Description: "管理员, 可以修改和删掉阿婆主的动态",
// 	},
// 	{
// 		ID: 2,
// 		Name: "moderator",
// 		Level: 3,
// 		Description: "版主, 可以修改阿婆主的动态",
// 	},
// 	{
// 		ID: 3,
// 		Name: "user",
// 		Level: 2,
// 		Description: "用户, 可以发动态和评论",
// 	},
// 	{
// 		ID: 4,
// 		Name: "guest",
// 		Level: 1,
// 		Description: "游客, 只能浏览动态",
// 	},
// }

var roleMap = map[string]int64{
	"admin": 4,
	"moderator": 3,
	"user": 2,
	"guest": 1,
}

var mapRole = map[int64]string{
	4: "admin",
	3: "moderator",
	2: "user",
	1: "guest",
}

