package handlers

import (
	"OwnGameWeb/internal/services"
	"fmt"
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
	email, password := c.PostForm("email"), c.PostForm("password")
	err := h.service.SignIn(email, password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	name, email, password := c.PostForm("name"), c.PostForm("email"), c.PostForm("password")
	err := h.service.SignUp(name, email, password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
}

func (h *AuthHandler) SignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func (h *AuthHandler) SignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{})
}

func (h *AuthHandler) RecoverPassword(c *gin.Context) {
	email := c.PostForm("email")
	err := h.service.RecoverPassword(email)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "successfully request",
	})

}

func (h *AuthHandler) RecoverPasswordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recovery.html", gin.H{})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
}
