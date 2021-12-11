package kafka

import (
	"context"
	"encoding/json"
	"log"

	"example.com/services/payment/internal/message_broker"
	"example.com/services/payment/internal/models"
	"github.com/Shopify/sarama"
)

const (
	paymentTopic = "cache"
)

type (
	PaymentBroker struct {
		syncProducer  sarama.SyncProducer
		consumerGroup sarama.ConsumerGroup

		consumerHandler *paymentConsumeHandler
		clientID        string
	}

	paymentConsumeHandler struct {
		pb *PaymentBroker
		ctx   context.Context
		ready chan bool
	}
)

func NewPaymentBroker(clientID string) messagebroker.PaymentBroker {
	return &PaymentBroker{
		clientID: clientID,
		consumerHandler: &paymentConsumeHandler{
			ready: make(chan bool),
		},
	}
}

func (c *PaymentBroker) Connect(ctx context.Context, brokers []string) error {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 10
	producerConfig.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return err
	}
	c.syncProducer = syncProducer

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(brokers, c.clientID, consumerConfig)
	if err != nil {
		return err
	}
	c.consumerGroup = consumerGroup

	c.consumerHandler.ctx = ctx

	go func() {
		for {
			if err := c.consumerGroup.Consume(ctx, []string{paymentTopic}, c.consumerHandler); err != nil {
				log.Panicf("[error] from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			c.consumerHandler.ready = make(chan bool)
		}
	}()

	<-c.consumerHandler.ready

	return nil
}

func (c *PaymentBroker) Close() error {
	if err := c.syncProducer.Close(); err != nil {
		return err
	}

	if err := c.consumerGroup.Close(); err != nil {
		return err
	}

	return nil
}


func (c *PaymentBroker) Check(from, to string, amount float64) error {
	
	return nil
}

func (c *PaymentBroker) Pay(from, to string, amount float64) error {
	
	return nil
}

func (c *paymentConsumeHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *paymentConsumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *paymentConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		msgResp := new(models.MessageResponse)
		if err := json.Unmarshal(msg.Value, msgResp); err != nil {
			return err
		}

		switch msgResp.Response {
		case models.Good:
			// c.pb.Pay()
			//good
		case models.Bad:
			//bad
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
