package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

func PanicHandler(next http.Handler, internalServerErrorHandler func(http.ResponseWriter, *http.Request, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				internalServerErrorHandler(w, r, errors.New("internal server error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func LogRequestHandler(next http.Handler, logger logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI()))
		next.ServeHTTP(w, r)
	})
}
