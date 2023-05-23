package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/broker/kafka/kafkaModels"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
	"log"
	"os"
	"os/signal"
	"sync"
)

type Message struct {
	Order kafkaModels.OrderByOffice `json:"Order"`
}

type Consumer struct {
	consumer  sarama.Consumer
	orderRepo repositories.Order
}

func NewConsumer(brokers []string, orderRepo repositories.Order) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  consumer,
		orderRepo: orderRepo,
	}, nil
}

func (c *Consumer) Consume(topic string, dataChan chan<- models.ProductsFromOrdersResponse) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions for topic %s: %s", topic, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Error creating consumer for partition %d: %s", partition, err)
				return
			}

			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					log.Printf("Error closing partition consumer for partition %d: %s", partition, err)
				}
			}()

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt)

			for {
				select {
				case msg := <-partitionConsumer.Messages():
					response, err := decodeMessage(msg.Value)
					if err != nil {
						log.Printf("Error decoding JSON: %s", err)
						continue
					}
					orderUUID := uuid.New()
					err = c.orderRepo.CreateOrder(orderUUID, response.Order)
					if err != nil {
						logger.Infof("Error creating order with kafka message: %s", err)
					}
					logger.Info("Создал заказ")
					go func() {
						err = c.orderRepo.GetDataForStatisticFromOrder(orderUUID, dataChan)
						if err != nil {
							logger.Infof("Error get statistic data: %s", err)
						}
					}()
					logger.Info("Запустил метод получения статы")
				case err := <-partitionConsumer.Errors():
					logger.Infof("Error consuming partition %d: %s", partition, err)
				case <-signals:
					return
				}
			}
		}(partition)
	}
	wg.Wait()

	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func decodeMessage(data []byte) (*Message, error) {
	var orders Message
	err := json.Unmarshal(data, &orders)
	if err != nil {
		return nil, err
	}

	return &orders, nil
}
