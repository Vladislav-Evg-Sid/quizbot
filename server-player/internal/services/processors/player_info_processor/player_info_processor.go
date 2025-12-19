package playerinfoprocessor

type playerService interface {
}

type PlayerInfoProcessor struct {
	playerService playerService
}

func NewStudentsInfoProcessor(playerService playerService) *PlayerInfoProcessor {
	return &PlayerInfoProcessor{
		playerService: playerService,
	}
}
