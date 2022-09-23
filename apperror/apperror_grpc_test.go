package apperror

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMessageForGrpcStatus(t *testing.T) {
	tests := []struct {
		name string
		err  *status.Status
		want string
	}{
		{
			name: "notfound",
			err:  status.New(codes.NotFound, "notfound"),
			want: `{"code":"ERR_NOTFOUND","message":"notfound"}`,
		},
		{
			name: "forbidden",
			err:  status.New(codes.PermissionDenied, "forbidden"),
			want: `{"code":"ERR_FORBIDDEN","message":"forbidden"}`,
		},
		{
			name: "unauthorize",
			err:  status.New(codes.Unauthenticated, "unauthorize"),
			want: `{"code":"ERR_UNAUTHORIZED","message":"unauthorize"}`,
		},
		{
			name: "badrequest",
			err:  status.New(codes.InvalidArgument, "badrequest"),
			want: `{"code":"ERR_BAD_REQUEST","message":"badrequest"}`,
		},
		{
			name: "internal",
			err:  status.New(codes.Internal, "internal"),
			want: `{"code":"ERR_INTERNAL","message":"internal"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MessageForGrpcStatus(tt.err); got != tt.want {
				t.Errorf("MessageForGrpcStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapGrpcError(t *testing.T) {
	errNotfound := &appError{
		Status:  codes.NotFound,
		Code:    ErrNotFound,
		Message: "notfound",
	}
	errGeneral := &appError{
		Status:  codes.Internal,
		Code:    ErrInternal,
		Message: "errorMsg",
	}
	tests := []struct {
		name     string
		handler  grpc.UnaryHandler
		wantResp interface{}
		wantErr  error
	}{
		{
			name: "OK response",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "OK", nil
			},
			wantResp: "OK",
			wantErr:  nil,
		},
		{
			name: "app error response",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, NotFound("notfound")
			},
			wantErr: status.Error(errNotfound.StatusCode(), errNotfound.ToJSON()),
		},
		{
			name: "general error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, errors.New("errorMsg")
			},
			wantErr: status.Error(errGeneral.Status, errGeneral.ToJSON()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := WrapGrpcError(context.Background(), nil, nil, tt.handler)
			if err == nil && tt.wantErr == nil {
				if !reflect.DeepEqual(gotResp, tt.wantResp) {
					t.Errorf("WrapGrpcError() = %v, want %v", gotResp, tt.wantResp)
				}
				return
			}
			if (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("WrapGrpcError() error = %v, want %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_isHTTPCodeAllowed(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "ok",
			statusCode: 400,
			want:       true,
		},
		{
			name:       "not ok",
			statusCode: 501,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHTTPCodeAllowed(tt.statusCode); got != tt.want {
				t.Errorf("isHTTPCodeAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatRestError(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		err        error
		wantBody   []byte
		wantStatus int
	}{
		{
			name:       "format error",
			ctx:        context.Background(),
			err:        status.Error(codes.Internal, "internal error"),
			wantBody:   []byte(`{"code":"ERR_INTERNAL","message":"internal error"}`),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "bypass error",
			ctx: runtime.NewServerMetadataContext(context.Background(), runtime.ServerMetadata{
				HeaderMD: metadata.Pairs(errorNormalizedFlag, "OK"),
			}),
			err:        status.Error(codes.NotFound, `Bypass`),
			wantBody:   []byte(`Bypass`),
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			FormatRestError(tt.ctx, nil, nil, w, nil, tt.err)
			bytes := w.Body.Bytes()
			if !reflect.DeepEqual(tt.wantBody, bytes) {
				t.Errorf("FormatRestError() Got body= %v, want %v", string(bytes), string(tt.wantBody))
			}
			if tt.wantStatus != w.Code {
				t.Errorf("FormatRestError() Got status= %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
