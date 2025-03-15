package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"errors"
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

func (d *DbController) GetPack(packId int) (*models.QuestionPack, error) {
	var questionPack models.QuestionPack
	err := pgxscan.Get(context.Background(), d.conn, &questionPack, `
		SELECT id, title, filename, owner
		FROM "question_pack"
		WHERE id = $1
		LIMIT 1;
	`, packId)

	if err != nil {
		return nil, err
	}

	return &questionPack, err
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

func (d *DbController) AddPack(userId int, filename string) error {
	tx, err := d.conn.Begin(context.Background())
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	if err != nil {
		return errors.New("transaction start failed")
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO "question_pack" (title, filename, owner) VALUES ($1, $2, $3);`,
		"Новый пак", filename, userId,
	)

	if err != nil {
		return errors.New("database add pack failed")
	}

	_, err = d.conn.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_append(packs, currval(pg_get_serial_sequence('question_pack', 'id')))
			WHERE id=$1`,
		userId,
	)
	if err != nil {
		return errors.New("database add pack to user failed")
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return errors.New("transaction commit failed")
	}

	return nil
}

func (d *DbController) GetUserPacks(userId int) (*[]models.QuestionPack, error) {

	var packs []models.QuestionPack

	rows, err := d.conn.Query(context.Background(), `
        SELECT  p.id, title, filename, p.owner  FROM "user" u
			JOIN LATERAL unnest(packs) AS pack_id ON true
			JOIN question_pack p ON p.id = pack_id
			WHERE u.id = $1;
    `, userId)

	if err != nil {
		return nil, errors.New("query error")
	}

	err = pgxscan.ScanAll(&packs, rows)

	if err != nil {
		return nil, err
	}

	return &packs, err
}
