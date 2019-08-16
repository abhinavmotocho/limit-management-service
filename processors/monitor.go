package processors

import (
	"kafka-helper-go/requests"
)

var priceMap map[string]float64

func init() {
	priceMap = make(map[string]float64)
}

//PushToDatastore stores a limit order in the datastore.
func PushToDatastore(req requests.OrderRequest) error {
	//TODO put order in datastore and assign LMS instance to track the order
	return nil
}

//UpdatePrice updates the current market price for each currency
func UpdatePrice(m map[string]float64) error {
	for k, v := range m {
		priceMap[k] = v
		checkOrders(k)
	}
	return nil
}

func checkOrders(coin string) {
	//TODO access datastore and check if there are limit orders near current market price.
	//if order is ready to execute, call PushLimitOrderToRoute(order) where order is of type requests.OrderRequest
}
