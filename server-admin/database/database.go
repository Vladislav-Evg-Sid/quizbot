package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"quiz-server-admin/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("✅ Connected to Admin PostgreSQL database")
	return createTables()
}

func createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        telegram_id BIGINT UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL,
        username VARCHAR(100),
        is_admin BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS rating (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
		topic_id INTEGER NOT NULL,
        score SMALLINT NOT NULL,
		completion_time INTEGER,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users and rating table: %v", err)
	}

	log.Println("✅ Admin database tables created/verified")
	return nil
}

func CreateOrUpdateUser(user *models.User) error {
	query := `
    INSERT INTO users (telegram_id, name, username, is_admin)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (telegram_id) 
    DO UPDATE SET name = $2, username = $3
    RETURNING id, is_admin
    `

	err := DB.QueryRow(query, user.TelegramID, user.Name, user.Username, user.IsAdmin).Scan(&user.ID, &user.IsAdmin)
	return err
}

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, telegram_id, name, username, is_admin FROM users WHERE telegram_id = $1`

	err := DB.QueryRow(query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Name, &user.Username, &user.IsAdmin,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}
