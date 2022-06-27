// Package fieldmask provides support manipulating field masks.
package fieldmask

import (
	"reflect"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/ahiho/gocandy/resource"
)

func TestMask_Contains(t *testing.T) {
	type fields struct {
		Fields   []string
		Resource resource.Resource
	}
	type args struct {
		field string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "mask contains the first item",
			fields: fields{
				Fields: []string{"first", "second", "third"},
			},
			args: args{
				field: "first",
			},
			want: true,
		},
		{
			name: "mask contains the last item",
			fields: fields{
				Fields: []string{"first", "second", "third"},
			},
			args: args{
				field: "third",
			},
			want: true,
		},
		{
			name: "mask doesn't contain the field",
			fields: fields{
				Fields: []string{"first", "second", "third"},
			},
			args: args{
				field: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Mask{
				Fields:   tt.fields.Fields,
				Resource: tt.fields.Resource,
			}
			if got := f.Contains(tt.args.field); got != tt.want {
				t.Errorf("Mask.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testResource struct {
	resource.Resource
}

func (t *testResource) RRN() string {
	panic("not implemented") // TODO: Implement
}

func (t *testResource) ResourceVersion() uint64 {
	panic("not implemented") // TODO: Implement
}

func (t *testResource) EventChannel() string {
	panic("not implemented") // TODO: Implement
}

func (t *testResource) IsFieldOutputOnly(field string) bool {
	return field == "first"
}

func TestMask_RemoveOutputOnly(t *testing.T) {
	type fields struct {
		Fields   []string
		Resource resource.Resource
	}
	tests := []struct {
		name            string
		fields          fields
		shouldBeRemoved []string
		want            []string
	}{
		{
			name: "output field should be removed",
			fields: fields{
				Fields:   []string{"first", "second", "third"},
				Resource: &testResource{},
			},
			shouldBeRemoved: []string{"first"},
			want:            []string{"second", "third"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Mask{
				Fields:   tt.fields.Fields,
				Resource: tt.fields.Resource,
			}
			f.RemoveOutputOnly()
			for _, f1 := range tt.want {
				for _, f2 := range tt.shouldBeRemoved {
					if f1 == f2 {
						t.Errorf("Mask.RemoveOutputOnly() %v should be removed", f2)
					}
				}
			}
			if !reflect.DeepEqual(tt.want, f.Fields) {
				t.Errorf("Mask.RemoveOutputOnly() = %v, want %v", f.Fields, tt.want)
			}
		})
	}
}

func TestNewResource(t *testing.T) {
	type args struct {
		mask *field_mask.FieldMask
		res  resource.Resource
	}
	tests := []struct {
		name    string
		args    args
		want    *Mask
		wantErr bool
	}{
		{
			name: "valid field mask",
			args: args{
				mask: &field_mask.FieldMask{
					Paths: []string{""},
				},
				res: &testResource{},
			},
			want: &Mask{
				Fields:   []string{""},
				Resource: &testResource{},
			},
			wantErr: false,
		},
		{
			name: "reach to max size",
			args: args{
				mask: &field_mask.FieldMask{
					Paths: []string{
						"one", "two", "three", "four", "five", "six", "seven", "eight", "nice", "ten",
						"one", "two", "three", "four", "five", "six", "seven", "eight", "nice", "ten",
						"one", "two", "three", "four", "five", "six", "seven", "eight", "nice", "ten",
						"one", "two", "three",
					},
				},
				res: &testResource{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResource(tt.args.mask, tt.args.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
