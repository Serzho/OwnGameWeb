package routes

import (
	"OwnGameWeb/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterOverviewRoutes(r *gin.Engine, h *handlers.OverviewHandler) *gin.RouterGroup {
	group := r.Group("/")

	group.GET("/", h.IndexPage)

	return group
}
