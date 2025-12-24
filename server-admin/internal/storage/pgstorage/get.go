package pgstorage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Vladislav-Evg-Sid/quizbot/server-admin/internal/models"
	"github.com/pkg/errors"
)

func (storage *PGstorage) GetUserPermissions(ctx context.Context, tg_id int64) (*models.Permissions, error) {
	query := storage.getQueryGetUserPermissions(tg_id)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "quering error")
	}

	var permissions models.Permissions
	defer rows.Close()
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		permissions.Permissions = append(permissions.Permissions, p)
	}
	return &permissions, nil

}

func (storage *PGstorage) getQueryGetUserPermissions(tg_id int64) squirrel.Sqlizer {
	q := squirrel.Select(
		fmt.Sprintf(`%s.%s`, permisions_tableName, permisions_NameColumnName),
	).From(
		permisions_tableName,
	).Join(
		fmt.Sprintf(`%s ON %s.%s = %s.%s`,
			rolePermision_tableName, rolePermision_tableName, rolePermision_PermisionIDColumnName, permisions_tableName, permisions_IDColumnName,
		),
	).Join(
		fmt.Sprintf(`%s ON %s.%s = %s.%s`,
			roles_tableName, roles_tableName, roles_IDColumnName, rolePermision_tableName, rolePermision_RoleIDColumnName,
		),
	).Join(
		fmt.Sprintf(`%s ON %s.%s = %s.%s`,
			users_tableName, users_tableName, users_RoleIDColumnName, roles_tableName, roles_IDColumnName,
		),
	).Where(squirrel.Eq{
		fmt.Sprintf("%s.%s", users_tableName, users_TelegramIDColumnName): tg_id,
	}).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetUserIDByTelegramID(ctx context.Context, tg_id int64) (int, error) {
	query := getQueryGetUserIDByTelegramID(tg_id)
	queryText, args, err := query.ToSql()
	if err != nil {
		return -1, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return -1, errors.Wrap(err, "quering error")
	}
	var user_id int
	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&user_id)
	defer rows.Close()
	if err != nil {
		return -1, errors.Wrap(err, "failed to scan row")
	}
	return user_id, nil
}

func getQueryGetUserIDByTelegramID(tg_id int64) squirrel.Sqlizer {
	q := squirrel.Select(
		users_IDColumnName,
	).From(
		users_tableName,
	).Where(
		squirrel.Eq{
			users_TelegramIDColumnName: tg_id,
		},
	).PlaceholderFormat(squirrel.Dollar)
	return q
}
