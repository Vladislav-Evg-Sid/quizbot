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
