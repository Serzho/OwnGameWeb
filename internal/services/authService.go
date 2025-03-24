package services

import (
	"log/slog"
	"strings"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(email, password string) (int, error)
	SignUp(name, email, password string) error
	RecoverPassword(email string) error
}

type AuthServiceImpl struct {
	dbController *database.DBController
	Cfg          *config.Config
}

func NewAuthService(c *database.DBController, config *config.Config) *AuthServiceImpl {
	return &AuthServiceImpl{dbController: c, Cfg: config}
}

func (s *AuthServiceImpl) SignIn(email, password string) (int, error) {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	slog.Info("Getting user by email", "email", preparedEmail)
	user, err := s.dbController.GetUserByEmail(preparedEmail)
	if err != nil {
		slog.Warn("User not found", "email", preparedEmail, "err", err)
		return 0, ErrInvalidEmail
	}

	slog.Info("Comparing password", "password", password, "user", user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		slog.Warn("Invalid password", "password", password, "user", user, "err", err)
		return 0, ErrIncorrectPassword
	}

	slog.Info("Successfully logged in", "email", email, "user", user)
	return user.ID, nil
}

func (s *AuthServiceImpl) SignUp(name, email, password string) error {
	preparedEmail := strings.ToLower(strings.TrimSpace(email))
	preparedPassword := strings.TrimSpace(password)

	slog.Info("Creating user", "email", preparedEmail, "user", preparedEmail)
	_, err := s.dbController.GetUserByEmail(preparedEmail)
	if err == nil {
		slog.Warn("User already exists", "email", preparedEmail)
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), 5)
	if err != nil {
		slog.Warn("Error generating password", "err", err)
		return ErrIncorrectPassword
	}

	slog.Info("Creating user", "email", preparedEmail, "name", name)
	err = s.dbController.AddUser(name, preparedEmail, string(hashedPassword))
	if err != nil {
		slog.Warn("Error creating user", "err", err, "email", email, "name", name)
		return ErrCreatingUser
	}

	slog.Info("Successfully created user", "email", email, "name", name)
	return nil
}

func (s *AuthServiceImpl) RecoverPassword(_ string) error {
	return ErrNotImplemented
}
