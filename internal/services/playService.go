package services

import (
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"
	"encoding/json"
	"github.com/pkg/errors"
)

type PlayService struct {
	dbController *database.DbController
}

func NewPlayService(c *database.DbController) *PlayService {
	return &PlayService{dbController: c}
}

func (s *PlayService) RemovePlayer(gameId, masterId, playerId int) error {
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		return errors.New("game not found")
	}

	if game.MasterId != masterId && masterId != playerId {
		return errors.New("master id not match")
	}

	err = s.dbController.RemovePlayer(gameId, playerId)
	if err != nil {
		return errors.New("cannot remove player")
	}
	return nil
}

func (s *PlayService) CancelGame(userId, gameId int) error {
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		return errors.New("game not found")
	}

	if game.MasterId != userId {
		return errors.New("master id not match")
	}

	err = s.dbController.DeleteGame(gameId)
	if err != nil {
		return errors.New("cannot delete game")
	}

	return nil
}

func (s *PlayService) StartGame(userId, gameId int) error {
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		return errors.New("game not found")
	}

	if game.MasterId != userId {
		return errors.New("master id not match")
	}

	err = s.dbController.SetGameStatus(gameId, "firststage")
	if err != nil {
		return errors.New("cannot set game status")
	}

	return nil
}

func (s *PlayService) GetGameInfo(gameId, userId int) (string, error) {
	game, err := s.dbController.GetGame(gameId)
	if err != nil {
		return "", errors.New("game not found")
	}

	players := make([]models.PlayerJson, 0, len(game.PlayersIds))
	for _, id := range game.PlayersIds {
		user, err := s.dbController.GetUser(id)
		if err != nil {
			return "", errors.New("user not found")
		}
		players = append(players, models.PlayerJson{Id: id, Name: user.Name})
	}

	isHost := userId == game.MasterId
	gameInfo := &models.GameInfoJson{Title: game.Title, Players: players, MaxPlayers: game.MaxPlayers, IsHost: isHost, InviteCode: game.InviteCode}
	result, err := json.Marshal(gameInfo)
	if err != nil {
		return "", errors.New("json marshal error")
	}
	return string(result), nil
}
