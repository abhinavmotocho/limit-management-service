package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getHealth(rw http.ResponseWriter, req *http.Request) {
	logger := log.WithContext(req.Context())
	logger.Debug("Received health check request ..Returning 200 status")
	rw.WriteHeader(http.StatusOK)
}
