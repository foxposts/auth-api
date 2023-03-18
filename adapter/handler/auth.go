package handler

import (
	"encoding/json"
	"github.com/foxposts/auth-api/adapter/context"
	"github.com/foxposts/auth-api/application/usecase"
	"log"
	"net/http"
)

type AuthHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	VerifyHandler(w http.ResponseWriter, r *http.Request)
	RefreshHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	us usecase.AuthService
}

func NewAuthHandler(u usecase.AuthService) AuthHandler {
	return &authHandler{u}
}

func (a *authHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var l context.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			createResponse(w, http.StatusBadRequest, context.ApiError{Message: "invalid json: cannot decode json"})
		} else {
			token, apiErr := a.us.Login(l)
			if apiErr != nil {
				createResponse(w, apiErr.Code, apiErr.AsMessage())
			} else {
				createResponse(w, http.StatusOK, *token)
			}
		}
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(405)
	}
}

func (a *authHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		urlParams := make(map[string]string)
		for k := range r.URL.Query() {
			urlParams[k] = r.URL.Query().Get(k)
		}

		if urlParams["token"] != "" {
			apiErr := a.us.Verify(urlParams)
			if apiErr != nil {
				createResponse(w, apiErr.Code, apiErr.AsMessage())
			} else {
				createResponse(w, http.StatusOK, authorizedResponse())
			}
		} else {
			createResponse(w, http.StatusForbidden, notAuthorizedResponse("not found token"))
		}
	default:
		w.WriteHeader(405)
	}
}

func (a *authHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var refreshTokenRequest context.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
			createResponse(w, http.StatusBadRequest, context.ApiError{Message: "invalid json: cannot decode json"})
		} else {
			tokens, apiErr := a.us.Refresh(refreshTokenRequest)
			if apiErr != nil {
				createResponse(w, apiErr.Code, apiErr.AsMessage())
			} else {
				createResponse(w, http.StatusOK, *tokens)
			}
		}
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(405)
	}
}

func (a *authHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var accessTokenRequest context.AccessTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&accessTokenRequest); err != nil {
			createResponse(w, http.StatusBadRequest, context.ApiError{Message: "invalid json: cannot decode json"})
		}
		apiErr := a.us.Logout(accessTokenRequest)
		if apiErr != nil {
			createResponse(w, apiErr.Code, apiErr.AsMessage())
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(405)
	}
}

func createResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatalf("cannot encode response data")
	}
}

func notAuthorizedResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      message,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
