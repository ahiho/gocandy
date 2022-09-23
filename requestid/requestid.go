package requestid

import (
	"net/http"

	"github.com/google/uuid"
)

func InjectRequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Grpc-Metadata-x-request-id", uuid.New().String())
		handler.ServeHTTP(w, r)
	})
}
