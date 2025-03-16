package handlers

import (
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
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
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	c.HTML(http.StatusOK, "waitingroom.html", gin.H{})
}

func (h *PlayHandler) GameInfo(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	gameInfo, err := h.service.GetGameInfo(gameId.(int), userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot get game info",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gameInfo,
	})
}

func (h *PlayHandler) StartGame(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	err := h.service.StartGame(userId.(int), gameId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot start game",
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/play/playroom")
}

func (h *PlayHandler) CancelGame(c *gin.Context) {
	gameId, exists := c.Get("gameId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/signin")
		return
	}

	err := h.service.CancelGame(userId.(int), gameId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot cancel game",
		})
	}

	c.Redirect(http.StatusTemporaryRedirect, "/main")
}

func (h *PlayHandler) RemovePlayer(c *gin.Context) {
	jsonMap, err := utils.ParseJsonRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}
	playerId, exists := jsonMap["playerId"].(float64)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "player id not found",
		})
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "cannot remove player",
		})
	}

	gameId, exists := c.Get("gameId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	err = h.service.RemovePlayer(gameId.(int), userId.(int), int(playerId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (h *PlayHandler) LeaveGame(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "cannot remove player",
		})
	}

	gameId, exists := c.Get("gameId")
	if !exists {
		c.Redirect(http.StatusTemporaryRedirect, "/main")
		return
	}

	err := h.service.RemovePlayer(gameId.(int), userId.(int), userId.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
