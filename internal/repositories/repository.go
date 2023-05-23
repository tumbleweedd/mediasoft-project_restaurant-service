package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/broker/kafka/kafkaModels"
	"time"
)

type Menu interface {
	CreateMenu(menu models.Menu, salads, garnishes, meats, soups, drinks, desserts []string) error
	GetMenu(onDateMenu time.Time) (*models.Menu, []*models.Product, error)
}

type Order interface {
	GetUpToDateOrderList() ([]*models.RestaurantOrderItem, []*models.OrdersByCompanyRows, error)
	CreateOrder(orderUUID uuid.UUID, order kafkaModels.OrderByOffice) error
	GetDataForStatisticFromOrder(orderUUID uuid.UUID, dataChanel chan<- models.ProductsFromOrdersResponse) error
}

type Product interface {
	CreateProduct(product models.Product) error
	GetProductList() ([]*models.Product, error)
}

type Repository struct {
	Menu
	Order
	Product
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Menu:    NewMenuRepository(db),
		Order:   NewOrderRepository(db),
		Product: NewProductRepository(db),
	}
}
