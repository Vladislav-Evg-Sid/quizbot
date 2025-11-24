package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"quiz-server-player/models"

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
	err = createTables()
	if err != nil {
		return err
	}
	return addTopicsData()
}

func createTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS topics (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100) NOT NULL UNIQUE,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		topic_id INTEGER,
		question_text TEXT,
		answers JSONB,
		correct_option_index SMALLINT NOT NULL,
		hard_level VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	log.Println("✅ Users database tables created/verified")
	return nil
}

func addTopicsData() error {
	query := `
    INSERT INTO topics (title) VALUES 
		('История России'),
		('Всемирная история'),
		('География мира'),
		('Биология'),
		('Химия'),
		('Физика'),
		('Математика'),
		('Программирование'),
		('Музыка'),
		('Кинематограф'),
		('Живопись'),
		('Спортивные события'),
		('Кулинария'),
		('Путешествия'),
		('Космос'),
		('Техника'),
		('Литература'),
		('Мультфильмы'),
		('Видеоигры')
	ON CONFLICT (title)
	DO NOTHING;`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert topics: %v", err)
	}

	log.Println("✅ Topics inserted in table")
	return nil
}

// func CreateOrUpdateUser(user *models.User) error {
// 	query := `
//     INSERT INTO users (telegram_id, name, username, is_admin)
//     VALUES ($1, $2, $3, $4)
//     ON CONFLICT (telegram_id)
//     DO UPDATE SET name = $2, username = $3
//     RETURNING id, is_admin
//     `

// 	err := DB.QueryRow(query, user.TelegramID, user.Name, user.Username, user.IsAdmin).Scan(&user.ID, &user.IsAdmin)
// 	return err
// }

func GetAllTopics() ([]*models.ActiveTopics, error) {
	var topics []*models.ActiveTopics
	query := `SELECT id, title FROM topics WHERE is_active;`

	rows, err := DB.Query(query)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer rows.Close()
	for rows.Next() {
		topic := &models.ActiveTopics{}
		err := rows.Scan(&topic.ID, &topic.Title)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, err
}

// func GetUserByTelegramID(telegramID int64) (*models.User, error) {
// 	user := &models.User{}
// 	query := `SELECT id, telegram_id, name, username, is_admin FROM users WHERE telegram_id = $1`

// 	err := DB.QueryRow(query, telegramID).Scan(
// 		&user.ID, &user.TelegramID, &user.Name, &user.Username, &user.IsAdmin,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}

// 	return user, err
// }
