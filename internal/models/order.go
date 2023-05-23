package models

import "github.com/google/uuid"

type Order struct {
	OrderUUID  uuid.UUID `json:"order_uuid" db:"uuid"`
	UserUUID   uuid.UUID `json:"user_uuid" db:"user_uuid"`
	OfficeUUID uuid.UUID `json:"office_uuid" db:"office_uuid"`
}

type OrderItem struct {
	OrderItemUUID uuid.UUID `json:"uuid" db:"uuid"`
	Count         int       `json:"count" db:"count"`
	ProductUUID   uuid.UUID `json:"product_uuid" db:"product_uuid"`
	OrderUUID     uuid.UUID `json:"order_uuid" db:"order_uuid"`
}

// --- For restaurant response

type RestaurantOrderItem struct {
	ProductUUID uuid.UUID `db:"product_uuid"`
	ProductName string    `db:"product_name"`
	Count       int       `db:"count"`
}

type OrderByOffice struct {
	OfficeUUID    uuid.UUID `db:"office_uuid"`
	OfficeName    string    `db:"office_name"`
	OfficeAddress string    `db:"office_address"`
	Order         []*RestaurantOrderItem
}

type ResponseBody struct {
	TotalOrders         []*RestaurantOrderItem
	TotalOrdersByOffice []*OrderByOffice
}

type OrdersByCompanyRows struct {
	OfficeUUID    uuid.UUID `db:"office_uuid"`
	OfficeName    string    `db:"office_name"`
	OfficeAddress string    `db:"office_address"`
	ProductUUID   uuid.UUID `db:"product_uuid"`
	ProductName   string    `db:"product_name"`
	Count         int       `db:"count"`
}
