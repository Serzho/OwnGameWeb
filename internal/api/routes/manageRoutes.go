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
	group.GET("/profile/info", h.ProfileInfo)
	group.GET("/getallpacks", h.GetAllPacks)
	group.GET("/getserverpacks", h.GetServerPacks)
	group.GET("/downloadpack/:id", h.DownloadPack)
	group.GET("/getpack/:id", h.GetPack)

	group.POST("/joingame", h.JoinGame)
	group.POST("/addpack", h.AddPack)
	group.POST("/addserverpack/:id", h.AddServerPack)

	group.PUT("/profile/update", h.UpdateProfile)
	group.PUT("/updatepacktitle/:id", h.UpdatePackTitle)
	group.PUT("/updatepackfile/:id", h.UpdatePackContent)

	group.DELETE("/deletepack/:id", h.DeletePack)

	return group
}
