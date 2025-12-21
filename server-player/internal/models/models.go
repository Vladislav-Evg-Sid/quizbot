package models

type ActiveTopics struct {
	ID    int64  `json:"id"`
	Title string `json:"telegram_id"`
}

type Question struct {
	Text         string   `json:"text"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correct_index"`
	Level        string   `json:"level"`
}
