package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBController struct {
	pool *pgxpool.Pool
}

func NewDBController(cfg *config.Config) *DBController {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic(err)
	}

	return &DBController{pool: pool}
}

func (d *DBController) Close() {
	d.pool.Close()
}

func (d *DBController) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), d.pool, &user, `
        SELECT id, name, email, password 
        FROM "user"
        WHERE email = $1
        LIMIT 1;
    `, email)
	if err != nil {
		slog.Warn("Error getting user by email", "error", err, "email", email)
		return nil, ErrGetUserByEmail
	}

	slog.Info("User found", "email", email, "user", user)
	return &user, nil
}

func (d *DBController) GetUser(userID int) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), d.pool, &user, `
        SELECT id, name, email, password 
        FROM "user"
        WHERE id = $1
        LIMIT 1;
    `, userID)
	if err != nil {
		slog.Warn("Error getting user", "error", err, "id", userID)
		return nil, ErrGetUser
	}

	slog.Info("User found", "id", userID, "user", user)
	return &user, nil
}

func (d *DBController) GetPack(packID int) (*models.QuestionPack, error) {
	var questionPack models.QuestionPack
	err := pgxscan.Get(context.Background(), d.pool, &questionPack, `
		SELECT id, title, filename, owner
		FROM "question_pack"
		WHERE id = $1
		LIMIT 1;
	`, packID)
	if err != nil {
		slog.Warn("Error getting pack", "error", err, "id", packID)
		return nil, ErrGetPack
	}

	slog.Info("QuestionPack found", "id", packID, "questionPack", questionPack)
	return &questionPack, nil
}

func (d *DBController) DeletePack(packID int) error {
	slog.Info("Deleting pack", "id", packID)
	transaction, err := d.pool.Begin(context.Background())
	defer func() {
		_ = transaction.Rollback(context.Background())
	}()

	if err != nil {
		slog.Warn("Error begin transaction", "error", err)
		return errors.New("transaction start failed")
	}

	_, err = transaction.Exec(context.Background(),
		`DELETE FROM "question_pack" WHERE id = $1;`, packID,
	)
	if err != nil {
		slog.Warn("Error delete pack", "error", err)
		return errors.New("database delete pack failed")
	}

	_, err = transaction.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_remove(packs, $1::int)
			WHERE packs @> ARRAY[$1::int]`, packID,
	)
	if err != nil {
		slog.Warn("Error removing pack from user", "error", err)
		return errors.New("removing packs from user failed")
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		slog.Warn("Error commit transaction", "error", err, "id", packID)
		return ErrTransactionCommit
	}

	slog.Info("Pack deleted", "id", packID)
	return nil
}

func (d *DBController) GetPassword(email string) (string, error) {
	user, err := d.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	return user.Password, nil
}

func (d *DBController) AddUser(name, email, password string) error {
	_, err := d.pool.Exec(
		context.Background(),
		`INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3);`,
		name, email, password,
	)
	if err != nil {
		slog.Warn("Error adding user", "error", err, "email", email, "name", name, "password", password)
		return ErrInsertUser
	}

	slog.Info("Added user", "email", email, "name", name, "password", password)
	return nil
}

func (d *DBController) UpdateUser(user *models.User) error {
	_, err := d.pool.Exec(
		context.Background(),
		`UPDATE "user" SET password = $1, name = $2 WHERE id = $3;`,
		user.Password, user.Name, user.ID,
	)
	if err != nil {
		slog.Warn("Error updating user", "error", err, "user", user)
		return ErrInsertUser
	}

	slog.Info("Updated user", "user", user)
	return nil
}

func (d *DBController) AddGame(title string, inviteCode string, userID int, maxPlayers int, sampleID int) (int, error) {
	var gameID int
	row := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "game" 
    	(title, status, invite_code, start_time, master_id, players_ids, max_players, sample, used_questions)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`,
		title, "created", inviteCode, time.Now(), userID, "{}", maxPlayers, sampleID, "{}",
	).Scan(&gameID)

	if row != nil {
		slog.Warn("Error adding game", "error", row, "title", title, "inviteCode", inviteCode,
			"UserId", userID, "maxPlayers", maxPlayers, "sample", sampleID)
		return 0, ErrInsertGame
	}

	slog.Info("Added game", "gameID", gameID)
	return gameID, nil
}

func (d *DBController) GetCurrentGameByMasterID(masterID int) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE master_id = $1 AND status != 'archieved'
        LIMIT 1;
    `, masterID)
	if err != nil {
		slog.Warn("Error getting game by masterID", "error", err, "masterID", masterID)
		return nil, ErrGetGameByMasterID
	}

	slog.Info("Game found", "masterID", masterID, "game", game)
	return &game, nil
}

func (d *DBController) AddPack(userID int, filename string) error {
	transaction, err := d.pool.Begin(context.Background())
	defer func() {
		_ = transaction.Rollback(context.Background())
	}()

	if err != nil {
		slog.Warn("Error begin transaction", "error", err, "id", userID, "filename", filename)
		return errors.New("transaction start failed")
	}

	_, err = transaction.Exec(context.Background(),
		`INSERT INTO "question_pack" (title, filename, owner) VALUES ($1, $2, $3);`,
		"Новый пак", filename, userID,
	)
	if err != nil {
		slog.Warn("Error adding pack", "error", err, "id", userID, "filename", filename)
		return errors.New("database add pack failed")
	}

	_, err = transaction.Exec(
		context.Background(),
		`UPDATE "user" SET packs = array_append(packs, currval(pg_get_serial_sequence('question_pack', 'id')))
			WHERE id=$1`,
		userID,
	)
	if err != nil {
		slog.Warn("Error adding pack to user", "error", err, "id", userID, "filename", filename)
		return errors.New("database add pack to user failed")
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		slog.Warn("Error commit transaction", "error", err, "id", userID, "filename", filename)
		return errors.New("transaction commit failed")
	}

	slog.Info("Added pack", "id", userID, "filename", filename)
	return nil
}

func (d *DBController) GetUserPacks(userID int) (*[]models.QuestionPack, error) {
	var packs []models.QuestionPack

	rows, err := d.pool.Query(context.Background(), `
        SELECT  p.id, title, filename, p.owner  FROM "user" u
			JOIN LATERAL unnest(packs) AS pack_id ON true
			JOIN question_pack p ON p.id = pack_id
			WHERE u.id = $1;
    `, userID)
	if err != nil {
		slog.Warn("Error getting user packs", "error", err, "id", userID)
		return nil, errors.New("query error")
	}

	err = pgxscan.ScanAll(&packs, rows)
	if err != nil {
		slog.Warn("Error getting user packs", "error", err, "id", userID)
		return nil, ErrGetUserPacks
	}

	slog.Info("Get user packs", "id", userID, "packs", packs)
	return &packs, nil
}

func (d *DBController) AddSample(sample *models.QuestionSample) (int, error) {
	var sampleID int
	row := d.pool.QueryRow(
		context.Background(),
		`INSERT INTO "question_sample" (pack, content) VALUES ($1::int, $2) RETURNING id;`,
		sample.Pack, sample.Content,
	).Scan(&sampleID)

	if row != nil {
		slog.Warn("Error adding sample", "error", row, "sample", sample)
		return 0, ErrInsertSample
	}

	slog.Info("Added sample", "sampleID", sampleID)
	return sampleID, nil
}

func (d *DBController) GetInvites() ([]string, error) {
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

func (d *DBController) GetGame(gameID int) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE id = $1
        LIMIT 1;
    `, gameID)
	if err != nil {
		slog.Warn("Error getting game", "error", err, "id", gameID)
		return nil, ErrGetGame
	}

	slog.Info("Game found", "id", gameID)
	return &game, nil
}

func (d *DBController) GetGameByInviteCode(code string) (*models.Game, error) {
	var game models.Game
	err := pgxscan.Get(context.Background(), d.pool, &game, `
        SELECT * FROM "game"
        WHERE invite_code = $1 and status = 'created'
        LIMIT 1;
    `, code)
	if err != nil {
		slog.Warn("Error getting game by invite code", "error", err, "code", code)
		return nil, ErrGetGameByInviteCode
	}

	slog.Info("Get game by invite code", "game", game, "code", code)
	return &game, nil
}

func (d *DBController) JoinGame(userID, gameID int) error {
	_, err := d.pool.Exec(
		context.Background(),
		`UPDATE game SET players_ids = array_append(players_ids, $1::int)
 			WHERE id = $2::int;`,
		userID, gameID,
	)
	if err != nil {
		slog.Warn("Error joining game", "error", err, "id", gameID, "userID", userID)
		return errors.New("join game failed")
	}

	slog.Info("Joined game", "id", gameID, "userID", userID)
	return nil
}

func (d *DBController) SetGameStatus(gameID int, status string) error {
	_, err := d.pool.Exec(
		context.Background(),
		`UPDATE game SET status = $1 
 			WHERE id = $2::int;`,
		status, gameID,
	)
	if err != nil {
		slog.Warn("Error setting game status", "error", err, "id", gameID, "status", status)
		return ErrSetGameStatus
	}

	slog.Info("Set game status", "id", gameID, "status", status)
	return nil
}

func (d *DBController) DeleteGame(gameID int) error {
	_, err := d.pool.Exec(context.Background(),
		`DELETE FROM "game" WHERE id = $1;`, gameID,
	)
	if err != nil {
		slog.Warn("Error deleting game", "error", err, "id", gameID)
		return ErrDeleteGame
	}

	slog.Info("Deleted game", "id", gameID)
	return nil
}

func (d *DBController) RemovePlayer(gameID, userID int) error {
	err := d.pool.QueryRow(
		context.Background(),
		`UPDATE game SET players_ids = array_remove(players_ids, $1::int)
 			WHERE id = $2::int;`,
		userID, gameID,
	)
	if err != nil {
		slog.Warn("Error removing player", "error", err, "id", gameID, "userID", userID)
		return errors.New("update game players failed")
	}

	slog.Info("Removed player", "id", gameID, "userID", userID)
	return nil
}

func (d *DBController) GetSample(sampleID int) (*models.QuestionSample, error) {
	var sample models.QuestionSample

	err := pgxscan.Get(context.Background(), d.pool, &sample, `
        SELECT * FROM "question_sample"
        WHERE id = $1
        LIMIT 1;
    `, sampleID)
	if err != nil {
		slog.Warn("Error getting sample", "error", err, "id", sampleID)

		return nil, ErrGetSample
	}

	return &sample, nil
}
