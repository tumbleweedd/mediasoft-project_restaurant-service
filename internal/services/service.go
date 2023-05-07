package services

import (
	"context"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
)

type Menu interface {
	CreateMenu(context.Context, *restaurant.CreateMenuRequest) (*restaurant.CreateMenuResponse, error)
	GetMenu(context.Context, *restaurant.GetMenuRequest) (*restaurant.GetMenuResponse, error)
}

type Order interface {
	GetUpToDateOrderList(context.Context, *restaurant.GetUpToDateOrderListRequest) (*restaurant.GetUpToDateOrderListResponse, error)
}

type Product interface {
	CreateProduct(context.Context, *restaurant.CreateProductRequest) (*restaurant.CreateProductResponse, error)
	GetProductList(context.Context, *restaurant.GetProductListRequest) (*restaurant.GetProductListResponse, error)
}

type Service struct {
	Menu
	Order
	Product
	restaurant.UnsafeMenuServiceServer
	restaurant.UnsafeOrderServiceServer
	restaurant.UnsafeProductServiceServer
}

func NewService(r *repositories.Repository) *Service {
	return &Service{
		Menu:    NewMenuService(r.Menu),
		Order:   NewOrderService(r.Order),
		Product: NewProductService(r.Product),
	}
}
