package query

type CondType uint

const (
	Equals CondType = iota
	NotEquals
	GreaterEqualThan
	LessEqualThan
	LessThan
	GreaterThan
)

type Cond struct {
	Operator CondType
	Field    string
	Value    interface{}
}

func Equal(field string, value interface{}) Cond {
	return Cond{
		Operator: Equals,
		Field:    field,
		Value:    value,
	}
}

func GreaterEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: GreaterEqualThan,
		Field:    field,
		Value:    value,
	}
}

func LessEqual(field string, value interface{}) Cond {
	return Cond{
		Operator: LessEqualThan,
		Field:    field,
		Value:    value,
	}
}

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
	default:
		return "<="
	}
}
