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

func (s *ManageService) CreateGame(_ int, _ string, _ string, _ int) (int, error) {
	return -1, errors.New("not implemented")
}
