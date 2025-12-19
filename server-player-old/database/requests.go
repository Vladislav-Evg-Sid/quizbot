package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"quiz-server-player/models"
)

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
func GetRandomQuestions(topicName string) ([]models.Question, int, error) {
	var questions []models.Question
	var questionsEasy []models.Question
	var questionsMedium []models.Question
	var questionsHard []models.Question

	query := `
		SELECT
			topics.id
		FROM topics
		WHERE topics.title LIKE $1;`

	var topic_id int

	err := DB.QueryRow(query, "%"+topicName+"%").Scan(&topic_id)

	if err == sql.ErrNoRows {
		return nil, 0, fmt.Errorf("error with quiz topic")
	}

	if err != nil {
		return nil, 0, err
	}

	query = `
		SELECT
			questions.question_text,
			questions.answers,
			questions.correct_option_index,
			questions.hard_level
		FROM questions
		WHERE questions.topic_id = $1;`

	rows, err := DB.Query(query, topic_id)

	if err == sql.ErrNoRows {
		return nil, 0, fmt.Errorf("not questions in topic")
	}

	defer rows.Close()
	for rows.Next() {
		topic := &models.Question{}
		var answersJSON []byte
		err := rows.Scan(&topic.Text, &answersJSON, &topic.CorrectIndex, &topic.Level)
		if err != nil {
			return nil, 0, err
		}

		if err := json.Unmarshal(answersJSON, &topic.Options); err != nil {
			return nil, 0, err
		}

		switch {
		case topic.Level == "простой":
			questionsEasy = append(questionsEasy, *topic)
		case topic.Level == "средний":
			questionsMedium = append(questionsMedium, *topic)
		case topic.Level == "сложный":
			questionsHard = append(questionsHard, *topic)
		default:
			return nil, 0, fmt.Errorf("unnown hard level: %v", topic.Level)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Перемешиваем вопросы каждого уровня
	rand.Shuffle(len(questionsEasy), func(i, j int) {
		questionsEasy[i], questionsEasy[j] = questionsEasy[j], questionsEasy[i]
	})
	rand.Shuffle(len(questionsMedium), func(i, j int) {
		questionsMedium[i], questionsMedium[j] = questionsMedium[j], questionsMedium[i]
	})
	rand.Shuffle(len(questionsHard), func(i, j int) {
		questionsHard[i], questionsHard[j] = questionsHard[j], questionsHard[i]
	})

	// Берем нужное количество каждого уровня
	if len(questionsEasy) >= 4 {
		questions = append(questions, questionsEasy[:4]...)
	} else {
		return nil, 0, fmt.Errorf("not enough easy questions")
	}

	if len(questionsMedium) >= 3 {
		questions = append(questions, questionsMedium[:3]...)
	} else {
		return nil, 0, fmt.Errorf("not enough medium questions")
	}

	if len(questionsHard) >= 3 {
		questions = append(questions, questionsHard[:3]...)
	} else {
		return nil, 0, fmt.Errorf("not enough hard questions")
	}

	return questions, topic_id, err
}
