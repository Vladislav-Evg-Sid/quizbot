package playerserviceapi

import (
	"context"
	"log"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	players_api "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
	"github.com/samber/lo"
)

func (s *PlayerServiceAPI) GetAllTopics(ctx context.Context, req *players_api.GetAllTopicsRequest) (*players_api.GetAllTopicsResponse, error) {
	log.Print("Received request")

	responce, err := s.playerService.GetAllTopics(ctx)
	if err != nil {
		return &players_api.GetAllTopicsResponse{}, err
	}
	return &players_api.GetAllTopicsResponse{
		Topics: mapTopicsByResponce(responce),
	}, nil
}

func mapTopicsByResponce(topicsInfo []*models.ActiveTopics) []*players_api.Topic {
	return lo.Map(topicsInfo, func(t *models.ActiveTopics, _ int) *players_api.Topic {
		return &players_api.Topic{
			Id:    t.ID,
			Title: t.Title,
		}
	})
}
