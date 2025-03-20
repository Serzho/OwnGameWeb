package handlers

import "github.com/pkg/errors"

var (
	ErrNotFoundInContext = errors.New("not found in context")
	ErrIncorrectType     = errors.New("incorrect type")
	ErrKeyNotFoundInJSON = errors.New("key not found in json")
)
