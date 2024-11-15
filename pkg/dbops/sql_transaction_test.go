package dbops

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "Success",
			args: args{db: &sql.DB{}},
			want: &Transaction{db: &sql.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewTransaction(tt.args.db)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTransaction_Transaction(t *testing.T) {
	type args struct {
		ctx context.Context
		fn  func(ctx context.Context) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *Transaction
	}{
		{
			name: "ErrorBegin",
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					return nil
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *Transaction {
				db, mock, _ := sqlmock.New()

				mock.ExpectBegin().WillReturnError(assert.AnError)

				return &Transaction{
					db: db,
				}
			},
		},
		{
			name: "ErrorExecuteFunction",
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					return assert.AnError
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *Transaction {
				db, mock, _ := sqlmock.New()

				mock.ExpectBegin()
				mock.ExpectRollback()

				return &Transaction{
					db: db,
				}
			},
		},
		{
			name: "ErrorExecuteFunctionAndErrorRollback",
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					return assert.AnError
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) *Transaction {
				db, mock, _ := sqlmock.New()

				mock.ExpectBegin()
				mock.ExpectRollback().WillReturnError(assert.AnError)

				return &Transaction{
					db: db,
				}
			},
		},
		{
			name: "SuccessButErrorCommit",
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					return nil
				},
			},
			wantErr: nil,
			mockFn: func(a args) *Transaction {
				db, mock, _ := sqlmock.New()

				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(assert.AnError)

				return &Transaction{
					db: db,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				fn: func(ctx context.Context) error {
					panic(1)
				},
			},
			wantErr: nil,
			mockFn: func(a args) *Transaction {
				db, mock, _ := sqlmock.New()

				mock.ExpectBegin()
				mock.ExpectCommit()

				return &Transaction{
					db: db,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tr := tt.mockFn(tt.args)

			didPanic := false
			defer func() {
				if r := recover(); r != nil {
					didPanic = true
				}
				log.Println(didPanic)
			}()

			// Call the transaction function
			err := tr.Transaction(tt.args.ctx, tt.args.fn)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
