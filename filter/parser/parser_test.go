package parser

import (
	"reflect"
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
				s:   NewLexer([]byte("duckhue01")),
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
		test := `(status LIKE home AND (status IN ("To Do", "In Progress", "Closed") AND artifact=art1)) OR metric > 0.98`
		parser := NewParser(test)
		parser.ParserHelper()
		parser.ParserToGroups()
		query, _ := parser.ParserToSQL()
		queryWant := "(status LIKE ? AND (status IN(?,?,?) AND artifact = ?)) OR metric > ? "
		g.Expect(query).To(Equal(queryWant))
	})
}
