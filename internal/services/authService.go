package services

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
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
	slog.Info("Getting user by email", "email", preparedEmail)
	user, err := s.dbController.GetUserByEmail(preparedEmail)

	if err != nil {
		slog.Warn("User not found", "email", preparedEmail, "err", err)
		return 0, errors.New("invalid email")
	}

	slog.Info("Comparing password", "password", password, "user", user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		slog.Warn("Invalid password", "password", password, "user", user, "err", err)
		return 0, errors.New("incorrect email or password")
	}

	slog.Info("Successfully logged in", "email", email, "user", user)
	return user.Id, nil
}

func (s *AuthService) SignUp(name, email, password string) error {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	preparedPassword := strings.TrimSpace(password)

	slog.Info("Creating user", "email", preparedEmail, "user", preparedEmail)
	_, err := s.dbController.GetUserByEmail(preparedEmail)
	if err == nil {
		slog.Warn("User already exists", "email", preparedEmail)
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), 5)
	if err != nil {
		slog.Warn("Error generating password", "err", err)
		return errors.New("incorrect password")
	}

	slog.Info("Creating user", "email", preparedEmail, "name", name)
	err = s.dbController.AddUser(name, preparedEmail, string(hashedPassword))
	if err != nil {
		slog.Warn("Error creating user", "err", err, "email", email, "name", name)
		return err
	}

	slog.Info("Successfully created user", "email", email, "name", name)
	return nil
}

func (s *AuthService) RecoverPassword(_ string) error {
	return errors.New("not implemented")
}
