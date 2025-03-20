package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/services"

	"github.com/gin-gonic/gin"
)

type PlayHandler struct {
	service *services.PlayService
}

func getIntFromContext(c *gin.Context, key string) (int, error) {
	value, exists := c.Get(key)
	if !exists {
		return -1, ErrNotFoundInContext
	}

	intValue, ok := value.(int)
	if !ok {
		return -1, ErrIncorrectType
	}

	return intValue, nil
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

func (h *PlayHandler) MasterRoomPage(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("Check is master", "userID", userID, "gameID", gameID)

	isMaster, err := h.service.CheckIsMaster(userID, gameID)
	if err != nil {
		slog.Warn("Check is master failed", "err", err, "userID", userID, "gameID", gameID)
	}

	if !isMaster {
		slog.Warn("userID is not master", "userID", userID, "gameID", gameID)
		c.Redirect(http.StatusTemporaryRedirect, "/main")

		return
	}

	c.HTML(http.StatusOK, "masterroom.html", gin.H{})
}

func (h *PlayHandler) GameInfo(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("GetGameInfo", "gameID", gameID, "userID", userID)

	gameInfo, err := h.service.GetGameInfo(gameID, userID)
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
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("StartGame", "gameID", gameID, "userID", userID)

	err = h.service.StartGame(userID, gameID)
	if err != nil {
		slog.Warn("Error StartGame", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot start game",
		})

		return
	}

	slog.Info("Successfully start game", "gameID", gameID, "userID", userID)
	c.JSON(http.StatusOK, gin.H{})
}

func (h *PlayHandler) CancelGame(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("CancelGame", "gameID", gameID, "userID", userID)

	err = h.service.CancelGame(userID, gameID)
	if err != nil {
		slog.Warn("Error CancelGame", "err", err, "gameID", gameID, "userID", userID)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot cancel game",
		})
	}

	slog.Info("Successfully cancel game", "gameID", gameID, "userID", userID)
	c.Redirect(http.StatusTemporaryRedirect, "/main")
}

func (h *PlayHandler) RemovePlayer(c *gin.Context) {
	jsonMap, err := utils.ParseJSONRequest(c)
	if err != nil {
		slog.Warn("Error ParseJSONRequest", "err", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fmt.Sprintf("%v", err),
		})

		return
	}

	playerID, exists := jsonMap["playerID"].(float64)
	if !exists {
		slog.Warn("PlayerId not found in context")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "player id not found",
		})
	}

	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("RemovePlayer", "gameID", gameID, "userID", userID, "playerID", int(playerID))

	err = h.service.RemovePlayer(gameID, userID, int(playerID))
	if err != nil {
		slog.Warn("Error RemovePlayer", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})

		return
	}

	slog.Info("Successfully remove player", "gameID", gameID, "userID", userID, "playerID", int(playerID))
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (h *PlayHandler) LeaveGame(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("RemovePlayer", "gameID", gameID, "userID", userID)

	err = h.service.RemovePlayer(gameID, userID, userID)
	if err != nil {
		slog.Warn("Error RemovePlayer", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "cannot remove player",
		})

		return
	}

	slog.Info("Successfully leaving player", "gameID", gameID, "userID", userID)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (h *PlayHandler) GetQuestions(c *gin.Context) {
	userID, err := getIntFromContext(c, "userID")
	if err != nil {
		slog.Warn("Error get userID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	gameID, err := getIntFromContext(c, "gameID")
	if err != nil {
		slog.Warn("Error get gameID", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("Getting questions", "gameID", gameID, "userID", userID)

	sampleContent, err := h.service.GetSampleContent(gameID, userID)
	if err != nil {
		slog.Warn("Error GetSampleContent", "err", err, "gameID", gameID, "userID", userID)
		c.JSON(http.StatusBadRequest, gin.H{})

		return
	}

	slog.Info("Successfully getting questions", "sampleContent", sampleContent,
		"gameID", gameID, "userID", userID,
	)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": sampleContent,
	})
}
