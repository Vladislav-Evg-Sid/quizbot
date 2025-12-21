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
			%v BIGINT UNIQUE NOT NULL,
			%v VARCHAR(100) NOT NULL,
			%v VARCHAR(100),
			%v BOOLEAN DEFAULT FALSE,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v INTEGER NOT NULL,
			%v INTEGER NOT NULL,
			%v SMALLINT NOT NULL,
			%v INTEGER,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		users_tableName, users_IDColumnName, users_TelegramIDColumnName, users_NameColumnName, users_UsernameColumnName, users_IsAdminColumnName, users_CreatedAtColumnName,
		rating_tableName, rating_IDColumnName, rating_UserIDColumnName, rating_TopicIDColumnName, rating_ScoreColumnName, rating_CompletionTimeColumnName, rating_CreatedAtColumnName, rating_UpdatedAtColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "inition tables")
	}
	return nil
}
