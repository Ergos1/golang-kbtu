package messagebroker

type PaymentBroker interface {
	BrokerWithClient
	Remove(key interface{}) error
	Purge() error
}
