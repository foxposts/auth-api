package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const ACCESS_TOKEN_EXPIRATION_DURATION = time.Hour
const REFRESH_TOKEN_EXPIRATION_DURATION = time.Hour * 24 * 30

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

type AccessTokenClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func (c RefreshTokenClaims) GenerateAccessTokenClaims() AccessTokenClaims {
	return AccessTokenClaims{
		UserID: c.UserID,
		Email:  c.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}

func (c AccessTokenClaims) GenerateRefreshTokenClaims() RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: "refresh",
		UserID:    c.UserID,
		Email:     c.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_EXPIRATION_DURATION).Unix(),
		},
	}
}
