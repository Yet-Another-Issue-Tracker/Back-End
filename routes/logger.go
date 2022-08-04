package routes

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"handler":   name,
			"timestamp": time.Since(start),
		}).Debug()
	})
}
