package models

type User struct {
	ID         int64  `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	IsAdmin    bool   `json:"is_admin"`
}

type StartRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
}

type GreetResponse struct {
	Success  bool   `json:"success"`
	Greeting string `json:"greeting"`
	UserType string `json:"user_type"`
}

type StartResponse struct {
	Success bool `json:"success"`
	User    User `json:"user"`
}
