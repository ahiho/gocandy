package validator

import (
	"reflect"
	"testing"
	"time"

	fp "github.com/ahiho/gocandy/filterparser"
)

func TestOperators_Contains(t *testing.T) {
	type args struct {
		op fp.Operator
	}
	tests := []struct {
		name string
		o    Operators
		args args
		want bool
	}{
		{
			name: "return true if operator is the first item",
			o:    Operators{fp.OpContains, fp.OpEqual, fp.OpGreater, fp.OpGreaterOrEqual, fp.OpIn, fp.OpLess, fp.OpLessOrEqual, fp.OpNotEqual, fp.OpRange},
			args: args{
				op: fp.OpContains,
			},
			want: true,
		},
		{
			name: "return true if operator is the last item",
			o:    Operators{fp.OpContains, fp.OpEqual, fp.OpGreater, fp.OpGreaterOrEqual, fp.OpIn, fp.OpLess, fp.OpLessOrEqual, fp.OpNotEqual, fp.OpRange},
			args: args{
				op: fp.OpRange,
			},
			want: true,
		},
		{
			name: "return true if operator is in the operators slice",
			o:    Operators{fp.OpContains, fp.OpEqual, fp.OpGreater, fp.OpGreaterOrEqual, fp.OpIn, fp.OpLess, fp.OpLessOrEqual, fp.OpNotEqual, fp.OpRange},
			args: args{
				op: fp.OpEqual,
			},
			want: true,
		},
		{
			name: "return false if operator is not in the operators slice",
			o:    Operators{fp.OpEqual, fp.OpGreater, fp.OpGreaterOrEqual, fp.OpIn, fp.OpLess, fp.OpLessOrEqual, fp.OpNotEqual, fp.OpRange},
			args: args{
				op: fp.OpContains,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Contains(tt.args.op); got != tt.want {
				t.Errorf("Operators.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRules_Field(t *testing.T) {
	type fields struct {
		Fields map[string]Field
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Field
		want1  bool
	}{
		{
			name: "invalid key",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				key: "a",
			},
			want:  Field{Type: 0, Container: fp.ContainerNone},
			want1: false,
		},
		{
			name: "If the input is NOT a nested field name, the top level field must NOT be a map",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				key: "firstName",
			},
			want:  Field{Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}},
			want1: false,
		},
		{
			name: "key without dot character",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				key: "firstName",
			},
			want:  Field{Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}},
			want1: true,
		},
		{
			name: "key with dot character",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				key: "firstName.us",
			},
			want:  Field{Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}},
			want1: true,
		},
		{
			name: "key with multiple dots character",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				key: "firstName.us.abc",
			},
			want:  Field{Type: fp.TypeBool, Container: fp.ContainerMap, Operators: Operators{fp.OpEqual}},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rules{
				Fields: tt.fields.Fields,
			}
			got, got1 := r.Field(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rules.Field() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Rules.Field() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRules_Validate(t *testing.T) {
	type fields struct {
		Fields map[string]Field
	}
	type args struct {
		f fp.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "validate a valid filer",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeString, Container: fp.ContainerNone, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				f: fp.Filter{Conditions: []fp.Condition{{
					Field: "firstName",
					Not:   false,
					Op:    fp.OpEqual,
					Values: []interface{}{
						"abc",
					},
				}}},
			},
			wantErr: false,
		},
		{
			name: "validate an invalid filer",
			fields: fields{
				Fields: map[string]Field{"firstName": {Type: fp.TypeString, Container: fp.ContainerNone, Operators: Operators{fp.OpEqual}}, "lastName": {Type: fp.TypeBool, Container: fp.ContainerArray, Operators: Operators{fp.OpEqual}}},
			},
			args: args{
				f: fp.Filter{Conditions: []fp.Condition{{
					Field: "firstName",
					Not:   false,
					Op:    fp.OpGreater,
					Values: []interface{}{
						"abc",
					},
				}}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rules{
				Fields: tt.fields.Fields,
			}
			if err := r.Validate(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Rules.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestField_Ops(t *testing.T) {
	type fields struct {
		Type      fp.Type
		Container fp.Container
		Operators Operators
	}
	tests := []struct {
		name   string
		fields fields
		want   Operators
	}{
		{
			name: "unknown type",
			fields: fields{
				Container: fp.ContainerArray,
			},
			want: nil,
		},
		{
			name: "type bool",
			fields: fields{
				Type: fp.TypeBool,
			},
			want: []fp.Operator{fp.OpEqual, fp.OpNotEqual},
		},
		{
			name: "type int",
			fields: fields{
				Type: fp.TypeInt,
			},
			want: []fp.Operator{fp.OpEqual, fp.OpNotEqual, fp.OpIn, fp.OpGreater, fp.OpLess, fp.OpGreaterOrEqual, fp.OpLessOrEqual, fp.OpRange},
		},
		{
			name: "type uint",
			fields: fields{
				Type: fp.TypeUInt,
			},
			want: []fp.Operator{fp.OpEqual, fp.OpNotEqual, fp.OpIn, fp.OpGreater, fp.OpLess, fp.OpGreaterOrEqual, fp.OpLessOrEqual, fp.OpRange},
		},
		{
			name: "type string",
			fields: fields{
				Type: fp.TypeString,
			},
			want: []fp.Operator{fp.OpEqual, fp.OpNotEqual, fp.OpIn, fp.OpContains},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field{
				Type:      tt.fields.Type,
				Container: tt.fields.Container,
				Operators: tt.fields.Operators,
			}
			if got := f.Ops(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field.Ops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateValue(t *testing.T) {
	type args struct {
		c     fp.Condition
		field Field
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "validate uint",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{"abc"}, Field: "firstName"},
				field: Field{Type: fp.TypeUInt, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: true,
		},
		{
			name: "validate int",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{1.2}, Field: "firstName"},
				field: Field{Type: fp.TypeInt, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: true,
		},
		{
			name: "validate wrong bool",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{"true"}, Field: "firstName"},
				field: Field{Type: fp.TypeBool, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: true,
		},
		{
			name: "validate correct bool",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{true}, Field: "firstName"},
				field: Field{Type: fp.TypeBool, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: false,
		},
		{
			name: "validate float",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{1.2}, Field: "firstName"},
				field: Field{Type: fp.TypeFloat, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: false,
		},
		{
			name: "validate wrong timestamp",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{1.2}, Field: "firstName"},
				field: Field{Type: fp.TypeTimestamp, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: true,
		},
		{
			name: "validate correct timestamp",
			args: args{
				c:     fp.Condition{Not: false, Values: []interface{}{time.Now()}, Field: "firstName"},
				field: Field{Type: fp.TypeTimestamp, Container: fp.ContainerNone, Operators: []fp.Operator{fp.OpEqual}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValue(tt.args.c, tt.args.field); (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkSameType(t *testing.T) {
	type args struct {
		cond fp.Condition
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "check difference type",
			args: args{
				cond: fp.Condition{Field: "firstName", Not: false, Op: fp.OpEqual, Values: []interface{}{1, "a"}},
			},
			wantErr: true,
		},
		{
			name: "check same type",
			args: args{
				cond: fp.Condition{Field: "firstName", Not: false, Op: fp.OpEqual, Values: []interface{}{1, 2}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkSameType(tt.args.cond); (err != nil) != tt.wantErr {
				t.Errorf("checkSameType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
