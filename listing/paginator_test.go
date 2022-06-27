package listing

import (
	"reflect"
	"testing"

	"github.com/ahiho/gocandy/gormx/model"
)

func Test_splitRequest(t *testing.T) {
	type args struct {
		req Request
	}
	tests := []struct {
		name    string
		args    args
		want    CommonState
		want1   []byte
		wantErr bool
	}{
		{
			name: "empty page token",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "district",
					PageToken:  []byte{},
				},
			},
			want: CommonState{
				Collection: "district",
				Knobs: Knobs{
					ShowDeleted: false,
					PageSize:    10,
					Filter:      "name",
					OrderBy:     "name",
				},
			},
			want1:   nil,
			wantErr: false,
		},
		{
			name: "invalid page token",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "district",
					PageToken:  []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'},
				},
			},
			want: CommonState{
				Collection: "district",
				Knobs: Knobs{
					ShowDeleted: false,
					PageSize:    10,
					Filter:      "name",
					OrderBy:     "name",
				},
			},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "valid page token",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams",
					PageToken:  []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"name","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
			},
			want: CommonState{
				Collection: "v1/teams",
				Knobs: Knobs{
					ShowDeleted: false,
					PageSize:    10,
					Filter:      "name",
					OrderBy:     "name",
				},
			},
			want1:   []byte(`{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}`),
			wantErr: false,
		},
		{
			name: "parameters change",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams",
					PageToken:  []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
			},
			want: CommonState{
				Collection: "v1/teams",
				Knobs: Knobs{
					ShowDeleted: false,
					PageSize:    10,
					Filter:      "name",
					OrderBy:     "name",
				},
			},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "invalid page token -> wrong collection",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams2",
					PageToken:  []byte(`{"p":{"c":"v1/teams1","k":{"ShowDeleted":false,"PageSize":10,"Filter":"","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
			},
			want: CommonState{
				Collection: "v1/teams2",
				Knobs: Knobs{
					ShowDeleted: false,
					PageSize:    10,
					Filter:      "name",
					OrderBy:     "name",
				},
			},
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitRequest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("splitRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		req         Request
		commonState *CommonState
		implState   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "list parameter changed",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams2",
					PageToken:  []byte(`{"p":{"c":"v1/teams1","k":{"ShowDeleted":false,"PageSize":10,"Filter":"","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
				commonState: &CommonState{
					Collection: "v1/teams2",
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
				},
				implState: []byte(`{"p":{"c":"v1/teams1","k":{"ShowDeleted":false,"PageSize":10,"Filter":"","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
			},
			wantErr: true,
		},
		{
			name: "unmarshal impl state",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams",
					PageToken:  []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"name","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
				commonState: &CommonState{
					Collection: "v1/teams",
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
				},
				implState: []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"name","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
			},
			wantErr: true,
		},
		{
			name: "page size <= 0",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    0,
						Filter:      "name",
						OrderBy:     "name",
					},
					Collection: "v1/teams",
					PageToken:  []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"name","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
				commonState: &CommonState{
					Collection: "v1/teams",
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    10,
						Filter:      "name",
						OrderBy:     "name",
					},
				},
				implState: []byte(`{"p":{"c":"v1/teams","k":{"ShowDeleted":false,"PageSize":10,"Filter":"name","OrderBy":"name"}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
			},
			wantErr: true,
		},
		{
			name: "invalid page token",
			args: args{
				req: Request{
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    1,
					},
					Collection: "v1/teams",
					PageToken:  []byte(`{"p":{"c":"","k":{"ShowDeleted":false,"PageSize":1,"Filter":"","OrderBy":""}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
				},
				commonState: &CommonState{
					Collection: "v1/teams",
					Knobs: Knobs{
						ShowDeleted: false,
						PageSize:    1,
					},
				},
				implState: []byte(`{"p":{"c":"","k":{"ShowDeleted":false,"PageSize":1,"Filter":"","OrderBy":""}},"s":{"d":"2021-04-07T08:57:11.69826Z","i":"1541336602343641088"}}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args.req, tt.args.commonState, tt.args.implState); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type hasntPagination struct {
	CommonState
}

func (l *hasntPagination) Finish() {

}

func (l *hasntPagination) HasNextPage() bool {
	return false
}

func (l *hasntPagination) ImplState() interface{} {
	return nil
}

func (l *hasntPagination) ModelHook(*model.Common) {

}

type hasPagination struct {
	CommonState
}

func (l *hasPagination) Finish() {

}

func (l *hasPagination) HasNextPage() bool {
	return true
}

func (l *hasPagination) ImplState() interface{} {
	return nil
}

func (l *hasPagination) ModelHook(*model.Common) {

}

func TestFinish(t *testing.T) {
	type args struct {
		p Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "doesn't have next page",
			args: args{
				p: &hasntPagination{
					CommonState: CommonState{
						Collection: "",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "has next page",
			args: args{
				p: &hasPagination{
					CommonState: CommonState{
						Collection: "",
					},
				},
			},
			want:    []byte(`{"p":{"c":"","k":{"ShowDeleted":false,"PageSize":0,"Filter":"","OrderBy":""}},"s":null}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Finish(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Finish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Finish() = %v, want %v", got, tt.want)
			}
		})
	}
}
