package main

import (
	"crypto/tls"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	kafkahelper "kafka-helper-go"
	khgutils "kafka-helper-go/utils"
	"limit-management-service/processors"
	"limit-management-service/service"
)

func main() {
	loadConfig()
	log.Info("Listening to new orders ...")
	kafkahelper.Init(viper.GetString(khgutils.KAFKABROKER))
	go processors.ListenToNewOrdersFromSRS()
	router := service.NewRouter()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	log.Fatal(http.ListenAndServe(":8082", router))
}

func loadConfig() {
	viper.SetConfigName("app")      // no need to include file extension
	viper.AddConfigPath("./config") // set the path of your config file

	err := viper.ReadInConfig()

	if err != nil {
		log.WithError(err).Fatalln("error opening config file")
	}

	viper.AutomaticEnv()
	khgutils.InitLog("OSS")
}
