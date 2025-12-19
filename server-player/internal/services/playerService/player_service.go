package playerService

import "context"

type PlayerStorage interface {
}

type PlayerService struct {
	playerStorage PlayerStorage
}

func NewPlayerService(ctx context.Context, playerStorage PlayerStorage) *PlayerService {
	return &PlayerService{
		playerStorage: playerStorage,
	}
}
