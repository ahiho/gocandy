package parser

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func lexerHelper(lex *Lexer) ([]Token, []string) {
	var tokens []Token
	var literals []string
	for {
		_, token, val := lex.Scan()
		tokens = append(tokens, token)
		literals = append(literals, val)
		if token == EOF {
			return tokens, literals
		}
	}
}

func TestLexer(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("scan into tokens succeeds", func(t *testing.T) {
		s := `(name=\"\" AND artifact=art1) OR metric > 0.98`
		lexer := NewLexer([]byte(s))
		_, literals := lexerHelper(lexer)
		fmt.Println(literals)
		g.Expect(literals).To(Equal([]string{"(", "name", "=", "", "AND", "artifact", "=", "art1", ")", "OR", "metric", ">", "0.98", ""}))
	})
}
