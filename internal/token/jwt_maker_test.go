package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/sjxiang/social/internal/utils"
	
)


func TestJWTMaker(t *testing.T) {
	maker := NewJWTMaker(utils.RandomString(32), "gua")
	
	want := Payload{
		Email: utils.RandomEmail(),
		Role:  utils.RandomRole(),
	}

	token, err := maker.CreateToken(want, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	
	have, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, have)

	require.Equal(t, have.Email, want.Email)
	require.Equal(t, have.Role, want.Role)
}

func TestExpiredJWTToken(t *testing.T) {
	maker := NewJWTMaker(utils.RandomString(32), "gua")
	
	want := Payload{
		Email: utils.RandomEmail(),
		Role:  utils.RandomRole(),
	}

	token, err := maker.CreateToken(want, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	have, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, have)
}
