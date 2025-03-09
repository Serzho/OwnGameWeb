package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type DbController struct {
	conn *pgx.Conn
}

func NewDbController(cfg *config.Config) *DbController {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic(err)
	}

	return &DbController{conn: conn}
}

func (d *DbController) GetUser(_ string) (*models.User, error) {
	return &models.User{}, errors.New("not implemented")
}

func (d *DbController) GetPassword(_ string) (string, error) {
	return "", errors.New("not implemented")
}

func (d *DbController) AddUser(_, _, _ string) error {
	return errors.New("not implemented")
}

func (d *DbController) Close() error {
	return d.conn.Close(context.Background())
}
