package services

import (
	"context"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
)

type OrderService struct {
	orderRepo repositories.Order
}

func NewOrderService(orderRepo repositories.Order) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (os *OrderService) GetUpToDateOrderList(ctx context.Context, request *restaurant.GetUpToDateOrderListRequest) (*restaurant.GetUpToDateOrderListResponse, error) {
	//TODO implement me
	panic("implement me")
}
