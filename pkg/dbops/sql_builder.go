// Package dbops provides abstractions and utility functions for database operations.
//
// This package defines interfaces for querying and executing SQL commands,
// and includes helper functions to simplify common database operations.
package dbops

import (
	"errors"
	"strconv"
	"strings"
)

// ErrMissingField is an error if the query is missing select columns or table name.
var ErrMissingField = errors.New("select columns or table name is missing")

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

// PaginationType specifies the type of pagination algorithm.
type PaginationType int

const (
	// Noop indicates no operation on pagination.
	Noop PaginationType = iota
	// LimitOffset indicates traditional LIMIT/OFFSET pagination.
	LimitOffset
	// CursorBased indicates cursor-based pagination.
	CursorBased
)

const (
	// DefaultPaginationLimit is the default number of items per page if no limit is provided.
	DefaultPaginationLimit = 10
)

// pagination holds the pagination-related parameters.
type pagination struct {
	pType        PaginationType
	limit        int
	offset       int
	cursor       any
	cursorColumn string
}

// QueryBuilder constructs SQL SELECT queries with WHERE, ORDER BY, and other clauses.
type QueryBuilder struct {
	pagination       pagination
	placeholder      Placeholder
	placeholderIndex int
	selectCols       string
	table            string
	whereClauses     []string
	orderClauses     []string
	args             []any
}

// New creates a new QueryBuilder instance with the specified placeholder type.
// If no placeholder is provided, the default is QuestionMark.
func New(placeholder ...Placeholder) *QueryBuilder {
	return &QueryBuilder{
		placeholder:      defaultPlaceholder(placeholder),
		pagination:       pagination{pType: Noop, limit: DefaultPaginationLimit},
		placeholderIndex: 0,
	}
}

func defaultPlaceholder(placeholder []Placeholder) Placeholder {
	if len(placeholder) == 0 {
		return QuestionMark
	}

	return placeholder[0]
}

// reset clears the internal state of the QueryBuilder, preparing it for reuse.
func (qb *QueryBuilder) reset() {
	qb.selectCols = ""
	qb.table = ""
	qb.whereClauses = nil
	qb.orderClauses = nil
	qb.args = nil
	qb.placeholderIndex = 0
	qb.pagination = pagination{pType: Noop, limit: DefaultPaginationLimit}
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
// If no columns are provided, it defaults to selecting all columns ("*").
func (qb *QueryBuilder) Select(cols ...string) *QueryBuilder {
	qb.selectCols = "*"
	if len(cols) > 0 {
		qb.selectCols = strings.Join(cols, ",")
	}

	return qb
}

// From specifies the table to select from in the SQL query.
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.table = table

	return qb
}

// Where adds a WHERE clause to the SQL query with the provided key-value pairs.
// The key-value pairs must be provided in pairs, where the key is the column name
// and the value is the corresponding value to match.
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
// This allows filtering rows where the specified column matches any of the provided values.
func (qb *QueryBuilder) WhereIn(column string, vals ...string) *QueryBuilder {
	placeholders := make([]string, len(vals))
	for i, v := range vals {
		placeholders[i] = qb.getPlaceholder()
		qb.args = append(qb.args, v)
	}
	qb.whereClauses = append(qb.whereClauses, column+" IN("+strings.Join(placeholders, ",")+")")

	return qb
}

// Limit sets the maximum number of rows to return in the SQL query.
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	if limit > 0 {
		qb.pagination.limit = limit
	}

	return qb
}

// Offset sets the offset for the rows to return in the SQL query.
// This method enables limit-offset based pagination.
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.pagination.pType = LimitOffset
	qb.pagination.offset = offset

	return qb
}

// Seek enables cursor-based pagination by specifying the column and cursor value.
// The query will return rows where the column's value is greater than the cursor value.
func (qb *QueryBuilder) Seek(column string, value any) *QueryBuilder {
	qb.pagination.pType = CursorBased
	qb.pagination.cursorColumn = column
	qb.pagination.cursor = value

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
// It applies the specified WHERE, ORDER BY, and pagination clauses.
// The method returns the SQL query string and a slice of arguments to be passed to the SQL driver.
// Else will return an error if the query cannot be constructed due to missing select columns or table name.
func (qb *QueryBuilder) ToSQL() (string, []any, error) {
	if qb.selectCols == "" || qb.table == "" {
		return "", nil, ErrMissingField
	}

	query := "SELECT " + qb.selectCols + " FROM " + qb.table
	if len(qb.whereClauses) > 0 {
		query += " WHERE " + strings.Join(qb.whereClauses, " AND ")
	}

	if qb.pagination.pType == CursorBased && qb.pagination.cursor != "" {
		query += " AND " + qb.pagination.cursorColumn + " > " + qb.getPlaceholder()
		qb.args = append(qb.args, qb.pagination.cursor)
	}

	if len(qb.orderClauses) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderClauses, ", ")
	}

	switch qb.pagination.pType {
	case LimitOffset:
		if qb.pagination.limit <= 0 {
			qb.pagination.limit = DefaultPaginationLimit
		}
		query += " LIMIT " + strconv.Itoa(qb.pagination.limit)
		query += " OFFSET " + strconv.Itoa(qb.pagination.offset)
	case CursorBased:
		if qb.pagination.limit <= 0 {
			qb.pagination.limit = DefaultPaginationLimit
		}
		query += " LIMIT " + strconv.Itoa(qb.pagination.limit)
	default:
		// Optional to satisfying "exhaustive" linter
	}

	sqlArgs := qb.args
	qb.reset()

	return query, sqlArgs, nil
}
