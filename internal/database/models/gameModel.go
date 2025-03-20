package models

import "time"

type Game struct {
	ID         int
	Title      string
	Status     string
	InviteCode string
	StartTime  time.Time
	MasterID   int
	PlayersIDs []int
	MaxPlayers int
	Sample     int
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
}
