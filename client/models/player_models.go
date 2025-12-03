package models

type ActiveTopics struct {
	ID    int64  `json:"id"`
	Title string `json:"telegram_id"`
}

type AllTopicsResponse struct {
	Success bool            `json:"success"`
	Topics  []*ActiveTopics `json:"topics"`
}

type TenQuestionsResponse struct {
	Success   bool       `json:"success"`
	Questions []Question `json:"questions"`
	TopicId   int        `json:"topic_id"`
}
