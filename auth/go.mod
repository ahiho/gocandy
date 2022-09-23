module github.com/ahiho/gocandy/auth

go 1.19

require (
	ahiho/gocandy/utils v0.0.0
	firebase.google.com/go/v4 v4.8.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	google.golang.org/api v0.97.0
	google.golang.org/grpc v1.49.0
)

replace ahiho/gocandy/utils v0.0.0 => ../utils

require (
	cloud.google.com/go v0.102.0 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/firestore v1.6.1 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/storage v1.22.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/appengine/v2 v2.0.1 // indirect
	google.golang.org/genproto v0.0.0-20220624142145-8cd45d7dbd1f // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
