package playerserviceapi

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
	"github.com/samber/lo"
)

func (s *PlayerServiceAPI) GetTenQuestionsByTopic(ctx context.Context, req *players_api.GetTenQuestionsByTopicRequest) (*players_api.GetTenQuestionsByTopicResponse, error) {
	topic_name := req.TopicName
	log.Printf("Received request, %v", topic_name)

	response, topic_id, err := s.playerService.GetTenQuestionsByTopic(ctx, topic_name)
	if err != nil {
		return &players_api.GetTenQuestionsByTopicResponse{}, err
	}

	response, err = s.SelectTenQuestions(response)
	if err != nil {
		return &players_api.GetTenQuestionsByTopicResponse{}, err
	}

	return &players_api.GetTenQuestionsByTopicResponse{
		Questions: mapTenQuestionsByResponce(response),
		TopicId:   int64(topic_id),
	}, nil
}

func (s *PlayerServiceAPI) SelectTenQuestions(questions []*models.Question) ([]*models.Question, error) {
	var tenQuestions []*models.Question
	var questionsEasy []*models.Question
	var questionsMedium []*models.Question
	var questionsHard []*models.Question

	// Разбиваем вопросы по уровню
	for _, q := range questions {
		switch {
		case q.Level == "простой":
			questionsEasy = append(questionsEasy, q)
		case q.Level == "средний":
			questionsMedium = append(questionsMedium, q)
		case q.Level == "сложный":
			questionsHard = append(questionsHard, q)
		}
	}

	// Перемешиваем вопросы
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
		tenQuestions = append(tenQuestions, questionsEasy[:4]...)
	} else {
		return nil, fmt.Errorf("not enough easy questions")
	}
	if len(questionsMedium) >= 3 {
		tenQuestions = append(tenQuestions, questionsMedium[:3]...)
	} else {
		return nil, fmt.Errorf("not enough medium questions")
	}
	if len(questionsHard) >= 3 {
		tenQuestions = append(tenQuestions, questionsHard[:3]...)
	} else {
		return nil, fmt.Errorf("not enough hard questions")
	}

	return tenQuestions, nil
}

func mapTenQuestionsByResponce(answersInfo []*models.Question) []*players_api.Question {
	return lo.Map(answersInfo, func(a *models.Question, _ int) *players_api.Question {
		return &players_api.Question{
			Text:         a.Text,
			Options:      a.Options,
			CorrectIndex: int64(a.CorrectIndex),
			Level:        a.Level,
		}
	})
}
