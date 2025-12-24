package pgstorage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/pkg/errors"
)

func (storage *PGstorage) SetUserData(ctx context.Context, tg_id int64, name, username string) (*models.User, error) {
	query := storage.getQuerySetUserData(tg_id, name, username)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	var user models.User
	user.TelegramID = tg_id
	user.Name = name
	user.Username = username

	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&user.ID, &user.RoleID)

	if err != nil {
		return nil, errors.Wrap(err, "quering error")
	}

	return &user, nil
}

func (storage *PGstorage) getQuerySetUserData(tg_id int64, name, username string) squirrel.Sqlizer {
	q := squirrel.Insert(
		users_tableName,
	).Columns(
		users_TelegramIDColumnName,
		users_NameColumnName,
		users_UsernameColumnName,
		users_RoleIDColumnName,
	).Values(
		tg_id,
		name,
		username,
		2,
	).Suffix(fmt.Sprintf(`
			ON CONFLICT (%v) 
			DO UPDATE SET %v = EXCLUDED.%v, %v = EXCLUDED.%v
			RETURNING %v, %v
		`,
		users_TelegramIDColumnName, users_NameColumnName, users_NameColumnName, users_UsernameColumnName, users_UsernameColumnName, users_IDColumnName, users_RoleIDColumnName,
	)).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) UpsertQuizResult(ctx context.Context, quizResult *models.QuizRequest, user_id int) error {
	query := storage.getQueryUpsertQuizResult(quizResult, user_id)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		err = errors.Wrap(err, "exeс query")
	}
	return err
}

func (storage *PGstorage) getQueryUpsertQuizResult(quizResult *models.QuizRequest, user_id int) squirrel.Sqlizer {
	q := squirrel.Insert(
		rating_tableName,
	).Columns(
		rating_UserIDColumnName, rating_TopicIDColumnName, rating_ScoreColumnName, rating_CompletionTimeColumnName,
	).Values(
		user_id, quizResult.ThemaID, quizResult.Score, quizResult.Time,
	).Suffix(fmt.Sprintf(`
			ON CONFLICT (%v, %v)
			DO UPDATE SET
				%v = EXCLUDED.%v,
				%v = EXCLUDED.%v
		`,
		rating_UserIDColumnName, rating_TopicIDColumnName, rating_ScoreColumnName, rating_ScoreColumnName, rating_CompletionTimeColumnName, rating_CompletionTimeColumnName,
	)).PlaceholderFormat(squirrel.Dollar)
	return q
}
