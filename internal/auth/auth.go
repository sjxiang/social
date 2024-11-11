package auth

import (
	"errors"
	"time"
)

var (
	ErrTokenExpiry = errors.New("token is expired")
)

type Authenticator interface {
	GenerateToken(userID int64, duration time.Duration) (string, error)
	ValidateToken(AccessToken string) (int64, error)
}