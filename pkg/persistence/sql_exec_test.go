package persistence

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	type args struct {
		ctx           context.Context
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
				ctx:    context.Background(),
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
				ctx: context.Background(),
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").WillReturnError(assert.AnError)

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
				ctx: context.Background(),
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").WillReturnResult(sqlmock.NewErrorResult(assert.AnError))

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
			name: "ErrorZeroRowsAffected",
			args: args{
				ctx: context.Background(),
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").WillReturnResult(sqlmock.NewResult(0, 0))

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
				ctx: context.Background(),
				execer: func() *sql.DB {
					db, mock, _ := sqlmock.New()

					mock.ExpectExec("DELETE FROM test WHERE id = ?").WillReturnResult(sqlmock.NewResult(1, 1))

					return db
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db := tt.args.execer()
			err := Exec(tt.args.ctx, db, tt.args.queryProvider, tt.args.feedback...)
			assert.Equal(t, tt.wantErr, err)
			if db != nil {
				db.Close()
			}
		})
	}
}
