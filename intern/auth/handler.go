package auth

import (
	"demo-1/configs"
	"demo-1/pkg/handle"
	"demo-1/pkg/res"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	Config *configs.Config
}

type AuthHandlerWithDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerWithDeps) {
	handler := AuthHandler{
		Config: deps.Config,
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
		fmt.Println(payload)
		data := LoginResponse{
			Token: "1234",
		}
		res.Json(w, data, 200)
	}
}

func (*AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := handle.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		fmt.Println(payload)
	}
}
