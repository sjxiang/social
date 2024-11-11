package main

import (
	"context"
	"net/http"

	"github.com/sjxiang/social/internal/db"
)

type ctxKey string

const (
	userKey ctxKey = "user"
	roleKey ctxKey = "role"
	postKey ctxKey = "post"
)

func setUserToContext(ctx context.Context, arg db.User) context.Context {
	return context.WithValue(ctx, userKey, arg)
}

func getUserFromContext(r *http.Request) db.User {
	v, ok := r.Context().Value(userKey).(db.User)
	if !ok {
		return db.User{}
	}

	return v
}
