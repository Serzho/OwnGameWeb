package routes

import (
	"OwnGameWeb/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, h *handlers.AuthHandler) *gin.RouterGroup {
	group := r.Group("/auth")

	group.GET("/signin", h.SignInPage)
	group.GET("/signup", h.SignUpPage)
	group.GET("/recoverPassword", h.RecoverPasswordPage)

	group.POST("/signin", h.SignIn)
	group.POST("/signup", h.SignUp)
	group.POST("/recoverPassword", h.RecoverPassword)
	group.POST("/signout", h.SignOut)

	return group
}
