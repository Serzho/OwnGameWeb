package services

import "github.com/pkg/errors"

var (
	ErrNotImplemented      = errors.New("not implemented")
	ErrCreatingUser        = errors.New("creating user failed")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrDeletePack          = errors.New("delete pack failed")
	ErrDeletePackFile      = errors.New("delete pack file failed")
	ErrNotOwner            = errors.New("not owner")
	ErrGetPack             = errors.New("get pack failed")
	ErrAddPack             = errors.New("add pack failed")
	ErrFileSave            = errors.New("file save failed")
	ErrAddGame             = errors.New("add game failed")
	ErrGenerateInvite      = errors.New("generate invite failed")
	ErrAddSample           = errors.New("add sample failed")
	ErrPlayerAlreadyInGame = errors.New("player already in game")
	ErrFindGame            = errors.New("find game failed")
	ErrJoinGame            = errors.New("join game failed")
	ErrGenerateSample      = errors.New("generate sample failed")
	ErrGetUserData         = errors.New("get user data failed")
	ErrMarshalUserData     = errors.New("marshal user data failed")
	ErrHashingPassword     = errors.New("hashing password failed")
	ErrUpdateUser          = errors.New("update user failed")
)
