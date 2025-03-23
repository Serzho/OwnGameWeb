package main

import (
	"fmt"
	"log/slog"
	"os"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/handlers"
	"OwnGameWeb/internal/api/middleware"
	"OwnGameWeb/internal/api/routes"
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	cfg := config.Load()

	level := cfg.Global.LoggerLevel

	logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			slog.Error("Can't close log file", "err", err)
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

	err := os.MkdirAll(cfg.Global.CsvPath, 0o755)
	if err != nil {
		slog.Error("cannot create csv path", "path", cfg.Global.CsvPath, "err", err)
	}

	router.LoadHTMLGlob("./web/html/*.html")
	router.Static("static", "./web/static")
	router.Use(gin.Recovery(), middleware.Logger())

	dbController := database.NewDBController(cfg)
	defer dbController.Close()

	authService := services.NewAuthService(dbController, cfg)
	manageService := services.NewManageService(dbController, cfg)
	playService := services.NewPlayService(dbController)

	authHandler := handlers.NewAuthHandler(authService)
	manageHandler := handlers.NewManageHandler(manageService)
	playHandler := handlers.NewPlayHandler(playService)
	overviewHandler := handlers.NewOverviewHandler()

	routes.RegisterManageRoutes(router, manageHandler, middleware.Auth(cfg))
	routes.RegisterPlayRoutes(router, playHandler, middleware.Auth(cfg), middleware.Play())
	routes.RegisterAuthRoutes(router, authHandler)
	routes.RegisterOverviewRoutes(router, overviewHandler)

	url := fmt.Sprintf("%s:%d", cfg.Server.URL, cfg.Server.Port)
	slog.Info("Routes was mounted. Starting server...", "url", url)

	err = router.Run(url)
	if err != nil {
		panic(err)
	}
}
