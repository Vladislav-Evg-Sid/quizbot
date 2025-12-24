package quizresultprocessor

type playerService interface {
}

type QuizResultProcessor struct {
	playerService playerService
}

func NewQuizResultProcessor(playerService playerService) *QuizResultProcessor {
	return &QuizResultProcessor{
		playerService: playerService,
	}
}
