package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

func SelectRandomValues(input []int, count int) ([]int, error) {
	if len(input) < count {
		return nil, ErrNotEnoughValuesToSelect
	}

	slice := make([]int, len(input))
	copy(slice, input)

	for i := range slice {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice[:count], nil
}

func GeneratePackFilename() (string, error) {
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

	return "", ErrFilenameGeneration
}

func GenerateInviteCode(codesList []string) (string, error) {
	codesMap := make(map[string]bool)
	for _, code := range codesList {
		codesMap[code] = true
	}

	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 6)

	for range 10000 {
		for i := range code {
			code[i] = letterBytes[rand.IntN(len(letterBytes))]
		}

		_, exists := codesMap[string(code)]
		if exists {
			continue
		}

		return string(code), nil
	}

	return "", ErrGenerateInviteCode
}
