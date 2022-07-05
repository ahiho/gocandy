package internal

import (
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		input string
		expr  *Expression
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "string with dot (.) character in the name",
			args: args{
				input: `props.key1 = "abc"`,
				expr:  &Expression{},
			},
		},
		{
			name: "string with bracket in the name",
			args: args{
				input: `[emails]: "email1"`,
				expr:  &Expression{},
			},
		},
		{
			name: "string with bracket and dot (.) in the name",
			args: args{
				input: `[emails].status: "ACTIVE"`,
				expr:  &Expression{},
			},
		},
		{
			name: "string with dot (.) inside the bracket in the name",
			args: args{
				input: `[emails.status]: "ACTIVE"`,
				expr:  &Expression{},
			},
			wantErr: true,
		},
		{
			name: "string with without close bracket in the name",
			args: args{
				input: `[emails: "ACTIVE"`,
				expr:  &Expression{},
			},
			wantErr: true,
		},
		{
			name: "string with without close bracket in the name 2",
			args: args{
				input: `[emails "ACTIVE"`,
				expr:  &Expression{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Parse(tt.args.input, tt.args.expr); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
