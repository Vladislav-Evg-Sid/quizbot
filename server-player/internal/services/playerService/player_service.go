package playerService

import "context"

type PlayerStorage interface { // TODO: прописать интерфейс для взаимодействия с БД
}

type PlayerService struct { // TODO: Прописать ограничения (например, на длинну имени)
	playerStorage PlayerStorage
}

func NewPlayerService(ctx context.Context, playerStorage PlayerStorage) *PlayerService {
	return &PlayerService{
		playerStorage: playerStorage,
	}
}
