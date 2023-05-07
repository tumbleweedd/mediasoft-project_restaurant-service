package app

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/services"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/database/postgres"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
	"net"
	"os"
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
	svc := services.NewService(repo)

	restaurant.RegisterMenuServiceServer(s, svc)
	restaurant.RegisterOrderServiceServer(s, svc)
	restaurant.RegisterProductServiceServer(s, svc)

	if err := s.Serve(lis); err != nil {
		logger.Errorf("Failed to serve: %v", err)
		return
	}

}
