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
	config := configs.LoadConfig()
	db := db.NewDb(config)
	router := http.NewServeMux()

	//Repositoris
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerWithDeps{
		Config:      config,
		AuthServise: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerWithDeps{
		LinkRepository: linkRepository,
	})

	//Midlewairs
	stack := middleware.Chain(
		middleware.Cors,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on 8081 port")
	server.ListenAndServe()
}
