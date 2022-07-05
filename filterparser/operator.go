package filterparser

// Operator represents a comparison operator between a field and one or more constants in a filter expression.
//
// A zero value for Operator is invalid, passing it to functions in this library causes undefined behavior or panic.
type Operator uint

const (
	OpEqual Operator = iota + 1
	OpNotEqual
	OpLess
	OpGreater
	OpLessOrEqual
	OpGreaterOrEqual
	OpContains
	OpIn
	OpRange
)

var operatorMap = map[string]Operator{
	"=":  OpEqual,
	"!=": OpNotEqual,
	"<":  OpLess,
	">":  OpGreater,
	"<=": OpLessOrEqual,
	">=": OpGreaterOrEqual,

	":":  OpContains,
	"[]": OpRange,
	"in": OpIn,
}

func (o Operator) String() string {
	for k, v := range operatorMap {
		if v == o {
			return k
		}
	}

	return "unknown"
}

// OperatorFromString returns an Operator based on a string.
//
// The following operators are available:
//  =, !=, <, >, <=, >=, in, :, []
// The last two are the contains operator and the range operator, respectively.
//
// The second return value indicates whether str was recognized.
func OperatorFromString(str string) (o Operator, ok bool) {
	op, ok := operatorMap[str]
	return op, ok
}
