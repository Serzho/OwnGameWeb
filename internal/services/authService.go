package services

import (
	"OwnGameWeb/internal/database"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthService struct {
	dbController *database.DbController
}

func NewAuthService(c *database.DbController) *AuthService {
	return &AuthService{dbController: c}
}

func (s *AuthService) SignIn(email, password string) error {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	targetPassword, err := s.dbController.GetPassword(preparedEmail)

	if err != nil {
		return errors.New("invalid email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(targetPassword), []byte(password))
	if err != nil {
		return errors.New("incorrect email or password")
	}

	return nil
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
