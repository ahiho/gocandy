package auth

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ahiho/gocandy/utils"
)

var apiKeys []string

func AddAPIKey(key string) {
	if key == "" {
		return
	}
	if !utils.IsInStringArr(apiKeys, key) {
		apiKeys = append(apiKeys, key)
	}
}

func VerifyAPIKey(ctx context.Context) (context.Context, error) {
	apiKey := metautils.ExtractIncoming(ctx).Get("x-api-key")
	if apiKey == "" || !utils.IsInStringArr(apiKeys, apiKey) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid api key")
	}
	return ctx, nil
}
