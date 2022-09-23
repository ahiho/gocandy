package deviceid

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

func DeviceID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get("x-device-id")
}
