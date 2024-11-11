package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey string
	iss       string  // 发行者
}

func NewJWTAuthenticator(secret, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss}
}



func (a *JWTAuthenticator) GenerateToken(userID int64, duration time.Duration) (string, error) {
	
	expires := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"id":      userID,
		"expires": expires,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token with secret, %v", err)
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(AccessToken string) (int64, error) {

	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])  // 签名算法不对
		}

		return []byte(a.secretKey), nil
	}

	token, err := jwt.Parse(AccessToken, keyFunc)
	if err!= nil {
		return 0, fmt.Errorf("failed to parse token, %v", err)  // 内容被篡改（签名和密钥对不上）
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if!ok {
		return 0, fmt.Errorf("failed to get claims from token")  // 类型断言失败
	}

	expires, ok := claims["expires"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to get expires from claims")
	}

	userID, ok := claims["id"].(int64)
	if !ok {
		return 0, fmt.Errorf("failed to get userID from claims")
	}

	if time.Now().Unix() > int64(expires) {
		return 0, ErrTokenExpiry // 过期
	}

	return userID, nil
}