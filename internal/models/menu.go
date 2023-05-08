package models

import (
	"github.com/google/uuid"
	"time"
)

type Menu struct {
	MenuUuid        uuid.UUID `json:"menu_uuid" db:"uuid"`
	OnDate          time.Time `json:"on_date" db:"on_date"`
	OpeningRecordAt time.Time `json:"opening_record_at" db:"opening_record_at"`
	ClosingRecordAt time.Time `json:"closing_record_at" db:"closing_record_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type MenuProduct struct {
	Id          int       `json:"-"`
	MenuUUID    uuid.UUID `json:"menu_uuid" db:"menu_uuid"`
	ProductUUID uuid.UUID `json:"product_uuid" db:"product_uuid"`
}
