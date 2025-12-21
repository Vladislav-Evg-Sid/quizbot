package pgstorage

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/pkg/errors"
)

func (storage *PGstorage) GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error) {
	query := storage.getQueryForGetAllTopics()
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "quering error")
	}
	var topics []*models.ActiveTopics
	defer rows.Close()
	for rows.Next() {
		var t models.ActiveTopics
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		topics = append(topics, &t)
	}
	return topics, nil
}

func (storage *PGstorage) getQueryForGetAllTopics() squirrel.Sqlizer {
	q := squirrel.Select(topics_IDColumnName, topics_TitleColumnName).
		From(topics_tableName).
		Where(squirrel.Eq{topics_IsActiveColumnName: true}).
		PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetTopicIdByName(ctx context.Context, topic_name string) (int, error) {
	query := storage.getQueryGetTopicIdByName(topic_name)
	queryText, args, err := query.ToSql()
	if err != nil {
		return -1, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return -1, errors.Wrap(err, "quering error")
	}
	var topic_id int
	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&topic_id)
	defer rows.Close()
	if err != nil {
		return -1, errors.Wrap(err, "failed to scan row")
	}
	return topic_id, nil
}

func (storage *PGstorage) getQueryGetTopicIdByName(topic_name string) squirrel.Sqlizer {
	q := squirrel.Select(topics_IDColumnName).
		From(topics_tableName).
		Where(squirrel.Eq{topics_TitleColumnName: topic_name}).
		PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetTenQuestionsByTopicID(ctx context.Context, topic_id int) ([]*models.Question, error) {
	query := storage.getQueryGetTenQuestionsByTopicID(topic_id)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "quering error")
	}
	var questions []*models.Question
	defer rows.Close()
	for rows.Next() {
		var q models.Question
		var answersJSON []byte
		if err := rows.Scan(&q.Text, &answersJSON, &q.CorrectIndex, &q.Level); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		if err := json.Unmarshal(answersJSON, &q.Options); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal answers")
		}
		questions = append(questions, &q)
	}
	return questions, nil
}

func (storage *PGstorage) getQueryGetTenQuestionsByTopicID(topic_id int) squirrel.Sqlizer {
	q := squirrel.Select(
		questions_QuestionTextColumnName,
		questions_AnswersColumnName,
		questions_CorrectOptionIndexColumnName,
		questions_HardLevelColumnName,
	).From(
		questions_tableName,
	).Where(
		squirrel.Eq{questions_TopicIDColumnName: topic_id},
	).PlaceholderFormat(squirrel.Dollar)
	return q
}
