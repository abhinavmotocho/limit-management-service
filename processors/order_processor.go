package processors

import (
	"context"
	"encoding/json"
	"fmt"
	kafkahelper "kafka-helper-go"
	"kafka-helper-go/requests"
	khgutils "kafka-helper-go/utils"
	"limit-management-service/utils"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

//ListenToNewOrdersFromSRS persists new orders to database and forwards to oes service topic
func ListenToNewOrdersFromSRS() {
	env := viper.GetString(utils.ENVIRONMENT)
	ordersTopic := viper.GetString(fmt.Sprintf("%s.%s", env, utils.SRSToLMSTopic))
	consumerGrp := viper.GetString(fmt.Sprintf("%s.%s", env, utils.LMSConsumerGrp))

	processNewOrdersFromKafka(ordersTopic, consumerGrp)
}

//ProcessNewOrdersFromKafka listens to a topic
func processNewOrdersFromKafka(consumerTopic string, consumerGroup string) {
	data := make(chan kafka.Message, 1)
	wait := make(chan bool, 1)
	c, err := kafkahelper.SubscribeTopic(context.Background(), consumerTopic, consumerGroup, data, wait)
	if err != nil {
		log.WithError(err).Error("Unable to listen to srs topic")
		return
	}

	for {
		select {
		case v := <-data:
			reqID := uuid.New().String()
			ctx := context.WithValue(context.Background(), khgutils.REQUESTID, reqID)
			logger := log.WithContext(ctx)
			var t requests.OrderRequest
			err := json.Unmarshal(v.Value, &t)
			if err != nil {
				logger.WithError(err).Error("Unable to decode kafka message")
			} else {
				logger.WithFields(log.Fields{
					"OrderID": t.OrderID,
				}).Debugln("Received SRS Order")

				if err := PushToDatastore(t); err != nil {
					logger.WithError(err).Error("Error pushing limit order to datastore: " + t.OrderID)
				}
			}
			if ktp, err := c.Consumer.Commit(); err != nil {
				logger.WithError(err).Error("Error in commiting offsets for SRS Orders")
			} else {
				logger.Debugln(ktp)
			}
			wait <- true
		}
	}
}

//PushLimitOrderToRoute pushes placed orders to kafka for status check
func PushLimitOrderToRoute(req requests.OrderRequest) error {
	env := viper.GetString(utils.ENVIRONMENT)
	srsTopic := viper.GetString(fmt.Sprintf("%s.%s", env, utils.LMSToSRSTopic))

	p, err := kafkahelper.GetProducer()
	if err != nil {
		log.WithError(err).Error("Unable to create producer for srs topic")
		return err
	}

	serializedMsg, _ := json.Marshal(req)
	p.Producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &srsTopic, Partition: kafka.PartitionAny},
		Value: serializedMsg}
	p.Producer.Flush(100)

	return nil
}
