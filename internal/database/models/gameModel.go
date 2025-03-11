package models

import "time"

type Game struct {
	Id         int
	Title      string
	Status     string
	InviteCode string
	StartTime  time.Time
	MasterId   int
	UsersIds   []int
}
