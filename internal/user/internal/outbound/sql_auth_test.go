package outbound

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLUser(t *testing.T) {
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

func TestSQLUser_FindUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
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
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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
				ctx: context.Background(),
				id:  1,
			},
			want: &domain.User{
				ID:       1,
				Name:     "full name",
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "name", "email", "password"}).
					AddRow(1, "full name", "test@test.com", "password")

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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

			got, err := s.FindUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLUser_FindUserByEmail(t *testing.T) {
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
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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
				Name:     "full name",
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "name", "email", "password").
					From("users").
					Where(goqu.Ex{"email": a.email}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "name", "email", "password"}).
					AddRow(1, "full name", "test@test.com", "password")

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
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

func TestSQLUser_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		user map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLUser, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx:  context.Background(),
				user: map[string]any{"id": 10, "name": "fullname"},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"name": a.user["name"]}).
					Where(goqu.Ex{"id": a.user["id"]}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

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
				ctx:  context.Background(),
				user: map[string]any{"id": 10, "name": "fullname"},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"name": a.user["name"]}).
					Where(goqu.Ex{"id": a.user["id"]}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

			err := s.Update(tt.args.ctx, tt.args.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLUser_UpdatePassword(t *testing.T) {
	type args struct {
		ctx  context.Context
		user domain.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLUser, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				user: domain.User{
					ID:       10,
					Password: "new_password",
				},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"password": a.user.Password}).
					Where(goqu.Ex{"id": a.user.ID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

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
				ctx: context.Background(),
				user: domain.User{
					ID:       10,
					Password: "new_password",
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"password": a.user.Password}).
					Where(goqu.Ex{"id": a.user.ID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

			err := s.UpdatePassword(tt.args.ctx, tt.args.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLUser_DeleteTokenByAccess(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLUser, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx:   context.Background(),
				token: "token",
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Delete("tokens").
					Where(goqu.Ex{"access_token": a.token}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

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
				token: "token",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLUser, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Delete("tokens").
					Where(goqu.Ex{"access_token": a.token}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

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

			err := s.DeleteTokenByAccess(tt.args.ctx, tt.args.token)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
