package models

import (
	"github.com/google/uuid"
	"time"
)

type ProductType string

const (
	UNSPECIFIED ProductType = "PRODUCT_TYPE_UNSPECIFIED"
	SALAD                   = "PRODUCT_TYPE_SALAD"
	GARNISH                 = "PRODUCT_TYPE_GARNISH"
	MEAT                    = "PRODUCT_TYPE_MEAT"
	SOUP                    = "PRODUCT_TYPE_SOUP"
	DRINK                   = "PRODUCT_TYPE_DRINK"
	DESSERT                 = "PRODUCT_TYPE_DESSERT"
)

func (pt ProductType) String() string {
	return string(pt)
}

type Product struct {
	ProductUUID uuid.UUID   `json:"product_uuid" db:"uuid"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Type        ProductType `json:"type" db:"type"`
	Weight      int32       `json:"weight" db:"weight"`
	Price       float64     `json:"price" db:"price"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}
