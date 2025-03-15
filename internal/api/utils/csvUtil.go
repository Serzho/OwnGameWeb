package utils

import (
	"OwnGameWeb/config"
	"OwnGameWeb/internal/database/models"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"mime/multipart"
	"os"
	"strconv"
)

func SavePackGame(cfg *config.Config, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Header.Get("Content-Type") != "text/csv" {
		return "", errors.New("invalid file type")
	}

	filename, err := GeneratePackFilename()

	if err != nil {
		return "", errors.New("filename generation error")
	}

	filepath := fmt.Sprintf("%s%s", cfg.Global.CsvPath, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", errors.New("create file error")
	}

	_, err = io.Copy(dst, file)

	if err != nil {
		_ = dst.Close()
		_ = file.Close()
		return "", errors.New("save file error")
	}

	_ = file.Close()
	_ = dst.Close()

	return filename, nil
}

func DeletePackGame(filename string, cfg *config.Config) error {
	err := os.Remove(fmt.Sprintf("%s%s", cfg.Global.CsvPath, filename))
	if err != nil {
		return errors.New("deleting file failed")
	}
	return nil
}

func ParseQuestions(filename string) (map[int]map[int]models.ThemeJson, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("open file failed")
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("close file failed")
		}
	}()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1 // Разрешаем разное количество полей

	records, err := reader.ReadAll()

	if err != nil {
		return nil, errors.New("error reading csv")
	}

	rounds := make(map[int]map[int]models.ThemeJson)

	for _, record := range records[1:] {
		if record[0] != "1" && record[0] != "2" && record[0] != "3" {
			continue
		}
		round, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, errors.New("invalid type of field")
		}

		themeId, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, errors.New("invalid type of field")
		}

		questionId, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, errors.New("invalid type of field")
		}

		level, err := strconv.Atoi(record[4])
		if err != nil {
			return nil, errors.New("invalid type of field")
		}

		themes, ok := rounds[round]
		if !ok {
			rounds[round] = map[int]models.ThemeJson{}
			themes = rounds[round]
		}

		th, ok := themes[themeId]
		if !ok {
			th = models.ThemeJson{
				Title:     record[1],
				Questions: make([]models.QuestionJson, 0, 10),
			}
		}

		th.Questions = append(th.Questions, models.QuestionJson{
			QuestionId: questionId,
			Level:      level,
			Type:       record[5],
			Attachment: record[6],
			Question:   record[7],
			Answer:     record[8],
		})
		themes[themeId] = th

	}

	return rounds, nil
}

func GenerateSample(pack *models.QuestionPack, cfg *config.Config) (*models.QuestionSample, error) {
	seqToSlice := func(seq iter.Seq[int]) []int {
		slice := make([]int, 0, 10)
		for val := range seq {
			slice = append(slice, val)
		}
		return slice
	}

	rounds, err := ParseQuestions(fmt.Sprintf("%s%s", cfg.Global.CsvPath, pack.Filename))
	if err != nil {
		return nil, errors.New("parse questions failed")
	}

	firstRoundThemes, ok := rounds[1]
	if !ok {
		return nil, errors.New("first round theme not found")
	}

	firstSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(firstRoundThemes)), 5)
	if err != nil {
		return nil, errors.New("select random values failed")
	}

	secondRoundThemes, ok := rounds[2]
	if !ok {
		return nil, errors.New("second round theme not found")
	}

	secondSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(secondRoundThemes)), 5)
	if err != nil {
		return nil, errors.New("select random values failed")
	}

	thirdRoundThemes, ok := rounds[3]
	if !ok {
		return nil, errors.New("third round theme not found")
	}

	thirdSelectedThemes, err := SelectRandomValues(seqToSlice(maps.Keys(thirdRoundThemes)), 1)
	if err != nil {
		return nil, errors.New("select random values failed")
	}

	var firstRound []*models.ThemeJson
	var secondRound []*models.ThemeJson
	var finalRound *models.ThemeJson

	for _, themeId := range firstSelectedThemes {
		theme, ok := firstRoundThemes[themeId]
		if !ok {
			return nil, errors.New("first round theme not found")
		}
		questionIndList := make([]int, 0, len(theme.Questions))
		for i := range theme.Questions {
			questionIndList = append(questionIndList, i)
		}

		selectedQuestions, err := SelectRandomValues(questionIndList, 5)
		if err != nil {
			return nil, errors.New("select random values failed")
		}
		questionList := make([]models.QuestionJson, 0, 5)
		for _, ind := range selectedQuestions {
			questionList = append(questionList, theme.Questions[ind])
		}

		firstRound = append(firstRound, &models.ThemeJson{Title: theme.Title, Questions: questionList})
	}

	for _, themeId := range secondSelectedThemes {
		theme, ok := secondRoundThemes[themeId]
		if !ok {
			return nil, errors.New("first round theme not found")
		}
		questionIndList := make([]int, 0, len(theme.Questions))
		for i := range theme.Questions {
			questionIndList = append(questionIndList, i)
		}

		selectedQuestions, err := SelectRandomValues(questionIndList, 5)
		if err != nil {
			return nil, errors.New("select random values failed")
		}
		questionList := make([]models.QuestionJson, 0, 5)
		for _, ind := range selectedQuestions {
			questionList = append(questionList, theme.Questions[ind])
		}

		secondRound = append(secondRound, &models.ThemeJson{Title: theme.Title, Questions: questionList})
	}

	themeId := thirdSelectedThemes[0]
	theme, ok := thirdRoundThemes[themeId]
	if !ok {
		return nil, errors.New("first round theme not found")
	}
	questionIndList := make([]int, 0, len(theme.Questions))
	for i := range theme.Questions {
		questionIndList = append(questionIndList, i)
	}

	selectedQuestions, err := SelectRandomValues(questionIndList, 1)
	if err != nil {
		return nil, errors.New("select random values failed")
	}
	questionList := make([]models.QuestionJson, 0, 5)
	for _, ind := range selectedQuestions {
		questionList = append(questionList, theme.Questions[ind])
	}

	finalRound = &models.ThemeJson{Title: theme.Title, Questions: questionList}

	content, err := json.Marshal(
		models.QuestionSampleJson{FirstRound: firstRound, SecondRound: secondRound, FinalRound: finalRound})
	if err != nil {
		return nil, errors.New("marshal json failed")
	}

	return &models.QuestionSample{
		Id:      0,
		Pack:    pack.Id,
		Content: string(content),
	}, nil
}
