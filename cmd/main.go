package main

import (
	"demo-1/configs"
	"demo-1/intern/auth"
	"demo-1/intern/link"
	"demo-1/intern/user"
	"demo-1/pkg/db"
	"demo-1/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	//подключение к основным сущностям с которыми будет вестись работа
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()

	//Repositoris (действия с бд)
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	//Services (бизнес логика пришедших запросов после получения дто)
	authServise := auth.NewAuthService(userRepository)

	//Handlers (обработчики принимающие все необходимые зависимости для дальнеших вз)
	auth.NewAuthHandler(router, auth.AuthHandlerWithDeps{
		Config:      config,
		AuthServise: authServise,
	})
	link.NewLinkHandler(router, link.LinkHandlerWithDeps{
		Config:         config,
		LinkRepository: linkRepository,
	})

	//Midlewairs (cтэк промежуточных обработчиков)
	stack := middleware.Chain(
		middleware.Cors,
		middleware.Logging,
	)

	//сервер где каждый запрос проходит stack middleware
	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on 8081 port")
	server.ListenAndServe()
}
