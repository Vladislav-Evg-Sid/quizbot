package playerinfoupsertconsumer

type playerInfoProcessor interface { // TODO: Прописать интерфейс
}

type PlayerInfoUpsertConsumer struct { // TODO: Прописать дополнительные переменные
	playerInfoProcessor playerInfoProcessor
}

func NewPlayerInfoUpsertConsumer(playerInfoProcessor playerInfoProcessor) *PlayerInfoUpsertConsumer { // TODO: Добавить получение кафки
	return &PlayerInfoUpsertConsumer{
		playerInfoProcessor: playerInfoProcessor,
	}
}
