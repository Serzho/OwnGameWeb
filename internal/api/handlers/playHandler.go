package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type PlayHandler struct {
	service *services.PlayService
}

func NewPlayHandler(s *services.PlayService) *PlayHandler {
	return &PlayHandler{service: s}
}

func (h *PlayHandler) WaitingRoomPage(c *gin.Context) {
	_, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	c.HTML(http.StatusOK, "waitingroom.html", gin.H{})
}

func (h *PlayHandler) GameInfo(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("userId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	slog.Info("GetGameInfo", "gameId", gameId.(int), "userId", userId.(int))
	gameInfo, err := h.service.GetGameInfo(gameId.(int), userId.(int))
	if err != nil {
		slog.Warn("Error GetGameInfo", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot get game info",
		})
		return
	}

	slog.Info("Successfully get game info", "gameInfo", gameInfo)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gameInfo,
	})
}

func (h *PlayHandler) StartGame(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("userId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	slog.Info("StartGame", "gameId", gameId.(int), "userId", userId.(int))
	err := h.service.StartGame(userId.(int), gameId.(int))
	if err != nil {
		slog.Warn("Error StartGame", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot start game",
		})
		return
	}

	slog.Info("Successfully start game", "gameId", gameId.(int), "userId", userId.(int))
	c.Redirect(http.StatusTemporaryRedirect, "/play/playroom")
}

func (h *PlayHandler) CancelGame(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("userId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	slog.Info("CancelGame", "gameId", gameId.(int), "userId", userId.(int))
	err := h.service.CancelGame(userId.(int), gameId.(int))
	if err != nil {
		slog.Warn("Error CancelGame", "err", err, "gameId", gameId.(int), "userId", userId.(int))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot cancel game",
		})
	}

	slog.Info("Successfully cancel game", "gameId", gameId.(int), "userId", userId.(int))
	c.Redirect(http.StatusTemporaryRedirect, "/main")
}

func (h *PlayHandler) RemovePlayer(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		slog.Warn("Error ParseJsonRequest", "err", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	playerId, exists := jsonMap["playerId"].(float64)
	if !exists {
		slog.Warn("PlayerId not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "player id not found",
		})
	}
	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("userId not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "cannot remove player",
		})
	}

	gameId, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	slog.Info("RemovePlayer", "gameId", gameId.(int), "userId", userId.(int), "playerId", int(playerId))
	err = h.service.RemovePlayer(gameId.(int), userId.(int), int(playerId))
	if err != nil {
		slog.Warn("Error RemovePlayer", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})
		return
	}

	slog.Info("Successfully remove player", "gameId", gameId.(int), "userId", userId.(int), "playerId", int(playerId))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (h *PlayHandler) LeaveGame(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		slog.Warn("userId not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "cannot remove player",
		})
	}

	gameId, exists := c.Get("gameId")
	if !exists {
		slog.Warn("gameId not found in context")
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	slog.Info("RemovePlayer", "gameId", gameId.(int), "userId", userId.(int))
	err := h.service.RemovePlayer(gameId.(int), userId.(int), userId.(int))
	if err != nil {
		slog.Warn("Error RemovePlayer", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})
		return
	}

	slog.Info("Successfully leaving player", "gameId", gameId.(int), "userId", userId.(int))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
