package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterManageRoutes(r *gin.Engine, _ *handlers.ManageHandler) *gin.RouterGroup {
	group := r.Group("/")
	{
		//users.GET("/:id", h.GetUser)
		//users.POST("/", h.CreateUser)
	}
	return group
}
