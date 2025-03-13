package models

type QuestionSample struct {
	Id            int
	Pack          int
	Themes        []int
	Questions     [][]int
	FinalQuestion int
}
