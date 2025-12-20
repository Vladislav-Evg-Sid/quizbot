package playerService

import (
	"context"

	"github.com/Vladislav-Evg-Sid/quizbot/server-player/internal/models"
)

type PlayerStorage interface { // TODO: прописать интерфейс для взаимодействия с БД
	GetAllTopics(ctx context.Context) ([]*models.ActiveTopics, error)
}

type PlayerService struct { // TODO: Прописать ограничения (например, на длинну имени)
	playerStorage PlayerStorage
}

func NewPlayerService(ctx context.Context, playerStorage PlayerStorage) *PlayerService {
	return &PlayerService{
		playerStorage: playerStorage,
	}
}
