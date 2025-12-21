package models

import "time"

type Question struct {
	Text         string   `json:"text"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correct_index"`
	Level        string   `json:"level"`
}

type GameSession struct {
	SessionID            string     `json:"session_id"`
	UserID               int64      `json:"user_id"`
	TopicID              int        `json:"topic_id"`
	CurrentQuestionIndex int        `json:"current_question_index"`
	Score                int        `json:"score"`
	TotalTime            int        `json:"total_time"`
	StartedAt            time.Time  `json:"started_at"`
	Questions            []Question `json:"questions"`
}
