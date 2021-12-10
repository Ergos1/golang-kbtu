package messagebroker

type CacheBroker interface {
	BrokerWithClient
	Remove(key interface{}) error
	Purge() error
}
