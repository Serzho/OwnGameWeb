package main

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/handlers"
	"OwnGameWeb/internal/api/middleware"
	"OwnGameWeb/internal/api/routes"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	cfg := config.Load()

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

	err := router.Run(fmt.Sprintf("%s:%d", cfg.Server.Url, cfg.Server.Port))
	if err != nil {
		panic(err)
	}
}
