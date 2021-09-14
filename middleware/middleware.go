package middleware

import (
	"net/http"

	"github.com/bgoldovsky/casher/app/logger"
)

// Logging Логирует входящие запросы
func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.WithField("method", r.Method).WithField("path", r.URL.Path).Info("request")
		f(w, r)
	}
}
