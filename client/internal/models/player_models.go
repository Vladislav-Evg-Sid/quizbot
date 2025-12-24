package models

type ActiveTopics struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type AllTopicsResponse struct {
	Topics []*ActiveTopics `json:"topics"`
}

type TenQuestionsResponse struct {
	Success   bool       `json:"success"`
	Questions []Question `json:"questions"`
	TopicId   int        `json:"topic_id"`
}

type FinishQuizRequest struct {
	TgID    int64 `json:"tg_id"`
	ThemaID int32 `json:"thema_id"`
	Score   int32 `json:"score"`
	Time    int32 `json:"time"`
}
