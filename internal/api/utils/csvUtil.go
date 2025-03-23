package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"maps"
	"mime/multipart"
	"os"
	"strconv"

	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
)

func SavePackGame(cfg *config.Config, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Header.Get("Content-Type") != "text/csv" {
		return "", ErrInvalidFileType
	}

	filename, err := GeneratePackFilename()
	if err != nil {
		return "", ErrFilenameGeneration
	}

	filepath := fmt.Sprintf("%s%s", cfg.Global.CsvPath, filename)

	slog.Info("Create Pack File: ", "filename", filepath)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", ErrCreatingFile
	}

	slog.Info("Writing File: ", "filename", filepath)

	_, err = io.Copy(dst, file)
	if err != nil {
		_ = dst.Close()
		_ = file.Close()

		return "", ErrSaveFile
	}

	_ = file.Close()
	_ = dst.Close()

	slog.Info("Successfully pack game saved")

	return filename, nil
}

func DeletePackGame(filename string, cfg *config.Config) error {
	err := os.Remove(fmt.Sprintf("%s%s", cfg.Global.CsvPath, filename))
	if err != nil {
		return ErrDeleteFile
	}

	return nil
}

func parseRecord(record []string) (int, int, int, int, error) {
	if record[0] != "1" && record[0] != "2" && record[0] != "3" {
		return 0, 0, 0, 0, ErrEmptyRecord
	}

	round, err := strconv.Atoi(record[0])
	if err != nil {
		return 0, 0, 0, 0, ErrInvalidFieldType
	}

	themeID, err := strconv.Atoi(record[2])
	if err != nil {
		return 0, 0, 0, 0, ErrInvalidFieldType
	}

	questionID, err := strconv.Atoi(record[3])
	if err != nil {
		return 0, 0, 0, 0, ErrInvalidFieldType
	}

	level, err := strconv.Atoi(record[4])
	if err != nil {
		return 0, 0, 0, 0, ErrInvalidFieldType
	}

	return round, themeID, questionID, level, nil
}

func ParseQuestions(filename string) (map[int]map[int]models.ThemeJSON, error) {
	slog.Info("Parse Questions: ", "filename", filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, ErrOpenFile
	}

	defer func() {
		err := file.Close()
		if err != nil {
			slog.Error("Close file failed")
		}
	}()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	slog.Info("Reading file: ", "filename", filename)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, ErrReadingCsv
	}

	rounds := make(map[int]map[int]models.ThemeJSON)

	slog.Info("Parse questions: ", "records", len(records)-1)

	for _, record := range records[1:] {
		round, themeID, questionID, level, err := parseRecord(record)
		if err != nil {
			slog.Warn("Parse record", "record", record, "err", err)
			continue
		}

		themes, ok := rounds[round]
		if !ok {
			rounds[round] = map[int]models.ThemeJSON{}
			themes = rounds[round]
		}

		theme, ok := themes[themeID]
		if !ok {
			theme = models.ThemeJSON{
				Title:     record[1],
				Questions: make([]models.QuestionJSON, 0, 10),
			}
		}

		theme.Questions = append(theme.Questions, models.QuestionJSON{
			QuestionID: questionID,
			Level:      level,
			Type:       record[5],
			Attachment: record[6],
			Question:   record[7],
			Answer:     record[8],
		})
		themes[themeID] = theme
	}

	slog.Info("Successfully parsed questions")

	return rounds, nil
}

func selectThemes(
	firstRoundThemes, secondRoundThemes, thirdRoundThemes map[int]models.ThemeJSON,
) ([]int, []int, []int, error) {
	firstSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(firstRoundThemes)), 5)
	if err != nil {
		return nil, nil, nil, ErrSelectRandomValues
	}

	secondSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(secondRoundThemes)), 5)
	if err != nil {
		return nil, nil, nil, ErrSelectRandomValues
	}

	thirdSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(thirdRoundThemes)), 1)
	if err != nil {
		return nil, nil, nil, ErrSelectRandomValues
	}

	return firstSelectedThemes, secondSelectedThemes, thirdSelectedThemes, nil
}

func seqToSlice(seq iter.Seq[int]) []int {
	slice := make([]int, 0, 10)
	for val := range seq {
		slice = append(slice, val)
	}

	return slice
}

func getThemes(csvPath, filename string) (
	map[int]models.ThemeJSON, map[int]models.ThemeJSON, map[int]models.ThemeJSON, error,
) {
	rounds, err := ParseQuestions(fmt.Sprintf("%s%s", csvPath, filename))
	if err != nil {
		return nil, nil, nil, ErrParseQuestions
	}

	firstRoundThemes, ok := rounds[1]
	if !ok {
		return nil, nil, nil, ErrThemeNotFound
	}

	secondRoundThemes, ok := rounds[2]
	if !ok {
		return nil, nil, nil, ErrThemeNotFound
	}

	thirdRoundThemes, ok := rounds[3]
	if !ok {
		return nil, nil, nil, ErrThemeNotFound
	}

	return firstRoundThemes, secondRoundThemes, thirdRoundThemes, nil
}

func handleTheme(roundThemes map[int]models.ThemeJSON, themeID int) (*models.ThemeJSON, error) {
	theme, ok := roundThemes[themeID]
	if !ok {
		return nil, ErrThemeNotFound
	}

	questionIndList := make([]int, 0, len(theme.Questions))
	for i := range theme.Questions {
		questionIndList = append(questionIndList, i)
	}

	selectedQuestions, err := SelectRandomValues(questionIndList, 5)
	if err != nil {
		return nil, ErrSelectRandomValues
	}

	questionList := make([]models.QuestionJSON, 0, 5)
	for _, ind := range selectedQuestions {
		questionList = append(questionList, theme.Questions[ind])
	}

	return &models.ThemeJSON{Title: theme.Title, Questions: questionList}, nil
}

func GenerateSample(pack *models.QuestionPack, cfg *config.Config) (*models.QuestionSample, error) {
	slog.Info("Generating Sample: ", "packId", pack.ID)

	firstRoundThemes, secondRoundThemes, thirdRoundThemes, err := getThemes(cfg.Global.CsvPath, pack.Filename)
	if err != nil {
		return nil, ErrGetThemes
	}

	firstSelectedThemes, secondSelectedThemes, thirdSelectedThemes, err := selectThemes(
		firstRoundThemes, secondRoundThemes, thirdRoundThemes)
	if err != nil {
		return nil, ErrSelectThemes
	}

	firstRound := make([]*models.ThemeJSON, 0, 20)

	secondRound := make([]*models.ThemeJSON, 0, 20)

	var finalRound *models.ThemeJSON

	for _, themeID := range firstSelectedThemes {
		theme, err := handleTheme(firstRoundThemes, themeID)
		if err != nil {
			return nil, ErrHandleTheme
		}

		firstRound = append(firstRound, theme)
	}

	for _, themeID := range secondSelectedThemes {
		theme, err := handleTheme(secondRoundThemes, themeID)
		if err != nil {
			return nil, ErrHandleTheme
		}

		secondRound = append(secondRound, theme)
	}

	theme, ok := thirdRoundThemes[thirdSelectedThemes[0]]
	if !ok {
		return nil, ErrThemeNotFound
	}

	questionIndList := make([]int, 0, len(theme.Questions))
	for i := range theme.Questions {
		questionIndList = append(questionIndList, i)
	}

	selectedQuestions, err := SelectRandomValues(questionIndList, 1)
	if err != nil {
		return nil, ErrSelectRandomValues
	}

	questionList := make([]models.QuestionJSON, 0, 5)
	for _, ind := range selectedQuestions {
		questionList = append(questionList, theme.Questions[ind])
	}

	finalRound = &models.ThemeJSON{Title: theme.Title, Questions: questionList}

	content, err := json.Marshal(
		models.QuestionSampleJSON{FirstRound: firstRound, SecondRound: secondRound, FinalRound: finalRound})
	if err != nil {
		return nil, ErrMarshalJSON
	}

	slog.Info("Successfully generated sample: ", "packId", pack.ID)

	return &models.QuestionSample{
		ID:      0,
		Pack:    pack.ID,
		Content: string(content),
	}, nil
}
