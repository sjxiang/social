package token

import (
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)


var ErrExpiredToken = errors.New("token is expired")


type Maker interface {
	CreateToken(payload Payload, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}


type Payload struct {
	jwt.RegisteredClaims   

	// 不要放敏感信息
	Email  string    `json:"email"`
	Role   string    `json:"role"`
}



/*

	例如
		'Authenticator' 
		'Auth2Claims'


	实现, 通过组合的方式

 */