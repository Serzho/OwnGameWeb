package services

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthService struct {
	dbController *database.DbController
	Cfg          *config.Config
}

func NewAuthService(c *database.DbController, config *config.Config) *AuthService {
	return &AuthService{dbController: c, Cfg: config}
}

func (s *AuthService) SignIn(email, password string) (int, error) {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.dbController.GetUser(preparedEmail)

	if err != nil {
		return 0, errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, errors.New("incorrect email or password")
	}

	return user.Id, nil
}

func (s *AuthService) SignUp(name, email, password string) error {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	preparedPassword := strings.TrimSpace(password)

	_, err := s.dbController.GetUser(preparedEmail)
	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), 5)
	if err != nil {
		return errors.New("incorrect password")
	}

	err = s.dbController.AddUser(name, preparedEmail, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RecoverPassword(_ string) error {
	return errors.New("not implemented")
}
