package apperror

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errorNormalizedFlag      = "x-rpc-err-normalized"
	AllowedHTTPErrorStatuses = []int{400, 401, 403, 404}
)

func MessageForGrpcStatus(err *status.Status) string {
	code := err.Code()
	msg := err.Message()
	// nolint: exhaustive // We don't care other status.
	switch code {
	case codes.InvalidArgument:
		return fmt.Sprintf(`{"code":"ERR_BAD_REQUEST","message":"%v"}`, msg)
	case codes.Unauthenticated:
		return fmt.Sprintf(`{"code":"ERR_UNAUTHORIZED","message":"%v"}`, msg)
	case codes.PermissionDenied:
		return fmt.Sprintf(`{"code":"ERR_FORBIDDEN","message":"%v"}`, msg)
	case codes.NotFound:
		return fmt.Sprintf(`{"code":"ERR_NOTFOUND","message":"%v"}`, msg)
	default:
		return fmt.Sprintf(`{"code":"ERR_INTERNAL","message":"%v"}`, msg)
	}
}

func WrapGrpcError(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	var msg AppError
	ok := errors.As(err, &msg)
	if ok {
		_ = grpc.SetHeader(ctx, metadata.Pairs(errorNormalizedFlag, "OK"))
		return nil, status.Error(msg.StatusCode(), msg.ToJSON())
	}

	return nil, err
}

func FormatRestError(ctx context.Context, sm *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	md, _ := runtime.ServerMetadataFromContext(ctx)

	w.Header().Add("Content-Type", "application/json")
	grpcErr := status.Convert(err)

	statusCode := runtime.HTTPStatusFromCode(grpcErr.Code())
	if !isHTTPCodeAllowed(statusCode) {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	var bytes []byte
	if len(md.HeaderMD[errorNormalizedFlag]) > 0 {
		bytes = []byte(grpcErr.Message())
	} else {
		bytes = []byte(MessageForGrpcStatus(grpcErr))
	}
	_, _ = w.Write(bytes)
}

func isHTTPCodeAllowed(statusCode int) bool {
	for _, val := range AllowedHTTPErrorStatuses {
		if statusCode == val {
			return true
		}
	}
	return false
}
