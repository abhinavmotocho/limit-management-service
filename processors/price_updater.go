package processors

import (
	"context"
	"encoding/json"
	"fmt"
	kafkahelper "kafka-helper-go"
	khgutils "kafka-helper-go/utils"
	"limit-management-service/utils"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//ListenToNewOrdersFromSRS persists new orders to database and forwards to oes service topic
func ListenToPriceUpdates() {
	env := viper.GetString(utils.ENVIRONMENT)
	priceTopic := viper.GetString(fmt.Sprintf("%s.%s", env, utils.PriceTopic))
	consumerGrp := viper.GetString(fmt.Sprintf("%s.%s", env, utils.LMSConsumerGrp))

	processNewOrdersFromKafka(priceTopic, consumerGrp)
}

//ProcessNewOrdersFromKafka listens to a topic
func processPriceUpdatesFromKafka(consumerTopic string, consumerGroup string) {
	data := make(chan kafka.Message, 1)
	wait := make(chan bool, 1)
	c, err := kafkahelper.SubscribeTopic(context.Background(), consumerTopic, consumerGroup, data, wait)
	if err != nil {
		log.WithError(err).Error("Unable to listen to price topic")
		return
	}

	for {
		select {
		case v := <-data:
			reqID := uuid.New().String()
			ctx := context.WithValue(context.Background(), khgutils.REQUESTID, reqID)
			logger := log.WithContext(ctx)
			//CHANGE THE STRUCT
			var t map[string]float64
			err := json.Unmarshal(v.Value, &t)
			if err != nil {
				logger.WithError(err).Error("Unable to decode kafka message")
			} else {
				logger.Debugln("Received Price Update")
				if err := UpdatePrice(t); err != nil {
					logger.WithError(err).Error("Error in reading price update.")
				}
			}
			if ktp, err := c.Consumer.Commit(); err != nil {
				logger.WithError(err).Error("Error in commiting offsets for price updates")
			} else {
				logger.Debugln(ktp)
			}
			wait <- true
		}
	}
}
