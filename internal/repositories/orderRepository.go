package repositories

import "github.com/jmoiron/sqlx"

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) GetUpToDateOrderList() {
	//TODO implement me
	panic("implement me")
}
