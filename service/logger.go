package service

import (
	"context"
	"kafka-helper-go/utils"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Logger will log the request
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := getRequestID(r.Header)
		ctx := context.WithValue(r.Context(), utils.REQUESTID, reqID)

		log.WithFields(log.Fields{
			"request-id": reqID,
			"Method":     r.Method,
			"RequestUri": r.RequestURI,
			"Name":       name,
			"StartTime":  start,
		})

		inner.ServeHTTP(w, r.WithContext(ctx))

		log.WithFields(log.Fields{
			"request-id":     reqID,
			"TimeSinceStart": time.Since(start),
		})
	})
}

func getRequestID(headers http.Header) string {
	reqID := headers.Get("x-request-id")
	if len(strings.TrimSpace(reqID)) == 0 {
		reqID = uuid.New().String()
	}
	return reqID
}
