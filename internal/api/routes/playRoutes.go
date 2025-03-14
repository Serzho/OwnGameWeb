package routes

import (
	"OwnGameWeb/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPlayRoutes(r *gin.Engine, _ *handlers.PlayHandler, m gin.HandlerFunc) *gin.RouterGroup {
	group := r.Group("/play")
	group.Use(m)

	return group
}
