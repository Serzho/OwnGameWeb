package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, _ *handlers.AuthHandler) *gin.RouterGroup {
	group := r.Group("/auth")
	{
		//users.GET("/:id", h.GetUser)
		//users.POST("/", h.CreateUser)
	}
	return group
}
