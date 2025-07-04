package outbound

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLPayment(t *testing.T) {
	type args struct {
		db  *sql.DB
		qu  goqu.DialectWrapper
		tel *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *SQLPayment
	}{
		{
			name: "SuccessMySQL",
			args: args{
				db:  &sql.DB{},
				qu:  goqu.Dialect(dbops.MySQLDriver),
				tel: telemetry.NewTelemetry(),
			},
			want: &SQLPayment{
				db:        &sql.DB{},
				qu:        goqu.Dialect(dbops.MySQLDriver),
				telemetry: telemetry.NewTelemetry(),
			},
		},
		{
			name: "SuccessPostgres",
			args: args{
				db:  &sql.DB{},
				qu:  goqu.Dialect(dbops.MySQLDriver),
				tel: telemetry.NewTelemetry(),
			},
			want: &SQLPayment{
				db:        &sql.DB{},
				qu:        goqu.Dialect(dbops.MySQLDriver),
				telemetry: telemetry.NewTelemetry(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSQLPayment(tt.args.db, tt.args.qu, tt.args.tel)
			assert.Equal(t, tt.want.db, got.db)
			assert.Equal(t, tt.want.qu, got.qu)
			assert.Equal(t, tt.want.telemetry, got.telemetry)
		})
	}
}

func TestSQLPayment_FindAccountByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Account
		wantErr error
		mockFn  func(a args) (*SQLPayment, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx:    context.Background(),
				userID: 19,
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "balance").
					From("accounts").
					Where(goqu.Ex{"user_id": a.userID}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:    context.Background(),
				userID: 19,
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "balance").
					From("accounts").
					Where(goqu.Ex{"user_id": a.userID}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: 19,
			},
			want: &domain.Account{
				ID:       20,
				UserID:   19,
				Balanace: decimal.NewFromFloat32(100.17),
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "balance").
					From("accounts").
					Where(goqu.Ex{"user_id": a.userID}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "user_id", "balance"}).
					AddRow(20, 19, 100.17)

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()

			got, err := s.FindAccountByUserID(tt.args.ctx, tt.args.userID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLPayment_FindTopupByReferenceID(t *testing.T) {
	type args struct {
		ctx   context.Context
		refID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Topup
		wantErr error
		mockFn  func(a args) (*SQLPayment, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx:   context.Background(),
				refID: "uuid",
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "transaction_id", "reference_id", "amount").
					From("topups").
					Where(goqu.Ex{"reference_id": a.refID}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "NotFound",
			args: args{
				ctx:   context.Background(),
				refID: "uuid",
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "transaction_id", "reference_id", "amount").
					From("topups").
					Where(goqu.Ex{"reference_id": a.refID}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				refID: "uuid",
			},
			want: &domain.Topup{
				ID:            11,
				TransactionID: 12,
				ReferenceID:   "uuid",
				Amount:        decimal.NewFromFloat32(234.457),
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "transaction_id", "reference_id", "amount").
					From("topups").
					Where(goqu.Ex{"reference_id": a.refID}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "transaction_id", "reference_id", "amount"}).
					AddRow(11, 12, "uuid", 234.457)

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()

			got, err := s.FindTopupByReferenceID(tt.args.ctx, tt.args.refID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLPayment_SaveTopup(t *testing.T) {
	type args struct {
		ctx context.Context
		t   domain.Topup
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLPayment, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				t: domain.Topup{
					ID:            10,
					TransactionID: 20,
					ReferenceID:   "uuid",
					Amount:        decimal.NewFromInt(100),
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("topups").
					Cols("id", "transaction_id", "reference_id", "amount").
					Vals([]any{a.t.ID, a.t.TransactionID, a.t.ReferenceID, a.t.Amount}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "ErrorNoRowsAffected",
			args: args{
				ctx: context.Background(),
				t: domain.Topup{
					ID:            10,
					TransactionID: 20,
					ReferenceID:   "uuid",
					Amount:        decimal.NewFromInt(100),
				},
			},
			wantErr: domain.ErrAccountNoRowsAffected,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("topups").
					Cols("id", "transaction_id", "reference_id", "amount").
					Vals([]any{a.t.ID, a.t.TransactionID, a.t.ReferenceID, a.t.Amount}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				t: domain.Topup{
					ID:            10,
					TransactionID: 20,
					ReferenceID:   "uuid",
					Amount:        decimal.NewFromInt(100),
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("topups").
					Cols("id", "transaction_id", "reference_id", "amount").
					Vals([]any{a.t.ID, a.t.TransactionID, a.t.ReferenceID, a.t.Amount}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()

			err := s.SaveTopup(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLPayment_SaveTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		t   domain.Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLPayment, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				t: domain.Transaction{
					ID:       10,
					UserID:   20,
					Amount:   decimal.NewFromInt(100),
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "remark",
					CreateAt: time.Time{},
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("transactions").
					Cols("id", "user_id", "amount", "type", "status", "remark", "created_at").
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.Amount,
						a.t.Type,
						a.t.Status,
						a.t.Remark,
						a.t.CreateAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "ErrorNoRowsAffected",
			args: args{
				ctx: context.Background(),
				t: domain.Transaction{
					ID:       10,
					UserID:   20,
					Amount:   decimal.NewFromInt(100),
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "remark",
					CreateAt: time.Time{},
				},
			},
			wantErr: domain.ErrTransactionNoRowsAffected,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("transactions").
					Cols("id", "user_id", "amount", "type", "status", "remark", "created_at").
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.Amount,
						a.t.Type,
						a.t.Status,
						a.t.Remark,
						a.t.CreateAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				t: domain.Transaction{
					ID:       10,
					UserID:   20,
					Amount:   decimal.NewFromInt(100),
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "remark",
					CreateAt: time.Time{},
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("transactions").
					Cols("id", "user_id", "amount", "type", "status", "remark", "created_at").
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.Amount,
						a.t.Type,
						a.t.Status,
						a.t.Remark,
						a.t.CreateAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()

			err := s.SaveTransaction(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLPayment_UpdateAccount(t *testing.T) {
	type args struct {
		ctx  context.Context
		data map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLPayment, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				data: map[string]any{
					"id":      uint64(11),
					"balance": decimal.NewFromFloat(100),
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("accounts").
					Set(map[string]any{
						"balance": decimal.NewFromFloat(100),
					}).
					Where(goqu.Ex{"id": uint64(11)}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				data: map[string]any{
					"id":      uint64(11),
					"balance": decimal.NewFromFloat(100),
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLPayment, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("accounts").
					Set(map[string]any{
						"balance": decimal.NewFromFloat(100),
					}).
					Where(goqu.Ex{"id": uint64(11)}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLPayment{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()

			err := s.UpdateAccount(tt.args.ctx, tt.args.data)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
