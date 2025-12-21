package models

type User struct {
	ID         int64  `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	IsAdmin    bool   `json:"is_admin"`
}
