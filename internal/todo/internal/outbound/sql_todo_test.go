package outbound

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/stretchr/testify/assert"
)

func testconvertArgs(args []any) []driver.Value {
	var driverArgs []driver.Value

	for _, arg := range args {
		driverArgs = append(driverArgs, arg)
	}

	return driverArgs
}

func TestNewSQLTodo(t *testing.T) {
	type args struct {
		db *sql.DB
		qu goqu.DialectWrapper
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
			got := NewSQLTodo(tt.args.db, tt.args.qu)
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
			args:    args{ctx: context.TODO(), todo: domain.Todo{}},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name:    "SuccessButNoAffected",
			args:    args{ctx: context.TODO(), todo: domain.Todo{}},
			wantErr: domain.ErrTodoNotCreated,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(0, 0))

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), todo: domain.Todo{}},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Insert("todos").
					Cols("id", "user_id", "title", "description", "status").
					Vals([]any{a.todo.ID, a.todo.UserID, a.todo.Title, a.todo.Description, a.todo.Status}).
					Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
			args:    args{ctx: context.TODO(), id: 1},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Delete("todos").
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), id: 1},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Delete("todos").
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
			args:    args{ctx: context.TODO(), id: 1},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name: "Success",
			args: args{ctx: context.TODO(), id: 1},
			want: &domain.Todo{
				ID:          1,
				UserID:      11,
				Title:       "title test",
				Description: "description test",
				Status:      domain.TodoStatusInProgress,
			},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				row := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "status"}).
					AddRow(1, 11, "title test", "description test", "IN_PROGRESS")

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).
					WillReturnRows(row)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
		in  map[string]string
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
			args:    args{ctx: context.TODO(), in: make(map[string]string)},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Select("id", "user_id", "title", "description", "status").
					From("todos").Prepared(true).ToSQL()

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name: "Success",
			args: args{ctx: context.TODO(), in: make(map[string]string)},
			want: []domain.Todo{
				{
					ID:          1,
					UserID:      12,
					Title:       "title test",
					Description: "description test",
					Status:      domain.TodoStatusDrop,
				}, {
					ID:          2,
					UserID:      13,
					Title:       "title test 2",
					Description: "description test 2",
					Status:      domain.TodoStatusInitiate,
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

				mock.ExpectQuery(query).WithArgs(testconvertArgs(args)...).
					WillReturnRows(rows)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
		sts domain.TodoStatus
	}
	tests := []struct {
		name    string
		mockFn  func(a args) (*SQLTodo, func() error)
		args    args
		wantErr error
	}{
		{
			name:    "Error",
			args:    args{ctx: context.TODO(), id: 1, sts: domain.TodoStatusDrop},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").
					Set(map[string]any{"status": a.sts}).
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), id: 1, sts: domain.TodoStatusDone},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").
					Set(map[string]any{"status": a.sts}).
					Where(goqu.Ex{"id": a.id}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
			args:    args{ctx: context.TODO(), todo: domain.Todo{}},
			wantErr: assert.AnError,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").Set(map[string]any{
					"user_id":     a.todo.UserID,
					"title":       a.todo.Title,
					"description": a.todo.Description,
					"status":      a.todo.Status,
				}).Where(goqu.Ex{"id": a.todo.ID}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnError(assert.AnError)

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), todo: domain.Todo{}},
			wantErr: nil,
			mockFn: func(a args) (*SQLTodo, func() error) {
				db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

				query, args, _ := goqu.Dialect(dbops.PostgresDriver).Update("todos").Set(map[string]any{
					"user_id":     a.todo.UserID,
					"title":       a.todo.Title,
					"description": a.todo.Description,
					"status":      a.todo.Status,
				}).Where(goqu.Ex{"id": a.todo.ID}).Prepared(true).ToSQL()

				mock.ExpectExec(query).WithArgs(testconvertArgs(args)...).
					WillReturnResult(sqlmock.NewResult(1, 1))

				return &SQLTodo{db: db, qu: goqu.Dialect(dbops.PostgresDriver)}, db.Close
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
