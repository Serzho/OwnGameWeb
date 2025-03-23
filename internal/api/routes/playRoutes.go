package routes

import (
	"OwnGameWeb/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterPlayRoutes(r *gin.Engine, h *handlers.PlayHandler, m ...gin.HandlerFunc) *gin.RouterGroup {
	group := r.Group("/play")
	group.Use(m...)

	group.GET("/waitingroom", h.WaitingRoomPage)
	group.GET("/playerroom", h.PlayerRoomPage)
	group.GET("/masterroom", h.MasterRoomPage)
	group.GET("/gameinfo", h.GameInfo)
	group.GET("/questions", h.GetQuestions)

	group.POST("/start", h.StartGame)
	group.POST("/leave", h.LeaveGame)

	group.DELETE("/game", h.CancelGame)
	group.DELETE("/removeplayer", h.RemovePlayer)

	return group
}
