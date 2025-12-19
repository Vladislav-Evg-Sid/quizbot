package playerinfoupsertconsumer

type playerInfoProcessor interface {
}

type PlayerInfoUpsertConsumer struct {
	playerInfoProcessor playerInfoProcessor
}

func NewPlayerInfoUpsertConsumer(playerInfoProcessor playerInfoProcessor) *PlayerInfoUpsertConsumer {
	return &PlayerInfoUpsertConsumer{
		playerInfoProcessor: playerInfoProcessor,
	}
}
