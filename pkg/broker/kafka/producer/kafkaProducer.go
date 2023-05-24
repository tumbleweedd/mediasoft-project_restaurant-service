package producer

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
)

type KafkaProducer struct {
	p sarama.AsyncProducer
}

func NewProducer(broker string) (*KafkaProducer, error) {
	producer, err := sarama.NewAsyncProducer([]string{broker}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		p: producer,
	}, nil
}

func (p *KafkaProducer) StartProduce(order models.ProductsFromOrdersResponse, done chan struct{}, topic string) {
	msg := order

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Infof("Failed to marshal message, err : %s\n", err)
		return
	}

	select {
	case <-done:
		return
	case p.p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}:
		logger.Info("Product send do kafka")
	case err := <-p.p.Errors():
		logger.Infof("Failed to send message to Kafka, err: %s, msg: %s\n", err, msgBytes)
	}
}

func (p *KafkaProducer) FormatBuffer(dataChanel <-chan models.ProductsFromOrdersResponse, done chan struct{}, topic string) {
	logger.Info("Запустил FormatBuffer")
	for {
		select {
		case <-done:
			return
		case product, ok := <-dataChanel:
			if !ok {
				p.StartProduce(product, done, topic)
				fmt.Println(product.ProductUUID)

				return
			}
			p.StartProduce(product, done, topic)
		}
	}
}

func (p *KafkaProducer) Close() error {
	if p != nil {
		return p.p.Close()
	}
	return nil
}
