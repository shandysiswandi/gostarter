package sqlkit

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/telemetry/logger"
)

const (
	MySQLDriver    = "mysql"
	PostgresDriver = "postgres"
)

var (
	// ErrZeroRowsAffected is returned when an update, insert, or delete operation affects zero rows.
	ErrZeroRowsAffected = errors.New("no rows affected by an update, insert, or delete")

	// ErrScanRow is returned when scanning a row into the field type fails.
	ErrScanRow = errors.New("failed to scan column into field type")

	// ErrDestNotPointer is ...
	ErrDestNotPointer = errors.New("destination must be a pointer")

	// ErrDestNotSupport is ...
	ErrDestNotSupport = errors.New("destination not support yet")
)

type Expression = goqu.Expression
type Ex = goqu.Ex

type Model interface {
	Table() string
}

type DB struct {
	db  *sql.DB
	log logger.Logger
	qb  *goqu.Database
}

func New(driver string, db *sql.DB, log logger.Logger) *DB {
	return &DB{
		db:  db,
		log: log,
		qb:  goqu.New(driver, db),
	}
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Tx() Tx {
	return d
}

func (d *DB) Querier(ctx context.Context) Querier {
	tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx)
	if !ok {
		return d.db
	}

	return tx
}

func (d *DB) Execer(ctx context.Context) Execer {
	if tx, ok := ctx.Value(contextKeySQLTx{}).(*sql.Tx); ok {
		return tx
	}

	return d.db
}

func (d *DB) Scan(ctx context.Context, dest any, query string, args ...any) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return ErrDestNotPointer
	}

	elem := val.Elem()
	isSlice := elem.Kind() == reflect.Slice
	isStruct := elem.Kind() == reflect.Struct
	isUint64 := elem.Kind() == reflect.Uint64

	if !isSlice && !isStruct && !isUint64 {
		return ErrDestNotSupport
	}

	rows, err := d.Querier(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		d.log.Error(ctx, "error when query", err, logger.KeyVal("query", query))

		return err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			d.log.Error(ctx, "error when close rows", err)
		}
	}()

	// Handle count case (dest is *uint64)
	if elem.Kind() == reflect.Uint64 {
		return d.scanCount(rows, elem)
	}

	structType := d.getStructType(elem, isSlice)
	if structType == nil {
		return ErrDestNotSupport
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	fieldMap := d.mapDBTagsToStruct(structType)

	// Handle multiple rows into a slice of structs
	if isSlice {
		for rows.Next() {
			structVal := reflect.New(structType).Elem()
			if err := d.scanRow(rows, structVal, fieldMap, columns); err != nil {
				return err
			}
			elem.Set(reflect.Append(elem, structVal))
		}
		return rows.Err()
	}

	// Handle row into a struct
	if rows.Next() {
		return d.scanRow(rows, elem, fieldMap, columns)
	}

	return sql.ErrNoRows
}

func (d *DB) scanCount(rows *sql.Rows, dest reflect.Value) error {
	dest.SetUint(0)

	if rows.Next() {
		var count uint64
		if err := rows.Scan(&count); err != nil {
			return err
		}
		dest.SetUint(count)

		return nil
	}

	return sql.ErrNoRows
}

func (d *DB) getStructType(elem reflect.Value, isSlice bool) reflect.Type {
	if isSlice {
		if elem.Type().Elem().Kind() == reflect.Struct {
			return elem.Type().Elem()
		}

		return nil // is not slice of structs
	}

	return elem.Type()
}

func (d *DB) mapDBTagsToStruct(structType reflect.Type) map[string]int {
	fieldMap := make(map[string]int)

	for i := range structType.NumField() {
		if tag := structType.Field(i).Tag.Get("db"); tag != "" {
			fieldMap[tag] = i
		}
	}

	return fieldMap
}

func (d *DB) scanRow(rows *sql.Rows, val reflect.Value, fm map[string]int, cols []string) error {
	fields := make([]any, len(cols))
	fieldPtrs := make([]any, len(cols))

	for i, colName := range cols {
		if fieldIndex, ok := fm[colName]; ok {
			field := val.Field(fieldIndex)
			if field.CanAddr() {
				fields[i] = field.Addr().Interface()
				fieldPtrs[i] = &fields[i]
			}
		} else {
			var dummy any // Ignore unmapped columns
			fieldPtrs[i] = &dummy
		}
	}

	return rows.Scan(fieldPtrs...)
}
