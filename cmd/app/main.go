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

	router.Use(gin.Recovery(), middleware.Logger())

	authService := services.NewAuthService()
	manageService := services.NewManageService()
	playService := services.NewPlayService()

	authHandler := handlers.NewAuthHandler(authService)
	manageHandler := handlers.NewManageHandler(manageService)
	playHandler := handlers.NewPlayHandler(playService)
	overviewHandler := handlers.NewOverviewHandler()

	manageGroup := routes.RegisterManageRoutes(router, manageHandler)
	playGroup := routes.RegisterPlayRoutes(router, playHandler)
	routes.RegisterAuthRoutes(router, authHandler)
	routes.RegisterOverviewRoutes(router, overviewHandler)

	manageGroup.Use(middleware.Auth(cfg))
	playGroup.Use(middleware.Auth(cfg))

	err := router.Run(fmt.Sprintf("%s:%d", cfg.Server.Url, cfg.Server.Port))
	if err != nil {
		panic(err)
	}
}
