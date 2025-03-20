package services

import (
	"encoding/json"
	"log/slog"

	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"

	"github.com/pkg/errors"
)

type PlayService struct {
	dbController *database.DBController
}

func NewPlayService(c *database.DBController) *PlayService {
	return &PlayService{dbController: c}
}

func (s *PlayService) RemovePlayer(gameID, masterID, playerID int) error {
	slog.Info("Search game", "gameID", gameID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err)
		return errors.New("game not found")
	}

	if game.MasterID != masterID && masterID != playerID {
		slog.Warn("MasterID is not match", "game", game, "masterID", masterID, "playerID", playerID)
		return errors.New("master id not match")
	}

	slog.Info("Remove player", "playerID", playerID, "gameID", gameID)
	err = s.dbController.RemovePlayer(gameID, playerID)
	if err != nil {
		slog.Warn("Error removing player", "playerID", playerID, "gameID", gameID, "error", err)
		return errors.New("cannot remove player")
	}

	slog.Info("Successfully remove player", "playerID", playerID, "gameID", gameID)
	return nil
}

func (s *PlayService) CancelGame(userID, gameID int) error {
	slog.Info("Search game", "gameID", gameID, "userID", userID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err)
		return errors.New("game not found")
	}

	if game.MasterID != userID {
		slog.Warn("MasterID is not match", "game", game, "userID", userID)
		return errors.New("master id not match")
	}

	slog.Info("Delete game", "gameID", gameID)
	err = s.dbController.DeleteGame(gameID)
	if err != nil {
		slog.Warn("Error deleting game", "id", gameID, "error", err)
		return errors.New("cannot delete game")
	}

	slog.Info("Successfully delete game", "gameID", gameID)
	return nil
}

func (s *PlayService) StartGame(userID, gameID int) error {
	slog.Info("Search game", "gameID", gameID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err)
		return errors.New("game not found")
	}

	if game.MasterID != userID {
		slog.Warn("MasterID is not match", "game", game, "userID", userID)
		return errors.New("master id not match")
	}

	slog.Info("Setting game status", "gameID", gameID, "status", "firststage")
	err = s.dbController.SetGameStatus(gameID, "firststage")
	if err != nil {
		slog.Warn("Error setting game status", "gameID", gameID, "error", err)
		return errors.New("cannot set game status")
	}

	slog.Info("Successfully set game status", "game", game, "userID", userID)
	return nil
}

func (s *PlayService) GetGameInfo(gameID, userID int) (string, error) {
	slog.Info("Search game", "gameID", gameID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err, "userID", userID)
		return "", errors.New("game not found")
	}

	players := make([]models.PlayerJSON, 0, len(game.PlayersIDs))
	for _, playerID := range game.PlayersIDs {
		slog.Info("Search player", "playerId", playerID)
		user, err := s.dbController.GetUser(playerID)
		if err != nil {
			slog.Warn("Error getting user", "id", playerID, "error", err)
			return "", errors.New("user not found")
		}
		players = append(players, models.PlayerJSON{ID: playerID, Name: user.Name})
	}

	isHost := userID == game.MasterID
	gameInfo := &models.GameInfoJSON{
		Title: game.Title, Players: players, MaxPlayers: game.MaxPlayers, IsHost: isHost,
		InviteCode: game.InviteCode,
	}
	slog.Info("Marshalling gameinfo", "userID", userID, "gameinfo", gameInfo)
	result, err := json.Marshal(gameInfo)
	if err != nil {
		slog.Warn("Error marshalling gameinfo", "error", err)
		return "", errors.New("json marshal error")
	}

	slog.Info("Successfully returning gameinfo", "userID", userID, "result", result)
	return string(result), nil
}

func (s *PlayService) CheckIsMaster(userID, gameID int) (bool, error) {
	slog.Info("Search game", "gameID", gameID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err)
	}

	return game.MasterID == userID, nil
}

func (s *PlayService) GetSampleContent(gameID, userID int) (string, error) {
	slog.Info("Search game", "gameID", gameID, "userID", userID)
	game, err := s.dbController.GetGame(gameID)
	if err != nil {
		slog.Warn("Error getting game", "id", gameID, "error", err)
		return "", errors.New("game not found")
	}

	if game.MasterID != userID {
		slog.Warn("MasterID is not match", "game", game, "userID", userID)
		return "", errors.New("master id not match")
	}

	slog.Info("Search sample", "game", game)
	sample, err := s.dbController.GetSample(game.Sample)
	if err != nil {
		slog.Warn("Error getting sample", "game", game, "error", err)
		return "", errors.New("sample not found")
	}

	return sample.Content, nil
}
