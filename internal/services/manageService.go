package services

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/api/utils"
	"OwnGameWeb/internal/database"
	"OwnGameWeb/internal/database/models"
	"errors"
	"fmt"
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
	game, err := s.dbController.GetGameByInviteCode(code)
	if err != nil {
		return 0, errors.New("cannot find game")
	}

	err = s.dbController.JoinGame(userId, game.Id)

	if err != nil {
		return 0, errors.New("failed to join game")
	}
	return game.Id, nil
}

func (s *ManageService) CreateGame(userId int, packId int, title string, maxPlayers int) (int, error) {
	_, err := s.dbController.GetCurrentGameByMasterId(userId)
	if err == nil {
		return 0, errors.New("player already playing")
	}

	pack, err := s.dbController.GetPack(packId)

	if err != nil {
		return 0, errors.New("pack not found")
	}

	sample, err := utils.GenerateSample(pack, s.Cfg)
	if err != nil {
		return 0, errors.New("generate sample failed")
	}

	sampleId, err := s.dbController.AddSample(sample)
	if err != nil {
		return 0, errors.New("add sample failed")
	}

	invitesList, err := s.dbController.GetInvites()
	if err != nil {
		return 0, errors.New("get invites failed")
	}

	inviteCode, err := utils.GenerateInviteCode(invitesList)
	if err != nil {
		return 0, errors.New("generate invite code failed")
	}
	gameId, err := s.dbController.AddGame(title, inviteCode, userId, maxPlayers, sampleId)

	if err != nil {
		return 0, err
	}
	return gameId, nil

}

func (s *ManageService) AddPack(userId int, file multipart.File, header *multipart.FileHeader) error {
	filename, err := utils.SavePackGame(s.Cfg, file, header)
	if err != nil {
		return errors.New("file save failed")
	}

	err = s.dbController.AddPack(userId, filename)
	if err != nil {
		return errors.New("database add pack failed")
	}

	return nil
}

func (s *ManageService) GetAllPacks(userId int) (*[]models.QuestionPackJson, error) {
	userPacks, err := s.dbController.GetUserPacks(userId)
	if err != nil {
		return nil, errors.New("get user packs failed")
	}

	jsonPacks := make([]models.QuestionPackJson, 0, 5)

	for _, val := range *userPacks {
		jsonPacks = append(jsonPacks, models.QuestionPackJson{Id: val.Id, Title: val.Title, IsOwner: userId == val.Owner})
	}

	return &jsonPacks, nil
}

func (s *ManageService) GetPackFile(packId int) (string, error) {
	pack, err := s.dbController.GetPack(packId)
	if err != nil {
		return "", errors.New("get pack from database failed")
	}

	filename := fmt.Sprintf("%s%s", s.Cfg.Global.CsvPath, pack.Filename)

	return filename, nil
}
func (s *ManageService) DeletePack(userId int, packId int) error {
	pack, err := s.dbController.GetPack(packId)
	if err != nil {
		return errors.New("get pack from database failed")
	}

	if pack.Owner != userId {
		return errors.New("you are not the owner of this pack")
	}

	err = s.dbController.DeletePack(packId)
	if err != nil {
		return errors.New("database delete failed")
	}

	err = utils.DeletePackGame(pack.Filename, s.Cfg)
	if err != nil {
		return err
	}

	return nil
}
