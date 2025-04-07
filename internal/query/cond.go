package query

// CondType is a type for conditionals
type CondType uint

const (
	// Equals is =
	Equals CondType = iota

	// NotEquals is !=
	NotEquals

	// GreaterEqualThan is >=
	GreaterEqualThan

	// LessEqualThan is <=
	LessEqualThan

	// LessThan is <
	LessThan

	// GreaterThan is >
	GreaterThan
)

// Cond is a structure for conditional
type Cond struct {
	Operator CondType
	Field    string
	Value    interface{}
}

// Equal creates conditional for Equals
func Equal(field string, value interface{}) Cond {
	return Cond{
		Operator: Equals,
		Field:    field,
		Value:    value,
	}
}

// GreaterEqual creates conditional for GreaterEqualThan
func GreaterEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: GreaterEqualThan,
		Field:    field,
		Value:    value,
	}
}

// LessEqual creates conditional for LessEqualThan
func LessEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: LessEqualThan,
		Field:    field,
		Value:    value,
	}
}

// NotEqual creates conditional for NotEquals
func NotEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: NotEquals,
		Field:    field,
		Value:    value,
	}
}

func (c *Cond) String() string {
	switch c.Operator {
	case Equals:
		return "="
	case NotEquals:
		return "<>"
	case GreaterEqualThan:
		return ">="
	case GreaterThan:
		return ">"
	case LessEqualThan:
		return "<="
	default:
		return "<"
	}
}
