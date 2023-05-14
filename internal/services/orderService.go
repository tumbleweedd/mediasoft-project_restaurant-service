package services

import (
	"context"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/rabbitmq"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
)

type OrderService struct {
	orderRepo repositories.Order
	rabbitmq  *rabbitmq.RabbitMQConn
}

func NewOrderService(orderRepo repositories.Order, rabbitmq *rabbitmq.RabbitMQConn) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		rabbitmq:  rabbitmq,
	}
}

func (os *OrderService) GetUpToDateOrderList(ctx context.Context, request *restaurant.GetUpToDateOrderListRequest) (*restaurant.GetUpToDateOrderListResponse, error) {
	return nil, nil
}
