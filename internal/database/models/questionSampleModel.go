package models

type QuestionSample struct {
	ID      int
	Pack    int
	Content string
}

type QuestionJSON struct {
	QuestionID int    `json:"questionid"`
	Level      int    `json:"level"`
	Type       string `json:"type"`
	Attachment string `json:"attachment"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
}

type ThemeJSON struct {
	Title     string         `json:"title"`
	Questions []QuestionJSON `json:"questions"`
}

type QuestionSampleJSON struct {
	FirstRound  []*ThemeJSON `json:"firstround"`
	SecondRound []*ThemeJSON `json:"secondround"`
	FinalRound  *ThemeJSON   `json:"finalround"`
}
