package app

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/services"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/broker/kafka/consumer"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/broker/kafka/producer"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/database/postgres"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	orderTopic    = "orders"
	sendStatTopic = "prodStat"
	broker        = "192.168.0.109:9092"
	bufferSize    = 6
)

func Run() {
	if err := godotenv.Load(); err != nil {
		logger.Errorf("Error getting env, %v", err)
	}
	db, err := postgres.NewPostgresDB(&postgres.Config{
		PgHost:         os.Getenv("DB_HOST"),
		PgPort:         os.Getenv("DB_PORT"),
		PgUser:         os.Getenv("DB_USER"),
		PgPwd:          os.Getenv("DB_PASSWORD"),
		PgDBName:       os.Getenv("DB_NAME"),
		PgDBSchemaName: os.Getenv("DB_SCHEMA_NAME"),
		PgSSLMode:      os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logger.Errorf("failed to initialize db: %s", err.Error())
		return
	}

	lis, err := net.Listen("tcp", os.Getenv("PORT"))
	if err != nil {
		logger.Errorf("Failed to listing: %v", err)
		return
	}

	logger.Infof("Auth on %s", os.Getenv("PORT"))

	s := grpc.NewServer()
	repo := repositories.NewRepository(db)

	dataChannel := make(chan models.ProductsFromOrdersResponse)
	done := make(chan struct{})
	defer close(done)
	defer close(dataChannel)

	kafkaConsumer, err := consumer.NewConsumer([]string{broker}, repo)
	if err != nil {
		log.Println("Failed to kafka conn: ", err)
	}
	defer kafkaConsumer.Close()

	kafkaProducer, err := producer.NewProducer(broker)
	if err != nil {
		log.Fatalln("Failed to kafka conn: ", err)
	}
	defer kafkaProducer.Close()

	consumeOrders(kafkaConsumer, orderTopic, dataChannel)
	produceStat(kafkaProducer, dataChannel, done, sendStatTopic, bufferSize)

	svc := services.NewService(repo)
	restaurant.RegisterMenuServiceServer(s, svc)
	restaurant.RegisterOrderServiceServer(s, svc)
	restaurant.RegisterProductServiceServer(s, svc)

	if err := s.Serve(lis); err != nil {
		logger.Errorf("Failed to serve: %v", err)
		return
	}
}

func consumeOrders(kafkaConsumer *consumer.Consumer, topic string, dataChannel chan<- models.ProductsFromOrdersResponse) {
	go func() {
		err := kafkaConsumer.Consume(topic, dataChannel)
		if err != nil {
			log.Printf("Error consuming Kafka topic: %s", err)
		}
	}()
	log.Printf("Started listening to Kafka topic: %s", topic)
}

func produceStat(kafkaProducer *producer.KafkaProducer, dataChannel <-chan models.ProductsFromOrdersResponse, done chan struct{}, topic string, bufferSize int) {
	go kafkaProducer.FormatBuffer(dataChannel, done, sendStatTopic, bufferSize)
}
