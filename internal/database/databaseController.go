package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DbController struct {
	pool *pgxpool.Pool
}

func NewDbController(cfg *config.Config) *DbController {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic(err)
	}

	return &DbController{pool: pool}
}

func (d *DbController) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), d.pool, &user, `
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

func (d *DbController) GetUser(userId int) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), d.pool, &user, `
        SELECT id, name, email, password 
        FROM "user"
        WHERE id = $1
        LIMIT 1;
    `, userId)

	if err != nil {
		return nil, err
	}

	return &user, err
}
func (d *DbController) GetPack(packId int) (*models.QuestionPack, error) {
	var questionPack models.QuestionPack
	err := pgxscan.Get(context.Background(), d.pool, &questionPack, `
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

func (d *DbController) DeletePack(packId int) error {
	tx, err := d.pool.Begin(context.Background())
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	if err != nil {
		return errors.New("transaction start failed")
	}

	_, err = tx.Exec(context.Background(),
		`DELETE FROM "question_pack" WHERE id = $1;`, packId,
	)

	if err != nil {
		return errors.New("database delete pack failed")
	}

	_, err = d.pool.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_remove(packs, $1::int)
			WHERE packs @> ARRAY[$1::int]`, packId,
	)
	if err != nil {
		return errors.New("removing packs from user failed")
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return errors.New("transaction commit failed")
	}

	return nil
}

func (d *DbController) GetPassword(email string) (string, error) {
	user, err := d.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	return user.Password, nil
}

func (d *DbController) AddUser(name, email, password string) error {
	_, err := d.pool.Exec(
		context.Background(),
		`INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3);`,
		name, email, password,
	)

	return err
}

func (d *DbController) AddGame(title string, inviteCode string, userId int, maxPlayers int, sampleId int) (int, error) {

	var id int
	err := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "game" (title, status, invite_code, start_time, master_id, players_ids, max_players, sample)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		title, "created", inviteCode, time.Now(), userId, "{}", maxPlayers, sampleId,
	).Scan(&id)

	if err != nil {
		return 0, errors.New("cannot insert game")
	}

	return id, nil
}

func (d *DbController) GetCurrentGameByMasterId(masterId int) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE master_id = $1 AND status != 'archieved'
        LIMIT 1;
    `, masterId)

	if err != nil {
		return nil, err
	}

	return &game, err
}

func (d *DbController) Close() {
	d.pool.Close()
}

func (d *DbController) AddPack(userId int, filename string) error {
	tx, err := d.pool.Begin(context.Background())
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

	_, err = d.pool.Exec(
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

	rows, err := d.pool.Query(context.Background(), `
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

func (d *DbController) AddSample(sample *models.QuestionSample) (int, error) {
	var id int
	err := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "question_sample" (pack, content) VALUES ($1::int, $2) RETURNING id;`,
		sample.Pack, sample.Content,
	).Scan(&id)

	if err != nil {
		return 0, errors.New("insert sample failed")
	}

	return id, nil
}

func (d *DbController) GetInvites() ([]string, error) {
	var invites []string
	rows, err := d.pool.Query(
		context.Background(),
		`SELECT invite_code FROM "game" WHERE status = 'created';`,
	)
	if err != nil {
		return nil, errors.New("get invites failed")
	}

	err = pgxscan.ScanAll(&invites, rows)
	if err != nil {
		return nil, errors.New("scan invites failed")
	}

	return invites, nil
}

func (d *DbController) GetGame(gameId int) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE id = $1
        LIMIT 1;
    `, gameId)

	if err != nil {
		return nil, err
	}

	return &game, err
}

func (d *DbController) GetGameByInviteCode(code string) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE invite_code = $1 and status = 'created'
        LIMIT 1;
    `, code)

	if err != nil {
		return nil, err
	}

	return &game, err
}

func (d *DbController) JoinGame(userId, gameId int) error {
	err := d.pool.QueryRow(
		context.Background(),
		`UPDATE game SET players_ids = array_append(players_ids, $1::int)
 			WHERE id = $2::int;`,
		userId, gameId,
	)
	if err != nil {
		return errors.New("join game failed")
	}
	return nil

}

func (d *DbController) SetGameStatus(gameId int, status string) error {
	err := d.pool.QueryRow(
		context.Background(),
		`UPDATE game SET status = $1 
 			WHERE id = $2::int;`,
		status, gameId,
	)

	if err != nil {
		return errors.New("update game status failed")
	}
	return nil
}

func (d *DbController) DeleteGame(gameId int) error {
	_, err := d.pool.Exec(context.Background(),
		`DELETE FROM "game" WHERE id = $1;`, gameId,
	)

	if err != nil {
		return errors.New("delete game failed")
	}

	return nil
}

func (d *DbController) RemovePlayer(gameId, userId int) error {
	err := d.pool.QueryRow(
		context.Background(),
		`UPDATE game SET players_ids = array_remove(players_ids, $1::int)
 			WHERE id = $2::int;`,
		userId, gameId,
	)

	if err != nil {
		return errors.New("update game players failed")
	}
	return nil
}
