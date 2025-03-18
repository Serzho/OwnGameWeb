package database

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
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

func (d *DbController) Close() {
	d.pool.Close()
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
		slog.Warn("Error getting user by email", "error", err, "email", email)
		return nil, err
	}

	slog.Info("User found", "email", email, "user", user)
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
		slog.Warn("Error getting user", "error", err, "id", userId)
		return nil, err
	}

	slog.Info("User found", "id", userId, "user", user)
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
		slog.Warn("Error getting pack", "error", err, "id", packId)
		return nil, err
	}

	slog.Info("QuestionPack found", "id", packId, "questionPack", questionPack)
	return &questionPack, err
}

func (d *DbController) DeletePack(packId int) error {
	slog.Info("Deleting pack", "id", packId)
	tx, err := d.pool.Begin(context.Background())
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	if err != nil {
		slog.Warn("Error begin transaction", "error", err)
		return errors.New("transaction start failed")
	}

	_, err = tx.Exec(context.Background(),
		`DELETE FROM "question_pack" WHERE id = $1;`, packId,
	)

	if err != nil {
		slog.Warn("Error delete pack", "error", err)
		return errors.New("database delete pack failed")
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_remove(packs, $1::int)
			WHERE packs @> ARRAY[$1::int]`, packId,
	)
	if err != nil {
		slog.Warn("Error removing pack from user", "error", err)
		return errors.New("removing packs from user failed")
	}

	err = tx.Commit(context.Background())

	if err != nil {
		slog.Warn("Error commit transaction", "error", err, "id", packId)
		return errors.New("transaction commit failed")
	}

	slog.Info("Pack deleted", "id", packId)
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

	if err != nil {
		slog.Warn("Error adding user", "error", err, "email", email, "name", name, "password", password)
		return errors.New("database add user failed")
	}

	slog.Info("Added user", "email", email, "name", name, "password", password)
	return nil
}

func (d *DbController) AddGame(title string, inviteCode string, userId int, maxPlayers int, sampleId int) (int, error) {
	var id int
	row := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "game" (title, status, invite_code, start_time, master_id, players_ids, max_players, sample)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		title, "created", inviteCode, time.Now(), userId, "{}", maxPlayers, sampleId,
	).Scan(&id)

	if row != nil {
		slog.Warn("Error adding game", "error", row, "title", title, "inviteCode", inviteCode, "UserId", userId, "maxPlayers", maxPlayers, "sample", sampleId)
		return 0, errors.New("cannot insert game")
	}

	slog.Info("Added game", "id", id)
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
		slog.Warn("Error getting game by masterId", "error", err, "masterId", masterId)
		return nil, err
	}

	slog.Info("Game found", "masterId", masterId, "game", game)
	return &game, err
}

func (d *DbController) AddPack(userId int, filename string) error {
	tx, err := d.pool.Begin(context.Background())
	defer func() {
		_ = tx.Rollback(context.Background())
	}()

	if err != nil {
		slog.Warn("Error begin transaction", "error", err, "id", userId, "filename", filename)
		return errors.New("transaction start failed")
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO "question_pack" (title, filename, owner) VALUES ($1, $2, $3);`,
		"Новый пак", filename, userId,
	)

	if err != nil {
		slog.Warn("Error adding pack", "error", err, "id", userId, "filename", filename)
		return errors.New("database add pack failed")
	}

	_, err = tx.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_append(packs, currval(pg_get_serial_sequence('question_pack', 'id')))
			WHERE id=$1`,
		userId,
	)
	if err != nil {
		slog.Warn("Error adding pack to user", "error", err, "id", userId, "filename", filename)
		return errors.New("database add pack to user failed")
	}

	err = tx.Commit(context.Background())

	if err != nil {
		slog.Warn("Error commit transaction", "error", err, "id", userId, "filename", filename)
		return errors.New("transaction commit failed")
	}

	slog.Info("Added pack", "id", userId, "filename", filename)
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
		slog.Warn("Error getting user packs", "error", err, "id", userId)
		return nil, errors.New("query error")
	}

	err = pgxscan.ScanAll(&packs, rows)

	if err != nil {
		slog.Warn("Error getting user packs", "error", err, "id", userId)
		return nil, errors.New("error getting user packs")
	}

	slog.Info("Get user packs", "id", userId, "packs", packs)
	return &packs, err
}

func (d *DbController) AddSample(sample *models.QuestionSample) (int, error) {
	var id int
	row := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "question_sample" (pack, content) VALUES ($1::int, $2) RETURNING id;`,
		sample.Pack, sample.Content,
	).Scan(&id)

	if row != nil {
		slog.Warn("Error adding sample", "error", row, "sample", sample)
		return 0, errors.New("insert sample failed")
	}

	slog.Info("Added sample", "id", id)
	return id, nil
}

func (d *DbController) GetInvites() ([]string, error) {
	var invites []string
	rows, err := d.pool.Query(
		context.Background(),
		`SELECT invite_code FROM "game" WHERE status = 'created';`,
	)
	if err != nil {
		slog.Warn("Error getting invites", "error", err)
		return nil, errors.New("get invites failed")
	}

	err = pgxscan.ScanAll(&invites, rows)
	if err != nil {
		slog.Warn("Error getting invites", "error", err)
		return nil, errors.New("scan invites failed")
	}

	slog.Info("Get invites", "invites", invites)
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
		slog.Warn("Error getting game", "error", err, "id", gameId)
		return nil, err
	}

	slog.Info("Game found", "id", gameId)
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
		slog.Warn("Error getting game by invite code", "error", err, "code", code)
		return nil, err
	}

	slog.Info("Get game by invite code", "game", game, "code", code)
	return &game, err
}

func (d *DbController) JoinGame(userId, gameId int) error {
	_, err := d.pool.Exec(
		context.Background(),
		`UPDATE game SET players_ids = array_append(players_ids, $1::int)
 			WHERE id = $2::int;`,
		userId, gameId,
	)
	if err != nil {
		slog.Warn("Error joining game", "error", err, "id", gameId, "userId", userId)
		return errors.New("join game failed")
	}

	slog.Info("Joined game", "id", gameId, "userId", userId)
	return nil

}

func (d *DbController) SetGameStatus(gameId int, status string) error {
	_, err := d.pool.Exec(
		context.Background(),
		`UPDATE game SET status = $1 
 			WHERE id = $2::int;`,
		status, gameId,
	)

	if err != nil {
		slog.Warn("Error setting game status", "error", err, "id", gameId, "status", status)
		return errors.New("update game status failed")
	}

	slog.Info("Set game status", "id", gameId, "status", status)
	return nil
}

func (d *DbController) DeleteGame(gameId int) error {
	_, err := d.pool.Exec(context.Background(),
		`DELETE FROM "game" WHERE id = $1;`, gameId,
	)

	if err != nil {
		slog.Warn("Error deleting game", "error", err, "id", gameId)
		return errors.New("delete game failed")
	}

	slog.Info("Deleted game", "id", gameId)
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
		slog.Warn("Error removing player", "error", err, "id", gameId, "userId", userId)
		return errors.New("update game players failed")
	}

	slog.Info("Removed player", "id", gameId, "userId", userId)
	return nil
}

func (d *DbController) GetSample(sampleId int) (*models.QuestionSample, error) {
	var sample models.QuestionSample

	err := pgxscan.Get(context.Background(), d.pool, &sample, `
        SELECT * FROM "question_sample"
        WHERE id = $1
        LIMIT 1;
    `, sampleId)

	if err != nil {
		slog.Warn("Error getting sample", "error", err, "id", sampleId)
		return nil, err
	}

	return &sample, nil
}
