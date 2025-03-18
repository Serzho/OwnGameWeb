package services

import (
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"
	"encoding/json"
	"github.com/pkg/errors"
	"log/slog"
)

type PlayService struct {
	dbController *database.DbController
}

func NewPlayService(c *database.DbController) *PlayService {
	return &PlayService{dbController: c}
}

func (s *PlayService) RemovePlayer(gameId, masterId, playerId int) error {
	slog.Info("Search game", "gameId", gameId)
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		slog.Warn("Error getting game", "id", gameId, "error", err)
		return errors.New("game not found")
	}

	if game.MasterId != masterId && masterId != playerId {
		slog.Warn("MasterId is not match", "game", game, "masterId", masterId, "playerId", playerId)
		return errors.New("master id not match")
	}

	slog.Info("Remove player", "playerId", playerId, "gameId", gameId)
	err = s.dbController.RemovePlayer(gameId, playerId)
	if err != nil {
		slog.Warn("Error removing player", "playerId", playerId, "gameId", gameId, "error", err)
		return errors.New("cannot remove player")
	}

	slog.Info("Successfully remove player", "playerId", playerId, "gameId", gameId)
	return nil
}

func (s *PlayService) CancelGame(userId, gameId int) error {
	slog.Info("Search game", "gameId", gameId, "userId", userId)
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		slog.Warn("Error getting game", "id", gameId, "error", err)
		return errors.New("game not found")
	}

	if game.MasterId != userId {
		slog.Warn("MasterId is not match", "game", game, "userId", userId)
		return errors.New("master id not match")
	}

	slog.Info("Delete game", "gameId", gameId)
	err = s.dbController.DeleteGame(gameId)
	if err != nil {
		slog.Warn("Error deleting game", "id", gameId, "error", err)
		return errors.New("cannot delete game")
	}

	slog.Info("Successfully delete game", "gameId", gameId)
	return nil
}

func (s *PlayService) StartGame(userId, gameId int) error {
	slog.Info("Search game", "gameId", gameId)
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		slog.Warn("Error getting game", "id", gameId, "error", err)
		return errors.New("game not found")
	}

	if game.MasterId != userId {
		slog.Warn("MasterId is not match", "game", game, "userId", userId)
		return errors.New("master id not match")
	}

	slog.Info("Setting game status", "gameId", gameId, "status", "firststage")
	err = s.dbController.SetGameStatus(gameId, "firststage")
	if err != nil {
		slog.Warn("Error setting game status", "gameId", gameId, "error", err)
		return errors.New("cannot set game status")
	}

	slog.Info("Successfully set game status", "game", game, "userId", userId)
	return nil
}

func (s *PlayService) GetGameInfo(gameId, userId int) (string, error) {
	slog.Info("Search game", "gameId", gameId)
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		slog.Warn("Error getting game", "id", gameId, "error", err, "userId", userId)
		return "", errors.New("game not found")
	}

	players := make([]models.PlayerJson, 0, len(game.PlayersIds))
	for _, id := range game.PlayersIds {
		slog.Info("Search player", "playerId", id)
		user, err := s.dbController.GetUser(id)
		if err != nil {
			slog.Warn("Error getting user", "id", id, "error", err)
			return "", errors.New("user not found")
		}
		players = append(players, models.PlayerJson{Id: id, Name: user.Name})
	}

	isHost := userId == game.MasterId
	gameInfo := &models.GameInfoJson{Title: game.Title, Players: players, MaxPlayers: game.MaxPlayers, IsHost: isHost, InviteCode: game.InviteCode}
	slog.Info("Marshalling gameinfo", "userId", userId, "gameinfo", gameInfo)
	result, err := json.Marshal(gameInfo)
	if err != nil {
		slog.Warn("Error marshalling gameinfo", "error", err)
		return "", errors.New("json marshal error")
	}

	slog.Info("Successfully returning gameinfo", "userId", userId, "result", result)
	return string(result), nil
}
