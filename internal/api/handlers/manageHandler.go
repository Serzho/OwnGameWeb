package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
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

func (h *ManageHandler) AddPack(c *gin.Context) {
	file, header, err := c.Request.FormFile("packFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "csv file is required",
		})
		return
	}

	userId := c.GetInt("userId")

	err = h.service.AddPack(userId, file, header)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) CreateGame(c *gin.Context) {

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "user id not found",
		})
		return
	}

	maxPlayers, err := strconv.Atoi(c.PostForm("maxPlayers"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "maxPlayers is required",
		})
		return
	}
	title := c.PostForm("title")

	gameId, err := h.service.CreateGame(userId.(int), "", title, maxPlayers)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "gameId": gameId})
}

func (h *ManageHandler) GetAllPacks(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "user id not found",
		})
		return
	}

	packs, err := h.service.GetAllPacks(userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not get all packs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "packs": packs})
}

func (h *ManageHandler) DownloadPack(c *gin.Context) {
	packId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})
		return
	}

	filepath, err := h.service.GetPackFile(packId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot get pack file",
		})
		return
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot read pack file",
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=pack")
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	_, err = c.Writer.Write(content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot write pack file",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Download file successfully",
	})
}

func (h *ManageHandler) DeletePack(c *gin.Context) {
	packId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "user id not found",
		})
		return
	}

	err = h.service.DeletePack(userId.(int), packId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not delete pack",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
