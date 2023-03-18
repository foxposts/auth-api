package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Refresh struct {
	UserID int64  `db:"user_id"`
	Email  string `db:"email"`
}

func (l Refresh) SetClaimsForRefresh() AccessTokenClaims {
	return AccessTokenClaims{
		UserID: l.UserID,
		Email:  l.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}

func (l Refresh) SetRefreshClaimsForRefresh() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refresh",
		UserID:    l.UserID,
		Email:     l.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}
