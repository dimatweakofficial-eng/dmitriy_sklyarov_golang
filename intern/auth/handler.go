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
	//хэндлер принимающий все нужные зависимости для работы с его методами
	handler := AuthHandler{
		Config:      deps.Config,
		AuthServise: deps.AuthServise,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//запрос в дто (заносим данные json запроса в переменную при успешной валидации)
		payload, err := handle.HandleBody[LoginRequest](w, r)
		if err != nil {
			return
		}
		//через сервис бизнес логика запроса (ищем модель в бд, если логин есть и пароль совпадает возвращаем email)
		email, err := handler.AuthServise.Login(payload.Email, payload.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		//генерим jwt для защищенных роутов, проверки авторизации и доступа к своим ссылкам (емаил аккаунта совпадает значит и токен совпадет с тем что при регестрации)
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
		//запрос в дто
		payload, err := handle.HandleBody[RegisterRequest](w, r)
		if err != nil {
			return
		}
		//через сервис бизнес логика запроса (нет ли уже такого емаил, если нет добавляем в бд шифруя пароль)
		email, err := handler.AuthServise.Register(payload.Email, payload.Password, payload.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		//генерим jwt для защищенных роутов c линком, проверка авторизации,  доступа к своим ссылкам (емаил аккаунта совпадает значит и токен совпадет с тем что при регестрации)
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
