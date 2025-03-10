package handlers

import (
	"OwnGameWeb/internal/services"
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
