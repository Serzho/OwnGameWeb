package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterManageRoutes(r *gin.Engine, h *handlers.ManageHandler, m gin.HandlerFunc) *gin.RouterGroup {
	group := r.Group("/")
	group.Use(m)

	group.GET("/creategame", h.CreateGamePage)
	group.POST("/creategame", h.CreateGame)
	group.GET("/joingame", h.JoinGamePage)
	group.GET("/main", h.MainPage)
	group.GET("/packeditor", h.PackEditorPage)
	group.GET("/profile", h.ProfilePage)
	group.GET("/getallpacks", h.GetAllPacks)
	group.GET("/downloadpack/:id", h.DownloadPack)

	group.POST("/joingame", h.JoinGame)
	group.POST("/addpack", h.AddPack)

	group.DELETE("/deletepack/id", h.DeletePack)

	return group
}
