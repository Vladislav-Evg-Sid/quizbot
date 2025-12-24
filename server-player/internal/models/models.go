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

type QuizRequest struct {
	TgID    int64 `json:"tg_id"`
	ThemaID int32 `json:"thema_id"`
	Score   int32 `json:"score"`
	Time    int32 `json:"time"`
}
