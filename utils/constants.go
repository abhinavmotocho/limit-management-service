package utils

type key int

const (
	//ENVIRONMENT access key
	ENVIRONMENT = "ENVIRONMENT"

	//SRSToLMSTopic kafka topic to receive orders sent by smart routing service
	SRSToLMSTopic = "kafka.SRS_ORDERS_LMS"

	//LMSConsumerGrp consumer group for new orders
	LMSConsumerGrp = "kafka.LMS_CONSUMER_GRP"

	//LMSToSRSTopic kafka topic to send orders to smart routing service
	LMSToSRSTopic = "kafka.LMS_ORDERS_SRS"

	//PriceTopic kafka topic to receive live coin price data
	PriceTopic = "kafka.TRADE_BIDASK_TOPIC"
)
