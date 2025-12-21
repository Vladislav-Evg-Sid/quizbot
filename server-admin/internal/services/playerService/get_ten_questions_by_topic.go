package playerService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	"github.com/pkg/errors"
)

func (s *PlayerService) GetTenQuestionsByTopic(ctx context.Context, topic_name string) ([]*models.Question, int, error) {
	topic_id, err := s.playerStorage.GetTopicIdByName(ctx, topic_name)
	if err != nil {
		return nil, -1, errors.Wrap(err, "founding topic id by topic name")
	}
	questions, err := s.playerStorage.GetTenQuestionsByTopicID(ctx, topic_id)
	if err != nil {
		return nil, -1, err
	}
	return questions, topic_id, nil
}
