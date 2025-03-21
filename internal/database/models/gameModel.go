package models

import "time"

type Game struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	Status     string    `db:"status"`
	InviteCode string    `db:"invite_code"`
	StartTime  time.Time `db:"start_time"`
	MasterID   int       `db:"master_id"`
	PlayersIDs []int     `db:"players_ids"`
	MaxPlayers int       `db:"max_players"`
	Sample     int       `db:"sample"`
}

type PlayerJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GameInfoJSON struct {
	Title      string       `json:"title"`
	Players    []PlayerJSON `json:"players"`
	MaxPlayers int          `json:"maxPlayers"`
	IsHost     bool         `json:"isHost"`
	InviteCode string       `json:"inviteCode"`
	Status     string       `json:"status"`
}
