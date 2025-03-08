package main

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/handlers"
	"OwnGameWeb/internal/api/middleware"
	"OwnGameWeb/internal/api/routes"
	"OwnGameWeb/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	_ = config.Load()

	router.Use(gin.Recovery(), middleware.Logger(), middleware.Auth())

	authService := services.NewAuthService()
	manageService := services.NewManageService()
	playService := services.NewPlayService()

	authHandler := handlers.NewAuthHandler(authService)
	manageHandler := handlers.NewManageHandler(manageService)
	playHandler := handlers.NewPlayHandler(playService)

	routes.RegisterManageRoutes(router, manageHandler)
	routes.RegisterPlayRoutes(router, playHandler)
	routes.RegisterAuthRoutes(router, authHandler)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
