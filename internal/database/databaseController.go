package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (d *DbController) GetUser(email string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), d.conn, &user, `
        SELECT id, name, email, password 
        FROM "user"
        WHERE email = $1
        LIMIT 1;
    `, email)

	return &user, err
}

func (d *DbController) GetPassword(email string) (string, error) {
	user, err := d.GetUser(email)
	if err != nil {
		return "", err
	}

	return user.Password, nil
}

func (d *DbController) AddUser(name, email, password string) error {
	_, err := d.GetUser(email)
	if err == nil {
		return errors.New("email already exists")
	}

	_, err = d.conn.Exec(
		context.Background(),
		`INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3);`,
		name, email, password,
	)

	return err
}

func (d *DbController) Close() error {
	return d.conn.Close(context.Background())
}
