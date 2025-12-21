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

	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&user.ID, &user.IsAdmin)

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
		users_IsAdminColumnName,
	).Values(
		tg_id,
		name,
		username,
		false,
	).Suffix(fmt.Sprintf(`
			ON CONFLICT (%v) 
			DO UPDATE SET %v = EXCLUDED.%v, %v = EXCLUDED.%v
			RETURNING %v, %v
		`,
		users_TelegramIDColumnName, users_NameColumnName, users_NameColumnName, users_UsernameColumnName, users_UsernameColumnName, users_IDColumnName, users_IsAdminColumnName,
	)).PlaceholderFormat(squirrel.Dollar)
	return q
}
