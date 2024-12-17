package outbound

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLTodo(t *testing.T) {
	type args struct {
		db  *sql.DB
		qu  goqu.DialectWrapper
		tel *telemetry.Telemetry
	}
	tests := []struct {
		name string
		args args
		want *SQLTodo
	}{
		{
			name: "SuccessMySQL",
			args: args{
				db: &sql.DB{},
				qu: goqu.Dialect(dbops.MySQLDriver),
			},
			want: &SQLTodo{db: &sql.DB{}, qu: goqu.Dialect(dbops.MySQLDriver)},
		},
		{
			name: "SuccessPostgres",
			args: args{
				db: &sql.DB{},
				qu: goqu.Dialect(dbops.PostgresDriver),
			},
			want: &SQLTodo{db: &sql.DB{}, qu: goqu.Dialect(dbops.PostgresDriver)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSQLTodo(tt.args.db, tt.args.qu, tt.args.tel)
			assert.Equal(t, tt.want.db, got.db)
			assert.Equal(t, tt.want.qu, got.qu)
		})
	}
}

func TestSQLTodo_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		todo domain.Todo
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) (*SQLTodo, func() error)
	}{
		{
			name:    "Error",
			args:    args{ctx: context.Background(), todo: domain.Todo{}},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name:    "SuccessButNoAffected",
			args:    args{ctx: context.Background(), todo: domain.Todo{}},
			wantErr: domain.ErrTodoNotCreated,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.Background(), todo: domain.Todo{}},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			err := s.Create(tt.args.ctx, tt.args.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLTodo_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		wantErr error
	}{
		{
			name:    "Error",
			args:    args{ctx: context.Background(), id: 1},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Delete("todos").
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.Background(), id: 1},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Delete("todos").
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			err := s.Delete(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLTodo_Find(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		want    *domain.Todo
		wantErr error
	}{
		{
			name:    "Error",
			args:    args{ctx: context.Background(), id: 1},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{ctx: context.Background(), id: 1},
			want: &domain.Todo{
				ID:          1,
				UserID:      11,
				Title:       "title test",
				Description: "description test",
				Status:      enum.New(domain.TodoStatusInProgress),
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status"}).
					AddRow(1, 11, "title test", "description test", "IN_PROGRESS")

				mock.ExpectQuery(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(row)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Find(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLTodo_Fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		in  map[string]any
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		want    []domain.Todo
		wantErr error
	}{
		{
			name:    "Error",
			args:    args{ctx: context.Background(), in: make(map[string]any)},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "SuccessWithoutFilter",
			args: args{ctx: context.Background(), in: make(map[string]any)},
			want: []domain.Todo{
				{
					ID:          1,
					UserID:      12,
					Title:       "title test",
					Description: "description test",
					Status:      enum.New(domain.TodoStatusDrop),
				}, {
					ID:          2,
					UserID:      13,
					Title:       "title test 2",
					Description: "description test 2",
					Status:      enum.New(domain.TodoStatusInitiate),
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Prepared(true).ToSQL()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status"}).
					AddRow(1, 12, "title test", "description test", "DROP").
					AddRow(2, 13, "title test 2", "description test 2", "INITIATE")

				mock.ExpectQuery(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(rows)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "SuccessWithFilter",
			args: args{ctx: context.Background(), in: map[string]any{
				"cursor": uint64(1),
				"limit":  int(1),
				"status": enum.New(domain.TodoStatusDone),
			}},
			want: []domain.Todo{
				{
					ID:          1,
					UserID:      12,
					Title:       "title test",
					Description: "description test",
					Status:      enum.New(domain.TodoStatusDrop),
				},
				{
					ID:          2,
					UserID:      13,
					Title:       "title test 2",
					Description: "description test 2",
					Status:      enum.New(domain.TodoStatusInitiate),
				},
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).
					Select("id", "user_id", "title", "description", "status").
					From("todos").
					Where(goqu.Ex{"id": goqu.Op{"gt": 1}}).
					Where(goqu.Ex{"status": "DONE"}).
					Limit(2).
					Prepared(true).
					ToSQL()

				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status"}).
					AddRow(1, 12, "title test", "description test", "DROP").
					AddRow(2, 13, "title test 2", "description test 2", "INITIATE")

				mock.ExpectQuery(query).
					WithArgs(dbops.AnyToValue(args)...).
					WillReturnRows(rows)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			got, err := s.Fetch(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSQLTodo_UpdateStatus(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
		sts enum.Enum[domain.TodoStatus]
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		wantErr error
	}{
		{
			name: "Error",
			args: args{
				ctx: context.Background(),
				id:  1,
				sts: enum.New(domain.TodoStatusDrop),
			},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").
					Set(map[string]any{"status": a.sts}).
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  1,
				sts: enum.New(domain.TodoStatusDone),
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").
					Set(map[string]any{"status": a.sts}).
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			err := s.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.sts)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestSQLTodo_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		todo domain.Todo
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		wantErr error
	}{
		{
			name:    "Error",
			args:    args{ctx: context.Background(), todo: domain.Todo{}},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").Set(map[string]any{
					"user_id":     a.todo.UserID,
					"title":       a.todo.Title,
					"description": a.todo.Description,
					"status":      a.todo.Status,
				}).Where(goqu.Ex{"id": a.todo.ID}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.Background(), todo: domain.Todo{}},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").Set(map[string]any{
					"user_id":     a.todo.UserID,
					"title":       a.todo.Title,
					"description": a.todo.Description,
					"status":      a.todo.Status,
				}).Where(goqu.Ex{"id": a.todo.ID}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(dbops.AnyToValue(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{
					db:  db,
					qu:  goqu.Dialect(dbops.PostgresDriver),
					tel: telemetry.NewTelemetry(),
				}, db.Close
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, dbMockCloser := tt.mockFn(tt.args)
			defer dbMockCloser()
			err := s.Update(tt.args.ctx, tt.args.todo)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
