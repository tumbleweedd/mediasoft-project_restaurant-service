package models

import "github.com/google/uuid"

type Ordr struct {
	OrderUUID uuid.UUID `json:"order_uuid" db:"uuid"`
	UserUUID  uuid.UUID `json:"user_uuid" db:"user_uuid"`
}

type OrdrItm struct {
	ID          int       `json:"id" db:"id"`
	Count       int       `json:"count" db:"count"`
	ProductUUID uuid.UUID `json:"product_uuid" db:"product_uuid"`
	OrderUUID   uuid.UUID `json:"order_uuid" db:"order_uuid"`
}
