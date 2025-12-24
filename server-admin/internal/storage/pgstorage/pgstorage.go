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

	err = storage.initDataRoles()
	if err != nil {
		return nil, err
	}

	err = storage.initDataPermisions()
	if err != nil {
		return nil, err
	}

	err = storage.initDataRolePermision()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v VARCHAR(100) NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v VARCHAR(100) NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v INTEGER NOT NULL,
			%v INTEGER NOT NULL,
			FOREIGN KEY (%v) REFERENCES %v(%v) ON DELETE CASCADE,
			FOREIGN KEY (%v) REFERENCES %v(%v) ON DELETE CASCADE,
			UNIQUE(%v, %v)
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v BIGINT UNIQUE NOT NULL,
			%v INTEGER NOT NULL,
			%v VARCHAR(100) NOT NULL,
			%v VARCHAR(100),
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (%v) REFERENCES %v(%v) ON DELETE RESTRICT
		);
		CREATE TABLE IF NOT EXISTS %v (
			%v SERIAL PRIMARY KEY,
			%v INTEGER NOT NULL,
			%v INTEGER NOT NULL,
			%v SMALLINT NOT NULL,
			%v INTEGER,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			%v TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (%v) REFERENCES %v(%v) ON DELETE CASCADE,
			UNIQUE(%v, %v)
		);`,
		roles_tableName, roles_IDColumnName, roles_NameColumnName,

		permisions_tableName, permisions_IDColumnName, permisions_NameColumnName,

		rolePermision_tableName, rolePermision_IDColumnName, rolePermision_RoleIDColumnName, rolePermision_PermisionIDColumnName,
		rolePermision_RoleIDColumnName, roles_tableName, roles_IDColumnName,
		rolePermision_PermisionIDColumnName, permisions_tableName, permisions_IDColumnName,
		rolePermision_RoleIDColumnName, rolePermision_PermisionIDColumnName,

		users_tableName, users_IDColumnName, users_TelegramIDColumnName, users_RoleIDColumnName,
		users_NameColumnName, users_UsernameColumnName, users_CreatedAtColumnName,
		users_RoleIDColumnName, roles_tableName, roles_IDColumnName,

		rating_tableName, rating_IDColumnName, rating_UserIDColumnName, rating_TopicIDColumnName,
		rating_ScoreColumnName, rating_CompletionTimeColumnName, rating_CreatedAtColumnName, rating_UpdatedAtColumnName,
		rating_UserIDColumnName, users_tableName, users_IDColumnName,
		rating_UserIDColumnName, rating_TopicIDColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "inition tables")
	}
	return nil
}

func (s *PGstorage) initDataRoles() error {
	sql := fmt.Sprintf(`
		INSERT INTO %v (%v, %v) VALUES
			(1, 'admin'),
			(2, 'player')
		ON conflict (%v)
		DO UPDATE SET
			%v = EXCLUDED.%v;
		`,
		roles_tableName, roles_IDColumnName, roles_NameColumnName,
		roles_IDColumnName, roles_NameColumnName, roles_NameColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "add roles data")
	}
	return nil
}

func (s *PGstorage) initDataPermisions() error {
	sql := fmt.Sprintf(`
		INSERT INTO %v (%v, %v) VALUES
			(1, 'canReadTopics'),
			(2, 'canReadQuestions'),
			(3, 'canPlayInQuiz'),
			(4, 'canReadRatings'),
			(5, 'canUpdateTopics'),
			(6, 'canDeleteTopics'),
			(7, 'canCreateTopics'),
			(8, 'canUpdateQuestions'),
			(9, 'canDeleteQuestions'),
			(10, 'canCreateQuestions')
		ON conflict (%v)
		DO UPDATE SET
			%v = EXCLUDED.%v;
		`, // TODO: 7, 8, 9 не обрабатываются
		permisions_tableName, permisions_IDColumnName, permisions_NameColumnName,
		permisions_IDColumnName, permisions_NameColumnName, permisions_NameColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "add permisions data")
	}
	return nil
}

func (s *PGstorage) initDataRolePermision() error {
	sql := fmt.Sprintf(`
		INSERT INTO %v (%v, %v) VALUES
			(1, 1),
			(2, 1),
			(1, 2),
			(2, 3),
			(1, 4),
			(2, 4),
			(1, 5),
			(1, 6),
			(1, 7),
			(1, 8),
			(1, 9),
			(1, 10)
		ON conflict (%v, %v)
		DO NOTHING;
		`,
		rolePermision_tableName, rolePermision_RoleIDColumnName, rolePermision_PermisionIDColumnName,
		rolePermision_RoleIDColumnName, rolePermision_PermisionIDColumnName,
	)
	_, err := s.db.Exec(context.Background(), sql)
	if err != nil {
		return errors.Wrap(err, "add role_permision data")
	}
	return nil
}
