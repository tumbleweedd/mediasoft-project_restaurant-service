package services

import (
	"context"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/repositories"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	orderItems, ordersByCompanyRows, err := os.orderRepo.GetUpToDateOrderList()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &models.ResponseBody{
		TotalOrders: orderItems,
	}

	ordersByCompanyMap := make(map[string]*models.OrderByOffice)
	for _, row := range ordersByCompanyRows {
		officeUUID := row.OfficeUUID
		if _, ok := ordersByCompanyMap[officeUUID.String()]; !ok {
			ordersByCompanyMap[officeUUID.String()] = &models.OrderByOffice{
				OfficeUUID:    row.OfficeUUID,
				OfficeName:    row.OfficeName,
				OfficeAddress: row.OfficeAddress,
				Order:         []*models.RestaurantOrderItem{},
			}
		}
		totalOrder := &models.RestaurantOrderItem{
			ProductUUID: row.ProductUUID,
			ProductName: row.ProductName,
			Count:       row.Count,
		}
		ordersByCompanyMap[officeUUID.String()].Order = append(ordersByCompanyMap[officeUUID.String()].Order, totalOrder)
	}

	for _, value := range ordersByCompanyMap {
		response.TotalOrdersByOffice = append(response.TotalOrdersByOffice, value)
	}

	result := getResponse(response)

	return &restaurant.GetUpToDateOrderListResponse{
		TotalOrders:          result.TotalOrders,
		TotalOrdersByCompany: result.TotalOrdersByCompany,
	}, nil
}

func getResponse(response *models.ResponseBody) *restaurant.GetUpToDateOrderListResponse {
	totalOrders := make([]*restaurant.Order, 0, len(response.TotalOrders))
	totalOrdersByOffice := make([]*restaurant.OrdersByOffice, 0, len(response.TotalOrdersByOffice))

	for _, item := range response.TotalOrders {
		totalOrders = append(totalOrders, &restaurant.Order{
			ProductId:   item.ProductUUID.String(),
			ProductName: item.ProductName,
			Count:       int64(item.Count),
		})
	}

	for _, item := range response.TotalOrdersByOffice {
		ordersByOffice := &restaurant.OrdersByOffice{
			CompanyId:     item.OfficeUUID.String(),
			OfficeName:    item.OfficeName,
			OfficeAddress: item.OfficeAddress,
			Result:        make([]*restaurant.Order, len(item.Order)),
		}

		for i, orderItem := range item.Order {
			order := &restaurant.Order{
				ProductId:   orderItem.ProductUUID.String(),
				ProductName: orderItem.ProductName,
				Count:       int64(orderItem.Count),
			}
			ordersByOffice.Result[i] = order
		}

		totalOrdersByOffice = append(totalOrdersByOffice, ordersByOffice)
	}

	result := &restaurant.GetUpToDateOrderListResponse{
		TotalOrders:          totalOrders,
		TotalOrdersByCompany: totalOrdersByOffice,
	}
	return result
}
