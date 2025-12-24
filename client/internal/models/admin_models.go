package models

type StartRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
}

type StartResponse struct {
	Success bool `json:"success"`
	User    struct {
		ID          string   `json:"id"`
		TelegramID  string   `json:"telegram_id"`
		Name        string   `json:"name"`
		Username    string   `json:"username"`
		Permissions []string `json:"permissions"`
	} `json:"user"`
}

type GreetResponse struct {
	Success  bool   `json:"success"`
	Greeting string `json:"greeting"`
	UserType string `json:"user_type"`
}

type Permissions struct {
	Permissions []string `json:"permissions"`
}
