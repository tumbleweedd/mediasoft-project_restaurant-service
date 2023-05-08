package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MenuService struct {
	menuRepo repositories.Menu
}

func NewMenuService(menuRepo repositories.Menu) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

func (ms *MenuService) CreateMenu(ctx context.Context, request *restaurant.CreateMenuRequest) (*restaurant.CreateMenuResponse, error) {
	menu := models.Menu{
		MenuUuid:        uuid.New(),
		OnDate:          request.OnDate.AsTime(),
		OpeningRecordAt: request.OpeningRecordAt.AsTime(),
		ClosingRecordAt: request.ClosingRecordAt.AsTime(),
	}

	err := ms.menuRepo.CreateMenu(
		menu, request.Salads, request.Garnishes, request.Meats, request.Soups, request.Drinks, request.Desserts,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &restaurant.CreateMenuResponse{}, nil
}

func (ms *MenuService) GetMenu(ctx context.Context, request *restaurant.GetMenuRequest) (*restaurant.GetMenuResponse, error) {
	menu, products, err := ms.menuRepo.GetMenu(request.OnDate.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := &restaurant.Menu{
		Uuid:            menu.MenuUuid.String(),
		OnDate:          timestamppb.New(menu.OnDate),
		OpeningRecordAt: timestamppb.New(menu.OpeningRecordAt),
		ClosingRecordAt: timestamppb.New(menu.ClosingRecordAt),
		CreatedAt:       timestamppb.New(menu.CreatedAt),
	}

	for _, product := range products {
		addProductToMenu(data, product)
	}

	return &restaurant.GetMenuResponse{Menu: data}, nil
}

func addProductToMenu(menu *restaurant.Menu, product *models.Product) {
	p := &restaurant.Product{
		Uuid:        product.ProductUUID.String(),
		Name:        product.Name,
		Description: product.Description,
		Type:        restaurant.ProductType(restaurant.ProductType_value[product.Type.String()]),
		Weight:      product.Weight,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
	}

	switch p.Type {
	case restaurant.ProductType_PRODUCT_TYPE_SALAD:
		menu.Salads = append(menu.Salads, p)
	case restaurant.ProductType_PRODUCT_TYPE_GARNISH:
		menu.Garnishes = append(menu.Garnishes, p)
	case restaurant.ProductType_PRODUCT_TYPE_MEAT:
		menu.Meats = append(menu.Meats, p)
	case restaurant.ProductType_PRODUCT_TYPE_SOUP:
		menu.Soups = append(menu.Soups, p)
	case restaurant.ProductType_PRODUCT_TYPE_DRINK:
		menu.Drinks = append(menu.Drinks, p)
	case restaurant.ProductType_PRODUCT_TYPE_DESSERT:
		menu.Desserts = append(menu.Desserts, p)
	}
}
