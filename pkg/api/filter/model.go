package filter

const (
	OperatorNotEq   = "neq"
	OperatorEq      = "eq"
	OperatorGt      = "gt"
	OperatorGte     = "gte"
	OperatorLt      = "lt"
	OperatorLte     = "lte"
	OperatorOverlap = "overlap"
)

type Options struct {
	Fields []Field
	Limit  uint64
	Offset uint64
	Query  string
}

type Field struct {
	Name     string
	Operator string
	Value    interface{}
}
