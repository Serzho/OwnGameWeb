package models

type QuestionPack struct {
	Id       int    `db:"id"`
	Title    string `db:"title"`
	Filename string `db:"filename"`
	Owner    int    `db:"owner"`
}

type QuestionPackJson struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	IsOwner bool   `json:"IsOwner"`
}

type QuestionPacksJson struct {
	Packs []QuestionPackJson `json:"packs"`
}
