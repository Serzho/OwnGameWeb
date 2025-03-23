package models

type User struct {
	ID          int    `db:"id"`
	Email       string `db:"email"`
	Name        string `db:"name"`
	Password    string `db:"password"`
	Packs       []int  `db:"packs"`
	PlayedGames int    `db:"played_games"`
	WonGames    int    `db:"won_games"`
}

type UserDataJSON struct {
	Name        string `json:"name"`
	PlayedGames int    `json:"playedGames"`
	WonGames    int    `json:"wonGames"`
}
