package services

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
)

type ManageService struct {
	dbController *database.DbController
	Cfg          *config.Config
}

func NewManageService(c *database.DbController, cfg *config.Config) *ManageService {
	return &ManageService{dbController: c, Cfg: cfg}
}

func (s *ManageService) JoinGame(code string, userId int) (int, error) {
	slog.Info("Search game by code", "code", code, "userId", userId)
	game, err := s.dbController.GetGameByInviteCode(code)
	if err != nil {
		slog.Warn("Can't find game by code", "code", code, "userId", userId, "err", err)
		return 0, errors.New("cannot find game")
	}

	slog.Info("Join game by code", "code", code, "userId", userId)
	err = s.dbController.JoinGame(userId, game.Id)

	if err != nil {
		slog.Warn("Can't join game", "code", code, "userId", userId, "err", err)
		return 0, errors.New("failed to join game")
	}

	slog.Info("Successfully join game", "code", code, "userId", userId, "game", game)
	return game.Id, nil
}

func (s *ManageService) CreateGame(userId int, packId int, title string, maxPlayers int) (int, error) {
	slog.Info("Search current game by master id", "userId", userId)
	_, err := s.dbController.GetCurrentGameByMasterId(userId)
	if err == nil {
		slog.Warn("Player already in game", "userId", userId)
		return 0, errors.New("player already playing")
	}

	slog.Info("Get pack by packId", "packId", packId)
	pack, err := s.dbController.GetPack(packId)

	if err != nil {
		slog.Warn("Can't find pack by id", "packId", packId, "err", err)
		return 0, errors.New("pack not found")
	}

	slog.Info("Generating sample", "pack", pack, "userId", userId)
	sample, err := utils.GenerateSample(pack, s.Cfg)
	if err != nil {
		slog.Warn("Can't generate sample", "pack", pack, "userId", userId, "err", err)
		return 0, errors.New("generate sample failed")
	}

	slog.Info("Adding sample to database", "userId", userId, "sample", sample)
	sampleId, err := s.dbController.AddSample(sample)
	if err != nil {
		slog.Warn("Can't add sample to database", "userId", userId, "sample", sample, "err", err)
		return 0, errors.New("add sample failed")
	}

	slog.Info("GetInvites", "userId", userId)
	invitesList, err := s.dbController.GetInvites()
	if err != nil {
		slog.Warn("Can't get invites", "userId", userId, "err", err)
		return 0, errors.New("get invites failed")
	}

	slog.Info("Generating invite code", "inviteList", invitesList, "userId", userId)
	inviteCode, err := utils.GenerateInviteCode(invitesList)
	if err != nil {
		slog.Warn("Can't generate invite code", "userId", userId, "err", err)
		return 0, errors.New("generate invite code failed")
	}

	slog.Info("Adding game to database", "userId", userId, "title", title, "maxPlayers", maxPlayers, "sample", sample, "inviteCode", inviteCode)
	gameId, err := s.dbController.AddGame(title, inviteCode, userId, maxPlayers, sampleId)

	if err != nil {
		slog.Warn("Can't add game to database", "userId", userId, "err", err)
		return 0, err
	}

	slog.Info("Successfully created game", "userId", userId, "gameId", gameId)
	return gameId, nil

}

func (s *ManageService) AddPack(userId int, file multipart.File, header *multipart.FileHeader) error {
	slog.Info("Saving pack file", "userId", userId, "header", header)
	filename, err := utils.SavePackGame(s.Cfg, file, header)
	if err != nil {
		slog.Warn("Can't save pack file", "userId", userId, "err", err)
		return errors.New("file save failed")
	}

	slog.Info("Adding pack to database", "userId", userId, "file", filename)
	err = s.dbController.AddPack(userId, filename)
	if err != nil {
		slog.Warn("Error adding pack to database", "userId", userId, "err", err)
		return errors.New("database add pack failed")
	}

	slog.Info("Successfully added pack to database", "userId", userId, "file", filename)
	return nil
}

func (s *ManageService) GetAllPacks(userId int) (*[]models.QuestionPackJson, error) {
	slog.Info("Get all packs", "userId", userId)
	userPacks, err := s.dbController.GetUserPacks(userId)
	if err != nil {
		slog.Warn("Can't get all packs", "userId", userId, "err", err)
		return nil, errors.New("get user packs failed")
	}

	jsonPacks := make([]models.QuestionPackJson, 0, 5)

	for _, val := range *userPacks {
		jsonPacks = append(jsonPacks, models.QuestionPackJson{Id: val.Id, Title: val.Title, IsOwner: userId == val.Owner})
	}

	slog.Info("Successfully get all packs", "userId", userId, "jsonPacks", jsonPacks)
	return &jsonPacks, nil
}

func (s *ManageService) GetPackFile(packId int) (string, error) {
	slog.Info("Get pack file", "packId", packId)
	pack, err := s.dbController.GetPack(packId)
	if err != nil {
		slog.Warn("Can't find pack by id", "packId", packId, "err", err)
		return "", errors.New("get pack from database failed")
	}

	filename := fmt.Sprintf("%s%s", s.Cfg.Global.CsvPath, pack.Filename)

	slog.Info("Successfully getting file from csv", "file", filename)
	return filename, nil
}

func (s *ManageService) DeletePack(userId int, packId int) error {
	slog.Info("Get pack", "userId", userId, "packId", packId)
	pack, err := s.dbController.GetPack(packId)
	if err != nil {
		slog.Warn("Can't find pack by id", "packId", packId, "err", err)
		return errors.New("get pack from database failed")
	}

	if pack.Owner != userId {
		slog.Warn("userId is not owner of pack", "userId", userId, "pack", pack)
		return errors.New("you are not the owner of this pack")
	}

	slog.Info("Deleting pack", "userId", userId, "packId", packId)
	err = s.dbController.DeletePack(packId)
	if err != nil {
		slog.Warn("Can't delete pack", "userId", userId, "packId", packId, "err", err)
		return errors.New("database delete failed")
	}

	slog.Info("Deleting pack file", "pack", pack)
	err = utils.DeletePackGame(pack.Filename, s.Cfg)
	if err != nil {
		slog.Warn("Can't delete pack file", "packId", pack.Id, "err", err)
		return err
	}

	return nil
}
