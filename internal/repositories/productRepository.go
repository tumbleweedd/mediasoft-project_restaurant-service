package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) CreateProduct(product models.Product) error {
	const query = `insert into restaurant.product (uuid, name, description, type, weight, price, created_at) 
					values ($1, $2, $3, $4, $5, $6, current_timestamp)`

	_, err := pr.db.Exec(query,
		product.ProductUUID, product.Name, product.Description, product.Type, product.Weight, product.Price,
	)

	return err
}

func (pr *ProductRepository) GetProductList() ([]*models.Product, error) {
	const query = `select uuid, name, description, type, weight, price, created_at from restaurant.product`

	var products []*models.Product

	err := pr.db.Select(&products, query)

	return products, err
}
