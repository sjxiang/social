package main

import (
	"encoding/base64"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/sjxiang/social/internal/token"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		accessToken := parts[1]
		// 验证
		payload, err := app.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.User.GetByEmail(ctx, payload.Email)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		
		// 携带
		ctx = context.WithValue(ctx, userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
				return
			}

			// check the credentials
			username := app.config.Auth.Basic.Username
			password := app.config.Auth.Basic.Password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}


// 认证
func Authenticate() {

}

// 授权
func Authorize() {

}


// 限流
func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
			app.rateLimitExceededResponse(w, r, retryAfter.String())
			return
		}
	
		next.ServeHTTP(w, r)
	})
}