// Lexer tokens

package parser

// Token is the type of a single token.
type Token int

// TokenInfo stores relevant information about the token during scanning.
type TokenInfo struct {
	Token   Token
	Literal string
}

const (
	ILLEGAL Token = iota
	EOF
	APPEND
	ASSIGN
	COMMA
	EQUAL
	GTE
	GREATER
	LESS
	LPAREN
	LTE
	NOT
	NOT_EQUALS
	OR
	AND
	IN
	LIKE
	RPAREN
	EQUALS
	WS
	HASH

	NAME
	NUMBER
	STRING
	REGEX

	// Brackets
	OPEN_BRACKET
	CLOSED_BRACKET
)

var keywordTokens = map[string]Token{}

// KeywordToken returns the token associated with the given keyword
// string, or ILLEGAL if given name is not a keyword.
func KeywordToken(name string) Token {
	return keywordTokens[name]
}

var TokenNames = map[Token]string{
	ILLEGAL: "<illegal>",
	EOF:     "EOF",

	WS:         "WS",
	ASSIGN:     "=",
	COMMA:      ",",
	GTE:        ">=",
	GREATER:    ">",
	LESS:       "<",
	LPAREN:     "(",
	LTE:        "<=",
	NOT:        "!",
	NOT_EQUALS: "!=",
	EQUAL:      "!=",
	RPAREN:     ")",
	EQUALS:     "==",
	LIKE:       "LIKE",
	AND:        "AND",
	OR:         "OR",
	IN:         "IN",

	NAME:   "name",
	NUMBER: "number",
	STRING: "string",
	REGEX:  "regex",

	OPEN_BRACKET:   "(",
	CLOSED_BRACKET: ")",
}

// String returns the string name of this token.
func (t Token) String() string {
	return TokenNames[t]
}
