package playerinfoprocessor

type playerService interface { // TODO: Прописать интерфейс
}

type PlayerInfoProcessor struct {
	playerService playerService
}

func NewPlayersInfoProcessor(playerService playerService) *PlayerInfoProcessor {
	return &PlayerInfoProcessor{
		playerService: playerService,
	}
}
