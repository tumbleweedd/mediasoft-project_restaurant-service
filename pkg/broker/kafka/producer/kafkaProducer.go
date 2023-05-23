package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
)

type Message struct {
	Product []models.ProductsFromOrdersResponse `json:"result"`
}

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

func (p *KafkaProducer) StartProduce(done chan struct{}, topic string, buffer []models.ProductsFromOrdersResponse) {
	msg := Message{
		Product: buffer,
	}

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

func (p *KafkaProducer) FormatBuffer(dataChanel <-chan models.ProductsFromOrdersResponse, done chan struct{}, topic string, bufferSize int) {
	buffer := make([]models.ProductsFromOrdersResponse, 0, bufferSize)
	logger.Info("Запустил FormatBuffer")

	for {
		select {
		case <-done:
			return
		case product, ok := <-dataChanel:
			if !ok {
				if len(buffer) > 0 {
					p.StartProduce(done, topic, buffer)
				}
				return
			}

			buffer = append(buffer, product)
			if len(buffer) >= bufferSize {
				p.StartProduce(done, topic, buffer)
				buffer = buffer[:0]
			}
		}
	}
}

func (p *KafkaProducer) Close() error {
	if p != nil {
		return p.p.Close()
	}
	return nil
}
