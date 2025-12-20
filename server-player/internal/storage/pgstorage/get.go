package pgstorage

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/pkg/errors"
)

func (storage *PGstorage) GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error) {
	query := storage.getQuery()
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "quering error")
	}
	var topics []*models.ActiveTopics
	for rows.Next() {
		var t models.ActiveTopics
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		topics = append(topics, &t)
	}
	return topics, nil
}

func (storage *PGstorage) getQuery() squirrel.Sqlizer {
	q := squirrel.Select(topics_IDColumnName, topics_TitleColumnName).
		From(topics_tableName).
		Where(squirrel.Eq{topics_IsActiveColumnName: true}).
		PlaceholderFormat(squirrel.Dollar)
	return q
}
