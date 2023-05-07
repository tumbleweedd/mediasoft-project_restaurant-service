package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
)

type Menu interface {
	CreateMenu(menu models.Menu, salads, garnishes, meats, soups, drinks, desserts []string) error
	GetMenu()
}

type Order interface {
	GetUpToDateOrderList()
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
