package models

type Permissions struct {
	Permissions []string `json:"permissions"`
}

type User struct {
	ID         int64  `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	RoleID     int32  `json:"role_id"`
}

type QuizRequest struct {
	TgID    int64 `json:"tg_id"`
	ThemaID int32 `json:"thema_id"`
	Score   int32 `json:"score"`
	Time    int32 `json:"time"`
}
