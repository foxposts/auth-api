package persistence

import (
	"context"
	"database/sql"
	"fmt"
	mycontext "github.com/foxposts/auth-api/adapter/context"
	"github.com/foxposts/auth-api/domain/model"
	"github.com/foxposts/auth-api/domain/repository"
	"github.com/foxposts/auth-api/domain/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type authPersistence struct {
	mysqlClient *sql.DB
	redisClient *redis.Client
}

func NewAuthRedisPersistence(mysqlClient *sql.DB, redisClient *redis.Client) repository.AuthRepository {
	return &authPersistence{mysqlClient, redisClient}
}

func (d authPersistence) FindBy(requestEmail string, requestPassword string) (*service.Login, *mycontext.ApiError) {
	var user model.User

	err := d.mysqlClient.QueryRow(`SELECT id, email, password FROM users WHERE email = ?`, requestEmail).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, mycontext.AuthenticationError("invalid login input")
		} else {
			mycontext.UnexpectedError(fmt.Sprintf("%v : unexpected while scanning users", err))
		}
	}
	er := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestPassword))
	if er != nil {
		return nil, mycontext.AuthenticationError("invalid login input")
	}
	return &service.Login{UserID: user.ID, Email: user.Email}, nil
}

var ctx = context.Background()

func (d authPersistence) Save(userID int64, accessToken string, refreshToken string) *mycontext.ApiError {
	refreshKey := "refresh-" + strconv.Itoa(int(userID))
	err := d.redisClient.Set(ctx, strconv.Itoa(int(userID)), accessToken, service.ACCESS_TOKEN_EXPIRATION_DURATION).Err()
	if err != nil {
		return mycontext.UnexpectedError(fmt.Sprintf("%v :unexpected while setting access token", err))
	}
	er := d.redisClient.Set(ctx, refreshKey, refreshToken, service.REFRESH_TOKEN_EXPIRATION_DURATION).Err()
	if er != nil {
		return mycontext.UnexpectedError(fmt.Sprintf("%v : unexpected while setting refresh token", err))
	}

	return nil
}

func (d authPersistence) FetchAccessToken(userID int64) (string, *mycontext.ApiError) {
	accessToken, err := d.redisClient.Get(ctx, strconv.Itoa(int(userID))).Result()
	if err == redis.Nil {
		return "", mycontext.NotFoundError("no such access token")
	}
	if err != nil {
		return "", mycontext.NotFoundError(fmt.Sprintf("%v : unexpected while scanning accessToken", err))
	}
	return accessToken, nil
}

func (d authPersistence) FetchRefreshToken(userID int64) (string, *mycontext.ApiError) {
	refreshKey := "refresh-" + strconv.Itoa(int(userID))
	refreshToken, err := d.redisClient.Get(ctx, refreshKey).Result()
	if err == redis.Nil {
		return "", mycontext.NotFoundError("no such refresh token")
	}
	if err != nil {
		return "", mycontext.NotFoundError(fmt.Sprintf("%v : unexpected while scanning accessToken", err))
	}
	return refreshToken, nil
}

func (d authPersistence) DeleteJWT(userID int64) *mycontext.ApiError {
	ar, err := d.redisClient.Del(ctx, strconv.Itoa(int(userID))).Result()
	if err != nil {
		return mycontext.UnexpectedError(fmt.Sprintf("%v : unexpected while deleting access token", err))
	}
	if ar != 1 {
		mycontext.StatusError("expired access token")
	}
	refreshKey := "refresh-" + strconv.Itoa(int(userID))
	rr, err := d.redisClient.Del(ctx, refreshKey).Result()
	if err != nil {
		return mycontext.UnexpectedError(fmt.Sprintf("%v :unexpected while deleting refresh token", err))
	}
	if rr != 1 {
		return mycontext.StatusError("expired refresh token")
	}
	return nil
}

func (d authPersistence) FindByID(userID int64) (*service.Refresh, *mycontext.ApiError) {
	var id int64
	var email string
	err := d.mysqlClient.QueryRow("SELECT id, email  FROM users WHERE id = ?", userID).Scan(&id, &email)
	if err == sql.ErrNoRows {
		return nil, mycontext.NotFoundError("invalid login input")
	} else {
		mycontext.UnexpectedError(fmt.Sprintf("%v : unexpected while scanning users", err))
	}
	return &service.Refresh{UserID: id, Email: email}, nil
}
