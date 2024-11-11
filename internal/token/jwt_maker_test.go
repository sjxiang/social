package token

import (
	"time"
	"errors"
	"testing"
)

func TestCreateToken(t *testing.T) {
	payload := Payload{
		Email: "gua@vip.cn",
		Role:  "admin",
	}

	maker := NewJWTMaker("8xEMrWkBARcDDYQ", "JUEJIN")

	token, err := maker.CreateToken(payload, time.Minute*3)
	if err!= nil {
		t.Fatal(err)
	}

	t.Log(token)
}

func TestVerifyToken(t *testing.T) {
	maker := NewJWTMaker("8xEMrWkBARcDDYQ", "JUEJIN")
	
	payload, err := maker.VerifyToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJKVUVKSU4iLCJleHAiOjE3MzEzNTg1MTAsImlhdCI6MTczMTM1ODMzMCwiZW1haWwiOiJndWFAdmlwLmNuIiwicm9sZSI6ImFkbWluIn0.IQSJ9IR0ELuGDfXwnvfr7O0bCUbLMKH7WxqxbNhOyrE")
	if errors.Is(err, ErrTokenExpiry) {
		t.Fatal("过期了")
	}
	if err!= nil {
		t.Fatal(err)
	} 
	
	t.Log(payload)
}