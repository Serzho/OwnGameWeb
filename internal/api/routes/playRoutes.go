package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPlayRoutes(r *gin.Engine, _ *handlers.PlayHandler) *gin.RouterGroup {
	group := r.Group("/play")
	{
		//users.GET("/:id", h.GetUser)
		//users.POST("/", h.CreateUser)
	}
	return group
}
