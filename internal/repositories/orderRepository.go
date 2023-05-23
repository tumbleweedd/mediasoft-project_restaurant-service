package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/broker/kafka/kafkaModels"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) GetUpToDateOrderList() ([]*models.RestaurantOrderItem, []*models.OrdersByCompanyRows, error) {
	const getTotalOrdersQuery = `select oi.product_uuid as product_uuid, p.name as product_name, sum(oi.count) as count
									from restaurant.order_items oi
											 join restaurant.product p on oi.product_uuid = p.uuid
									group by oi.product_uuid, p.name`

	var orderItems []*models.RestaurantOrderItem
	err := or.db.Select(&orderItems, getTotalOrdersQuery)
	if err != nil {
		return nil, nil, err
	}

	const ordersByOfficeQuery = `select 
    							   ood.uuid        as office_uuid,
								   ood.name        as office_name,
								   ood.address     as office_address,
								   oi.product_uuid as product_uuid,
								   p.name          as product_name,
								   sum(oi.count)        as count
							from restaurant.orders o
									 join restaurant.offices_of_delivery ood on ood.uuid = o.office_uuid
									 join restaurant.order_items oi on o.uuid = oi.order_uuid
									 join restaurant.product p on oi.product_uuid = p.uuid
							group by ood.uuid, ood.name, ood.address, oi.product_uuid, p.name`
	var ordersByCompanyRows []*models.OrdersByCompanyRows
	if err := or.db.Select(&ordersByCompanyRows, ordersByOfficeQuery); err != nil {
		return nil, nil, err
	}

	return orderItems, ordersByCompanyRows, nil
}

func (or *OrderRepository) CreateOrder(orderUUID uuid.UUID, order kafkaModels.OrderByOffice) error {
	tx, err := or.db.Begin()
	if err != nil {
		return err
	}

	const createOfficeByDeliveryQuery = `insert into restaurant.offices_of_delivery (uuid, name, address) 
									VALUES ($1, $2, $3) on conflict do nothing `
	_, err = tx.Exec(createOfficeByDeliveryQuery, order.OfficeUUID, order.OfficeName, order.OfficeAddress)
	if err != nil {
		tx.Rollback()
		return err
	}

	const createOrderQuery = `insert into restaurant.orders (uuid, user_uuid, office_uuid)
							VALUES ($1, $2, $3)`
	_, err = tx.Exec(createOrderQuery, orderUUID, order.UserUUID, order.OfficeUUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	const createOrderItemQuery = `insert into restaurant.order_items (count, product_uuid, order_uuid) 
									values ($1, $2, $3)`
	for _, items := range [][]*kafkaModels.OrderItemByOffice{order.Salads, order.Garnishes, order.Meats, order.Soups, order.Drinks, order.Desserts} {
		for _, item := range items {
			_, err := tx.Exec(createOrderItemQuery, item.Count, item.ProductUUID, orderUUID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (or *OrderRepository) GetDataForStatisticFromOrder(orderUUID uuid.UUID, dataChanel chan<- models.ProductsFromOrdersResponse) error {
	const query = `select oi.product_uuid, oi.count, p.name, p.price, p.type, o.created_at
						from restaurant.order_items oi
								 join restaurant.product p on p.uuid = oi.product_uuid
								 join restaurant.orders o on o.uuid = oi.order_uuid
						where oi.order_uuid=$1`

	rows, err := or.db.Query(query, orderUUID)
	if err != nil {
		return err
	}

	logger.Info("Сделал запрос на получение статистики")

	for rows.Next() {
		var productStat models.ProductsFromOrders

		err := rows.Scan(
			&productStat.ProductUUID,
			&productStat.Count,
			&productStat.ProductName,
			&productStat.Price,
			&productStat.ProductType,
			&productStat.CreatedAt,
		)

		productStatResponse := formatResponseToChannel(&productStat)

		logger.Info("Прочитал стату")
		if err != nil {
			return err
		}
		logger.Info("Стою перед каналом (для отправки в него)")
		dataChanel <- productStatResponse
		logger.Info("Отправил в канал стату")
	}

	return nil
}

func formatResponseToChannel(productStat *models.ProductsFromOrders) models.ProductsFromOrdersResponse {
	var productStatResponse models.ProductsFromOrdersResponse

	createdAtFormat := productStat.CreatedAt.Format("2006.01.02 15:04:05")

	productStatResponse.ProductUUID = productStat.ProductUUID
	productStatResponse.Count = productStat.Count
	productStatResponse.ProductName = productStat.ProductName
	productStatResponse.Price = productStat.Price
	productStatResponse.ProductType = productStat.ProductType
	productStatResponse.CreatedAt = createdAtFormat

	return productStatResponse
}
