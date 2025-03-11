package utils

import (
	"encoding/csv"
	"errors"
	"github.com/gin-gonic/gin"
)

func ParsePackGame(c *gin.Context) (string, error) {
	file, header, err := c.Request.FormFile("*.csv")
	if err != nil {
		return "", errors.New("CSV file is required")
	}

	if header.Header.Get("Content-Type") != "text/csv" {
		return "", errors.New("invalid file type")
	}

	reader := csv.NewReader(file)
	_, err = reader.ReadAll()
	if err != nil {
		return "", errors.New("error reading CSV")
	}
	return "", nil
}
