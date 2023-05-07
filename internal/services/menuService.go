package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//TODO implement me
	panic("implement me")
}
