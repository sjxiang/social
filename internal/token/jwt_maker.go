package token

import (
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string  // 密钥
	iss       string  // 发行方
}

func NewJWTMaker(secretKey, iss string) Maker {
	return &JWTMaker{secretKey: secretKey, iss: iss}
}

func (maker *JWTMaker) CreateToken(payload Payload, duration time.Duration) (string, error) {
	
	claims := &Payload{
		Email: payload.Email,
		Role:  payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    maker.iss,
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(duration)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(accessToken string) (*Payload, error) {
	
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}

	payload := &Payload{}

	_, err := jwt.ParseWithClaims(accessToken, payload, keyFunc)
	if err != nil {
		if strings.HasPrefix(err.Error(), "token has invalid claims: token is expired") {
			return nil,  ErrTokenExpiry   // 过期
		}

		return nil, err  // 内容被篡改; 例, 签名和密钥对不上...
	}

	return payload, nil
}
