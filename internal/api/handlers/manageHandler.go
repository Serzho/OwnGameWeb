package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
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
		slog.Warn("Error parsing json request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	userId, exists := c.Get("userId")

	if !exists {
		slog.Warn("UserId not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "user id not found",
		})
	}

	code := jsonMap["code"].(string)

	slog.Info("JoinGame", "code", code, "userId", userId)
	gameId, err := h.service.JoinGame(code, userId.(int))

	if err != nil {
		slog.Warn("Error joining game", "code", code, "userId", userId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("JwtCreate", "userId", userId.(int), "gameId", gameId)
	token, err := utils.JwtCreate(userId.(int), gameId, h.service.Cfg.Global.SecretPhrase)
	if err != nil {
		slog.Warn("Error creating token", "userId", userId, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	slog.Info("Success create token", "userId", userId.(int), "gameId", gameId)

	c.SetCookie("token", token, 60*60*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) AddPack(c *gin.Context) {
	file, header, err := c.Request.FormFile("packFile")
	if err != nil {
		slog.Warn("Error parsing form file", "file", file, "header", header, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "csv file is required",
		})
		return
	}

	userId := c.GetInt("userId")

	slog.Info("AddPack", "userId", userId, "file", file, "header", header)
	err = h.service.AddPack(userId, file, header)

	if err != nil {
		slog.Warn("Error adding pack", "userId", userId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("Success adding pack", "userId", userId)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) CreateGame(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		slog.Warn("Error parsing json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	title := jsonMap["title"].(string)
	maxPlayers := int(jsonMap["maxPlayers"].(float64))
	packId := int(jsonMap["packId"].(float64))

	userId := c.GetInt("userId")

	slog.Info("CreateGame", "title", title, "maxPlayers", maxPlayers, "packId", packId, "userId", userId)
	gameId, err := h.service.CreateGame(userId, packId, title, maxPlayers)

	if err != nil {
		slog.Warn("Error creating game", "userId", userId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("JwtCreate", "userId", userId, "gameId", gameId)
	token, err := utils.JwtCreate(userId, gameId, h.service.Cfg.Global.SecretPhrase)
	if err != nil {
		slog.Warn("Error creating token", "userId", userId, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	slog.Info("Success create game", "userId", userId, "gameId", gameId)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) GetAllPacks(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("UserId not found in context")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "user id not found",
		})
		return
	}

	slog.Info("GetAllPacks", "userId", userId.(int))
	packs, err := h.service.GetAllPacks(userId.(int))
	if err != nil {
		slog.Warn("Error getting all packs", "userId", userId.(int), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not get all packs",
		})
		return
	}

	slog.Info("Success get all packs", "userId", userId.(int))
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "packs": packs})
}

func (h *ManageHandler) DownloadPack(c *gin.Context) {
	packId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Error parsing id", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})
		return
	}

	slog.Info("GetPackFile", "id", packId)
	filepath, err := h.service.GetPackFile(packId)
	if err != nil {
		slog.Warn("Error getting pack file", "id", packId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot get pack file",
		})
		return
	}

	slog.Info("Reading file", "path", filepath)
	content, err := os.ReadFile(filepath)
	if err != nil {
		slog.Warn("Error reading file", "path", filepath, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot read pack file",
		})
		return
	}

	slog.Info("Writing file", "path", filepath)
	c.Header("Content-Disposition", "attachment; filename=pack")
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	_, err = c.Writer.Write(content)
	if err != nil {
		slog.Warn("Error writing file", "path", filepath, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot write pack file",
		})
	}

	slog.Info("Success writing file", "path", filepath)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Download file successfully",
	})
}

func (h *ManageHandler) DeletePack(c *gin.Context) {
	packId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Error parsing id", "packId", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("UserId not found in context")
		c.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "user id not found",
		})
		return
	}

	slog.Info("DeletePack", "id", packId, "userId", userId.(int))
	err = h.service.DeletePack(userId.(int), packId)

	if err != nil {
		slog.Warn("Error deleting pack", "id", packId, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not delete pack",
		})
	}

	slog.Info("Successfully deleted pack", "id", packId)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
