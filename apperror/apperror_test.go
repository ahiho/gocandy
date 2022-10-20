package apperror

import (
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
)

func Test_appError_Error(t *testing.T) {
	type fields struct {
		Status   codes.Code
		Code     string
		Message  string
		InnerErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "OK",
			fields: fields{
				Status:  codes.NotFound,
				Code:    ErrNotFound,
				Message: "notfound",
			},
			want: "apperror:code=ERR_NOTFOUND;msg=notfound",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				Status:   tt.fields.Status,
				Code:     tt.fields.Code,
				Message:  tt.fields.Message,
				InnerErr: tt.fields.InnerErr,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("appError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appError_StackTrace(t *testing.T) {
	type fields struct {
		Status   codes.Code
		Code     string
		Message  string
		InnerErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "OK",
			fields: fields{},
		},
	}
	EnableStackTrace.Store(true)
	defer EnableStackTrace.Store(false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				Status:   tt.fields.Status,
				Code:     tt.fields.Code,
				Message:  tt.fields.Message,
				InnerErr: tt.fields.InnerErr,
			}
			if got := e.StackTrace(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appError.StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appError_StatusCode(t *testing.T) {
	type fields struct {
		Status   codes.Code
		Code     string
		Message  string
		InnerErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   codes.Code
	}{
		{
			name: "OK",
			fields: fields{
				Status: codes.Unauthenticated,
			},
			want: codes.Unauthenticated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				Status:   tt.fields.Status,
				Code:     tt.fields.Code,
				Message:  tt.fields.Message,
				InnerErr: tt.fields.InnerErr,
			}
			if got := e.StatusCode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("appError.StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appError_ToJSON(t *testing.T) {
	type fields struct {
		Status   codes.Code
		Code     string
		Message  string
		InnerErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "to json",
			fields: fields{
				Code:    ErrInternal,
				Message: "internal",
			},
			want: `{"code":"ERR_INTERNAL","message":"internal"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &appError{
				Status:   tt.fields.Status,
				Code:     tt.fields.Code,
				Message:  tt.fields.Message,
				InnerErr: tt.fields.InnerErr,
			}
			if got := e.ToJSON(); got != tt.want {
				t.Errorf("appError.ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr string
	}{
		{
			name:    "not found",
			msg:     "Record not found",
			wantErr: `apperror:code=ERR_NOTFOUND;msg=Record not found`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NotFound(tt.msg); err.Error() != tt.wantErr {
				t.Errorf("NotFound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotFoundWithCode(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "not found with code",
			args: args{
				code: "SENDER_NOT_FOUND",
				msg:  "Sender not found",
			},
			wantErr: "apperror:code=SENDER_NOT_FOUND;msg=Sender not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NotFoundWithCode(tt.args.code, tt.args.msg); err.Error() != tt.wantErr {
				t.Errorf("NotFoundWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotFoundWithCodeE(t *testing.T) {
	type args struct {
		code string
		msg  string
		err  error
	}
	tests := []struct {
		name         string
		args         args
		wantErr      string
		wantInnerErr string
	}{
		{
			name: "not found with code",
			args: args{
				code: "SENDER_NOT_FOUND",
				msg:  "Sender not found",
				err:  errors.New("inner err"),
			},
			wantErr:      "apperror:code=SENDER_NOT_FOUND;msg=Sender not found",
			wantInnerErr: "inner err",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NotFoundWithCodeE(tt.args.code, tt.args.msg, tt.args.err)
			if err.Error() != tt.wantErr {
				t.Errorf("NotFoundWithCodeE() error = %v, wantErr %v", err, tt.wantErr)
			}
			var appE = &appError{}
			ok := errors.As(err, &appE)
			if !ok || appE.InnerErr.Error() != tt.wantInnerErr {
				t.Errorf("NotFoundWithCodeE() inner error = %v, wantInnerErr %v", appE.InnerErr, tt.wantInnerErr)
			}
		})
	}
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr string
	}{
		{
			name:    "unauthorized",
			msg:     "Not login",
			wantErr: `apperror:code=ERR_UNAUTHORIZED;msg=Not login`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unauthorized(tt.msg); err.Error() != tt.wantErr {
				t.Errorf("Unauthorized() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnauthorizedWithCode(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "Unauthorized with code",
			args: args{
				code: "INVALID_PASSWORD",
				msg:  "Wrong password",
			},
			wantErr: "apperror:code=INVALID_PASSWORD;msg=Wrong password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnauthorizedWithCode(tt.args.code, tt.args.msg); err.Error() != tt.wantErr {
				t.Errorf("UnauthorizedWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnauthorizedWithCodeE(t *testing.T) {
	type args struct {
		code string
		msg  string
		err  error
	}
	tests := []struct {
		name         string
		args         args
		wantErr      string
		wantInnerErr string
	}{
		{
			name: "Unauthorized with code",
			args: args{
				code: "INVALID_PASSWORD",
				msg:  "Invalid pass",
				err:  errors.New("inner error"),
			},
			wantErr:      "apperror:code=INVALID_PASSWORD;msg=Invalid pass",
			wantInnerErr: "inner error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnauthorizedWithCodeE(tt.args.code, tt.args.msg, tt.args.err)
			if err.Error() != tt.wantErr {
				t.Errorf("UnauthorizedWithCodeE() error = %v, wantErr %v", err, tt.wantErr)
			}
			var appE = &appError{}
			ok := errors.As(err, &appE)
			if !ok || appE.InnerErr.Error() != tt.wantInnerErr {
				t.Errorf("UnauthorizedWithCodeE() inner error = %v, wantInnerErr %v", appE.InnerErr, tt.wantInnerErr)
			}
		})
	}
}

func TestForbidden(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr string
	}{
		{
			name:    "forbidden",
			msg:     "Can not access",
			wantErr: "apperror:code=ERR_FORBIDDEN;msg=Can not access",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Forbidden(tt.msg)
			if err.Error() != tt.wantErr {
				t.Errorf("Forbidden() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForbiddenWithCode(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "forbidden with code",
			args: args{
				code: "NOT_ALLOW",
				msg:  "Not allow",
			},
			wantErr: "apperror:code=NOT_ALLOW;msg=Not allow",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ForbiddenWithCode(tt.args.code, tt.args.msg)
			if err.Error() != tt.wantErr {
				t.Errorf("ForbiddenWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForbiddenWithCodeE(t *testing.T) {
	type args struct {
		code string
		msg  string
		err  error
	}
	tests := []struct {
		name      string
		args      args
		wantErr   string
		wantInErr string
	}{
		{
			name: "forbidden code with e",
			args: args{
				code: "DENIED",
				msg:  "Access denied",
				err:  errors.New("denied"),
			},
			wantErr:   "apperror:code=DENIED;msg=Access denied",
			wantInErr: "denied",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ForbiddenWithCodeE(tt.args.code, tt.args.msg, tt.args.err)
			if err.Error() != tt.wantErr {
				t.Errorf("ForbiddenWithCodeE() error = %v, wantErr %v", err, tt.wantErr)
			}
			var appE = &appError{}
			ok := errors.As(err, &appE)
			if !ok || appE.InnerErr.Error() != tt.wantInErr {
				t.Errorf("ForbiddenWithCodeE() inner error = %v, wantInnerErr %v", appE.InnerErr, tt.wantInErr)
			}
		})
	}
}

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantErr string
	}{
		{
			name:    "bad request",
			msg:     "Bad request",
			wantErr: "apperror:code=ERR_BAD_REQUEST;msg=Bad request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BadRequest(tt.msg)
			if err.Error() != tt.wantErr {
				t.Errorf("BadRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBadRequestWithCode(t *testing.T) {
	type args struct {
		code string
		msg  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "bad request",
			args: args{
				code: "INVALID_RQ",
				msg:  "Invalid",
			},
			wantErr: "apperror:code=INVALID_RQ;msg=Invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BadRequestWithCode(tt.args.code, tt.args.msg)
			if err.Error() != tt.wantErr {
				t.Errorf("BadRequestWithCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBadRequestWithCodeE(t *testing.T) {
	type args struct {
		code string
		msg  string
		err  error
	}
	tests := []struct {
		name      string
		args      args
		wantErr   string
		wantInErr string
	}{
		{
			name: "bad request code e",
			args: args{
				code: "BAD_RQ",
				msg:  "Bad rq",
				err:  errors.New("bad"),
			},
			wantErr:   "apperror:code=BAD_RQ;msg=Bad rq",
			wantInErr: "bad",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BadRequestWithCodeE(tt.args.code, tt.args.msg, tt.args.err)
			if err.Error() != tt.wantErr {
				t.Errorf("BadRequestWithCodeE() error = %v, wantErr %v", err, tt.wantErr)
			}
			var appE = &appError{}
			ok := errors.As(err, &appE)
			if !ok || appE.InnerErr.Error() != tt.wantInErr {
				t.Errorf("BadRequestWithCodeE() inner error = %v, wantInnerErr %v", appE.InnerErr, tt.wantInErr)
			}
		})
	}
}

func TestInternalError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *appError
	}{
		{
			name: "internal err",
			err:  errors.New("Some error"),
			want: &appError{
				Status:   codes.Internal,
				Code:     ErrInternal,
				Message:  "Some error",
				InnerErr: errors.New("Some error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InternalError(tt.err); got.Error() != tt.want.Error() {
				t.Errorf("InternalError() = %v, want %v", got, tt.want)
			}
		})
	}
}
