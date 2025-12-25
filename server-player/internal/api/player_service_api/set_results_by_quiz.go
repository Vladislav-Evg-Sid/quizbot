package playerserviceapi

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
	players_api "github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/pb/players_api"
)

func (a *PlayerServiceAPI) SetResultsByQuiz(ctx context.Context, req *players_api.SetResultsByQuizRequest) (*players_api.SetResultsByQuizResponce, error) {
	quizReq := &models.QuizRequest{
		TgID:    req.TgId,
		ThemaID: req.ThemaId,
		Score:   req.Score,
		Time:    req.Time,
	}
	err := a.playerService.SetResultsByQuiz(ctx, quizReq)
	if err != nil {
		return nil, err
	}
	return &players_api.SetResultsByQuizResponce{}, nil
}
