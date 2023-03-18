package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/foxposts/auth-api/adapter/context"
	"os"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *context.ApiError) {
	signedString, err := t.token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", context.UnexpectedError("cannot create access token")
	}
	return signedString, nil
}

func (t AuthToken) NewRefreshToken() (string, *context.ApiError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshClaims := c.GenerateRefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", context.UnexpectedError("cannot create refresh token")
	}
	return signedString, nil
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}
