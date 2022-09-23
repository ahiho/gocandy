package deviceid

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

type Device string

var (
	deviceID Device = "x-device-id"
)

func TestDeviceID(t *testing.T) {
	c := context.Background()
	m := map[string]string{
		"x-device-id": "2f4c158bc0d7bacc",
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty request",
			args: args{
				ctx: context.WithValue(c, deviceID, ""),
			},
			want: "",
		},
		{
			name: "Get correct value",
			args: args{
				ctx: metadata.NewIncomingContext(c, metadata.New(m)),
			},
			want: "2f4c158bc0d7bacc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeviceID(tt.args.ctx); got != tt.want {
				t.Errorf("DeviceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
