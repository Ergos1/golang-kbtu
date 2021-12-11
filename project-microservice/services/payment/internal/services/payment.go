package services

import (
	messagebroker "example.com/services/payment/internal/message_broker"
	"example.com/services/payment/internal/store"
)

type PaymentService struct {
	store  store.Store
	broker messagebroker.MessageBroker
}

func NewPaymentService(store store.Store, broker messagebroker.MessageBroker) *PaymentService {
	return &PaymentService{
		store:  store,
		broker: broker,
	}
}

/*
	1.Check wallets money
	2.if good
		change it
		then return nil
	3.if bad
		return error
*/
func (ps *PaymentService) Pay(from, to string, amount float64) error {
	return ps.broker.Payment().Pay(from, to, amount)
}
