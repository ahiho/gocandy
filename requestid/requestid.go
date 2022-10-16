package requestid

import (
	"net/http"
)

type RequestIDGenerator func() string

func InjectRequestID(handler http.Handler, generator RequestIDGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Grpc-Metadata-x-request-id", generator())
		handler.ServeHTTP(w, r)
	})
}
