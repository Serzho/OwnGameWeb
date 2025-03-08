package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterManageRoutes(r *gin.Engine, _ *handlers.ManageHandler) {
	_ = r.Group("/")
	{
		//users.GET("/:id", h.GetUser)
		//users.POST("/", h.CreateUser)
	}
}
