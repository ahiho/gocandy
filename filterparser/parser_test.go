// Package filterparser provides a parser for filter expressions.
package filterparser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	var a int64 = 1
	var b int64 = 2
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    *Filter
		wantErr bool
	}{
		{
			name: "1 field contain condition",
			args: args{
				input: "name : \"abc\"",
			},
			want: &Filter{
				Conditions: []Condition{{Field: "name", Not: false, Op: OpContains, Values: []interface{}{"abc"}}},
			},
			wantErr: false,
		},
		{
			name: "2 fields contain condition",
			args: args{
				input: "name : \"abc\" type = 1.2",
			},
			want: &Filter{
				Conditions: []Condition{
					{Field: "name", Not: false, Op: OpContains, Values: []interface{}{"abc"}},
					{Field: "type", Not: false, Op: OpEqual, Values: []interface{}{1.2}},
				},
			},
			wantErr: false,
		},
		{
			name: "2 fields with one field is integer",
			args: args{
				input: "name : \"abc\" type = 1",
			},
			want: &Filter{
				Conditions: []Condition{
					{Field: "name", Not: false, Op: OpContains, Values: []interface{}{"abc"}},
					{Field: "type", Not: false, Op: OpEqual, Values: []interface{}{a}},
				},
			},
			wantErr: false,
		},
		{
			name: "2 fields with IN filter",
			args: args{
				input: "name : \"abc\" type IN (1, 2)",
			},
			want: &Filter{
				Conditions: []Condition{
					{Field: "name", Not: false, Op: OpContains, Values: []interface{}{"abc"}},
					{Field: "type", Not: false, Op: OpIn, Values: []interface{}{a, b}},
				},
			},
			wantErr: false,
		},
		{
			name: "2 fields with range filter",
			args: args{
				input: "name : \"abc\" type :[1,2]",
			},
			want: &Filter{
				Conditions: []Condition{
					{Field: "name", Not: false, Op: OpContains, Values: []interface{}{"abc"}},
					{Field: "type", Not: false, Op: OpRange, Values: []interface{}{a, b}},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong input",
			args: args{
				input: "name : \"abc\" type [] 1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
