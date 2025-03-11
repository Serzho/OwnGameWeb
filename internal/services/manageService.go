package services

import (
	"OwnGameWeb/internal/database"
	"errors"
)

type ManageService struct {
	dbController *database.DbController
}

func NewManageService(c *database.DbController) *ManageService {
	return &ManageService{dbController: c}
}

func (s *ManageService) JoinGame(_ string) error {
	return errors.New("not implemented")
}

func (s *ManageService) CreateGame(userId int, _ string, title string, maxPlayers int) (int, error) {
	_, err := s.dbController.GetCurrentGameByMasterId(userId)
	if err == nil {
		return 0, errors.New("player already playing")
	}

	// TODO: сделать генерацию кода приглашения
	err = s.dbController.AddGame(title, "000000", userId, maxPlayers)
	if err != nil {
		return 0, err
	}
	return -1, nil // TODO:  СДЕЛАТЬ ПОЛУЧЕНИЕ ID игры

}
