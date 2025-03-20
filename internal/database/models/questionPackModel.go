package models

type QuestionPack struct {
	ID       int    `db:"id"`
	Title    string `db:"title"`
	Filename string `db:"filename"`
	Owner    int    `db:"owner"`
}

type QuestionPackJSON struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	IsOwner bool   `json:"isOwner"`
}

type QuestionPacksJSON struct {
	Packs []QuestionPackJSON `json:"packs"`
}
