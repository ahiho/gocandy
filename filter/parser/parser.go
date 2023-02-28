package parser

import (
	"errors"
)

// Parser represents a parser, including a scanner and the underlying raw input.
// It also contains a small buffer to allow for two unscans.
type Parser struct {
	s        *Lexer
	raw      string
	literals []string
	groups   [][]string
}

// NewParser returns a new instance of Parser.
func NewParser(s string) *Parser {
	return &Parser{s: NewLexer([]byte(s)), raw: s}
}

// Parse takes the raw string and returns the root node of the AST.
func (p *Parser) Parse() (Node, error) {
	operation, err := p.parseOperation()
	if err != nil {
		return nil, err
	}
	// Try to peel like an onion.
	gate := operation.(*Operation).Gate
	if gate != "" && operation.(*Operation).RightNode == nil {
		return nil, errors.New("found open gate")
	}
	for gate == "" {
		operation = operation.(*Operation).LeftNode
		if operation == nil {
			return nil, errors.New("got nil operation")
		}
		if operation.Type() == EXPRESSION {
			break
		}
		gate = operation.(*Operation).Gate
	}
	return operation, nil
}

func (p *Parser) parseOperation() (Node, error) {
	op := &Operation{
		LeftNode:  nil,
		Gate:      "",
		RightNode: nil,
	}

	return op, nil
}

func (p *Parser) ParserHelper() {
	for {
		_, token, val := p.s.Scan()
		if IsTokenComparator(token.String()) && token.String() != val {
			p.literals = append(p.literals, token.String())
		}
		p.literals = append(p.literals, val)
		if token == EOF {
			return
		}
	}
}

func (p *Parser) ParserToGroups() {
	literals := p.literals
	group := []string{}
	count := 0
	for i, lit := range literals {
		group = append(group, lit)
		if lit == OPEN_BRACKET.String() {
			count += 1
		}
		if lit == CLOSED_BRACKET.String() {
			count -= 1
			if count == 0 {
				p.groups = append(p.groups, group)
				group = []string{}
			}
		}
		if i == len(literals)-1 {
			p.groups = append(p.groups, group)
		}
	}
}

func (p *Parser) ParserToSQL() (query string, values []string) {
	isWhereIn := false
	var valueIn string
	for _, group := range p.groups {
		for key, val := range group {
			if key > 0 && val != COMMA.String() && group[key-1] != OPEN_BRACKET.String() && group[key-1] != CLOSED_BRACKET.String() && IsTokenComparator(group[key-1]) {
				if val != "" {
					if key < len(group)-1 && group[key+1] == CLOSED_BRACKET.String() {
						query = query + "?"
					} else {
						query = query + "? "
					}
					values = append(values, val)
				}
			} else if val == IN.String() {
				query = query + val
				isWhereIn = true
			} else if isWhereIn {
				if val != OPEN_BRACKET.String() && val != CLOSED_BRACKET.String() && val != "" {
					if val != COMMA.String() {
						if valueIn == "" {
							valueIn = valueIn + val
						} else {
							valueIn = valueIn + " " + val
						}

					} else {
						values = append(values, valueIn)
						valueIn = ""
						query = query + "?,"
					}
				} else if val == CLOSED_BRACKET.String() {
					isWhereIn = false
					values = append(values, valueIn)
					query = query + "?" + val + " "
				} else {
					query = query + val
				}
			} else {
				if (key < len(group)-1 && group[key+1] == CLOSED_BRACKET.String()) || val == OPEN_BRACKET.String() {
					query = query + val
				} else if val != "" {
					query = query + val + " "
				}
			}
		}
	}
	return
}
