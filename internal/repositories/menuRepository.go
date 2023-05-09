package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"time"
)

type MenuRepository struct {
	db *sqlx.DB
}

func NewMenuRepository(db *sqlx.DB) *MenuRepository {
	return &MenuRepository{
		db: db,
	}
}

func (mr *MenuRepository) CreateMenu(menu models.Menu, salads, garnishes, meats, soups, drinks, desserts []string) error {
	tx, err := mr.db.Beginx()
	if err != nil {
		return err
	}

	var menuUUID uuid.UUID
	const createMenuQuery = `insert into restaurant.menu (uuid, on_date, opening_record_at, closing_record_at) 
								VALUES ($1, $2, $3, $4) returning uuid`
	row := tx.QueryRow(createMenuQuery, menu.MenuUuid, menu.OnDate, menu.OpeningRecordAt, menu.ClosingRecordAt)
	if err = row.Scan(&menuUUID); err != nil {
		tx.Rollback()
		return err
	}

	const createMenuProductQuery = `insert into restaurant.menu_product(menu_uuid, order_uuid) 
										VALUES ($1, $2) `
	for _, products := range [][]string{salads, garnishes, meats, soups, drinks, desserts} {
		for _, productUUID := range products {
			_, err = tx.Exec(createMenuProductQuery, menuUUID, productUUID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (mr *MenuRepository) GetMenu(menuOnDate time.Time) (*models.Menu, []*models.Product, error) {
	const getMenuQuery = `select uuid, on_date, opening_record_at, closing_record_at, created_at 
						  from restaurant.menu where date(on_date)=date($1)`

	var menu models.Menu
	var products []*models.Product

	err := mr.db.Get(&menu, getMenuQuery, menuOnDate)
	if err != nil {
		return nil, nil, err
	}

	const getProductsInMenuQuery = `select p.uuid, p.name, p.description, p.type, p.weight, p.price, p.created_at
    								from restaurant.product p
										join restaurant.menu_product mp on p.uuid = mp.order_uuid
									where mp.menu_uuid=$1`
	err = mr.db.Select(&products, getProductsInMenuQuery, menu.MenuUuid.String())

	return &menu, products, err
}
