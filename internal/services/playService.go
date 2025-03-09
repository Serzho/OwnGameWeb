package services

import "OwnGameWeb/internal/database"

type PlayService struct {
	dbController *database.DbController
}

func NewPlayService(c *database.DbController) *PlayService {
	return &PlayService{dbController: c}
}
