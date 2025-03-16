package main

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/handlers"
	"OwnGameWeb/internal/api/middleware"
	"OwnGameWeb/internal/api/routes"
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := gin.Default()
	cfg := config.Load()

	err := os.MkdirAll(cfg.Global.CsvPath, 0755)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	router.LoadHTMLGlob("./web/html/*.html")
	router.Static("static", "./web/static")
	router.Use(gin.Recovery(), middleware.Logger())

	dbController := database.NewDbController(cfg)

	defer dbController.Close()
	authService := services.NewAuthService(dbController, cfg)
	manageService := services.NewManageService(dbController, cfg)
	playService := services.NewPlayService(dbController)

	authHandler := handlers.NewAuthHandler(authService)
	manageHandler := handlers.NewManageHandler(manageService)
	playHandler := handlers.NewPlayHandler(playService)
	overviewHandler := handlers.NewOverviewHandler()

	routes.RegisterManageRoutes(router, manageHandler, middleware.Auth(cfg))
	routes.RegisterPlayRoutes(router, playHandler, middleware.Auth(cfg))
	routes.RegisterAuthRoutes(router, authHandler)
	routes.RegisterOverviewRoutes(router, overviewHandler)

	err = router.Run(fmt.Sprintf("%s:%d", cfg.Server.Url, cfg.Server.Port))
	if err != nil {
		panic(err)
	}
}
