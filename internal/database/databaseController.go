package database

import (
	"OwnGameWeb/internal/database/models"
	"errors"
)

type DbController struct {
}

func NewDbController() *DbController { return &DbController{} }

func (d *DbController) GetUser(_ string) (*models.User, error) {
	return &models.User{}, errors.New("not implemented")
}

func (d *DbController) GetPassword(_ string) (string, error) {
	return "", errors.New("not implemented")
}

func (d *DbController) AddUser(_, _, _ string) error {
	return errors.New("not implemented")
}
