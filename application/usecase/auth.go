package usecase

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/foxposts/auth-api/adapter/context"
	"github.com/foxposts/auth-api/domain/repository"
	"github.com/foxposts/auth-api/domain/service"
	"os"
)

type AuthService interface {
	Login(request context.LoginRequest) (*context.LoginResponse, *context.ApiError)
	Verify(urlParams map[string]string) *context.ApiError
	Refresh(request context.RefreshTokenRequest) (*context.RefreshTokenResponse, *context.ApiError)
	Logout(request context.AccessTokenRequest) *context.ApiError
}

type authService struct {
	repo repository.AuthRepository
}

func NewUserUseCase(repo repository.AuthRepository) AuthService {
	return &authService{repo}
}

func (s authService) Login(req context.LoginRequest) (*context.LoginResponse, *context.ApiError) {
	var login *service.Login
	var apiErr *context.ApiError
	login, apiErr = s.repo.FindBy(req.Email, req.Password)
	if apiErr != nil {
		return nil, apiErr
	}

	claims := login.SetClaimsForLogin()
	authToken := service.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, apiErr = authToken.NewAccessToken(); apiErr != nil {
		return nil, apiErr
	}
	if refreshToken, apiErr = authToken.NewRefreshToken(); apiErr != nil {
		return nil, apiErr
	}
	if apiErr := s.repo.Save(login.UserID, accessToken, refreshToken); apiErr != nil {
		return nil, apiErr
	}
	return &context.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s authService) Verify(urlParams map[string]string) *context.ApiError {

	jwtToken, err := jwtTokenFromString(urlParams["token"])
	if err != nil {
		return context.AuthenticationError("cannot convert string to JWT")
	}
	if jwtToken.Valid {
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return context.UnexpectedError("cannot extract claims")
		}

		userID := int64(claims["user_id"].(float64))
		if !ok {
			return context.UnexpectedError("cannot extract userID")
		}

		storeAccessToken, apiErr := s.repo.FetchAccessToken(userID)
		if apiErr != nil {
			return apiErr
		}
		if jwtToken.Raw == storeAccessToken {
			return nil
		}

	}
	return context.AuthenticationError("invalid token")
}

func (s authService) Refresh(request context.RefreshTokenRequest) (*context.RefreshTokenResponse, *context.ApiError) {
	jwtToken, err := jwtTokenFromString(request.RefreshToken)
	if err != nil {
		return nil, context.AuthenticationError("cannot convert string to JWT")
	}

	if jwtToken.Valid {
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, context.UnexpectedError("cannot extract claims")
		}

		userID := int64(claims["user_id"].(float64))
		if !ok {
			return nil, context.UnexpectedError("cannot extract userID")
		}

		storeRefreshToken, apiErr := s.repo.FetchRefreshToken(userID)
		if apiErr != nil {
			return nil, apiErr
		}
		if jwtToken.Raw != storeRefreshToken {
			return nil, context.AuthenticationError("not match stored refreshToken")
		}

		// consider changed email
		refresh, apiErr := s.repo.FindByID(userID)
		if apiErr != nil {
			return nil, apiErr
		}
		refreshClaims := refresh.SetClaimsForRefresh()
		authToken := service.NewAuthToken(refreshClaims)

		var accessToken, refreshToken string
		if accessToken, apiErr = authToken.NewAccessToken(); apiErr != nil {
			return nil, apiErr
		}
		if refreshToken, apiErr = authToken.NewRefreshToken(); apiErr != nil {
			return nil, apiErr
		}
		if apiErr := s.repo.Save(refresh.UserID, accessToken, refreshToken); apiErr != nil {
			return nil, apiErr
		}
		return &context.RefreshTokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

	}
	return nil, context.AuthenticationError("invalid token")
}

func (s authService) Logout(request context.AccessTokenRequest) *context.ApiError {
	jwtToken, err := jwtTokenFromString(request.AccessToken)
	if err != nil {
		return context.UnexpectedError("cannot convert string to JWT")
	}

	if jwtToken.Valid {
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return context.UnexpectedError("cannot extract claims")
		}

		userID := int64(claims["user_id"].(float64))
		if !ok {
			return context.UnexpectedError("cannot extract userID")
		}
		storeAccessToken, apiErr := s.repo.FetchAccessToken(userID)
		if apiErr != nil {
			return apiErr
		}
		if jwtToken.Raw != storeAccessToken {
			return context.AuthenticationError("not match stored accessToken")
		}

		apiEr := s.repo.DeleteJWT(userID)
		if apiEr != nil {
			return apiEr
		}

		return nil

	}
	return context.AuthenticationError("invalid token")
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
