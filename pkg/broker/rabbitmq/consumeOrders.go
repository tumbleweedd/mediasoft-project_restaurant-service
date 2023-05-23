package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/internal/models"
	"github.com/tumbleweedd/mediasoft-intership/restaraunt-service/pkg/logger"
)

func (r *RabbitMQConn) ConsumeOrders() error {
	msgs, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order models.Ordr
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				logger.Infof("Ошибка декодирования заказа из очереди: %s", err)
				d.Ack(false)
				continue
			}

			d.Ack(false)
			fmt.Println(order.OrderUUID)
		}
	}()

	logger.Infof(" [*] Ожидание заказов из очереди %s. Для выхода нажмите CTRL+C", r.queueName)
	<-forever

	return nil
}
