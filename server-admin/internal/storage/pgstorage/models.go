package pgstorage

const (
	users_tableName            = "users"
	users_IDColumnName         = "id"
	users_TelegramIDColumnName = "telegram_id"
	users_NameColumnName       = "name"
	users_UsernameColumnName   = "username"
	users_IsAdminColumnName    = "is_admin"
	users_CreatedAtColumnName  = "created_at"

	rating_tableName                = "rating"
	rating_IDColumnName             = "id"
	rating_UserIDColumnName         = "user_id"
	rating_TopicIDColumnName        = "topic_id"
	rating_ScoreColumnName          = "score"
	rating_CompletionTimeColumnName = "completion_time"
	rating_CreatedAtColumnName      = "created_at"
	rating_UpdatedAtColumnName      = "updated_at"
)
