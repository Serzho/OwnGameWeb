package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		slog.Warn("SignIn: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	email, password := jsonMap["email"].(string), jsonMap["password"].(string)

	slog.Info("SignIn", "email", email, "password", password)
	userId, err := h.service.SignIn(email, password)

	if err != nil {
		slog.Warn("SignIn: Failed to sign in", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// TODO: сделать поиск игры
	gameId := -1
	slog.Info("Creating jwt token", "userId", userId, "gameId", gameId)
	token, err := utils.JwtCreate(userId, gameId, h.service.Cfg.Global.SecretPhrase)

	if err != nil {
		slog.Warn("SignIn: Failed to create jwt token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("Successfully signed in", "userId", userId)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		slog.Warn("SignUp: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	name, email, password := jsonMap["name"].(string), jsonMap["email"].(string), jsonMap["password"].(string)
	slog.Info("SignUp", "name", name, "email", email)
	err = h.service.SignUp(name, email, password)

	if err != nil {
		slog.Warn("SignUp: Failed to sign in", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("Successfully signed up", "name", name, "email", email)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *AuthHandler) SignUpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func (h *AuthHandler) SignInPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", gin.H{})
}

func (h *AuthHandler) RecoverPassword(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		slog.Warn("RecoverPassword: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	email := jsonMap["email"].(string)
	slog.Info("RecoverPassword", "email", email)
	err = h.service.RecoverPassword(email)

	if err != nil {
		slog.Warn("RecoverPassword: Failed to recover password", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("Successfully recovered password", "email", email)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *AuthHandler) RecoverPasswordPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recovery.html", gin.H{})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
}
