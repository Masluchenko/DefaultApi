package main

import (
	"demo/api/configs"
	"demo/api/internal/auth"
	"demo/api/internal/link"
	"demo/api/internal/user"
	"demo/api/pkg/db"
	"demo/api/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	//Service
	authService := auth.NewAuthService(userRepository)

	//Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	//Middleares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.IsAuthed,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}
	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
