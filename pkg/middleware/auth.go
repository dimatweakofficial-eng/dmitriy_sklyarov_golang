package middleware

import (
	"context"
	"demo-1/configs"
	"demo-1/pkg/jwt"
	"net/http"
	"strings"
)

type Key string

const (
	ContextEmailKey Key = "ContextEmailKey"
)

func WriteUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer") {
			WriteUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(authedHeader, "Bearer ")

		//достаем ключ доступа и с помощью него достаем емаил из payload
		jwtSecret := jwt.NewJwt(config.Token.TokenName)
		isValid, jwtData := jwtSecret.Parse(token)
		if !isValid {
			WriteUnauthorized(w)
			return
		}
		//создаем контекст для передачи информации о емаил следующим обработчикам
		ctx := context.WithValue(r.Context(), ContextEmailKey, jwtData.Email)
		//модифицируем передаваемый запрос новым контекстом
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
