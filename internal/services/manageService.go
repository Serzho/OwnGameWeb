package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"
)

type ManageService struct {
	dbController *database.DBController
	Cfg          *config.Config
}

func NewManageService(c *database.DBController, cfg *config.Config) *ManageService {
	return &ManageService{dbController: c, Cfg: cfg}
}

func (s *ManageService) AddServerPack(userID, packID int) error {
	err := s.dbController.AddServerPack(userID, packID)
	if err != nil {
		return ErrAddServerPack
	}

	return nil
}

func (s *ManageService) GetUserData(userID int) (string, error) {
	user, err := s.dbController.GetUser(userID)
	if err != nil {
		slog.Warn("Error getting user", "userID", userID, "error", err)
		return "", ErrGetUserData
	}

	userData := models.UserDataJSON{Name: user.Name, PlayedGames: user.PlayedGames, WonGames: user.WonGames}
	result, err := json.Marshal(userData)
	if err != nil {
		slog.Warn("Error marshalling userdata", "error", err, "user", user, "userData", userData)
		return "", ErrMarshalUserData
	}

	return string(result), nil
}

func (s *ManageService) UpdateUserData(userID int, oldPassword string, newPassword string, newName string) error {
	user, err := s.dbController.GetUser(userID)
	if err != nil {
		slog.Warn("Error getting user", "userID", userID, "error", err)
		return ErrGetUserData
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		slog.Warn("Invalid password", "userID", userID, "error", err)
		return ErrIncorrectPassword
	}

	preparedPassword := strings.TrimSpace(newPassword)
	if len(preparedPassword) != 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(preparedPassword), 5)
		if err != nil {
			slog.Warn("Error hashing password", "error", err, "userID", userID)
			return ErrHashingPassword
		}

		slog.Info("Updating password", "UserID", userID)
		user.Password = string(hashedPassword)
	}

	if len(newName) != 0 {
		slog.Info("Updating username", "UserID", userID)
		user.Name = newName
	}

	err = s.dbController.UpdateUser(user)
	if err != nil {
		slog.Warn("Error updating user", "userID", userID, "error", err, "user", user)
		return ErrUpdateUser
	}

	return nil
}

func (s *ManageService) JoinGame(code string, userID int) (int, error) {
	slog.Info("Search game by code", "code", code, "userID", userID)
	game, err := s.dbController.GetGameByInviteCode(code)
	if err != nil {
		slog.Warn("Can't find game by code", "code", code, "userID", userID, "err", err)
		return 0, ErrFindGame
	}

	slog.Info("Join game by code", "code", code, "userID", userID)
	err = s.dbController.JoinGame(userID, game.ID)
	if err != nil {
		slog.Warn("Can't join game", "code", code, "userID", userID, "err", err)
		return 0, ErrJoinGame
	}

	slog.Info("Successfully join game", "code", code, "userID", userID, "game", game)
	return game.ID, nil
}

func (s *ManageService) CreateGame(userID int, packID int, title string, maxPlayers int) (int, error) {
	slog.Info("Search current game by master id", "userID", userID)
	_, err := s.dbController.GetCurrentGameByMasterID(userID)
	if err == nil {
		slog.Warn("Player already in game", "userID", userID)
		return 0, ErrPlayerAlreadyInGame
	}

	slog.Info("Get pack by packID", "packID", packID)
	pack, err := s.dbController.GetPack(packID)
	if err != nil {
		slog.Warn("Can't find pack by id", "packID", packID, "err", err)
		return 0, ErrGetPack
	}

	slog.Info("Generating sample", "pack", pack, "userID", userID)
	sample, err := utils.GenerateSample(pack, s.Cfg)
	if err != nil {
		slog.Warn("Can't generate sample", "pack", pack, "userID", userID, "err", err)
		return 0, ErrGenerateSample
	}

	slog.Info("Adding sample to database", "userID", userID, "sample", sample)
	sampleID, err := s.dbController.AddSample(sample)
	if err != nil {
		slog.Warn("Can't add sample to database", "userID", userID, "sample", sample, "err", err)
		return 0, ErrAddSample
	}

	slog.Info("GetInvites", "userID", userID)
	invitesList, err := s.dbController.GetInvites()
	if err != nil {
		slog.Warn("Can't get invites", "userID", userID, "err", err)
		return 0, ErrGenerateInvite
	}

	slog.Info("Generating invite code", "inviteList", invitesList, "userID", userID)
	inviteCode, err := utils.GenerateInviteCode(invitesList)
	if err != nil {
		slog.Warn("Can't generate invite code", "userID", userID, "err", err)
		return 0, ErrGenerateInvite
	}

	slog.Info("Adding game to database", "userID", userID, "title", title, "maxPlayers", maxPlayers,
		"sample", sample, "inviteCode", inviteCode)

	gameID, err := s.dbController.AddGame(title, inviteCode, userID, maxPlayers, sampleID)
	if err != nil {
		slog.Warn("Can't add game to database", "userID", userID, "err", err)
		return 0, ErrAddGame
	}

	slog.Info("Successfully created game", "userID", userID, "gameID", gameID)
	return gameID, nil
}

func (s *ManageService) AddPack(userID int, file multipart.File, header *multipart.FileHeader) error {
	slog.Info("Saving pack file", "userID", userID, "header", header)
	filename, err := utils.SavePackGame(s.Cfg, file, header)
	if err != nil {
		slog.Warn("Can't save pack file", "userID", userID, "err", err)
		return ErrFileSave
	}

	slog.Info("Adding pack to database", "userID", userID, "file", filename)
	err = s.dbController.AddPack(userID, filename)
	if err != nil {
		slog.Warn("Error adding pack to database", "userID", userID, "err", err)
		return ErrAddPack
	}

	slog.Info("Successfully added pack to database", "userID", userID, "file", filename)
	return nil
}

func (s *ManageService) GetAllPacks(userID int) (*[]models.QuestionPackJSON, error) {
	slog.Info("Get all packs", "userID", userID)
	userPacks, err := s.dbController.GetUserPacks(userID)
	if err != nil {
		slog.Warn("Can't get all packs", "userID", userID, "err", err)
		return nil, ErrGetPack
	}

	jsonPacks := make([]models.QuestionPackJSON, 0, 5)

	for _, val := range *userPacks {
		jsonPacks = append(jsonPacks, models.QuestionPackJSON{ID: val.ID, Title: val.Title, IsOwner: userID == val.Owner})
	}

	slog.Info("Successfully get all packs", "userID", userID, "jsonPacks", jsonPacks)
	return &jsonPacks, nil
}

func (s *ManageService) GetServerPacks(userID int) (*[]models.QuestionPackJSON, error) {
	slog.Info("Get server packs", "userID", userID)
	userPacks, err := s.dbController.GetServerPacks(userID)
	if err != nil {
		slog.Warn("Can't get server packs", "userID", userID, "err", err)
		return nil, ErrGetPack
	}

	jsonPacks := make([]models.QuestionPackJSON, 0, 5)

	for _, val := range *userPacks {
		jsonPacks = append(jsonPacks, models.QuestionPackJSON{ID: val.ID, Title: val.Title, IsOwner: userID == val.Owner})
	}

	slog.Info("Successfully get server packs", "userID", userID, "jsonPacks", jsonPacks)
	return &jsonPacks, nil
}

func (s *ManageService) GetPack(packID int) (string, string, error) {
	slog.Info("Get pack file", "packID", packID)
	pack, err := s.dbController.GetPack(packID)
	if err != nil {
		slog.Warn("Can't find pack by id", "packID", packID, "err", err)
		return "", "", ErrGetPack
	}

	filename := fmt.Sprintf("%s%s", s.Cfg.Global.CsvPath, pack.Filename)

	slog.Info("Successfully getting file from csv", "file", filename)

	slog.Info("Reading file", "path", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		slog.Warn("Error reading file", "path", filename, "error", err)
		return "", "", ErrReadingFile
	}

	return string(content), pack.Title, nil
}

func (s *ManageService) UpdatePackContent(packID, userID int, content string) error {
	slog.Info("Get pack file", "packID", packID)
	pack, err := s.dbController.GetPack(packID)
	if err != nil {
		slog.Warn("Can't find pack by id", "packID", packID, "err", err)
		return ErrGetPack
	}

	if pack.Owner != userID {
		slog.Warn("user is not owner", "userID", userID, "packID", packID)
		return ErrUserNotOwner
	}

	filename := fmt.Sprintf("%s%s", s.Cfg.Global.CsvPath, pack.Filename)

	slog.Info("Updating file in csv", "file", filename)

	err = utils.UpdateFile([]byte(content), filename)
	if err != nil {
		slog.Warn("Error updating file in csv", "file", filename, "err", err)
		return ErrUpdatePack
	}

	slog.Info("Successfully pack updated")

	return nil
}

func (s *ManageService) UpdatePackTitle(packID, userID int, newTitle string) error {
	slog.Info("Get pack file", "packID", packID)
	pack, err := s.dbController.GetPack(packID)
	if err != nil {
		slog.Warn("Can't find pack by id", "packID", packID, "err", err)
		return ErrGetPack
	}

	if pack.Owner != userID {
		slog.Warn("user is not owner", "userID", userID, "packID", packID)
		return ErrUserNotOwner
	}

	err = s.dbController.UpdatePackTitle(packID, newTitle)
	if err != nil {
		slog.Warn("Can't update pack title", "packID", packID, "err", err)
		return ErrUpdatePack
	}

	return nil
}

func (s *ManageService) DeletePack(userID int, packID int) error {
	slog.Info("Get pack", "userID", userID, "packID", packID)
	pack, err := s.dbController.GetPack(packID)
	if err != nil {
		slog.Warn("Can't find pack by id", "packID", packID, "err", err)
		return ErrGetPack
	}

	if pack.Owner != userID {
		slog.Warn("userID is not owner of pack", "userID", userID, "pack", pack)
		return ErrNotOwner
	}

	slog.Info("Deleting pack", "userID", userID, "packID", packID)
	err = s.dbController.DeletePack(packID)
	if err != nil {
		slog.Warn("Can't delete pack", "userID", userID, "packID", packID, "err", err)
		return ErrDeletePack
	}

	slog.Info("Deleting pack file", "pack", pack)
	err = utils.DeletePackGame(pack.Filename, s.Cfg)
	if err != nil {
		slog.Warn("Can't delete pack file", "packID", pack.ID, "err", err)
		return ErrDeletePackFile
	}

	return nil
}
