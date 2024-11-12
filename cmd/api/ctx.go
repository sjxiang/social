package main

import (
	"fmt"
	"net/http"

	"github.com/sjxiang/social/internal/data"
)

type ctxKey string

const (
	userKey ctxKey = "user"
	roleKey ctxKey = "role"
)


func getUserFromContext(r *http.Request) (data.User, error) {
	v, ok := r.Context().Value(userKey).(data.User)

	if !ok {
		return data.User{}, fmt.Errorf("user not found")  // 类型断言失败
	}

	return v, nil 
}
