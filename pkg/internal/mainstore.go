package internal

type mainstore interface {
	Welcome()
	RegisterHandler()
	LoginHandler()
	TransactionHistoryHandler()
	ExchangeRatesHandler()
}
