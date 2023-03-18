package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Login struct {
	UserID int64  `db:"user_id"`
	Email  string `db:"email"`
}

func (l Login) SetClaimsForLogin() AccessTokenClaims {
	return AccessTokenClaims{
		UserID: l.UserID,
		Email:  l.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}

func (l Login) SetRefreshClaimsForLogin() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refresh",
		UserID:    l.UserID,
		Email:     l.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}
