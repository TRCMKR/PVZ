package query

// CondType ...
type CondType uint

const (
	// Equals ...
	Equals CondType = iota

	// NotEquals ...
	NotEquals

	// GreaterEqualThan ...
	GreaterEqualThan

	// LessEqualThan ...
	LessEqualThan

	// LessThan ...
	LessThan

	// GreaterThan ...
	GreaterThan
)

// Cond ...
type Cond struct {
	Operator CondType
	Field    string
	Value    interface{}
}

// Equal ...
func Equal(field string, value interface{}) Cond {
	return Cond{
		Operator: Equals,
		Field:    field,
		Value:    value,
	}
}

// GreaterEqual ...
func GreaterEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: GreaterEqualThan,
		Field:    field,
		Value:    value,
	}
}

// LessEqual ...
func LessEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: LessEqualThan,
		Field:    field,
		Value:    value,
	}
}

// NotEqual ...
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
