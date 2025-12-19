package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStore(connString string) (*PGstorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	storage := &PGstorage{
		db: db,
	}
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v VARCHAR(100) NOT NULL UNIQUE,
			%v BOOLEAN DEFAULT TRUE,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v INTEGER,
			%v TEXT UNIQUE,
			%v JSONB,
			%v SMALLINT NOT NULL,
			%v VARCHAR(20) NOT NULL,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		topics_tableName, topics_IDColumnName, topics_TitleColumnName, topics_IsActiveColumnName, topics_CreateAtColumnName,
		questions_tableName, questions_IDColumnName, questions_TopicIDColumnName, questions_QuestionTextColumnName, questions_AnswersColumnName, questions_CorrectOptionIndexColumnName, questions_HardLevelColumnName, questions_CreatedAtColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "inition tables")
	}
	return nil
}
