package kafkaModels

import "github.com/google/uuid"

type OrderItemByOffice struct {
	Count       int       `json:"Count"`
	ProductUUID uuid.UUID `json:"ProductUUID"`
}

type OrderByOffice struct {
	UserUUID      uuid.UUID            `json:"UserUUID"`
	OfficeUUID    uuid.UUID            `json:"Office_uuid"`
	OfficeName    string               `json:"Office_name"`
	OfficeAddress string               `json:"Office_address"`
	Salads        []*OrderItemByOffice `json:"Salads,omitempty"`
	Garnishes     []*OrderItemByOffice `json:"Garnishes,omitempty"`
	Meats         []*OrderItemByOffice `json:"Meats,omitempty"`
	Soups         []*OrderItemByOffice `json:"Soups,omitempty"`
	Drinks        []*OrderItemByOffice `json:"Drinks,omitempty"`
	Desserts      []*OrderItemByOffice `json:"Desserts,omitempty"`
}
