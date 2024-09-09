package outbound

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shandysiswandi/gostarter/internal/region/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/config/mocker"
	"github.com/stretchr/testify/assert"
)

func testconvertArgs(args []any) []driver.Value {
	var driverArgs []driver.Value

	for _, arg := range args {
		driverArgs = append(driverArgs, arg)
	}

	return driverArgs
}

func TestNewSQLRegion(t *testing.T) {
	type args struct {
		db     *sql.DB
		config config.Config
	}
	tests := []struct {
		name string
		args func() args
		want *SQLRegion
	}{
		{
			name: "SuccessMySQL",
			args: func() args {
				mcfg := mocker.NewMockConfig(t)

				mcfg.EXPECT().GetString("database.driver").Return("mysql")

				return args{
					db:     &sql.DB{},
					config: mcfg,
				}
			},
			want: &SQLRegion{db: &sql.DB{}, qu: goqu.Dialect("mysql"), config: &mocker.MockConfig{}},
		},
		{
			name: "SuccessPostgres",
			args: func() args {
				mcfg := mocker.NewMockConfig(t)

				mcfg.EXPECT().GetString("database.driver").Return("postgres")

				return args{
					db:     &sql.DB{},
					config: mcfg,
				}
			},
			want: &SQLRegion{db: &sql.DB{}, qu: goqu.Dialect("postgres"), config: &mocker.MockConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			arg := tt.args()
			got := NewSQLRegion(arg.db, arg.config)
			assert.Equal(t, tt.want.db, got.db)
			assert.Equal(t, tt.want.qu, got.qu)
		})
	}
}

func TestSQLRegion_Provinces(t *testing.T) {
	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.Province
		wantErr error
		mockFn  func(a args) (*SQLRegion, func() error)
	}{

		{
			name:    "ErrorMySQL",
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, _, _ := goqu.Dialect("mysql").
					Select("id", "name").From("provinces").
					Limit(10).Prepared(true).ToSQL()
				mock.ExpectQuery(query).WithArgs(10).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "ErrorPostgres",
			args:    args{ctx: context.TODO()},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, _, _ := goqu.Dialect("postgres").
					Select("id", "name").From("provinces").
					Limit(10).Prepared(true).ToSQL()
				mock.ExpectQuery(query).WithArgs(10).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
		{
			name:    "SuccessMySQL",
			args:    args{ctx: context.TODO(), ids: []string{"1", "2"}},
			want:    []domain.Province{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("provinces").
					Where(exp.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "SuccessPostgres",
			args:    args{ctx: context.TODO(), ids: []string{"1", "2"}},
			want:    []domain.Province{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("provinces").
					Where(exp.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Provinces(tt.args.ctx, tt.args.ids...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLRegion_Cities(t *testing.T) {
	type args struct {
		ctx context.Context
		pID string
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.City
		wantErr error
		mockFn  func(a args) (*SQLRegion, func() error)
	}{
		{
			name:    "ErrorMySQL",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("cities").
					Where(goqu.Ex{"province_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "ErrorPostgres",
			args:    args{ctx: context.TODO(), ids: []string{"1"}},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("cities").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
		{
			name:    "SuccessMySQL",
			args:    args{ctx: context.TODO(), ids: []string{"1", "2"}},
			want:    []domain.City{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("cities").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "SuccessPostgres",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    []domain.City{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("cities").
					Where(goqu.Ex{"province_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Cities(tt.args.ctx, tt.args.pID, tt.args.ids...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLRegion_Districts(t *testing.T) {
	type args struct {
		ctx context.Context
		pID string
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.District
		wantErr error
		mockFn  func(a args) (*SQLRegion, func() error)
	}{
		{
			name:    "ErrorMySQL",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("districts").
					Where(goqu.Ex{"city_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "ErrorPostgres",
			args:    args{ctx: context.TODO(), ids: []string{"1"}},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("districts").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
		{
			name:    "SuccessMySQL",
			args:    args{ctx: context.TODO(), ids: []string{"1", "2"}},
			want:    []domain.District{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("districts").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "SuccessPostgres",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    []domain.District{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("districts").
					Where(goqu.Ex{"city_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Districts(tt.args.ctx, tt.args.pID, tt.args.ids...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLRegion_Villages(t *testing.T) {
	type args struct {
		ctx context.Context
		pID string
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []domain.Village
		wantErr error
		mockFn  func(a args) (*SQLRegion, func() error)
	}{
		{
			name:    "ErrorMySQL",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("villages").
					Where(goqu.Ex{"district_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "ErrorPostgres",
			args:    args{ctx: context.TODO(), ids: []string{"1"}},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("villages").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
		{
			name:    "SuccessMySQL",
			args:    args{ctx: context.TODO(), ids: []string{"1", "2"}},
			want:    []domain.Village{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("mysql").
					Select("id", "name").From("villages").
					Where(goqu.Ex{"id": a.ids}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("mysql"),
				}, db.Close
			},
		},
		{
			name:    "SuccessPostgres",
			args:    args{ctx: context.TODO(), pID: "100"},
			want:    []domain.Village{{ID: "1", Name: "test 1"}, {ID: "2", Name: "test 2"}},
			wantErr: nil,
			mockFn: func(a args) (*SQLRegion, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect("postgres").
					Select("id", "name").From("villages").
					Where(goqu.Ex{"district_id": a.pID}).
					Limit(10).Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("1", "test 1").
					AddRow("2", "test 2")
				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).WillReturnRows(rows)

				return &SQLRegion{
					db: db,
					qu: goqu.Dialect("postgres"),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Villages(tt.args.ctx, tt.args.pID, tt.args.ids...)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
