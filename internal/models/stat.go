package models

import "time"

type ProductsFromOrders struct {
	Count       uint32    `json:"count"`
	ProductUUID string    `json:"product_uuid"`
	Price       float64   `json:"price"`
	ProductName string    `json:"product_name"`
	ProductType string    `json:"product_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductsFromOrdersResponse struct {
	Count       uint32  `json:"count"`
	ProductUUID string  `json:"product_uuid"`
	Price       float64 `json:"price"`
	ProductName string  `json:"product_name"`
	ProductType string  `json:"product_type"`
	CreatedAt   string  `json:"created_at"`
}
