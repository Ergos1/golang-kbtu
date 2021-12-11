package messagebroker

type PaymentBroker interface {
	BrokerWithClient
	Check(from, to string, amount float64) error
	Pay(from, to string, amount float64) error
}
