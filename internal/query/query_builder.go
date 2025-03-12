package query

import (
	"fmt"
	"strconv"
	"strings"
)

type SelectQuery struct {
	from         string
	wheres       []string
	orderBy      string
	desc         bool
	limit        int
	offset       int
	currentIndex int
	args         []interface{}
}

type Param func(*SelectQuery)

func BuildSelectQuery(table string, params ...Param) (string, []interface{}) {
	s := SelectQuery{
		from:         table,
		wheres:       []string{},
		orderBy:      "",
		desc:         false,
		limit:        0,
		offset:       0,
		currentIndex: 1,
		args:         []interface{}{},
	}

	for _, param := range params {
		param(&s)
	}

	sb := strings.Builder{}
	sb.WriteString("SELECT * FROM ")
	sb.WriteString(s.from)
	if len(s.args) != 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(s.wheres, " AND "))
	}

	if s.orderBy != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(s.orderBy)
	}

	if s.desc {
		sb.WriteString(" DESC ")
	}

	if s.limit > 0 {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(s.limit))
	}

	if s.offset > 0 {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.Itoa(s.offset))
	}

	sb.WriteByte(';')

	return sb.String(), s.args
}

func Where(conds ...Cond) Param {
	return func(s *SelectQuery) {
		sb := strings.Builder{}
		for _, cond := range conds {
			sb.WriteString(cond.Field + " ")
			sb.WriteString(cond.String())
			sb.WriteString(fmt.Sprintf(" $%d", s.currentIndex))
			s.currentIndex++
			s.args = append(s.args, cond.Value)
			s.wheres = append(s.wheres, sb.String())
			sb.Reset()
		}
	}
}

func OrderBy(field string) Param {
	return func(s *SelectQuery) {
		s.orderBy = field
	}
}

func Desc(flag bool) Param {
	return func(s *SelectQuery) {
		s.desc = flag
	}
}

func Limit(limit int) Param {
	return func(s *SelectQuery) {
		s.limit = limit
	}
}

func Offset(offset int) Param {
	return func(s *SelectQuery) {
		s.offset = offset
	}
}
