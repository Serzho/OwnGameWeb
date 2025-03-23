package utils

import "github.com/pkg/errors"

var (
	ErrMarshalJSON             = errors.New("marshal json failed")
	ErrUnmarshalJSON           = errors.New("unmarshal json failed")
	ErrSelectRandomValues      = errors.New("select random values failed")
	ErrThemeNotFound           = errors.New("theme not found")
	ErrParseQuestions          = errors.New("parse questions failed")
	ErrInvalidFieldType        = errors.New("invalid field type")
	ErrReadingCsv              = errors.New("reading csv failed")
	ErrOpenFile                = errors.New("open file failed")
	ErrDeleteFile              = errors.New("delete file failed")
	ErrSaveFile                = errors.New("save file failed")
	ErrFilenameGeneration      = errors.New("filename generation failed")
	ErrInvalidFileType         = errors.New("invalid file type")
	ErrReadingRequestBody      = errors.New("reading request body failed")
	ErrCreatingToken           = errors.New("creating token failed")
	ErrJWTParse                = errors.New("jwt parse failed")
	ErrNotEnoughValuesToSelect = errors.New("not enough values to select")
	ErrGenerateInviteCode      = errors.New("generate invite code failed")
	ErrCreatingFile            = errors.New("creating file failed")
	ErrEmptyRecord             = errors.New("empty record")
	ErrSelectThemes            = errors.New("select themes failed")
	ErrGetThemes               = errors.New("get themes failed")
	ErrHandleTheme             = errors.New("handle theme failed")
	ErrWritingFile             = errors.New("writing file failed")
)
