package models

import "time"

type Game struct {
	Id         int
	Title      string
	Status     string
	InviteCode string
	StartTime  time.Time
	MasterId   int
	PlayersIds []int
	MaxPlayers int
	Sample     string
}

type PlayerJson struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GameInfoJson struct {
	Title      string       `json:"title"`
	Players    []PlayerJson `json:"players"`
	MaxPlayers int          `json:"maxPlayers"`
	IsHost     bool         `json:"isHost"`
}
