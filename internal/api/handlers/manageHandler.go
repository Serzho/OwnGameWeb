package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ManageHandler struct {
	service *services.ManageService
}

func NewManageHandler(s *services.ManageService) *ManageHandler {
	return &ManageHandler{service: s}
}

func (h *ManageHandler) CreateGamePage(c *gin.Context) {
	c.HTML(http.StatusOK, "creategame.html", gin.H{})
}

func (h *ManageHandler) JoinGamePage(c *gin.Context) {
	c.HTML(http.StatusOK, "joingame.html", gin.H{})
}

func (h *ManageHandler) MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", gin.H{})
}

func (h *ManageHandler) PackEditorPage(c *gin.Context) {
	c.HTML(http.StatusOK, "packeditor.html", gin.H{})
}

func (h *ManageHandler) ProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", gin.H{})
}

func (h *ManageHandler) JoinGame(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	code := jsonMap["code"].(string)
	err = h.service.JoinGame(code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
