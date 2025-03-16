package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPlayRoutes(r *gin.Engine, h *handlers.PlayHandler, m gin.HandlerFunc) *gin.RouterGroup {
	group := r.Group("/play")
	group.Use(m)

	group.GET("/waitingroom", h.WaitingRoomPage)
	group.GET("/gameinfo", h.GameInfo)

	group.POST("/start", h.StartGame)

	group.DELETE("/game", h.CancelGame)
	group.DELETE("/removeplayer", h.RemovePlayer)

	return group
}
