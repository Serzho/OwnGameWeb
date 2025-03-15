package models

type QuestionSample struct {
	Id      int
	Pack    int
	Content string
}

type QuestionJson struct {
	QuestionId int    `json:"question_id"`
	Level      int    `json:"level"`
	Type       string `json:"type"`
	Attachment string `json:"attachment"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
}

type ThemeJson struct {
	Title     string         `json:"title"`
	Questions []QuestionJson `json:"questions"`
}

type QuestionSampleJson struct {
	FirstRound  []*ThemeJson `json:"first_round"`
	SecondRound []*ThemeJson `json:"second_round"`
	FinalRound  *ThemeJson   `json:"final_round"`
}
