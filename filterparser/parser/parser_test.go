package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewParser(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *Parser
	}{
		{
			name: "new parser",
			args: args{
				s: "duckhue01",
			},
			want: &Parser{
				s:   NewLexer(strings.NewReader("duckhue01")),
				raw: "duckhue01",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParser(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parse empty quoted string", func(t *testing.T) {
		test := "(name=\"\" AND artifact=art1) OR metric > 0.98"
		parser := NewParser(test)
		firstExpression := &Operation{
			LeftNode: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "",
			},
			Gate: "AND",
			RightNode: &Expression{
				Field:      "artifact",
				Comparator: "=",
				Value:      "art1",
			},
		}
		secondExpression := &Operation{
			LeftNode: firstExpression,
			Gate:     "OR",
			RightNode: &Expression{
				Field:      "metric",
				Comparator: ">",
				Value:      "0.98",
			},
		}
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node).To(Equal(secondExpression))
	})
	t.Run("parsing OR/AND is case insensitive", func(t *testing.T) {
		test := "name=model1 AND version=2.0"
		parser := NewParser(test)
		node, err := parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("AND"))

		test = "name=model1 and version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("AND"))

		test = "name=model1 OR version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("OR"))

		test = "name=model1 or version=2.0"
		parser = NewParser(test)
		node, err = parser.Parse()
		g.Expect(err).To(BeNil())
		g.Expect(node.(*Operation).Gate).To(Equal("OR"))
	})
	t.Run("parser fails for invalid queries with missing comparator", func(t *testing.T) {
		test := "name"
		parser := NewParser(test)
		_, err := parser.Parse()
		g.Expect(err).ToNot(BeNil())

		test = "name=default OR age"
		parser = NewParser(test)
		_, err = parser.Parse()
		g.Expect(err).ToNot(BeNil())

		test = ""
		parser = NewParser(test)
		_, err = parser.Parse()
		g.Expect(err).ToNot(BeNil())
	})
	t.Run("parser fails for invalid queries with open gate", func(t *testing.T) {
		test := "name=default AND"
		parser := NewParser(test)
		_, err := parser.Parse()
		g.Expect(err).ToNot(BeNil())
	})
}
func TestParser_parseOperation(t *testing.T) {
	g := NewGomegaWithT(t)

	t.Run("parse operations badly formatted return errors", func(t *testing.T) {
		tests := []string{
			"name=max AND AND artifact=wow",
			"name=max artifact=wow",
			")(name = max)",
		}
		for _, test := range tests {
			parser := NewParser(test)
			_, err := parser.parseOperation()
			g.Expect(err).ToNot(BeNil(), fmt.Sprintf("failed case: `%s`", test))
		}
	})

	tests := []struct {
		name    string
		query   string
		want    Node
		wantErr bool
	}{
		{
			name:  "parse operations correctly formatted succeeds",
			query: "name=max AND artifact%art1",
			want: &Operation{
				LeftNode: &Expression{
					Field:      "name",
					Comparator: "=",
					Value:      "max",
				},
				Gate: "AND",
				RightNode: &Expression{
					Field:      "artifact",
					Comparator: "%",
					Value:      "art1",
				},
			},
			wantErr: false,
		},
		{
			name:  "parse operations correctly formatted succeeds",
			query: "(name=max AND artifact=art1) OR metric > 0.98",
			want: &Operation{
				LeftNode: &Operation{
					LeftNode: &Expression{
						Field:      "name",
						Comparator: "=",
						Value:      "max",
					},
					Gate: "AND",
					RightNode: &Expression{
						Field:      "artifact",
						Comparator: "=",
						Value:      "art1",
					},
				},
				Gate: "OR",
				RightNode: &Expression{
					Field:      "metric",
					Comparator: ">",
					Value:      "0.98",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.query)
			got, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseExpression(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    Node
		wantErr bool
	}{
		{
			name:  "parse expression succeeds",
			query: "name=max",
			want: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "max",
			},
			wantErr: false,
		},
		{
			name:  "parse expression succeeds with whitespace",
			query: "name=  max ",
			want: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "max",
			},
			wantErr: false,
		},
		{
			name:    "parse expression fails invalid",
			query:   "name==dog",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "parse operation when just bracketed expression succeeds",
			query: "(name=max)",
			want: &Expression{
				Field:      "name",
				Comparator: "=",
				Value:      "max",
			},
			wantErr: false,
		},
		{
			name:  "parse operation with metrics/tags format",
			query: "(metrics[metric-name_1]>= 0.98)",
			want: &Expression{
				Field:      "metrics[metric-name_1]",
				Comparator: ">=",
				Value:      "0.98",
			},
			wantErr: false,
		},
		{
			name:  "parse operation camelCase",
			query: "TaskType=classification",
			want: &Expression{
				Field:      "TaskType",
				Comparator: "=",
				Value:      "classification",
			},
			wantErr: false,
		},
		{
			name:  "parse operation quoted string",
			query: "(Name=\"Iris Classifier\")",
			want: &Expression{
				Field:      "Name",
				Comparator: "=",
				Value:      "Iris Classifier",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.query)
			got, err := parser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzParser(f *testing.F) {
	testcases := []string{">=!=", "name=default OR age", "< <= = != AND OR and or", "1  !=   \"2\"", "(Name=\"Iris Classifier\")"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		parser := NewParser(s)
		node, err := parser.Parse()
		if err == nil {
			_, nodeIsOp := node.(*Operation)
			_, nodeIsExpr := node.(*Expression)
			if !nodeIsOp && !nodeIsExpr {
				t.Errorf("node must be either op or expression")
			}
		}
	})
}
