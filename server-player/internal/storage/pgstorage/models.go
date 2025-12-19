package pgstorage

const (
	topics_tableName          = "topics"
	topics_IDColumnName       = "id"
	topics_TitleColumnName    = "title"
	topics_IsActiveColumnName = "is_active"
	topics_CreateAtColumnName = "created_at"

	questions_tableName                    = "questions"
	questions_IDColumnName                 = "id"
	questions_TopicIDColumnName            = "topic_id"
	questions_QuestionTextColumnName       = "question_text"
	questions_AnswersColumnName            = "answers"
	questions_CorrectOptionIndexColumnName = "correct_option_index"
	questions_HardLevelColumnName          = "hard_level"
	questions_CreatedAtColumnName          = "created_at"
)
