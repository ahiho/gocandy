package grpcmdinjector

import (
	"net/http"

	"github.com/ahiho/gocandy/utils"
)

var (
	bypassHeaders []string
)

const (
	MetadataHeaderPrefix = "Grpc-Metadata-"
)

func AllowHeaders(headers ...string) {
	for _, h := range headers {
		hkey := http.CanonicalHeaderKey(h)
		if !utils.IsInStringArr(bypassHeaders, hkey) {
			bypassHeaders = append(bypassHeaders, hkey)
		}
	}
}

func WrapGrpcMD(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v string
		for _, h := range bypassHeaders {
			v = r.Header.Get(h)
			if v != "" {
				r.Header.Set(MetadataHeaderPrefix+h, v)
			}
		}
		handler.ServeHTTP(w, r)
	})
}
