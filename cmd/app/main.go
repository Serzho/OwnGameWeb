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
	"log/slog"
	"os"
)

func main() {
	router := gin.New()
	cfg := config.Load()

	var level slog.Level
	switch cfg.Global.LoggerLevel {
	case -4:
		level = slog.LevelDebug
	case 0:
		level = slog.LevelInfo
	case 4:
		level = slog.LevelWarn
	case 8:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println("cannot close log file")
		}
	}(logFile)

	logger := slog.New(
		slog.NewJSONHandler(
			logFile,
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)
	slog.SetDefault(logger)

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

	url := fmt.Sprintf("%s:%d", cfg.Server.Url, cfg.Server.Port)
	slog.Info("Routes was mounted. Starting server...", "url", url)

	err = router.Run(url)
	if err != nil {
		panic(err)
	}
}
