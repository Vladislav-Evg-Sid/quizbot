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
