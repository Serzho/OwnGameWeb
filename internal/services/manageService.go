package services

import (
	"OwnGameWeb/internal/database"
)

type ManageService struct {
	dbController *database.DbController
}

func NewManageService(c *database.DbController) *ManageService {
	return &ManageService{dbController: c}
}

func (s *ManageService) JoinGame(_ string) error {
	return nil
}
