package admininfoupsertconsumer

type adminInfoProcessor interface { // TODO: Прописать интерфейс
}

type AdminInfoUpsertConsumer struct {
	adminInfoProcessor adminInfoProcessor
}

func NewAdminInfoUpsertConsumer(adminInfoProcessor adminInfoProcessor) *AdminInfoUpsertConsumer { // TODO: Добавить получение кафки
	return &AdminInfoUpsertConsumer{
		adminInfoProcessor: adminInfoProcessor,
	}
}
