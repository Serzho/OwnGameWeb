package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterOverviewRoutes(r *gin.Engine) *gin.RouterGroup {
	group := r.Group("/auth")
	{
		//users.GET("/:id", h.GetUser)
		//users.POST("/", h.CreateUser)
	}
	return group
}
