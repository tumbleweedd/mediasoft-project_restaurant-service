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

type ProductService struct {
	productRepo repositories.Product
}

func NewProductService(productRepo repositories.Product) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, request *restaurant.CreateProductRequest) (*restaurant.CreateProductResponse, error) {
	officeUUID := uuid.New()
	product := models.Product{
		ProductUUID: officeUUID,
		Name:        request.Name,
		Description: request.Description,
		Type:        models.ProductType(request.Type.String()),
		Weight:      request.Weight,
		Price:       request.Price,
	}

	err := ps.productRepo.CreateProduct(product)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &restaurant.CreateProductResponse{}, nil
}

func (ps *ProductService) GetProductList(ctx context.Context, request *restaurant.GetProductListRequest) (*restaurant.GetProductListResponse, error) {
	products, err := ps.productRepo.GetProductList()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*restaurant.Product, 0, len(products))

	for _, product := range products {

		data = append(data, &restaurant.Product{
			Uuid:        product.ProductUUID.String(),
			Name:        product.Name,
			Description: product.Description,
			Type:        restaurant.ProductType(restaurant.ProductType_value[product.Type.String()]),
			Weight:      product.Weight,
			Price:       product.Price,
			CreatedAt:   timestamppb.New(product.CreatedAt),
		})
	}

	return &restaurant.GetProductListResponse{Result: data}, nil
}
