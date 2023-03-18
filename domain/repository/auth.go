package repository

import (
	"github.com/foxposts/auth-api/adapter/context"
	"github.com/foxposts/auth-api/domain/service"
)

type AuthRepository interface {
	FindBy(requestEmail string, requestPassword string) (*service.Login, *context.ApiError)
	Save(userID int64, accessToken string, refreshToken string) *context.ApiError
	FetchAccessToken(userID int64) (string, *context.ApiError)
	FetchRefreshToken(userID int64) (string, *context.ApiError)
	DeleteJWT(userID int64) *context.ApiError
	FindByID(userID int64) (*service.Refresh, *context.ApiError)
}
