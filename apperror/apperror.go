package apperror

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

// based on github/pkg/errors and grpc

const (
	ErrBadRequest   = "ERR_BAD_REQUEST"
	ErrUnauthorized = "ERR_UNAUTHORIZED"
	ErrForbidden    = "ERR_FORBIDDEN"
	ErrNotFound     = "ERR_NOTFOUND"
	ErrInternal     = "ERR_INTERNAL"
)

type Err struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace,omitempty"`
	InnerError error  `json:"-"`
}

type AppError interface {
	Error() string
	StatusCode() codes.Code
	StackTrace() string
	ToJSON() string
	Inner() error
}

type appError struct {
	Status   codes.Code
	Code     string
	Message  string
	InnerErr error
	Stack    errors.StackTrace
}

func (e *appError) Error() string {
	return fmt.Sprintf("apperror:code=%v;msg=%v", e.Code, e.Message)
}

func (e *appError) StackTrace() string {
	if len(e.Stack) == 0 {
		return ""
	}
	return fmt.Sprintf("%+v", e.Stack)
}

func (e *appError) StatusCode() codes.Code {
	return e.Status
}

func (e *appError) ToJSON() string {
	err := &Err{
		Code:       e.Code,
		Message:    e.Message,
		StackTrace: e.StackTrace(),
	}
	bytes, _ := json.Marshal(err)
	return string(bytes)
}

func (e *appError) Inner() error {
	return e.InnerErr
}

func NotFound(msg string) error {
	return makeError(codes.NotFound, ErrNotFound, msg, nil)
}

func NotFoundWithCode(code string, msg string) error {
	return makeError(codes.NotFound, code, msg, nil)
}

func NotFoundWithCodeE(code string, msg string, err error) error {
	return makeError(codes.NotFound, code, msg, err)
}

func Unauthorized(msg string) error {
	return makeError(codes.Unauthenticated, ErrUnauthorized, msg, nil)
}

func UnauthorizedWithCode(code string, msg string) error {
	return makeError(codes.Unauthenticated, code, msg, nil)
}

func UnauthorizedWithCodeE(code string, msg string, err error) error {
	return makeError(codes.Unauthenticated, code, msg, err)
}

func Forbidden(msg string) error {
	return makeError(codes.PermissionDenied, ErrForbidden, msg, nil)
}

func ForbiddenWithCode(code string, msg string) error {
	return makeError(codes.PermissionDenied, code, msg, nil)
}

func ForbiddenWithCodeE(code string, msg string, err error) error {
	return makeError(codes.PermissionDenied, code, msg, err)
}

func BadRequest(msg string) error {
	return makeError(codes.InvalidArgument, ErrBadRequest, msg, nil)
}

func BadRequestWithCode(code string, msg string) error {
	return makeError(codes.InvalidArgument, code, msg, nil)
}

func BadRequestWithCodeE(code string, msg string, err error) error {
	return makeError(codes.InvalidArgument, code, msg, err)
}

func InternalError(err error) AppError {
	msg := ""
	if err == nil {
		msg = "Undefined error"
	} else {
		msg = err.Error()
	}
	return makeError(codes.Internal, ErrInternal, msg, err)
}

func makeError(statusCode codes.Code, appErrCode string, msg string, err error) *appError {
	return &appError{
		Status:   statusCode,
		Code:     appErrCode,
		Message:  msg,
		InnerErr: err,
		Stack:    stackTrace(),
	}
}

func stackTrace() errors.StackTrace {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(4, pcs[:]) - 1
	var st = pcs[0:n]
	f := make([]errors.Frame, len(st))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame(st[i])
	}
	return f
}
