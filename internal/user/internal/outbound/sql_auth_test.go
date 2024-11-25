package outbound

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func testconvertArgs(args []any) []driver.Value {
	var driverArgs []driver.Value

	for _, arg := range args {
		driverArgs = append(driverArgs, arg)
	}

	return driverArgs
}

func TestNewSQLAuth(t *testing.T) {
	type args struct {
		db  *sql.DB
		qu  goqu.DialectWrapper
		tel *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *SQLUser
	}{
		{
			name: "SuccessMySQL",
			args: args{
				db:  &sql.DB{},
				qu:  goqu.Dialect(dbops.MySQLDriver),
				tel: telemetry.NewTelemetry(),
			},
			want: &SQLUser{
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
			want: &SQLUser{
				db:        &sql.DB{},
				qu:        goqu.Dialect(dbops.MySQLDriver),
				telemetry: telemetry.NewTelemetry(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSQLUser(tt.args.db, tt.args.qu, tt.args.tel)
			assert.Equal(t, tt.want.db, got.db)
			assert.Equal(t, tt.want.qu, got.qu)
			assert.Equal(t, tt.want.telemetry, got.telemetry)
		})
	}
}

func TestSQLAuth_FindUserByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr error
		mockFn  func(a args) (*SQLUser, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx:   context.Background(),
				email: "test@test.com",
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLUser{
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
				email: "test@test.com",
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(testconvertArgs(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLUser{
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
				email: "test@test.com",
			},
			want: &domain.User{
				ID:       1,
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "password")

				mock.ExpectQuery(query).
					WithArgs(testconvertArgs(args)...).
					WillReturnRows(row)

				return &SQLUser{
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

			got, err := s.FindUserByEmail(tt.args.ctx, tt.args.email)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
