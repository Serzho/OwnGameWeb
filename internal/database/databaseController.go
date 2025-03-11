package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"time"
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

	if err != nil {
		return nil, err
	}

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
	_, err := d.conn.Exec(
		context.Background(),
		`INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3);`,
		name, email, password,
	)

	return err
}

func (d *DbController) AddGame(title string, inviteCode string, userId int, maxPlayers int) error {
	_, err := d.conn.Exec(
		context.Background(),
		`INSERT INTO "game" (title, status, invite_code, start_time, master_id, players_ids, max_players)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		title, "created", inviteCode, time.Now(), userId, "ARRAY[]", maxPlayers,
	)

	return err
}

func (d *DbController) GetCurrentGameByMasterId(masterId int) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.conn, &game, `
        SELECT * FROM "game"
        WHERE master_id = $1 AND status != 'archieved'
        LIMIT 1;
    `, masterId)

	if err != nil {
		return nil, err
	}

	return &game, err
}

func (d *DbController) Close() error {
	return d.conn.Close(context.Background())
}
