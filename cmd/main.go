package main

import (
	"demo-1/configs"
	"demo-1/intern/auth"
	"demo-1/intern/link"
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

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerWithDeps{
		Config: config,
	})
	link.NewLinkHandler(router, link.LinkHandlerWithDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Handler: middleware.Logging(router),
		Addr:    ":8081",
	}

	fmt.Println("Server is listening on 8081 port")
	server.ListenAndServe()
}
