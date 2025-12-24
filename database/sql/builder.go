package sql

import (
	"fmt"
	"strings"
)

// Builder is a SQL query builder
type Builder struct {
	selectFields []string
	fromTable    string
	joins        []string
	whereClause  []string
	groupBy      []string
	having       []string
	orderBy      []string
	limit        int
	offset       int
	args         []interface{}
}

// New creates a new SQL builder
func New() *Builder {
	return &Builder{
		selectFields: make([]string, 0),
		joins:        make([]string, 0),
		whereClause:  make([]string, 0),
		groupBy:      make([]string, 0),
		having:       make([]string, 0),
		orderBy:      make([]string, 0),
		args:         make([]interface{}, 0),
	}
}

// Select specifies the fields to select
func (b *Builder) Select(fields ...string) *Builder {
	b.selectFields = append(b.selectFields, fields...)
	return b
}

// From specifies the table to query from
func (b *Builder) From(table string) *Builder {
	b.fromTable = table
	return b
}

// Join adds a JOIN clause
func (b *Builder) Join(table, on string) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("JOIN %s ON %s", table, on))
	return b
}

// LeftJoin adds a LEFT JOIN clause
func (b *Builder) LeftJoin(table, on string) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("LEFT JOIN %s ON %s", table, on))
	return b
}

// RightJoin adds a RIGHT JOIN clause
func (b *Builder) RightJoin(table, on string) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("RIGHT JOIN %s ON %s", table, on))
	return b
}

// Where adds a WHERE clause
func (b *Builder) Where(condition string, args ...interface{}) *Builder {
	b.whereClause = append(b.whereClause, condition)
	b.args = append(b.args, args...)
	return b
}

// WhereOr adds an OR condition to WHERE clause
func (b *Builder) WhereOr(condition string, args ...interface{}) *Builder {
	if len(b.whereClause) > 0 {
		b.whereClause[len(b.whereClause)-1] = fmt.Sprintf("(%s OR %s)", b.whereClause[len(b.whereClause)-1], condition)
	} else {
		b.whereClause = append(b.whereClause, condition)
	}
	b.args = append(b.args, args...)
	return b
}

// GroupBy adds a GROUP BY clause
func (b *Builder) GroupBy(fields ...string) *Builder {
	b.groupBy = append(b.groupBy, fields...)
	return b
}

// Having adds a HAVING clause
func (b *Builder) Having(condition string, args ...interface{}) *Builder {
	b.having = append(b.having, condition)
	b.args = append(b.args, args...)
	return b
}

// OrderBy adds an ORDER BY clause
func (b *Builder) OrderBy(field string) *Builder {
	b.orderBy = append(b.orderBy, field)
	return b
}

// OrderByDesc adds a DESC ORDER BY clause
func (b *Builder) OrderByDesc(field string) *Builder {
	b.orderBy = append(b.orderBy, fmt.Sprintf("%s DESC", field))
	return b
}

// Limit sets the LIMIT clause
func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit
	return b
}

// Offset sets the OFFSET clause
func (b *Builder) Offset(offset int) *Builder {
	b.offset = offset
	return b
}

// Build builds the SQL query string
func (b *Builder) Build() (string, []interface{}) {
	var query strings.Builder

	// SELECT
	if len(b.selectFields) > 0 {
		query.WriteString("SELECT ")
		query.WriteString(strings.Join(b.selectFields, ", "))
	} else {
		query.WriteString("SELECT *")
	}

	// FROM
	if b.fromTable != "" {
		query.WriteString(" FROM ")
		query.WriteString(b.fromTable)
	}

	// JOINs
	if len(b.joins) > 0 {
		for _, join := range b.joins {
			query.WriteString(" ")
			query.WriteString(join)
		}
	}

	// WHERE
	if len(b.whereClause) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(b.whereClause, " AND "))
	}

	// GROUP BY
	if len(b.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(b.groupBy, ", "))
	}

	// HAVING
	if len(b.having) > 0 {
		query.WriteString(" HAVING ")
		query.WriteString(strings.Join(b.having, " AND "))
	}

	// ORDER BY
	if len(b.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(b.orderBy, ", "))
	}

	// LIMIT
	if b.limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
	}

	// OFFSET
	if b.offset > 0 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", b.offset))
	}

	return query.String(), b.args
}

// String returns the SQL query string (for debugging)
func (b *Builder) String() string {
	query, _ := b.Build()
	return query
}

// Reset resets the builder for reuse
func (b *Builder) Reset() *Builder {
	b.selectFields = make([]string, 0)
	b.fromTable = ""
	b.joins = make([]string, 0)
	b.whereClause = make([]string, 0)
	b.groupBy = make([]string, 0)
	b.having = make([]string, 0)
	b.orderBy = make([]string, 0)
	b.limit = 0
	b.offset = 0
	b.args = make([]interface{}, 0)
	return b
}
