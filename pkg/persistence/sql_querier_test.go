package persistence

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// testStruct is a placeholder struct for tests
type testStruct struct {
	ID   int
	Name string
}

// MockRow implementation for testStruct
func (r *testStruct) ScanColumn() []any {
	return []any{&r.ID, &r.Name}
}

func TestSQLGet(t *testing.T) {
	type args struct {
		ctx           context.Context
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
				ctx:     context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name")
					mock.ExpectQuery("SELECT id, name FROM test WHERE id = ?").WillReturnRows(rows)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test WHERE id = ?", []any{1}, nil
				},
			},
			want:    &testStruct{ID: 1, Name: "name"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := tt.args.querier()
			got, err := SQLGet[testStruct](tt.args.ctx, db, tt.args.queryProvider)
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
		ctx           context.Context
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name").AddRow(true, "anme2")
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
				ctx: context.Background(),
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name").AddRow(2, "anme2").RowError(1, assert.AnError)
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
				ctx: context.Background(),
				querier: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "name")
					mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(rows)

					return db
				},
				queryProvider: func() (string, []any, error) {
					return "SELECT id, name FROM test", nil, nil
				},
			},
			want:    []testStruct{{ID: 1, Name: "name"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := tt.args.querier()
			got, err := SQLGets[testStruct](tt.args.ctx, db, tt.args.queryProvider)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			if db != nil {
				db.Close()
			}
		})
	}
}
