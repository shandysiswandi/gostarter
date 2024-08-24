// Package dbops provides abstractions and utility functions for database operations.
//
// This package defines interfaces for querying and executing SQL commands,
// and includes helper functions to simplify common database operations.
package dbops

import (
	"strconv"
	"strings"
)

// OrderDirection specifies the direction of ordering in a SQL query.
type OrderDirection int

const (
	// ASC indicates ascending order.
	ASC OrderDirection = iota
	// DESC indicates descending order.
	DESC
)

// Placeholder specifies the type of placeholder to use in SQL queries.
type Placeholder int

const (
	// QuestionMark indicates the use of "?" as a placeholder in SQL queries.
	QuestionMark Placeholder = iota
	// Dollar indicates the use of "$" followed by an index as a placeholder in SQL queries.
	Dollar
)

// QueryBuilder constructs SQL SELECT queries with WHERE, ORDER BY, and other clauses.
type QueryBuilder struct {
	placeholder      Placeholder
	placeholderIndex int
	selectCols       string
	table            string
	whereClauses     []string
	orderClauses     []string
	args             []any
}

// New creates a new QueryBuilder instance with the specified placeholder type.
func New(placeholder ...Placeholder) *QueryBuilder {
	if len(placeholder) == 0 {
		return &QueryBuilder{
			placeholder: QuestionMark,
		}
	}

	return &QueryBuilder{
		placeholder: placeholder[0],
	}
}

// reset clears the internal state of the QueryBuilder, preparing it for reuse.
func (qb *QueryBuilder) reset() {
	qb.selectCols = ""
	qb.table = ""
	qb.whereClauses = nil
	qb.orderClauses = nil
	qb.args = nil
	qb.placeholderIndex = 0
}

// getPlaceholder returns the current placeholder string based on the placeholder type.
func (qb *QueryBuilder) getPlaceholder() string {
	if qb.placeholder == Dollar {
		qb.placeholderIndex++

		return "$" + strconv.Itoa(qb.placeholderIndex)
	}

	return "?"
}

// Select specifies the columns to select in the SQL query.
func (qb *QueryBuilder) Select(cols string) *QueryBuilder {
	qb.selectCols = cols

	return qb
}

// From specifies the table to select from in the SQL query.
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.table = table

	return qb
}

// Where adds a WHERE clause to the SQL query with the provided key-value pairs.
func (qb *QueryBuilder) Where(kv ...string) *QueryBuilder {
	if len(kv)%2 != 0 {
		return qb
	}

	for i := 0; i < len(kv); i += 2 {
		qb.whereClauses = append(qb.whereClauses, kv[i]+"="+qb.getPlaceholder())
		qb.args = append(qb.args, kv[i+1])
	}

	return qb
}

// WhereIn adds a WHERE IN clause to the SQL query for the specified column and values.
func (qb *QueryBuilder) WhereIn(column string, vals ...string) *QueryBuilder {
	placeholders := make([]string, len(vals))
	for i, v := range vals {
		placeholders[i] = qb.getPlaceholder()
		qb.args = append(qb.args, v)
	}
	qb.whereClauses = append(qb.whereClauses, column+" IN("+strings.Join(placeholders, ",")+")")

	return qb
}

// OrderBy adds an ORDER BY clause to the SQL query for the specified columns and direction.
func (qb *QueryBuilder) OrderBy(od OrderDirection, columns ...string) *QueryBuilder {
	direction := "ASC"
	if od == DESC {
		direction = "DESC"
	}

	for _, column := range columns {
		qb.orderClauses = append(qb.orderClauses, column+" "+direction)
	}

	return qb
}

// ToSQL generates the SQL query string and the arguments to be used with it.
func (qb *QueryBuilder) ToSQL() (string, []any) {
	if qb.selectCols == "" || qb.table == "" {
		return "", nil
	}

	query := "SELECT " + qb.selectCols + " FROM " + qb.table
	if len(qb.whereClauses) > 0 {
		query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
	}

	if len(qb.orderClauses) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderClauses, ", ")
	}

	sqlArgs := qb.args
	qb.reset()

	return query, sqlArgs
}
