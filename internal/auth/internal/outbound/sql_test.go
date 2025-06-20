package outbound

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLAuth(t *testing.T) {
	type args struct {
		db        *sqlkit.DB
		telemetry *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *SQL
	}{
		{
			name: "Success",
			args: args{
				db:        &sqlkit.DB{},
				telemetry: telemetry.NewTelemetry(),
			},
			want: &SQL{
				db:        &sqlkit.DB{},
				telemetry: telemetry.NewTelemetry(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSQL(tt.args.db, tt.args.telemetry)
			assert.Equal(t, tt.want.db, got.db)
			assert.Equal(t, tt.want.telemetry, got.telemetry)
		})
	}
}

func TestSQL_UserByEmail(t *testing.T) {
	tel := telemetry.NewTelemetry()

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr error
		mockFn  func(a args) (*SQL, func() error)
	}{
		{
			name: "Success",
			args: args{
				ctx:   context.Background(),
				email: "test@test.com",
			},
			want: &domain.User{
				ID:       1,
				Name:     "name",
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQL, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))

				// "Query: actual sql: \"SELECT \"email\", \"id\", \"name\", \"password\", \"verified_at\" FROM \"users\" WHERE (\"email\" = 'test@test.com') LIMIT 1\" does not equal to expected \"\""

				query := "SELECT \"email\", \"id\", \"name\", \"password\", \"verified_at\" FROM \"users\" WHERE (\"email\" = '" + a.email + "') LIMIT 1"

				row := sqlmock.
					NewRows([]string{"id", "name", "email", "password", "verified_at"}).
					AddRow(1, "name", "test@test.com", "password", nil)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(row)

				return &SQL{
					db:        sqlkit.New("mysql", db, tel.Logger()),
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

			got, err := s.UserByEmail(tt.args.ctx, tt.args.email)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

/*


func TestSQL_SaveUser(t *testing.T) {
	type args struct {
		ctx context.Context
		u   domain.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				u:   domain.User{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("users").
					Cols("id", "name", "email", "password").
					Vals([]any{a.u.ID, a.u.Name, a.u.Email, a.u.Password}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				u:   domain.User{},
			},
			wantErr: domain.ErrUserNotCreated,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("users").
					Cols("id", "name", "email", "password").
					Vals([]any{a.u.ID, a.u.Name, a.u.Email, a.u.Password}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLAuth{
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
				u:   domain.User{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("users").
					Cols("id", "name", "email", "password").
					Vals([]any{a.u.ID, a.u.Name, a.u.Email, a.u.Password}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.SaveUser(tt.args.ctx, tt.args.u)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_SaveAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		a   domain.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				a:   domain.Account{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("accounts").
					Cols("id", "user_id").
					Vals([]any{a.a.ID, a.a.UserID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				a:   domain.Account{},
			},
			wantErr: domain.ErrAccountNotCreated,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("accounts").
					Cols("id", "user_id").
					Vals([]any{a.a.ID, a.a.UserID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLAuth{
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
				a:   domain.Account{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("accounts").
					Cols("id", "user_id").
					Vals([]any{a.a.ID, a.a.UserID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.SaveAccount(tt.args.ctx, tt.args.a)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_UpdateUserPassword(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   uint64
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx:  context.Background(),
				id:   10,
				pass: "password",
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"password": a.pass}).
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				id:   10,
				pass: "password",
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("users").
					Set(map[string]any{"password": a.pass}).
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.UpdateUserPassword(tt.args.ctx, tt.args.id, tt.args.pass)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_FindTokenByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		uid uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Token
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx: context.Background(),
				uid: 10,
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				uid: 10,
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLAuth{
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
				uid: 10,
			},
			want: &domain.Token{
				ID:               10,
				UserID:           101,
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiredAt:  time.Time{},
				RefreshExpiredAt: time.Time{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					}).
					AddRow(
						10,
						101,
						"access_token",
						"refresh_token",
						time.Time{},
						time.Time{},
					)

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLAuth{
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

			got, err := s.FindTokenByUserID(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLAuth_FindTokenByRefresh(t *testing.T) {
	type args struct {
		ctx context.Context
		ref string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Token
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx: context.Background(),
				ref: "refresh_token",
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"refresh_token": a.ref}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				ref: "refresh_token",
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"refresh_token": a.ref}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLAuth{
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
				ref: "refresh_token",
			},
			want: &domain.Token{
				ID:               10,
				UserID:           101,
				AccessToken:      "access_token",
				RefreshToken:     "refresh_token",
				AccessExpiredAt:  time.Time{},
				RefreshExpiredAt: time.Time{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					From("tokens").
					Where(goqu.Ex{"refresh_token": a.ref}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					}).
					AddRow(
						10,
						101,
						"access_token",
						"refresh_token",
						time.Time{},
						time.Time{},
					)

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLAuth{
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

			got, err := s.FindTokenByRefresh(tt.args.ctx, tt.args.ref)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLAuth_SaveToken(t *testing.T) {
	type args struct {
		ctx context.Context
		t   domain.Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				t:   domain.Token{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("tokens").
					Cols(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.AccessToken,
						a.t.RefreshToken,
						a.t.AccessExpiredAt,
						a.t.RefreshExpiredAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				t:   domain.Token{},
			},
			wantErr: domain.ErrTokenNoRowsAffected,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("tokens").
					Cols(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.AccessToken,
						a.t.RefreshToken,
						a.t.AccessExpiredAt,
						a.t.RefreshExpiredAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLAuth{
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
				t:   domain.Token{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("tokens").
					Cols(
						"id",
						"user_id",
						"access_token",
						"refresh_token",
						"access_expires_at",
						"refresh_expires_at",
					).
					Vals([]any{
						a.t.ID,
						a.t.UserID,
						a.t.AccessToken,
						a.t.RefreshToken,
						a.t.AccessExpiredAt,
						a.t.RefreshExpiredAt,
					}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.SaveToken(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_UpdateToken(t *testing.T) {
	type args struct {
		ctx context.Context
		t   domain.Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				t:   domain.Token{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("tokens").
					Set(map[string]any{
						"user_id":            a.t.UserID,
						"access_token":       a.t.AccessToken,
						"refresh_token":      a.t.RefreshToken,
						"access_expires_at":  a.t.AccessExpiredAt,
						"refresh_expires_at": a.t.RefreshExpiredAt,
					}).
					Where(goqu.Ex{"id": a.t.ID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				t:   domain.Token{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Update("tokens").
					Set(map[string]any{
						"user_id":            a.t.UserID,
						"access_token":       a.t.AccessToken,
						"refresh_token":      a.t.RefreshToken,
						"access_expires_at":  a.t.AccessExpiredAt,
						"refresh_expires_at": a.t.RefreshExpiredAt,
					}).
					Where(goqu.Ex{"id": a.t.ID}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.UpdateToken(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_FindPasswordResetByUserID(t *testing.T) {
	type args struct {
		ctx context.Context
		uid uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PasswordReset
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx: context.Background(),
				uid: 1,
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				uid: 1,
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLAuth{
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
				uid: 10,
			},
			want: &domain.PasswordReset{
				ID:        10,
				UserID:    101,
				Token:     "token",
				ExpiresAt: time.Time{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"user_id": a.uid}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "user_id", "token", "expires_at"}).
					AddRow(10, 101, "token", time.Time{})

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLAuth{
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

			got, err := s.FindPasswordResetByUserID(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLAuth_FindPasswordResetByToken(t *testing.T) {
	type args struct {
		ctx context.Context
		t   string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PasswordReset
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenQuery",
			args: args{
				ctx: context.Background(),
				t:   "token",
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"token": a.t}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				t:   "token",
			},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"token": a.t}).
					Prepared(true).
					ToSQL()

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(sql.ErrNoRows)

				return &SQLAuth{
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
				t:   "token",
			},
			want: &domain.PasswordReset{
				ID:        10,
				UserID:    101,
				Token:     "token",
				ExpiresAt: time.Time{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Select("id", "user_id", "token", "expires_at").
					From("password_resets").
					Where(goqu.Ex{"token": a.t}).
					Prepared(true).
					ToSQL()

				row := sqlmock.
					NewRows([]string{"id", "user_id", "token", "expires_at"}).
					AddRow(10, 101, "token", time.Time{})

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLAuth{
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

			got, err := s.FindPasswordResetByToken(tt.args.ctx, tt.args.t)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLAuth_SavePasswordReset(t *testing.T) {
	type args struct {
		ctx context.Context
		ps  domain.PasswordReset
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name: "ErrorWhenExec",
			args: args{
				ctx: context.Background(),
				ps:  domain.PasswordReset{},
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("password_resets").
					Cols("id", "user_id", "token", "expires_at").
					Vals([]any{a.ps.ID, a.ps.UserID, a.ps.Token, a.ps.ExpiresAt}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
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
				ps:  domain.PasswordReset{},
			},
			wantErr: domain.ErrPasswordResetNotCreated,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("password_resets").
					Cols("id", "user_id", "token", "expires_at").
					Vals([]any{a.ps.ID, a.ps.UserID, a.ps.Token, a.ps.ExpiresAt}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLAuth{
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
				ps:  domain.PasswordReset{},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Insert("password_resets").
					Cols("id", "user_id", "token", "expires_at").
					Vals([]any{a.ps.ID, a.ps.UserID, a.ps.Token, a.ps.ExpiresAt}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.SavePasswordReset(tt.args.ctx, tt.args.ps)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLAuth_DeletePasswordReset(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLAuth, func() error)
	}{
		{
			name:    "ErrorWhenExec",
			args:    args{ctx: context.Background(), id: 1},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Delete("password_resets").
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLAuth{
					db:        db,
					qu:        goqu.Dialect(dbops.MySQLDriver),
					telemetry: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.Background(), id: 1},
			wantErr: nil,
			mockFn: func(a args) (*SQLAuth, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.MySQLDriver).
					Delete("password_resets").
					Where(goqu.Ex{"id": a.id}).
					Prepared(true).
					ToSQL()

				mock.ExpectExec(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLAuth{
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

			err := s.DeletePasswordReset(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
*/
