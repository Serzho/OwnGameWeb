package database

import "github.com/pkg/errors"

var (
	ErrGetSample           = errors.New("get sample failed")
	ErrDeleteGame          = errors.New("delete game failed")
	ErrSetGameStatus       = errors.New("set game status failed")
	ErrGetGameByInviteCode = errors.New("get game by invite code failed")
	ErrGetGame             = errors.New("get game failed")
	ErrInsertSample        = errors.New("insert sample failed")
	ErrGetUserPacks        = errors.New("get user packs failed")
	ErrGetGameByMasterID   = errors.New("get game by master id failed")
	ErrInsertGame          = errors.New("insert game failed")
	ErrInsertUser          = errors.New("insert user failed")
	ErrTransactionCommit   = errors.New("transaction commit failed")
	ErrGetUserByEmail      = errors.New("get user by email failed")
	ErrGetUser             = errors.New("get user failed")
	ErrGetPack             = errors.New("get pack failed")
	ErrGetServerPacks      = errors.New("get server packs failed")
	ErrAddServerPack       = errors.New("add server pack failed")
)
