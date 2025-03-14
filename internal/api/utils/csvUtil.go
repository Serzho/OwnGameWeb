package utils

import (
	"OwnGameWeb/config"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"os"
	"time"
)

func getPackFilename() (string, error) {
	for range 1000000 {
		hashInp := fmt.Sprintf("%d%s", rand.IntN(10000000), time.Now())
		h := sha256.New()
		h.Write([]byte(hashInp))
		hashBytes := h.Sum(nil)

		filename := hex.EncodeToString(hashBytes)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename, nil
		}
	}
	return "", errors.New("cannot generate correct filename")
}

func SavePackGame(cfg *config.Config, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Header.Get("Content-Type") != "text/csv" {
		return "", errors.New("invalid file type")
	}

	filename, err := getPackFilename()

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
