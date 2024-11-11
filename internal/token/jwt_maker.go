package token

import (
	"time"

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

	token, err := jwt.ParseWithClaims(accessToken, &Payload{}, keyFunc)
	if err != nil {
		return nil, err  // 内容被篡改（签名和密钥对不上）、过期
	}

	payload, ok := token.Claims.(*Payload)
	if !(ok && token.Valid) {
		return nil, err  // 类型断言失败
	}
	
	return payload, nil
}