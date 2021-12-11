package kafka

import (
	"context"

	messagebroker "example.com/services/payment/internal/message_broker"
)

type Broker struct {
	brokers  []string
	clientID string

	paymentBroker messagebroker.PaymentBroker
}

func NewBroker(brokers []string, clientID string) messagebroker.MessageBroker {
	return &Broker{
		brokers:  brokers,
		clientID: clientID,
	}
}

func (b *Broker) Connect(ctx context.Context) error {
	brokers := []messagebroker.BrokerWithClient{b.Payment()}
	for _, broker := range brokers {
		if err := broker.Connect(ctx, b.brokers); err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) Close() error {
	brokers := []messagebroker.BrokerWithClient{b.Payment()}
	for _, broker := range brokers {
		if err := broker.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Broker) Payment() messagebroker.PaymentBroker {
	if b.paymentBroker == nil {
		b.paymentBroker = NewPaymentBroker(b.clientID)
	}

	return b.paymentBroker
}
