package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"

	"github.com/gin-gonic/gin"
)

type ManageHandler struct {
	service *services.ManageService
}

func NewManageHandler(s *services.ManageService) *ManageHandler {
	return &ManageHandler{service: s}
}

func getIntFromJSON(jsonMap map[string]interface{}, key string) (int, error) {
	value, exists := jsonMap[key]
	if !exists {
		return -1, ErrKeyNotFoundInJSON
	}

	floatValue, ok := value.(float64)
	if !ok {
		return -1, ErrIncorrectType
	}

	return int(floatValue), nil
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
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("Error parsing json request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	code, ok := jsonMap["code"].(string)
	if !ok {
		slog.Warn("Error get code", "err", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("JoinGame", "code", code, "userID", userID)

	gameID, err := h.service.JoinGame(code, userID)
	if err != nil {
		slog.Warn("Error joining game", "code", code, "userID", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("JwtCreate", "userID", userID, "gameID", gameID)

	token, err := utils.JwtCreate(userID, gameID, h.service.Cfg.Global.SecretPhrase)
	if err != nil {
		slog.Warn("Error creating token", "userID", userID, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("Success create token", "userID", userID, "gameID", gameID)

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

	userID := c.GetInt("userID")

	slog.Info("AddPack", "userID", userID, "file", file, "header", header)

	err = h.service.AddPack(userID, file, header)
	if err != nil {
		slog.Warn("Error adding pack", "userID", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("Success adding pack", "userID", userID)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) CreateGame(c *gin.Context) {
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("Error parsing json request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	title, ok := jsonMap["title"].(string)
	if !ok {
		slog.Warn("Title is not string", "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	maxPlayers, err := getIntFromJSON(jsonMap, "maxPlayers")
	if err != nil {
		slog.Warn("Error parsing maxPlayers", "error", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	packID, err := getIntFromJSON(jsonMap, "packID")
	if err != nil {
		slog.Warn("Error parsing packID", "error", err, "map", jsonMap)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	userID := c.GetInt("userID")

	slog.Info("CreateGame", "title", title, "maxPlayers", maxPlayers, "packID", packID, "userID", userID)

	gameID, err := h.service.CreateGame(userID, packID, title, maxPlayers)
	if err != nil {
		slog.Warn("Error creating game", "userID", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("JwtCreate", "userID", userID, "gameID", gameID)

	token, err := utils.JwtCreate(userID, gameID, h.service.Cfg.Global.SecretPhrase)
	if err != nil {
		slog.Warn("Error creating token", "userID", userID, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	slog.Info("Success create game", "userID", userID, "gameID", gameID)
	c.SetCookie("token", token, 60*60*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func (h *ManageHandler) GetAllPacks(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("GetAllPacks", "userID", userID)

	packs, err := h.service.GetAllPacks(userID)
	if err != nil {
		slog.Warn("Error getting all packs", "userID", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not get all packs",
		})

		return
	}

	slog.Info("Success get all packs", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "packs": packs})
}

func (h *ManageHandler) DownloadPack(c *gin.Context) {
	packID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Error parsing id", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})

		return
	}

	slog.Info("GetPackFile", "id", packID)

	filepath, err := h.service.GetPackFile(packID)
	if err != nil {
		slog.Warn("Error getting pack file", "id", packID, "error", err)
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
	c.Header("Accept-Length", strconv.Itoa(len(content)))

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
	packID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Warn("Error parsing id", "packID", c.Param("id"), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "pack id is invalid",
		})

		return
	}

	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("DeletePack", "id", packID, "userID", userID)

	err = h.service.DeletePack(userID, packID)
	if err != nil {
		slog.Warn("Error deleting pack", "id", packID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "can not delete pack",
		})
	}

	slog.Info("Successfully deleted pack", "id", packID)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
