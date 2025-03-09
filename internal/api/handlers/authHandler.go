package handlers

import (
	"OwnGameWeb/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) SignUp(_ *gin.Context) {}

func (h *AuthHandler) SignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func (h *AuthHandler) SignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{})
}

func (h *AuthHandler) RecoverPassword(_ *gin.Context) {
}

func (h *AuthHandler) RecoverPasswordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recovery.html", gin.H{})
}
