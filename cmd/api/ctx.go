package main

import (
	"context"
	"net/http"

	"github.com/sjxiang/social/internal/data"
)

type ctxKey string

const (
	userKey ctxKey = "user"
	roleKey ctxKey = "role"
)


func setUserToContext(ctx context.Context, arg data.User) context.Context {
	return context.WithValue(ctx, userKey, arg)
}

func getUserFromContext(r *http.Request) data.User {
	v, ok := r.Context().Value(userKey).(data.User)
	if !ok {
		return data.User{}
	}

	return v
}
