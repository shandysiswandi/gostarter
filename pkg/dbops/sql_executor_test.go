package dbops

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// testStruct is a placeholder struct for tests
type testStruct struct {
	id   int
	name string
}

// MockRow implementation for testStruct
func (r *testStruct) ScanColumn() []any {
	return []any{&r.id, &r.name}
}

func TestExec(t *testing.T) {
	SetVerbose(true)

	type args struct {
		ctx           func() context.Context
		execer        func() *sql.DB
		queryProvider func() (string, []any, error)
		feedback      []bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "ErrorQueryProvider",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				execer: func() *sql.DB { return nil },
				queryProvider: func() (string, []any, error) {
					return "", []any{}, assert.AnError
				},
			},
			wantErr: assert.AnError,
		},
		{
			name: "ErrorExec",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").
						WillReturnError(assert.AnError)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "DELETE FROM test WHERE id = ?", []any{1}, nil
				},
			},
			wantErr: assert.AnError,
		},
		{
			name: "ErrorGetRowsAffected",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").
						WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "DELETE FROM test WHERE id = ?", []any{1}, nil
				},
				feedback: []bool{true},
			},
			wantErr: assert.AnError,
		},
		{
			name: "ErrorGetRowsAffectedWithoutFeedback",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").
						WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "DELETE FROM test WHERE id = ?", []any{1}, nil
				},
				feedback: nil,
			},
			wantErr: assert.AnError,
		},
		{
			name: "ErrorZeroRowsAffected",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").
						WillReturnResult(sqlmock.NewResult(0, 0))

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "DELETE FROM test WHERE id = ?", []any{1}, nil
				},
				feedback: []bool{true},
			},
			wantErr: ErrZeroRowsAffected,
		},
		{
			name: "Success",
			args: args{
				ctx: func() context.Context {
					db, mock, _ := sqlmock.New()
					mock.ExpectBegin()
					tx, _ := db.Begin()

					ctx := context.Background()
					ctx = context.WithValue(ctx, contextKeySQLTx{}, tx)

					mock.ExpectExec("DELETE FROM test WHERE id = ?").
						WillReturnResult(sqlmock.NewResult(1, 1))

					return ctx
				},
				execer: func() *sql.DB {
					return nil
				},
				queryProvider: func() (string, []any, error) {
					return "DELETE FROM test WHERE id = ?", []any{1}, nil
				},
				feedback: []bool{true},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := tt.args.execer()
			err := Exec(tt.args.ctx(), db, tt.args.queryProvider, tt.args.feedback...)
			assert.Equal(t, tt.wantErr, err)
			if db != nil {
				db.Close()
			}
		})
	}
}

func TestSQLGet(t *testing.T) {
	type args struct {
		ctx           func() context.Context
		querier       func() *sql.DB
		queryProvider func() (string, []any, error)
	}
	tests := []struct {
		name    string
		args    args
		want    *testStruct
		wantErr error
	}{
		{
			name: "ErrorQueryProvider",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB { return nil },
				queryProvider: func() (string, []any, error) {
					return "", []any{}, assert.AnError
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
		{
			name: "ErrorScanNoRows",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectQuery("SELECT id, name FROM test WHERE id = ?").WillReturnError(sql.ErrNoRows)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test WHERE id = ?", []any{1}, nil
				},
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "ErrorScanInternal",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectQuery("SELECT id, name FROM test WHERE id = ?").WillReturnError(assert.AnError)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test WHERE id = ?", []any{1}, nil
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
		{
			name: "Success",
			args: args{
				ctx: func() context.Context {
					db, mock, _ := sqlmock.New()
					mock.ExpectBegin()
					tx, _ := db.Begin()

					ctx := context.Background()
					ctx = context.WithValue(ctx, contextKeySQLTx{}, tx)

					rows := sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "name")

					mock.ExpectQuery("SELECT id, name FROM test WHERE id = ?").
						WillReturnRows(rows)

					return ctx
				},
				querier: func() *sql.DB {
					return nil
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test WHERE id = ?", []any{1}, nil
				},
			},
			want:    &testStruct{id: 1, name: "name"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := tt.args.querier()
			got, err := SQLGet[testStruct](tt.args.ctx(), db, tt.args.queryProvider)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			if db != nil {
				db.Close()
			}
		})
	}
}

func TestSQLGets(t *testing.T) {
	type args struct {
		ctx           func() context.Context
		querier       func() *sql.DB
		queryProvider func() (string, []any, error)
	}
	tests := []struct {
		name    string
		args    args
		want    []testStruct
		wantErr error
	}{
		{
			name: "ErrorQueryProvider",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					return nil
				},
				queryProvider: func() (string, []any, error) {
					return "", nil, assert.AnError
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
		{
			name: "ErrorQuery",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectQuery("SELECT id, name FROM test").WillReturnError(assert.AnError)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test", nil, nil
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
		{
			name: "ErrorScanTypeInternal",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name").AddRow(true, "name2")
					mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(rows)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test", nil, nil
				},
			},
			want:    nil,
			wantErr: ErrScanRow,
		},
		{
			name: "ErrorScanInternal",
			args: args{
				ctx: func() context.Context {
					return context.Background()
				},
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name").AddRow(2, "name2").RowError(1, assert.AnError)
					mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(rows)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test", nil, nil
				},
			},
			want:    nil,
			wantErr: assert.AnError,
		},
		{
			name: "Success",
			args: args{
				ctx: func() context.Context {
					db, mock, _ := sqlmock.New()
					mock.ExpectBegin()
					tx, _ := db.Begin()

					ctx := context.Background()
					ctx = context.WithValue(ctx, contextKeySQLTx{}, tx)

					rows := sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "name")

					mock.ExpectQuery("SELECT id, name FROM test").
						WillReturnRows(rows)

					return ctx
				},
				querier: func() *sql.DB {
					return nil
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test", nil, nil
				},
			},
			want:    []testStruct{{id: 1, name: "name"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := tt.args.querier()
			got, err := SQLGets[testStruct](tt.args.ctx(), db, tt.args.queryProvider)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			if db != nil {
				db.Close()
			}
		})
	}
}
