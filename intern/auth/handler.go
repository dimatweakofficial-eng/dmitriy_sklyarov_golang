package auth

import (
	"demo-1/configs"
	"demo-1/pkg/handle"
	"demo-1/pkg/jwt"
	"demo-1/pkg/res"
	"net/http"
)

type AuthHandler struct {
	Config      *configs.Config
	AuthServise *AuthServise
}

type AuthHandlerWithDeps struct {
	*configs.Config
	*AuthServise
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerWithDeps) {
	handler := AuthHandler{
		Config:      deps.Config,
		AuthServise: deps.AuthServise,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := handle.HandleBody[LoginRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthServise.Login(payload.Email, payload.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJwt(handler.Config.Token.TokenName).Create(jwt.JwtData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := handle.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthServise.Register(payload.Email, payload.Password, payload.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJwt(handler.Config.Token.TokenName).Create(jwt.JwtData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}
}
