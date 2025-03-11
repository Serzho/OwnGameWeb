package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterManageRoutes(r *gin.Engine, h *handlers.ManageHandler) *gin.RouterGroup {
	group := r.Group("/")

	group.GET("/creategame", h.CreateGamePage)
	group.POST("/creategame", h.CreateGame)
	group.GET("/joingame", h.JoinGamePage)
	group.GET("/main", h.MainPage)
	group.GET("/packeditor", h.PackEditorPage)
	group.GET("/profile", h.ProfilePage)

	group.POST("/joingame", h.JoinGame)

	return group
}
