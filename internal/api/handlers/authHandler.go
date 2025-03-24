package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"OwnGameWeb/config"

	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
	Cfg     *config.Config
}

func NewAuthHandler(s services.AuthService, c *config.Config) *AuthHandler {
	return &AuthHandler{service: s, Cfg: c}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("SignIn: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	email, emailOk := jsonMap["email"].(string)
	password, passwordOk := jsonMap["password"].(string)

	if !emailOk || !passwordOk {
		slog.Warn("SignIn: Failed to parse json request")
	}

	slog.Info("SignIn", "email", email, "password", password)

	userID, err := h.service.SignIn(email, password)
	if err != nil {
		slog.Warn("SignIn: Failed to sign in", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	// TODO: сделать поиск игры
	gameID := -1
	slog.Info("Creating jwt token", "userID", userID, "gameID", gameID)

	token, err := utils.JwtCreate(userID, gameID, h.Cfg.Global.SecretPhrase)
	if err != nil {
		slog.Warn("SignIn: Failed to create jwt token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("Successfully signed in", "userID", userID)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("SignUp: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	name, ok := jsonMap["name"].(string)
	if !ok {
		slog.Warn("Error get name", "err", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	email, ok := jsonMap["email"].(string)
	if !ok {
		slog.Warn("Error get email", "err", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	password, ok := jsonMap["password"].(string)
	if !ok {
		slog.Warn("Error get password", "err", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

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
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("RecoverPassword: Failed to parse json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	email, ok := jsonMap["email"].(string)
	if !ok {
		slog.Warn("Error get email", "err", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

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
