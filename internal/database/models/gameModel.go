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
}
